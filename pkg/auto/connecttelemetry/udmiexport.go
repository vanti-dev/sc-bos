package connecttelemetry

import (
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	"github.com/smart-core-os/sc-bos/pkg/auto/udmi"
	"github.com/smart-core-os/sc-bos/pkg/dbo"
	"github.com/smart-core-os/sc-bos/pkg/proto/metadatapb"
	"github.com/smart-core-os/sc-bos/pkg/proto/meterpb"
	"github.com/smart-core-os/sc-bos/pkg/proto/resourceusepb"
)

// UDMI event subfolders, the trailing topic segment Connect classifies on.
const (
	subfolderPointset  = "pointset"
	subfolderDiscovery = "discovery"
)

// eventTopic builds the UDMI event topic for a device under the Connect ingest
// grammar: <prefix>/devices/<deviceRef>/events/<subfolder>.
//
// NB deviceRef is currently the raw Smart Core device name, which contains '/'.
// That collides with topic segmentation and is a known, deliberately-parked gap
// (see docs/connect-telemetry-ingest.md) — do not rely on the segment count here.
func eventTopic(prefix, deviceRef, subfolder string) string {
	return strings.TrimRight(prefix, "/") + "/devices/" + deviceRef + "/events/" + subfolder
}

func pointsetTopic(prefix, deviceRef string) string {
	return eventTopic(prefix, deviceRef, subfolderPointset)
}

func discoveryTopic(prefix, deviceRef string) string {
	return eventTopic(prefix, deviceRef, subfolderDiscovery)
}

// meterTelemetry maps a meter reading to a UDMI pointset event. The point key is
// chosen by naming: DBO standard field names (energy_accumulator, …) make the
// building-config translation an identity mapping; raw names (usage/produced) leave
// the name→field mapping to the consumer.
//
// The reading instant is MeterReading.end_time, falling back to now when unset.
func meterTelemetry(r *meterpb.MeterReading, support *meterpb.MeterReadingSupport, now time.Time, naming dbo.Naming) udmi.PointsetEvent {
	ts := now
	if r.GetEndTime() != nil {
		ts = r.GetEndTime().AsTime()
	}
	points := udmi.PointsEvent{}
	// Meter readings are float32 (the Meter trait's type). For a cumulative
	// accumulator that is a resolution ceiling once the running total grows large
	// (~7 significant digits; reached sooner for meters reporting in Wh) — a known,
	// upstream-bound limitation tracked in docs/connect-telemetry-ingest.md.
	// Deliberately kept as float32 rather than widened to float64: encoding/json
	// emits the shortest decimal that round-trips the float32, whereas a float64
	// conversion would surface the float32's binary artifacts on the wire.
	for _, f := range dbo.MeterFields(support.GetUsageUnit(), support.GetProducedUnit()) {
		value := r.GetUsage()
		if f.Exported {
			value = r.GetProduced()
		}
		points[f.PointName(naming)] = udmi.PointValue{PresentValue: value}
	}
	return udmi.PointsetEvent{Timestamp: ts, Version: udmi.PointsetVersion, Points: points}
}

// meterInventory contributes the meter's declared points to discovery, keyed by the
// same name mode as the telemetry (see meterTelemetry). Meters are read-only, so no
// point is marked writable. Units carry the raw device unit string.
func meterInventory(support *meterpb.MeterReadingSupport, naming dbo.Naming) map[string]udmi.MetadataPoint {
	inv := map[string]udmi.MetadataPoint{}
	for _, f := range dbo.MeterFields(support.GetUsageUnit(), support.GetProducedUnit()) {
		inv[f.PointName(naming)] = udmi.MetadataPoint{Units: f.RawUnit}
	}
	return inv
}

// Units declared in discovery for ResourceUse points. ResourceUse leaves them
// implicit: utilisation is a 0-100 percentage and usage/limit are byte counts.
const (
	unitPercent = "percent"
	unitBytes   = "bytes"
)

// resourceUseTelemetry maps a node's ResourceUse to UDMI points. The names are
// DBO-style fields from Connect's catalogue, which extends DBO with the `cpu`,
// `memory` and `disk` components and the `datasize` measurement (SCC-978); they
// are the same in both naming modes. Unset optional fields produce no point.
// Disks are keyed by mount point (see diskInstance), so their list order doesn't
// matter.
//
// ResourceUse carries no timestamp; the caller stamps the reading with now.
func resourceUseTelemetry(r *resourceusepb.ResourceUse) udmi.PointsEvent {
	points := udmi.PointsEvent{}
	forEachResourceUsePoint(r, func(name string, value any, _ string, _ string) {
		points[name] = udmi.PointValue{PresentValue: value}
	})
	return points
}

// resourceUseInventory declares the points resourceUseTelemetry emits for r, with
// their units. Disk points carry the raw mount point as the model-point ref.
func resourceUseInventory(r *resourceusepb.ResourceUse) map[string]udmi.MetadataPoint {
	inv := map[string]udmi.MetadataPoint{}
	forEachResourceUsePoint(r, func(name string, _ any, units, ref string) {
		inv[name] = udmi.MetadataPoint{Units: units, Ref: ref}
	})
	return inv
}

// forEachResourceUsePoint calls fn for every set ResourceUse field we export, so
// telemetry and discovery can't drift apart. ref is the disk mount point, or empty.
func forEachResourceUsePoint(r *resourceusepb.ResourceUse, fn func(name string, value any, units, ref string)) {
	if cpu := r.GetCpu(); cpu != nil && cpu.Utilization != nil {
		fn("cpu_percentage_sensor", cpu.GetUtilization(), unitPercent, "")
	}
	if mem := r.GetMemory(); mem != nil {
		if mem.Utilization != nil {
			fn("memory_percentage_sensor", mem.GetUtilization(), unitPercent, "")
		}
		if mem.Usage != nil {
			fn("memory_datasize_sensor", mem.GetUsage(), unitBytes, "")
		}
		if mem.Limit != nil {
			fn("memory_datasize_capacity", mem.GetLimit(), unitBytes, "")
		}
	}
	for _, d := range r.GetDisks() {
		mount := d.GetMountPoint()
		n := diskInstance(mount)
		if d.Utilization != nil {
			fn(fmt.Sprintf("disk_percentage_sensor_%d", n), d.GetUtilization(), unitPercent, mount)
		}
		if d.Usage != nil {
			fn(fmt.Sprintf("disk_datasize_sensor_%d", n), d.GetUsage(), unitBytes, mount)
		}
		if d.Limit != nil {
			fn(fmt.Sprintf("disk_datasize_capacity_%d", n), d.GetLimit(), unitBytes, mount)
		}
	}
}

// diskInstance is a disk's DBO numeric instance suffix: the fnv32a hash of its raw
// mount point, so "/" is 705468254. A field name's tokens must all be catalogue
// subfields, so a mount can't be spelt into the name, and a list index would move
// one disk's history onto another when a mount appears or goes. The hash keeps
// each mount on its own series in any order; discovery's ref names the mount.
func diskInstance(mount string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(mount))
	return h.Sum32()
}

// discoveryEvent builds the UDMI device-metadata (discovery) message from the
// device's Smart Core metadata and the declared point inventory.
//
// system.name prefers the appearance title (falling back to the device name),
// tags carry the building subsystem selector, and location carries the floor.
// site is intentionally left empty here: on the Connect path org/site identity is
// supplied by broker enrichment, not by the payload.
func discoveryEvent(meta *metadatapb.Metadata, inventory map[string]udmi.MetadataPoint, now time.Time) udmi.MetadataEvent {
	ev := udmi.MetadataEvent{
		Timestamp: now,
		Version:   udmi.PointsetVersion,
	}
	if meta != nil {
		if a := meta.GetAppearance(); a != nil {
			ev.System.Name = a.GetTitle()
			ev.System.Description = a.GetDescription()
		}
		if ev.System.Name == "" {
			ev.System.Name = meta.GetName()
		}
		if sub := meta.GetMembership().GetSubsystem(); sub != "" {
			ev.System.Tags = append(ev.System.Tags, sub)
		}
		if loc := meta.GetLocation(); loc != nil && loc.GetFloor() != "" {
			ev.System.Location = &udmi.MetadataLocation{Floor: loc.GetFloor()}
		}
	}
	if len(inventory) > 0 {
		ev.Pointset = &udmi.MetadataPointset{Points: inventory}
	}
	return ev
}

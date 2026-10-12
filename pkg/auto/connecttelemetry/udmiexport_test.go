package connecttelemetry

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-bos/pkg/auto/udmi"
	"github.com/smart-core-os/sc-bos/pkg/dbo"
	"github.com/smart-core-os/sc-bos/pkg/proto/metadatapb"
	"github.com/smart-core-os/sc-bos/pkg/proto/meterpb"
	"github.com/smart-core-os/sc-bos/pkg/proto/resourceusepb"
)

func TestEventTopics(t *testing.T) {
	const ref = "van/uk/brum/ugs/meters/elec-main"

	assert.Equal(t,
		"tlm/devices/van/uk/brum/ugs/meters/elec-main/events/pointset",
		pointsetTopic("tlm", ref))
	assert.Equal(t,
		"tlm/devices/van/uk/brum/ugs/meters/elec-main/events/discovery",
		discoveryTopic("tlm", ref))

	// a trailing slash on the prefix must not double up.
	assert.Equal(t,
		"tlm/devices/m1/events/pointset",
		pointsetTopic("tlm/", "m1"))
}

func TestMeterTelemetry(t *testing.T) {
	end := time.Date(2026, 6, 22, 10, 15, 0, 0, time.UTC)
	now := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("usage and produced map to DBO energy fields", func(t *testing.T) {
		r := &meterpb.MeterReading{Usage: 123.45, Produced: 67.89, EndTime: timestamppb.New(end)}
		support := &meterpb.MeterReadingSupport{UsageUnit: "kWh", ProducedUnit: "kWh"}

		ev := meterTelemetry(r, support, now, dbo.NamingDBO)

		assert.Equal(t, end, ev.Timestamp, "timestamp should be the reading end_time")
		assert.Equal(t, udmi.PointsetVersion, ev.Version)
		require.Contains(t, ev.Points, dbo.FieldEnergyAccumulator)
		require.Contains(t, ev.Points, dbo.FieldExportedEnergyAccumulator)
		assert.Equal(t, float32(123.45), ev.Points[dbo.FieldEnergyAccumulator].PresentValue)
		assert.Equal(t, float32(67.89), ev.Points[dbo.FieldExportedEnergyAccumulator].PresentValue)
	})

	t.Run("falls back to now when end_time unset", func(t *testing.T) {
		r := &meterpb.MeterReading{Usage: 1}
		ev := meterTelemetry(r, &meterpb.MeterReadingSupport{UsageUnit: "kWh"}, now, dbo.NamingDBO)
		assert.Equal(t, now, ev.Timestamp)
	})

	t.Run("consumption-only meter omits exported when support has no producedUnit", func(t *testing.T) {
		r := &meterpb.MeterReading{Usage: 10, Produced: 0, EndTime: timestamppb.New(end)}
		support := &meterpb.MeterReadingSupport{UsageUnit: "kWh"} // no ProducedUnit

		ev := meterTelemetry(r, support, now, dbo.NamingDBO)

		require.Contains(t, ev.Points, dbo.FieldEnergyAccumulator)
		assert.NotContains(t, ev.Points, dbo.FieldExportedEnergyAccumulator, "no producedUnit ⇒ no exported series")
	})

	t.Run("without support, only the energy accumulator is emitted", func(t *testing.T) {
		ev := meterTelemetry(&meterpb.MeterReading{Usage: 10, Produced: 5}, nil, now, dbo.NamingDBO)
		require.Contains(t, ev.Points, dbo.FieldEnergyAccumulator)
		assert.NotContains(t, ev.Points, dbo.FieldExportedEnergyAccumulator,
			"without support (no declared producedUnit) exported energy is not claimed")
	})

	t.Run("water meter emits water_volume_accumulator from usage", func(t *testing.T) {
		r := &meterpb.MeterReading{Usage: 42, EndTime: timestamppb.New(end)}
		ev := meterTelemetry(r, &meterpb.MeterReadingSupport{UsageUnit: "m³"}, now, dbo.NamingDBO)
		require.Contains(t, ev.Points, dbo.FieldWaterVolumeAccumulator)
		assert.NotContains(t, ev.Points, dbo.FieldEnergyAccumulator, "a water meter is not energy")
		assert.Equal(t, float32(42), ev.Points[dbo.FieldWaterVolumeAccumulator].PresentValue)
	})

	t.Run("raw naming emits raw Smart Core point names", func(t *testing.T) {
		r := &meterpb.MeterReading{Usage: 123.45, Produced: 67.89, EndTime: timestamppb.New(end)}
		support := &meterpb.MeterReadingSupport{UsageUnit: "kWh", ProducedUnit: "kWh"}

		ev := meterTelemetry(r, support, now, dbo.NamingRaw)

		require.Contains(t, ev.Points, dbo.RawPointUsage)
		require.Contains(t, ev.Points, dbo.RawPointProduced)
		assert.NotContains(t, ev.Points, dbo.FieldEnergyAccumulator, "raw mode must not emit DBO names")
		assert.Equal(t, float32(123.45), ev.Points[dbo.RawPointUsage].PresentValue)
		assert.Equal(t, float32(67.89), ev.Points[dbo.RawPointProduced].PresentValue)
	})
}

func TestMeterInventory(t *testing.T) {
	t.Run("usage and produced units (raw strings)", func(t *testing.T) {
		inv := meterInventory(&meterpb.MeterReadingSupport{UsageUnit: "kWh", ProducedUnit: "kWh"}, dbo.NamingDBO)
		require.Contains(t, inv, dbo.FieldEnergyAccumulator)
		require.Contains(t, inv, dbo.FieldExportedEnergyAccumulator)
		// discovery carries the raw device unit string; DBO unit-name mapping is a
		// building-config concern.
		assert.Equal(t, "kWh", inv[dbo.FieldEnergyAccumulator].Units)
		assert.Equal(t, "kWh", inv[dbo.FieldExportedEnergyAccumulator].Units)
		// meters are read-only: nothing writable.
		assert.False(t, inv[dbo.FieldEnergyAccumulator].Writable)
		assert.False(t, inv[dbo.FieldExportedEnergyAccumulator].Writable)
	})

	t.Run("energy accumulator only when no producedUnit", func(t *testing.T) {
		inv := meterInventory(&meterpb.MeterReadingSupport{UsageUnit: "kWh"}, dbo.NamingDBO)
		require.Contains(t, inv, dbo.FieldEnergyAccumulator)
		assert.NotContains(t, inv, dbo.FieldExportedEnergyAccumulator)
	})

	t.Run("nil support still lists the energy accumulator", func(t *testing.T) {
		inv := meterInventory(nil, dbo.NamingDBO)
		require.Contains(t, inv, dbo.FieldEnergyAccumulator)
		assert.Empty(t, inv[dbo.FieldEnergyAccumulator].Units)
	})

	t.Run("water meter lists water_volume_accumulator with raw m³", func(t *testing.T) {
		inv := meterInventory(&meterpb.MeterReadingSupport{UsageUnit: "m³"}, dbo.NamingDBO)
		require.Contains(t, inv, dbo.FieldWaterVolumeAccumulator)
		assert.NotContains(t, inv, dbo.FieldEnergyAccumulator)
		assert.Equal(t, "m³", inv[dbo.FieldWaterVolumeAccumulator].Units)
	})

	t.Run("raw naming lists raw point names with native units", func(t *testing.T) {
		inv := meterInventory(&meterpb.MeterReadingSupport{UsageUnit: "kWh", ProducedUnit: "kWh"}, dbo.NamingRaw)
		require.Contains(t, inv, dbo.RawPointUsage)
		require.Contains(t, inv, dbo.RawPointProduced)
		assert.NotContains(t, inv, dbo.FieldEnergyAccumulator)
		assert.Equal(t, "kWh", inv[dbo.RawPointUsage].Units)
	})
}

func TestDiscoveryEvent(t *testing.T) {
	now := time.Date(2026, 6, 22, 10, 15, 0, 0, time.UTC)
	inv := map[string]udmi.MetadataPoint{
		dbo.FieldEnergyAccumulator: {Units: "kWh"},
	}

	meta := &metadatapb.Metadata{
		Name: "van/uk/brum/ugs/meters/elec-main",
		Appearance: &metadatapb.Metadata_Appearance{
			Title:       "Main Electrical Meter",
			Description: "Building main power meter",
		},
		Location:   &metadatapb.Metadata_Location{Floor: "03", Zone: "electrical-room"},
		Membership: &metadatapb.Metadata_Membership{Subsystem: "hvac"},
	}

	ev := discoveryEvent(meta, inv, now)

	assert.Equal(t, now, ev.Timestamp)
	assert.Equal(t, udmi.PointsetVersion, ev.Version)
	assert.Equal(t, "Main Electrical Meter", ev.System.Name)
	assert.Equal(t, "Building main power meter", ev.System.Description)
	assert.Equal(t, []string{"hvac"}, ev.System.Tags)
	require.NotNil(t, ev.System.Location)
	assert.Equal(t, "03", ev.System.Location.Floor)

	// site is left empty (Connect enrichment supplies it) and must be omitted, not
	// emitted as an invalid "site":"".
	b, err := json.Marshal(ev)
	require.NoError(t, err)
	assert.NotContains(t, string(b), `"site"`, "empty site must be omitted from discovery")
	require.NotNil(t, ev.Pointset)
	require.Contains(t, ev.Pointset.Points, dbo.FieldEnergyAccumulator)
	assert.Equal(t, "kWh", ev.Pointset.Points[dbo.FieldEnergyAccumulator].Units)
}

func TestDiscoveryEventFallsBackToDeviceName(t *testing.T) {
	now := time.Date(2026, 6, 22, 10, 15, 0, 0, time.UTC)

	// no appearance, no location, no membership.
	ev := discoveryEvent(&metadatapb.Metadata{Name: "meter-1"}, nil, now)

	assert.Equal(t, "meter-1", ev.System.Name)
	assert.Empty(t, ev.System.Description)
	assert.Empty(t, ev.System.Tags)
	assert.Nil(t, ev.System.Location)
	assert.Nil(t, ev.Pointset, "no inventory ⇒ no pointset block")
}

// Disk instances are the fnv32a hash of the mount point, pinned here (and in
// Connect's ingest-transform spec) so a change to the hashing can't pass unseen.
const (
	instRoot   = "705468254"  // "/"
	instVarLib = "3086875876" // "/var/lib/bos"
	instC      = "1065500364" // `C:\`
)

func TestResourceUseTelemetry(t *testing.T) {
	t.Run("full reading maps to utilisation and byte points", func(t *testing.T) {
		r := &resourceusepb.ResourceUse{
			Cpu:    &resourceusepb.CpuUse{Utilization: new(float32(12.5)), CorePercent: []float32{10, 15}},
			Memory: &resourceusepb.MemoryUse{Usage: new(uint64(4 << 30)), Limit: new(uint64(16 << 30)), Utilization: new(float32(25))},
			Disks: []*resourceusepb.DiskUse{
				{MountPoint: "/", Usage: new(uint64(100)), FreeBytes: new(uint64(300)), Limit: new(uint64(400)), Utilization: new(float32(25))},
				{MountPoint: "/var/lib/bos", Usage: new(uint64(50)), Limit: new(uint64(200)), Utilization: new(float32(25))},
			},
			Network: &resourceusepb.NetworkUse{ConnectionCount: new(uint64(7))},
		}

		assert.Equal(t, udmi.PointsEvent{
			"cpu_percentage_sensor":                {PresentValue: float32(12.5)},
			"memory_percentage_sensor":             {PresentValue: float32(25)},
			"memory_datasize_sensor":               {PresentValue: uint64(4 << 30)},
			"memory_datasize_capacity":             {PresentValue: uint64(16 << 30)},
			"disk_percentage_sensor_" + instRoot:   {PresentValue: float32(25)},
			"disk_datasize_sensor_" + instRoot:     {PresentValue: uint64(100)},
			"disk_datasize_capacity_" + instRoot:   {PresentValue: uint64(400)},
			"disk_percentage_sensor_" + instVarLib: {PresentValue: float32(25)},
			"disk_datasize_sensor_" + instVarLib:   {PresentValue: uint64(50)},
			"disk_datasize_capacity_" + instVarLib: {PresentValue: uint64(200)},
		}, resourceUseTelemetry(r), "core_percent, free_bytes and connection_count are not exported")
	})

	t.Run("unset optionals are omitted", func(t *testing.T) {
		r := &resourceusepb.ResourceUse{
			Cpu:    &resourceusepb.CpuUse{},
			Memory: &resourceusepb.MemoryUse{Usage: new(uint64(1))},
			Disks:  []*resourceusepb.DiskUse{{MountPoint: "/", Utilization: new(float32(0))}},
		}
		assert.Equal(t, udmi.PointsEvent{
			"memory_datasize_sensor":             {PresentValue: uint64(1)},
			"disk_percentage_sensor_" + instRoot: {PresentValue: float32(0)},
		}, resourceUseTelemetry(r), "a set zero is kept, an unset field is not")
	})

	t.Run("nil reading has no points", func(t *testing.T) {
		assert.Empty(t, resourceUseTelemetry(nil))
	})
}

func TestResourceUseInventory(t *testing.T) {
	r := &resourceusepb.ResourceUse{
		Cpu:    &resourceusepb.CpuUse{Utilization: new(float32(1))},
		Memory: &resourceusepb.MemoryUse{Usage: new(uint64(1)), Limit: new(uint64(2)), Utilization: new(float32(50))},
		Disks:  []*resourceusepb.DiskUse{{MountPoint: `C:\`, Usage: new(uint64(1)), Limit: new(uint64(2)), Utilization: new(float32(50))}},
	}

	assert.Equal(t, map[string]udmi.MetadataPoint{
		"cpu_percentage_sensor":           {Units: "percent"},
		"memory_percentage_sensor":        {Units: "percent"},
		"memory_datasize_sensor":          {Units: "bytes"},
		"memory_datasize_capacity":        {Units: "bytes"},
		"disk_percentage_sensor_" + instC: {Units: "percent", Ref: `C:\`},
		"disk_datasize_sensor_" + instC:   {Units: "bytes", Ref: `C:\`},
		"disk_datasize_capacity_" + instC: {Units: "bytes", Ref: `C:\`},
	}, resourceUseInventory(r))

	assert.Empty(t, resourceUseInventory(nil), "no reading yet ⇒ no inventory")
}

func TestDiskInstance(t *testing.T) {
	assert.Equal(t, instRoot, strconv.FormatUint(uint64(diskInstance("/")), 10))
	assert.Equal(t, instVarLib, strconv.FormatUint(uint64(diskInstance("/var/lib/bos")), 10))
	assert.Equal(t, instC, strconv.FormatUint(uint64(diskInstance(`C:\`)), 10))
	assert.NotEqual(t, diskInstance("/var-lib"), diskInstance("/var/lib"),
		"mounts differing only in punctuation stay distinct")

	// The same mount keeps its points whatever order the disks are listed in.
	a := &resourceusepb.DiskUse{MountPoint: "/var-lib", Utilization: new(float32(1))}
	b := &resourceusepb.DiskUse{MountPoint: "/var/lib", Utilization: new(float32(2))}
	forward := resourceUseTelemetry(&resourceusepb.ResourceUse{Disks: []*resourceusepb.DiskUse{a, b}})
	reverse := resourceUseTelemetry(&resourceusepb.ResourceUse{Disks: []*resourceusepb.DiskUse{b, a}})
	assert.Equal(t, forward, reverse, "names must not depend on disk order")
	assert.Len(t, forward, 2)
}

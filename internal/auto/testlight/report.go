package testlight

import (
	"encoding/csv"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/timshannon/bolthold"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
)

type ReportEntry struct {
	Name                   string
	LastUpdate             time.Time
	Faults                 []gen.EmergencyLightFault
	LatestFunctionTestPass time.Time
	LatestDurationTestPass time.Time
}

func GenerateReport(db *bolthold.Store) ([]ReportEntry, error) {
	data := make(map[string]*ReportEntry)

	// get the latest status for all lights we know about
	err := db.ForEach(nil, func(record *LatestStatusRecord) error {
		data[record.Name] = &ReportEntry{
			Name:       record.Name,
			LastUpdate: record.LastUpdate,
			Faults:     record.Faults,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// find the latest function and duration test for each light
	query := bolthold.Where("Kind").In(FunctionTestPassEvent, DurationTestPassEvent)
	aggregates, err := db.FindAggregate(&EventRecord{}, query, "Name", "Kind")
	if err != nil {
		return nil, err
	}
	for _, aggregate := range aggregates {
		var (
			name   string
			kind   EventKind
			latest EventRecord
		)
		aggregate.Group(&name, &kind)
		aggregate.Max("Timestamp", &latest)

		entry, ok := data[name]
		if !ok {
			// we don't know about this device
			continue
		}

		switch kind {
		case FunctionTestPassEvent:
			entry.LatestFunctionTestPass = latest.Timestamp
		case DurationTestPassEvent:
			entry.LatestDurationTestPass = latest.Timestamp
		default:
			// other values were excluded by the query, this can't happen unless there is a bug
			panic("unexpected EventRecord.Kind")
		}
	}

	// convert into a slice sorted by name
	result := make([]ReportEntry, 0, len(data))
	for _, entry := range data {
		result = append(result, *entry)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result, nil
}

func WriteReportCSV(dst io.Writer, report []ReportEntry, header bool) (err error) {
	writer := csv.NewWriter(dst)
	if header {
		err = writer.Write([]string{"Name", "Last Update", "Faults", "Latest Function Test Pass", "Latest Duration Test Pass"})
		if err != nil {
			return err
		}
	}

	for _, entry := range report {
		var faultStrings []string
		for _, fault := range entry.Faults {
			faultStrings = append(faultStrings, fault.String())
		}

		var latestFunctionString, latestDurationString string
		if !entry.LatestFunctionTestPass.IsZero() {
			latestFunctionString = entry.LatestFunctionTestPass.Format("2006-01-02")
		}
		if !entry.LatestDurationTestPass.IsZero() {
			latestDurationString = entry.LatestDurationTestPass.Format("2006-01-02")
		}

		line := []string{
			entry.Name,
			entry.LastUpdate.Format(time.RFC3339),
			strings.Join(faultStrings, " "),
			latestFunctionString,
			latestDurationString,
		}

		err = writer.Write(line)
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

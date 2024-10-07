package air_quality

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
)

func Seed(ctx context.Context, db *pgxpool.Pool, name string, lookBack time.Duration) error {
	current := time.Now()
	now := current

	source := fmt.Sprintf("%s[%s]", name, trait.AirQualitySensor)

	for current.After(now.Add(-lookBack)) {

		co2 := rand.Float32() * 2000
		voc := rand.Float32() * 500
		airPressure := rand.Float32()
		comfort := traits.AirQuality_Comfort(rand.Intn(2) + 1)
		infection := rand.Float32() * 100
		score := rand.Float32() * 100
		particulate1 := rand.Float32()
		particulate25 := rand.Float32()
		particulate10 := rand.Float32()
		airChange := rand.Float32() * 5

		payload, err := proto.Marshal(&traits.AirQuality{
			CarbonDioxideLevel:       &co2,
			VolatileOrganicCompounds: &voc,
			AirPressure:              &airPressure,
			Comfort:                  comfort,
			InfectionRisk:            &infection,
			Score:                    &score,
			ParticulateMatter_1:      &particulate1,
			ParticulateMatter_10:     &particulate10,
			ParticulateMatter_25:     &particulate25,
			AirChangePerHour:         &airChange,
		})

		if err != nil {
			return err
		}
		cmd, err := db.Exec(ctx, `INSERT INTO history (source, create_time, payload) VALUES ($1, $2, $3)`, source, current.Format(time.RFC3339Nano), payload)

		if err != nil {
			return err
		}

		if cmd.RowsAffected() == 0 {
			return pgx.ErrNoRows
		}

		current = current.Add(-time.Duration(rand.Intn(60)) * time.Minute)
	}
	return nil
}

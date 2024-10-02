package occupancy

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
)

func Seed(ctx context.Context, db *pgxpool.Pool, name string, lookBack time.Duration) error {
	current := time.Now()
	now := current

	source := fmt.Sprintf("%s[%s]", name, trait.OccupancySensor)

	for current.After(now.Add(-lookBack)) {

		payload, err := proto.Marshal(&traits.Occupancy{
			PeopleCount:     int32(rand.Intn(50)),
			StateChangeTime: timestamppb.New(current),
			Confidence:      1,
		})

		if err != nil {
			return err
		}

		cmd, err := db.Exec(ctx, "INSERT INTO history (source, create_time, payload) VALUES ($1, $2, $3)", source, current.Format(time.RFC3339Nano), payload)

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

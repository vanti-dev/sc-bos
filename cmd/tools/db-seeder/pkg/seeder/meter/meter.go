package meter

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/proto"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/meter"
)

func Seed(ctx context.Context, db *pgxpool.Pool, name string, lookBack time.Duration) error {
	now := time.Now()
	current := now.Add(-lookBack)

	source := fmt.Sprintf("%s[%s]", name, meter.TraitName)

	incremental := rand.Float32() * 1_000

	for current.Before(now) {
		incremental = incremental + rand.Float32()*1_000
		payload, err := proto.Marshal(&gen.MeterReading{
			Usage: incremental,
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

		current = current.Add(time.Duration(rand.Intn(60)) * time.Minute)

	}
	return nil
}

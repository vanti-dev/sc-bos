package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/pkg/history/pgxstore"
)

func SeedOccupancy(ctx context.Context, db *pgxpool.Pool, name string, lookBack time.Duration) error {
	current := time.Now()
	now := current

	source := fmt.Sprintf("%s[%s]", name, trait.OccupancySensor)

	store, err := pgxstore.SetupStoreFromPool(ctx, source, db)
	if err != nil {
		return err
	}

	for current.After(now.Add(-lookBack)) {

		payload, err := proto.Marshal(&traits.Occupancy{
			PeopleCount:     int32(rand.Intn(50)),
			StateChangeTime: timestamppb.New(current),
			Confidence:      1,
		})

		if err != nil {
			return err
		}

		_, _, err = store.Insert(ctx, current, payload)

		if err != nil {
			return err
		}

		current = current.Add(-time.Duration(rand.Intn(60)) * time.Minute)
	}

	return nil
}

package auto

import (
	"math/rand"
	"time"
)

func float32Between(min, max float32) float32 {
	return min + (max-min)*rand.Float32()
}

func oneOf[T any](vals ...T) T {
	return vals[rand.Intn(len(vals))]
}

func ptr[T any](v T) *T {
	return &v
}

func durationBetween(min, max time.Duration) time.Duration {
	return time.Duration(rand.Intn(int(max-min)) + int(min))
}

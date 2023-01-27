package services

import (
	"golang.org/x/exp/slices"
)

type Option func(a *Api)

func WithKnownTypes(t ...string) Option {
	return func(a *Api) {
		a.knownTypes = append(a.knownTypes, t...)
	}
}

func WithKnownTypesFromMapKeys[M ~map[string]T, T any](m M) Option {
	return func(a *Api) {
		a.knownTypes = slices.Grow(a.knownTypes, len(m))
		for k, _ := range m {
			a.knownTypes = append(a.knownTypes, k)
		}
	}
}

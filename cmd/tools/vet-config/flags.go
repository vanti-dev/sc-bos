package main

import (
	"strings"
)

type stringList []string

func (p *stringList) String() string {
	return strings.Join(*p, ",")
}

func (p *stringList) Set(value string) error {
	*p = append(*p, strings.Split(value, ",")...)
	return nil
}

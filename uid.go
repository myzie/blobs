package main

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

var entropy *rand.Rand

func init() {
	entropy = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func uid() string {
	return ulid.MustNew(ulid.Now(), entropy).String()
}

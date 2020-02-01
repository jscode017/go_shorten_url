package main

import (
	"math/rand"
	"time"
)

func RandonIntFromTime() int {
	source := rand.NewSource(time.Now().UnixNano() / 1e6)
	random := rand.New(source)
	randomInteger := random.Intn(1e5)
	return randomInteger
}

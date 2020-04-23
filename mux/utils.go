package mux

import (
	"math/rand"
	"time"
)

//SIDUint32 simple session id generator
func SIDUint32() uint32 {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	return seededRand.Uint32()
}

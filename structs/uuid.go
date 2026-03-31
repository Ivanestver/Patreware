package structs

import "math/rand/v2"

type UUID uint64

func GenerateUUID() UUID {
	return UUID(rand.Uint64())
}

func (id UUID) Equals(other UUID) bool {
	return id == other
}

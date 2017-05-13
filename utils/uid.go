package utils

import (
	"hash/fnv"
	"math/rand"
)

var SHORT_UID_LENGTH int = 7

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func ShortUid(toHash *string) string {
	return randomString(SHORT_UID_LENGTH)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

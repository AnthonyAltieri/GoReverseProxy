package utils

import (
	"time"
	"hash/fnv"
	"math/rand"
	"fmt"
)

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func ShortUid(toHash *string) string {
	rand.Seed(time.Now().Unix())
	hashed := rand.Uint32()
	if toHash != nil {
		hashed = hash(*toHash)
	}
	currentTime := time.Now().Unix()
	doubleHash := fmt.Sprintf("%d", hash(fmt.Sprintf("%d:%d", hashed, currentTime)))
	return doubleHash
}

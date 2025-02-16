package utils

import (
	"crypto/sha1"
	"encoding/binary"
	"hash/fnv"
	"os"
	"strconv"
)

func GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func GenerateUniqueHashUrl(url string) uint32 {
	hasher := fnv.New32a()
	hasher.Write([]byte(url))
	return hasher.Sum32()
}

func ChordHash(nodeAddress string, mBits int) uint64 {
	hash := sha1.Sum([]byte(nodeAddress))
	hashBytes := hash[:8]
	id := binary.BigEndian.Uint64(hashBytes)
	mask := uint64(1<<mBits - 1)
	return id & mask
}

func Between(x, a, b uint64) bool {
	return (a < x && x < b) || (b < x && x < a)
}

func BetweenRightInclusive(x, a, b uint64) bool {
	return (a < x && x <= b) || (b < x && x <= a)
}

package utils

import (
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

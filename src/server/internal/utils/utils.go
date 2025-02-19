package utils

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"os"
	"strconv"
	"strings"
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

func BetweenRightInclusive(x, a, b uint64) bool { // use this instead of Between
	fmt.Println("Asked if ", x, " is between right inclusive", a, " and ", b)
	if a < b {
		return a < x && x <= b
	}
	fmt.Println("WRAP AROUND")
	return a < x || x <= b
}

func Between(x, a, b uint64) bool {
	fmt.Println("Asked if ", x, " is between", a, " and ", b)
	if a < b {
		return a < x && x < b
	}
	fmt.Println("WRAP AROUND")
	return a < x || x < b
}

func IpAddress(addrWithPort string) string {
	return strings.Split(addrWithPort, ":")[0]
}

func ChangePort(ip string, port string) string {
	return IpAddress(ip) + ":" + port
}

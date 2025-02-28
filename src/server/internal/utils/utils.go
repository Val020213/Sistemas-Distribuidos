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

var hardCodeHashes = map[string]uint64{
	// "10.0.10.11": 2,
	// "10.0.10.12": 5,
	// "10.0.10.13": 7,
	// "10.0.10.14": 4,
	// "10.0.10.15": 3,
}

func ChordHash(nodeAddress string, mBits int) uint64 {
	if id, ok := hardCodeHashes[nodeAddress]; ok {
		return id
	}

	hash := sha1.Sum([]byte(nodeAddress))
	hashBytes := hash[:8]
	id := binary.BigEndian.Uint64(hashBytes)
	mask := uint64(1<<mBits - 1)
	return id & mask
}

func BetweenRightInclusive(x, a, b uint64) bool { // use this instead of Between
	if a < b {
		return a < x && x <= b
	}
	return a < x || x <= b
}

func Between(x, a, b uint64) bool {
	if a < b {
		return a < x && x < b
	}
	return a < x || x < b
}

func IpAddress(addrWithPort string) string {
	return strings.Split(addrWithPort, ":")[0]
}

func ChangePort(ip string, port string) string {
	return IpAddress(ip) + ":" + port
}

func RedPrint(a ...any) (n int, err error) {
	return fmt.Fprint(os.Stderr, "\033[31m", a, "\033[0m")
}

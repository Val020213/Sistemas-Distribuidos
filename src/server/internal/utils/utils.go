package utils

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"os"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
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

func GetFilterBetweenRightIncusive(a, b uint64) bson.M {
	if a < b {
		return bson.M{
			"$and": []bson.M{
				{"key": bson.M{"$gt": a}},
				{"key": bson.M{"$lte": b}},
			},
		}
	}
	return bson.M{
		"$or": []bson.M{
			{"key": bson.M{"$gt": a}},
			{"key": bson.M{"$lte": b}},
		},
	}
}

func GetNegativeFilterBetweenRightIncusive(a, b uint64) bson.M {
	return bson.M{
		"$not": GetFilterBetweenRightIncusive(a, b),
	}
}

// if n.Id != successor.Id && utils.Between(key, predecessorId, successor.Id) {
// 	replicated = append(replicated, ToPbData(&cData, key))
// }

func GetFilterBetween(a, b uint64) bson.M {
	if a < b {
		return bson.M{
			"$and": []bson.M{
				{"key": bson.M{"$gt": a}},
				{"key": bson.M{"$lt": b}},
			},
		}
	}
	return bson.M{
		"$or": []bson.M{
			{"key": bson.M{"$gt": a}},
			{"key": bson.M{"$lt": b}},
		},
	}
}

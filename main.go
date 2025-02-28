package main

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
)

func ChordHash(nodeAddress string, mBits int) uint64 {
	hash := sha1.Sum([]byte(nodeAddress))
	hashBytes := hash[:8]
	id := binary.BigEndian.Uint64(hashBytes)
	mask := uint64(1<<mBits - 1)
	return id & mask
}

func main() {
	value := "10.0.10.13"
	bits := 3

	fmt.Println(ChordHash(value, bits))
}

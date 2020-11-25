package main

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
)

func main() {
	filter := NewBloomFilter(20, []HashingFunction{HashFun1, HashFun2})
	filter.Add("hi")

	fmt.Printf("Is hi in the set: ")
	fmt.Println(filter.InSet("hi"))
	fmt.Printf("Is hey in the set: ")
	fmt.Println(filter.InSet("hey"))
}

func HashFun1(item string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(item))
	return h.Sum64()
}

func HashFun2(item string) uint64 {
	return binary.BigEndian.Uint64(fnv.New128().Sum([]byte(item)))
}

type HashingFunction func(in string) uint64

func NewBloomFilter(bitArrayLength int, hashingFunctions []HashingFunction) *BloomFilter {
	b := new(BloomFilter)
	b.BitArray = make([]bool, bitArrayLength)
	b.HashingFunctions = hashingFunctions
	return b
}

type BloomFilter struct {
	BitArray []bool
	HashingFunctions []HashingFunction
}

func (filter BloomFilter) Add(item string) {
	for _, hashFn := range filter.HashingFunctions {
		hasVal := hashFn(item)
		k := int(hasVal) % len(filter.BitArray)
		filter.BitArray[k] = true
	}
}

func (filter BloomFilter) InSet(item string) string {
	if filter.DefinitelyMissingFromSet(item) {
		return "No"
	}
	return "Maybe"
}

func (filter BloomFilter) DefinitelyMissingFromSet(item string) bool {
	for _, hashFn := range filter.HashingFunctions {
		hasVal := hashFn(item)
		k := int(hasVal) % len(filter.BitArray)
		if !filter.BitArray[k] {
			return true
		}
	}

	return false
}
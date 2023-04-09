package bloom

import (
	"fmt"
	"math"
)

type HashFunction interface {
	Hash(s string) uint32
}
type Filter struct {
	bitArray            []bool
	numHashFunctions    int
	expectedNumElements int
	hashFunction        HashFunction
}

func NewFilter(expectedNumElements int, bitArraySize int, hashFunction HashFunction) *Filter {
	bitArray := make([]bool, bitArraySize)
	numHashFunctions := int(math.Round(float64(bitArraySize) / float64(expectedNumElements) * math.Log(2)))
	return &Filter{
		bitArray:            bitArray,
		numHashFunctions:    numHashFunctions,
		expectedNumElements: expectedNumElements,
		hashFunction:        hashFunction,
	}
}

func (*Filter) HashFunc(s string) uint32 {
	hash := uint32(2166136261)
	for i := 0; i < len(s); i++ {
		hash = (hash * 16777619) ^ uint32(s[i])
	}
	return hash
}

func (filter *Filter) addToFilter(s string) {
	for i := 0; i < filter.numHashFunctions; i++ {
		hash := filter.hashFunction.Hash(fmt.Sprintf("%d%s", i, s))
		index := hash % uint32(len(filter.bitArray))
		filter.bitArray[index] = true
	}
}

func (filter *Filter) check(s string) bool {
	for i := 0; i < filter.numHashFunctions; i++ {
		if filter.bitArray[filter.hashFunction.Hash(fmt.Sprintf("%d%s", i, s))] == false {
			return false
		}
	}
	return true
}

func (filter *Filter) GetFalsePositiveProbability() float64 {
	k := float64(filter.numHashFunctions)
	n := float64(filter.expectedNumElements)
	m := float64(len(filter.bitArray))
	return math.Pow(1-math.Exp(-k*n/m), k)
}

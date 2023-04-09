package bloom

import (
	"errors"
	"fmt"
	"math"
)

type HashFunction interface {
	Hash(s string) uint32
}

type DefaultHash struct {
}

func (*DefaultHash) Hash(s string) uint32 {
	hash := uint32(2166136261)
	for i := 0; i < len(s); i++ {
		hash = (hash * 16777619) ^ uint32(s[i])
	}
	return hash
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

func NewFilterWithDefaultHash(expectedNumElements int, bitArraySize int) *Filter {
	return NewFilter(expectedNumElements, bitArraySize, &DefaultHash{})
}

func CalculateArraySize(expectedNumElements int, falsePositiveProbability float64) (int, error) {
	if falsePositiveProbability > 1 || falsePositiveProbability <= 0 {
		return 0, errors.New("falsePositiveProbability must be greater than 0 and smaller than 1")
	}
	return int(math.Ceil(-1 * (float64(expectedNumElements) * math.Log(falsePositiveProbability)) / math.Pow(math.Log(2), 2))), nil
}

func (filter *Filter) AddToFilter(s string) {
	if filter.hashFunction == nil {
		panic("hashFunction is nil")
	}
	for i := 0; i < filter.numHashFunctions; i++ {
		hash := filter.hashFunction.Hash(fmt.Sprintf("%d%s", i, s))
		index := hash % uint32(len(filter.bitArray))
		filter.bitArray[index] = true
	}
}

func (filter *Filter) Check(s string) bool {
	for i := 0; i < filter.numHashFunctions; i++ {
		if filter.bitArray[filter.hashFunction.Hash(fmt.Sprintf("%d%s", i, s))%uint32(len(filter.bitArray))] == false {
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

package bloom

import (
	"errors"
	"fmt"
	"log"
	"math"
	"task/pkg/bloom/abstraction"
)

type Filter struct {
	bitArray              []bool
	numOfHashFunctions    int
	expectedNumOfElements int
	hashFunction          abstraction.HashFunction
}

func NewFilter(expectedNumElements int, bitArraySize int, hashFunction abstraction.HashFunction) *Filter {
	bitArray := make([]bool, bitArraySize)
	numHashFunctions := int(float64(bitArraySize) / float64(expectedNumElements) * math.Log(2))
	return &Filter{
		bitArray:              bitArray,
		numOfHashFunctions:    numHashFunctions,
		expectedNumOfElements: expectedNumElements,
		hashFunction:          hashFunction,
	}
}

func NewFilterWithDefaultHash(expectedNumElements int, bitArraySize int) *Filter {
	return NewFilter(expectedNumElements, bitArraySize, &DefaultHash{})
}

func (f *Filter) AddToFilter(s string) {
	if f.hashFunction == nil {
		log.Println("hashFunction is nil")
		return
	}
	if f.Check(s) {
		log.Println("barcode is already in filter")
		return
	} else {
		for i := 0; i < f.numOfHashFunctions; i++ {
			hash := f.hashFunction.Hash(fmt.Sprintf("%d%s", i, s))
			index := hash % uint32(len(f.bitArray))
			f.bitArray[index] = true
		}
	}
}

func (f *Filter) Check(s string) bool {
	if f.hashFunction == nil {
		log.Println("hashFunction is nil")
		return false
	}

	for i := 0; i < f.numOfHashFunctions; i++ {
		hash := f.hashFunction.Hash(fmt.Sprintf("%d%s", i, s))
		index := hash % uint32(len(f.bitArray))
		if !f.bitArray[index] {
			return false
		}
	}

	return true
}

func CalculateArraySize(expectedNumElements int, falsePositiveProbability float64) (int, error) {
	if falsePositiveProbability > 1 || falsePositiveProbability <= 0 {
		return 0, errors.New("falsePositiveProbability must be greater than 0 and smaller than 1")
	}

	return int(math.Ceil(-1 * (float64(expectedNumElements) * math.Log(falsePositiveProbability)) / math.Pow(
		math.Log(2), 2))), nil
}

func (f *Filter) GetFalsePositiveProbability() float64 {
	k := f.numOfHashFunctions
	n := f.expectedNumOfElements
	m := len(f.bitArray)
	return math.Pow(1-math.Exp(float64(-k*n/m)), float64(k))
}

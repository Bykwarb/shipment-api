package bloom

import (
	"fmt"
	"testing"
)

func TestFilter(t *testing.T) {
	expectedNumElements := 100
	bitArraySize, err := CalculateArraySize(expectedNumElements, 0.01)
	if err != nil {
		t.Errorf("Error calculating bit array size: %v", err)
	}
	filter := NewFilterWithDefaultHash(expectedNumElements, bitArraySize)

	for i := 0; i < expectedNumElements; i++ {
		if filter.Check(fmt.Sprint(i)) {
			t.Errorf("Element %d is not expected to be in the filter", i)
		}
	}

	for i := 0; i < expectedNumElements; i++ {
		filter.AddToFilter(fmt.Sprint(i))
	}

	for i := 0; i < expectedNumElements; i++ {
		if !filter.Check(fmt.Sprint(i)) {
			t.Errorf("Element %d is expected to be in the filter", i)
		}
	}

	for i := expectedNumElements; i < expectedNumElements+100; i++ {
		if filter.Check(fmt.Sprint(i)) {
			t.Errorf("Element %d is not expected to be in the filter", i)
		}
	}
}

func TestCalculateArraySize(t *testing.T) {

	testCases := []struct {
		expectedNumElements      int
		falsePositiveProbability float64
		expectedBitArraySize     int
	}{
		{100, 0.01, 959},
		{1000, 0.001, 14378},
		{10000, 0.0001, 191702},
	}

	for _, test := range testCases {
		bitArraySize, err := CalculateArraySize(test.expectedNumElements, test.falsePositiveProbability)
		if err != nil {
			t.Errorf("Error calculating bit array size: %v", err)
		}
		if bitArraySize != test.expectedBitArraySize {
			t.Errorf("Bit array size calculation failed for %d elements and %f false positive probability. Expected %d, got %d", test.expectedNumElements, test.falsePositiveProbability, test.expectedBitArraySize, bitArraySize)
		}
	}
}

func TestFilter_GetFalsePositiveProbability(t *testing.T) {
	expectedNumElements := 100
	bitArraySize, err := CalculateArraySize(expectedNumElements, 0.01)
	if err != nil {
		t.Errorf("Error calculating bit array size: %v", err)
	}
	filter := NewFilterWithDefaultHash(expectedNumElements, bitArraySize)

	for i := 0; i < expectedNumElements; i++ {
		filter.AddToFilter(fmt.Sprint(i))
	}

	falsePositiveProbability := filter.GetFalsePositiveProbability()
	expectedFalsePositiveProbability := 0.01
	if falsePositiveProbability < expectedFalsePositiveProbability-0.01 || falsePositiveProbability > expectedFalsePositiveProbability+0.01 {
		t.Errorf("False positive probability calculation failed. Expected %f, got %f", expectedFalsePositiveProbability, falsePositiveProbability)
	}
}

func TestDefaultHash(t *testing.T) {
	testCases := []struct {
		input    string
		expected uint32
	}{
		{"hello", 3069866343},
		{"world", 2609808943},
		{"", 2166136261},
		{"test", 3157003241},
	}

	hashFunc := &DefaultHash{}

	for _, tc := range testCases {
		output := hashFunc.Hash(tc.input)
		if output != tc.expected {
			t.Errorf("DefaultHash(%q) = %d; want %d", tc.input, output, tc.expected)
		}
	}
}

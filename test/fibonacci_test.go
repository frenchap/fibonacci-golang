package test

import (
	"math/big"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestMemoryStoreExpectedValues(t *testing.T) {
	testMeta := NewTestMeta()

	threeHundred := big.NewInt(0)
	threeHundred, succeeded := threeHundred.SetString("222232244629420445529739893461909967206666939096499764990979600", 10)
	if !succeeded {
		t.Fatal("Error creating big.int for 300th fibonacci element")
	}
	expectedValues := map[int]big.Int{
		0:   *big.NewInt(0),
		1:   *big.NewInt(1),
		2:   *big.NewInt(1),
		3:   *big.NewInt(2),
		4:   *big.NewInt(3),
		5:   *big.NewInt(5),
		6:   *big.NewInt(8),
		13:  *big.NewInt(233),
		300: *threeHundred,
	}

	for k, expectedValue := range expectedValues {
		actualValue, err := testMeta.TestMemoryStore.GetValue(k)

		if err != nil {
			t.Fatal("Error retrieving value", err)
		}
		if actualValue.Y.Cmp(&expectedValue) != 0 {
			t.Fatalf("Value incorrect: should be %s but was %s", expectedValue.String(), actualValue.Y.String())
		}
	}
}

func TestMemoryStoreGetIntermediateValueCount(t *testing.T) {
	testMeta := NewTestMeta()

	intermediateCountTestMap := map[*big.Int]int{big.NewInt(377): 13, big.NewInt(378): 14}

	for value, expectedCount := range intermediateCountTestMap {
		actualCount := testMeta.TestMemoryStore.GetIntermediateValueCount(value)
		if actualCount != expectedCount {
			t.Fatalf("GetIntermediateValueCount incorrect: should be %d but was %d", expectedCount, actualCount)
		}
	}
}

func TestMemoryStoreUpperBound(t *testing.T) {
	testMeta := NewTestMeta()

	expectedMax := 1
	actualMax := testMeta.TestMemoryStore.GetMax()
	if actualMax != expectedMax {
		t.Fatalf("Max stored value incorrect: should be %d but was %d", expectedMax, actualMax)
	}

	maxFibonacci, err := testMeta.TestMemoryStore.GetValue(testMeta.TestUpperBound)
	if err != nil {
		t.Fatal("Error getting max fibonacci: ", err)
	}

	expectedMax = testMeta.TestUpperBound
	actualMax = testMeta.TestMemoryStore.GetMax()
	if actualMax != expectedMax {
		t.Fatalf("Max stored value incorrect: should be %d but was %d", expectedMax, actualMax)
	}

	logrus.Infof("Max fibonacci: %d, Time cost(ms): %d", maxFibonacci.Y, int64(maxFibonacci.TimeCost/time.Millisecond))
}

func TestFibonacciMemoryStoreValuePastUpperBound(t *testing.T) {
	testMeta := NewTestMeta()

	_, err := testMeta.TestMemoryStore.GetValue(testMeta.TestUpperBound + 1)
	if err == nil {
		t.Fatal("Value over upper bound should throw error")
	}

}

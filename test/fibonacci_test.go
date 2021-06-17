package test

import (
	"fmt"
	"testing"
	"time"
)

type FibonacciElement struct {
	Y        int
	TimeCost time.Duration
}

type IFibonacciStore interface {
	GetMax() int
	GetValue(x int) (FibonacciElement, error)
}

type MemoryFibonacciStore struct {
	UpperBound int
	Sequence   map[int]FibonacciElement
}

func NewMemoryFibonacciStore(upperBound int) MemoryFibonacciStore {
	return MemoryFibonacciStore{
		UpperBound: upperBound,
		Sequence: map[int]FibonacciElement{
			0: {0, 0},
			1: {1, 0},
		},
	}
}

func (thisStore MemoryFibonacciStore) GetMax() int {
	return len(thisStore.Sequence) - 1
}

func (thisStore MemoryFibonacciStore) GetValue(x int) (FibonacciElement, error) {
	if value, valueFound := thisStore.Sequence[x]; valueFound {
		return value, nil
	}
	if x > thisStore.UpperBound {
		return FibonacciElement{0, 0}, fmt.Errorf("requested value of %d above upper bound of %d", x, thisStore.UpperBound)
	}

	start := time.Now()
	for i := 2; i <= x; i++ {
		thisStore.Sequence[i] = FibonacciElement{thisStore.Sequence[i-2].Y + thisStore.Sequence[i-1].Y, time.Since(start)}
	}

	returnThis := thisStore.Sequence[x]
	return returnThis, nil
}

type FibonacciGenerator struct {
}

func (thisGenerator FibonacciGenerator) CalculateFibonacci(x int) int {
	return 0
}

func TestFibonacciMemoryStore(t *testing.T) {
	testUpperBound := 1000
	var testStore IFibonacciStore = NewMemoryFibonacciStore(testUpperBound)

	expectedMax := 1

	actualMax := testStore.GetMax()

	if actualMax != expectedMax {
		t.Fatalf("Max stored value incorrect: should be %d but was %d", expectedMax, actualMax)
	}

	_, err := testStore.GetValue(testUpperBound + 1)
	if err == nil {
		t.Fatal("Value over upper bound should throw error")
	}

	expectedValues := map[int]int{
		0:  0,
		1:  1,
		2:  1,
		3:  2,
		4:  3,
		5:  5,
		6:  8,
		13: 233,
	}

	for k, expectedValue := range expectedValues {
		actualValue, err := testStore.GetValue(k)

		if err != nil {
			t.Fatal("Error retrieving value", err)
		}
		if actualValue.Y != expectedValue {
			t.Fatalf("Value incorrect: should be %d but was %d", expectedValue, actualValue)
		}
	}
}

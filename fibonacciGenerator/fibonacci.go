package fibonacciGenerator

import (
	"fmt"
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

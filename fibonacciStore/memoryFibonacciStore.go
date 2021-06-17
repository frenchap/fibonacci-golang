package fibonacciStore

import (
	"fmt"
	"math/big"
	"time"
)

type MemoryFibonacciStore struct {
	UpperBound int
	Sequence   map[int]*FibonacciElement
}

func NewMemoryFibonacciStore(upperBound int) MemoryFibonacciStore {
	return MemoryFibonacciStore{
		UpperBound: upperBound,
		Sequence: map[int]*FibonacciElement{
			0: {big.NewInt(0), 0},
			1: {big.NewInt(1), 0},
		},
	}
}

func (thisStore MemoryFibonacciStore) GetMax() int {
	return len(thisStore.Sequence) - 1
}

func (thisStore MemoryFibonacciStore) GetValue(x int) (*FibonacciElement, error) {
	if value, valueFound := thisStore.Sequence[x]; valueFound {
		return value, nil
	}
	if x > thisStore.UpperBound {
		return nil, fmt.Errorf("requested value of %d above upper bound of %d", x, thisStore.UpperBound)
	}

	start := time.Now()
	for i := 2; i <= x; i++ {
		nextValue := big.NewInt(0)
		nextValue = nextValue.Add(thisStore.Sequence[i-2].Y, thisStore.Sequence[i-1].Y)
		nextElement := FibonacciElement{nextValue, time.Since(start)}
		thisStore.Sequence[i] = &nextElement
	}

	returnThis := thisStore.Sequence[x]
	return returnThis, nil
}

func (thisStore MemoryFibonacciStore) GetIntermediateValueCount(y *big.Int) int {
	return 0
}

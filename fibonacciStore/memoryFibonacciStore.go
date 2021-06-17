package fibonacciStore

import (
	"fmt"
	"math/big"
	"time"

	"github.com/sirupsen/logrus"
)

type MemoryFibonacciStore struct {
	UpperBound  int
	Sequence    map[int]*FibonacciElement
	GoldenRatio *big.Float
}

func NewMemoryFibonacciStore(upperBound int) MemoryFibonacciStore {
	goldenRatio := big.NewFloat(0)
	goldenRatio, succeeded := goldenRatio.SetString("1.618033988749895")
	if !succeeded {
		logrus.Error("Error setting golden ratio from string")
	}

	return MemoryFibonacciStore{
		UpperBound: upperBound,
		Sequence: map[int]*FibonacciElement{
			0: {big.NewInt(0), 0},
			1: {big.NewInt(1), 0},
		},
		GoldenRatio: goldenRatio,
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

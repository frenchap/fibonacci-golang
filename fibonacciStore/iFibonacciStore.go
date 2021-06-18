package fibonacciStore

import "math/big"

type IFibonacciStore interface {
	GetMax() int
	GetValue(x int) (*FibonacciElement, error)
	GetIntermediateValueCount(y *big.Int) int
	ClearStore()
}

package fibonacciStore

import (
	"math/big"
	"time"
)

type FibonacciElement struct {
	Y        *big.Int
	TimeCost time.Duration
}

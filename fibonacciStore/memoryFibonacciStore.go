package fibonacciStore

import (
	"flag"
	"fmt"
	"math/big"
	"time"

	"github.com/sirupsen/logrus"
	"robpike.io/ivy/config"
	"robpike.io/ivy/exec"
	"robpike.io/ivy/value"
)

type MemoryFibonacciStore struct {
	UpperBound  int
	Sequence    map[int]*FibonacciElement
	GoldenRatio *big.Float
	Conf        config.Config
}

func NewMemoryFibonacciStore(upperBound int) MemoryFibonacciStore {
	goldenRatio := big.NewFloat(0)
	goldenRatio, succeeded := goldenRatio.SetString("1.618033988749895")
	if !succeeded {
		logrus.Error("Error setting golden ratio from string")
	}

	conf := config.Config{}
	if flag.Lookup("format") == nil {
		conf.SetFormat(*flag.String("format", "", "use `fmt` as format for printing numbers; empty sets default format"))
	}
	if flag.Lookup("maxbits") == nil {
		conf.SetMaxBits(*flag.Uint("maxbits", 1e9, "maximum size of an integer, in bits; 0 means no limit"))
	}
	if flag.Lookup("maxdigits") == nil {
		conf.SetMaxDigits(*flag.Uint("maxdigits", 0, "above this many `digits`, integers print as floating point; 0 disables"))
	}
	if flag.Lookup("origin") == nil {
		conf.SetOrigin(*flag.Int("origin", 1, "set index origin to `n` (must be 0 or 1)"))
	}

	return MemoryFibonacciStore{
		UpperBound: upperBound,
		Sequence: map[int]*FibonacciElement{
			0: {big.NewInt(0), 0},
			1: {big.NewInt(1), 0},
		},
		GoldenRatio: goldenRatio,
		Conf:        conf,
	}
}

func (thisStore MemoryFibonacciStore) GetMax() int {
	return len(thisStore.Sequence) - 1
}

func (thisStore MemoryFibonacciStore) ClearStore() {

	for k := range thisStore.Sequence {
		delete(thisStore.Sequence, k)
	}

	thisStore.Sequence[0] = &FibonacciElement{big.NewInt(0), 0}
	thisStore.Sequence[1] = &FibonacciElement{big.NewInt(1), 0}

}

func (thisStore MemoryFibonacciStore) GetValue(x int) (*FibonacciElement, error) {

	if x > thisStore.UpperBound {
		return nil, fmt.Errorf("requested value of %d above upper bound of %d", x, thisStore.UpperBound)
	}

	if value, valueFound := thisStore.Sequence[x]; valueFound {
		return value, nil
	}

	start := time.Now()

	for i := 2; i <= x; i++ {
		nextValue := big.NewInt(0)
		nextValue = nextValue.Add(thisStore.Sequence[i-2].Y, thisStore.Sequence[i-1].Y)
		nextElement := FibonacciElement{nextValue, time.Since(start).Microseconds()}
		thisStore.Sequence[i] = &nextElement
	}

	returnThis := thisStore.Sequence[x]
	return returnThis, nil
}

//given a potential fibonacci result, this will return the number of results less than that value
func (thisStore MemoryFibonacciStore) GetIntermediateValueCount(y *big.Int) int {

	//The formula for the oridinal, given a possible result is Xn = (ln(y * sqrt(5))) / ln(g)
	//where Xn is the ordinal, y is the given value, and g is the golden ratio (1.618033...)
	floatY := new(big.Float).SetInt(y)

	//Build numerator
	floatYSqrt5 := new(big.Float).Sqrt(big.NewFloat(5))
	floatYSqrt5 = floatYSqrt5.Mul(floatY, floatYSqrt5)

	context := exec.NewContext(&thisStore.Conf)

	floatLnYSqrt5 := context.EvalUnary("log", value.BigFloat{floatYSqrt5})

	//Build denominator
	floatLne := context.EvalUnary("log", value.BigFloat{thisStore.GoldenRatio})

	//Evaluate and take the floor (gives the next lower value)
	floatX := context.EvalBinary(floatLnYSqrt5, "/", floatLne)
	floatXFloor := context.EvalUnary("floor", floatX)

	//Convert to a big int
	bigIntX, succeeded := new(big.Int).SetString(floatXFloor.Sprint(&thisStore.Conf), 10)
	if !succeeded {
		logrus.Error("Error converting back to big int")
	}

	//Convert back to a native int and return
	return int(bigIntX.Int64()) + 1
}

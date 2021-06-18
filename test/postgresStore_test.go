package test

import (
	"math/big"
	"testing"

	"github.com/frenchap/fibonacci-golang/fibonacciStore"
)

func TestPostgresStoreClearStore(t *testing.T) {
	testMeta := setup()

	var testPostgresStore fibonacciStore.IFibonacciStore
	testPostgresStore, err := fibonacciStore.NewPostgresFibonacciStore(testMeta.TestUpperBound, testMeta.Dialect, testMeta.DataSourceName)
	if err != nil {
		t.Fatal("Error creating new postgres store: ", err)
	}

	expectedMax := testMeta.TestUpperBound
	testPostgresStore.GetValue(testMeta.TestUpperBound)
	actualMax := testPostgresStore.GetMax()

	if expectedMax != actualMax {
		t.Fatalf("Max stored value incorrect: should be %d but was %d", expectedMax, actualMax)
	}

	testPostgresStore.ClearStore()

	expectedMax = 1 //all values cleared except 0 and 1; store requires them for GetValue
	actualMax = testPostgresStore.GetMax()
	if expectedMax != actualMax {
		t.Fatalf("Max stored value incorrect after ClearStore: should be %d but was %d", expectedMax, actualMax)
	}

	teardown(testMeta)
}

func TestPostgresStoreGetIntermediateValueCount(t *testing.T) {
	testMeta := setup()

	var testPostgresStore fibonacciStore.IFibonacciStore
	testPostgresStore, err := fibonacciStore.NewPostgresFibonacciStore(testMeta.TestUpperBound, testMeta.Dialect, testMeta.DataSourceName)
	if err != nil {
		t.Fatal("Error creating new postgres store: ", err)
	}

	intermediateCountTestMap := map[*big.Int]int{big.NewInt(377): 14, big.NewInt(378): 15, big.NewInt(102334155): 40, big.NewInt(102334156): 41, big.NewInt(120): 12}

	for value, expectedCount := range intermediateCountTestMap {
		actualCount := testPostgresStore.GetIntermediateValueCount(value)
		if actualCount != expectedCount {
			t.Fatalf("GetIntermediateValueCount incorrect: should be %d but was %d", expectedCount, actualCount)
		}
	}

	teardown(testMeta)
}

func TestPostgresStoreExpectedValues(t *testing.T) {
	testMeta := setup()

	var testPostgresStore fibonacciStore.IFibonacciStore
	testPostgresStore, err := fibonacciStore.NewPostgresFibonacciStore(testMeta.TestUpperBound, testMeta.Dialect, testMeta.DataSourceName)
	if err != nil {
		t.Fatal("Error creating new postgres store: ", err)
	}

	threeHundred := big.NewInt(0)
	threeHundred, succeeded := threeHundred.SetString("222232244629420445529739893461909967206666939096499764990979600", 10)
	if !succeeded {
		t.Fatal("Error creating big.int for 300th fibonacci element")
	}

	ninetyThree := big.NewInt(0)
	ninetyThree, succeeded = ninetyThree.SetString("12200160415121876738", 10)
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
		92:  *big.NewInt(7540113804746346429),
		93:  *ninetyThree,
		300: *threeHundred,
	}

	for k, expectedValue := range expectedValues {
		actualValue, err := testPostgresStore.GetValue(k)

		if err != nil {
			t.Fatal("Error retrieving value", err)
		}
		if actualValue.Y.Cmp(&expectedValue) != 0 {
			t.Fatalf("Value incorrect: should be %s but was %s", expectedValue.String(), actualValue.Y.String())
		}
	}

	teardown(testMeta)
}

package test

import (
	"testing"

	"gitlab.com/frenchap/fibonacci-golang/fibonacciGenerator"
)

func TestFibonacciMemoryStore(t *testing.T) {
	testUpperBound := 1000
	var testStore fibonacciGenerator.IFibonacciStore = fibonacciGenerator.NewMemoryFibonacciStore(testUpperBound)

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

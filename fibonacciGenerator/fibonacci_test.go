package fibonacciGenerator

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	logrus.Infoln("Starting setup")

	jsonFile, err := os.Open("./../.env.local.json")
	if err != nil {
		logrus.Error("Error reading local env file", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]string
	json.Unmarshal([]byte(byteValue), &result)

	for key, element := range result {

		os.Setenv(key, element)
	}

	logrus.Infoln("Exiting setup")
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

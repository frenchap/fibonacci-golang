package test

import "github.com/frenchap/fibonacci-golang/fibonacciStore"

type TestMeta struct {
	DbUser     string
	DbPassword string
	DbName     string
	DbPort     string

	DataSourceName string

	TestUpperBound  int
	TestMemoryStore fibonacciStore.IFibonacciStore
}

func NewTestMeta() TestMeta {
	upperBound := 500000 //500000 took about 5 seconds to calculate and made my processor fan kick on.
	return TestMeta{
		DbUser:          "postgres",
		DbPassword:      "12345-luggage-combo",
		DbName:          "postgres",
		DbPort:          "5432",
		DataSourceName:  "postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		TestUpperBound:  upperBound,
		TestMemoryStore: fibonacciStore.NewMemoryFibonacciStore(upperBound),
	}
}

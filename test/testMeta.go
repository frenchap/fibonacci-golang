package test

import (
	"fmt"

	"github.com/frenchap/fibonacci-golang/fibonacciStore"
	"github.com/ory/dockertest/v3"
)

type TestMeta struct {
	DbUser     string
	DbPassword string
	DbName     string
	DbPort     string

	DataSourceName string

	TestUpperBound  int
	TestMemoryStore fibonacciStore.IFibonacciStore

	Dialect string

	Pool     *dockertest.Pool
	Resource *dockertest.Resource
}

func NewTestMeta() *TestMeta {
	upperBound := 50000 //500000 took about 5 seconds to calculate and made my processor fan kick on.
	dialect := "postgres"
	dsnString := "postgres://%s:%s@localhost:%s/%s?sslmode=disable"
	dbUser := "postgres"
	dbPassword := "12345-luggage-combo"
	dbName := "postgres"
	dbPort := "5432"
	dsn := fmt.Sprintf(dsnString, dbUser, dbPassword, dbPort, dbName)

	returnThis := TestMeta{
		DbUser:          dbUser,
		DbPassword:      dbPassword,
		DbName:          dbName,
		DbPort:          dbPort,
		DataSourceName:  dsn,
		TestUpperBound:  upperBound,
		TestMemoryStore: fibonacciStore.NewMemoryFibonacciStore(upperBound),
		Dialect:         dialect,
	}

	return &returnThis
}

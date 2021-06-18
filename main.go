package main

import (
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strconv"

	"github.com/frenchap/fibonacci-golang/fibonacciStore"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

func main() {

	//500000 took about 5 seconds to calculate and made my processor fan kick on with the memory store.
	//75000 took over five minutes to calculate using the postgresStore
	//Using 50000 for live build
	upperBoundString := os.Getenv("FI_API_UPPER_BOUND")
	upperBound, err := strconv.Atoi(upperBoundString)
	if err != nil {
		logrus.Fatalf("Error converting upper bound (%s) to integer: ", upperBoundString)
	}
	dbUser := os.Getenv("FI_API_DB_USER")
	dbPassword := os.Getenv("FI_API_DB_PASSWORD")
	dbName := os.Getenv("FI_API_DB_NAME")
	dbPort := os.Getenv("FI_API_DB_PORT")

	dsnString := "postgres://%s:%s@localhost:%s/%s?sslmode=disable"
	dialect := "postgres"
	dsn := fmt.Sprintf(dsnString, dbUser, dbPassword, dbPort, dbName)

	var postgresFibonacciStore fibonacciStore.IFibonacciStore
	postgresFibonacciStore, err = fibonacciStore.NewPostgresFibonacciStore(upperBound, dialect, dsn)
	if err != nil {
		logrus.Fatal("Error creating new postgres store: ", err)
	}

	mw := io.MultiWriter(os.Stdout)
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
		Output: mw,
	}))

	e.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})

	e.GET("/values/:x", func(c echo.Context) error {
		x := c.Param("x")
		i, err := strconv.Atoi(x)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "parameter must be numeric string")
		}
		value, err := postgresFibonacciStore.GetValue(i)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "error retrieving value")
		}
		return c.JSON(http.StatusOK, value)

	})

	e.DELETE("/values", func(c echo.Context) error {
		postgresFibonacciStore.ClearStore()
		return c.JSON(http.StatusOK, "Store cleared")
	})

	e.GET("/ordinals/:y", func(c echo.Context) error {
		y := c.Param("y")
		yInt, succeeded := new(big.Int).SetString(y, 10)
		if !succeeded {
			return c.JSON(http.StatusBadRequest, "parameter must be numeric string")
		}
		return c.JSON(http.StatusOK, postgresFibonacciStore.GetIntermediateValueCount(yInt))
	})

	e.Logger.Fatal(e.Start(":8080"))
}

package test

import (
	"io"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/frenchap/fibonacci-golang/fibonacciStore"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func TestAPIStartup(t *testing.T) {

	testMeta := setup()

	var testPostgresStore fibonacciStore.IFibonacciStore
	testPostgresStore, err := fibonacciStore.NewPostgresFibonacciStore(testMeta.TestUpperBound, testMeta.Dialect, testMeta.DataSourceName)
	if err != nil {
		t.Fatal("Error creating new postgres store: ", err)
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
		value, err := testPostgresStore.GetValue(i)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "error retrieving value")
		}
		return c.JSON(http.StatusOK, value)

	})

	e.DELETE("/values", func(c echo.Context) error {
		testPostgresStore.ClearStore()
		return c.JSON(http.StatusOK, "Store cleared")
	})

	e.GET("/ordinals/:y", func(c echo.Context) error {
		y := c.Param("y")
		yInt, succeeded := new(big.Int).SetString(y, 10)
		if !succeeded {
			return c.JSON(http.StatusBadRequest, "parameter must be numeric string")
		}
		return c.JSON(http.StatusOK, testPostgresStore.GetIntermediateValueCount(yInt))
	})

	e.Logger.Fatal(e.Start(":8080"))

	teardown(testMeta)
}

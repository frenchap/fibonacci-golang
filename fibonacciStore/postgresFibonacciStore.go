package fibonacciStore

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"time"

	"github.com/sirupsen/logrus"
	"robpike.io/ivy/config"
	"robpike.io/ivy/exec"
	"robpike.io/ivy/value"
)

type fibonacciModel struct {
	id    int
	value string
	cost  int
}

type PostgresFibonacciStore struct {
	UpperBound  int
	GoldenRatio *big.Float
	Conf        config.Config
	Db          *sql.DB
}

func NewPostgresFibonacciStore(upperBound int, dialect string, dsn string) (*PostgresFibonacciStore, error) {

	db, err := sql.Open(dialect, dsn)
	if err != nil {
		logrus.Error("Error creating db: ", err)
		return nil, err
	}
	ctx := context.Background()

	createQuery :=
		"CREATE TABLE IF NOT EXISTS fibonacci (" +
			"id INT PRIMARY KEY, " +
			"value text, " +
			"cost int " +
			");"

	time.Sleep(2 * time.Second)
	createStatement, err := db.PrepareContext(ctx, createQuery)
	if err != nil {
		logrus.Error("Error preparing create statement: ", err)
		return nil, err
	}
	defer createStatement.Close()

	_, err = createStatement.ExecContext(ctx)
	if err != nil {
		logrus.Error("Error creating new postgres fibonacci store: ", err)
		return nil, err
	}

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
		conf.SetMaxBits(*flag.Uint("maxbits", 0, "maximum size of an integer, in bits; 0 means no limit"))
	}
	if flag.Lookup("maxdigits") == nil {
		conf.SetMaxDigits(*flag.Uint("maxdigits", 0, "above this many `digits`, integers print as floating point; 0 disables"))
	}
	if flag.Lookup("origin") == nil {
		conf.SetOrigin(*flag.Int("origin", 1, "set index origin to `n` (must be 0 or 1)"))
	}
	conf.SetFormat("%.0f")

	newStore := PostgresFibonacciStore{
		UpperBound:  upperBound,
		GoldenRatio: goldenRatio,
		Conf:        conf,
		Db:          db,
	}

	newStore.truncateTable()
	newStore.insertStartValues()

	return &newStore, nil
}

func (thisStore PostgresFibonacciStore) ClearStore() {
	thisStore.truncateTable()
	thisStore.insertStartValues()
}

//given a potential fibonacci result, this will return the number of results less than that value
func (thisStore PostgresFibonacciStore) GetIntermediateValueCount(y *big.Int) int {
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

func (thisStore PostgresFibonacciStore) GetMax() int {
	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel()

	id := 0

	row := thisStore.Db.QueryRowContext(ctx, "SELECT id FROM fibonacci ORDER BY id DESC LIMIT 1")

	err := row.Scan(&id)

	if err == nil {
		return id
	}

	if err == sql.ErrNoRows {

		return 0
	}

	logrus.Error("Error getting max value for x (%d)", err)
	return 0
}

func (thisStore PostgresFibonacciStore) GetValue(x int) (*FibonacciElement, error) {

	if x > thisStore.UpperBound {
		return nil, fmt.Errorf("requested value of %d above upper bound of %d", x, thisStore.UpperBound)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel()

	id := 0
	y := ""
	cost := int64(0)

	row := thisStore.Db.QueryRowContext(ctx, "SELECT id, value, cost FROM fibonacci where id = $1", x)

	err := row.Scan(&id, &y, &cost)

	if err == nil {
		y, succeeded := new(big.Int).SetString(y, 10)

		if !succeeded {
			return nil, errors.New(fmt.Sprintf("Error setting big int from string %s", y))
		}
		return &FibonacciElement{
			Y:        y,
			TimeCost: cost,
		}, nil
	}

	if err == sql.ErrNoRows {

		fibonacciValuesMap := make(map[int]FibonacciElement, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
		defer cancel()

		rows, err := thisStore.Db.QueryContext(ctx, "SELECT id, value, cost FROM fibonacci ORDER BY id DESC LIMIT 2")
		if err != nil {
			logrus.Error("Erro querying database: ", err)
			return nil, err
		}
		defer rows.Close()

		startHere := 0
		for rows.Next() {
			fm := new(fibonacciModel)
			err = rows.Scan(
				&fm.id,
				&fm.value,
				&fm.cost,
			)

			if err != nil {
				logrus.Error("Error iterating values: ", err)
				return nil, err
			}
			yVal, succeeded := new(big.Int).SetString(fm.value, 10)
			if !succeeded {
				return nil, errors.New(fmt.Sprintf("Error convert db value (%s) to bigint", fm.value))
			}
			if startHere == 0 {
				startHere = fm.id + 1
			}
			fibonacciValuesMap[fm.id] = FibonacciElement{
				Y:        yVal,
				TimeCost: int64(fm.cost),
			}
		}

		start := time.Now()

		var nextElement *FibonacciElement
		for i := startHere; i <= x; i++ {
			nextValue := big.NewInt(0)
			nextValue = nextValue.Add(fibonacciValuesMap[i-2].Y, fibonacciValuesMap[i-1].Y)
			nextElement = &FibonacciElement{nextValue, time.Since(start).Microseconds()}
			fibonacciValuesMap[i] = *nextElement

			ctx, cancel = context.WithTimeout(context.Background(), 240*time.Second)
			defer cancel()

			insertQuery := "INSERT INTO fibonacci (id, value, cost) VALUES ($1, $2, $3)"
			insertStatement, err := thisStore.Db.PrepareContext(ctx, insertQuery)

			if err != nil {
				logrus.Error("Error preparing insert statement: ", err)
				return nil, err
			}
			defer insertStatement.Close()

			_, err = insertStatement.ExecContext(ctx, i, nextElement.Y.String(), time.Since(start).Microseconds())
			if err != nil {
				logrus.Error("Error executing insert statement: ", err)
				return nil, err
			}
		}

		return nextElement, nil
	}

	logrus.Errorf("Error getting value for x (%d)", x, err)
	return nil, err
}

func (thisStore PostgresFibonacciStore) truncateTable() {
	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel()

	truncateQuery := "TRUNCATE TABLE fibonacci"
	truncateStatement, err := thisStore.Db.PrepareContext(ctx, truncateQuery)
	if err != nil {
		logrus.Error("Error preparing truncate queyr: ", err)
	}
	defer truncateStatement.Close()

	_, err = truncateStatement.ExecContext(ctx)
	if err != nil {
		logrus.Error("Error executign truncate query: ", err)
	}

	//"TRUNCATE TABLE fibonacci;" +
	//"INSERT INTO TABLE fibonacci (id, value) " +
	//"VALUES " +
	//"(0, '0', '0'), " +
	//"(1, '1', '1');"
}

func (thisStore PostgresFibonacciStore) insertStartValues() {

	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel()

	insertQuery := "INSERT INTO fibonacci (id, value, cost) VALUES (0, '0', 0), (1, '1', 1);"
	insertStatement, err := thisStore.Db.PrepareContext(ctx, insertQuery)

	if err != nil {
		logrus.Error("Error preparing insert statement: ", err)
	}
	defer insertStatement.Close()

	_, err = insertStatement.ExecContext(ctx)
	if err != nil {
		logrus.Error("Error executing insert statement: ", err)
	}

}

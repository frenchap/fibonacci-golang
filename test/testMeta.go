package test

type TestMeta struct {
	dbUser     string
	dbPassword string
	dbName     string
	dbPort     string

	dataSourceName string
}

func NewTestMeta() TestMeta {
	return TestMeta{
		dbUser:         "postgres",
		dbPassword:     "12345-luggage-combo",
		dbName:         "postgres",
		dbPort:         "5432",
		dataSourceName: "postgres://%s:%s@localhost:%s/%s?sslmode=disable",
	}

}

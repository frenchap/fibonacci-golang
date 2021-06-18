package test

import (
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sirupsen/logrus"
)

// func TestMain(m *testing.M) {
// 	setup()
// 	code := m.Run()
// 	os.Exit(code)
// }

func setup() *TestMeta {
	logrus.Infoln("Starting setup")

	testMeta := NewTestMeta()

	pool, err := dockertest.NewPool("")
	if err != nil {
		logrus.Fatal("Error creating new pool: ", err)
	}

	testMeta.Pool = pool

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.3",
		Env: []string{
			"POSTGRES_USER=" + testMeta.DbUser,
			"POSTGRES_PASSWORD=" + testMeta.DbPassword,
			"POSTGRES_DB=" + testMeta.DbName,
		},
		ExposedPorts: []string{testMeta.DbPort},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: testMeta.DbPort},
			},
		},
	}

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		logrus.Fatal("Error running pool: ", err)
	}

	testMeta.Resource = resource
	logrus.Infof("Resource: %+v", resource)

	logrus.Info("Exiting setup")

	return testMeta
}

func teardown(testMeta *TestMeta) {

	logrus.Infoln("Entering teardown")

	if err := testMeta.Pool.Purge(testMeta.Resource); err != nil {
		logrus.Fatal("Error purging resource: ", err)
	}

	logrus.Infoln("Exiting teardown")
}

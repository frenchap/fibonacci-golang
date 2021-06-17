package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	logrus.Infoln("Starting setup")

	testMeta := NewTestMeta()

	pool, err := dockertest.NewPool("")
	if err != nil {
		logrus.Error("Error creating new pool: ", err)
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.3",
		Env: []string{
			"POSTGRES_USER=" + testMeta.dbUser,
			"POSTGRES_PASSWORD=" + testMeta.dbPassword,
			"POSTGRES_DB=" + testMeta.dbName,
		},
		ExposedPorts: []string{testMeta.dbPort},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: testMeta.dbPort},
			},
		},
	}

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("Error running pool: %s", err)
	}

	testMeta.dataSourceName = fmt.Sprintf(testMeta.dataSourceName, testMeta.dbUser, testMeta.dbPassword, testMeta.dbPort, testMeta.dbName)

	logrus.Infof("Resource: %+v", resource)

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

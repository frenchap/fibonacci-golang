package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		logrus.Infof("Error running pool: %+v", err)
	}

	testMeta.DataSourceName = fmt.Sprintf(testMeta.DataSourceName, testMeta.DbUser, testMeta.DbPassword, testMeta.DbPort, testMeta.DbName)

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

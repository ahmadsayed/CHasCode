package ceas

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

type Testcase struct {
	appname    string
	yamlPath   string
	dockerPath string
}

func setupTestCase(t *testing.T) func(t *testing.T) {
	t.Log("setup test case")
	CreateCluster(1, 3)
	return func(t *testing.T) {
		t.Log("teardown delete the whole cluster")
		DeleteCluster()
	}
}

func setupSubTest(t *testing.T, tc Testcase) func(t *testing.T) {
	t.Log("setup sub test")
	BuildAndLoadDocker(tc.dockerPath, tc.appname)
	DeployApp(tc.yamlPath)
	return func(t *testing.T) {
		t.Log("teardown Remove App")
		RemoveApp(tc.yamlPath)
	}
}

func TestAddtionFirst(t *testing.T) {
	cases := []Testcase{
		{
			"nodejsdemo:latest",
			"C:\\Users\\AHMEDSAYEDHASSANABDE\\nodejs-demo\\deployment.yaml",
			"C:\\Users\\AHMEDSAYEDHASSANABDE\\nodejs-demo",
		},
	}
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)
	for _, tc := range cases {
		t.Run(tc.appname, func(t *testing.T) {
			teardownSubTest := setupSubTest(t, tc)
			defer teardownSubTest(t)
			// Dummy test case for now, as I am still checking setup and teardown
			result := "{\"status\":\"UP\"}"
			resp, err := http.Get("http://localhost/health")
			if err != nil {
				fmt.Println(err.Error())
			}
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			val := string(body)
			if result != val {
				t.Fatalf("expected sum %v, but got %v", val, result)
			}
		})
	}
}

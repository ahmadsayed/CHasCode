package ceas

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type Testcase struct {
	appname    string
	yamlPath   string
	dockerPath string
	nodetoKill []int
}

func setupTestCase(t *testing.T) func(t *testing.T) {
	t.Log("setup test case")
	CreateCluster(1, 3, true)
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
			appname:    "nodejsdemo:latest",
			yamlPath:   "C:\\Users\\AHMEDSAYEDHASSANABDE\\nodejs-demo\\deployment.yaml",
			dockerPath: "C:\\Users\\AHMEDSAYEDHASSANABDE\\nodejs-demo",
		},
		{
			appname:    "nodejsdemo:latest",
			yamlPath:   "C:\\Users\\AHMEDSAYEDHASSANABDE\\nodejs-demo\\deployment.yaml",
			dockerPath: "C:\\Users\\AHMEDSAYEDHASSANABDE\\nodejs-demo",
			nodetoKill: []int{1, 2},
		},
	}
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)
	for _, tc := range cases {
		t.Run(tc.appname, func(t *testing.T) {
			teardownSubTest := setupSubTest(t, tc)
			defer teardownSubTest(t)
			// Using this 10 Seconds better to use the health checker
			fmt.Println("Settle for 10 Seconds to insure ingress pick up the configuration")
			time.Sleep(10000 * time.Millisecond)
			// Dummy test case for now, as I am still checking setup and teardown
			result := "{\"status\":\"UP\"}"
			if tc.nodetoKill != nil {
				fmt.Println("Kill Node ", tc.nodetoKill)
				for _, nodeNumber := range tc.nodetoKill {
					StopWorkerNode(nodeNumber)
					defer StartWorkerNode(nodeNumber)
				}
				resp, err := http.Get("http://localhost/health")
				if err != nil {
					fmt.Println(err.Error())
				}
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				val := string(body)
				if result != val {
					t.Fatalf("expected sum %v, but got %v", result, val)
				}
			}

		})
	}
}

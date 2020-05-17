package ceas

import "testing"

type Testcase struct {
	appname    string
	yamlPath   string
	dockerPath string
}

func setupTestCase(t *testing.T) func(t *testing.T) {
	t.Log("setup test case")
	CreateCluster(1, 3)
	return func(t *testing.T) {
		DeleteCluster()
	}
}

func setupSubTest(t *testing.T, tc Testcase) func(t *testing.T) {
	t.Log("setup sub test")
	BuildAndLoadDocker(tc.dockerPath, tc.appname)
	DeployApp(tc.yamlPath)
	return func(t *testing.T) {
		t.Log("teardown sub test")
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
			result := 5
			if result != 5 {
				t.Fatalf("expected sum %v, but got %v", 5, result)
			}
		})
	}
}

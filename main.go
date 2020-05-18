package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/CHasCode/ceas"
)

func main() {
	// Create Cluster 1 master nodes and 3 workers
	// Install Nignx Ingress on the control-plan
	//	ceas.CreateCluster(1, 3)
	// Build a docker from source code and load it to the cluster registry not required
	//	ceas.BuildAndLoadDocker("C:\\Users\\AHMEDSAYEDHASSANABDE\\nodejs-demo", "nodejsdemo:latest")
	// Deploy app on the cluster
	// Remove app from the cluster
	// Destroy the kubernetes cluster
	//ceas.DeleteCluster()

	ceas.DeployApp("C:\\Users\\AHMEDSAYEDHASSANABDE\\nodejs-demo\\deployment.yaml")

	resp, err := http.Get("http://localhost/health")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	result := "{\"status\":\"UP\"}"
	fmt.Println(string(body) == result)

	ceas.RemoveApp("C:\\Users\\AHMEDSAYEDHASSANABDE\\nodejs-demo\\deployment.yaml")

}

package main

import "github.com/CHasCode/ceas"

func main() {
	// Create Cluster 1 master nodes and 3 workers
	// Install Nignx Ingress on the control-plan
	ceas.CreateCluster(1, 3)
	// Build a docker from source code and load it to the cluster registry not required
	ceas.BuildAndLoadDocker("C:\\Users\\AHMEDSAYEDHASSANABDE\\nodejs-demo", "nodejsdemo:latest")
	// Deploy app on the cluster
	ceas.DeployApp("C:\\Users\\AHMEDSAYEDHASSANABDE\\nodejs-demo\\deployment.yaml")
	// Remove app from the cluster
	ceas.RemoveApp("C:\\Users\\AHMEDSAYEDHASSANABDE\\nodejs-demo\\deployment.yaml")
	// Destroy the kubernetes cluster
	ceas.DeleteCluster()
}

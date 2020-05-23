package ceas

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// ClusterConfig hold kind cluster configuration
type ClusterConfig struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
	Nodes      []Node `yaml:"nodes"`
}

// Node cluster topology
type Node struct {
	Role              string             `yaml:"role"`
	ExtraPortMappings []ExtraPortMapping `yaml:"extraPortMappings"`
}

// ExtraPortMappings extraPortMappings
type ExtraPortMapping struct {
	ContainerPort int    `yaml:"containerPort"`
	HostPort      int    `yaml:"hostPort"`
	ListenAddress string `yaml:"listenAddress"`
}

// CreateClusterConfig will create cluster config from numberOfMasters and numberOfWorkers and return clean up function
func CreateClusterConfig(numberOfMasters, numberOfWorkers int) func() {
	var buffer bytes.Buffer
	filename := "Cluster.yaml"
	totalNumberOfNodes := numberOfMasters + numberOfWorkers
	nodes := make([]Node, totalNumberOfNodes, totalNumberOfNodes)

	extraPortMappings := []ExtraPortMapping{
		{
			80,
			80,
			"0.0.0.0",
		},
	}
	for i := 0; i < numberOfMasters; i++ {
		nodes[i] = Node{
			Role:              "control-plane",
			ExtraPortMappings: extraPortMappings,
		}
	}
	for i := numberOfMasters; i < totalNumberOfNodes; i++ {
		nodes[i] = Node{Role: "worker"}
	}

	clusterConfig := ClusterConfig{
		Kind:       "Cluster",
		APIVersion: "kind.x-k8s.io/v1alpha4",
		Nodes:      nodes,
	}
	err := yaml.NewEncoder(&buffer).Encode(clusterConfig)
	if err != nil {
		fmt.Println(err)
	}
	_ = ioutil.WriteFile(filename, buffer.Bytes(), 0644)
	return func() {
		os.Remove(filename)
	}
}

// CheckIfAllNodesAreReady Check if all nodes in Ready State
func CheckIfAllNodesAreReady() bool {
	//poorman hack get retrun stream and check it does not contains NotReady ;)
	cmd := exec.Command("kubectl", "get", "nodes")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	return !strings.Contains(string(stdout), "NotReady")

}

// CheckIfKindClusterExists exec kind get clusters
func CheckIfKindClusterExists() bool {
	noClusterMessage := "No kind clusters found."
	cmd := exec.Command("kind", "get", "clusters")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	results := string(stdout)

	return !strings.Contains(noClusterMessage, results)
}

// DeleteClusterIfExists delete cluster if exists prompt the user to take actions
func DeleteClusterIfExists() {
	if CheckIfKindClusterExists() {
		fmt.Print("A Kind Cluster found Type y, to delete : ")
		reader := bufio.NewReader(os.Stdin)
		char, _, err := reader.ReadRune()
		if err != nil {
			fmt.Println(err)
		}
		if char == 'y' {
			fmt.Println("Deleting the Cluster")
			cmd := exec.Command("kind", "delete", "cluster")
			stdout, err := cmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(string(stdout))

		} else {
			fmt.Println("Please Delete the cluster manually and run again")
			os.Exit(0)
		}
	}
}

func portForwardWeave() {
	cmd := exec.Command("kubectl", "port-forward", "-n", "weave", "deployment/weave-scope-app", "4040")
	//cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)
}

//InstallVisulization
func InstallVisulization() {

	fmt.Println("Installing weavescope")
	cmd := exec.Command("kubectl", "apply", "-f", "https://cloud.weave.works/k8s/scope.yaml")
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Wait for weave to get to running state")
	for {
		cmd := exec.Command("kubectl", "get", "po", "-n", "weave")
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if strings.Contains(string(stdout), "1/1") && strings.Contains(string(stdout), "Running") {
			break
		}
		fmt.Println("Wait for weave to start ...")
		time.Sleep(10000 * time.Millisecond)

	}

	portForwardWeave()
	openbrowser("http://localhost:4040/#!/state/{%22topologyId%22:%22hosts%22}")

}

// CreateCluster create a kind cluster
func CreateCluster(numberOfMasters, numberOfWorkers int, visualization bool) {

	// Check if kind cluster exists delete if exists
	DeleteClusterIfExists()
	// Create ClusterConfig.yaml delete the file after done with defer
	defer CreateClusterConfig(numberOfMasters, numberOfWorkers)()

	// Create the cluster
	cmd := exec.Command("kind", "create", "cluster", "--wait", "30s", "--config", "Cluster.yaml")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(stdout))

	// Check if the cluster successfully created
	if CheckIfKindClusterExists() {
		fmt.Println("Successfully Created the Cluster, Wait for nodes to be ready")
	}

	// Check if all nodes in Ready State
	for {
		if CheckIfAllNodesAreReady() {
			break
		}
		fmt.Println("Nodes not yet Ready waiting for all nodes to be Ready .. ")
		time.Sleep(5000 * time.Millisecond)

	}

	// label control plan for ingress=ready
	cmd_label := exec.Command("kubectl", "label", "node", "kind-control-plane", "ingress-ready=true")
	stdout, err = cmd_label.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// install ingress nginx ingress controller
	//kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml
	fmt.Println("Install ingress nginx ingress controller ... ")
	cmd_ingress := exec.Command("kubectl", "apply", "-f", "https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml")
	stdout, err = cmd_ingress.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Wait for Ingress to be in in running state
	// kubectl get po -n ingress-nginx
	for {
		cmd := exec.Command("kubectl", "get", "po", "-n", "ingress-nginx")
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if strings.Contains(string(stdout), "1/1") && strings.Contains(string(stdout), "Running") {
			break
		}
		fmt.Println("Wait for Ingress to start ...")

		time.Sleep(10000 * time.Millisecond)

	}

	// Install weavescope

	if visualization {
		InstallVisulization()
	}

}

// DeleteCluster delete a kind cluster
func DeleteCluster() {
	cmd := exec.Command("kind", "delete", "cluster")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(stdout))
}

// CheckIfClusterCreated function to check the cluster
func CheckIfClusterCreated() bool {
	exec.Command("kind", "export", "kubeconfig")

	cmd := exec.Command("kubectl", "get", "nodes")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(stdout))
	return false
}

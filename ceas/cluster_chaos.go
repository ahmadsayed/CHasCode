package ceas

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func manageWorker(action string, node int) error {
	if node < 1 {
		return errors.New("Nodes starts with one")
	}

	//TODO: Add check if the node is >  the actual number of worker nodes

	workerNodeDocker := "kind-worker"
	if node != 1 {
		workerNodeDocker = fmt.Sprintf("%v%d", workerNodeDocker, node)
	}
	cmd := exec.Command("docker", action, workerNodeDocker)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	return nil

}

// StopSpecificWorkerNode, docker stop worker nodes named as follow by default kind-worker, kind-worker2, kind-worker3, ..... kind-workern
func StopSpecificWorkerNode(node int) error {
	return manageWorker("stop", node)
}

// StartSpecificWorkerNode, docker start the worker node container
func StartSpecificWorkerNode(node int) error {
	return manageWorker("start", node)
}

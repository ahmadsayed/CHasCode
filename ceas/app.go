package ceas

import (
	"fmt"
	"os"
	"os/exec"
	"time"
	"strings"
)

// BuildAndLoadDocker  provide projectFolder which contains Docker file, imagename the name used in the deployment
func BuildAndLoadDocker(projectFolder, imagename string) {
	// use Docker to build the project
	fmt.Println("Build and Load Docker")

	//docker build -t my-custom-image:unique-tag ./my-image-dir
	cmd := exec.Command("docker", "build", "-t", imagename, projectFolder)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	// use kind to load to all nodes
	//kind load docker-image my-custom-image:unique-tag
	cmd = exec.Command("kind", "load", "docker-image", imagename)
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}

func manageApp(yamlPath, action string) {
	cmd := exec.Command("kubectl", action, "-f", yamlPath)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}	
}

// DeployApp deploy add by pathing deployment file
func DeployApp(yamlPath string) {
	fmt.Println("Deploying Application")
	//kind load docker-image my-custom-image:unique-tag
	manageApp(yamlPath, "apply")

	for {
		cmd := exec.Command("kubectl", "get", "po", "-n", "default")
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if strings.Contains(string(stdout), "1/1") && strings.Contains(string(stdout), "Running") {
			break
		}
		fmt.Println("Wait for app to start ...")

		time.Sleep(5000 * time.Millisecond)

	}

}

// Removing Sepcific App 
func RemoveApp(yamlPath string) {
	fmt.Println("Removing Application")
	manageApp(yamlPath, "delete")
}	
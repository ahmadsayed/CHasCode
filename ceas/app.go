package ceas

import (
	"fmt"
	"os"
	"os/exec"

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
}

// Removing Sepcific App 
func RemoveApp(yamlPath string) {
	fmt.Println("Removing Application")
	manageApp(yamlPath, "delete")
}	
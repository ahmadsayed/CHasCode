# Chaos Engineering as Code 

For Chaos Engineering, the approach I am proposing, *instead of injecting failure to running Kubernetes Cluster*, instead of *to develop a different methodology by bringing up a Kubernetes Cluster dedicated for Chaos Testing and manipulate it programmatically*
This tool can be part of the DevOps pipeline, 
The cluster will be running with the same topology of the expected target environments, for example, Multiple masters, Multiple Workers.

The Functional Unit test suite can normally run but under different simulated hazardous environments.

[![!Demp](https://img.youtube.com/vi/moMnHb8y0U8/0.jpg)](https://www.youtube.com/watch?v=moMnHb8y0U8)

## The approach is as following 
- Using an imperative programming language, the current Repository using golang
- Use the strong foundation provided in the go Unit testing instead of providing declarative syntax using XML, JSON, or Yaml to be part of the Unit test suite.
- Run quickly in relatively limited resources, the current code able to build up one master, three workers nodes cluster run a simple test case ~  300-400 Seconds, - there are still many potential optimizations to speed up this - 
- Build libraries that help Chaos Engineer programmatically with imperative language manage cluster such as but not limited to 
  - Build a cluster with a specific number of masters and worker nodes
  - Build a cluster and test the application on a different version of Kubernetes even within the cluster 
  - Gracefully kill worker nodes
  - Deploy and un-deploy specific workloads


## This repository contains

This approach is mainly a methodology, not a framework, still trying to build a standard set of use cases, which can be used directly as a minimal base set of standard test cases, also using the libraries that can build, manage Kubernetes cluster, deploy application programmatically, as well as break the cluster drop worker nodes.


The following scenario planned to be executed 
- Create a cluster and deploy applications
- Randomly Kill and bring up worker nodes 
- Build a cluster with a different version of Kubernetes 
- Build a cluster with different version withing the cluster itself, for example, master and worker different version
- Build a cluster of different versions in worker nodes.

## Current Scope and suggested potential enhancements

The whole project now is around using Kubernetes in Docker [kind](https://kind.sigs.k8s.io/docs/user/quick-start/). Also, the current code is based on  exec kind and kubectl command; maybe a better approach is to replace it with a library-based approach like k8s go client, or via Kubernetes Rest API. Still, kubectl exec is a more straightforward approach, as well as it provides some isolation layer between the unit-testing libs, and underlaying the Kubernetes version.

## How to run 
Your machine needs to have docker and kind already installed on the executable Path.

Go and edit ceas_test.go according to your test case, define the number of worker nodes 
define normal state situation, then invoke 

Point our to the docker file and the deployment descriptor, in the helm chart you can generate a helm template.

run 
```
go test
```

This is still PoC work in Progress, in very early stage, in the upcoming day will introduce the Chaos Mode, to start killing worker nodes and test its impact on the  project

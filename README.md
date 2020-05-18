# Chaos Engineering as Code 

For Chaos Engineering the approach I am  trying to attempt in here, *instead of injecting failure to running server*, the approach in here is *to develop different methodology*
Is to build the cluster on the fly, may be as part of the DevOps pipeline
The cluster will be running with the same topology of the expected target enviroments for example  Multiple masters ,Multiple Workers

## The approach is as following 
- Using imperative programming language the current Repository using golang
- Use the strong foundation provided in the go Unit testing, instead of providing declartive syntax using XML, JSON or Yaml, so it can be integrated easily.
- Build libraries that help Chaos Engineer progrmatically with imperative language manage cluster such as not limied to 
  - Build a cluster with specific number of masters and worker nodes
  - Build a cluster and test the application on different version of kubernetes even within the cluster 
  - Gracefully kill worker nodes and check the effect 
- Run quickly in relatively limited resources, the current code able to build up 1 master, 3 workers nodes cluster an run simple test case ~  300 Seconds.

## This repository contains

This approach is mainly a methodolgy not a framework, still trying to build a standard set of use cases, which simply can use, also using the libraries that can help , manage kubernetes cluster , deploy application programmatically.


```
go test 
```

The following scenario planned to be executed 
- Create a cluster and deploy applications
- Randomly Kill and bring up worker nodes 
- Build a cluster with different version of Kubernetes 
- Build a cluster with different version withing the cluster itself for example master and worker different version
- Build a cluster different version in worker nodes.

## Current Scope and Furture works

The whole project now is around using Kubernetes in Docker [kind](https://kind.sigs.k8s.io/docs/user/quick-start/), also the current code is only around exec kind and kubectl command, may be a better approach is to replace it with library based approach like k8s go client or via kubernetes Rest API, still kubectl exec is now simpler approach, as well as it provides some sort of isolation layer between the unit-testing libs, and underlaying kubernetes version

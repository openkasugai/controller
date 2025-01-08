---
weight: 1
title: "Modify go.mod"
---
# Modify go.mod
## About This Procedure
* This document describes the areas of modification in the Kubeflow Pipelines Backend required to run the Custom Kubeflow SDK.

## Changes Made
* Edit `go.mod` to modify the source code referenced when building on the container in subsequent steps.

## Procedure
1. Clone the Kubeflow repository of the version to be modified to the Custom Kubeflow SDK execution environment. (*1)
```
$ git clone https://github.com/kubeflow/pipelines.git -b sdk-2.4.0
$ cd pipelines
```

2. Edit `go.mod`.
```
$ cd pipelines
$ vi go.mod
```

* Edit as follows.
```go
module github.com/kubeflow/pipelines
require (
	github.com/Masterminds/squirrel v0.0.0-20190107164353-fa735ea14f09
	github.com/VividCortex/mysqlerr v0.0.0-20170204212430-6c6b55f8796f
	github.com/argoproj/argo-workflows/v3 v3.3.10
	github.com/aws/aws-sdk-go v1.42.50
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/eapache/go-resiliency v1.2.0
	github.com/elazarl/goproxy v0.0.0-20181111060418-2ce16c963a8a // indirect
	github.com/emicklei/go-restful v2.16.0+incompatible // indirect
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5 // indirect
~~OMITTED~~
	sigs.k8s.io/controller-runtime v0.11.1
    sigs.k8s.io/yaml v1.3.0
+   gopkg.in/yaml.v2 v2.4.0
)

replace (
	k8s.io/kubernetes => k8s.io/kubernetes v1.11.1
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.2.9
+	github.com/kubeflow/pipelines/api => ./api
)

```

## Next Steps
[Modify argo.go](../modify-argo.go)

Additional Notes
(*1) The cloned repository is tagged, so no branch is set. To set a branch, execute the following command to create and apply any desired branch.
```
$ git switch -c <branch name>
```
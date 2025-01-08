---
weight: 3
title: "Modify visitor.go"
---
# Modify visitor.go
## About This Procedure
* Describes the areas of modification in the Kubeflow Pipelines Backend when running the Custom Kubeflow SDK.

## Changes Made
* Edit `visitor.go` to remove unused imports and functions.

## Prerequisites
* Ensure the following steps have been completed before performing this procedure.
    * [Modify argo.go](../modify-argo.go)

## Procedure
1. Open `visitor.go` in a text editor.
```
$ vi backend/src/v2/compiler/visitor.go
```

2. Edit the packages to import as follows.
```go
package compiler

import (
-	"bytes"
	"fmt"
-	"sort"

	"github.com/golang/protobuf/jsonpb"
	"github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec"
```

3. Delete various unused functions.
* Removal of the Accept function
 ```go
// visited first, then DAG components. When a DAG component is visited, it's
// guaranteed that all the components used in it have already been visited.
// * Each component is visited exactly once.

-func Accept(job *pipelinespec.PipelineJob, kubernetesSpec *pipelinespec.SinglePlatformSpec, v Visitor) error {
-	if job == nil {
-		return nil
~~OMITTED~~
-	}
-	return state.dfs(RootComponentName, spec.GetRoot())
-}

type pipelineDFS struct {
	spec           *pipelinespec.PipelineSpec
```

* Removal of the dfs function
```go
	visited map[string]bool
}

-func (state *pipelineDFS) dfs(name string, component *pipelinespec.ComponentSpec) error {
-	// each component is only visited once
-	// TODO(Bobgy): return an error when circular reference detected
~~OMITTED~~
-	// ready by the time the DAG component is visited.
-	return state.visitor.DAG(name, component, dag)
-}

func GetDeploymentConfig(spec *pipelinespec.PipelineSpec) (*pipelinespec.PipelineDeploymentConfig, error) {
```

* Removal of the GetDeploymentConfig function
```go
	visited map[string]bool
}

-func GetDeploymentConfig(spec *pipelinespec.PipelineSpec) (*pipelinespec.-PipelineDeploymentConfig, error) {
-	marshaler := jsonpb.Marshaler{}
-	buffer := new(bytes.Buffer)
~~OMITTED~~
-	}
-	return deploymentConfig, nil
-}

func GetPipelineSpec(job *pipelinespec.PipelineJob) (*pipelinespec.PipelineSpec, error) {
	// TODO(Bobgy): can we avoid this marshal to string step?
```

## Next Steps
[Modify resource_manager.go](../modify-resource_manager.go)
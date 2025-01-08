---
weight: 3
title: "visitor.go改修"
---
# visitor.go改修
## この手順について
* Custom Kubeflow SDKを実行するにあたりKubeflow Pipelines Backendの修正箇所について記載する。

## 対応内容
* `visitor.go`を編集し、使用しないimportと関数を削除する。

## 前提条件
* 当該手順を実施する前に以下手順を実施済みであること。
    * [argo.go改修](../modify-argo.go)

## 手順
1. `visitor.go`をテキストエディタで開く。
```
$ vi backend/src/v2/compiler/visitor.go
```

2. importするパッケージを以下のように編集する。
```go
package compiler

import (
-	"bytes"
	"fmt"
-	"sort"

	"github.com/golang/protobuf/jsonpb"
	"github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec"
```

3. 使用しない各種関数を削除する。
* Accept関数の削除

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

* dfs関数の削除
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

* GetDeploymentConfig関数の削除
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

## 次の手順について
[resource_manager.go改修](../modify-resource_manager.go)

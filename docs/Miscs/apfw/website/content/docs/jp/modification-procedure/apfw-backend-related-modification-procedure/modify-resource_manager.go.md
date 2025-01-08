---
weight: 4
title: "resource_manager.go改修"
---
# resource_manager.go改修
## この手順について
* Custom Kubeflow SDKを実行するにあたりKubeflow Pipelines Backendの修正箇所について記載する。

## 対応内容
* NamespaceをOpenKasugaiコントローラに合わせた設定にするため、`resource_manager.go`を編集しKubeflowで想定しているNamespaceの設定処理を削除する。

## 前提条件
* 当該手順を実施する前に以下手順を実施済みであること。
    * [visitor.go改修](../modify-visitor.go)

## 手順
1. `resource_manager.go`をテキストエディタで開く。
```
$ vi backend/src/apiserver/resource/resource_manager.go
```

2. Namespaceの設定処理を削除する。
```go
	if k8sNamespace == "" {
		return nil, util.NewInternalServerError(util.NewInvalidInputError("Namespace cannot be empty when creating an Argo workflow. Check if you have specified POD_NAMESPACE or try adding the parent namespace to the request"), "Failed to create a run due to empty namespace")
	}
-	executionSpec.SetExecutionNamespace(k8sNamespace)
	newExecSpec, err := r.getWorkflowClient(k8sNamespace).Create(ctx, executionSpec, v1.CreateOptions{})
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
```

## 次の手順について
[v2_template.go改修](../modify-v2_template.go)

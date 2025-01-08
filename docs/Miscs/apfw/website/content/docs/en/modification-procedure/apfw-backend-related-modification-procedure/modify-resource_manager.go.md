---
weight: 4
title: "Modify resource_manager.go"
---
# Modify resource_manager.go
## About This Procedure
* Describes the areas of modification in the Kubeflow Pipelines Backend to run the Custom Kubeflow SDK.

## Changes Made
* Edit `resource_manager.go` to remove the setting process of the Namespace to match the settings on the OpenKasugai controller and delete the processing of the Namespace assumed by Kubeflow.

## Prerequisites
* Make sure the following procedures have been completed before implementing this procedure.
    * [Modify visitor.go](../modify-visitor.go)

## Procedure
1. Open `resource_manager.go` in a text editor.
```
$ vi backend/src/apiserver/resource/resource_manager.go
```

2. Remove the Namespace setting process.
```go
	if k8sNamespace == "" {
		return nil, util.NewInternalServerError(util.NewInvalidInputError("Namespace cannot be empty when creating an Argo workflow. Check if you have specified POD_NAMESPACE or try adding the parent namespace to the request"), "Failed to create a run due to empty namespace")
	}
-	executionSpec.SetExecutionNamespace(k8sNamespace)
	newExecSpec, err := r.getWorkflowClient(k8sNamespace).Create(ctx, executionSpec, v1.CreateOptions{})
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
```

## Next Steps
[Modify v2_template.go](../modify-v2_template.go)
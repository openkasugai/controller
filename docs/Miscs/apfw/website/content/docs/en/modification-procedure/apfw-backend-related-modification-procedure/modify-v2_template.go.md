---
weight: 5
title: "v2_template.go Modification"
---
# Modify v2_template.go
## About This Procedure
* Describes the areas of modification in the Kubeflow Pipelines Backend when running the Custom Kubeflow SDK.

## Changes Made
* Edit `v2_template.go` to remove unused imports and functions.
* Change the setting value of `setDefaultServiceAccount` to `default` to set the service account for deploying custom resources from the Custom Kubeflow Backend to the `kubeflow` Namespace as `default`.

## Prerequisites
* Ensure the following steps have been completed before performing this procedure.
    * [Modify resource_manager.go](../modify-resource_manager.go)

## Procedure
1. Open `v2_template.go` in a text editor.
```
$ vi backend/src/apiserver/template/v2_template.go
```

2. Remove unnecessary imports.
```go
	"errors"
	"fmt"
	"io"

-	"regexp"
	"strings"

	structpb "github.com/golang/protobuf/ptypes/struct"
```

3. Delete unnecessary functions.
```go
			if err != nil {
				return nil, util.NewInvalidInputErrorWithDetails(ErrorInvalidPipelineSpec, fmt.Sprintf("invalid v2 pipeline spec: %s", err.Error()))
			}
-			if spec.GetSchemaVersion() != SCHEMA_VERSION_2_1_0 {
-				return nil, util.NewInvalidInputErrorWithDetails(ErrorInvalidPipelineSpec, fmt.Sprintf("KFP only supports schema version 2.1.0, but the pipeline spec has version %s", spec.
-			}
~~OMITTED~~
-			if spec.GetRoot() == nil {
-				return nil, util.NewInvalidInputErrorWithDetails(ErrorInvalidPipelineSpec, "invalid v2 pipeline spec: root component is empty")
-			}
			v2Spec.spec = &spec
		} else if IsPlatformSpecWithKubernetesConfig(valueBytes) {
			// Pick out the yaml document with platform spec
```

4. Align the permission of the `default` account in the `kubeflow` Namespace to create resources in each Namespace as set in [Create Namespace and Set Role](../../../deployment-procedure/configuration-on-dci-controller-node/create-namespace-and-set-role), and change the setting value of `setDefaultServiceAccount` to `default`.
```go
	if modelRun.Namespace != "" {
		executionSpec.SetExecutionNamespace(modelRun.Namespace)
	}
-	setDefaultServiceAccount(executionSpec, modelRun.ServiceAccount)
+	setDefaultServiceAccount(executionSpec, "default")
	// Disable istio sidecar injection if not specified
	executionSpec.SetAnnotationsToAllTemplatesIfKeyNotExist(util.AnnotationKeyIstioSidecarInject, util.AnnotationValueIstioSidecarInjectDisabled)
	// Add label to the workflow so it can be persisted by persistent agent later.
```

## Next Steps
[Modify install-go-licenses.sh](../modify-install-go-licenses.sh)
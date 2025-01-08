---
weight: 5
title: "v2_template.go改修"
---
# v2_template.go改修
## この手順について
* Custom Kubeflow SDKを実行するにあたりKubeflow Pipelines Backendの修正箇所について記載する。

## 対応内容
* `v2_template.go`を編集し、使用しないimportと関数を削除する。
* Custom Kubeflow Backendからカスタムリソースをデプロイする際のサービスアカウントを`kubeflow`Namespaceの`default`とするため、setDefaultServiceAccountの設定値を変更する。

## 前提条件
* 当該手順を実施する前に以下手順を実施済みであること。
    * [resource_manager.go改修](../modify-resource_manager.go)

## 手順
1. `v2_template.go`をテキストエディタで開く。
```
$ vi backend/src/apiserver/template/v2_template.go
```

2. 不要なimportを削除する。
```go
	"errors"
	"fmt"
	"io"

-	"regexp"
	"strings"

	structpb "github.com/golang/protobuf/ptypes/struct"
```

3. 不要な関数を削除する。
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

4. [Namespace作成とRole設定](../../../deployment-procedure/configuration-on-dci-controller-node/create-namespace-and-set-role)にて設定する、`kubeflow`Namespaceの`default`アカウントが各Namespaceにリソースを作成できる権限に合わせ、setDefaultServiceAccountの設定値を`default`に変更する。
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

## 次の手順について
[install-go-licenses.sh改修](../modify-install-go-licenses.sh)
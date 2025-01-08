---
weight: 8
title: "Modify Dockerfile"
---
# Modify Dockerfile
## About This Procedure
* Describes the modifications to the Kubeflow Pipelines Backend required to run the Custom Kubeflow SDK.

## Changes Made
* Edit the `Dockerfile` to exclude `github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec` from Go package license verification by `go-licenses`. (*1)

## Prerequisites
* Ensure the following steps have been completed before proceeding with this procedure.
    * [Modify Makefile](../modify-makefile)

## Procedure
1. Open the `Dockerfile` in a text editor.
```
$ vi backend/Dockerfile
```

2. Exclude `github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec` from Go package license verification.
```Dockerfile
 # Check licenses and comply with license terms.
 RUN ./hack/install-go-licenses.sh
 # First, make sure there's no forbidden license.
-RUN go-licenses check ./backend/src/apiserver
-RUN go-licenses csv ./backend/src/apiserver > /tmp/licenses.csv && \
+RUN go-licenses check --ignore github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec ./backend/src/apiserver
+RUN go-licenses csv --ignore github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec ./backend/src/apiserver > /tmp/licenses.csv && \
   diff /tmp/licenses.csv backend/third_party_licenses/apiserver.csv && \
-  go-licenses save ./backend/src/apiserver --save_path /tmp/NOTICES
+  go-licenses save --ignore github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec ./backend/src/apiserver 
--save_path /tmp/NOTICES
 
 # 2. Compile preloaded pipeline samples
 FROM python:3.7 as compiler
 ```

## Next Steps
[Delete Unnecessary Files](../delete-unnecessary-files)

## Notes
(*1) Due to the modifications in [Modify go.mod](../modify-go.mod) of this document, `go-licenses` excludes `github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec` from Go package license verification as the license information for it cannot be obtained.
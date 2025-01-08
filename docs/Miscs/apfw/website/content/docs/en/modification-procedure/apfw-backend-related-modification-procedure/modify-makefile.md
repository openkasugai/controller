---
weight: 7
title: "Modify Makefile"
---
# Modify Makefile
## About This Procedure
* Describes the areas of modification in the Kubeflow Pipelines Backend when running the Custom Kubeflow SDK.

## Changes Made
* Edit the Kubeflow Pipelines Backend `Makefile` to exclude `github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec` from Go package license verification by `go-licenses`. (*1)

## Prerequisites
* Ensure the following steps have been completed before performing this procedure.
    * [Modify install-go-licenses.sh](../modify-install-go-licenses.sh)

## Procedure
1. Edit the `Makefile`.
```
vi backend/Makefile
```

* Edit as follows.
```Makefile
# See README.md#updating-licenses-info section for more details.
.PHONY: license_apiserver
license_apiserver: $(BUILD)/apiserver
-       cd $(MOD_ROOT) && go-licenses csv ./backend/src/apiserver > $(CSV_PATH)/apiserver.csv
+       cd $(MOD_ROOT) && go-licenses csv --ignore github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec ./backend/src/apiserver > $(CSV_PATH)/apiserver.csv
.PHONY: license_persistence_agent
license_persistence_agent: $(BUILD)/persistence_agent
        cd $(MOD_ROOT) && go-licenses csv ./backend/src/agent/persistence > $(CSV_PATH)/persistence_agent.csv

```

## Next Steps
[Modify Dockerfile](../modify-dockerfile)

## Notes
(*1) Due to the [Modify go.mod](../modify-go.mod) in this document, the license information of `github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec` cannot be obtained, so it is excluded from Go package license verification by `go-licenses`.
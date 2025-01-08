---
weight: 9
title: "Delete Unnecessary Files"
---
# Delete Unnecessary Files
## About This Procedure
* Describes the areas of modification in the Kubeflow Pipelines Backend to run the Custom Kubeflow SDK.

## Changes Made
* Delete unnecessary files that are not used.

## Prerequisites
* Before implementing this procedure, make sure the following steps have been completed as part of the development environment setup.
    * [Modify Dockerfile](../modify-dockerfile)

#### Delete the following files under `backend/src/v2`.
There are functions that are no longer used for the modification, and if the files are not deleted, build errors will occur.
* compiler/argocompiler/argo_test.go
* compiler/argocompiler/dag.go
* compiler/visitor_test.go
* component/importer_launcher.go
* component/launcher_v2_test.go
* component/launcher_v2.go

## Next Steps
[Add OpenKasugai Library](../../apfw-sdk-related-modification-procedure/add-dci-library)
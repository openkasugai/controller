---
weight: 3
title: "Custom Kubeflow Backend Build"
---
# Custom Kubeflow Backend Build
## About This Procedure
* Perform a build after modifying Kubeflow Pipelines Backend.
* Build the `api-server` container, which converts IR YAML format from Custom Kubeflow to manifests and creates a container image.

## Changes Made
* Build the `api-server` container image.

## Prerequisites
### About the Build Environment
* Refer to the following for the necessary software and versions for the build.
    * [Environment Information and Prerequisites](../../../environment-information-and-prerequisites)

### Pre-Procedure
* It is assumed that the following procedures have been completed before this procedure.
    * [Build pipeline_spec.proto (Backend)](../build-pipeline_spec.proto)
    * [Updating Licenses Info](../updating-licenses-info)

## Procedure
1. Build the container image.
* Under the pipelines folder, build the container image with the following command.
* Specify any `TAG`.
```
$ cd pipelines
$ docker build -t api-server:<TAG> -f backend/Dockerfile .
```

2. Confirm that the `api-server` container image has been created.
```
$ docker images
REPOSITORY   TAG   　　    IMAGE ID   　　CREATED         SIZE
api-server   <TAG>  <CONTAINER ID>    <ELAPSED TIME>  　　219MB
```

## Next Steps
* [Build pipeline_spec.proto (SDK)](../../apfw-sdk-related-build-procedure/build-pipeline_spec.proto)
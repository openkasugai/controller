---
weight: 1
title: "Build pipeline_spec.proto (Backend)"
---
# Build pipeline_spec.proto (Backend)
## About This Procedure
* To build the Custom Kubeflow backend, you need to build `pipeline_spec.proto` beforehand to generate code.
* `pipeline_spec.proto` contains the interface definition of IR YAML, an intermediate representation for communication between Custom Kubeflow SDK and Custom Kubeflow backend, which is shared between them.
* Since Custom Kubeflow SDK outputs in Python and Custom Kubeflow backend in Go, the build procedures are different.
* Therefore, this procedure describes how to build the necessary `pipeline_spec.proto` before building the Custom Kubeflow backend and prepare `pipeline_spec.pb.go`.

## Changes Made
* Build `pipeline_spec.proto` to generate code.
* Confirm that `pipeline_spec.pb.go` is generated at the specified path.

## Prerequisites
### About the Build Environment
* To build `pipeline_spec.proto` according to this procedure, `Docker` must be installed in the build environment.
* Refer to the following for version information of the installed `Docker`: [Environment Information and Prerequisites](../../../environment-information-and-prerequisites)

### Pre-Procedure
* It is assumed that the following steps have been completed before this procedure is carried out.
    * [Modify Structures](../../../modification-procedure/apfw-sdk-related-modification-procedure/modify-structures)

## Procedure
1. Specify the API version to compile to.
```
$ export API_VERSION="v2beta1"
```

2. Navigate to the `api` directory.
```
$ cd pipelines/api
```

3. Execute the `make` command as follows.
```
$ make clean-go golang
```

Confirm that the compiled file is created at the following path.
`pipelines/api/v2alpha1/go/pipelinespec/pipeline_spec.pb.go`

## Next Steps
* [Updating Licenses Info](../updating-licenses-info)
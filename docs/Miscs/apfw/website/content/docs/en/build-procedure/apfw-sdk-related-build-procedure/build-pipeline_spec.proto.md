---
weight: 1
title: "Build pipeline_spec.proto (SDK)"
---
# Build pipeline_spec.proto (SDK)
## About This Procedure
* Each item output in YAML during Custom Kubeflow SDK execution is managed in `pipeline_spec.proto`.
* `pipeline_spec.proto` contains the interface definition of the intermediate representation YAML exchanged between Custom Kubeflow SDK and Custom Kubeflow backend, which is shared between them.
* When changing YAML items during Custom Kubeflow SDK modification, it is necessary to modify and build `pipeline_spec.proto`.
* Since Custom Kubeflow SDK outputs in Python and Custom Kubeflow backend in Go, the build procedures are different.
* Therefore, this procedure describes the steps to build the necessary `pipeline_spec.proto` for Custom Kubeflow SDK and prepare `pipeline_spec.pb.py`.

## Changes Made
* Build `pipeline_spec.proto` and generate code.
* Confirm that `pipeline_spec_pb2.py` is generated at the specified path.

## Prerequisites
### About the Build Environment
* To build `pipeline_spec.proto` based on this procedure, the build environment must have `protobuf-compiler` installed.
* Refer to the installation information for `protobuf-compiler` version details here: [Environment Information and Prerequisites](../../../environment-information-and-prerequisites)

### Pre-Procedure
* It is assumed that the following steps have been completed before implementing this procedure.
    * [Build Custom Kubeflow Backend](../../apfw-backend-related-build-procedure/build-apfw-backend)

## Procedure
1. Navigate to the `api` directory.
```
$ cd pipelines/api
```

2. Execute the following `make` command.
```
$ make clean-python python
```

Ensure that the compiled file is created at the following path:
`pipelines/api/v2alpha1/python/kfp/pipeline_spec/pipeline_spec_pb2.py`

## Next Steps
* [Create Namespace and Set Role](../../../deployment-procedure/configuration-on-dci-controller-node/create-namespace-and-set-role)
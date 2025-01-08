---
weight: 1
title: "Custom Kubeflow SDK Installation Procedure"
---
# Custom Kubeflow SDK Installation Procedure
## About This Procedure
* Describing the installation procedure for Custom Kubeflow SDK.

## Changes Made
* Executing the installation of Custom Kubeflow SDK.
* Build `pipeline_spec.proto` and replace the generated `pipeline_spec_pb2.py` with the existing `pipeline_spec_pb2.py`.

## Prerequisites
### About the Execution Environment
* To execute the installation of Custom Kubeflow SDK based on this procedure, it is necessary to have a constructed virtual environment `venv`.
* Refer to the following for the construction procedure of the `venv` to be installed.   
[Environment Information and Prerequisites](../../../environment-information-and-prerequisites)

### Pre-Procedure
* It is assumed that the following procedures have been completed before implementing this procedure.
    * [Deployment Procedures](../../../deployment-procedure/)

## Procedure
1. Switch to the virtual environment.
* Set `VENV_DIR` to the folder name specified when constructing the `venv` virtual environment.
```
$ source <VENV_DIR>/bin/activate
```

2. Confirm that the folder name is displayed at the front of the prompt, indicating that you have switched to the Python virtual execution environment.
```
<VENV_DIR>user_name@host:$
```

3. Move to the Custom Kubeflow SDK directory.
* Specify the directory where you cloned the Kubeflow repository in [Modify go.mod](../../../modification-procedure/apfw-backend-related-modification-procedure/modify-go.mod).
```
$ cd <DIR>/pipelines/sdk/python/kfp/
```

4. Install the Custom Kubeflow SDK.
```
$ pip install .
```

5. Confirm that the Custom Kubeflow SDK has been installed.
```
$ pip list | grep kfp
```

* Confirm the output as follows.
```
kfp                      2.4.0
kfp-pipeline-spec        0.2.2
kfp-server-api           2.0.5
```

6. Build `pipeline_spec.proto` and replace the generated `pipeline_spec_pb2.py` with the existing `pipeline_spec_pb2.py`.
```
$ cp <DIR>/pipelines/api/v2alpha1/python/kfp/pipeline_spec/pipeline_spec_pb2.py <VENV_DIR>/lib64/python3.8/site-packages/kfp/pipeline_spec/pipeline_spec_pb2.py
```

## Next Steps
* [Operation Confirmation by IR YAML Output](../operation-confirmation-by-iryaml-output)
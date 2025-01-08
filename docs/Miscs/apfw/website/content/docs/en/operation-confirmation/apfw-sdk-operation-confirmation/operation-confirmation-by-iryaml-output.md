---
weight: 2
title: "Operation Confirmation by IR YAML Output"
---
# Operation Confirmation by IR YAML Output
## About This Procedure
* This document describes the procedure for confirming the operation of Custom Kubeflow SDK by outputting IR YAML in an environment without a OpenKasugai controller.

## Changes Made
* Run the Custom Kubeflow SDK.
* Confirm that IR YAML is generated.

## Prerequisites
### About the Execution Environment
* To run the Custom Kubeflow SDK based on this procedure, the execution environment must have `python3` installed and a `venv` virtual environment set up.
* For information on the version of `python3` to install and the procedure for setting up `venv`, refer to the following:
    * [Environment Information and Prerequisites](../../../environment-information-and-prerequisites)

### Pre-Procedure
* It is assumed that the following procedures have been completed before this procedure is carried out.
    * [Deployment Procedures](../../../deployment-procedure/)

## Procedure
1. Switch to the virtual environment.
* Set `VENV_DIR` to the folder name specified when setting up the `venv` virtual environment.
```
$ source <VENV_DIR>/bin/activate
```

2. Confirm that the folder name is displayed at the front of the prompt, indicating that you have switched to the Python virtual execution environment.
```
<VENV_DIR>user_name@host:$
```

3. Execute the pipeline script.
* Specify the name of the pipeline script to execute in pipelines.
* This document provides an example of running the [Pipeline Sample](../../../pipeline-sample).
```
$ python <pipelines>.py
```

After running the script, confirm that an IR YAML file is created in the directory where the command was executed.

4. Check the contents of the IR YAML output.
* Confirm that the values set in the pipeline are reflected in the IR YAML settings. The settings corresponding to each output item are as follows.
| Pipeline | IR YAML |
| -------- | ------- |
| Setting values of FunctionItem | Setting values of Items |
| Setting values of Function | Setting values of the specified component name in components |
| Metadata such as Dataflow name | Setting values of metadata |
| Setting values of create_connection_task | Setting values of dag.task |

## Next Steps
* [Operation Confirmation by Dataflow Deployment](../operation-confirmation-by-dataflow-deployment)
---
weight: 3
title: "Operation Confirmation by Dataflow Deployment"
---
# Operation Confirmation by Dataflow Deployment
## About This Procedure
* This document describes the procedure for confirming the operation by deploying Dataflow as part of the Custom Kubeflow SDK operation confirmation.

## Changes Made
* Running the Custom Kubeflow SDK.
* Confirming that Dataflow is deployed.

## Prerequisites
### About the Execution Environment
* To run the Custom Kubeflow SDK based on this procedure, the execution environment must have `python3` installed and a `venv` virtual environment set up.
* To deploy Dataflow with the Custom Kubeflow SDK, the OpenKasugai controller must be operational.
* For information on the version of `python3` to install, setting up `venv`, and documents for setting up the OpenKasugai controller, refer to the following:
    * [Environment Information and Prerequisites](../../../environment-information-and-prerequisites)

### Pre-Procedure
* It is assumed that the following steps have been completed before implementing this procedure.
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
* Specify the name of the pipeline script to execute in `pipelines`.
* This document provides an example of executing the [Pipeline Sample](../../../pipeline-sample).
```
$ python <pipelines>.py
```

4. Confirm that the Configmap is deployed.
```
$ kubectl get configmap -A | grep -e NAME -e func | grep -v ^default
```

5. Confirm that the FunctionKind is deployed.
```
$ kubectl get functionkind -A
```

6. Confirm that the FunctionChain is deployed.
```
$ kubectl get functionchain -A
```

7. Confirm that Dataflow is deployed.
```
kubectl get dataflow -A
```
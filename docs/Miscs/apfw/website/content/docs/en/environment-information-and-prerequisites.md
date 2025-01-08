---
weight: 4
title: "Environment Information and Prerequisites"
---
# Environment Information and Prerequisites
## Overview
* This section describes the necessary environment information and prerequisites to follow the procedures outlined in this document.
* While this document assumes the OpenKasugai controller is installed, it is possible to proceed with some steps even if it is not installed.
  * If the OpenKasugai controller is not installed, the following steps cannot be performed:
    * [Create Namespace and Set Role](../deployment-procedure/configuration-on-dci-controller-node/create-namespace-and-set-role)
    * [Deploy Custom Kubeflow Backend](../deployment-procedure/apfw-backend-related-deployment-procedure/deploy-apfw-backend)
    * [Operation Confirmation by Dataflow Deployment](../operation-confirmation/apfw-sdk-operation-confirmation/operation-confirmation-by-dataflow-deployment)
  * If preparing the OpenKasugai controller environment, refer to the [OpenKasugai Controller Documentation](https://github.com/openkasugai/controller/blob/main/README.md#documentations) for setup.

## Environment Information
* The environment is as follows:
| OS/Software         | Version (*1) | 
| ------------ | ---- | 
| Ubuntu         | 20.04.2 LTS | 

## Prerequisites
* Ensure the following software is installed:
| Software         | Version (*1) | Installation Procedure |
| ------------ | ---- | ---- | 
| Python         | 3.8.10 | https://wiki.python.org/moin/BeginnersGuide/Download |
| Go             | 1.20.14 | https://go.dev/doc/install |
| libprotoc         | 3.6.1 | https://grpc.io/docs/protoc-installation/ |
| Docker         | 24.0.5 | https://docs.docker.com/engine/install/ubuntu/ |
| Harbor         | 2.9.1 | https://goharbor.io/docs/2.9.0/install-config/ |
* Ensure the venv virtual environment is set up.
| Software         | Setup Procedure | 
| ------------ | ---- |
| venv       | https://docs.python.org/3.8/library/venv.html |

## Next Steps
[Modify go.mod](../modification-procedure/apfw-backend-related-modification-procedure/modify-go.mod)

## Notes
(*1) The versions of the OS and software used at the time of document creation are listed.
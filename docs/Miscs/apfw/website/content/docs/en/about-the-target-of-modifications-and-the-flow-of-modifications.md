---
weight: 3
title: "About the Target of Modifications and the Flow of Modifications"
---
# About the Target of Modifications and the Flow of Modifications
## Overview
* Describing the target of modifications and the flow of modifications.

## Target of Modifications
* The target of modifications in this document is the `Kubeflow Pipelines` project with the `sdk-2.4.0` tag.

## Purpose of Modifications
* Using Kubeflow to make custom resources required for OpenKasugai Controller DataFlow execution deployable.
* To align with the specifications of custom resources required by OpenKasugai Controller, modify Kubeflow Pipelines Backend and Kubeflow Pipelines SDK.

## Configuration Diagram
![apfw arch](/images/apfw-dci_arch.jpg)
* The [Pipeline Sample](../pipeline-sample) used as a sample in this document corresponds to the `Python Pipeline` in the diagram.
* The `Python Pipeline` can represent information to deploy FunctionInfo (Configmap), FunctionKind, FunctionChain, and Dataflow in Python code.
* By modifying Kubeflow, the IR YAML format is changed to make it possible to deploy Dataflow and FunctionChain of OpenKasugai custom resources.

## Details of Modifications
* The details of modifications are as follows:
* For differences in resources due to version differences of OpenKasugai Controller, annotations are included for the parts where differences occur in the modifications.

### Kubeflow Pipelines Backend
| Item | Modifications |
| ---- | ------------- |
| [Modify go.mod](../modification-procedure/apfw-backend-related-modification-procedure/modify-go.mod) | Change the reference of the code at build time to a local file on the container. |
| [Modify argo.go](../modification-procedure/apfw-backend-related-modification-procedure/modify-argo.go) | Modify according to each resource to be deployed with Custom Kubeflow SDK. |
| [Modify visitor.go](../modification-procedure/apfw-backend-related-modification-procedure/modify-visitor.go) | Remove unused imports and functions. |
| [Modify resource_manager.go](../modification-procedure/apfw-backend-related-modification-procedure/modify-resource_manager.go) | Remove namespace configuration processing. |
| [Modify v2_template.go](../modification-procedure/apfw-backend-related-modification-procedure/modify-v2_template.go) | Delete unused imports and functions, and change the setting value of setDefaultServiceAccount. |
| [Modify Dockerfile](../modification-procedure/apfw-backend-related-modification-procedure/modify-dockerfile) | Exclude Go packages that can no longer obtain license information due to modifications from license verification. |
| [Delete Unnecessary Files](../modification-procedure/apfw-backend-related-modification-procedure/delete-unnecessary-files) | Delete unnecessary files that are not used. |

### Kubeflow Pipelines SDK
| Item | Modifications |
| ---- | ------------- |
| [Add OpenKasugai Library](../modification-procedure/apfw-sdk-related-modification-procedure/add-dci-library) | Create a new library for OpenKasugai Controller. (*1) |
| [Add Requirements](../modification-procedure/apfw-sdk-related-modification-procedure/add-requirements) | Create `requirements.in` and `requirements.txt` to reference the `kfp-pipeline-spec` modified in the steps of this document. |
| [Modify pipeline_spec.proto](../modification-procedure/apfw-sdk-related-modification-procedure/modify-pipeline_spec.proto) | Modify the data structure to match the format of IR YAML generated at Custom Kubeflow SDK runtime. |
| [Modify Components](../modification-procedure/apfw-sdk-related-modification-procedure/modify-components) | Add methods according to the component definition method of Custom Kubeflow SDK and add a process to create a pipeline_spec for OpenKasugai Controller environment. |
| [Modify Pipeline](../modification-procedure/apfw-sdk-related-modification-procedure/modify-pipeline) | Modify according to the specifications of the modified pipeline_spec for OpenKasugai Controller environment. |
| [Modify Compiler](../modification-procedure/apfw-sdk-related-modification-procedure/modify-compiler) | Modify according to the specifications of the modified pipeline_spec for OpenKasugai Controller environment. |
| [Modify Structures](../modification-procedure/apfw-sdk-related-modification-procedure/modify-structures) | Modify according to the specifications of the modified pipeline_spec for OpenKasugai Controller environment. |

## Flow of Modifications
* In this document, modifications are made and operational confirmation is performed in the following flow.
1. Modify Kubeflow Pipelines Backend and Kubeflow Pipelines SDK to create Custom Kubeflow Backend and Custom Kubeflow SDK.
2. Build pipeline_spec.proto and Custom Kubeflow Backend.
3. Configure OpenKasugai Controller and deploy Custom Kubeflow Backend.
4. Execute Custom Kubeflow SDK and perform operational confirmation.

## Next Steps
[Environment Information and Prerequisites](../environment-information-and-prerequisites)

## Notes
(*1) The library name `swb` in the library addition procedure refers to OpenKasugai.
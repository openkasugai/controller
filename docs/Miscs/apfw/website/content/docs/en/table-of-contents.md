---
weight: 1
title: "Table of Contents"
---
## Table of Contents

- [Overview](../overview)
- Table of Contents
- [About the Target of Modifications and the Flow of Modifications](../about-the-target-of-modifications-and-the-flow-of-modifications)
- [Environment Information and Prerequisites](../environment-information-and-prerequisites)
- Modification Procedure
    - Custom Kubeflow Backend Related Modification Procedure
        - [Modify go.mod](../modification-procedure/apfw-backend-related-modification-procedure/modify-go.mod)
        - [Modify argo.go](../modification-procedure/apfw-backend-related-modification-procedure/modify-argo.go)
        - [Modify visitor.go](../modification-procedure/apfw-backend-related-modification-procedure/modify-visitor.go)
        - [Modify resource_manager.go](../modification-procedure/apfw-backend-related-modification-procedure/modify-resource_manager.go)
        - [Modify v2_template.go](../modification-procedure/apfw-backend-related-modification-procedure/modify-v2_template.go)
        - [Modify install-go-licenses.sh](../modification-procedure/apfw-backend-related-modification-procedure/modify-install-go-licenses.sh)
        - [Modify Makefile](../modification-procedure/apfw-backend-related-modification-procedure/modify-makefile)
        - [Modify Dockerfile](../modification-procedure/apfw-backend-related-modification-procedure/modify-dockerfile)
        - [Delete Unnecessary Files](../modification-procedure/apfw-backend-related-modification-procedure/delete-unnecessary-files)
    - Custom Kubeflow SDK Related Modification Procedure
        - [Add OpenKasugai Library](../modification-procedure/apfw-sdk-related-modification-procedure/add-dci-library)
        - [Add Requirements](../modification-procedure/apfw-sdk-related-modification-procedure/add-requirements)
        - [Modify pipeline_spec.proto](../modification-procedure/apfw-sdk-related-modification-procedure/modify-pipeline_spec.proto)
        - [Modify Components](../modification-procedure/apfw-sdk-related-modification-procedure/modify-components)
        - [Modify Pipeline](../modification-procedure/apfw-sdk-related-modification-procedure/modify-pipeline)
        - [Modify Compiler](../modification-procedure/apfw-sdk-related-modification-procedure/modify-compiler)
        - [Modify Structures](../modification-procedure/apfw-sdk-related-modification-procedure/modify-structures)
- Build Procedure
    - Custom Kubeflow Backend Related Build Procedure
        - [Build pipeline_spec.proto (Backend)](../build-procedure/apfw-backend-related-build-procedure/build-pipeline_spec.proto)
        - [Updating Licenses Info](../build-procedure/apfw-backend-related-build-procedure/updating-licenses-info)
        - [Build Custom Kubeflow Backend](../build-procedure/apfw-backend-related-build-procedure/build-apfw-backend)
    - Custom Kubeflow SDK Related Build Procedure
        - [Build pipeline_spec.proto (SDK)](../build-procedure/apfw-sdk-related-build-procedure/build-pipeline_spec.proto)
- Deployment Procedure
    - Configuration on OpenKasugai Controller Node
        - [Create Namespace and Set Role](../deployment-procedure/configuration-on-dci-controller-node/create-namespace-and-set-role)
    - Custom Kubeflow Backend Related Deployment Procedure
        - [Deploy Custom Kubeflow Backend](../deployment-procedure/apfw-backend-related-deployment-procedure/deploy-apfw-backend)
- Operation Confirmation
    - Custom Kubeflow SDK Operation Confirmation
        - [Custom Kubeflow SDK Installation Procedure](../operation-confirmation/apfw-sdk-operation-confirmation/apfw-sdk-install)
        - [Operation Confirmation by IR YAML Output](../operation-confirmation/apfw-sdk-operation-confirmation/operation-confirmation-by-iryaml-output)
        - [Operation Confirmation by Dataflow Deployment](../operation-confirmation/apfw-sdk-operation-confirmation/operation-confirmation-by-dataflow-deployment)
- [Pipeline Sample](../pipeline-sample)
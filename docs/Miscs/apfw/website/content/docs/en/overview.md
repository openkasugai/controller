---
weight: 2
title: "Overview"
---
# Overview
## About This Document
* This is a guide for running inference/machine learning Pipelines (referred to as "Pipeline" here) on OpenKasugai controllers.
* This document provides an example of using [Kubeflow](https://www.kubeflow.org/) as one of the methods to execute Pipelines. (*1)
* To integrate Kubeflow with OpenKasugai controllers, you need to rewrite Kubeflow's source code to match the specifications of OpenKasugai controllers.
* Therefore, this document describes the following:
  * Areas of Kubeflow that need to be modified to match the specifications of OpenKasugai controllers and how to make those modifications.
  * How to execute and verify Kubeflow after modification.

## Target Audience
* Users of OpenKasugai controllers considering Pipeline execution.
* Developers of OpenKasugai controllers who want to understand the development content for executing Pipelines.

## What You Can Achieve with This Document
* Modification of Kubeflow to match the specifications of OpenKasugai controllers (*3)
* Execution of Pipelines and deployment of resources (Dataflow) for OpenKasugai controllers (*4)
* Execution of modified Kubeflow as a standalone and verification of proper operation using intermediate representation files.
* Application of modification steps to different versions of OpenKasugai controllers based on the above procedures (*5)

## Benefits of Setting Up the Pipeline Execution Environment
* Setting up the Pipeline execution environment based on this document has the following benefits:
  * Since Pipelines are described in Python, creating Pipelines is easy with Python knowledge.
  * By providing minimal information, resources can be deployed easily, reducing the learning cost for creating resources for OpenKasugai controllers.
  * It becomes possible to remotely execute Pipelines and create resources for OpenKasugai controllers via API.

## Term Definitions
The following terms will be used in the subsequent chapters:
| Term         | Definition | 
| ------------ | ---- | 
| [Kubeflow Pipelines SDK](https://www.kubeflow.org/docs/components/pipelines/)      | The original Kubeflow Pipeline SDK before modification. | 
| Custom Kubeflow SDK      | The Kubeflow Pipelines SDK being modified or already modified to match the specifications of OpenKasugai controllers. | 
| [Kubeflow Pipelines Backend](https://github.com/kubeflow/pipelines/tree/sdk-2.4.0/backend) | Components and code to launch Kubeflow on Kubernetes. | 
| Custom Kubeflow Backend | The Kubeflow Pipelines Backend being modified or already modified to match the specifications of OpenKasugai controllers. | 

## Next Steps
[About the Target of Modifications and the Flow of Modifications](../about-the-target-of-modifications-and-the-flow-of-modifications)

## Notes
(*1) Reasons for choosing Kubeflow include the following:
1. Environment Compatibility: It can be used on-premises and on various clouds, and it is highly compatible with Kubernetes.
2. Input/Output Ease: It allows you to describe any flow in Python and define flows from a WebUI.
3. Update Frequency: It has a relatively high community activity with around 10 updates per year.

(*2) The version of the OpenKasugai controller used for testing in this document is generally not disclosed, and the FPGA circuit and settings used for image processing are also not disclosed.   
(*3) 
* If you make modifications to Kubeflow as described in this document, it will not be possible to execute Pipelines that match the specifications of Kubeflow before modification.
* Due to the different IR YAML formats between Kubeflow and OpenKasugai controllers deploying custom resources, it is difficult to enable existing Pipeline execution without affecting existing operations, so no modification approach has been taken to avoid this.

(*4) The [Pipeline sample](../pipeline-sample) used as a sample in this document deploys Dataflow and creates an inference Pod. The inference process in the Pod uses a pre-built learning model.  
(*5) The purpose is to enable Pipeline execution on OpenKasugai controllers with different versions by publishing modification details for older versions of OpenKasugai controllers.
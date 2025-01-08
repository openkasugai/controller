---
weight: 3
title: "Modify pipeline_spec.proto"
---
# Modify pipeline_spec.proto
## About This Procedure
* Describes the areas of modification in the Kubeflow Pipelines SDK required to run the Custom Kubeflow SDK.

## Changes Made
* Edit `pipeline_spec.proto` to modify the structure of the pipeline_spec used by the Kfp Client for OpenKasugai. (*1)

## Prerequisites
* Ensure the following steps have been completed before proceeding with this procedure.
    * [Add Requirements](../add-requirements)

## Procedure
1. Open `pipeline_spec.proto` in a text editor.
```
$ vi api/v2alpha1/pipeline_spec.proto
```

2. Edit as follows.
```go
// The spec of a pipeline.
message PipelineSpec {
- // The metadata of the pipeline.
- PipelineInfo pipeline_info = 1;
-
- // The deployment config of the pipeline.
- // The deployment config can be extended to provide platform specific configs.
- google.protobuf.Struct deployment_spec = 7;
+ MetaSpec meta = 1;
+
+ // The map of name to definition of all components used in this pipeline.
+ map<string, ComponentSpec> components = 2;
+
+ // The definition of the main pipeline.  Execution of the pipeline is
+ // completed upon the completion of this component.
+ RootSpec root = 3;

  // The version of the sdk, which compiles the spec.
  string sdk_version = 4;

  // The version of the schema.
  string schema_version = 5;

- // The definition of the runtime parameter.
- message RuntimeParameter {
-   // Required field. The type of the runtime parameter.
-   PrimitiveType.PrimitiveTypeEnum type = 1;
-   // Optional field. Default value of the runtime parameter. If not set and
-   // the runtime parameter value is not provided during runtime, an error will
-   // be raised.
-   Value default_value = 2;
- }
-
- // The map of name to definition of all components used in this pipeline.
- map<string, ComponentSpec> components = 8;
-
- // The definition of the main pipeline.  Execution of the pipeline is
- // completed upon the completion of this component.
- ComponentSpec root = 9;
-
- // Optional field. The default root output directory of the pipeline.
- string default_pipeline_root = 10;
+ PipelineInfo pipeline_info = 6;
}

+message MetaSpec {
+ string dataflowName = 1;
+ string dataflowNamespace = 2;
+ string functionChainName = 3;
+ string functionChainNamespace = 4;
+ string functionKindNamespace = 5;
+ string connectionKindNamespace = 6;
+}

// Definition of a component.
message ComponentSpec {
- // Definition of the input parameters and artifacts of the component.
- ComponentInputsSpec input_definitions = 1;
- // Definition of the output parameters and artifacts of the component.
- ComponentOutputsSpec output_definitions = 2;
- // Either a DAG or a single execution.
- oneof implementation {
-   DagSpec dag = 3;
-   string executor_label = 4;
- }
-}
+ string namespace = 1;
+ map<string, interfaceSpec> info = 2;
+ string version = 3;
+
+ //Connection information.
+ //start
+ string startPointIP = 4;
+ int64 startPointPort = 5;
+ string startPointProtocol = 6;
+
+ //origin
+ int64 fromPort = 7;
+ int64 toPort = 8;
+ string connectionKind = 9;
+
+ //end
+ string endPointIP = 10;
+ int64 endPointPort = 11;
+ string endPointProtocol = 12;
+
+ message interfaceSpec{
+   map<string, ItemsSpec> items = 1;
+   message ItemsSpec{
+     string configName = 1;
+     int64 coreMin = 2;
+     int64 coreMax = 3;
+     int64 totalBase = 6;
+     int64 capacityTotalBase = 7;
+   }
+ }
+}
+
+message RootSpec {
+ DagSpec dag = 1;
+ ComponentInputsSpec input_definitions = 2;
+ ComponentOutputsSpec output_definitions = 3;
+}

// A DAG contains multiple tasks.
message DagSpec {
  // The tasks inside the dag.
  map<string, PipelineTaskSpec> tasks = 1;
-
- // Defines how the outputs of the dag are linked to the sub tasks.
- DagOutputsSpec outputs = 2;
}

// Definition of the output artifacts and parameters of the DAG component.
~~OMITTED~~
// The spec of a pipeline task.
message PipelineTaskSpec {
- // Basic info of a pipeline task.
- PipelineTaskInfo task_info = 1;
-
~~OMITTED~~
-   ArtifactIteratorSpec artifact_iterator = 9;
-   // Iterator to iterate over a parameter input.
-   ParameterIteratorSpec parameter_iterator = 10;
+ ComponentRef componentRef = 1;
+ repeated string dependentTasks = 2;
+
+ message ComponentRef {
+ string name = 1;
  }

  // User-configured task-level retry.
~~OMITTED~~
  // User-configured task-level retry.
  // Applicable only to component tasks.
  RetryPolicy retry_policy = 11;
-
- // Iterator related settings.
- message IteratorPolicy {
-   // The limit for the number of concurrent sub-tasks spawned by an iterator
-   // task. The value should be a non-negative integer. A value of 0 represents
-   // unconstrained parallelism.
-   int32 parallelism_limit = 1;
- }
-
- // Iterator related settings.
- IteratorPolicy iterator_policy = 12;
}

// The spec of an artifact iterator. It supports fan-out a workflow from a list
```

## Next Steps
[Modify Components](../modify-components)

## Notes
(*1) Depending on the version differences of the OpenKasugai controller, modifications must be made to align with the specifications.
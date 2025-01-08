---
weight: 5
title: "Modify Pipeline"
---
# Modify Pipeline
## About This Procedure
* This document describes the areas of modification in the Kubeflow Pipelines SDK to run the Custom Kubeflow SDK.

## Procedure
* Edit `pipeline_context.py` and `pipeline_task.py` to align with the modified pipeline_spec specifications for the OpenKasugai controller environment. (*1)

## Prerequisites
* Ensure the following steps have been completed before executing this procedure.
    * [Modify Components](../modify-components)

## Procedure
1. Open `pipeline_context.py` in a text editor.
```
$ vi sdk/python/kfp/dsl/pipeline_context.py
```

* Modify the arguments of the pipeline decorator.
```python
from kfp.dsl import tasks_group
from kfp.dsl import utils


def pipeline(func: Optional[Callable] = None,
             *,
-            name: Optional[str] = None,
-            description: Optional[str] = None,
-            pipeline_root: Optional[str] = None,
-            display_name: Optional[str] = None) -> Callable:
+            dataflowName: Optional[str] = None,
+            dataflowNamespace: Optional[str] = None,
+            functionChainName: Optional[str] = None,
+            functionChainNamespace: Optional[str] = None,
+            functionKindNamespace: Optional[str] = None,
+            connectionKindNamespace: Optional[str] = None,
+   ) -> Callable:
    """Decorator used to construct a pipeline.
    Example
```

* Modify the return value of the pipeline decorator.
* Remove pipeline_root as it specifies the path for input and output.
```python
    if func is None:
        return functools.partial(
            pipeline,
-           name=name,
-           description=description,
-           pipeline_root=pipeline_root,
-           display_name=display_name,
+           dataflowName=dataflowName,
+           dataflowNamespace=dataflowNamespace,
+           functionChainName=functionChainName,
+           functionChainNamespace=functionChainNamespace,
+           functionKindNamespace=functionKindNamespace,
+           connectionKindNamespace=connectionKindNamespace
        )

-   if pipeline_root:
-       func.pipeline_root = pipeline_root
-
    return component_factory.create_graph_component_from_func(
        func,
-       name=name,
-       description=description,
-       display_name=display_name,
+       dataflowName=dataflowName,
+       dataflowNamespace=dataflowNamespace,
+       functionChainName=functionChainName,
+       functionChainNamespace=functionChainNamespace,
+       functionKindNamespace=functionKindNamespace,
+       connectionKindNamespace=connectionKindNamespace
    )
```

2. Open `pipeline_task.py` in a text editor.
```
$ vi sdk/python/kfp/dsl/pipeline_task.py
```

* Add `connection_spec`.
```python
        self.importer_spec = None
        self.container_spec = None
        self.pipeline_spec = None
+       self.connection_spec = None
        self._ignore_upstream_failure_tag = False
        # platform_config for this primitive task; empty if task is for a graph component
        self.platform_config = {}
```

* In the definition of the modified component, the `implementation` specifying how to execute the component instance has been removed, so remove related processing.
```python
                (component_spec.implementation.container.args or [])):
                check_primitive_placeholder_is_used_for_correct_io_type(
                    inputs_dict, outputs_dict, arg)

-       if component_spec.implementation.container is not None:
-           validate_placeholder_types(component_spec)
-           self.container_spec = self._extract_container_spec_and_convert_placeholders(
-               component_spec=component_spec)
-       elif component_spec.implementation.importer is not None:
-           self.importer_spec = component_spec.implementation.importer
-           self.importer_spec.artifact_uri = args['uri']
-       else:
-           self.pipeline_spec = self.component_spec.implementation.graph
-
        self._outputs = {
            output_name: pipeline_channel.create_pipeline_channel(
```

## Next Steps
[Modify Compiler](../modify-compiler)

## Notes
(*1) Depending on the version differences of the OpenKasugai controller, it may be necessary to make modifications according to the specifications if they differ.
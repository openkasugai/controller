---
weight: 6
title: "Modify Compiler"
---
# Modify Compiler
## About This Procedure
* Describes the modifications to the Kubeflow Pipelines SDK required to run the Custom Kubeflow SDK.

## Changes Made
* Edit `compiler.py` and `pipeline_spec_builder.py` to align with the modified pipeline_spec for the OpenKasugai controller environment. (*1)

## Prerequisites
* Ensure the following steps have been completed before proceeding with this procedure.
    * [Modify Pipeline](../modify-pipeline)

## Procedure
1. Open `compiler.py` in a text editor.
```
$ vi sdk/python/kfp/compiler/compiler.py
```

* Remove unnecessary parts to match the format of the modified IR YAML.
```python

            builder.write_pipeline_spec_to_file(
                pipeline_spec=pipeline_spec,
-               pipeline_description=pipeline_func.description,
-               platform_spec=pipeline_func.platform_spec,
                package_path=package_path,
            )
```

2. Open `pipeline_spec_builder.py` in a text editor.
```
$ vi sdk/python/kfp/compiler/pipeline_spec_builder.py
```

* Remove related processes that are no longer needed to align with the changed structure of the IR YAML due to the pipeline_spec modification.
```python
        ]
        is_parent_component_root = (group_component_spec == pipeline_spec.root)

-       if isinstance(subgroup, pipeline_task.PipelineTask):
-           subgroup_task_spec = build_task_spec_for_task(
-               task=subgroup,
~~OMITTED~~
-   pipeline_spec.deployment_spec.update(
-       json_format.MessageToDict(deployment_config))
-
    # Surface metrics outputs to the top.
    populate_metrics_in_dag_outputs(
```

* Modify to match the new structure of pipeline_spec and dag.
```python
def create_pipeline_spec(
-   pipeline: pipeline_context.Pipeline,
+   pipeline: pipeline_context,
    component_spec: structures.ComponentSpec,
+   meta_spec: pipeline_spec_pb2.MetaSpec,
    pipeline_outputs: Optional[Any] = None,
) -> Tuple[pipeline_spec_pb2.PipelineSpec, pipeline_spec_pb2.PlatformSpec]:
    """Creates a pipeline spec object.
```

* Adjust to pass values of `meta` and `components` to `pipeline_spec`.
```python
    # Schema version 2.1.0 is required for kfp-pipeline-spec>0.1.13
    pipeline_spec.schema_version = '2.1.0'

-   pipeline_spec.root.CopyFrom(
-       _build_component_spec_from_component_spec_structure(component_spec))
+   pipeline_spec.meta.CopyFrom(meta_spec)
+
+   components_tmp = {"components":{}}
+   for task_k,task_v in pipeline.tasks.items():
+        components_tmp['components'].update(task_v.component_spec.component["components"])
+        del components_tmp['components']['inputDefinitions']
+        del components_tmp['components']['outputDefinitions']
+
+   json_string = json.dumps(components_tmp)
+   json_format.Parse(json_string, pipeline_spec)

    # TODO: add validation of returned outputs -- it's possible to return
    # an output from a task in a condition group, for example, which isn't
```

* Remove `platform_spec` in the modified IR YAML format as it is unnecessary.
* Modify the pipeline generation process.
```python
        condition_channels=condition_channels,
    )

-   platform_spec = pipeline_spec_pb2.PlatformSpec()
-   for group in all_groups:
-       build_spec_by_group(
~~OMITTED~~
-       structures_component_spec=component_spec)
-
-   return pipeline_spec, platform_spec
+   for group in all_groups:
+       for task in group.tasks:
+           pipeline_spec.root.dag.tasks[task.name].componentRef.name = task.name
+           for dependentTask in task._run_after:
+               pipeline_spec.root.dag.tasks[task.name].dependentTasks.append(dependentTask)
+   return pipeline_spec


def _validate_dag_output_types(
~~OMITTED~~

def write_pipeline_spec_to_file(
    pipeline_spec: pipeline_spec_pb2.PipelineSpec,
-   pipeline_description: Union[str, None],
-   platform_spec: pipeline_spec_pb2.PlatformSpec,
    package_path: str,
) -> None:
    """Writes PipelineSpec into a YAML or JSON (deprecated) file.
~~OMITTED~~
        platform_spec: The PlatformSpec.
    """
    pipeline_spec_dict = json_format.MessageToDict(pipeline_spec)
-   yaml_comments = extract_comments_from_pipeline_spec(pipeline_spec_dict,
-                                                       pipeline_description)
-   has_platform_specific_features = len(platform_spec.platforms) > 0

    if package_path.endswith('.json'):
        warnings.warn(
~~OMITTED~~
            stacklevel=2,
        )
        with open(package_path, 'w') as json_file:
-           if has_platform_specific_features:
-               raise ValueError(
-                   f'Platform-specific features are only supported when serializing to YAML. Argument for {"package_path"!r} has file extension {".json"!r}.'
-               )
            json.dump(pipeline_spec_dict, json_file, indent=2, sort_keys=True)

    elif package_path.endswith(('.yaml', '.yml')):
        with open(package_path, 'w') as yaml_file:
-           yaml_file.write(yaml_comments)
            documents = [pipeline_spec_dict]
-           if has_platform_specific_features:
-               documents.append(json_format.MessageToDict(platform_spec))
            yaml.dump_all(documents, yaml_file, sort_keys=True)

    else:
```

## Next Steps
[Modify Structures](../modify-structures)

## Notes
(*1) If the specifications differ due to version differences in the OpenKasugai controller, modifications to align with the specifications are necessary.
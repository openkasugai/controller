---
weight: 6
title: "Compiler改修"
---
# Compiler改修
## この手順について
* Custom Kubeflow SDKを実行するにあたり、Kubeflow Pipelines SDKの修正箇所について記載する。

## 対応内容
* `compiler.py`と`pipeline_spec_builder.py`を編集し、OpenKasugaiコントローラ環境用に改修したpipeline_specの仕様に合わせる。(※1) 

## 前提条件
* 当該手順を実施する前に以下手順を実施済みであること。
    * [Pipeline改修](../modify-pipeline)

## 手順
1. `compiler.py`をテキストエディタで開く。
```
$ vi sdk/python/kfp/compiler/compiler.py
```

* 改修後のIR YAMLのフォーマットにあわせ、不要箇所を削除する。
```python

            builder.write_pipeline_spec_to_file(
                pipeline_spec=pipeline_spec,
-               pipeline_description=pipeline_func.description,
-               platform_spec=pipeline_func.platform_spec,
                package_path=package_path,
            )
```

2. `pipeline_spec_builder.py`をテキストエディタで開く。
```
$ vi sdk/python/kfp/compiler/pipeline_spec_builder.py
```

* pipeline_specの改修により変わったIR YAMLの構造に合わせて、不要となった関連処理を削除
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

* 改修後のpipeline_spec型とdag型にあわせるため修正する。
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

* `meta`と`components`の値を`pipeline_spec`に渡すよう改修する。
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

* 改修後のIR YAMLのフォーマットには`platform_spec`は不要のため削除する。
* `pipeline`生成処理を改修する。
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

## 次の手順について
[structures改修](../modify-structures)

## 補足事項
(※1) OpenKasugaiコントローラのバージョン差分により仕様が異なる場合は、仕様に合わせた改修を実施する必要がある。
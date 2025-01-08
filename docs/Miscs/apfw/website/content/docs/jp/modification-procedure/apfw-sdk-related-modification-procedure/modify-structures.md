---
weight: 7
title: "structures改修"
---
# structures改修
## この手順について
* Custom Kubeflow SDKを実行するにあたり、Kubeflow Pipelines SDKの修正箇所について記載する。

## 対応内容
* `structures.py`を編集し、OpenKasugaiコントローラ環境用に[pipeline_spec.proto改修](../modify-pipeline_spec.proto)で改修したpipeline_specの仕様に合わせる。(※1) 

## 前提条件
* 当該手順を実施する前に以下手順を実施済みであること。
    * [Compiler改修](../modify-compiler)

## 手順
1. `structures.py`をテキストエディタで開く。
```
$ vi sdk/python/kfp/dsl/structures.py
```

* 改修後のIR YAMLのフォーマットにあわせ、不要箇所を削除する。
```python
@dataclasses.dataclass
class ComponentSpec:
    """The definition of a component.
~~OMITTED~~
            (container, importer) or a DAG consists of other components.
    """
    name: str
-   implementation: Implementation
-   description: Optional[str] = None
+   component: Optional[Dict[str, Any]] = None
    inputs: Optional[Dict[str, InputSpec]] = None
    outputs: Optional[Dict[str, OutputSpec]] = None
-   platform_spec: pipeline_spec_pb2.PlatformSpec = dataclasses.field(
-       default_factory=pipeline_spec_pb2.PlatformSpec)
-
~~OMITTED~~
-           inputs=inputs,
-           outputs=outputs,
-       )

    @classmethod
    def from_ir_dicts(
        cls,
        pipeline_spec_dict: dict,
-       platform_spec_dict: dict,
    ) -> 'ComponentSpec':
        """Creates a ComponentSpec from the PipelineSpec and PlatformSpec
        messages as dicts."""
~~OMITTED~~
                                return docstring
            return None

-       component_key = utils.sanitize_component_name(raw_name)
-       component_spec_dict = pipeline_spec_dict['components'].get(
-           component_key, pipeline_spec_dict['root'])
+       component_spec_dict = pipeline_spec_dict['components']

        inputs = inputs_dict_from_component_spec_dict(component_spec_dict)
        outputs = outputs_dict_from_component_spec_dict(component_spec_dict)

-       implementation = Implementation.from_pipeline_spec_dict(
-           pipeline_spec_dict, raw_name)
-
-       description = extract_description_from_command(
-           implementation.container.command or
-           []) if implementation.container else None
-
-       platform_spec = pipeline_spec_pb2.PlatformSpec()
-       json_format.ParseDict(platform_spec_dict, platform_spec)
-
        return ComponentSpec(
            name=raw_name,
-           implementation=implementation,
-           description=description,
+           component=pipeline_spec_dict,
            inputs=inputs,
            outputs=outputs,
-           platform_spec=platform_spec,
        )

    @classmethod
~~OMITTED~~
                component_yaml)
            return cls.from_v1_component_spec(v1_component)
        else:
-           component_spec = ComponentSpec.from_ir_dicts(
-               pipeline_spec_dict, platform_spec_dict)
(pipeline_spec_dict)
-           if not component_spec.description:
-               component_spec.description = extract_description(
-                   component_yaml=component_yaml)
+           component_spec = ComponentSpec.from_ir_dicts
            return component_spec

    def save_to_component_yaml(self, output_file: str) -> None:
```

## 次の手順について
[pipeline_spec.protoビルド(バックエンド)](../../../build-procedure/apfw-backend-related-build-procedure/build-pipeline_spec.proto)

## 補足事項
(※1) OpenKasugaiコントローラのバージョン差分により仕様が異なる場合は、仕様に合わせた改修を実施する必要がある。
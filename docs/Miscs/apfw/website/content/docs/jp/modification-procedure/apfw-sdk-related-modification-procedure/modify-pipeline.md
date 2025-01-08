---
weight: 5
title: "Pipeline改修"
---
# Pipeline改修
## この手順について
* Custom Kubeflow SDKを実行するにあたり、Kubeflow Pipelines SDKの修正箇所について記載する。

## 対応内容
* `pipeline_context.py`と`pipeline_task.py`を編集し、OpenKasugaiコントローラ環境用に改修したpipeline_specの仕様に合わせる。(※1)

## 前提条件
* 当該手順を実施する前に以下手順を実施済みであること。
    * [Components改修](../modify-components)

## 手順
1. `pipeline_context.py`をテキストエディタで開く。
```
$ vi sdk/python/kfp/dsl/pipeline_context.py
```

* pipelineデコレータの引数を変更する。

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

* pipelineデコレータの戻り値を変更する。
* pipeline_rootは入出力先のパスを指定する引数のため削除する。

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

2. `pipeline_task.py`をテキストエディタで開く。
```
$ vi sdk/python/kfp/dsl/pipeline_task.py
```

* `connection_spec`を追加する。

```python
        self.importer_spec = None
        self.container_spec = None
        self.pipeline_spec = None
+       self.connection_spec = None
        self._ignore_upstream_failure_tag = False
        # platform_config for this primitive task; empty if task is for a graph component
        self.platform_config = {}
```

* 改修後のcomponentの定義ではcomponentインスタンスの実行方法を指定する`implemention`が削除されているため関連処理を削除する。
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

## 次の手順について
[Compiler改修](../modify-compiler)

## 補足事項
(※1) OpenKasugaiコントローラのバージョン差分により仕様が異なる場合は、仕様に合わせた改修を実施する必要がある。
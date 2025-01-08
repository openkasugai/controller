---
weight: 4
title: "Components改修"
---
# Components改修
## この手順について
* Custom Kubeflow SDKを実行するにあたり、Kubeflow Pipelines SDKの修正箇所について記載する。

## 対応内容
* `load_yaml_utilities.py`を編集し、[pipeline_spec.proto改修](../modify-pipeline_spec.proto)で改修したcomponent定義にあわせたメソッドを追加する。(※1) 
* `dsl`配下のcomponents関連のファイルを編集し、OpenKasugaiコントローラ環境用のpipeline_spec作成処理を追加する。(※1) 

## 前提条件
* 当該手順を実施する前に以下手順を実施済みであること。
    * [pipeline_spec.proto改修](../modify-pipeline_spec.proto)

## 手順
1. `__init__.py`をテキストエディタで開く。
```
$ vi sdk/python/kfp/components/__init__.py
```

* `load_component`を追加する。
```python
# limitations under the License.

__all__ = [
+   'load_component',
    'load_component_from_file',
    'load_component_from_url',
    'load_component_from_text',
    'PythonComponent',
    'BaseComponent',
    'ContainerComponent',
    'YamlComponent',
]

+from kfp.components.load_yaml_utilities import load_component
from kfp.components.load_yaml_utilities import load_component_from_file
from kfp.components.load_yaml_utilities import load_component_from_text
from kfp.components.load_yaml_utilities import load_component_from_url
```

2. `load_yaml_utilities.py`をテキストエディタで開く。
```
$ vi sdk/python/kfp/components/load_yaml_utilities.py
```

* compnents名に`_`が含まれていた場合に`-`に変換する処理を追記する。
* OpenKasugai環境用に改修したcomponentの定義にあわせてload_componentメソッドを作成する。
```python
from kfp.dsl import structures
from kfp.dsl import yaml_component
import requests
+import yaml


+def __valid_component_param(text: str):
+   comp_dict = yaml.safe_load(text)
+   for k in comp_dict["components"].copy():
+       if k != "inputDefinitions" and k != "outputDefinitions":
+           if '_' in k:
+               changed_k = k.replace('_', '-').strip()
+               comp_dict["components"][changed_k] = comp_dict["components"].pop(k)
+
+   return yaml.dump(comp_dict)
+
+def load_component(text: str) -> yaml_component.YamlComponent:
+   """Loads a component from text.
+   Args:
+       text (str): Component YAML text.
+   Returns:
+       Component loaded from YAML.
+   """
+
+   text = __valid_component_param(text=text)
+
+   return yaml_component.YamlComponent(
+       component_spec=structures.ComponentSpec.from_yaml_documents(text),
+       component_yaml=text)
+
def load_component_from_text(text: str) -> yaml_component.YamlComponent:
    """Loads a component from text.
```

3. `base_component.py`をテキストエディタで開く。
```
$ vi sdk/python/kfp/dsl/base_component.py
```

* pipelineSpec型の戻り値を取得するために初期化処理を追加する。
```python
        """
        self.component_spec = component_spec
        self.name = component_spec.name
-       self.description = component_spec.description or None
+       self.pipeline_res = pipeline_spec_pb2.PipelineSpec()

        # Arguments typed as PipelineTaskFinalStatus are special arguments that
        # do not count as user inputs. Instead, they are reserved to for the
```

* 不要箇所を削除する。
```python
        with BlockPipelineTaskRegistration():
            return self.component_spec.to_pipeline_spec()

-   @property
-   def platform_spec(self) -> pipeline_spec_pb2.PlatformSpec:
-       """Returns the PlatformSpec of the component.
-
-       Useful when the component is a GraphComponent, else will be
-       empty per component_spec.platform_spec default.
-       """
-       return self.component_spec.platform_spec

    @abc.abstractmethod
    def execute(self, **kwargs):
```

4. `component_factory.py`をテキストエディタで開く。
```
$ vi sdk/python/kfp/dsl/component_factory.py
```

* `pipeline_spec_pb2`のインポートを追加する。
```python
from kfp.dsl.types import custom_artifact_types
from kfp.dsl.types import type_annotations
from kfp.dsl.types import type_utils
+from kfp.pipeline_spec import pipeline_spec_pb2

_DEFAULT_BASE_IMAGE = 'python:3.7'
SINGLE_OUTPUT_NAME = 'Output'
```

* compnents名の置換処理を`_`が含まれていた場合に`-`に変換する処理に修正する。
```python
REGISTERED_MODULES = None


def _python_function_name_to_component_name(name):
-   name_with_spaces = re.sub(' +', ' ', name.replace('_', ' ')).strip(' ')
-   return name_with_spaces[0].upper() + name_with_spaces[1:]
+   name_with_spaces = name.strip(' ').replace('_', '-')
+   return name_with_spaces


def make_index_url_options(pip_index_urls: Optional[List[str]]) -> str:
```

* 改修したcomponent定義にあわせて、不要な戻り値を削除する。
```python
    return structures.ComponentSpec(
        name=component_name,
-       description=description,
        inputs=name_to_input_spec or None,
        outputs=name_to_output_spec or None,
-       implementation=structures.Implementation(),
    )


```

* 改修したpipelineデコレータの引数にあわせて引数を変更する。
```python
def create_graph_component_from_func(
    func: Callable,
    name: Optional[str] = None,
-   description: Optional[str] = None,
-   display_name: Optional[str] = None,
+   dataflowName: Optional[str] = None,
+   dataflowNamespace: Optional[str] = None,
+   functionChainName: Optional[str] = None,
+   functionChainNamespace: Optional[str] = None,
+   functionKindNamespace: Optional[str] = None,
+   connectionKindNamespace: Optional[str] = None,    
) -> graph_component.GraphComponent:
    """Implementation for the @pipeline decorator.
    The decorator is defined under pipeline_context.py. See the
    decorator for the canonical documentation for this function.
    """
+   meta = pipeline_spec_pb2.MetaSpec()
+   meta.dataflowName = dataflowName
+   meta.dataflowNamespace = dataflowNamespace
+   meta.functionChainName = functionChainName
+   meta.functionChainNamespace = functionChainNamespace
+   meta.functionKindNamespace = functionKindNamespace
+   meta.connectionKindNamespace = connectionKindNamespace

    component_spec = extract_component_interface(
        func,
-       description=description,
        name=name,
    )
    return graph_component.GraphComponent(
        component_spec=component_spec,
        pipeline_func=func,
-       display_name=display_name,
+       meta_spec=meta,
    )


```

5. `graph_component.py`をテキストエディタで開く。
```
$ vi sdk/python/kfp/dsl/graph_component.py
```

* インポートに`Dict`と`Any`を追加する。
```python
"""Pipeline as a component (aka graph component)."""

import inspect
-from typing import Callable, Optional
+from typing import Callable, Optional, Dict, Any
import uuid

from kfp.compiler import pipeline_spec_builder as builder
```

* `pipeline_spec`作成処理である`create_pipeline_spec`メソッドに値を渡すため、コンストラクタに定義する。
```python
        self,
        component_spec: structures.ComponentSpec,
        pipeline_func: Callable,
+       meta_spec: pipeline_spec_pb2.MetaSpec(),
        display_name: Optional[str] = None,
    ):
        super().__init__(component_spec=component_spec)
```

* 改修後のIR YAMLのフォーマットに合わせて`platform_spec`の削除と`meta_spec`を追加する。
* 改修したcomponent定義には`description`と`implementation`がないため削除する。
* 戻り値を`create_pipeline_spec`関数で作成したpipeline_specに書き換える。
```python
        pipeline_group = dsl_pipeline.groups[0]
        pipeline_group.name = uuid.uuid4().hex

-       pipeline_spec, platform_spec = builder.create_pipeline_spec(
+       pipeline_spec = builder.create_pipeline_spec(
            pipeline=dsl_pipeline,
            component_spec=self.component_spec,
            pipeline_outputs=pipeline_outputs,
+           meta_spec=meta_spec,
        )

        pipeline_root = getattr(pipeline_func, 'pipeline_root', None)
        if pipeline_root is not None:
            pipeline_spec.default_pipeline_root = pipeline_root
        if display_name is not None:
            pipeline_spec.pipeline_info.display_name = display_name
-       if component_spec.description is not None:
-           pipeline_spec.pipeline_info.description = component_spec.description
-
-       self.component_spec.implementation.graph = pipeline_spec
-       self.component_spec.platform_spec = platform_spec
+       self.pipeline_res.CopyFrom(pipeline_spec)

    @property
    def pipeline_spec(self) -> pipeline_spec_pb2.PipelineSpec:
        """Returns the pipeline spec of the component."""
-       return self.component_spec.implementation.graph
+       return self.pipeline_res

    def execute(self, **kwargs):
        raise RuntimeError('Graph component has no local execution mode.')
```

## 次の手順について
[Pipeline改修](../modify-pipeline)

## 補足事項
(※1) OpenKasugaiコントローラのバージョン差分により仕様が異なる場合は、仕様に合わせた改修を実施する必要がある。
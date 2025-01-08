---
weight: 4
title: "Modify Components"
--- 
# Modify Components
## About This Procedure
* This document describes the modifications to the Kubeflow Pipelines SDK required to run the Custom Kubeflow SDK.

## Changes Made
* Edit `load_yaml_utilities.py` to add methods that align with the modified component definitions from [Modify pipeline_spec.proto](../modify-pipeline_spec.proto). (*1)
* Edit files under the `dsl` directory to add processing for creating pipeline_spec for the OpenKasugai controller environment. (*1)

## Prerequisites
* Ensure the following steps have been completed before proceeding with this procedure:
    * [Modify pipeline_spec.proto](../modify-pipeline_spec.proto)

## Procedure
1. Open `__init__.py` in a text editor.
```
$ vi sdk/python/kfp/components/__init__.py
```

* Add `load_component`.
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

2. Open `load_yaml_utilities.py` in a text editor.
```
$ vi sdk/python/kfp/components/load_yaml_utilities.py
```

* Add a process to convert `_` to `-` if it is included in the component name.
* Create the `load_component` method to align with the modified component definitions for the OpenKasugai environment.
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

3. Open `base_component.py` in a text editor.
```
$ vi sdk/python/kfp/dsl/base_component.py
```

* Add initialization to retrieve the pipelineSpec type.
```python
        """
        self.component_spec = component_spec
        self.name = component_spec.name
-       self.description = component_spec.description or None
+       self.pipeline_res = pipeline_spec_pb2.PipelineSpec()

        # Arguments typed as PipelineTaskFinalStatus are special arguments that
        # do not count as user inputs. Instead, they are reserved to for the
```

* Remove unnecessary parts.
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

4. Open `component_factory.py` in a text editor.
```
$ vi sdk/python/kfp/dsl/component_factory.py
```

* Add import of `pipeline_spec_pb2`.
```python
from kfp.dsl.types import custom_artifact_types
from kfp.dsl.types import type_annotations
from kfp.dsl.types import type_utils
+from kfp.pipeline_spec import pipeline_spec_pb2

_DEFAULT_BASE_IMAGE = 'python:3.7'
SINGLE_OUTPUT_NAME = 'Output'
```

* Modify the component name replacement process to convert `_` to `-`.
```python
REGISTERED_MODULES = None


def _python_function_name_to_component_name(name):
-   name_with_spaces = re.sub(' +', ' ', name.replace('_', ' ')).strip(' ')
-   return name_with_spaces[0].upper() + name_with_spaces[1:]
+   name_with_spaces = name.strip(' ').replace('_', '-')
+   return name_with_spaces


def make_index_url_options(pip_index_urls: Optional[List[str]]) -> str:
```

* Remove unnecessary return values according to the modified component definitions.
```python
    return structures.ComponentSpec(
        name=component_name,
-       description=description,
        inputs=name_to_input_spec or None,
        outputs=name_to_output_spec or None,
-       implementation=structures.Implementation(),
    )


```

* Update the arguments based on the modified pipeline decorator.
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

5. Open `graph_component.py` in a text editor.
```
$ vi sdk/python/kfp/dsl/graph_component.py
```

* Add `Dict` and `Any` to imports.
```python
"""Pipeline as a component (aka graph component)."""

import inspect
-from typing import Callable, Optional
+from typing import Callable, Optional, Dict, Any
import uuid

from kfp.compiler import pipeline_spec_builder as builder
```

* Define in the constructor to pass values to the `create_pipeline_spec` method, which creates the pipelineSpec.
```python
        self,
        component_spec: structures.ComponentSpec,
        pipeline_func: Callable,
+       meta_spec: pipeline_spec_pb2.MetaSpec(),
        display_name: Optional[str] = None,
    ):
        super().__init__(component_spec=component_spec)
```

* Remove `platform_spec` and add `meta_spec` to match the format of the modified IR YAML.
* Remove `description` and `implementation` based on the modified component definitions.
* Replace the return value with the pipeline_spec created by the `create_pipeline_spec` function.
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

## Next Steps
[Modify Pipeline](../modify-pipeline)

## Notes
(*1) Depending on the version differences of the OpenKasugai controller, it may be necessary to make modifications according to the specifications.
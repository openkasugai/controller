---
weight: 7
title: "Modify Structures"
---
# Modify Structures
## About This Procedure
* This document describes the modifications to the Kubeflow Pipelines SDK required to run the Custom Kubeflow SDK.

## Changes Made
* Edit `structures.py` to align with the specifications of the modified pipeline_spec from [Modify pipeline_spec.proto](../modify-pipeline_spec.proto). (*1)

## Prerequisites
* Ensure the following steps have been completed before proceeding with this procedure.
    * [Modify Compiler](../modify-compiler)

## Procedure
1. Open `structures.py` in a text editor.
```
$ vi sdk/python/kfp/dsl/structures.py
```

* Remove unnecessary parts to match the format of the modified IR YAML.
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

## Next Steps
[Build pipeline_spec.proto (Backend)](../../../build-procedure/apfw-backend-related-build-procedure/build-pipeline_spec.proto)

## Notes
(*1) If the specifications differ due to version differences in the OpenKasugai controller, modifications must be made to align with the specifications.
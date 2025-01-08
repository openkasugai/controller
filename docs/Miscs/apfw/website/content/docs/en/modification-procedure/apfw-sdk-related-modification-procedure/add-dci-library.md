---
weight: 1
title: "Add OpenKasugai Library"
---
# Add OpenKasugai Library
## About This Procedure
* Describes the modifications to the Kubeflow Pipelines SDK required to run the Custom Kubeflow SDK.

## Changes Made
* Create `__init__.py` and `swb.py` to add libraries for the OpenKasugai controller environment. (*1)

## Prerequisites
* Ensure the following steps have been completed before proceeding with this procedure.
    * [Delete Unnecessary Files](../../apfw-backend-related-modification-procedure/delete-unnecessary-files)

## Procedure
1. Create the swb directory.
```
$ mkdir sdk/python/kfp/swb
```

2. Create `__init__.py`.
```
$ vi sdk/python/kfp/swb/__init__.py
```

* Edit as follows.
```python
__all__ = [
    'swb',
]

from kfp.swb.swb import swb
```

3. Create `swb.py`.
```
$ vi sdk/python/kfp/swb/swb.py
```

* Edit as follows.
```python
from kfp import components
import yaml

class swb:
    def create_connection_task(name: str, fromPort: int, toPort: int, connectionKind: str):
        func={
            "components":{
                name:{
                    "fromPort": fromPort,
                    "toPort": toPort,
                    "connectionKind": connectionKind,
                },
                "inputDefinitions":{
                },
                "outputDefinitions":{
                },
            },
            "pipelineInfo":{
                "name": name
            }
        }

        connect_comp = components.load_component(yaml.dump(func))
        connect_task = connect_comp()

        return connect_task


    def start(startPointIP: str, startPointPort: int, startPointProtocol: str):
        func={
            "components":{
                "start":{
                    "startPointIP": startPointIP,
                    "startPointPort": startPointPort,
                    "startPointProtocol": startPointProtocol,
                },
                "inputDefinitions":{
                },
                "outputDefinitions":{
                },
            },
            "pipelineInfo":{
                "name": "start"
            }
        }

        start_comp = components.load_component(yaml.dump(func))
        start_task = start_comp()

        return start_task

    def end(endPointIP: str, endPointPort: int, endPointProtocol: str):
        func={
            "components":{
                "end":{
                    "endPointIP": endPointIP,
                    "endPointPort":endPointPort,
                    "endPointProtocol": endPointProtocol,
                },
                "inputDefinitions":{
                },
                "outputDefinitions":{
                },
            },
            "pipelineInfo":{
                "name": "end"
            }
        }

        end_comp = components.load_component(yaml.dump(func))
        end_task = end_comp()

        return end_task

    def FunctionItem(configName: str, coreMin: int, coreMax: int, totalBase: int, capacityTotalBase: int):
        items={
            "configName": configName,
            "coreMin": coreMin,
            "coreMax": coreMax,
            "totalBase": totalBase,
            "capacityTotalBase": capacityTotalBase
        }
        return items

    def Function(FunctionName, kindNameSpace, interface, FunctionItems, version):

        items = {}
        for key, value in FunctionItems.items():
            items[key] = value

        func={
            "components":{
                FunctionName:{
                    "namespace": kindNameSpace,
                    "info": { 
                        interface: {
                            "items": items,
                        },         
                    },
                    "version": version
                },
                "inputDefinitions":{
                    "parameters":{
                        "version":{
                            "parameterType": "STRING"
                        },
                    },
                },
                "outputDefinitions":{
                    "parameters":{
                        "Output":{
                            "parameterType": "NUMBER_INTEGER"
                        },
                    },
                },
            },
            "pipelineInfo":{
                "name": FunctionName
            }
        }

        return yaml.dump(func)

```

## Next Steps
[Add Requirements](../add-requirements)

## Notes
(*1) Depending on version differences in the OpenKasugai controller, modifications may be required to align with the specifications.
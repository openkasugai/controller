---
weight: 1
title: "OpenKasugaiライブラリ追加"
---
# OpenKasugaiライブラリ追加
## この手順について
* Custom Kubeflow SDKを実行するにあたり、Kubeflow Pipelines SDKの修正箇所について記載する。

## 対応内容
* `__init__.py`と`swb.py`を新規作成し、OpenKasugaiコントローラ環境用にライブラリを追加する。(※1)

## 前提条件
* 当該手順を実施する前に以下手順を実施済みであること。
    * [不要ファイル削除](../../apfw-backend-related-modification-procedure/delete-unnecessary-files)

## 手順
1. swbディレクトリを作成する。
```
$ mkdir sdk/python/kfp/swb
```

2. `__init__.py`を新規作成する。
```
$ vi sdk/python/kfp/swb/__init__.py
```

* 以下のように編集する。
```python
__all__ = [
    'swb',
]

from kfp.swb.swb import swb
```

3. `swb.py`を新規作成する。
```
$ vi sdk/python/kfp/swb/swb.py
```

* 以下のように編集する。
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

## 次の手順について
[requirements追加](../add-requirements)

## 補足事項
(※1) OpenKasugaiコントローラのバージョン差分により仕様が異なる場合は、仕様に合わせた改修を実施する必要がある。
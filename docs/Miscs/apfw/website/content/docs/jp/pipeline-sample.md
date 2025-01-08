---
weight: 9
title: "Pipelineサンプル"
---
# Pipelineサンプル
## 当該ドキュメントで使用するPipelineサンプル
* サンプルとして使用するPipelineのPythonコードは以下の通り
* 以下PipelineはDataflowをデプロイし推論用Podを作成する処理となる。
* `< >`で囲まれた変数は環境に合わせ書き換えて使用する。
  * pipelineはKubeflowのml-pipeline-uiに向けてAPIを送信する。
  * そのため、ml-pipeline-uiサービスIP及びポート番号を環境にあわせて設定する必要がある。

```python
from kfp import dsl
from kfp import compiler
from kfp import components
from kfp.client import Client
from kfp.swb import swb

decode = {}
filter_resize_high_infer = {}
high_infer = {}

decode['dev25gether'] = swb.FunctionItem(configName="fpgafunc-config-decode", coreMin=1, coreMax=1, totalBase=6, capacityTotalBase=20)

filter_resize_high_infer['dev25gether'] = swb.FunctionItem(configName="fpgafunc-config-filter-resize-high-infer", coreMin=1, coreMax=1, totalBase=8, capacityTotalBase=40)
filter_resize_high_infer['mem'] = swb.FunctionItem(configName="fpgafunc-config-filter-resize-high-infer", coreMin=1, coreMax=1, totalBase=8, capacityTotalBase=40)

high_infer['host100gether'] = swb.FunctionItem(configName="gpufunc-config-high-infer", coreMin=1, coreMax=1, totalBase=1, capacityTotalBase=15)
high_infer['mem'] = swb.FunctionItem(configName="gpufunc-config-high-infer", coreMin=1, coreMax=1, totalBase=1, capacityTotalBase=15)

function1 = swb.Function("decode", "wbfunc-imgproc", "alveo", decode, version="1.0.0")
function2 = swb.Function("filter-resize-high-infer", "wbfunc-imgproc", "alveo", filter_resize_high_infer, "1.0.0")
function3 = swb.Function("high-infer", "wbfunc-imgproc", "a100", high_infer, "1.0.0")

decode = components.load_component(function1)
filter_resize_high_infer = components.load_component(function2)
high_infer = components.load_component(function3)


@dsl.pipeline(dataflowName="df-fpgadec-highinf-01",
              dataflowNamespace="test01",
              functionChainName="decode-filter-resize-high-infer-chain",
              functionChainNamespace="chain-imgproc",
              functionKindNamespace="wbfunc-imgproc",
              connectionKindNamespace="default")
def my_pipeline():
    _start = swb.start("", 0, "")
    _connect001_task = swb.create_connection_task("connect_001", 0, 0, "auto").after(_start)
    _decode = decode(version="1.0.0").after(_connect001_task)
    _connect002_task = swb.create_connection_task(name="connect_002", fromPort=0, toPort=0, connectionKind="auto").after(_decode)
    _filter_resize_high_infer = filter_resize_high_infer(version="1.0.0").after(_connect002_task)
    _connect003_task = swb.create_connection_task(name="connect_003", fromPort=0, toPort=0, connectionKind="auto").after(_filter_resize_high_infer)
    _high_infer = high_infer(version="1.0.0").after(_connect003_task)
    _connect004_task = swb.create_connection_task(name="connect_004", fromPort=0, toPort=0, connectionKind="auto").after(_high_infer)
    _end = swb.end("", 0, "").after(_connect004_task)


compiler.Compiler().compile(my_pipeline, package_path='no-start-end-pipeline.yaml')

client = Client(host='http://<ml-pipeline-ui Service IP>:<Port Number>')
client.create_run_from_pipeline_package('no-start-end-pipeline.yaml')
```



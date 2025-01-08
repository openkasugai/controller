---
weight: 7
title: "Makefile改修"
---
# Makefile改修
## この手順について
* Custom Kubeflow SDKを実行するにあたりKubeflow Pipelines Backendの修正箇所について記載する。

## 対応内容
* Kubeflow Pipelines Backendの`Makefile`を編集し、`github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec`を`go-licenses`によるGoパッケージのライセンス確認から除外する。(※1)

## 前提条件
* 当該手順を実施する前に以下手順を実施済みであること。
    * [install-go-licenses.sh改修](../modify-install-go-licenses.sh)

## 手順
1. `Makefile`を編集する。
```
vi backend/Makefile
```
* 以下のように編集する。
```Makefile
# See README.md#updating-licenses-info section for more details.
.PHONY: license_apiserver
license_apiserver: $(BUILD)/apiserver
-       cd $(MOD_ROOT) && go-licenses csv ./backend/src/apiserver > $(CSV_PATH)/apiserver.csv
+       cd $(MOD_ROOT) && go-licenses csv --ignore github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec ./backend/src/apiserver > $(CSV_PATH)/apiserver.csv
.PHONY: license_persistence_agent
license_persistence_agent: $(BUILD)/persistence_agent
        cd $(MOD_ROOT) && go-licenses csv ./backend/src/agent/persistence > $(CSV_PATH)/persistence_agent.csv

```
## 次の手順について
[Dockerfile改修](../modify-dockerfile)

## 補足事項
(※1) 当該ドキュメントの[go.mod改修](../modify-go.mod)によって、`github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec`のライセンス情報を取得できなくなっているため、`go-licenses`によるGoパッケージのライセンス確認から除外する。
---
weight: 8
title: "Dockerfile改修"
---
# Dockerfile改修
## この手順について
* Custom Kubeflow SDKを実行するにあたりKubeflow Pipelines Backendの修正箇所について記載する。

## 対応内容
* `Dockerfile`を編集し、`github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec`を`go-licenses`によるGoパッケージのライセンス確認から除外する。(※1)

## 前提条件
* 当該手順を実施する前に以下手順を実施済みであること。
    * [Makefile改修](../modify-makefile)

## 手順
1. `Dockerfile`をテキストエディタで開く。
```
$ vi backend/Dockerfile
```

2. Goパッケージのライセンス確認から`github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec`を除外する。
```Dockerfile
 # Check licenses and comply with license terms.
 RUN ./hack/install-go-licenses.sh
 # First, make sure there's no forbidden license.
-RUN go-licenses check ./backend/src/apiserver
-RUN go-licenses csv ./backend/src/apiserver > /tmp/licenses.csv && \
+RUN go-licenses check --ignore github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec ./backend/src/apiserver
+RUN go-licenses csv --ignore github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec ./backend/src/apiserver > /tmp/licenses.csv && \
   diff /tmp/licenses.csv backend/third_party_licenses/apiserver.csv && \
-  go-licenses save ./backend/src/apiserver --save_path /tmp/NOTICES
+  go-licenses save --ignore github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec ./backend/src/apiserver 
--save_path /tmp/NOTICES
 
 # 2. Compile preloaded pipeline samples
 FROM python:3.7 as compiler
 ```

## 次の手順について
[不要ファイル削除](../delete-unnecessary-files)

## 補足事項
(※1) 当該ドキュメントの[go.mod改修](../modify-go.mod)によって、`github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec`のライセンス情報を取得できなくなっているため`go-licenses`によるGoパッケージのライセンス確認から除外する。
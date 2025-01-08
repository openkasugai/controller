---
weight: 6
title: "install-go-licenses.sh改修"
---
# install-go-licenses.sh改修
## この手順について
* Custom Kubeflow SDKを実行するにあたりKubeflow Pipelines Backendの修正箇所について記載する。

## 対応内容
* `install-go-licenses.sh`を編集し、インストールする`go-licenses`のバージョンを変更する。(※1)

## 前提条件
* 当該手順を実施する前に以下手順を実施済みであること。
    * [v2_template.go改修](../modify-v2_template.go)

## 手順
1. `install-go-licenses.sh`を編集する。
```
$ vi hack/install-go-licenses.sh
```
* 以下のように編集する。
```sh
set -ex

# TODO: update to a released version.
-go install github.com/google/go-licenses@d483853
+go install github.com/google/go-licenses@706b9c60
```

## 次の手順について
[Makefile改修](../modify-makefile)

## 補足事項
(※1) 
* `go-licenses`はGoパッケージの依存関係を分析し、使用されているライブラリと使用が許諾されるライセンスについて確認、出力ができる。
* バージョン変更することで、指定したGoパッケージを`go-licenses`による確認の対象外にできる`ignore`オプションが使用可能になる。
---
weight: 9
title: "不要ファイル削除"
---
# 不要ファイル削除
## この手順について
* Custom Kubeflow SDKを実行するにあたりKubeflow Pipelines Backendの修正箇所について記載する。

## 対応内容
* 使用しない不要なファイルを削除する。

## 前提条件
* 当該手順を実施する前に開発環境の準備として以下手順を実施済みであること。
    * [Dockerfile改修](../modify-dockerfile)

#### `backend/src/v2`配下の以下ファイルを削除する。
改修にあたり利用しなくなった関数が存在し、ファイルを削除しないとビルドエラーとなる。
* compiler/argocompiler/argo_test.go
* compiler/argocompiler/dag.go
* compiler/visitor_test.go
* component/importer_launcher.go
* component/launcher_v2_test.go
* component/launcher_v2.go

## 次の手順について
[OpenKasugaiライブラリ追加](../../apfw-sdk-related-modification-procedure/add-dci-library)
---
weight: 1
title: "pipeline_spec.protoビルド(バックエンド)"
---
# pipeline_spec.protoビルド(バックエンド)
## この手順について
* Custom Kubeflow Backendをビルドするには事前に`pipeline_spec.proto`をビルドしコードを生成する必要がある。
* `pipeline_spec.proto`はCustom Kubeflow SDKとCustom Kubeflow Backend間でやり取りをする中間表現であるIR YAMLのインターフェース定義が記述されており、Custom Kubeflow SDKとCustom Kubeflow Backendで共用する。
* Custom Kubeflow SDKはPython、Custom Kubeflow BackendはGoで出力する必要があるためビルド手順が異なる。
* そのため、この手順ではCustom Kubeflow Backendビルド前に必要な`pipeline_spec.proto`をビルドし、`pipeline_spec.pb.go`を用意する手順について記載する。

## 対応内容
* `pipeline_spec.proto`をビルドしコードを生成する。
* `pipeline_spec.pb.go`が指定のパスに生成されることを確認する。

## 前提条件
### ビルド環境について
* 当該手順をもとに`pipeline_spec.proto`をビルドする場合、ビルド環境には`Docker`がインストールされている必要がある。
* インストールする`Docker`のバージョン情報については以下を参照   
[環境情報・前提条件](../../../environment-information-and-prerequisites)

### 事前手順
* 当該手順実施前に以下手順が実施済みであることを前提とする。
    * [structures改修](../../../modification-procedure/apfw-sdk-related-modification-procedure/modify-structures)

## 手順
1. コンパイル先のAPIバージョンを指定する。
```
$ export API_VERSION="v2beta1"
```

2. `api`ディレクトリ配下に移動する。
```
$ cd pipelines/api
```

3. 以下、`make`コマンドを実行する。
```
$ make clean-go golang
```
コンパイル後ファイルが以下pathに作成されることを確認する。
`pipelines/api/v2alpha1/go/pipelinespec/pipeline_spec.pb.go`

## 次の手順について
* [ライセンス情報の更新](../updating-licenses-info)
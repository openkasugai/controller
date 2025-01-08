---
weight: 1
title: "pipeline_spec.protoビルド(SDK)"
---
# pipeline_spec.protoビルド(SDK)
## この手順について
* Custom Kubeflow SDK実行時にYAML出力される各項目は、`pipeline_spec.proto`で管理している。
* `pipeline_spec.proto`はCustom Kubeflow SDKとCustom Kubeflow Backend間でやり取りをする中間表現であるIR YAMLのインターフェース定義が記述されており、Custom Kubeflow SDKとCustom Kubeflow Backendで共用する。
* Custom Kubeflow SDK改修時にYAML項目を変更するときは、`pipeline_spec.proto`の改修およびビルドを行う必要がある。
* Custom Kubeflow SDKはPython、Custom Kubeflow BackendはGoで出力する必要があるためビルド手順が異なる。
* そのため、この手順ではCustom Kubeflow SDKに必要な`pipeline_spec.proto`をビルドし、`pipeline_spec.pb.py`を用意する手順について記載する。

## 対応内容
* `pipeline_spec.proto`をビルドしコードを生成する。
* `pipeline_spec_pb2.py`が指定のパスに生成されることを確認する。

## 前提条件
### ビルド環境について
* 当該手順をもとに`pipeline_spec.proto`をビルドする場合、ビルド環境には`protobuf-compiler`がインストールされている必要がある。
* インストールする`protobuf-compiler`のバージョン情報については以下を参照   
[環境情報・前提条件](../../../environment-information-and-prerequisites)

### 事前手順
* 当該手順実施前に以下手順が実施済みであることを前提とする。
    * [pipeline_spec.proto改修](../../apfw-backend-related-build-procedure/build-apfw-backend)

## 手順
1. `api`ディレクトリ配下に移動する。
```
$ cd pipelines/api
```

2. 以下、`make`コマンドを実行する。
```
$ make clean-python python
```

コンパイル後ファイルが以下pathに作成されることを確認する。
`pipelines/api/v2alpha1/python/kfp/pipeline_spec/pipeline_spec_pb2.py`

## 次の手順について
* [Namespace作成とRole設定](../../../deployment-procedure/configuration-on-dci-controller-node/create-namespace-and-set-role)

---
weight: 3
title: "Custom Kubeflow Backendビルド"
---
# Custom Kubeflow Backendビルド
## この手順について
* Kubeflow Pipelines Backendを改修後にビルドを実施する。
* Custom KubeflowからのIR YAML書式をマニフェストに変換する役割をになう`api-server`コンテナをビルドし、コンテナイメージを作成する。

## 対応内容
* `api-server`コンテナイメージのビルドを行う。

## 前提条件
### ビルド環境について
* ビルドに必要なソフトウェアおよびバージョンについては以下を参照
    * [環境情報・前提条件](../../../environment-information-and-prerequisites)

### 事前手順
* 当該手順実施前に以下手順が実施済みであることを前提とする。
    * [pipeline_spec.protoビルド(バックエンド)](../build-pipeline_spec.proto)
    * [ライセンス情報の更新](../updating-licenses-info)

## 手順
1. コンテナイメージをビルドする。
* pipelinesフォルダ配下で、以下コマンドでコンテナイメージをビルドする。
* `TAG`は任意のタグを指定する。
```
$ cd pipelines
$ docker build -t api-server:<TAG> -f backend/Dockerfile .
```

2. `api-server`のコンテナイメージが作成されることを確認する。
```
$ docker images
REPOSITORY   TAG   　　    IMAGE ID   　　CREATED         SIZE
api-server   <TAG>  <コンテナID>    <経過時間>  　　219MB
```

## 次の手順について
* [pipeline_spec.protoビルド(SDK)](../../apfw-sdk-related-build-procedure/build-pipeline_spec.proto)
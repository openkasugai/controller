---
weight: 3
title: "Dataflowのデプロイによる動作確認"
---
# Dataflowのデプロイによる動作確認
## この手順について
* Custom Kubeflow SDKの動作確認として、Dataflowのデプロイによる動作確認手順について記載する。

## 対応内容
* Custom Kubeflow SDKを実行する。
* Dataflowがデプロイされることを確認する。

## 前提条件
### 実行環境について
* 当該手順をもとにCustom Kubeflow SDKを実行する場合、実行環境には`python3`のインストール及び`venv`の仮想環境が構築されている必要がある。
* Custom Kubeflow SDKの実行でDataflowをデプロイするためにはOpenKasugaiコントローラが動作している必要がある。
* インストールする`python3`のバージョン情報と`venv`の構築手順及びOpenKasugaiコントローラ構築に使用するドキュメントについては以下を参照   
    * [環境情報・前提条件](../../../environment-information-and-prerequisites)

### 事前手順
* 当該手順実施前に以下手順が実施済みであることを前提とする。
    * [デプロイ手順](../../../deployment-procedure/)

## 手順
1. 仮想環境に切り替える。
* `VENV_DIR`は`venv`仮想環境を構築した時に指定したフォルダ名を設定する。
```
$ source <VENV_DIR>/bin/activate
```

2. プロンプトの前方にフォルダ名が表示されPython仮想実行環境に切り替わったことを確認する。
```
<VENV_DIR>user_name@host:$
```

3. パイプラインスクリプトを実行する。
* pipelinesは実行するパイプラインスクリプト名を指定する。
* 当該ドキュメントでは[Pipelineサンプル](../../../pipeline-sample)の実行を例を記載している。
```
$ python <pipelines>.py
```

4. Configmapがデプロイされていることを確認する。
```
$ kubectl get configmap -A | grep -e NAME -e func | grep -v ^default
```

5. FunctionKindがデプロイされていることを確認する。
```
$ kubectl get functionkind -A
```

6. FunctionChainがデプロイされていることを確認する。
```
$ kubectl get functionchain -A
```

7. Dataflowがデプロイされていることを確認する。
```
kubectl get dataflow -A
```
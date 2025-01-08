---
weight: 2
title: "IR YAMLの出力による動作確認"
---
# IR YAMLの出力による動作確認
## この手順について
* OpenKasugaiコントローラがない環境でのCustom Kubeflow SDKの動作確認として、IR YAMLの出力を確認する手順について記載する。

## 対応内容
* Custom Kubeflow SDKを実行する。
* IR YAMLが生成されることを確認する。

## 前提条件
### 実行環境について
* 当該手順をもとにCustom Kubeflow SDKを実行する場合、実行環境には`python3`のインストール及び`venv`の仮想環境が構築されている必要がある。
* インストールする`python3`のバージョン情報と`venv`の構築手順については以下を参照   
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

スクリプト実行後、IR YAMLファイルがコマンド実行時のディレクトリに作成されることを確認する。

4. IR YAMLの出力内容を確認する。
* pipelineにて設定した値がIR YAMLの設定値に反映していることを確認する。各出力項目ごとに紐づく設定値は以下の通り。

| pipeline | IR YAML |
| -------- | ------- |
| FunctionItemの設定値 | Itemsの設定値 |
| Functionの設定値 | components.設定したcomponent名の設定値 |
| Dataflow名等のメタデータ | metadataの設定値 |
| create_connection_taskの設定値 | dag.taskの設定値 |

## 次の手順について
* [Dataflowのデプロイによる動作確認](../operation-confirmation-by-dataflow-deployment)

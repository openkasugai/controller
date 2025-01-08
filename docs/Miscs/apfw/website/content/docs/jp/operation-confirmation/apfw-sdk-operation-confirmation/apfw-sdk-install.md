---
weight: 1
title: "Custom Kubeflow SDKインストール手順"
---
# Custom Kubeflow SDKインストール手順
## この手順について
* Custom Kubeflow SDKのインストール手順について記載する。

## 対応内容
* Custom Kubeflow SDKのインストールを実行する。
* `pipeline_spec.proto`をビルドし生成された`pipeline_spec_pb2.py`を既存の`pipeline_spec_pb2.py`置き換える。

## 前提条件
### 実行環境について
* 当該手順をもとにCustom Kubeflow SDKのインストールを実行する場合、`venv`の仮想環境が構築されている必要がある。
* インストールする`venv`の構築手順については以下を参照   
[環境情報・前提条件](../../../environment-information-and-prerequisites)

### 事前手順
* 当該手順実施前に以下手順が実施済みであることを前提とする。
    * [デプロイ手順](../../../deployment-procedure/)

## 手順
1. 仮想環境への切り替える。
* `VENV_DIR`は`venv`仮想環境が構築時に指定したフォルダ名を設定する。
```
$ source <VENV_DIR>/bin/activate
```

2. プロンプトの前方にフォルダ名が表示されPython仮想実行環境に切り替わったことを確認する。
```
<VENV_DIR>user_name@host:$
```

3. Custom Kubeflow SDKディレクトリに移動する。
* DIRは[go.mod改修](../../../modification-procedure/apfw-backend-related-modification-procedure/modify-go.mod)でKubeflowリポジトリをcloneしたディレクトリを指定する。
```
$ cd <DIR>/pipelines/sdk/python/kfp/
```

4. Custom Kubeflow SDKをインストールする。
```
$ pip install .
```

5. Custom Kubeflow SDKをインストールされていることを確認する。
```
$ pip list | grep kfp
```

* 以下の通り出力することを確認する。
```
kfp                      2.4.0
kfp-pipeline-spec        0.2.2
kfp-server-api           2.0.5
```

6. `pipeline_spec.proto`をビルドし生成された`pipeline_spec_pb2.py`を既存の`pipeline_spec_pb2.py`置き換える。
```
$ cp <DIR>/pipelines/api/v2alpha1/python/kfp/pipeline_spec/pipeline_spec_pb2.py <VENV_DIR>/lib64/python3.8/site-packages/kfp/pipeline_spec/pipeline_spec_pb2.py
```

## 次の手順について
* [IR YAMLの出力による動作確認](../operation-confirmation-by-iryaml-output)
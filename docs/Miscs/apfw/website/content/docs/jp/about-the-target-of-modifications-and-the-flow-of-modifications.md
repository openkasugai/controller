---
weight: 3
title: "改修対象と改修の流れについて"
---
# 改修対象と改修の流れについて
## 概要
* 改修対象と改修の流れについて記載する。

## 改修対象
* 当該ドキュメントの改修対象は`Kubeflow Pipelines`プロジェクトの`sdk-2.4.0`タグとする。

## 改修の目的
* Kuberflowを用いて、OpenKasugaiコントローラのDataFlow実行に要求するカスタムリソースをデプロイ可能とする。
* OpenKasugaiコントローラが要求するカスタムリソースの仕様に合わせるため、Kubeflowに含むKubeflow Pipelines Backend及びKubeflow Pipelines SDKを改修する。

## 構成図
![apfw arch](/images/apfw-dci_arch.jpg)
* 当該ドキュメントでサンプルとして使用する[Pipelineサンプル](../pipeline-sample)は、図の`Python Pipeline`にあたる。
* `Python Pipeline`はFunctionInfo（Configmap）, FunctionKind, FunctionChain, Dataflowをデプロイするための情報をPythonコードで表現できる。
* Kubeflowを改修することでIR YAMLのフォーマットを変更し、OpenKasugaiのカスタムリソースであるDataflowとFunctionChainをデプロイ可能にする。

## 改修内容
* 改修内容は以下の通り
* OpenKasugaiコントローラのバージョン差分によるリソースの違いで、改修内容に差分が発生する箇所については注釈を入れる。

### Kubeflow Pipelines Backend
| 項目 | 改修内容 |
| ---- | -------- |
| [go.mod改修](../modification-procedure/apfw-backend-related-modification-procedure/modify-go.mod) | ビルド時のコードの参照先をコンテナ上のローカルファイルに変更 |
| [argo.go改修](../modification-procedure/apfw-backend-related-modification-procedure/modify-argo.go) | Custom Kubeflow SDKでデプロイする各リソースに合わせて改修 |
| [visitor.go改修](../modification-procedure/apfw-backend-related-modification-procedure/modify-visitor.go) | 使用しないimportと関数を削除 |
| [resource_manager.go改修](../modification-procedure/apfw-backend-related-modification-procedure/modify-resource_manager.go) | namespaceの設定処理を削除 |
| [v2_template.go改修](../modification-procedure/apfw-backend-related-modification-procedure/modify-v2_template.go) | 使用しないimportと関数の削除、及びsetDefaultServiceAccountの設定値を変更 |
| [Dockerfile改修](../modification-procedure/apfw-backend-related-modification-procedure/modify-dockerfile) | 改修によりライセンス情報を取得できなくなっているGoパッケージをライセンス確認から除外 |
| [不要ファイル削除](../modification-procedure/apfw-backend-related-modification-procedure/delete-unnecessary-files) | 使用しない不要なファイルを削除 |

### Kubeflow Pipelines SDK
| 項目 | 改修内容 |
| ---- | -------- |
| [OpenKasugaiライブラリ追加](../modification-procedure/apfw-sdk-related-modification-procedure/add-dci-library) | OpenKasugaiコントローラ向けに新規ライブラリを作成(※1) |
| [requirements追加](../modification-procedure/apfw-sdk-related-modification-procedure/add-requirements) | 当該ドキュメントの手順にて改修した`kfp-pipeline-spec`を参照するよう`requirements.in`と`requirements.txt`を新規作成 |
| [pipeline_spec.proto改修](../modification-procedure/apfw-sdk-related-modification-procedure/modify-pipeline_spec.proto) | Custom Kubeflow SDK実行時に生成するIR YAMLのフォーマットにあわせてデータ構造を改修 |
| [Components改修](../modification-procedure/apfw-sdk-related-modification-procedure/modify-components) | Custom Kubeflow SDKのcomponent定義方法にあわせたメソッドの追加とOpenKasugaiコントローラ環境用のpipeline_spec作成処理を追加 |
| [Pipeline改修](../modification-procedure/apfw-sdk-related-modification-procedure/modify-pipeline) | OpenKasugaiコントローラ環境用に改修したpipeline_specの仕様に合わせて改修 |
| [Compiler改修](../modification-procedure/apfw-sdk-related-modification-procedure/modify-compiler) | OpenKasugaiコントローラ環境用に改修したpipeline_specの仕様に合わせて改修 |
| [structures改修](../modification-procedure/apfw-sdk-related-modification-procedure/modify-structures) | OpenKasugaiコントローラ環境用に改修したpipeline_specの仕様に合わせて改修 |

## 改修の流れ
* 当該ドキュメントでは以下の流れで改修から動作確認を実施する。
1. Kubeflow Pipelines Backend及びKubeflow Pipelines SDKを改修し、Custom Kubeflow Backend及びCustom Kubeflow SDKを作成
2. pipeline_spec.proto及びCustom Kubeflow Backendのビルド
3. OpenKasugaiコントローラへの設定とCustom Kubeflow Backendのデプロイ
4. Custom Kubeflow SDKを実行し動作確認

## 次の手順について
[環境情報・前提条件](../environment-information-and-prerequisites)

## 補足事項
(※1) ライブラリ追加手順中のライブラリ名`swb`はOpenKasugaiを指している。
---
weight: 4
title: "環境情報・前提条件"
---
# 環境情報・前提条件
## 概要
* 当該ドキュメントに記載の手順を実施するために必要な環境情報と前提条件について記載する。
* 当該ドキュメントではOpenKasugaiコントローラがインストールされていることを基本とするが、インストールされていない場合でも一部手順を除き実施が可能。
  * OpenKasugaiコントローラがインストールされていない場合、以下手順の実施は不可
    * [Namespace作成とRole設定](../deployment-procedure/configuration-on-dci-controller-node/create-namespace-and-set-role)
    * [Custom Kubeflow Backendデプロイ](../deployment-procedure/apfw-backend-related-deployment-procedure/deploy-apfw-backend)
    * [Dataflowのデプロイによる動作確認](../operation-confirmation/apfw-sdk-operation-confirmation/operation-confirmation-by-dataflow-deployment)
  * OpenKasugaiコントローラ環境を用意する場合は[OpenKasugaiコントローラドキュメント](https://github.com/openkasugai/controller/blob/main/README_jp.md#%E3%83%89%E3%82%AD%E3%83%A5%E3%83%A1%E3%83%B3%E3%83%88)を使用し構築する。

## 環境情報
* 環境は以下の通り
| OS・ソフトウェア         | バージョン(※1) | 
| ------------ | ---- | 
| Ubuntu         | 20.04.2 LTS | 

## 前提条件
* 以下のソフトウェアがインストールされていること
| ソフトウェア         | バージョン(※1) | インストール手順 |
| ------------ | ---- | ---- | 
| Python         | 3.8.10 | https://wiki.python.org/moin/BeginnersGuide/Download |
| Go             | 1.20.14 | https://go.dev/doc/install |
| libprotoc         | 3.6.1 | https://grpc.io/docs/protoc-installation/ |
| Docker         | 24.0.5 | https://docs.docker.com/engine/install/ubuntu/ |
| Harbor         | 2.9.1 | https://goharbor.io/docs/2.9.0/install-config/ |

* venvの仮想環境が構築されていること
| ソフトウェア         | 構築手順 | 
| ------------ | ---- |
| venv       | https://docs.python.org/3.8/library/venv.html |

## 次の手順について
[go.mod改修](../modification-procedure/apfw-backend-related-modification-procedure/modify-go.mod)

## 補足事項
(※1) ドキュメント作成時点で使用したOS・ソフトウェアのバージョンを記載している。
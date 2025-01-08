---
weight: 2
title: "概要"
---
# 概要
## このドキュメントについて
* OpenKasugaiコントローラに対し推論/機械学習Pipeline(以下「Pipeline」と表記する。)を実行するためのガイドラインである。
* 当該ドキュメントではPipeline実行の実現方法のひとつとして[Kubeflow](https://www.kubeflow.org/)を使用する例を記載する。(※1)
* KubeflowとOpenKasugaiコントローラの連携を行うためには、KubeflowのソースコードをOpenKasugaiコントローラの仕様にあわせて書き換える必要がある。
* そのため、当該ドキュメントでは以下について記述している。
  * OpenKasugaiコントローラ(※2)の仕様にあわせたKubeflowの改修箇所と改修方法
  * 改修後のKubeflowの実行と確認方法

## 対象読者
* Pipeline実行を検討しているOpenKasugaiコントローラの利用者
* Pipelineを実行するための開発内容を把握したいOpenKasugaiコントローラの開発者

## このドキュメントで実現できること
* OpenKasugaiコントローラの仕様にあわせたKubeflowの改修(※3)
* OpenKasugaiコントローラに対してのPipelineの実行およびリソース(Dataflow)のデプロイ(※4)
* 改修後のKubeflow単体での実行と中間表現ファイルによる正常動作確認
* 上記手順を参考とした異なるバージョンのOpenKasugaiコントローラへの改修手順の適用(※5)

## Pipeline実行環境を整備する利点
* 当該ドキュメントをもとにPipeline実行環境を整備すると以下の利点がある。
  * PythonでPipelineを記述するため、Pythonの知識があれば容易にPipeline作成ができる。
  * 必要最低限の情報を与えるだけでリソースがデプロイできるため、OpenKasugaiコントローラ用のリソース作成に関する学習コストを削減できる。
  * OpenKasugaiコントローラに対しAPIによる遠隔からのPipeline実行およびリソース作成が実現できるようになる。

## 用語定義
以降の章では以下用語定義をもとに記述を行う。
| 用語         | 定義 | 
| ------------ | ---- | 
| [Kubeflow Pipelines SDK](https://www.kubeflow.org/docs/components/pipelines/)      | 改修前のプレーンなKubeflow Pipeline SDK。 | 
| Custom Kubeflow SDK      | OpenKasugaiコントローラの仕様にあわせて改修途中。もしくは改修済みのKubeflow Pipelines SDK。 | 
| [Kubeflow Pipelines Backend](https://github.com/kubeflow/pipelines/tree/sdk-2.4.0/backend) | Kubernetes上でKubeflowを起動するためのコンポーネントやコード。 | 
| Custom Kubeflow Backend | OpenKasugaiコントローラの仕様にあわせて改修途中。もしくは改修済みのKubeflow Pipelines Backend。 | 

## 次の手順について
[改修対象と改修の流れについて](../about-the-target-of-modifications-and-the-flow-of-modifications)

## 補足事項
(※1) Kubeflowを選定した理由としては以下があげられる。
1. 環境適合性：オンプレミスや各種クラウドで利用でき、Kubernetesとの親和性も高いため。
2. 入出力容易性: Pythonによる任意のフローを記述でき、WebUIからのフロー定義が可能であるため。
3. 更新頻度: 年間更新回数10回程度とコミュニティの活性度が比較的高いため。

(※2) 当該ドキュメントで動作確認に使用したOpenKasugaiコントローラのバージョンは一般に未公開であり、映像処理に使用するFPGA回路及び設定に関しても未公開である。   
(※3) 
* 当該ドキュメントのKubeflowの改修を実施した場合、改修前のKubeflowの仕様に合わせたPipeline実行は不可となる。
* KubeflowとカスタムリソースをデプロイするOpenKasugaiコントローラではIR YAMLのフォーマットが異なり、既存のPipeline実行を可能にすることが難しいため、既存の動作に影響しないようにするという改修アプローチはとっていない。 

(※4) 当該ドキュメントでサンプルとして使用する[Pipelineサンプル](../pipeline-sample)はDataflowをデプロイし推論用Podを作成する処理となる。推論用Pod内の推論処理は事前に構築済みの学習モデルを使用する。  
(※5) 旧版OpenKasugaiコントローラに対応した改修内容を公開することでバージョンが異なるOpenKasugaiコントローラへのPipeline実行を実現できるようにすることを目的としている。
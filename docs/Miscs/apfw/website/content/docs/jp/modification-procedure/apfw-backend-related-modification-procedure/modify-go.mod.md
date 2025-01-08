---
weight: 1
title: "go.mod改修"
---
# go.mod改修
## この手順について
* Custom Kubeflow SDKを実行するにあたり、Kubeflow Pipelines Backendの修正箇所について記載する。

## 対応内容
* `go.mod`を編集し、後の手順にてコンテナ上でビルドする際に参照するソースコードを、当該ドキュメントで修正したものとするよう変更する。

## 手順
1. 改修対象となるバージョンのKubeflowリポジトリをCustom Kubeflow SDK実行環境にcloneする。(※1)
```
$ git clone https://github.com/kubeflow/pipelines.git -b sdk-2.4.0
$ cd pipelines
```

2. `go.mod`を編集する。
```
$ cd pipelines
$ vi go.mod
```

以下のように編集する。
```go
module github.com/kubeflow/pipelines
require (
	github.com/Masterminds/squirrel v0.0.0-20190107164353-fa735ea14f09
	github.com/VividCortex/mysqlerr v0.0.0-20170204212430-6c6b55f8796f
	github.com/argoproj/argo-workflows/v3 v3.3.10
	github.com/aws/aws-sdk-go v1.42.50
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/eapache/go-resiliency v1.2.0
	github.com/elazarl/goproxy v0.0.0-20181111060418-2ce16c963a8a // indirect
	github.com/emicklei/go-restful v2.16.0+incompatible // indirect
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5 // indirect
~~OMITTED~~
	sigs.k8s.io/controller-runtime v0.11.1
    sigs.k8s.io/yaml v1.3.0
+   gopkg.in/yaml.v2 v2.4.0
)

replace (
	k8s.io/kubernetes => k8s.io/kubernetes v1.11.1
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.2.9
+	github.com/kubeflow/pipelines/api => ./api
)

```

## 次の手順について
[argo.go改修](../modify-argo.go)

## 補足事項
(※1) cloneしたリポジトリはtagのためブランチが設定されていない。ブランチを設定する場合は以下コマンドを実行し任意のブランチを作成し適用する。
```
$ git switch -c <branch name>
```
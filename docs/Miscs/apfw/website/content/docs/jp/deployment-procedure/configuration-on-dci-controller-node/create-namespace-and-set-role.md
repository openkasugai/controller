---
weight: 1
title: "Namespace作成とRole設定"
---
# Namespace作成とRole設定
## この手順について
* Custom Kubeflow Backendからカスタムリソースをデプロイする際には、デプロイ対象となるNamespaceの作成とロール設定が必要となる。
* 当該ドキュメントでは[Pipelineサンプル](../../../pipeline-sample)の実行を例として`ClusterRole`および`ClusterRoleBinding`の設定を記載する。
* 設定対象となるアカウントは`kubeflow`Namespaceの`default`アカウントとする。

## 対応内容
* Pipelineサンプルを実行する前提で、カスタムリソースをデプロイする対象となるNamespaceを作成する。
* `kubeflow`Namespaceの`default`アカウントが各Namespaceにリソースを作成できる権限設定(`ClusterRole`および`ClusterRoleBinding`)を行う。

## 前提条件
* 当該手順はOpenKasugaiコントローラがインストール済みである環境に対し実施する手順である。
* 環境情報については以下を参照。   
  [環境情報・前提条件](../../../environment-information-and-prerequisites)

## 手順
1. 作成されているNamespaceを確認する。(※1)
```
$ kubectl get namespace
```

2. Namespaceを作成する。
実行するPipelineのリソースのデプロイ先にあわせてNamespaceを作成する。(※2)
```
$ kubectl create namespace test01
$ kubectl create namespace wbfunc-imgproc
$ kubectl create namespace chain-imgproc
```

3. `ClusterRole`のYamlファイルを作成する。
```
$ vi swb-role.yaml
```

以下、マニフェストファイルを作成する。
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: swb-role
rules:
- apiGroups:
  - example.com
  resources:
  - dataflows
  - functionchains
  - functionkinds
  verbs:
  - "*"
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - "*"
```

4. 作成した`ClusterRole`をデプロイする。
```
$ kubectl apply -f swb-role.yaml
```

5. ClusterRole設定を確認する。
```
$ kubectl get clusterrole
```
以下のようにリソースが作成されたことを確認する。
```
NAME                                                                   CREATED AT
~~OMITTED~~
swb-role                                                               2024-03-01T01:00:45Z
```

6. `ClusterRoleBinding`のYamlファイルを作成する。
```
$ vi swb-role-binding.yaml
```
以下、マニフェストファイルを作成する。
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: swb-role-binding
subjects:
- kind: ServiceAccount
  name: default
  namespace: kubeflow
roleRef:
  kind: ClusterRole
  name: swb-role
  apiGroup: rbac.authorization.k8s.io
```

7. 作成した`ClusterRoleBinding`をデプロイする。
```
$ kubectl apply -f swb-role-binding.yaml
```

8. ClusterRoleBinding設定を確認する。
```
$ kubectl get clusterrolebinding
```
以下のようにリソースが作成されたことを確認する。
```
NAME                                                   ROLE                                      AGE
~~OMITTED~~
swb-role-binding                                       ClusterRole/swb-role                                      161d
```

## 次の手順について
* [Custom Kubeflow Backendデプロイ](../../apfw-backend-related-deployment-procedure/deploy-apfw-backend)

## 補足
(※1) OpenKasugaiコントローラ環境構築手順の内容によっては、すでに対象となるNamespaceが作成されている場合がある。次手順で作成するNamespaceがすでに作成済みの場合、Namespace作成手順は不要となる。  
(※2) 作成するNamespaceは実行するシナリオに合わせて変更する必要がある。今回は[Pipelineサンプル](../../../pipeline-sample)に合わせてNamespaceを作成している。
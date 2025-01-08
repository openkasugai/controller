---
weight: 1
title: "Custom Kubeflow Backendデプロイ"
---
# Custom Kubeflow Backendデプロイ
## この手順について
* Custom KubeflowからのIR YAML書式をマニフェストに変換する役割をになう`api-server`コンテナをデプロイする。

## 対応内容
* `api-server`コンテナのデプロイを行う。

## 前提条件
### デプロイ環境について
* デプロイに必要なソフトウェアおよびバージョンについては以下を参照
    * [環境情報・前提条件](../../../environment-information-and-prerequisites)

### 改修対象
* 改修対象となるKubeflowについては以下を参照
    * [改修対象と改修の流れについて](../../../about-the-target-of-modifications-and-the-flow-of-modifications)

### 事前手順
* 当該手順実施前に以下手順が実施済みであることを前提とする。
    * [Custom Kubeflow Backendビルド](../../../build-procedure/apfw-backend-related-build-procedure/build-apfw-backend)

## 手順
1. コンテナレジストリにログインする。
* `IP`と`PORT`はHarborのIPアドレスとポート番号を設定する。
```
$ docker login <IP>:<PORT>
```

2. コンテナイメージにタグを付ける。
* TAGは[Custom Kubeflow Backendビルド](../../../build-procedure/apfw-backend-related-build-procedure/build-apfw-backend)でコンテナイメージを作成した際のタグを指定する。
```
$ docker tag api-server:<TAG> <IP>:<PORT>/kfp/api-server:<TAG>
```

3. コンテナイメージをレジストリにPushする。
```
$ docker push <IP>:<PORT>/kfp/api-server:<TAG>
```

4. コンテナイメージを差し替える。
* namespaceは環境に合わせて指定する。
```
$ kubectl edit deploy ml-pipeline -n <namespace>
```

* `image`の値をレジストリにpushしたコンテナイメージに書き換える。
```
        - name: OBJECTSTORECONFIG_SECRETACCESSKEY
          valueFrom:
            secretKeyRef:
              key: secretkey
              name: mlpipeline-minio-artifact
        image: <IP>:<PORT>/kfp/api-server:<TAG>
        imagePullPolicy: IfNotPresent
```

5. PodのSTATUSがRunningになることを確認する。
```
$ kubectl get po -n <namespace> -w
NAME                                               READY   STATUS    RESTARTS       AGE
~~OMITTED~~
ml-pipeline-XXXXXXXXXX-XXXXX                       1/1     Running   0              8d
```

## 次の手順について
* [Custom Kubeflow SDKインストール手順](../../../operation-confirmation/apfw-sdk-operation-confirmation/apfw-sdk-install)
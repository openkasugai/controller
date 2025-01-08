---
weight: 2
title: "ライセンス情報の更新"
---
# ライセンス情報の更新
## この手順について
* `go-licenses`を使用し、Custom Kubeflow Backendに含まれるGoパッケージのライセンス情報を更新する。

## 対応内容
* `go-licenses`をインストールする。
* Custom Kubeflow Backendに含まれるGoパッケージのライセンス情報を更新する。

## 前提条件
### インストール環境について
* 当該手順をもとに`go-licenses`をインストールする場合、インストール環境には`Go`がインストールされている必要がある。
* インストールする`Go`のバージョン情報については以下を参照   
    * [環境情報・前提条件](../../../environment-information-and-prerequisites)

### 事前手順
* 当該手順実施前に以下手順が実施済みであることを前提とする。
    * [install-go-licenses.sh改修](../../../modification-procedure/apfw-backend-related-modification-procedure/modify-install-go-licenses.sh)
    * [Makefile改修](../../../modification-procedure/apfw-backend-related-modification-procedure/modify-makefile)

## 手順
1. `install-go-licenses.sh`を実行し`go-licenses`をインストールする。
```
$ bash hack/install-go-licenses.sh
```

2. Custom Kubeflow Backendに含まれるGoパッケージのライセンス情報を更新する。
```
$ make -C backend/ license_apiserver
```

## 次の手順について
* [Custom Kubeflow Backendビルド](../build-apfw-backend)
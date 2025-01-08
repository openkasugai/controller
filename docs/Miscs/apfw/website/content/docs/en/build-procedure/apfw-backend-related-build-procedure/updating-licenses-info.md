---
weight: 2
title: "Updating Licenses Info"
---
# Updating Licenses Info
## About This Procedure
* Use `go-licenses` to update the license information of Go packages included in Custom Kubeflow Backend.

## Changes Made
* Install `go-licenses`.
* Update the license information of Go packages included in Custom Kubeflow Backend.

## Prerequisites
### About the Installation Environment
* To install `go-licenses` based on this procedure, the installation environment must have `Go` installed.
* Refer to the version information of the installed `Go` for installation.
    * [Environment Information and Prerequisites](../../../environment-information-and-prerequisites)

### Pre-Procedure
* It is assumed that the following procedures have been completed before this procedure is implemented.
    * [Modify install-go-licenses.sh](../../../modification-procedure/apfw-backend-related-modification-procedure/modify-install-go-licenses.sh)
    * [Modify Makefile](../../../modification-procedure/apfw-backend-related-modification-procedure/modify-makefile)

## Procedure
1. Run `install-go-licenses.sh` to install `go-licenses`.
```
$ bash hack/install-go-licenses.sh
```

2. Update the license information of Go packages included in Custom Kubeflow Backend.
```
$ make -C backend/ license_apiserver
```

## Next Steps
* [Build Custom Kubeflow Backend](../build-apfw-backend)
---
weight: 6
title: "Modify install-go-licenses.sh"
---
# Modify install-go-licenses.sh
## About This Procedure
* Describes the modification points of Kubeflow Pipelines Backend to run Custom Kubeflow SDK.

## Changes Made
* Edit `install-go-licenses.sh` to change the version of `go-licenses` to be installed. (*1)

## Prerequisites
* Make sure the following procedures have been completed before performing this procedure.
    * [Modify v2_template.go](../modify-v2_template.go)

## Procedure
1. Edit `install-go-licenses.sh`.
```
$ vi hack/install-go-licenses.sh
```

* Edit as follows.
```sh
set -ex

# TODO: update to a released version.
-go install github.com/google/go-licenses@d483853
+go install github.com/google/go-licenses@706b9c60
```

## Next Steps
[Modify Makefile](../modify-makefile)

## Notes
(*1) 
* `go-licenses` analyzes the dependencies of Go packages, confirms the libraries used, and the licenses that can be used.
* By changing the version, the `ignore` option can be used to exclude the specified Go package from the confirmation by `go-licenses`.
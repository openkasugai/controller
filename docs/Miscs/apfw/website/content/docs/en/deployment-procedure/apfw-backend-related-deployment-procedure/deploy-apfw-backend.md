---
weight: 1
title: "Deploy Custom Kubeflow Backend"
---
# Deploy Custom Kubeflow Backend
## About This Procedure
* Deploy the `api-server` container, which converts IR YAML format from Custom Kubeflow into manifests.

## Changes Made
* Perform the deployment of the `api-server` container.

## Prerequisites
### About the Deploy Environment
* Refer to the following for the necessary software and versions for deployment:
    * [Environment Information and Prerequisites](../../../environment-information-and-prerequisites)

### Target for Modifications
* Refer to the following for the Kubeflow that is the target for modifications:
    * [About the Target of Modifications and the Flow of Modifications](../../../about-the-target-of-modifications-and-the-flow-of-modifications)

### Pre-Procedure
* It is assumed that the following steps have been completed before implementing this procedure:
    * [Custom Kubeflow Backend Build](../../../build-procedure/apfw-backend-related-build-procedure/build-apfw-backend)

## Procedure
1. Log in to the container registry.
* Set the `IP` and `PORT` to the IP address and port number of Harbor.
```
$ docker login <IP>:<PORT>
```

2. Tag the container image.
* Specify the TAG as the tag used when creating the container image in [Custom Kubeflow Backend Build](../../../build-procedure/apfw-backend-related-build-procedure/build-apfw-backend).
```
$ docker tag api-server:<TAG> <IP>:<PORT>/kfp/api-server:<TAG>
```

3. Push the container image to the registry.
```
$ docker push <IP>:<PORT>/kfp/api-server:<TAG>
```

4. Replace the container image.
* Specify the namespace according to the environment.
```
$ kubectl edit deploy ml-pipeline -n <namespace>
```

* Replace the value of `image` with the container image pushed to the registry.
```
        - name: OBJECTSTORECONFIG_SECRETACCESSKEY
          valueFrom:
            secretKeyRef:
              key: secretkey
              name: mlpipeline-minio-artifact
        image: <IP>:<PORT>/kfp/api-server:<TAG>
        imagePullPolicy: IfNotPresent
```

5. Confirm that the Pod's STATUS is Running.
```
$ kubectl get po -n <namespace> -w
NAME                                               READY   STATUS    RESTARTS       AGE
~~OMITTED~~
ml-pipeline-XXXXXXXXXX-XXXXX                       1/1     Running   0              8d
```

## Next Steps
* [Custom Kubeflow SDK Installation Procedure](../../../operation-confirmation/apfw-sdk-operation-confirmation/apfw-sdk-install)
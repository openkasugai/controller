---
weight: 1
title: "Create Namespace and Set Role"
---
# Create Namespace and Set Role
## About This Procedure
* When deploying custom resources from the Custom Kubeflow backend, it is necessary to create the Namespace where the resources will be deployed and set the roles.
* This document describes the configuration of `ClusterRole` and `ClusterRoleBinding` using the execution of the [Pipeline Sample](../../../pipeline-sample) as an example.
* The target account for configuration is the `default` account in the `kubeflow` Namespace.

## Changes Made
* In preparation for executing the Pipeline Sample, create the Namespace where the custom resources will be deployed.
* Set permissions for the `default` account in the `kubeflow` Namespace to create resources in each Namespace (`ClusterRole` and `ClusterRoleBinding`).

## Prerequisites
* This procedure is performed on an environment where the OpenKasugai controller is installed.
* Refer to the environment information below.
  [Environment Information and Prerequisites](../../../environment-information-and-prerequisites)

## Procedure
1. Check the existing Namespaces. (*1)
```
$ kubectl get namespace
```

2. Create a Namespace.
Create a Namespace according to where the resources of the Pipeline will be deployed. (*2)
```
$ kubectl create namespace test01
$ kubectl create namespace wbfunc-imgproc
$ kubectl create namespace chain-imgproc
```

3. Create a Yaml file for `ClusterRole`.
```
$ vi swb-role.yaml
```

Create a manifest file as follows.
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

4. Deploy the created `ClusterRole`.
```
$ kubectl apply -f swb-role.yaml
```

5. Verify the ClusterRole configuration.
```
$ kubectl get clusterrole
```

Confirm that the resources have been created as follows.
```
NAME                                                                   CREATED AT
~~OMITTED~~
swb-role                                                               2024-03-01T01:00:45Z
```

6. Create a Yaml file for `ClusterRoleBinding`.
```
$ vi swb-role-binding.yaml
```

Create a manifest file as follows.
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

7. Deploy the created `ClusterRoleBinding`.
```
$ kubectl apply -f swb-role-binding.yaml
```

8. Verify the ClusterRoleBinding configuration.
```
$ kubectl get clusterrolebinding
```

Confirm that the resources have been created as follows.
```
NAME                                                   ROLE                                      AGE
~~OMITTED~~
swb-role-binding                                       ClusterRole/swb-role                                      161d
```

## Next Steps
* [Deploy Custom Kubeflow Backend](../../apfw-backend-related-deployment-procedure/deploy-apfw-backend)

## Note
(*1) Depending on the content of the OpenKasugai controller environment construction procedure, the Namespace to be created may already exist. If the Namespace to be created in the next step has already been created, the Namespace creation procedure is not necessary.
(*2) The Namespace to be created needs to be adjusted according to the scenario to be executed. In this case, the Namespace is created according to the [Pipeline Sample](../../../pipeline-sample).
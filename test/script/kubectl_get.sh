#!/bin/bash -x
# Copyright 2024 NTT Corporation , FUJITSU LIMITED

kubectl get childbs.example.com -A
kubectl get computeresources.example.com
kubectl get connectiontypes.example.com -A
kubectl get connectiontargets.example.com -A
kubectl get dataflows.example.com -A
kubectl get deviceinfoes.example.com -A
kubectl get ethernetconnections.example.com -A
kubectl get fpgafunctions.example.com -A
kubectl get fpgas.example.com -A
kubectl get functionchains.example.com -A
kubectl get functiontypes.example.com -A
kubectl get functiontargets.example.com -A
kubectl get gpufunctions.example.com -A
kubectl get pcieconnections.example.com -A
kubectl get wbconnections.example.com -A
kubectl get wbfunctions.example.com -A
kubectl get cpufunctions.example.com -A
kubectl get topologyinfos.example.com -A
kubectl get pod -A -o wide
kubectl get schedulingdata.example.com -A
kubectl get configmap -A

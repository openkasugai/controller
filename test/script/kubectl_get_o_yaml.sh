#!/bin/bash -x
# Copyright 2024 NTT Corporation , FUJITSU LIMITED

kubectl get childbs.example.com -A -o yaml
kubectl get computeresources.example.com -A -o yaml
kubectl get connectiontypes.example.com -A -o yaml
kubectl get connectiontargets.example.com -A -o yaml
kubectl get dataflows.example.com -A -o yaml
kubectl get deviceinfoes.example.com -A -o yaml
kubectl get ethernetconnections.example.com -A -o yaml
kubectl get fpgafunctions.example.com -A -o yaml
kubectl get fpgas.example.com -A -o yaml
kubectl get functionchains.example.com -A -o yaml
kubectl get functiontypes.example.com -A -o yaml
kubectl get functiontargets.example.com -A -o yaml
kubectl get gpufunctions.example.com -A -o yaml
kubectl get pcieconnections.example.com -A -o yaml
kubectl get wbconnections.example.com -A -o yaml
kubectl get wbfunctions.example.com -A -o yaml
kubectl get cpufunctions.example.com -A -o yaml
kubectl get topologyinfos.example.com -A -o yaml
kubectl get pod -n test01 -o yaml
kubectl get schedulingdata.example.com -A -o yaml
kubectl get cm -n wbfunc-imgproc -o yaml
kubectl get cm -n test01 -o yaml
kubectl get cm -o yaml

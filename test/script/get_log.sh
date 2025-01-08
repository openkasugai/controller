#!/bin/bash -x
# Copyright 2024 NTT Corporation , FUJITSU LIMITED

CRC_LOG_DIR=./logs

if [ -d $CRC_LOG_DIR ];then
    rm -f $CRC_LOG_DIR/*
else
    mkdir $CRC_LOG_DIR
fi

kubectl get pod -A |grep -v NAME |grep crc- | awk '{print "kubectl logs "$2" -n "$1" > logs/"$2".log"}' | bash
kubectl get pod -A |grep -v NAME |grep controller-manager | awk '{print "kubectl logs "$2" -n "$1" > logs/"$2".log"}' | bash
kubectl get pod -A |grep -v NAME |grep df- | awk '{print "kubectl logs "$2" -n "$1" > logs/"$2".log"}' | bash
kubectl get pod -A |grep -v NAME |grep wbconnection- | awk '{print "kubectl logs "$2" -n "$1" > logs/"$2".log"}' | bash
kubectl get pod -A |grep -v NAME |grep wbfunction- | awk '{print "kubectl logs "$2" -n "$1" > logs/"$2".log"}' | bash

./kubectl_get.sh |& tee $CRC_LOG_DIR/kubectl_get.log 
./kubectl_get_o_yaml.sh |& tee $CRC_LOG_DIR/kubectl_get_o_yaml.log 

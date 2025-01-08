#!/bin/bash -x
# Copyright 2024 NTT Corporation , FUJITSU LIMITED

K8S_SOFT_DIR=$HOME/controller
YAML_DIR=$HOME/controller/test/sample-data/sample-data-common/yaml
TAG=1.0.0

cd $K8S_SOFT_DIR/src/whitebox-k8s-flowctrl
make deploy IMG=localhost/whitebox-k8s-flowctrl:$TAG

kubectl create namespace test01
kubectl create namespace wbfunc-imgproc
kubectl create namespace chain-imgproc
kubectl create namespace cluster01
kubectl create namespace topologyinfo

kubectl apply -f $YAML_DIR/functioninfo.yaml 
sleep 5
kubectl apply -f $YAML_DIR/functiontype.yaml
sleep 5
kubectl apply -f $YAML_DIR/functionchain.yaml
sleep 5

cd $K8S_SOFT_DIR/src/WBFunction
make deploy IMG=localhost/wbfunction:$TAG

cd $K8S_SOFT_DIR/src/WBConnection
make deploy IMG=localhost/wbconnection:$TAG

kubectl apply -f $K8S_SOFT_DIR/src/DeviceInfo/config/samples/crc_deviceinfo_daemonset.yaml
sleep 3
kubectl apply -f $K8S_SOFT_DIR/src/PCIeConnection/config/samples/crc_pcieconnection_daemonset.yaml 
kubectl apply -f $K8S_SOFT_DIR/src/EthernetConnection/config/samples/crc_ethernetconnection_daemonset.yaml 
kubectl apply -f $K8S_SOFT_DIR/src/FPGAFunction/config/samples/crc_fpgafunction_daemonset.yaml 
kubectl apply -f $K8S_SOFT_DIR/src/GPUFunction/config/samples/crc_gpufunction_daemonset.yaml
kubectl apply -f $K8S_SOFT_DIR/src/CPUFunction/config/samples/crc_cpufunction_daemonset.yaml

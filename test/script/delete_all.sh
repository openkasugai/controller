#!/bin/bash -x
# Copyright 2025 NTT Corporation , FUJITSU LIMITED

YAML_DIR=$HOME/controller/test/sample-data/sample-data-common/yaml
WBFUNCTION_NAMESPACE=test01

kubectl get dataflow -A |grep -v NAME | awk '{system ("kubectl delete dataflow " $2 " -n " $1)}'
sleep 5
kubectl get computeresource -A |grep -v NAME | awk '{system ("kubectl delete computeresource " $2 " -n " $1)}'
kubectl get fpga -A |grep -v NAME | awk '{system ("kubectl delete fpga " $2 " -n " $1)}'

kubectl delete -f $YAML_DIR/functioninfo.yaml
kubectl delete -f $YAML_DIR/functiontype.yaml
kubectl delete -f $YAML_DIR/functionchain.yaml

sleep 5

## Delete WBFunction related finalizer
for PREFIX in fpga cpu gpu wb; do
  FUNCTION_NAMES=$(kubectl get ${PREFIX}functions -n $WBFUNCTION_NAMESPACE --no-headers -o custom-columns=":metadata.name")
  for FUNCTION_NAME in $FUNCTION_NAMES; do
    kubectl patch ${PREFIX^^}Function $FUNCTION_NAME -n $WBFUNCTION_NAMESPACE --type=json -p '[{"op": "remove", "path": "/metadata/finalizers"}]'
  done
done

K8S_SOFT_DIR=$HOME/controller
kubectl delete -f $K8S_SOFT_DIR/src/DeviceInfo/config/samples/crc_deviceinfo_daemonset.yaml
kubectl delete -f $K8S_SOFT_DIR/src/PCIeConnection/config/samples/crc_pcieconnection_daemonset.yaml
kubectl delete -f $K8S_SOFT_DIR/src/EthernetConnection/config/samples/crc_ethernetconnection_daemonset.yaml
kubectl delete -f $K8S_SOFT_DIR/src/FPGAFunction/config/samples/crc_fpgafunction_daemonset.yaml
kubectl delete -f $K8S_SOFT_DIR/src/GPUFunction/config/samples/crc_gpufunction_daemonset.yaml
kubectl delete -f $K8S_SOFT_DIR/src/CPUFunction/config/samples/crc_cpufunction_daemonset.yaml
cd $K8S_SOFT_DIR/src/WBFunction && make undeploy
cd $K8S_SOFT_DIR/src/WBConnection && make undeploy
cd $K8S_SOFT_DIR/src/whitebox-k8s-flowctrl && make undeploy

CHILDBS_NAMES=$(kubectl get childbs -n default --no-headers -o custom-columns=":metadata.name")
for CHILDBS_NAME in $CHILDBS_NAMES; do
  kubectl patch childbs $CHILDBS_NAME --type=json -p '[{"op": "remove", "path": "/metadata/finalizers"}]'
done
kubectl get childbs -A |grep -v NAME | awk '{system ("kubectl delete childbs " $2 " -n " $1)}'

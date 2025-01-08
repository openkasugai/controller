#!/bin/sh

echo 
echo "============================================"
echo "FunctionTargets"
echo "============================================"
kubectl get functiontargets.example.com -A

echo 
echo "============================================"
echo "DataFlow"
echo "============================================"
kubectl get dataflows.example.com -A

echo 
echo "============================================"
echo "DataFlow (detail)"
echo "============================================"
kubectl get dataflow -o yaml test

# echo 
# echo "============================================"
# echo "WBFunction (detail)"
# echo "============================================"
# kubectl get wbfunctions.example.com -o yaml

# echo 
# echo "============================================"
# echo "WBConnection (detail)"
# echo "============================================"
# kubectl get wbconnections.example.com -o yaml

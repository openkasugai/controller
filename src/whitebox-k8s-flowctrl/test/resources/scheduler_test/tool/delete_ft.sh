#!/bin/bash
# Copyright 2024 NTT Corporation , FUJITSU LIMITED

kubectl get functiontargets.example.com -A |grep -v NAME | awk '{system ("kubectl delete functiontargets.example.com " $2 " -n " $1)}'

# SCRIPT_DIR=$(cd $(dirname $0); pwd)

# files="./functiontargets/functiontarget-$1_FT*.yaml"
# for file in $files; do
#     kubectl delete -f $file
# done

#!/bin/bash
# Copyright 2024 NTT Corporation , FUJITSU LIMITED

SCRIPT_DIR=$(cd $(dirname $0); pwd)

files="./functiontargets/functiontarget-$1_FT*.yaml"
for file in $files; do
    kubectl apply -f $file
    $SCRIPT_DIR/../tool/update_status.sh apply -f $file
done

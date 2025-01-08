#!/usr/bin/bash
# Copyright 2024 NTT Corporation , FUJITSU LIMITED

PROXY_PORT=55555

# Directory where Python scripts are stored
CMD_DIR="$(dirname $0)/python"

PARSE_ARGS_CMD="${CMD_DIR}/parse_args.py"
GET_STATUS_YAML_CMD="${CMD_DIR}/get_status_yaml.py"
YAML2JSON_CMD="${CMD_DIR}/yaml2json.py"
GET_RESOURCE_CMD=${CMD_DIR}/"get_resource_info_from_yaml.py"

args=(`${PARSE_ARGS_CMD} $@`)
if [ $? -ne 0 ]; then
	exit 1
fi

if [[ ${args[0]} = usage* ]]; then
	${PARSE_ARGS_CMD} $@
	exit 0
fi

# Set file name
if [ ${#args[@]} -eq 1 ]; then
	filename=${args[0]}
elif [ ${#args[@]} -eq 4 ]; then
	filename=${args[3]}
fi

# Set kind, name, namespace
if [ ${#args[@]} -eq 1 ]; then
	resourceinfo=($(${GET_RESOURCE_CMD} -f ${filename}))
	if [ $? -ne 0 ]; then
		exit 1
	fi
	kind=${resourceinfo[0]}
	name=${resourceinfo[1]}
	namespace=${resourceinfo[2]}
else
	kind=${args[0]}
	name=${args[1]}
	namespace=${args[2]}
fi

if [ ${#args[@]} -eq 3 ]; then
	status_json=`cat | ${GET_STATUS_YAML_CMD} | ${YAML2JSON_CMD}`
else
	status_json=`${GET_STATUS_YAML_CMD} -f ${filename} | ${YAML2JSON_CMD}`
fi
if [ $? -eq 1 ];  then
	echo -e "$status_json"
	exit 1
fi

echo "kind="$kind
kind=$(echo $kind | sed -r 's/.+/\L\0/')
if [[ $kind = *s ]]; then
	kind=${kind}es
else
	kind=${kind}s
fi

group=$(kubectl get crd | grep $kind | awk '{print $1}' | sed -r 's/'$kind'\.(.+)/\1/')

kubectl proxy --port=$PROXY_PORT &
pid=$!

sleep 1

curl -v -XPATCH -H "Content-Type: application/merge-patch+json" -d "
$(echo $status_json)
" http://localhost:${PROXY_PORT}/apis/${group}/v1/namespaces/${namespace}/${kind}/${name}/status

kill $!

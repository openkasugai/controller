#!/usr/bin/bash -x
if [ $1 == "delete" ]; then
	kubectl delete configmap servicer-mgmt-info
	kubectl delete configmap functionkindmap
	kubectl delete configmap connectionkindmap
	kubectl delete configmap region-unique-info
	kubectl delete configmap function-unique-info
	kubectl delete configmap fpgafunc-config-filter-resize-high-infer
	kubectl delete configmap fpgafunc-config-filter-resize-low-infer
	kubectl delete configmap cpufunc-config-decode
	kubectl delete configmap cpufunc-config-glue-fdma-to-tcp
	kubectl delete configmap cpufunc-config-copy-branch
	kubectl delete configmap cpufunc-config-filter-resize-high-infer
	kubectl delete configmap cpufunc-config-filter-resize-low-infer
	kubectl delete configmap gpufunc-config-high-infer
	kubectl delete configmap gpufunc-config-low-infer
elif [ $1 == "create" ]; then
	kubectl create configmap servicer-mgmt-info --from-file=servicer-mgmt-info.json
	kubectl create configmap functionkindmap --from-file=functionkindmap.json
	kubectl create configmap connectionkindmap --from-file=connectionkindmap.json
	kubectl create configmap region-unique-info --from-file=region-unique-info.json
	kubectl create configmap function-unique-info --from-file=function-unique-info.json
	kubectl create configmap fpgafunc-config-filter-resize-high-infer --from-file=fpgafunc-config-filter-resize-high-infer.json
	kubectl create configmap fpgafunc-config-filter-resize-low-infer --from-file=fpgafunc-config-filter-resize-low-infer.json
	kubectl create configmap cpufunc-config-decode --from-file=cpufunc-config-decode.json
	kubectl create configmap cpufunc-config-glue-fdma-to-tcp --from-file=cpufunc-config-glue-fdma-to-tcp.json
	kubectl create configmap cpufunc-config-copy-branch --from-file=cpufunc-config-copy-branch.json
	kubectl create configmap cpufunc-config-filter-resize-high-infer --from-file=cpufunc-config-filter-resize-high-infer.json
	kubectl create configmap cpufunc-config-filter-resize-low-infer --from-file=cpufunc-config-filter-resize-low-infer.json
	kubectl create configmap gpufunc-config-high-infer --from-file=gpufunc-config-high-infer.json
	kubectl create configmap gpufunc-config-low-infer --from-file=gpufunc-config-low-infer.json
	kubectl create namespace test01
elif [ $1 == "get" ]; then
	kubectl get configmap
else
	echo "parameter error. (1st-param:create or delete or get)"
fi


#!/usr/bin/env bash
set +e
echo -e "==== CPUFunction ===="
golangci-lint run ./CPUFunction/...
RESULT_CPUFunction=$?
echo -e "\n\n==== DeviceInfo ===="
golangci-lint run ./DeviceInfo/...
RESULT_DeviceInfo=$?
echo -e "\n\n==== EthernetConnection ===="
golangci-lint run ./EthernetConnection/...
RESULT_EthernetConnection=$?
echo -e "\n\n==== FPGAFunction ===="
golangci-lint run ./FPGAFunction/...
RESULT_FPGAFunction=$?
echo -e "\n\n==== GPUFunction ===="
golangci-lint run ./GPUFunction/...
RESULT_GPUFunction=$?
echo -e "\n\n==== PCIeConnection ===="
golangci-lint run ./PCIeConnection/...
RESULT_PCIeConnection=$?
echo -e "\n\n==== WBConnection ===="
golangci-lint run ./WBConnection/...
RESULT_WBConnection=$?
echo -e "\n\n==== WBFunction ===="
golangci-lint run ./WBFunction/...
RESULT_WBFunction=$?
echo -e "\n\n==== whitebox-k8w-flowctrl ===="
golangci-lint run ./whitebox-k8s-flowctrl/...
RESULT_WBCTRL=$?
echo -e "\n\n==== tools/FPGAReconfigurationTool ===="
golangci-lint run ./tools/FPGAReconfigurationTool/...
RESULT_tools_FPGAReconfigurationTool=$?
echo -e "\n\n==== tools/FPGAClearCheckTool/FPGACheckPerDF ===="
golangci-lint run ./tools/FPGAClearCheckTool/FPGACheckPerDF/...
RESULT_tools_FPGAClearCheckTool_FPGACheckPerDF=$?
echo -e "\n\n==== tools/gpu_info ===="
golangci-lint run ./tools/gpu_info/...
RESULT_tools_gpu_info=$?
echo -e "\n\n==== tools/InfoCollector ===="
golangci-lint run ./tools/InfoCollector/...
RESULT_tools_InfoCollector=$?

ERR=$(($RESULT_CPUFunction + \
    $RESULT_DeviceInfo + \
    $RESULT_EthernetConnection + \
    $RESULT_FPGAFunction + \
    $RESULT_GPUFunction + \
    $RESULT_PCIeConnection + \
    $RESULT_WBConnection + \
    $RESULT_WBFunction + \
    $RESULT_WBCTRL + \
    $RESULT_tools_FPGAReconfigurationTool + \
    $RESULT_tools_FPGAClearCheckTool_FPGACheckPerDF + \
    $RESULT_tools_gpu_info + \
    $RESULT_tools_InfoCollector
))

if [ $ERR -eq 0 ]; then
    echo -e "\n\nAll linter passed."
    exit 0
else
    echo -e "\n\nPlease check the error(s)."
    exit 1
fi

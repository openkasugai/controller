/*
Copyright 2024 NTT Corporation, FUJITSU LIMITED

InfoCollector functions
*/

package infocollect

import (
	"context"
	"encoding/json"
	cm "example.com/InfoCollector/pkg/configmap"
	examplecomv1 "example.com/InfoCollector/pkg/reference"
	"fmt"
	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"os/exec"
	"regexp"
	ctrl "sigs.k8s.io/controller-runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"

	// #cgo CFLAGS: -I/usr/local/include/fpgalib/
	// #cgo CFLAGS: -I/usr/include/
	// #cgo LDFLAGS: -L/usr/local/lib/fpgalib/ -lfpga
	// #cgo LDFLAGS: -L/usr/local/lib/fpgalib/ -lfpgadb
	// #cgo LDFLAGS: -L/usr/lib/gcc/x86_64-linux-gnu/9/ -lstdc++
	// #cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu/ -lpciaccess
	// #include <libfpgactl.h>
	// #include <libfunction.h>
	// #include <libfpgadb.h>
	// #include <stdlib.h>
	"C"
)

const (
	inputFilePath      = "infrainfo/"
	fileList           = "premadefilelist.json"
	regionUniqueInfo   = "region-unique-info"
	functionUniqueInfo = "function-unique-info"
	decodeChildBs      = "decode-childbs"
	resizeChildBs      = "filter-resize-childbs"
	servicerMgmtInfo   = "servicer-mgmt-info"
)

const (
	returnError = -1
)

var cpuNum int32 = 1

/* Provisional support (dynamic reconfiguration not supported) */
var gFunctionNameMap map[string][]FunctionNameMap
var gDeviceTypeMap []DeviceTypeMap

/* Provisional support (dynamic reconfiguration not supported) */
// FuncName conversion function
func convFunctionName(
	ctx context.Context,
	height uint32,
	width uint32,
	functionName string) string {

	const (
		functionNameMap = "functionnamemap.json"
	)

	var err error
	err = nil
	var bytes []uint8

	logger := ctxzap.Extract(ctx)

	for doWhile := 0; doWhile < 1; doWhile++ {

		if "decode" == functionName {
			break
		}

		if 0 == len(gFunctionNameMap["sizeList"]) {
			// Get the FuncName conversion file
			bytes, err = os.ReadFile(inputFilePath + functionNameMap)
			if err != nil {
				logger.Error("unable to readfile. FileName = "+functionNameMap, zap.Error(err))
				break
			}
			// JSON decoding
			if err = json.Unmarshal(bytes, &gFunctionNameMap); err != nil {
				logger.Error("unable to unmarshal. FileName = "+functionNameMap, zap.Error(err))
				break
			}
		}

		// Convert FuncName based on size
		for i := 0; i < len(gFunctionNameMap["sizeList"]); i++ {
			sizeList := gFunctionNameMap["sizeList"][i]
			if sizeList.Height == height && sizeList.Width == width {
				functionName = functionName + sizeList.FunctionName
				break
			}
		}
	}

	return functionName
}

/* Provisional support (dynamic reconfiguration not supported) */
// DeviceType conversion function
func convDeviceType(ctx context.Context, deviceType string) string {

	const (
		deviceTypeMap = "devicetypemap.json"
	)

	var err error
	err = nil
	var bytes []uint8

	logger := ctxzap.Extract(ctx)

	for doWhile := 0; doWhile < 1; doWhile++ {

		if 0 == len(gDeviceTypeMap) {
			// Get the DeviceType file
			bytes, err = os.ReadFile(inputFilePath + deviceTypeMap)
			if err != nil {
				logger.Error("unable to readfile. FileName = "+deviceTypeMap, zap.Error(err))
				break
			}
			// JSON decoding
			if err = json.Unmarshal(bytes, &gDeviceTypeMap); err != nil {
				logger.Error("unable to unmarshal. FileName = "+deviceTypeMap, zap.Error(err))
				break
			}
		}

		// Convert DeviceType
		for i := 0; i < len(gDeviceTypeMap); i++ {
			deviceTypeTable := gDeviceTypeMap[i]
			if deviceTypeTable.InputDeviceType == deviceType {
				deviceType = deviceTypeTable.OutputDeviceType
				break
			}
		}
	}

	return deviceType
}

// Pre-deployment information acquisition function
func GetJsonFile(
	ctx context.Context,
	pRegionSpecifics *[]RegionSpecificInfo,
	pFPGACatalogs *[]FPGACatalog,
	pFunctionDedicatedDecodeInfo *map[string][]FunctionDedicatedInfo,
	pFunctionDedicatedResizeInfo *map[string][]FunctionDedicatedInfo,
	pServicerMgmtInfo *[]ServicerMgmtInfo) error {

	var filemap map[string]interface{}
	var err error
	err = nil
	var bytes []uint8

	logger := ctxzap.Extract(ctx)

	for doWhile := 0; doWhile < 1; doWhile++ {
		// Get the predeployed file name
		if bytes, err = os.ReadFile(inputFilePath + fileList); err != nil {
			logger.Error("unable to readfile. FileName = "+fileList, zap.Error(err))
			break
		}
		// JSON decoding
		if err = json.Unmarshal(bytes, &filemap); err != nil {
			logger.Error("unable to unmarshal. FileName = "+fileList, zap.Error(err))
			break
		}
		for key, val := range filemap {
			filePath := inputFilePath + val.(string)
			if bytes, err = os.ReadFile(filePath); err != nil {
				logger.Error("unable to readfile. FileName = "+val.(string), zap.Error(err))
				break
			}
			if regionUniqueInfo == key {
				// Get domain-specific information
				err = json.Unmarshal(bytes, &pRegionSpecifics)
			} else if functionUniqueInfo == key {
				// Function-specific information - Get common attributes
				err = json.Unmarshal(bytes, &pFPGACatalogs)
			} else if decodeChildBs == key {
				// Function-specific information - get decode-only attributes
				err = json.Unmarshal(bytes, &pFunctionDedicatedDecodeInfo)
			} else if resizeChildBs == key {
				// Function-specific information - get filter/resize-specific attributes
				err = json.Unmarshal(bytes, &pFunctionDedicatedResizeInfo)
			} else if servicerMgmtInfo == key {
				// Get address-specific information
				err = json.Unmarshal(bytes, &pServicerMgmtInfo)
			}

			if nil != err {
				logger.Error("unable to unmarshal. FileName = "+val.(string), zap.Error(err))
				break
			}
		}
	}
	return err
}

// Device information acquisition function
func GetDeviceInfo(
	ctx context.Context,
	pDevices *[]DeviceBasicInfo) int {

	var err error
	err = nil
	var ret = 0
	var retCInt C.int
	var deviceListCChar **C.char
	var fpgaDeviceUUIDs []string
	var deviceUUID string
	var deviceID C.uint
	var fpgaDeviceUserInfo C.fpga_device_user_info_t
	var gpuIDs []uint
	var gpuDeviceInfo dcgm.Device
	var cpuModelName []byte
	var cpuNumByte []byte
	var cpuDeviceType string

	logger := ctxzap.Extract(ctx)

	for doWhile := 0; doWhile < 1; doWhile++ {

		// Get the node name
		if cm.GMyNodeName, err = os.Hostname(); err != nil {
			logger.Error("HostName Get Error.", zap.Error(err))
			ret = returnError
			break
		}
		// Call the device scan function
		retCInt = C.fpga_scan_devices()
		if 0 > retCInt {
			logger.Error("fpga_scan_devices() err = " + strconv.Itoa(int(retCInt)))
		} else if 0 == retCInt {
			logger.Info("fpga_scan_devices() get fpga num = " + strconv.Itoa(int(retCInt)))
		} else {
			logger.Info("fpga_scan_devices() get fpga num = " + strconv.Itoa(int(retCInt)))

			// Call the device list acquisition function
			retCInt = C.fpga_get_device_list((***C.char)(unsafe.Pointer(&deviceListCChar))) //nolint:gocritic // suspicious indentical LHS and RHS for next block "==". QUESTION: why?
			if 0 > retCInt {
				logger.Error("fpga_get_device_list() err = " + strconv.Itoa(int(retCInt)))
			} else {
				logger.Info("fpga_get_device_list() ret = " + strconv.Itoa(int(retCInt)))

				// Get the FPGADeviceUUID and its number from the return value of the device list acquisition function
				deviceListHeadPointer := *(**byte)(unsafe.Pointer(deviceListCChar))

				// Get the number of characters in the first FPGADeviceUUID in device_list
				n := 0
				for ptr := deviceListHeadPointer; *ptr != 0; n++ {
					deviceUUID = deviceUUID + string(*ptr)
					ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + 1))
				}
				// Get the first FPGADeviceUUID in device_list
				fpgaDeviceUUIDs = append(fpgaDeviceUUIDs, deviceUUID)

				// Repeat for each array in device_list
				for addressIncrementNum := 0x08; ; addressIncrementNum += 0x08 {

					deviceUUID = ""

					nextHeadPointer := *(**byte)(unsafe.Pointer(uintptr(
						unsafe.Pointer(deviceListCChar)) + uintptr(addressIncrementNum)))
					// Check if the next element in device_list has FPGADeviceUUID
					if 0x20 != (uintptr(unsafe.Pointer(nextHeadPointer)) -
						(uintptr(unsafe.Pointer(deviceListHeadPointer)))) {
						break
					}
					n := 0
					for ptr := nextHeadPointer; *ptr != 0; n++ {
						deviceUUID = deviceUUID + string(*ptr)
						ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + 1))
					}

					fpgaDeviceUUIDs = append(fpgaDeviceUUIDs, deviceUUID)
					deviceListHeadPointer = nextHeadPointer
				}
			}

			// Get FPGA device information
			for _, fpgaDeviceUUID := range fpgaDeviceUUIDs {

				fpgaDeviceUUIDCString := C.CString(fpgaDeviceUUID)
				defer C.free(unsafe.Pointer(fpgaDeviceUUIDCString))
				// Call the device ID acquisition function
				retCInt = C.fpga_get_dev_id(fpgaDeviceUUIDCString, &deviceID)
				if 0 > retCInt {
					logger.Error("fpga_get_dev_id() err = " + strconv.Itoa(int(retCInt)))
					break
				} else {
					logger.Info("fpga_get_dev_id() ret = " + strconv.Itoa(int(retCInt)))
				}
				// Call to get device information
				retCInt = C.fpga_get_device_info(deviceID, (*C.fpga_device_user_info_t)(
					unsafe.Pointer(&fpgaDeviceUserInfo)))
				if 0 > retCInt {
					logger.Error("fpga_get_device_info() err = " + strconv.Itoa(int(retCInt)))
					break
				} else {
					logger.Info("fpga_get_device_info() ret = " + strconv.Itoa(int(retCInt)))
				}

				var deviceBasicData DeviceBasicInfo
				deviceUUID := fpgaDeviceUUID
				deviceFilePath := C.GoString((*C.char)(unsafe.Pointer(&fpgaDeviceUserInfo.device_file_path[0])))
				deviceVendor := C.GoString((*C.char)(unsafe.Pointer(&fpgaDeviceUserInfo.vendor[0])))
				pciDomain := int32(fpgaDeviceUserInfo.pcie_bus.domain)
				pciBus := int32(fpgaDeviceUserInfo.pcie_bus.bus)
				pciDevice := int32(fpgaDeviceUserInfo.pcie_bus.device)
				pciFunction := int32(fpgaDeviceUserInfo.pcie_bus.function)

				// Add FPGA information to the device information structure
				deviceBasicData.NodeName = cm.GMyNodeName
				deviceBasicData.DeviceUUID = &deviceUUID
				deviceBasicData.DeviceFilePath = &deviceFilePath
				deviceBasicData.DeviceType = C.GoString((*C.char)(
					unsafe.Pointer(&fpgaDeviceUserInfo.device_type[0])))
				deviceBasicData.DeviceIndex = int32(fpgaDeviceUserInfo.device_index)
				deviceBasicData.ParentID = fmt.Sprintf("%08s", strconv.FormatInt(
					int64(fpgaDeviceUserInfo.bitstream_id.parent), 16))
				deviceBasicData.SubDeviceSpecRef = fmt.Sprintf("%08s", strconv.FormatInt(
					int64(fpgaDeviceUserInfo.bitstream_id.child), 16))
				deviceBasicData.DeviceVendor = &deviceVendor
				deviceBasicData.PCIDomain = &pciDomain
				deviceBasicData.PCIBus = &pciBus
				deviceBasicData.PCIDevice = &pciDevice
				deviceBasicData.PCIFunction = &pciFunction
				deviceBasicData.DeviceCategory = "FPGA"

				var dummyChildBsID C.uint32_t = 0xfffffff0
				var parentIDUint C.uint32_t
				var childIDUint C.uint32_t
				uuidCString := C.CString(deviceUUID)
				defer C.free(unsafe.Pointer(uuidCString))
				// Call the dummy data enable function
				ret := C.fpga_db_enable_dummy_bitstream(uuidCString, nil, &dummyChildBsID)
				if 0 > ret {
					logger.Error("fpga_db_enable_dummy_bitstream() err = " +
						strconv.Itoa(int(ret)))
				} else {
					logger.Info("fpga_db_enable_dummy_bitstream() ret = " +
						strconv.Itoa(int(ret)))
				}

				// Call the bitstreamID acquisition function
				ret = C.fpga_db_get_bitstream_id(uuidCString, &parentIDUint, &childIDUint)
				if 0 > ret {
					logger.Error("fpga_db_get_bitstream_id() err = " +
						strconv.Itoa(int(ret)))
				} else {
					logger.Info("fpga_db_get_bitstream_id() ret = " +
						strconv.Itoa(int(ret)))
				}

				// Call the dummy data disable function
				ret = C.fpga_db_disable_dummy_bitstream(uuidCString)
				if 0 > ret {
					logger.Error("fpga_db_disable_dummy_bitstream() err = " +
						strconv.Itoa(int(ret)))
				} else {
					logger.Info("fpga_db_disable_dummy_bitstream() ret = " +
						strconv.Itoa(int(ret)))
				}

				deviceBasicData.ParentID = fmt.Sprintf("%08s", strconv.FormatInt(
					int64(parentIDUint), 16))
				deviceBasicData.SubDeviceSpecRef = fmt.Sprintf("%08s", strconv.FormatInt(
					int64(childIDUint), 16))

				*pDevices = append(*pDevices, deviceBasicData)
			}

			// Call the device list release function
			retCInt = C.fpga_release_device_list((**C.char)(
				unsafe.Pointer(deviceListCChar)))
			if 0 > retCInt {
				logger.Error("fpga_release_device_list() err = " +
					strconv.Itoa(int(retCInt)))
			} else {
				logger.Info("fpga_release_device_list() ret = " +
					strconv.Itoa(int(retCInt)))
			}
		}

		// Call DCGM initialization function
		cleanup, dcgmerr := dcgm.Init(dcgm.Embedded)
		if nil != dcgmerr {
			logger.Info("dcgm.Init() Error but Maybe there are NOT any GPU.", zap.Error(dcgmerr))
		} else {
			defer cleanup()

			// Call the GPUID acquisition function
			if gpuIDs, err = dcgm.GetSupportedDevices(); nil != err {
				logger.Error("GetSupportedDevices() Error.", zap.Error(err))
			}
			// Get GPU device information
			for _, gpuId := range gpuIDs {
				// Call device information acquisition function
				if gpuDeviceInfo, err = dcgm.GetDeviceInfo(gpuId); nil != err {
					logger.Error("GetDeviceInfo() Error.", zap.Error(err))
				}

				var deviceBasicData DeviceBasicInfo
				deviceUUID := gpuDeviceInfo.UUID

				// Add GPU information to the device information structure
				deviceBasicData.NodeName = cm.GMyNodeName
				deviceBasicData.DeviceUUID = &deviceUUID
				deviceBasicData.DeviceType = gpuDeviceInfo.Identifiers.Model
				deviceBasicData.DeviceIndex = int32(gpuDeviceInfo.GPU)
				deviceBasicData.ParentID = ""
				deviceBasicData.SubDeviceSpecRef = gpuDeviceInfo.Identifiers.Model
				deviceBasicData.DeviceCategory = "GPU"
				*pDevices = append(*pDevices, deviceBasicData)
			}
		}

		// CPU information acquisition process
		reg := "[:]"

		cpuModelName, err = exec.Command("sh", "-c",
			"grep model.name /proc/cpuinfo | sort -u").Output()
		if nil != err {
			logger.Info("unable to get cpuinfo but Maybe there are NOT any CPU.", zap.Error(err))
		} else {
			arr := regexp.MustCompile(reg).Split(string(cpuModelName), -1)
			for _, st := range arr {
				if true != strings.Contains(st, "model name") {
					cpuDeviceType = strings.Trim(st, "\n\r\t ")
				}
			}

			cpuNumByte, _ = exec.Command("sh", "-c",
				"fgrep 'physical id' /proc/cpuinfo | sort -u | wc -l").Output()
			cpuNumInt, _ := strconv.Atoi(string(cpuNumByte[:len(cpuNumByte)-1]))
			cpuNum = int32(cpuNumInt)

			var deviceBasicData DeviceBasicInfo
			deviceUUID := cm.GMyNodeName + "-cpu0"

			// Add CPU information to the device information structure
			deviceBasicData.NodeName = cm.GMyNodeName
			deviceBasicData.DeviceUUID = &deviceUUID
			deviceBasicData.DeviceType = cpuDeviceType
			deviceBasicData.DeviceIndex = 0
			deviceBasicData.ParentID = ""
			deviceBasicData.SubDeviceSpecRef = cpuDeviceType
			deviceBasicData.DeviceCategory = "CPU"
			*pDevices = append(*pDevices, deviceBasicData)
		}

		// MEM information acquisition process
		var deviceBasicData DeviceBasicInfo
		deviceUUID := cm.GMyNodeName + "-mem0"

		// Add CPU information to the device information structure
		deviceBasicData.NodeName = cm.GMyNodeName
		deviceBasicData.DeviceUUID = &deviceUUID
		deviceBasicData.DeviceType = "memory"
		deviceBasicData.DeviceIndex = 0
		deviceBasicData.ParentID = ""
		deviceBasicData.SubDeviceSpecRef = "memory"
		deviceBasicData.DeviceCategory = "MEM"
		*pDevices = append(*pDevices, deviceBasicData)
	}

	return ret
}

// fpga information creation function
func CreateFPGACR(
	ctx context.Context,
	mng ctrl.Manager,
	pRegionSpecifics *[]RegionSpecificInfo,
	pDevices *[]DeviceBasicInfo) error {

	var err error
	err = nil

	logger := ctxzap.Extract(ctx)
	c := mng.GetClient()

	// Repeat for the number of devices in the environment information
	for deviceIndex := 0; deviceIndex < len(*pDevices); deviceIndex++ {

		deviceData := (*pDevices)[deviceIndex]
		var fpgaCRData examplecomv1.FPGA
		var forError error
		var jsonstr []byte

		if "FPGA" != deviceData.DeviceCategory {
			continue
		}

		// Add information to the node and device information structure
		fpgaCRData.Spec.DeviceIndex = deviceData.DeviceIndex
		fpgaCRData.Spec.DeviceFilePath = *deviceData.DeviceFilePath
		fpgaCRData.Spec.DeviceUUID = *deviceData.DeviceUUID
		fpgaCRData.Spec.NodeName = deviceData.NodeName
		fpgaCRData.Spec.ParentBitstreamID = deviceData.ParentID
		fpgaCRData.Spec.PCIDomain = *deviceData.PCIDomain
		fpgaCRData.Spec.PCIBus = *deviceData.PCIBus
		fpgaCRData.Spec.PCIDevice = *deviceData.PCIDevice
		fpgaCRData.Spec.PCIFunction = *deviceData.PCIFunction
		fpgaCRData.Spec.Vendor = *deviceData.DeviceVendor

		// Repeat for the number of lists of domain-specific information
		for regionSpecificIndex := 0; regionSpecificIndex < len(*pRegionSpecifics); regionSpecificIndex++ {
			regionSpecificInfo := (*pRegionSpecifics)[regionSpecificIndex]
			functionTargets := regionSpecificInfo.FunctionTargets

			if deviceData.SubDeviceSpecRef != regionSpecificInfo.SubDeviceSpecRef {
				continue
			}

			// Repeat for the number of FunctionTargets in the domain-specific information
			for functionTargetIndex := 0; functionTargetIndex < len(functionTargets); functionTargetIndex++ {
				functionTargetData := functionTargets[functionTargetIndex]
				if nil != functionTargetData.MaxFunctions {
					fpgaCRData.Spec.ChildBitstreamID = &deviceData.SubDeviceSpecRef
					break
				}
			}
		}

		jsonstr, forError = json.Marshal(fpgaCRData.Spec)
		if nil != forError {
			logger.Error("FPGA.Spec unable to Marshal.", zap.Error(forError))
			break
		}
		forError = json.Unmarshal(jsonstr, &fpgaCRData.Status)
		if nil != forError {
			logger.Error("FPGA.Status unable to Unmarshal.", zap.Error(forError))
			break
		}
		fpgaCRData.Status.Status = examplecomv1.FPGAStatusReady

		crData := &unstructured.Unstructured{}
		crData.SetGroupVersionKind(schema.GroupVersionKind{
			Version: "example.com/v1",
			Kind:    "FPGA",
		})
		crData.SetName(strings.ToLower("fpga-" + fpgaCRData.Spec.DeviceUUID + "-" + fpgaCRData.Spec.NodeName))
		crData.SetNamespace("default")

		crData.UnstructuredContent()["spec"] = fpgaCRData.Spec
		crData.UnstructuredContent()["status"] = fpgaCRData.Status

		forError = c.Create(ctx, crData)
		if errors.IsAlreadyExists(forError) {
			logger.Info("FPGACR is exist.")
			forError = nil
		} else if nil != forError {
			logger.Error("Failed to create FPGACR.", zap.Error(forError))
		} else {
			logger.Info("Success to create FPGACR.")
		}

		if nil != forError {
			err = forError
		}
	}

	return err
}

// Node and device information creation function
func MakeNodeAndDeviceInfo(
	ctx context.Context,
	pDevices *[]DeviceBasicInfo,
	pNodeAndDevices *[]DeviceInfo) error {

	var err error
	err = nil

	logger := ctxzap.Extract(ctx)

	// Repeat for the number of devices in the environment information
	for deviceIndex := 0; deviceIndex < len(*pDevices); deviceIndex++ {

		deviceData := (*pDevices)[deviceIndex]
		var nodeAndDeviceData DeviceInfo

		// Add information to the node and device information structure
		nodeAndDeviceData.NodeName = deviceData.NodeName
		nodeAndDeviceData.DeviceFilePath = deviceData.DeviceFilePath
		nodeAndDeviceData.DeviceUUID = deviceData.DeviceUUID
		nodeAndDeviceData.DeviceType = deviceData.DeviceType
		nodeAndDeviceData.DeviceIndex = deviceData.DeviceIndex
		*pNodeAndDevices = append(*pNodeAndDevices, nodeAndDeviceData)
	}
	if nil == err {
		logger.Info("nodeanddeviceinfostructure created successfully")
	}
	return err
}

// Function to create deployment information within the device
func MakeInDeviceDeployInfo(
	ctx context.Context,
	pDevices *[]DeviceBasicInfo,
	pRegionSpecifics *[]RegionSpecificInfo,
	pInDeviceDeploys *[]DeviceRegionInfo) error {

	var err error
	err = nil

	logger := ctxzap.Extract(ctx)

	// Repeat for the number of devices in the device information
	for deviceIndex := 0; deviceIndex < len(*pDevices); deviceIndex++ {

		deviceData := (*pDevices)[deviceIndex]
		var InDeviceDeploy DeviceRegionInfo

		var bsConfigCChar *C.char
		var bsConfig BsConfigInfo
		var childBsRegions []ChildBsRegion

		defer C.free(unsafe.Pointer(bsConfigCChar))
		parentIDCString := C.CString(deviceData.ParentID)
		defer C.free(unsafe.Pointer(parentIDCString))
		childIDCString := C.CString(deviceData.SubDeviceSpecRef)
		defer C.free(unsafe.Pointer(childIDCString))

		if "MEM" == deviceData.DeviceCategory {
			continue
		}

		var regionNames []string
		var nopFlag bool = false

		// Repeat for the number of lists of domain-specific information
		for regionSpecificIndex := 0; regionSpecificIndex < len(*pRegionSpecifics); regionSpecificIndex++ {
			regionSpecificInfo := (*pRegionSpecifics)[regionSpecificIndex]
			functionTargets := regionSpecificInfo.FunctionTargets

			if deviceData.SubDeviceSpecRef != regionSpecificInfo.SubDeviceSpecRef {
				continue
			}
			// Repeat for the number of FunctionTargets in the domain-specific information
			for functionTargetIndex := 0; functionTargetIndex < len(functionTargets); functionTargetIndex++ {
				functionTargetData := functionTargets[functionTargetIndex]
				regionNames = append(regionNames, functionTargetData.RegionName)

				if nil == functionTargetData.MaxFunctions {
					nopFlag = true
				}
			}
		}

		if "FPGA" == deviceData.DeviceCategory {
			if false == nopFlag {

				// Call the FPGA function information acquisition function
				ret := C.fpga_db_get_device_config_by_bitstream_id(parentIDCString,
					childIDCString,
					(**C.char)(unsafe.Pointer(&bsConfigCChar))) //nolint:gocritic // suspicious indentical LHS and RHS for next block "==". QUESTION: why?
				if 0 > ret {
					logger.Error("fpga_db_get_device_config_by_bitstream_id() err = " +
						strconv.Itoa(int(ret)))
				} else {
					logger.Info("fpga_db_get_device_config_by_bitstream_id() ret = " +
						strconv.Itoa(int(ret)))
				}

				if 0 == ret {
					// Get the number of characters in the JSON data
					n := 0
					head := (*byte)(unsafe.Pointer(bsConfigCChar))
					for ptr := head; *ptr != 0; n++ {
						ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + 1))
					}

					// Convert to []bytes type
					bsConfigBytes := C.GoBytes(unsafe.Pointer(bsConfigCChar), C.int(n))

					// Get FPGA function information
					functionError := json.Unmarshal(bsConfigBytes, &bsConfig)
					if nil != functionError {
						logger.Error("unable to unmarshal fpga_db_get_device_config_by_bitstream_id() data.",
							zap.Error(functionError))
					}

					childBsRegions = bsConfig.ChildBitstreamIDs[0].Regions
				}
			}
		}

		// Repeat for the number of FunctionTargets(RegionName) in the device information
		for _, regionName := range regionNames {
			var functionTargetData RegionInDevice
			var functionDataSlice []SimpleFunctionInfraStruct

			var childBsModule ChildBsModule
			var functions []ChildBsFunctions
			// Repeat for the number of regions in the FPGA function information
			for regionSpecificIndex := 0; regionSpecificIndex < len(childBsRegions); regionSpecificIndex++ {

				if regionName != *childBsRegions[regionSpecificIndex].Name {
					continue
				}
				childBsModule = *childBsRegions[regionSpecificIndex].Modules
				functions = *childBsModule.Functions
			}

			// Repeat for the number of function lists
			for functionKernelIndex := 0; functionKernelIndex < len(functions); functionKernelIndex++ {
				function := functions[functionKernelIndex]
				var functionData SimpleFunctionInfraStruct

				// Add information to Functions in Function to create deployment information within the device
				functionData.FunctionName = (*(*function.Module)[0].Type)
				functionData.FunctionIndex = int32(functionKernelIndex)
				functionData.FrameworkKernelID = int32(*childBsModule.Chain.ID)
				functionData.PartitionName = strconv.Itoa(int(*function.ID))
				functionDataSlice =
					append(functionDataSlice, functionData)
			}
			// Add information to FunctionTargets in Function to create deployment information within the device
			functionTargetData.RegionName = regionName
			functionTargetData.Functions = &functionDataSlice
			InDeviceDeploy.FunctionTargets = append(InDeviceDeploy.FunctionTargets, functionTargetData)

		}

		// Function to create deployment information within the device
		InDeviceDeploy.NodeName = deviceData.NodeName
		InDeviceDeploy.DeviceFilePath = deviceData.DeviceFilePath
		InDeviceDeploy.DeviceUUID = deviceData.DeviceUUID
		InDeviceDeploy.SubDeviceSpecRef = deviceData.SubDeviceSpecRef
		*pInDeviceDeploys = append(*pInDeviceDeploys, InDeviceDeploy)
	}

	if nil == err {
		logger.Info("indevicedeployinfostructure created successfully")
	}
	return err
}

// Infrastructure configuration information creation function
func MakeInfrastructureInfo(
	ctx context.Context,
	mng ctrl.Manager,
	pNodeAndDevices *[]DeviceInfo,
	pInfrastructureInfo *map[string][]cm.DeviceInfo) error {

	var err error
	err = nil

	// logger := ctxzap.Extract(ctx)

	// Repeat for the number of devices in the node and device information
	for deviceIndex := 0; deviceIndex < len(*pNodeAndDevices); deviceIndex++ {

		nodeAndDeviceInfo := cm.DeviceInfo((*pNodeAndDevices)[deviceIndex])

		/* Provisional support (dynamic reconfiguration not supported) */
		// Call FuncName conversion function
		nodeAndDeviceInfo.DeviceType = convDeviceType(ctx, nodeAndDeviceInfo.DeviceType)

		// Add information to the infrastructure configuration information structure
		(*pInfrastructureInfo)["devices"] =
			append((*pInfrastructureInfo)["devices"], nodeAndDeviceInfo)
	}
	/* Supports multiple workers */
	for {
		var eventKind int
		err, eventKind = cm.GetConfigMap(ctx, mng, cm.CMInfraInfo)
		if nil != err {
			break
		}
		data := cm.ExistingInfraCMupdate(ctx, pInfrastructureInfo)
		err = cm.CreateConfigMap(ctx, mng, cm.CMInfraInfo, data, eventKind)
		if nil != err {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	return err
}

// Deployment area information creation function
func MakeDeployInfo(
	ctx context.Context,
	mng ctrl.Manager,
	pInDeviceDeploys *[]DeviceRegionInfo,
	pRegionSpecifics *[]RegionSpecificInfo,
	pFPGACatalogs *[]FPGACatalog,
	pDeployRegionInfo *map[string][]cm.DeviceRegionInfo) error {

	var err error
	err = nil

	// Repeat for the number of devices in the device deployment information
	for deviceIndex := 0; deviceIndex < len(*pInDeviceDeploys); deviceIndex++ {

		inDeviceDeployInfo := (*pInDeviceDeploys)[deviceIndex]
		inDeviceDeployFunctionTargets := inDeviceDeployInfo.FunctionTargets
		var deviceRegionInfo cm.DeviceRegionInfo

		var regionFunctionTargets []RegionInDevice
		// Repeat for the number of lists of domain-specific information
		for regionSpecificIndex := 0; regionSpecificIndex < len(*pRegionSpecifics); regionSpecificIndex++ {
			regionSpecificInfo := (*pRegionSpecifics)[regionSpecificIndex]
			regionSpecificFunctionTargets := regionSpecificInfo.FunctionTargets

			if inDeviceDeployInfo.SubDeviceSpecRef != regionSpecificInfo.SubDeviceSpecRef {
				continue
			}

			for functionTargetIndex := 0; functionTargetIndex < len(regionSpecificFunctionTargets); functionTargetIndex++ {
				// Get FunctionTargets for domain-specific information that matches ChildId
				regionSpecificFunctionTarget := regionSpecificFunctionTargets[functionTargetIndex]
				regionFunctionTargets = append(regionFunctionTargets, regionSpecificFunctionTarget)
			}
		}

		// Repeat for the number of FunctionTargets in the device deployment information
		for functionTargetIndex := 0; functionTargetIndex < len(inDeviceDeployFunctionTargets); functionTargetIndex++ {

			inDeviceDeployFunctionTarget := inDeviceDeployFunctionTargets[functionTargetIndex]
			inDeviceDeployFunctions := inDeviceDeployFunctionTarget.Functions
			var regionInDevice cm.RegionInDevice

			// Repeat for the number of Functions deployed in the device
			for functionIndex := 0; functionIndex < len(*inDeviceDeployFunctions); functionIndex++ {
				inDeviceDeployFunction := (*inDeviceDeployFunctions)[functionIndex]
				var simpleFunctionInfraStruct cm.SimpleFunctionInfraStruct

				// Repeat for the number of lists of function-specific information
				for functionSpecificIndex := 0; functionSpecificIndex < len(*pFPGACatalogs); functionSpecificIndex++ {
					functionSpecificInfo := (*pFPGACatalogs)[functionSpecificIndex]
					if inDeviceDeployFunction.FunctionName != functionSpecificInfo.FunctionName {
						continue
					}
					// Add information to the Functions of the deployment area information structure.
					simpleFunctionInfraStruct.FunctionIndex = &inDeviceDeployFunction.FunctionIndex
					simpleFunctionInfraStruct.PartitionName = inDeviceDeployFunction.PartitionName
					simpleFunctionInfraStruct.FunctionName = functionSpecificInfo.FunctionName
					simpleFunctionInfraStruct.MaxDataFlows = functionSpecificInfo.MaxDataFlows
					simpleFunctionInfraStruct.MaxCapacity = functionSpecificInfo.MaxCapacity
					regionInDevice.Functions = append(regionInDevice.Functions, simpleFunctionInfraStruct)
				}
			}

			// Repeat for the number of FunctionTargets of domain-specific information that matches ChildId
			for functionTargetIndex := 0; functionTargetIndex < len(regionFunctionTargets); functionTargetIndex++ {
				regionFunctionTarget := regionFunctionTargets[functionTargetIndex]
				if inDeviceDeployFunctionTarget.RegionName != regionFunctionTarget.RegionName {
					continue
				}

				// Add information to the FunctionTargets of the deployment area information structure.
				regionInDevice.RegionName = regionFunctionTarget.RegionName
				regionInDevice.RegionType = regionFunctionTarget.RegionType
				regionInDevice.MaxFunctions = regionFunctionTarget.MaxFunctions
				regionInDevice.MaxCapacity = regionFunctionTarget.MaxCapacity
				deviceRegionInfo.FunctionTargets =
					append(deviceRegionInfo.FunctionTargets, regionInDevice)
			}
		}

		// Add information to the deployment area information structure
		deviceRegionInfo.NodeName = inDeviceDeployInfo.NodeName
		deviceRegionInfo.DeviceFilePath = inDeviceDeployInfo.DeviceFilePath
		*deviceRegionInfo.DeviceUUID = *inDeviceDeployInfo.DeviceUUID
		(*pDeployRegionInfo)["devices"] = append((*pDeployRegionInfo)["devices"], deviceRegionInfo)
	}

	/* Supports multiple workers */
	for {
		var eventKind int
		err, eventKind = cm.GetConfigMap(ctx, mng, cm.CMDeployInfo)
		if nil != err {
			break
		}
		data := cm.ExistingDeployCMupdate(ctx, pDeployRegionInfo)
		err = cm.CreateConfigMap(ctx, mng, cm.CMDeployInfo, data, eventKind)
		if nil != err {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	return err
}

/* Provisional support (dynamic reconfiguration not supported) */
// Deployment area information creation function
func MakeDeployInfoConvFuncName(
	ctx context.Context,
	mng ctrl.Manager,
	pInDeviceDeploys *[]DeviceRegionInfo,
	pRegionSpecifics *[]RegionSpecificInfo,
	pFPGACatalogs *[]FPGACatalog,
	pDeployRegionInfo *map[string][]cm.DeviceRegionInfo) error {

	var err error
	err = nil
	var framesizeConfigCChar *C.char
	var framesizeConfig FrameSizeConfig

	logger := ctxzap.Extract(ctx)

	// Repeat for the number of devices in the device deployment information
	for deviceIndex := 0; deviceIndex < len(*pInDeviceDeploys); deviceIndex++ {

		inDeviceDeployInfo := (*pInDeviceDeploys)[deviceIndex]
		inDeviceDeployFunctionTargets := inDeviceDeployInfo.FunctionTargets
		var deviceRegionInfo cm.DeviceRegionInfo
		var deviceID C.uint

		var regionFunctionTargets []RegionInDevice
		// Repeat for the number of lists of domain-specific information
		for regionSpecificIndex := 0; regionSpecificIndex < len(*pRegionSpecifics); regionSpecificIndex++ {
			regionSpecificInfo := (*pRegionSpecifics)[regionSpecificIndex]
			regionSpecificFunctionTargets := regionSpecificInfo.FunctionTargets

			if inDeviceDeployInfo.SubDeviceSpecRef != regionSpecificInfo.SubDeviceSpecRef {
				continue
			}

			for functionTargetIndex := 0; functionTargetIndex < len(regionSpecificFunctionTargets); functionTargetIndex++ {
				// Get FunctionTargets for domain-specific information that matches ChildId
				regionSpecificFunctionTarget := regionSpecificFunctionTargets[functionTargetIndex]
				regionFunctionTargets = append(regionFunctionTargets, regionSpecificFunctionTarget)
			}
		}

		// Repeat for the number of FunctionTargets in the device deployment information
		for functionTargetIndex := 0; functionTargetIndex < len(inDeviceDeployFunctionTargets); functionTargetIndex++ {

			inDeviceDeployFunctionTarget := inDeviceDeployFunctionTargets[functionTargetIndex]
			inDeviceDeployFunctions := inDeviceDeployFunctionTarget.Functions
			var regionInDevice cm.RegionInDevice

			// Repeat for the number of Functions deployed in the device
			for functionIndex := 0; functionIndex < len(*inDeviceDeployFunctions); functionIndex++ {
				inDeviceDeployFunction := (*inDeviceDeployFunctions)[functionIndex]
				var deployRegionFunction cm.SimpleFunctionInfraStruct

				// Repeat for the number of lists of function-specific information
				for functionSpecificIndex := 0; functionSpecificIndex < len(*pFPGACatalogs); functionSpecificIndex++ {
					functionSpecificInfo := (*pFPGACatalogs)[functionSpecificIndex]
					if inDeviceDeployFunction.FunctionName != functionSpecificInfo.FunctionName {
						continue
					}
					// Add information to the Functions of the deployment area information structure.
					deployRegionFunction.FunctionIndex = &inDeviceDeployFunction.FunctionIndex
					deployRegionFunction.PartitionName = inDeviceDeployFunction.PartitionName
					deployRegionFunction.FunctionName = functionSpecificInfo.FunctionName
					deployRegionFunction.MaxDataFlows = functionSpecificInfo.MaxDataFlows
					deployRegionFunction.MaxCapacity = functionSpecificInfo.MaxCapacity

					fpgaDeviceUUIDCString := C.CString(*inDeviceDeployInfo.DeviceUUID)
					defer C.free(unsafe.Pointer(fpgaDeviceUUIDCString))
					// Call the device ID acquisition function
					ret := C.fpga_get_dev_id(fpgaDeviceUUIDCString, &deviceID)
					if 0 > ret {
						logger.Error("fpga_get_dev_id() err = " + strconv.Itoa(int(ret)))
						break
					} else {
						logger.Info("fpga_get_dev_id() ret = " + strconv.Itoa(int(ret)))
					}

					for doWhile := 0; doWhile < 1; doWhile++ {

						kernelID, _ := strconv.Atoi(deployRegionFunction.PartitionName)

						// Call the function registration function
						functionname := C.CString(functionSpecificInfo.FunctionName)
						defer C.free(unsafe.Pointer(functionname))
						ret := C.fpga_function_config(deviceID, C.uint(kernelID),
							functionname)
						if 0 > ret {
							logger.Error("fpga_function_config() err = " +
								strconv.Itoa(int(ret)))
							break
						} else {
							logger.Info("fpga_function_config() ret = " +
								strconv.Itoa(int(ret)))
						}

						// Call the FPGA frame size setting acquisition function
						ret = C.fpga_function_get(deviceID, C.uint(kernelID),
							(**C.char)(unsafe.Pointer(&framesizeConfigCChar))) //nolint:gocritic // suspicious indentical LHS and RHS for next block "==". QUESTION: why?
						if 0 > ret {
							logger.Error("fpga_function_get() err = " +
								strconv.Itoa(int(ret)))
							break
						} else {
							logger.Info("fpga_function_get() ret = " +
								strconv.Itoa(int(ret)))
						}

						// Get the number of characters in the JSON data
						n := 0
						head := (*byte)(unsafe.Pointer(framesizeConfigCChar))
						for ptr := head; *ptr != 0; n++ {
							ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + 1))
						}

						// Convert to []bytes type
						framesizeConfigByte := C.GoBytes(unsafe.Pointer(
							framesizeConfigCChar), C.int(n))

						// Get frame size information
						err := json.Unmarshal(framesizeConfigByte,
							&framesizeConfig)
						if nil != err {
							logger.Error("unable to unmarshal fpga_function_get() data.", zap.Error(err))
						}

						outputHeight := framesizeConfig.OutputHeight
						outputWidth := framesizeConfig.OutputWidth

						// Provisional support (dynamic reconfiguration not supported)
						// Call FuncName conversion function
						deployRegionFunction.FunctionName = convFunctionName(ctx, outputHeight,
							outputWidth, deployRegionFunction.FunctionName)
					}

					regionInDevice.Functions = append(regionInDevice.Functions, deployRegionFunction)
				}
			}

			// Repeat for the number of FunctionTargets of domain-specific information that matches ChildId
			for functionTargetIndex := 0; functionTargetIndex < len(regionFunctionTargets); functionTargetIndex++ {
				regionFunctionTarget := regionFunctionTargets[functionTargetIndex]
				if inDeviceDeployFunctionTarget.RegionName != regionFunctionTarget.RegionName {
					continue
				}

				// Add information to the FunctionTargets of the deployment area information structure.
				regionInDevice.RegionName = regionFunctionTarget.RegionName
				regionInDevice.RegionType = regionFunctionTarget.RegionType
				if nil != regionFunctionTarget.MaxFunctions {
					if "cpu" == inDeviceDeployFunctionTarget.RegionName {
						maxFunctions := *regionFunctionTarget.MaxFunctions * cpuNum
						regionInDevice.MaxFunctions = &maxFunctions
					} else {
						regionInDevice.MaxFunctions = regionFunctionTarget.MaxFunctions
					}
				}
				if nil != regionFunctionTarget.MaxCapacity {
					if "cpu" == inDeviceDeployFunctionTarget.RegionName {
						maxCapacity := *regionFunctionTarget.MaxCapacity * cpuNum
						regionInDevice.MaxCapacity = &maxCapacity
					} else {
						regionInDevice.MaxCapacity = regionFunctionTarget.MaxCapacity
					}
				}

				deviceRegionInfo.FunctionTargets =
					append(deviceRegionInfo.FunctionTargets, regionInDevice)
			}
		}

		deviceUUID := *inDeviceDeployInfo.DeviceUUID

		// Add information to the deployment area information structure
		deviceRegionInfo.NodeName = inDeviceDeployInfo.NodeName
		deviceRegionInfo.DeviceFilePath = inDeviceDeployInfo.DeviceFilePath
		deviceRegionInfo.DeviceUUID = &deviceUUID
		(*pDeployRegionInfo)["devices"] = append((*pDeployRegionInfo)["devices"], deviceRegionInfo)
	}

	/* Supports multiple workers */
	for {
		var eventKind int
		err, eventKind = cm.GetConfigMap(ctx, mng, cm.CMDeployInfo)
		if nil != err {
			break
		}
		data := cm.ExistingDeployCMupdate(ctx, pDeployRegionInfo)
		err = cm.CreateConfigMap(ctx, mng, cm.CMDeployInfo, data, eventKind)
		if nil != err {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	return err
}

// Circuit placement information creation function
func MakeFPGACatalogInfo(
	ctx context.Context,
	mng ctrl.Manager,
	pInDeviceDeploys *[]DeviceRegionInfo,
	pFunctionDedicatedDecodeInfo *map[string][]FunctionDedicatedInfo,
	pFunctionDedicatedResizeInfo *map[string][]FunctionDedicatedInfo,
	pServicerMgmtInfo *[]ServicerMgmtInfo,
	pFPGACatalogInfo *map[string][]cm.FPGACatalog) error {

	var err error
	err = nil

	// logger := ctxzap.Extract(ctx)

	// Repeat for the number of devices in the device deployment information
	for deviceIndex := 0; deviceIndex < len(*pInDeviceDeploys); deviceIndex++ {

		inDeviceDeployInfo := (*pInDeviceDeploys)[deviceIndex]
		inDeviceDeployFunctionTargets := inDeviceDeployInfo.FunctionTargets
		var fpgaCatalogInfo cm.FPGACatalog
		var appendFlag bool = true

		// Repeat for the number of FunctionTargets in the device deployment information
		for functionTargetIndex := 0; functionTargetIndex < len(inDeviceDeployFunctionTargets); functionTargetIndex++ {

			inDeviceDeployFunctionTarget := inDeviceDeployFunctionTargets[functionTargetIndex]
			inDeviceDeployFunctions := inDeviceDeployFunctionTarget.Functions
			var fpgaCatalogDeviceDetail cm.DeviceDetails
			var ip string

			if 0 == len(*inDeviceDeployFunctions) {
				appendFlag = false
				break
			}

			// Repeat for the number of Functions deployed in the device
			for functionIndex := 0; functionIndex < len(*inDeviceDeployFunctions); functionIndex++ {
				inDeviceDeployFunction := (*inDeviceDeployFunctions)[functionIndex]
				var fpgaCatalogFunctionData cm.FuncDataStruct
				var functionDedicatedInfos []FunctionDedicatedInfo

				// Decide whether to decode or filter/resize
				if "decode" == inDeviceDeployFunction.FunctionName {
					functionDedicatedInfos =
						(*pFunctionDedicatedDecodeInfo)["functionKernels"]
				} else {
					functionDedicatedInfos =
						(*pFunctionDedicatedResizeInfo)["functionKernels"]
				}

				var functionChannelIDs []int32
				// Function-specific information - Repeat for the number of lists of dedicated information
				for functionDedicatedIndex := 0; functionDedicatedIndex < len(functionDedicatedInfos); functionDedicatedIndex++ {
					functionDedicatedDetail := functionDedicatedInfos[functionDedicatedIndex]

					if inDeviceDeployFunction.PartitionName !=
						functionDedicatedDetail.PartitionName {
						continue
					}
					functionChannelIDs = functionDedicatedDetail.FunctionChannelIDList
					break
				}

				PartitionNameInt, _ := strconv.Atoi(inDeviceDeployFunction.PartitionName)
				PartitionName := int32(PartitionNameInt)

				// Add FuncData information to the circuit placement information structure
				fpgaCatalogFunctionData.FunctionIndex = inDeviceDeployFunction.FunctionIndex
				fpgaCatalogFunctionData.FunctionKernelID = PartitionName
				fpgaCatalogFunctionData.FrameworkKernelID = inDeviceDeployFunction.FrameworkKernelID
				fpgaCatalogFunctionData.FunctionChannelIDs = functionChannelIDs
				fpgaCatalogDeviceDetail.FuncData = append(fpgaCatalogDeviceDetail.FuncData, fpgaCatalogFunctionData)
			}

			// Repeat for the number of NodeNames in the servicer managemant information
			for _, servicerMgmtInfo := range *pServicerMgmtInfo {
				if cm.GMyNodeName != servicerMgmtInfo.NodeName {
					continue
				}
				networks := servicerMgmtInfo.NetworkInfo

				// Repeat for the number of networks in the servicer managemant information
				for networkindex := 0; networkindex < len(networks); networkindex++ {
					networkData := networks[networkindex]

					if *inDeviceDeployInfo.DeviceFilePath != networkData.DeviceFilePath {
						continue
					}

					servicerRegion :=
						"lane" + strconv.Itoa(int(networkData.LaneIndex))
					if inDeviceDeployFunctionTarget.RegionName != servicerRegion {
						continue
					}

					ip = networkData.IPAddress
					break
				}
			}

			// Add information to the Details of the circuit placement information structure
			fpgaCatalogDeviceDetail.RegionName = inDeviceDeployFunctionTarget.RegionName
			fpgaCatalogDeviceDetail.IPAddress = ip
			fpgaCatalogInfo.Details = append(fpgaCatalogInfo.Details, fpgaCatalogDeviceDetail)

		}

		if true != appendFlag {
			continue
		}

		// Add information to the circuit's deployment information structure
		fpgaCatalogInfo.NodeName = inDeviceDeployInfo.NodeName
		fpgaCatalogInfo.DeviceFilePath = *inDeviceDeployInfo.DeviceFilePath
		fpgaCatalogInfo.DeviceUUID = *inDeviceDeployInfo.DeviceUUID
		(*pFPGACatalogInfo)["devices"] =
			append((*pFPGACatalogInfo)["devices"], fpgaCatalogInfo)
	}

	/* Supports multiple workers */
	for {
		var eventKind int
		err, eventKind = cm.GetConfigMap(ctx, mng, cm.CMFPGACatalog)
		if nil != err {
			break
		}
		data := cm.ExistingFPGACatalogCMupdate(ctx, pFPGACatalogInfo)
		err = cm.CreateConfigMap(ctx, mng, cm.CMFPGACatalog, data, eventKind)
		if nil != err {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	return err
}

// Decode resource information creation function
func MakeDecodeInfo(
	ctx context.Context,
	mng ctrl.Manager,
	pFunctionDedicatedDecodeInfo *map[string][]FunctionDedicatedInfo,
	pFunctionDecodeInfo *map[string][]cm.FunctionDedicatedInfo) error {

	functionDedicatedInfos := (*pFunctionDedicatedDecodeInfo)["functionKernels"]
	var err error
	err = nil

	// logger := ctxzap.Extract(ctx)

	// Repeat for the number of FunctionKernels of Function-specific information-decode-only information
	for functionDedicatedIndex := 0; functionDedicatedIndex < len(functionDedicatedInfos); functionDedicatedIndex++ {

		var functionDecodeData cm.FunctionDedicatedInfo
		functionDedicatedDetail := functionDedicatedInfos[functionDedicatedIndex]
		functionDetails := functionDedicatedDetail.FunctionChannelIDs

		// Repeat for the number of FuncCHIds of Function-specific information-decode-only information
		for functionDetailIndex := 0; functionDetailIndex < len(functionDetails); functionDetailIndex++ {

			var functionDetailData cm.FunctionDetail
			functionDetail := functionDetails[functionDetailIndex]

			rxMakeFlag := true
			for protocol := range functionDetail.Rx.Protocol {
				if _, ok := functionDetailData.Rx.Protocol[protocol]; ok {
					rxMakeFlag = false
				}
			}
			if true == rxMakeFlag {
				functionDetailData.Rx.Protocol = make(map[string]cm.FPGAConnectionCatalogDetails)
			}

			txMakeFlag := true
			for protocol := range functionDetail.Tx.Protocol {
				if _, ok := functionDetailData.Tx.Protocol[protocol]; ok {
					txMakeFlag = false
				}
			}
			if true == txMakeFlag {
				functionDetailData.Tx.Protocol = make(map[string]cm.FPGAConnectionCatalogDetails)
			}

			functionDetailData.FunctionChannelID = functionDetail.FunctionChannelID
			for protocol := range functionDetail.Rx.Protocol {
				functionDetailData.Rx.Protocol[protocol] =
					cm.FPGAConnectionCatalogDetails(functionDetail.Rx.Protocol[protocol])
			}
			for protocol := range functionDetail.Tx.Protocol {
				functionDetailData.Tx.Protocol[protocol] =
					cm.FPGAConnectionCatalogDetails(functionDetail.Tx.Protocol[protocol])
			}

			// Add information to the decode resource information structure
			functionDecodeData.FunctionChannelIDs =
				append(functionDecodeData.FunctionChannelIDs, functionDetailData)

		}

		functionDecodeData.PartitionName = functionDedicatedDetail.PartitionName
		(*pFunctionDecodeInfo)["functionKernels"] =
			append((*pFunctionDecodeInfo)["functionKernels"], functionDecodeData)
	}

	/* Supports multiple workers */
	for {
		var eventKind int
		err, eventKind = cm.GetConfigMap(ctx, mng, cm.CMFPGADecode)
		if nil != err {
			break
		}
		err = cm.CreateConfigMap(ctx, mng, cm.CMFPGADecode, *pFunctionDecodeInfo, eventKind)
		if nil != err {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	return err
}

// filter/resize resource information creation function
func MakeFilterResizeInfo(
	ctx context.Context,
	mng ctrl.Manager,
	pFunctionDedicatedResizeInfo *map[string][]FunctionDedicatedInfo,
	pFunctionResizeInfo *map[string][]cm.FunctionDedicatedInfo) error {

	functionDedicatedInfos := (*pFunctionDedicatedResizeInfo)["functionKernels"]
	var err error
	err = nil

	// logger := ctxzap.Extract(ctx)

	// Function-specific information - Repeat for the number of FunctionKernels for filter/resize-specific information
	for functionDedicatedIndex := 0; functionDedicatedIndex < len(functionDedicatedInfos); functionDedicatedIndex++ {

		var functionResizeData cm.FunctionDedicatedInfo
		functionDedicatedDetail := functionDedicatedInfos[functionDedicatedIndex]
		functionDetails := functionDedicatedDetail.FunctionChannelIDs

		// Function-specific information - Repeat for the number of FuncCHIds in filter/resize-specific information
		for functionDetailIndex := 0; functionDetailIndex < len(functionDetails); functionDetailIndex++ {

			var functionDetailData cm.FunctionDetail
			functionDetail := functionDetails[functionDetailIndex]

			functionDetailData.FunctionChannelID = functionDetail.FunctionChannelID

			rxMakeFlag := true
			for protocol := range functionDetail.Rx.Protocol {
				if _, ok := functionDetailData.Rx.Protocol[protocol]; ok {
					rxMakeFlag = false
				}
			}
			if true == rxMakeFlag {
				functionDetailData.Rx.Protocol = make(map[string]cm.FPGAConnectionCatalogDetails)
			}

			txMakeFlag := true
			for protocol := range functionDetail.Tx.Protocol {
				if _, ok := functionDetailData.Tx.Protocol[protocol]; ok {
					txMakeFlag = false
				}
			}
			if true == txMakeFlag {
				functionDetailData.Tx.Protocol = make(map[string]cm.FPGAConnectionCatalogDetails)
			}

			for protocol := range functionDetail.Rx.Protocol {
				functionDetailData.Rx.Protocol[protocol] =
					cm.FPGAConnectionCatalogDetails(functionDetail.Rx.Protocol[protocol])
			}
			for protocol := range functionDetail.Tx.Protocol {
				functionDetailData.Tx.Protocol[protocol] =
					cm.FPGAConnectionCatalogDetails(functionDetail.Tx.Protocol[protocol])
			}

			// Add information to the filter/resize resource info structure
			functionResizeData.FunctionChannelIDs =
				append(functionResizeData.FunctionChannelIDs, functionDetailData)
		}

		functionResizeData.PartitionName = functionDedicatedDetail.PartitionName
		(*pFunctionResizeInfo)["functionKernels"] =
			append((*pFunctionResizeInfo)["functionKernels"], functionResizeData)
	}

	/* Supports multiple workers */
	for {
		var eventKind int
		err, eventKind = cm.GetConfigMap(ctx, mng, cm.CMFPGAFilterResize)
		if nil != err {
			break
		}
		err = cm.CreateConfigMap(ctx, mng, cm.CMFPGAFilterResize, *pFunctionResizeInfo, eventKind)
		if nil != err {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	return err
}

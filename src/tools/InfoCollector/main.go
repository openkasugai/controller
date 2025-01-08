/* Copyright 2024 NTT Corporation , FUJITSU LIMITED */

package main

import (
	"context"
	cm "example.com/InfoCollector/pkg/configmap"
	infocol "example.com/InfoCollector/pkg/infocollect"
	logging "example.com/InfoCollector/pkg/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {

	var err error
	var ret int
	var regionSpecifics []infocol.RegionSpecificInfo
	var functionSpecifics []infocol.FPGACatalog
	var functionDedicatedDecodeInfo map[string][]infocol.FunctionDedicatedInfo
	var functionDedicatedResizeInfo map[string][]infocol.FunctionDedicatedInfo
	var servicerMgmtInfo []infocol.ServicerMgmtInfo
	var devices []infocol.DeviceBasicInfo
	var nodeAndDevices []infocol.DeviceInfo
	var inDeviceDeploys []infocol.DeviceRegionInfo
	var infrastructureInfo map[string][]cm.DeviceInfo
	var deployRegionInfo map[string][]cm.DeviceRegionInfo
	// var fpgaCatalogInfo map[string][]cm.FPGACatalog
	// var functionDecodeInfo map[string][]cm.FunctionDedicatedInfo
	var functionResizeInfo map[string][]cm.FunctionDedicatedInfo

	infrastructureInfo = make(map[string][]cm.DeviceInfo)
	deployRegionInfo = make(map[string][]cm.DeviceRegionInfo)
	// fpgaCatalogInfo = make(map[string][]cm.FPGACatalog)
	// functionDecodeInfo = make(map[string][]cm.FunctionDedicatedInfo)
	functionResizeInfo = make(map[string][]cm.FunctionDedicatedInfo)

	ctx := context.Background()
	logger := logging.SettingZapLogger()
	ctx = ctxzap.ToContext(ctx, logger)
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{})
	if err != nil {
		os.Exit(1)
	}

	for doWhile := 0; doWhile < 1; doWhile++ {
		// Pre-deployment information acquisition function
		err = infocol.GetJsonFile(
			ctx,
			&regionSpecifics,
			&functionSpecifics,
			&functionDedicatedDecodeInfo,
			&functionDedicatedResizeInfo,
			&servicerMgmtInfo)
		if nil != err {
			logger.Error("Get JsonFile error", zap.Error(err))
			break
		}

		// Device information acquisition function
		ret = infocol.GetDeviceInfo(ctx, &devices)
		if 0 > ret {
			logger.Error("Create DeviceInfo error", zap.Error(err))
			break
		}
		logger.Debug("DeviceInfo Data", zap.Any("deviceInfo", devices))

		// fpga information creation function
		err = infocol.CreateFPGACR(
			ctx,
			mgr,
			&regionSpecifics,
			&devices)
		if nil != err {
			logger.Error("Create fpgacr error", zap.Error(err))
			// break
		}

		// Node and device information creation function
		err = infocol.MakeNodeAndDeviceInfo(
			ctx,
			&devices,
			&nodeAndDevices)
		if nil != err {
			logger.Error("Create NodeAndDeviceInfo error", zap.Error(err))
			break
		}
		logger.Debug("NodeAndDeviceInfo Data", zap.Any("nodeAndDeviceInfo", nodeAndDevices))

		// Function to create deployment information within the device
		err = infocol.MakeInDeviceDeployInfo(
			ctx,
			&devices,
			&regionSpecifics,
			&inDeviceDeploys)
		if nil != err {
			logger.Error("Create InDeviceDeployInfo error", zap.Error(err))
			break
		}
		logger.Debug("InDeviceDeployInfo Data", zap.Any("inDeviceDeployInfo", inDeviceDeploys))

		// Infrastructure configuration information creation function
		err = infocol.MakeInfrastructureInfo(
			ctx,
			mgr,
			&nodeAndDevices,
			&infrastructureInfo)
		if nil != err {
			logger.Error("Create or Update infrastructureinfo error", zap.Error(err))
			break
		}

		// Deployment area information creation function
		/*
			err = infocol.MakeDeployInfo(
				&inDeviceDeploys,
				&regionSpecifics,
				&functionSpecifics,
				&deployRegionInfo)
			if nil != err {
				logger.Error("Create or Update deployinfo error")
				break
			}
		*/
		/* Provisional support (dynamic reconfiguration not supported) */
		err = infocol.MakeDeployInfoConvFuncName(
			ctx,
			mgr,
			&inDeviceDeploys,
			&regionSpecifics,
			&functionSpecifics,
			&deployRegionInfo)
		if nil != err {
			logger.Error("Create or Update deployinfo error", zap.Error(err))
			break
		}

		// Circuit placement information creation function
		/*
			err = infocol.MakeFPGACatalogInfo(
				ctx,
				mgr,
				&inDeviceDeploys,
				&functionDedicatedDecodeInfo,
				&functionDedicatedResizeInfo,
				&servicerMgmtInfo,
				&fpgaCatalogInfo)
			if nil != err {
				logger.Error("Create or Update fpgacatalogmap error", zap.Error(err))
				break
			}
		*/

		// Decode resource information creation function
		/*
			err = infocol.MakeDecodeInfo(
				ctx,
				mgr,
				&functionDedicatedDecodeInfo,
				&functionDecodeInfo)
			if nil != err {
				logger.Error("Create or Update decode-ch error", zap.Error(err))
				break
			}
		*/

		// filter/resize resource information creation function
		err = infocol.MakeFilterResizeInfo(
			ctx,
			mgr,
			&functionDedicatedResizeInfo,
			&functionResizeInfo)
		if nil != err {
			logger.Error("Create or Update filter-resize-ch error", zap.Error(err))
			break
		}
	}
	if nil == err {
		logger.Info("infocollector was successful")
	}
}

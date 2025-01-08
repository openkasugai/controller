/* Copyright 2024 NTT Corporation , FUJITSU LIMITED */

package configmap

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

// Adding or overwriting existing infrastructure information CM
func ExistingInfraCMupdate(ctx context.Context, newData *map[string][]DeviceInfo) any {

	logger := ctxzap.Extract(ctx)

	getCMData := GInfrastructureInfo["devices"]
	newCMData := (*newData)["devices"]
	upData := make(map[string][]DeviceInfo)

	for j := 0; j < len(getCMData); j++ {
		if getCMData[j].NodeName == GMyNodeName {
			continue
		} else {
			upData["devices"] = append(upData["devices"], getCMData[j])
		}
	}
	upData["devices"] = append(upData["devices"], newCMData...)
	logger.Info("ExistingInfraCMupdate() success.")
	return upData
}

// Add or overwrite existing device deployment area information CM
func ExistingDeployCMupdate(ctx context.Context, newData *map[string][]DeviceRegionInfo) any {

	logger := ctxzap.Extract(ctx)

	getCMData := GDeployInfo["devices"]
	newCMData := (*newData)["devices"]
	upData := make(map[string][]DeviceRegionInfo)

	for j := 0; j < len(getCMData); j++ {
		if getCMData[j].NodeName == GMyNodeName {
			continue
		} else {
			upData["devices"] = append(upData["devices"], getCMData[j])
		}
	}
	upData["devices"] = append(upData["devices"], newCMData...)
	logger.Info("ExistingDeployCMupdate() success.")
	return upData
}

// Adding or overwriting existing FPGA catalog information CM
func ExistingFPGACatalogCMupdate(ctx context.Context, newData *map[string][]FPGACatalog) any {

	logger := ctxzap.Extract(ctx)

	getCMData := GFPGACatalogMap["devices"]
	newCMData := (*newData)["devices"]
	upData := make(map[string][]FPGACatalog)

	for j := 0; j < len(getCMData); j++ {
		if getCMData[j].NodeName == GMyNodeName {
			continue
		} else {
			upData["devices"] = append(upData["devices"], getCMData[j])
		}
	}
	upData["devices"] = append(upData["devices"], newCMData...)
	logger.Info("ExistingFPGACatalogCMupdate() success.")
	return upData
}

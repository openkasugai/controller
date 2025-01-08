/* Copyright 2024 NTT Corporation , FUJITSU LIMITED */

package configmap

import (
	"context"
	"encoding/json"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Get ConfigMap
func GetConfigMap(ctx context.Context,
	mng ctrl.Manager, cmName string) (err error, eventKind int) {

	logger := ctxzap.Extract(ctx)

	var cmData []byte
	var mapData map[string]string

	eventKind = UpdateEvent

	m := mng.GetAPIReader()

	tmpData := &unstructured.Unstructured{}
	tmpData.SetGroupVersionKind(schema.GroupVersionKind{
		Kind:    "ConfigMap",
		Version: "v1",
	})
	err = m.Get(ctx,
		client.ObjectKey{
			Namespace: "default",
			Name:      cmName,
		}, tmpData)
	if errors.IsNotFound(err) {
		eventKind = CreateEvent
		logger.Info("ConfigMap does not exist. ConfigName=" + cmName)
		err = nil
	} else if nil != err {
		logger.Error("unable to get ConfigMap.", zap.Error(err))
	} else {
		mapData, _, _ = unstructured.NestedStringMap(tmpData.Object, "data")
		for _, jsonRecord := range mapData {
			cmData = []byte(jsonRecord)
		}
		convertToJson(ctx, &cmData, cmName)
	}
	return err, eventKind
}

// Convert the Data part of the obtained ConfigMap into a structure.
func convertToJson(ctx context.Context, cmData *[]byte, cmName string) {

	logger := ctxzap.Extract(ctx)

	var err error
	var detailMap map[string][]FunctionDetail
	var detailMapFR map[string][]FunctionDetail

	err = nil

	switch cmName {
	case CMInfraInfo:
		err = json.Unmarshal(*cmData, &GInfrastructureInfo)
	case CMDeployInfo:
		err = json.Unmarshal(*cmData, &GDeployInfo)
	case CMFPGACatalog:
		err = json.Unmarshal(*cmData, &GFPGACatalogMap)
	case CMFPGADecode:
		err = json.Unmarshal(*cmData, &detailMap)
		if nil == err {
			GDecodeCH = detailMap["FunctionChannelIDs"]
		}
	case CMFPGAFilterResize:
		err = json.Unmarshal(*cmData, &detailMapFR)
		if nil == err {
			GFilterResizeCH = detailMapFR["FunctionChannelIDs"]
		}
	default:
		logger.Warn("convertToJson() no action.")
	}
	if nil != err {
		logger.Error("unable to unmarshal. ConfigMap="+cmName, zap.Error(err))
	} else {
		logger.Info("Success to unmarshal.")
	}
}

// Create ConfigMap
func CreateConfigMap(ctx context.Context,
	mng ctrl.Manager,
	cmName string,
	data any,
	eventKind int) error {

	logger := ctxzap.Extract(ctx)

	var err error

	for doWhile := 0; doWhile < 1; doWhile++ {

		createCR := &corev1.ConfigMap{}
		createCR.SetName(cmName)
		createCR.SetNamespace(CMNameSpace)

		jsonData, err := json.Marshal(data)

		if nil != err {
			logger.Error("unable to marshal. ConfigMap="+cmName, zap.Error(err))
			break
		}
		crData := map[string]string{
			cmName + ".json": string(jsonData),
		}

		createCR.Data = crData

		err = createOrUpdateCM(ctx, createCR, mng, eventKind)
		// Processing result
		if nil != err {
			logger.Error("createOrUpdateCM failed.", zap.Error(err))
		} else {
			logger.Info("createOrUpdateCM success.")
		}
	}
	return err
}

// Create or update a ConfigMap
func createOrUpdateCM(ctx context.Context,
	cmData *corev1.ConfigMap,
	mng ctrl.Manager,
	eventKind int) error {

	logger := ctxzap.Extract(ctx)

	var err error
	err = nil
	c := mng.GetClient()

	// Event detection
	if CreateEvent == eventKind {
		err = c.Create(ctx, cmData)
	} else if UpdateEvent == eventKind {
		err = c.Update(ctx, cmData)
	} else {
		logger.Error("Event Kind ERROR.")
	}

	// Processing result
	if nil != err {
		logger.Error("Create or Update ConfigMap failed.", zap.Error(err))
	} else {
		logger.Info("Create or Update ConfigMap success.")
	}
	return err
}

/* Functions for unit testing only *************************************************/
/* Get data from json file */
func GetJsonDataInfra(ctx context.Context, infrastructureInfo *map[string][]DeviceInfo) {

	logger := ctxzap.Extract(ctx)

	raw, err := ioutil.ReadFile("./ConfigFile/infrastructureinfo.json")
	if nil != err {
		logger.Error("unable to readfile.", zap.Error(err))
	} else {
		err = json.Unmarshal(raw, infrastructureInfo)
		if nil != err {
			logger.Error("unable to unmarshal.", zap.Error(err))
		}
	}
}

func GetJsonDataDep(ctx context.Context, deploy *map[string][]DeviceRegionInfo) {

	logger := ctxzap.Extract(ctx)

	raw, err := ioutil.ReadFile("./ConfigFile/deployinfo.json")
	if nil != err {
		logger.Error("unable to readfile.", zap.Error(err))
	} else {
		err = json.Unmarshal(raw, deploy)
		if nil != err {
			logger.Error("unable to unmarshal.", zap.Error(err))
		}
	}
}

func GetJsonDataFpgactlg(ctx context.Context, fpgaCatalogMap *map[string][]FPGACatalog) {

	logger := ctxzap.Extract(ctx)

	raw, err := ioutil.ReadFile("./ConfigFile/fpgacatalogmap.json")
	if nil != err {
		logger.Error("unable to readfile.", zap.Error(err))
	} else {
		err = json.Unmarshal(raw, fpgaCatalogMap)
		if nil != err {
			logger.Error("unable to unmarshal.", zap.Error(err))
		}
	}
}

func GetJsonDataFuncCH(ctx context.Context, funcCH *map[string][]FunctionDetail, cmName string) {

	logger := ctxzap.Extract(ctx)

	var filePath string

	if cmName == CMFPGADecode {
		filePath = "./ConfigFile/decode-ch.json"
	} else {
		filePath = "./ConfigFile/filter-resize-ch.json"
	}
	raw, err := ioutil.ReadFile(filePath)
	if nil != err {
		logger.Error("unable to readfile.", zap.Error(err))
	} else {
		err = json.Unmarshal(raw, funcCH)
		if nil != err {
			logger.Error("unable to unmarshal.", zap.Error(err))
		}
	}
}

/* Functions for unit testing only *************************************************/

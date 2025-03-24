/*
Copyright 2025 NTT Corporation, FUJITSU LIMITED

fpgareconfiguration functions
*/

package request

import (
	"context"
	"encoding/json"
	common "example.com/FPGAReconfigurationTool/pkg/common"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

func makeCRName(ctx context.Context,
	nodeName string,
	deviceFilePath string,
	pcrNameSpace *string,
	pcrName *string) error {

	var err error
	err = nil

	logger := ctxzap.Extract(ctx)

	*pcrNameSpace = os.Getenv("K8S_FPGARECONFIGURATION_NAMESPACE")
	if "" == *pcrNameSpace {
		*pcrNameSpace = "default"
	}

	var deviceUUID string
	deviceUUID = ""
	if true == strings.Contains(deviceFilePath, "/dev/xpcie_") {
		deviceUUID = strings.ReplaceAll(deviceFilePath, "/dev/xpcie_", "")
	}
	if "" == deviceUUID {
		logger.Info("DeviceFilePath error. DeviceFilePath=" + deviceFilePath)
		err = fmt.Errorf("UUID not found.  DeviceFilePath")
	}
	*pcrName = "manualfpgareconfig-" + nodeName + "-" + strings.ToLower(deviceUUID)

	return err
}

// check nodename function
func CheckNodeNameCR(
	ctx context.Context,
	mng ctrl.Manager,
	nodeName string) error {

	var err error

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		logger := ctxzap.Extract(ctx)

		m := mng.GetAPIReader()
		tmpData := &unstructured.Unstructured{}
		tmpData.SetGroupVersionKind(schema.GroupVersionKind{
			Version: "v1",
			Kind:    "Node",
		})
		err = m.Get(ctx, client.ObjectKey{
			Namespace: "default",
			Name:      nodeName}, tmpData)
		if errors.IsNotFound(err) {
			// If CR does not exist
			logger.Info(" unknown nodename error." +
				" NameSpace=default" +
				" Name=" + nodeName)
			break
		} else if err != nil {
			logger.Error("unable to fetch nodename."+
				" NameSpace=default"+
				" Name="+nodeName,
				zap.Error(err))
			break
		}
	}

	return err
}

// get fpga reconfiguration function
func getFPGAReconfigurationCR(
	ctx context.Context,
	mng ctrl.Manager,
	crNameSpace string,
	crName string,
	pcrFPGAReconfigurationData *FPGAReconfiguration) error {

	var err error

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		logger := ctxzap.Extract(ctx)

		m := mng.GetAPIReader()
		tmpData := &unstructured.Unstructured{}
		tmpData.SetGroupVersionKind(schema.GroupVersionKind{
			Version: "example.com/v1",
			Kind:    "FPGAReconfiguration",
		})
		err = m.Get(ctx, client.ObjectKey{
			Namespace: crNameSpace,
			Name:      crName}, tmpData)
		if errors.IsNotFound(err) {
			// If CR does not exist
			break
		} else if err != nil {
			logger.Error("unable to fetch FPGAReconfigurationCR."+
				" crNameSpace="+crNameSpace+
				" crName="+crName,
				zap.Error(err))
			break
		} else {
			var jsonstr []byte
			getstr, _, _ := unstructured.NestedString(tmpData.Object,
				"apiVersion")
			pcrFPGAReconfigurationData.TypeMeta.APIVersion = getstr
			getstr, _, _ = unstructured.NestedString(tmpData.Object, "kind")
			pcrFPGAReconfigurationData.TypeMeta.Kind = getstr
			getdata, _, _ := unstructured.NestedMap(tmpData.Object, "spec")
			jsonstr, err = json.Marshal(getdata)
			if nil != err {
				logger.Error("FPGAReconfiguration.Spec unable to Marshal.", zap.Error(err))
				break
			}
			err = json.Unmarshal(jsonstr, &pcrFPGAReconfigurationData.Spec)
			if nil != err {
				logger.Error("FPGAReconfiguration.Spec unable to Unmarshal.", zap.Error(err))
				break
			}
			getdata, _, _ = unstructured.NestedMap(tmpData.Object, "status")
			jsonstr, err = json.Marshal(getdata)
			if nil != err {
				logger.Error("FPGAReconfiguration.Status unable to Marshal.", zap.Error(err))
				break
			}
			err = json.Unmarshal(jsonstr, &pcrFPGAReconfigurationData.Status)
			if nil != err {
				logger.Error("FPGAReconfiguration.Status unable to Unmarshal.", zap.Error(err))
				break
			}
		}
	}
	return err
}

// create fpga reconfiguration function
func CreateFPGAReconfiguration(
	ctx context.Context,
	mng ctrl.Manager,
	nodeName string,
	deviceFilePath string,
	requestMode int32,
	laneSetting map[int32]string) error {

	var err error
	err = nil

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		logger := ctxzap.Extract(ctx)
		c := mng.GetClient()

		var crFPGAReconfigurationData FPGAReconfiguration

		err = DeleteFPGAReconfiguration(ctx,
			mng,
			nodeName,
			deviceFilePath)
		if nil != err {
			break
		}

		crFPGAReconfigurationData.Spec.NodeName = nodeName
		crFPGAReconfigurationData.Spec.DeviceFilePath = deviceFilePath
		setFlag := true
		if common.MODECHILDBSRESET == requestMode {
			crFPGAReconfigurationData.Spec.ChildBsResetFlag = &setFlag
		} else if common.MODEFPGARESET == requestMode {
			crFPGAReconfigurationData.Spec.FPGAResetFlag = &setFlag
		}
		if common.MODEMANUALSET == requestMode ||
			common.MODEFPGARESET == requestMode {
			// Set ConfigNames Parameter
			configNameCount := 0
			var configNames FPGAConfigNames
			for index, value := range laneSetting {
				configNames.LaneIndex = index
				configNames.ConfigName = value
				crFPGAReconfigurationData.Spec.ConfigNames =
					append(crFPGAReconfigurationData.Spec.ConfigNames, configNames)
				configNameCount++
			}
			if 0 == configNameCount {
				logger.Info("configname not found.")
				break
			}
		}

		crData := &unstructured.Unstructured{}
		crData.SetGroupVersionKind(schema.GroupVersionKind{
			Version: "example.com/v1",
			Kind:    "FPGAReconfiguration",
		})

		var crNameSpace string
		var crName string
		err = makeCRName(ctx, nodeName, deviceFilePath, &crNameSpace, &crName)
		if nil != err {
			logger.Info("FPGAReconfigurationCR namespace/name make error." + " deviceFilePath=" + deviceFilePath)
			break
		}
		crData.SetName(crName)
		crData.SetNamespace(crNameSpace)

		crData.UnstructuredContent()["spec"] = crFPGAReconfigurationData.Spec
		crData.UnstructuredContent()["status"] = crFPGAReconfigurationData.Status

		err = c.Create(ctx, crData)
		if errors.IsAlreadyExists(err) {
			logger.Info("FPGAReconfigurationCR is exist.")
		} else if nil != err {
			fmt.Println("Failed to create FPGAReconfigurationCR. ", err.Error())
			logger.Error("Failed to create FPGAReconfigurationCR.",
				zap.Error(err))
		} else {
			logger.Info("Success to create FPGAReconfigurationCR.")
		}
	}

	return err
}

// delete fpga reconfiguration function
func GetStatusFPGAReconfiguration(ctx context.Context,
	mng ctrl.Manager,
	nodeName string,
	deviceFilePath string) (string, error) {

	var err error
	err = nil
	var status string
	status = ""

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		logger := ctxzap.Extract(ctx)

		var crFPGAReconfigurationData FPGAReconfiguration
		var crNameSpace string
		var crName string

		err = makeCRName(ctx, nodeName, deviceFilePath, &crNameSpace, &crName)
		if nil != err {
			logger.Info("FPGAReconfigurationCR namespace/name make error." + " deviceFilePath=" + deviceFilePath)
			break
		}

		err = getFPGAReconfigurationCR(ctx,
			mng,
			crNameSpace,
			crName,
			&crFPGAReconfigurationData)
		if errors.IsNotFound(err) {
			err = nil
			break
		} else if nil != err {
			break
		}
		if FPGARECONFSTATUSSUCCEEDED == crFPGAReconfigurationData.Status.Status ||
			FPGARECONFSTATUSFAILED == crFPGAReconfigurationData.Status.Status {
			status = crFPGAReconfigurationData.Status.Status
		}

	}

	return status, err
}

// delete fpga reconfiguration function
func DeleteFPGAReconfiguration(
	ctx context.Context,
	mng ctrl.Manager,
	nodeName string,
	deviceFilePath string) error {

	var err error
	err = nil

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		logger := ctxzap.Extract(ctx)
		c := mng.GetClient()

		var crFPGAReconfigurationData FPGAReconfiguration
		var crNameSpace string
		var crName string

		err = makeCRName(ctx, nodeName, deviceFilePath, &crNameSpace, &crName)
		if nil != err {
			logger.Info("FPGAReconfigurationCR namespace/name make error." + " deviceFilePath=" + deviceFilePath)
			break
		}

		err = getFPGAReconfigurationCR(ctx,
			mng,
			crNameSpace,
			crName,
			&crFPGAReconfigurationData)
		if errors.IsNotFound(err) {
			err = nil
			break
		} else if nil != err {
			break
		}

		crData := &unstructured.Unstructured{}
		crData.SetGroupVersionKind(schema.GroupVersionKind{
			Version: "example.com/v1",
			Kind:    "FPGAReconfiguration",
		})

		crData.SetName(crName)
		crData.SetNamespace(crNameSpace)

		err = c.Delete(ctx, crData)
		if err != nil {
			logger.Error("faild to delete FPGAReconfigurationCR.",
				zap.Error(err))
		}
	}

	return err
}

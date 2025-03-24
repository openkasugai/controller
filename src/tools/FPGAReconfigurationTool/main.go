/* Copyright 2025 NTT Corporation , FUJITSU LIMITED */

package main

import (
	"context"
	common "example.com/FPGAReconfigurationTool/pkg/common"
	logging "example.com/FPGAReconfigurationTool/pkg/logging"
	paramcheck "example.com/FPGAReconfigurationTool/pkg/paramcheck"
	request "example.com/FPGAReconfigurationTool/pkg/request"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"strconv"
	"strings"
	"time"
)

const (
	WAIT_COUNTMAX = 30
	WAIT_TIME     = 3
)

func main() {
	ctx := context.Background()
	logger := logging.SettingZapLogger()
	ctx = ctxzap.ToContext(ctx, logger)
	mng, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{})
	if err != nil {
		os.Exit(1)
	}

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		var nodeName string
		var deviceFilePath string
		var requestMode int32
		var laneSetting = map[int32]string{}
		boolRet := paramcheck.ParamMain(ctx,
			&nodeName,
			&deviceFilePath,
			&requestMode,
			laneSetting)
		if false == boolRet {
			logger.Info("FPGAReconfiguration paramater error.")
			logger.Info("FPGAReconfiguration was create fail.")
			break
		}
		logger.Info("input parameter")
		logger.Info("    nodeName =" + nodeName)
		logger.Info("    deviceFilePath =" + deviceFilePath)
		logger.Info("    requestMode =" + strconv.Itoa(int(requestMode)))
		for i, v := range laneSetting {
			logger.Info("    laneSetting[" + strconv.Itoa(int(i)) + "] =" + v)
		}

		if common.MODEFPGARESET == requestMode {
			if 0 == len(laneSetting) {
				laneSetting[0] = "fpgafunc-config-filter-resize-high-infer"
				logger.Info("FPGA reset default laneSetting[0] =" + laneSetting[0])
			}
		}

		err = request.CheckNodeNameCR(ctx,
			mng,
			nodeName)
		if errors.IsNotFound(err) {
			logger.Info("FPGAReconfiguration was create fail.")
			break
		} else if nil != err {
			logger.Error("FPGAReconfiguration was create fail.", zap.Error(err))
			break
		}

		err = request.CreateFPGAReconfiguration(ctx,
			mng,
			nodeName,
			deviceFilePath,
			requestMode,
			laneSetting)
		if nil != err {
			if true == strings.Contains(err.Error(), "UUID not found.") {
				logger.Info("FPGAReconfiguration was create fail.")
				break
			}
			logger.Error("FPGAReconfiguration was create fail.", zap.Error(err))
			break
		}
		logger.Info("FPGAReconfiguration was create success.")

		for waitCnt := 0; WAIT_COUNTMAX > waitCnt; waitCnt++ {

			status, err := request.GetStatusFPGAReconfiguration(ctx,
				mng,
				nodeName,
				deviceFilePath)
			if nil != err {
				logger.Error("FPGAReconfiguration get status fail.", zap.Error(err))
				break
			}

			if request.FPGARECONFSTATUSSUCCEEDED == status {
				logger.Info("FPGAReconfiguration was execution success.")
				break
			} else if request.FPGARECONFSTATUSFAILED == status {
				logger.Info("FPGAReconfiguration was execution fail. See FPGAFunctionCRC log for details of failure.")
				break
			}

			time.Sleep(WAIT_TIME * time.Second)

			if (WAIT_COUNTMAX - 1) == waitCnt {
				logger.Info("FPGAReconfiguration was execution timeout.")
			}
		}

		err = request.DeleteFPGAReconfiguration(ctx,
			mng,
			nodeName,
			deviceFilePath)
		if nil != err {
			logger.Error("FPGAReconfiguration was delete fail.", zap.Error(err))
			break
		}

	}

	return
}

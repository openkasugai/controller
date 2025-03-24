/* Copyright 2025 NTT Corporation , FUJITSU LIMITED */

package paramcheck

import (
	"context"
	common "example.com/FPGAReconfigurationTool/pkg/common"
	"flag"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"os"
	"strconv"
)

func ParamMain(ctx context.Context,
	pNodeName *string,
	pDeviceFilePath *string,
	requestMode *int32,
	laneSetting map[int32]string) bool {

	boolRet := true

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		logger := ctxzap.Extract(ctx)

		var nodeName string
		nodeName = ""
		var deviceFilePath string
		deviceFilePath = ""
		var resetFlag string

		paramGet(&nodeName, &deviceFilePath, laneSetting, &resetFlag)
		if "" == nodeName {
			logger.Info("not found nodename parameter.")
			boolRet = false
			break
		}
		*pNodeName = nodeName
		if "" == deviceFilePath {
			logger.Info("not found devicefilepath parameter.")
			boolRet = false
			break
		}
		*pDeviceFilePath = deviceFilePath

		if "FPGA" == resetFlag {
			*requestMode = common.MODEFPGARESET
		} else if "ChildBs" == resetFlag {
			*requestMode = common.MODECHILDBSRESET
		} else if "default" == resetFlag {
			*requestMode = common.MODEMANUALSET
		} else {
			logger.Info("unknown reset parameter.")
			boolRet = false
			break
		}

		laneCnt := int32(0)
		for _, value := range laneSetting {
			if "" != value {
				laneCnt++
			}
		}
		if common.MODEMANUALSET == *requestMode {
			if int32(0) == laneCnt {
				logger.Info("not found configmap name parameter.")
				boolRet = false
				break
			}
		}
	}

	return boolRet
}

func paramGet(nodeName *string, deviceFilePath *string, laneSetting map[int32]string, resetFlag *string) {

	for i, v := range os.Args {
		if 0 == i {
		} else if 1 == i {
			*nodeName = v
		} else if 2 == i {
			*deviceFilePath = v
		} else {
			break
		}
	}

	if len(os.Args) >= 3 {
		var prm_resetflag *string
		var prm_lanesetting = map[int32]*string{}

		fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		for i := 0; i < 10; i++ {
			str := "l" + strconv.Itoa(i)
			prm_lanesetting[int32(i)] = fs.String(str, "", "ConfigName used to configure lane")
		}
		prm_resetflag = fs.String("reset", "default", "ConfigName used to configure lane")

		fs.Parse(os.Args[3:])

		if nil != prm_resetflag {
			*resetFlag = *prm_resetflag
		}

		for i, v := range prm_lanesetting {
			if nil != v && "" != *v {
				laneSetting[i] = *v
			}
		}
	}

	return
}

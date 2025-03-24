/*
Copyright 2025 NTT Corporation, FUJITSU LIMITED

FPGACheckPerDF functions
*/

package dfcheck

import (
	"context"
	"encoding/json"
	examplecomv1 "example.com/FPGACheckPerDF/api/v1"
	dfchaincheckercom "example.com/FPGACheckPerDF/pkg/common"
	"fmt"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"os/exec"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
	"syscall"
)

var gDfKind string
var gWbcKind string
var gWbfKind string
var gDfAPIVersion string
var gWbcAPIVersion string
var gWbfAPIVersion string

const (
	RETNORMAL      = 0
	RETCONFIGERROR = 1
	RETNOTFOUND    = 2
	RETGETERROR    = 4
)

func ShellDeleteAllMain(ctx context.Context, cCommand string) error {

	var err error

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		var cShellDirPath string
		cShellDirPath = dfchaincheckercom.SHELLFOLDER

		files, _ := ioutil.ReadDir(cShellDirPath)

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if !strings.HasPrefix(file.Name(), dfchaincheckercom.SHELLNAME) {
				continue
			}
			if !strings.HasSuffix(file.Name(), dfchaincheckercom.SHELLEXTENSION) {
				continue
			}
			cShellFilePath := filepath.Join(cShellDirPath, file.Name())

			fileDeleteFunc(ctx, cShellFilePath)
		}
	}
	return err
}

func LogDeleteAllMain(ctx context.Context, cCommand string) error {

	var err error

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		var cLogDirPath string
		cLogDirPath = dfchaincheckercom.LOGFOLDER

		files, _ := ioutil.ReadDir(cLogDirPath)

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if !strings.HasPrefix(file.Name(), dfchaincheckercom.LOGNAME) {
				continue
			}
			if !strings.HasSuffix(file.Name(), dfchaincheckercom.LOGEXTENSION) {
				continue
			}
			cLogFilePath := filepath.Join(cLogDirPath, file.Name())

			fileDeleteFunc(ctx, cLogFilePath)
		}
	}
	return err
}

// dataflow list collect function
func CollectMain(ctx context.Context,
	mng ctrl.Manager,
	cCommand string,
	cCommandMode string,
	cNodeName string,
	crNamespace string) int {

	retInt := RETNORMAL

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		err := LoadConfigMap(ctx, mng)
		if nil != err {
			fmt.Println("ConfigMap load error.")
			retInt = RETCONFIGERROR
			break
		}

		m := mng.GetAPIReader()

		crDfList := examplecomv1.DataFlowList{}
		err = m.List(ctx, &crDfList, &client.ListOptions{Namespace: crNamespace})
		if nil != err {
			fmt.Println("unable to fetch dataflow list custom resource. Namespace="+
				crNamespace+". ", err)
			break
		}

		for i := 0; i < len(crDfList.Items); i++ {
			crDfData := crDfList.Items[i]
			crDfName := crDfList.Items[i].ObjectMeta.Name
			crDfNamespace := crDfList.Items[i].ObjectMeta.Namespace
			if "default" != crNamespace {
				if crNamespace != crDfNamespace {
					continue
				}
			}
			_ = collectFunc(ctx, mng, cCommand, cCommandMode, cNodeName, crDfNamespace, crDfName, crDfData)
		}
		if 0 == len(crDfList.Items) {
			if "check" != cCommandMode {
				fmt.Println("get dataflow is nothing.")
			}
			break
		}
	}

	return retInt
}

func ShellExecMain(ctx context.Context, cCommand string) error {

	var err error

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		var cShellDirPath string
		cShellDirPath = dfchaincheckercom.SHELLFOLDER

		files, _ := ioutil.ReadDir(cShellDirPath)

		if 0 == len(files) {
			fmt.Println("check dataflow is nothing. because dataflow is no change.")
			break
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if !strings.HasPrefix(file.Name(), dfchaincheckercom.SHELLNAME) {
				continue
			}
			if !strings.HasSuffix(file.Name(), dfchaincheckercom.SHELLEXTENSION) {
				continue
			}
			cShellFilePath := filepath.Join(cShellDirPath, file.Name())
			cmd := exec.Command(cShellFilePath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if nil != err {
				fmt.Println("shell file exec error."+
					" path="+cShellFilePath+". ", err)
				break
			}
		}
	}
	return err
}

// dataflow collect function
func collectFunc(ctx context.Context,
	mng ctrl.Manager,
	cCommand string,
	cCommandMode string,
	cNodeName string,
	crDfNamespace string,
	crDfName string,
	crDf examplecomv1.DataFlow) int {

	gDfKind = "dataflow"
	gWbcKind = "WBConnection"
	gWbfKind = "WBFunction"
	gDfAPIVersion = "example.com/v1"
	gWbcAPIVersion = "example.com/v1"
	gWbfAPIVersion = "example.com/v1"

	retCode := RETNORMAL
	hitFlag := false
	endFlag := false
	errFlag := false
	crMapTableIdx := 0
	fpgaCnt := 0
	myNodeCnt := 0
	copyFlag := false
	copyCnt := 0
	procCnt := 0

	var err error
	var crMapTable map[int]dfchaincheckercom.CrTable
	crMapTable = make(map[int]dfchaincheckercom.CrTable)
	var curCrTable dfchaincheckercom.CrTable

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		// var crDf examplecomv1.DataFlow

		for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
			// get DataflowCR.
			curCrTable = dfchaincheckercom.CrTable{}
			curCrTable.APIVersion = gDfAPIVersion
			curCrTable.Kind = gDfKind
			curCrTable.Namespace = crDfNamespace
			curCrTable.Name = crDfName
			curCrTable.CrData.CrDf = crDf
			curCrTable.Status = crDf.Status.Status
		}
		crMapTable[crMapTableIdx] = curCrTable
		crMapTableIdx++
		if RETNORMAL != retCode {
			errFlag = true
			break
		}

		funcFlag := false

		nextCRFuncKey := "wb-start-of-chain"
		nextCRName := "wb-start-of-chain"
		for {
			retCode = RETNORMAL
			if !funcFlag {
				hitFlag = false
				// search start chain
				var crWBConn examplecomv1.WBConnection
				var iSchConnIdx int
				for iSchConnIdx = 0; iSchConnIdx < len(crDf.Status.ScheduledConnections); iSchConnIdx++ {

					dfSchConValue := crDf.Status.ScheduledConnections[iSchConnIdx]

					if dfSchConValue.From.FunctionKey != nextCRFuncKey {
						continue
					}
					if procCnt > copyCnt {
						copyCnt++
						continue
					}
					nextCRFuncKey = dfSchConValue.To.FunctionKey

					wbconnName := crDfName + "-wbconnection-" +
						dfSchConValue.From.FunctionKey + "-" + dfSchConValue.To.FunctionKey

					for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
						// get WBConnectionCR.
						curCrTable = dfchaincheckercom.CrTable{}
						curCrTable.APIVersion = gWbcAPIVersion
						curCrTable.Kind = gWbcKind
						curCrTable.Namespace = crDfNamespace
						curCrTable.Name = wbconnName

						err = getWBConnData(ctx, mng,
							wbconnName, crDfNamespace, gWbcKind, gWbcAPIVersion,
							&crWBConn)
						if errors.IsNotFound(err) {
							curCrTable.Status = "NotFound"
							retCode = RETGETERROR
							break
						} else if nil != err {
							fmt.Println(
								"get WBConnection CustomResource error."+
									" Name="+wbconnName+". ", err)
							curCrTable.Status = "GetError"
							retCode = RETGETERROR
							break
						}
						curCrTable.CrData.CrWBConn = crWBConn
						curCrTable.Status = string(crWBConn.Status.Status)

						nextCRName = crWBConn.Spec.To.WBFunctionRef.Name
						funcFlag = true
						hitFlag = true
						if strings.HasPrefix(
							dfSchConValue.From.FunctionKey, "wb-start-of-chain") {
							break
						} else if strings.HasPrefix(
							dfSchConValue.To.FunctionKey, "wb-end-of-chain") {
							endFlag = true
							break
						} else {
							if strings.HasPrefix(
								dfSchConValue.To.FunctionKey, "copy-branch-main") {
								copyFlag = true
							}

							ret := convWBConnType(ctx,
								&curCrTable.CrData.CrWBConn,
								&curCrTable.KindSub,
								&curCrTable.APIVersionSub)
							if !ret {
								break
							}

							curCrTable.ChainFlag = true
							if "EthernetConnection" == curCrTable.KindSub {
								err = getEthConnData(ctx, mng,
									wbconnName, crDfNamespace,
									curCrTable.KindSub,
									curCrTable.APIVersionSub,
									&curCrTable.CrData.CrEthConn)
								if errors.IsNotFound(err) {
									curCrTable.StatusSub = "NotFound"
									retCode = RETNOTFOUND
									break
								} else if nil != err {
									fmt.Println(
										"get EthernetConnection CustomResource error."+
											" Name="+wbconnName+". ", err)
									curCrTable.StatusSub = "GetError"
									retCode = RETGETERROR
									break
								}
								curCrTable.StatusSub =
									curCrTable.CrData.CrEthConn.Status.Status

							} else if "PCIeConnection" == curCrTable.KindSub {
								err = getPCIeConnData(ctx, mng,
									wbconnName, crDfNamespace,
									curCrTable.KindSub,
									curCrTable.APIVersionSub,
									&curCrTable.CrData.CrPCIeConn)
								if errors.IsNotFound(err) {
									curCrTable.StatusSub = "NotFound"
									retCode = RETNOTFOUND
									break
								} else if nil != err {
									fmt.Println(
										"get PCIeConnection CustomResource error."+
											" Name="+wbconnName+". ", err)
									curCrTable.StatusSub = "GetError"
									retCode = RETGETERROR
									break
								}
								curCrTable.StatusSub =
									curCrTable.CrData.CrPCIeConn.Status.Status
							} else {
								fmt.Println(
									"unknow Connection CustomResource Kind error." +
										" Name=" + wbconnName +
										" Kind=" + curCrTable.KindSub + ". ")
								curCrTable.StatusSub = "Unknow"
								retCode = RETGETERROR
								break
							}
						}
					}
					if hitFlag {
						crMapTable[crMapTableIdx] = curCrTable
						crMapTableIdx++
						break
					} else if RETGETERROR == retCode || RETNOTFOUND == retCode {
						crMapTable[crMapTableIdx] = curCrTable
						crMapTableIdx++
						break
					}
				}
				if !hitFlag {
					if len(crDf.Status.ScheduledConnections) == iSchConnIdx {
						if strings.HasPrefix(nextCRFuncKey, "wb-start-of-chain") {
							crPrevTable := crMapTable[0]
							crPrevTable.Status = "NoStartPt"
							crMapTable[0] = crPrevTable
						} else if strings.HasPrefix(nextCRFuncKey, "copy-branch-main") {
							hitFlag = true
							endFlag = true
						} else if !strings.HasPrefix(nextCRFuncKey, "wb-end-of-chain") {
							crPrevTable := crMapTable[0]
							crPrevTable.Status = "NoEndPt"
							crMapTable[0] = crPrevTable
						}
					}
					break
				}
				if endFlag {
					procCnt++
					if copyFlag {
						nextCRFuncKey = "copy-branch-main"
						funcFlag = false
						endFlag = false
						copyCnt = 0
					} else {
						break
					}
				}
			} else {

				var crWBFunc examplecomv1.WBFunction
				var curCrTable dfchaincheckercom.CrTable

				for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
					curCrTable = dfchaincheckercom.CrTable{}
					curCrTable.APIVersion = gWbfAPIVersion
					curCrTable.Kind = gWbfKind
					curCrTable.Namespace = crDfNamespace
					curCrTable.Name = nextCRName

					err = getWBFuncData(ctx, mng,
						nextCRName, crDfNamespace, gWbfKind, gWbfAPIVersion,
						&crWBFunc)
					if errors.IsNotFound(err) {
						curCrTable.Status = "NotFound"
						retCode = RETNOTFOUND
						break
					} else if nil != err {
						fmt.Println(
							"get WBFunction CustomResource error."+
								" Name="+nextCRName+". ", err)
						curCrTable.Status = "GetError"
						retCode = RETGETERROR
						break
					}
					if nil != err {
						break
					}
					curCrTable.Status = string(crWBFunc.Status.Status)
					curCrTable.CrData.CrWBFunc = crWBFunc

					ret := convWBFuncType(ctx,
						&crWBFunc, &curCrTable.KindSub,
						&curCrTable.APIVersionSub)
					if !ret {
						break
					}
					curCrTable.ChainFlag = true
					if "FPGAFunction" == curCrTable.KindSub {
						err = getFPGAFuncData(ctx, mng,
							nextCRName, crDfNamespace,
							curCrTable.KindSub, curCrTable.APIVersionSub,
							&curCrTable.CrData.CrFPGAFunc)
						if errors.IsNotFound(err) {
							curCrTable.StatusSub = "NotFound"
							retCode = RETNOTFOUND
							break
						} else if nil != err {
							fmt.Println(
								"get FPGAFunction CustomResource error."+
									" Name="+nextCRName+". ", err)
							curCrTable.StatusSub = "GetError"
							retCode = RETGETERROR
							break
						}
						curCrTable.StatusSub =
							curCrTable.CrData.CrFPGAFunc.Status.Status
						fpgaCnt++
						if cNodeName == curCrTable.CrData.CrFPGAFunc.Spec.NodeName {
							myNodeCnt++
							curCrTable.FPGAPosition = dfchaincheckercom.FPGAPOSITIONPERSONAL
							if 2 < crMapTableIdx {
								if "FPGAFunction" != crMapTable[crMapTableIdx-2].KindSub {
									crPrevTable := crMapTable[crMapTableIdx-2]
									crPrevTable.FPGAPosition =
										dfchaincheckercom.FPGAPOSITIONPREVIOUS
									crMapTable[crMapTableIdx-2] = crPrevTable
								}
							}
						}
					} else if "CPUFunction" == curCrTable.KindSub {
						err = getCPUFuncData(ctx, mng,
							nextCRName, crDfNamespace,
							curCrTable.KindSub, curCrTable.APIVersionSub,
							&curCrTable.CrData.CrCPUFunc)
						if errors.IsNotFound(err) {
							curCrTable.StatusSub = "NotFound"
							retCode = RETNOTFOUND
							break
						} else if nil != err {
							fmt.Println(
								"get CPUFunction CustomResource error."+
									" Name="+nextCRName+". ", err)
							curCrTable.StatusSub = "GetError"
							retCode = RETGETERROR
							break
						}
						curCrTable.StatusSub =
							curCrTable.CrData.CrCPUFunc.Status.Status
						if 2 < crMapTableIdx {
							if "FPGAFunction" == crMapTable[crMapTableIdx-2].KindSub &&
								dfchaincheckercom.FPGAPOSITIONPERSONAL == crMapTable[crMapTableIdx-2].FPGAPosition {
								curCrTable.FPGAPosition = dfchaincheckercom.FPGAPOSITIONNEXT
							}
						}
					} else if "GPUFunction" == curCrTable.KindSub {
						err = getGPUFuncData(ctx, mng,
							nextCRName, crDfNamespace,
							curCrTable.KindSub, curCrTable.APIVersionSub,
							&curCrTable.CrData.CrGPUFunc)
						if errors.IsNotFound(err) {
							curCrTable.StatusSub = "NotFound"
							retCode = RETNOTFOUND
							break
						} else if nil != err {
							fmt.Println(
								"get GPUFunction CustomResource error."+
									" Name="+nextCRName+". ", err)
							curCrTable.StatusSub = "GetError"
							retCode = RETGETERROR
							break
						}
						curCrTable.StatusSub =
							curCrTable.CrData.CrGPUFunc.Status.Status
						if 2 < crMapTableIdx {
							if "FPGAFunction" == crMapTable[crMapTableIdx-2].KindSub &&
								dfchaincheckercom.FPGAPOSITIONPERSONAL == crMapTable[crMapTableIdx-2].FPGAPosition {
								curCrTable.FPGAPosition = dfchaincheckercom.FPGAPOSITIONNEXT
							}
						}
					} else {
						fmt.Println(
							"unknow Connection CustomResource Kind error." +
								" Name=" + nextCRName +
								" Kind=" + curCrTable.KindSub + ". ")
						curCrTable.StatusSub = "Unknow"
						retCode = RETGETERROR
						break
					}
				}
				funcFlag = false
				crMapTable[crMapTableIdx] = curCrTable
				crMapTableIdx++
			}
			if RETNORMAL != retCode && RETNOTFOUND != retCode {
				break
			}
			if RETNOTFOUND == retCode {
				errFlag = true
			}
		}
		if !hitFlag {
			break
		}
		if endFlag {
			break
		}
	}
	if "get" == cCommandMode {
		prStatus := "OK"
		if errFlag {
			prStatus = "NG"
		} else if !hitFlag {
			prStatus = "NG"
		} else if !endFlag {
			prStatus = "NG"
		} else if 0 == fpgaCnt {
			prStatus = "NotFPGA"
		} else if 0 == myNodeCnt {
			prStatus = "OtherNode"
		}
		fmt.Println("Namespace:" + crMapTable[0].Namespace +
			"  Name:" + crMapTable[0].Name +
			"  " + prStatus)
		if 0 == retCode && hitFlag && endFlag && 0 != fpgaCnt && 0 != myNodeCnt {
			cShellFilePath, err := shellMakeFunc(ctx, cCommand, crMapTable)
			if nil == err {
				os.Chmod(cShellFilePath, 0775)
			}
		}
	}
	if "check" == cCommandMode {
		if 0 == retCode && hitFlag && endFlag && 0 != fpgaCnt && 0 != myNodeCnt {
			shellDeleteFunc(ctx, cCommand, crMapTable)
		}
	}
	if "debug" == cCommandMode {
		printFunc(ctx, crMapTable)
	}
	return retCode
}

func shellMakeFunc(ctx context.Context,
	cCommand string,
	crMapTable map[int]dfchaincheckercom.CrTable) (string, error) {

	var err error
	var crDfName string
	fpgaProcFlag := false

	crDfName = crMapTable[0].Name
	var cShellFilePath string
	var cLogFilePath string
	var cFPGACommand string
	cFPGACommand = dfchaincheckercom.FPGACOMMANDPATH
	cShellFilePath = dfchaincheckercom.SHELLFOLDER +
		dfchaincheckercom.SHELLNAME + crDfName +
		dfchaincheckercom.SHELLEXTENSION
	cLogFilePath = dfchaincheckercom.LOGFOLDER +
		dfchaincheckercom.LOGNAME + crDfName +
		dfchaincheckercom.LOGEXTENSION

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop

		defaultUmask := syscall.Umask(0)
		err := os.MkdirAll(dfchaincheckercom.SHELLFOLDER, 0775)
		syscall.Umask(defaultUmask)
		if nil != err {
			fmt.Println("shell dir make error."+
				" dir="+dfchaincheckercom.SHELLFOLDER+". ", err)
			break
		}

		f, err := os.Create(cShellFilePath)
		if nil != err {
			fmt.Println("shell file create error."+
				" file="+cShellFilePath+". ", err)
			break
		}
		defer f.Close()

		recordData := "#!/usr/bin/bash\n"
		_, err = f.WriteString(recordData)
		if nil != err {
			fmt.Println("shell header file write error."+
				" file="+cShellFilePath+". ", err)
			break
		}
		recordData = "echo Namespace:" + crMapTable[0].Namespace +
			"  Name:" + crMapTable[0].Name +
			" > " + cLogFilePath + "\n"
		_, err = f.WriteString(recordData)
		if nil != err {
			fmt.Println("Namespace & Name shell file write error."+
				" file="+cShellFilePath+". ", err)
			break
		}

		for crMapTableIndex := 0; crMapTableIndex < len(crMapTable); crMapTableIndex++ {
			crTable := crMapTable[crMapTableIndex]

			if dfchaincheckercom.FPGAPOSITIONNONE == crTable.FPGAPosition {
				continue
			}
			if "CPUFunction" == crTable.KindSub ||
				"GPUFunction" == crTable.KindSub {

				prGress := "ingress"
				if fpgaProcFlag {
					fpgaProcFlag = false
					prGress = "egress "
				}
				if "CPUFunction" == crTable.KindSub {
					if nil != crTable.CrData.CrCPUFunc.Spec.SharedMemory {
						recordData = "echo \"  #LLDMA " + prGress +
							" Command-Result:" + cFPGACommand +
							" -k " + crTable.CrData.CrCPUFunc.Spec.SharedMemory.CommandQueueID + "\"" +
							" >> " + cLogFilePath + "\n"
						_, err = f.WriteString(recordData)
						if nil != err {
							fmt.Println("echo LLDMA "+prGress+" shell file write error."+
								" file="+cShellFilePath+". ", err)
							break
						}
						recordData = cFPGACommand +
							" -k " + crTable.CrData.CrCPUFunc.Spec.SharedMemory.CommandQueueID +
							" >> " + cLogFilePath + "\n"
						_, err = f.WriteString(recordData)
						if nil != err {
							fmt.Println("LLDMA "+prGress+" shell file write error."+
								" file="+cShellFilePath+". ", err)
							break
						}
					}
				} else {
					if nil != crTable.CrData.CrGPUFunc.Spec.SharedMemory {
						recordData = "echo \"  #LLDMA " + prGress +
							" Command-Result:" + cFPGACommand +
							" -k " + crTable.CrData.CrGPUFunc.Spec.SharedMemory.CommandQueueID + "\"" +
							" >> " + cLogFilePath + "\n"
						_, err = f.WriteString(recordData)
						if nil != err {
							fmt.Println("echo LLDMA "+prGress+" shell file write error."+
								" file="+cShellFilePath+". ", err)
							break
						}
						recordData = cFPGACommand +
							" -k " + crTable.CrData.CrGPUFunc.Spec.SharedMemory.CommandQueueID +
							" >> " + cLogFilePath + "\n"
						_, err = f.WriteString(recordData)
						if nil != err {
							fmt.Println("LLDMA "+prGress+" shell file write error."+
								" file="+cShellFilePath+". ", err)
							break
						}
					}
				}
			} else if "FPGAFunction" == crTable.KindSub {

				fpgaProcFlag = true

				cDeviceFilePath := crTable.CrData.CrFPGAFunc.Spec.AcceleratorIDs[0].ID
				var cDeviceUUID string
				cDeviceUUID = ""
				if true == strings.Contains(cDeviceFilePath, "/dev/xpcie_") {
					cDeviceUUID = strings.ReplaceAll(cDeviceFilePath, "/dev/xpcie_", "")
				}

				cRegionName := crTable.CrData.CrFPGAFunc.Spec.RegionName
				var cLaneNo string
				cLaneNo = ""
				if true == strings.Contains(cRegionName, "lane") {
					cLaneNo = strings.ReplaceAll(cRegionName, "lane", "")
				}
				cFunctionChannelID :=
					strconv.Itoa(int(crTable.CrData.CrFPGAFunc.Status.FunctionChannelID))

				if "" != cDeviceUUID && "" != cLaneNo {

					recordData = "echo \"  #CHAIN " + "ingress" +
						" Command-Result:" + cFPGACommand +
						" -d " + cDeviceUUID +
						" -l " + cLaneNo +
						" -f " + cFunctionChannelID +
						" --dir ingress\"" +
						" >> " + cLogFilePath + "\n"
					_, err = f.WriteString(recordData)
					if nil != err {
						fmt.Println("echo CHAIN ingress shell file write error."+
							" file="+cShellFilePath+". ", err)
						break
					}
					recordData = cFPGACommand +
						" -d " + cDeviceUUID +
						" -l " + cLaneNo +
						" -f " + cFunctionChannelID +
						" --dir ingress" +
						" >> " + cLogFilePath + "\n"
					_, err = f.WriteString(recordData)
					if nil != err {
						fmt.Println("CHAIN ingress shell file write error."+
							" file="+cShellFilePath+". ", err)
						break
					}

					recordData = "echo \"  #CHAIN " + "egress " +
						" Command-Result:" + cFPGACommand +
						" -d " + cDeviceUUID +
						" -l " + cLaneNo +
						" -f " + cFunctionChannelID +
						" --dir egress\"" +
						" >> " + cLogFilePath + "\n"
					_, err = f.WriteString(recordData)
					if nil != err {
						fmt.Println("echo CHAIN egress shell file write error."+
							" file="+cShellFilePath+". ", err)
						break
					}
					recordData = cFPGACommand +
						" -d " + cDeviceUUID +
						" -l " + cLaneNo +
						" -f " + cFunctionChannelID +
						" --dir egress" +
						" >> " + cLogFilePath + "\n"
					_, err = f.WriteString(recordData)
					if nil != err {
						fmt.Println("CHAIN egress shell file write error."+
							" file="+cShellFilePath+". ", err)
						break
					}
				}
			}
		}
		if nil == err {
			recordData = "grep -e Namespace: -e Result " + cLogFilePath
			_, err = f.WriteString(recordData)
			if nil != err {
				fmt.Println("result get shell file write error."+
					" file="+cShellFilePath+". ", err)
				break
			}
		}
	}
	return cShellFilePath, err
}

func shellDeleteFunc(ctx context.Context,
	cCommand string,
	crMapTable map[int]dfchaincheckercom.CrTable) {

	crDfName := crMapTable[0].Name

	var cShellFilePath string
	cShellFilePath = dfchaincheckercom.SHELLFOLDER +
		dfchaincheckercom.SHELLNAME + crDfName +
		dfchaincheckercom.SHELLEXTENSION

	fileDeleteFunc(ctx, cShellFilePath)

	return
}

func fileDeleteFunc(ctx context.Context,
	cFilePath string) {
	if "/" != cFilePath {
		defer os.Remove(cFilePath)
	}
	return
}

// dataflow print function
func printFunc(ctx context.Context,
	crMapTable map[int]dfchaincheckercom.CrTable) {

	fmt.Println("KIND   STATUS     NAMESPACE  NAME")

	for crMapTableIndex := 0; crMapTableIndex < len(crMapTable); crMapTableIndex++ {
		crTable := crMapTable[crMapTableIndex]
		prKind := ""
		if gDfKind == crTable.Kind {
			prKind = "DF    "
		} else if gWbcKind == crTable.Kind {
			prKind = " WBC  "
		} else if gWbfKind == crTable.Kind {
			prKind = " WBF  "
		} else {
			prKind = " ???  "
		}
		prStatus := crTable.Status + "          "
		prStatus = prStatus[:10]
		prNamespace := crTable.Namespace + "          "
		prNamespace = prNamespace[:10]
		prName := crTable.Name
		fmt.Println(prKind + " " + prStatus + " " + prNamespace + " " + prName)
		if crTable.ChainFlag {
			prKindSub := ""
			if "EthernetConnection" == crTable.KindSub {
				prKindSub = "  ETH "
			} else if "PCIeConnection" == crTable.KindSub {
				prKindSub = "  PCIE"
			} else if "FPGAFunction" == crTable.KindSub {
				prKindSub = "  FPGA"
			} else if "CPUFunction" == crTable.KindSub {
				prKindSub = "  CPU "
			} else if "GPUFunction" == crTable.KindSub {
				prKindSub = "  GPU "
			} else {
				prKindSub = "  ????"
			}
			prStatusSub := crTable.StatusSub + "          "
			prStatus = prStatusSub[:10]
			fmt.Println(prKindSub + " " + prStatus + " " + prNamespace + " " + prName)
		}
	}
}

// Create Connection CR attribute information and obtain it
func convWBConnType(
	ctx context.Context,
	pCRData *examplecomv1.WBConnection,
	pConnectionKind *string,
	pConnectionAPIVersion *string) bool {

	var ret bool = false

	for count := 0; count < len(gConnectionKindMap); count++ {
		if gConnectionKindMap[count].ConnectionMethod ==
			pCRData.Spec.ConnectionMethod {
			*pConnectionKind = gConnectionKindMap[count].ConnectionCRKind
			*pConnectionAPIVersion = gWbcAPIVersion
		}
	}

	if *pConnectionAPIVersion == "" || *pConnectionKind == "" {
		fmt.Println("Convert Resource Kind Error." +
			" kind=" + *pConnectionKind +
			" apivesion=" + *pConnectionAPIVersion)
	} else {
		ret = true
	}
	return ret
}

// Create CR attribute information and obtain it
func convWBFuncType(
	ctx context.Context,
	pCRData *examplecomv1.WBFunction,
	pFunctionKind *string,
	pFunctionAPIVersion *string) bool {

	var ret bool = false

	for count := 0; count < len(gFunctionKindMap); count++ {
		if gFunctionKindMap[count].DeviceType == pCRData.Spec.DeviceType {
			*pFunctionKind = gFunctionKindMap[count].FunctionCRKind
			*pFunctionAPIVersion = gWbfAPIVersion
		}
	}
	if *pFunctionAPIVersion == "" || *pFunctionKind == "" {
		fmt.Println("Convert Resource Kind Error." +
			" kind=" + *pFunctionKind +
			" apivesion=" + *pFunctionAPIVersion)
	} else {
		ret = true
	}
	return ret
}

// get wbconnection function
func getWBConnData(ctx context.Context,
	mng ctrl.Manager,
	crName string,
	crNamespace string,
	crKind string,
	crVeraion string,
	pDataCR *examplecomv1.WBConnection) error {

	var err error

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		m := mng.GetAPIReader()
		err = m.Get(ctx, client.ObjectKey{
			Namespace: crNamespace,
			Name:      crName}, pDataCR)
		if errors.IsNotFound(err) {
			break
		} else if nil != err {
			fmt.Println("unable to fetch WBConnection custom resource."+
				" Namespace="+crNamespace+
				" Name="+crName+". ", err)
			break
		}
	}
	return err
}

// get wbfunction function
func getWBFuncData(ctx context.Context,
	mng ctrl.Manager,
	crName string,
	crNamespace string,
	crKind string,
	crVeraion string,
	pDataCR *examplecomv1.WBFunction) error {

	var err error

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		m := mng.GetAPIReader()
		err = m.Get(ctx, client.ObjectKey{
			Namespace: crNamespace,
			Name:      crName}, pDataCR)
		if errors.IsNotFound(err) {
			break
		} else if nil != err {
			fmt.Println("unable to fetch WBFunction custom resource."+
				" Namespace="+crNamespace+
				" Name="+crName+". ", err)
			break
		}
	}
	return err
}

// get ethernet connection function
func getEthConnData(ctx context.Context,
	mng ctrl.Manager,
	crName string,
	crNamespace string,
	crKind string,
	crVeraion string,
	pDataCR *examplecomv1.EthernetConnection) error {

	var err error

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		m := mng.GetAPIReader()
		err = m.Get(ctx, client.ObjectKey{
			Namespace: crNamespace,
			Name:      crName}, pDataCR)
		if errors.IsNotFound(err) {
			break
		} else if nil != err {
			fmt.Println("unable to fetch EthernetConnection custom resource."+
				" Namespace="+crNamespace+
				" Name="+crName+". ", err)
			break
		}
	}
	return err
}

// get PCIe connection function
func getPCIeConnData(ctx context.Context,
	mng ctrl.Manager,
	crName string,
	crNamespace string,
	crKind string,
	crVeraion string,
	pDataCR *examplecomv1.PCIeConnection) error {

	var err error

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		m := mng.GetAPIReader()
		err = m.Get(ctx, client.ObjectKey{
			Namespace: crNamespace,
			Name:      crName}, pDataCR)
		if errors.IsNotFound(err) {
			break
		} else if nil != err {
			fmt.Println("unable to fetch PCIeConnection custom resource."+
				" Namespace="+crNamespace+
				" Name="+crName+". ", err)
			break
		}
	}
	return err
}

// get FPGAfunction function
func getFPGAFuncData(ctx context.Context,
	mng ctrl.Manager,
	crName string,
	crNamespace string,
	crKind string,
	crVeraion string,
	pDataCR *examplecomv1.FPGAFunction) error {

	var err error

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		m := mng.GetAPIReader()
		err = m.Get(ctx, client.ObjectKey{
			Namespace: crNamespace,
			Name:      crName}, pDataCR)
		if errors.IsNotFound(err) {
			break
		} else if nil != err {
			fmt.Println("unable to fetch FPGAFunction custom resource."+
				" Namespace="+crNamespace+
				" Name="+crName+". ", err)
			break
		}
	}
	return err
}

// get cpufunction function
func getCPUFuncData(ctx context.Context,
	mng ctrl.Manager,
	crName string,
	crNamespace string,
	crKind string,
	crVeraion string,
	pDataCR *examplecomv1.CPUFunction) error {

	var err error

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		m := mng.GetAPIReader()
		err = m.Get(ctx, client.ObjectKey{
			Namespace: crNamespace,
			Name:      crName}, pDataCR)
		if errors.IsNotFound(err) {
			break
		} else if nil != err {
			fmt.Println("unable to fetch CPUFunction custom resource."+
				" Namespace="+crNamespace+
				" Name="+crName+". ", err)
			break
		}
	}
	return err
}

// get gpufunction function
func getGPUFuncData(ctx context.Context,
	mng ctrl.Manager,
	crName string,
	crNamespace string,
	crKind string,
	crVeraion string,
	pDataCR *examplecomv1.GPUFunction) error {

	var err error

	for doWhile := 0; doWhile < 1; doWhile++ { //nolint:staticcheck // SA4008: This loop is intentionally only a one-time loop
		m := mng.GetAPIReader()
		err = m.Get(ctx, client.ObjectKey{
			Namespace: crNamespace,
			Name:      crName}, pDataCR)
		if errors.IsNotFound(err) {
			break
		} else if nil != err {
			fmt.Println("unable to fetch GPUFunction custom resource."+
				" Namespace="+crNamespace+
				" Name="+crName+". ", err)
			break
		}
	}
	return err
}

type ConnectionKindMap struct {
	ConnectionMethod string `json:"connectionMethod"`
	ConnectionCRKind string `json:"connectionCRKind"`
}
type FunctionKindMap struct {
	DeviceType     string `json:"deviceType"`
	FunctionCRKind string `json:"functionCRKind"`
}

var gConnectionKindMap []ConnectionKindMap
var gFunctionKindMap []FunctionKindMap

type ConfigTable struct {
	name string
}

var configLoadTable = []ConfigTable{
	{
		"connectionkindmap",
	}, {
		"functionkindmap",
	},
}

// Load ConfigMap (for main)
func LoadConfigMap(ctx context.Context, mng ctrl.Manager) error {

	var cfgdata []byte
	var err error

	for _, record := range configLoadTable {
		err = getConfigMap(ctx, mng, record.name, &cfgdata)
		if nil != err {
			break
		}
		if "connectionkindmap" == record.name {
			err = json.Unmarshal(cfgdata, &gConnectionKindMap)
		}
		if "functionkindmap" == record.name {
			err = json.Unmarshal(cfgdata, &gFunctionKindMap)
		}
		if nil != err {
			fmt.Println("unable to Unmarshal."+
				" ConfigMap="+record.name+". ", err)
			break
		}
	}
	return err
}

// Get ConfigMap
func getConfigMap(ctx context.Context,
	mng ctrl.Manager,
	cfgname string, cfgdata *[]byte) error {

	m := mng.GetAPIReader()

	var mapdata map[string]string

	tmpData := &unstructured.Unstructured{}
	tmpData.SetGroupVersionKind(schema.GroupVersionKind{
		Kind:    "ConfigMap",
		Version: "v1",
	})

	// Get a ConfigMap by namespace/name
	err := m.Get(context.Background(),
		client.ObjectKey{
			Namespace: "default",
			Name:      cfgname,
		},
		tmpData)
	if errors.IsNotFound(err) {
		fmt.Println("ConfigMap does not exist."+
			" ConfigName="+cfgname+". ", err)
	} else if nil != err {
		fmt.Println("ConfigMap unable to fetch."+
			" ConfigName="+cfgname+". ", err)
	} else {
		mapdata, _, _ = unstructured.NestedStringMap(tmpData.Object, "data")
		for _, jsonrecord := range mapdata {
			*cfgdata = []byte(jsonrecord)
		}
	}
	return err
}

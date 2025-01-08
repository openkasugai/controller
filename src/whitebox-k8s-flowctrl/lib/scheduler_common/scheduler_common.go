package scheduler_common

import (
	"reflect"
	"sort"
	"strings"

	ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
)

var EndPointKeyword = "wb-end-of-chain"
var StartPointKeyword = "wb-start-of-chain"

// TODO: support other types
const (
	UNDEFINED = iota
	DEVICE
	NODE
	RACK
	FLOOR
	OTHER
)

type FuncKeyInfoStruct struct {
	FuncKey          string
	PrevConnTypeName string
	NextConnTypeName string
}

// struct that holds order and connection of functions
type FuncOrderInfoStruct struct {
	FuncKey       string
	FromFunctions []string
	ToFunctions   []string
	Depth         int32
}

func StrSliceCheck(slice []string, target string) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}

func PutStrSliceNoDupe(slice []string, target string) []string {
	if !StrSliceCheck(slice, target) {
		slice = append(slice, target)
	}
	return slice
}

func GetDereferencedValueOrZeroValue[T any](pointerAddr *T) T {
	var zeroValue T
	if pointerAddr == nil {
		return zeroValue
	} else {
		return *pointerAddr
	}
}

func GetOrElseDefaultValue[T any](value, defaultValue T) T {
	var zeroValue T
	if reflect.DeepEqual(value, zeroValue) {
		return defaultValue
	}
	return value
}

func StringToConst(s string) int {
	switch s {
	case "device":
		return DEVICE
	case "node":
		return NODE
	case "rack":
		return RACK
	case "floor":
		return FLOOR
	default:
		// TODO: fix this to support OTHER pattern
		return OTHER
	}
}

func ConstToString(i int) string {
	switch i {
	case DEVICE:
		return "device"
	case NODE:
		return "node"
	case RACK:
		return "rack"
	case FLOOR:
		return "floor"
	default:
		// TODO: fix this to support OTHER pattern
		return "other"
	}
}

func CheckSortedFunc(checkFuncInfo []FuncKeyInfoStruct, targetStr string) bool {
	for _, funcInfo := range checkFuncInfo {
		if funcInfo.FuncKey == targetStr {
			return true
		}
	}
	return false
}

func GetStartEndPointAvailableInterfaceType(s string) []string {
	switch s {
	case StartPointKeyword:
		return []string{"host100gether"}
	case EndPointKeyword:
		return []string{"host100gether"}
	default:
		return []string{}
	}
}

func GetFieldValuesFromStructs[T any](
	targetSliceOfStructs interface{}, targetFieldName string) []T {

	structs := reflect.ValueOf(targetSliceOfStructs)

	values := make([]T, 0, structs.Len())

	for i := 0; i < structs.Len(); i++ {
		aStruct := structs.Index(i)
		field := aStruct.FieldByName(targetFieldName)
		values = append(values, field.Interface().(T))
	}

	return values
}

// // Sort the functionKeys of FunctionChain.Spec.Functions from start point to end point.
// func GetSortedFunctionInstanceName(fc *ntthpcv1.FunctionChain) []string {

// 	var sortedFuncInstNames []string
// 	var prevFunc string

// // Get the first Function
// 	for _, fcConValue := range fc.Spec.Connections {
// 		if fcConValue.From.FunctionKey == StartPointKeyword {
// 			prevFunc = fcConValue.To.FunctionKey
// 		}
// 	}

// 	isConfirmed := make(map[string]bool)

// 	for sortCnt := 0; sortCnt < len(fc.Spec.Functions); sortCnt++ {
// 		for _, fcConValue := range fc.Spec.Connections {

// 			target := fcConValue.From.FunctionKey
// 			if _, ok := isConfirmed[target]; ok {
// 				continue
// 			}

// 			if target == prevFunc {
// 				sortedFuncInstNames = append(sortedFuncInstNames, target)
// 				prevFunc = fcConValue.To.FunctionKey
// 				isConfirmed[target] = true
// 				break
// 			}

// 		}
// 	}

// 	return sortedFuncInstNames
// }

// Sort FunctionChain.Spec.Functions from start point to end point.
func GetSortedFunctionInstanceOrder(fc *ntthpcv1.FunctionChain) []FuncOrderInfoStruct {
	funcOrderInfoStructMap := make(map[string]*FuncOrderInfoStruct)

	// add connections into each FuncOrderInfoStruct
	for _, connection := range fc.Spec.Connections {
		fromfuncOrderInfoStruct, fromFuncExists := funcOrderInfoStructMap[connection.From.FunctionKey]
		if !fromFuncExists {
			fromfuncOrderInfoStruct = &FuncOrderInfoStruct{FuncKey: connection.From.FunctionKey, FromFunctions: []string{}, ToFunctions: []string{}, Depth: 0}
			funcOrderInfoStructMap[connection.From.FunctionKey] = fromfuncOrderInfoStruct
		}
		tofuncOrderInfoStruct, toFuncExists := funcOrderInfoStructMap[connection.To.FunctionKey]
		if !toFuncExists {
			tofuncOrderInfoStruct = &FuncOrderInfoStruct{FuncKey: connection.To.FunctionKey, FromFunctions: []string{}, ToFunctions: []string{}, Depth: 0}
			funcOrderInfoStructMap[connection.To.FunctionKey] = tofuncOrderInfoStruct
		}
		fromfuncOrderInfoStruct.ToFunctions = append(fromfuncOrderInfoStruct.ToFunctions, connection.To.FunctionKey)
		tofuncOrderInfoStruct.FromFunctions = append(tofuncOrderInfoStruct.FromFunctions, connection.From.FunctionKey)
	}

	// exclude FuncOrderInfoStruct from funcOrderInfoStructMap with partial match for "wb-start-of-chain" or "wb-end-of-chain"
	for funcKey := range funcOrderInfoStructMap {
		if strings.HasPrefix(funcKey, StartPointKeyword) || strings.HasPrefix(funcKey, EndPointKeyword) {
			delete(funcOrderInfoStructMap, funcKey)
		}
	}

	// update the depth for each FuncOrderInfoStruct
	var updateDepths func(*FuncOrderInfoStruct, int32)
	updateDepths = func(FuncOrderInfoStruct *FuncOrderInfoStruct, currentDepth int32) {
		if FuncOrderInfoStruct.Depth < currentDepth {
			FuncOrderInfoStruct.Depth = currentDepth
		}
		for _, nextFuncKey := range FuncOrderInfoStruct.ToFunctions {
			nextfuncOrderInfoStruct, exists := funcOrderInfoStructMap[nextFuncKey]
			if exists {
				updateDepths(nextfuncOrderInfoStruct, currentDepth+1)
			}
		}
	}

	// find and update the depth starting from the start function
	for _, FuncOrderInfoStruct := range funcOrderInfoStructMap {
		isStartFunc := false
		for _, fromFunc := range FuncOrderInfoStruct.FromFunctions {
			if strings.HasPrefix(fromFunc, StartPointKeyword) {
				isStartFunc = true
				break
			}
		}
		if isStartFunc {
			updateDepths(FuncOrderInfoStruct, 1)
		}
	}

	// sort funcOrderInfoStructMap by depth
	sortedfuncOrderInfoList := make([]FuncOrderInfoStruct, 0, len(funcOrderInfoStructMap))
	for _, FuncOrderInfoStruct := range funcOrderInfoStructMap {
		sortedfuncOrderInfoList = append(sortedfuncOrderInfoList, *FuncOrderInfoStruct)
	}
	sort.Slice(sortedfuncOrderInfoList, func(i, j int) bool {
		if sortedfuncOrderInfoList[i].Depth != sortedfuncOrderInfoList[j].Depth {
			return sortedfuncOrderInfoList[i].Depth < sortedfuncOrderInfoList[j].Depth
		}
		return sortedfuncOrderInfoList[i].FuncKey < sortedfuncOrderInfoList[j].FuncKey
	})

	return sortedfuncOrderInfoList
}

func SearchFunction(fi int32, functions *[]ntthpcv1.FunctionCapStruct) *ntthpcv1.FunctionCapStruct {
	for _, f := range *functions {
		if f.FunctionIndex == fi {
			return &f
		}
	}
	return nil
}

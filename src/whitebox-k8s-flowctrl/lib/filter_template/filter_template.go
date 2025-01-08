/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package filter_template

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1" //nolint:stylecheck // ST1019: intentional import as another name
	v1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"       //nolint:stylecheck // ST1019: intentional import as another name
	. "github.com/compsysg/whitebox-k8s-flowctrl/lib/scheduler_common"
	_ "github.com/lib/pq"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	ErrFirstConfigMapNotExist       = errors.New("First ConfigMap Doesn't Not Exist")
	ErrSecondConfigMapNotExist      = errors.New("Second ConfigMap Doesn't Not Exist")
	ErrFilterStatusIsNotFiltering   = errors.New("Filter status is not Filtering")
	ErrCurrentFilterIsNotThisFilter = errors.New("Current filter isn't this filter.")
	errInvalidConfigMapParamter     = errors.New("INvalid ConfigMap Parameter.")
)

const (
	FinishStatus    = "Finish"
	FilteringStatus = "Filtering"
	FailedStatus    = "Failed"
)

// Strategy Reserved Key
const (
	ReferenceParameterKey = "referenceParameter"

	MaxTypeCombinationKey   = "maxTypeCombination"
	MaxTargetCombinationKey = "maxTargetCombination"
)

// UserRequirement Reserved Key
const (
	RequestNodeNamesKey        = "requestNodeNames"
	RequestFunctionTargetKey   = "requestFunctionTargets"
	RequestConnectionTargetKey = "requestConnectionTargets"
	RequestDeviceTypesKey      = "requestDeviceTypes"
	RequestConnectionTypesKey  = "requestConnectionTypes"
	RequestSameNodeGroupsKey   = "requestSameNodeGroups"
)

// UserRequirement Reserved Key for Namespace
const (
	// Default value of FunctionTarget namespace
	defaultFtNameSpace = "default"
	// Default value of ConnectionTarget namespace
	defaultCtNameSpace = "default"
	// Key name when reading FunctionTarget namespace from userRequirement
	ftNameSpaceKey = "functionTargetNameSpace"
	// Key name when reading ConnectionTarget namespace from userRequirement
	ctNameSpaceKey = "connectionTargetNameSpace"
)

type FunctionTargetFilter struct {
	NodeNames            *[]string
	DeviceTypes          *[]string
	FunctionNames        *[]string
	FunctionIndexes      *[]string
	FunctionTargets      *[]string
	RegionNames          *[]string
	RegionTypes          *[]string
	IncludesNotAvailable *bool
}

type functionTargetFilter struct {
	Filter *FunctionTargetFilter
}

type ConnectionTargetFilter struct {
	NodeNames            *[]string
	DeviceTypes          *[]string
	ConnectionTarget     *[]string
	IncludesNotAvailable *bool
}

type connectionTargetFilter struct {
	Filter *ConnectionTargetFilter
}

type DeviceTypeFilter struct {
	DeviceTypes *[]string
}

type deviceTypeFilter struct {
	Filter *DeviceTypeFilter
}

type ConnectionTypeFilter struct {
	ConnectionType *[]string
}

type connectionTypeFilter struct {
	Filter *ConnectionTypeFilter
}

type parseConfigMap struct {
	Prefix string
	Key    string
	Value  string
}

type prefixCheckerStruct struct {
	StringCurrentIndex string
	FilterName         string
}

type independentMembers struct {
	Name            string
	UserRequirement map[string]string
	Strategy        map[string]string
	IsDBConnected   bool
	DBConnection    *db_ConnectionInfoStruct
}

type FilterTemplate struct {
	client.Client
	Scheme         *runtime.Scheme
	Recorder       record.EventRecorder
	RequeueTimeSec int
	initialized    bool
	independent    map[string]independentMembers
	independentMtx sync.Mutex
}

type FunctionIndexStruct struct {
	FunctionTarget *v1.FunctionTarget
	FunctionIndex  *int32
}

// Public functions
// Initialize filter.
// Fetch resources {SchedulingData, UserRequirement, Strategy}
func (r *FilterTemplate) Initialize(
	ctx context.Context,
	req ctrl.Request,
	filterNameCandidates ...string,
) (
	sd *ntthpcv1.SchedulingData,
	df *ntthpcv1.DataFlow,
	err error,
) {

	l := log.FromContext(ctx)

	// Get SchedulingData
	sd, err = r.getSchedulingData(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	// Check if the SchedulingData Status is Filtering
	if !r.checkSchedulingDataStatus(sd) {
		l.Info(fmt.Sprintf("SchedulingData status is not Filtering. Status is %v. Index : %v. len : %v", sd.Status.Status, sd.Status.CurrentFilterIndex, len(sd.Spec.FilterPipeline)))
		return nil, nil, ErrFilterStatusIsNotFiltering
	}

	// Check if the current filter is included in the input filter candidates
	currentFilterName := sd.Spec.FilterPipeline[*sd.Status.CurrentFilterIndex]
	includeCurrentFilter := false
	for _, name := range filterNameCandidates {
		if currentFilterName == name {
			includeCurrentFilter = true
			break
		}
	}
	if !includeCurrentFilter {
		return nil, nil, ErrCurrentFilterIsNotThisFilter
	}

	r.independentMtx.Lock()
	if !r.initialized {
		r.independent = make(map[string]independentMembers)
		r.initialized = true
	}
	r.independentMtx.Unlock()

	// fetch Dataflow
	df, err = r.fetchDataflow(ctx, sd)
	if err != nil {
		return sd, nil, err
	}

	strategy, usrreq, err := r.fetchUserParameters(ctx, df, *sd.Status.CurrentFilterIndex, currentFilterName)
	if err != nil {
		return sd, df, err
	}

	r.independentMtx.Lock()
	r.independent[r.generateKey(sd)] = independentMembers{
		Name:            currentFilterName,
		UserRequirement: usrreq,
		Strategy:        strategy,
		IsDBConnected:   false,
		DBConnection:    &db_ConnectionInfoStruct{},
	}
	r.independentMtx.Unlock()

	return sd, df, nil

}

func (r *FilterTemplate) GetName(sd *ntthpcv1.SchedulingData) string {
	r.independentMtx.Lock()
	name := r.independent[r.generateKey(sd)].Name
	r.independentMtx.Unlock()
	return name
}

func (r *FilterTemplate) GetFunctionChain(
	df *ntthpcv1.DataFlow,
) (*ntthpcv1.FunctionChain, error) {
	// Currently there are no errors, but we will keep this as a return value to allow for the possibility of adding value checks, etc.
	return df.Status.FunctionChain, nil
}

func (r *FilterTemplate) GetFunctionTypes(
	df *ntthpcv1.DataFlow,
) ([]*ntthpcv1.FunctionType, error) {
	return df.Status.FunctionType, nil
}

func (r *FilterTemplate) FetchFunctionTargets(
	ctx context.Context,
	req ctrl.Request,
	sd *ntthpcv1.SchedulingData,
	df *ntthpcv1.DataFlow,
	filter *FunctionTargetFilter,
) (retFTs []*ntthpcv1.FunctionTarget, err error) {

	// Get the default namespace value from a constant
	ftNameSpace := defaultFtNameSpace

	// Check the namespace of FunctionTarget in userRequirement
	if r.DoesUserRequirementHaveKey(sd, ftNameSpaceKey) {
		r.LoadUserRequirementParameter(sd, ftNameSpaceKey, &ftNameSpace)
	}

	var ftl ntthpcv1.FunctionTargetList = ntthpcv1.FunctionTargetList{}
	if err = r.List(ctx, &ftl, &client.ListOptions{Namespace: ftNameSpace}); err != nil {
		return retFTs, err
	}

	var compFTs *[]*v1.FunctionTarget
	retFTs = make([]*v1.FunctionTarget, len(ftl.Items))
	compFTs = &retFTs
	for i, ft := range ftl.Items {
		var getFt ntthpcv1.FunctionTarget = ft
		retFTs[i] = &getFt
	}

	if filter == nil {
		filter = &FunctionTargetFilter{IncludesNotAvailable: valToAddr[bool](false)}
	}

	filterExecuter := functionTargetFilter{Filter: filter}

	if filter.FunctionTargets != nil {
		filterExecuter.FilterFunctionTarget(ctx, compFTs)
	}

	if filter.FunctionNames != nil {
		filterExecuter.FilterFunctionNames(ctx, compFTs)
	}

	if filter.RegionTypes != nil {
		filterExecuter.FilterRegionTypes(ctx, compFTs)
	}

	if filter.NodeNames != nil {
		filterExecuter.FilterNodeNames(ctx, compFTs)
	}

	if filter.DeviceTypes != nil {
		filterExecuter.FilterDeviceTypes(ctx, compFTs)
	}

	if !(filter.IncludesNotAvailable != nil && *filter.IncludesNotAvailable) {
		filterExecuter.FilterAvailableFunctionTargets(compFTs)
	}

	if filter.RegionNames != nil {
		filterExecuter.FilterRegionNames(ctx, compFTs)
	}

	retFTs = *compFTs
	return retFTs, nil
}

func (r *FilterTemplate) FetchFunctionIndexStructs(
	ctx context.Context,
	fts *[]*v1.FunctionTarget,
	filter *FunctionTargetFilter,
) (retFIs []FunctionIndexStruct) {

	var compFIs *[]FunctionIndexStruct
	retFIs = make([]FunctionIndexStruct, 0)
	compFIs = &retFIs

	for _, ft := range *fts {
		retFIs = append(retFIs, FunctionIndexStruct{ft, nil})
		for _, function := range ft.Status.Functions {
			functionIndex := function.FunctionIndex
			retFIs = append(retFIs, FunctionIndexStruct{ft, &functionIndex})
		}
	}

	if filter == nil {
		filter = &FunctionTargetFilter{IncludesNotAvailable: valToAddr[bool](false)}
	}

	filterExecuter := functionTargetFilter{Filter: filter}

	if filter.FunctionIndexes != nil {
		filterExecuter.FilterFunctionIndexes(ctx, compFIs)
	}

	if !(filter.IncludesNotAvailable != nil && *filter.IncludesNotAvailable) {
		filterExecuter.FilterAvailableFunctionIndexStructs(compFIs)
	}

	retFIs = *compFIs
	return retFIs
}

func (r *FilterTemplate) FetchConnectionTarget(
	ctx context.Context,
	req ctrl.Request,
	sd *ntthpcv1.SchedulingData,
	df *ntthpcv1.DataFlow,
	filter *ConnectionTargetFilter,
) (retCTs []*ntthpcv1.ConnectionTarget, err error) {

	// Get the default namespace value from a constant
	ctNameSpace := defaultCtNameSpace

	// Check the namespace of FunctionTarget in userRequirement
	if r.DoesUserRequirementHaveKey(sd, ctNameSpaceKey) {
		r.LoadUserRequirementParameter(sd, ctNameSpaceKey, &ctNameSpace)
	}

	ctl := &ntthpcv1.ConnectionTargetList{}
	if err = r.List(ctx, ctl, &client.ListOptions{Namespace: ctNameSpace}); err != nil {
		return retCTs, err
	}

	var compFTs *[]*v1.ConnectionTarget
	retCTs = make([]*v1.ConnectionTarget, len(ctl.Items))
	compFTs = &retCTs
	for i, ct := range ctl.Items {
		var getFt ntthpcv1.ConnectionTarget = ct
		retCTs[i] = &getFt
	}

	if filter == nil {
		return retCTs, nil
	}

	filterExecuter := connectionTargetFilter{Filter: filter}

	if filter.ConnectionTarget != nil {
		filterExecuter.FilterConnectionTarget(ctx, compFTs)
	}

	if filter.NodeNames != nil {
		filterExecuter.FilterNodeNames(ctx, compFTs)
	}

	if filter.DeviceTypes != nil {
		filterExecuter.FilterDeviceTypes(ctx, compFTs)
	}

	retCTs = *compFTs
	return retCTs, nil
}

func (r *FilterTemplate) FetchDeviceTypes(
	ctx context.Context,
	req ctrl.Request,
	sd *ntthpcv1.SchedulingData,
	df *ntthpcv1.DataFlow,
	filter *DeviceTypeFilter,
) (retDTs []string, err error) {

	fc := df.Status.FunctionChain

	rtCands := make([]string, 0)
	fTyps := df.Status.FunctionType

	// obtain DeviceTypeCandidates per function instance

	isAdded := make(map[string]bool)
	for _, function := range fc.Spec.Functions {

		if _, ok := isAdded[function.FunctionName]; ok {
			continue
		}

		for _, ft := range fTyps {
			if function.FunctionName == ft.Spec.FunctionName {
				rtCands = append(rtCands, ft.Status.RegionTypeCandidates...)
				isAdded[function.FunctionName] = true
				break
			}
		}
	}

	l := log.FromContext(context.Background())
	l.Info(fmt.Sprintf("rtCands %v", rtCands))

	fTgts, err := r.FetchFunctionTargets(ctx, req, sd, df, nil)
	if err != nil {
		return nil, err
	}

	for _, ft := range fTgts {
		l.Info(fmt.Sprintf("ft.Status %v", ft.Status))
	}

	retDTs = make([]string, 0)

	alreadyAppend := make(map[string]bool)
	for _, rtCand := range rtCands {
		for _, ft := range fTgts {
			if rtCand == ft.Status.RegionType {
				dt := ft.Status.DeviceType
				if _, ok := alreadyAppend[dt]; !ok {
					retDTs = append(retDTs, dt)
					alreadyAppend[dt] = true
				}
			}
		}
	}

	if filter == nil {
		return retDTs, nil
	}

	var compFTs *[]string
	compFTs = &retDTs
	filterExecuter := deviceTypeFilter{Filter: filter}

	if filter.DeviceTypes != nil {
		filterExecuter.FilterDeviceTypes(ctx, compFTs)
	}

	retDTs = *compFTs

	l.Info("retDTs", "retDTs", retDTs)

	return retDTs, nil
}

func (r *FilterTemplate) FetchConnectionTypes(
	ctx context.Context,
	req ctrl.Request,
	sd *ntthpcv1.SchedulingData,
	df *ntthpcv1.DataFlow,
	filter *ConnectionTypeFilter,
) (retCTs []*ntthpcv1.ConnectionType, err error) {

	fc := df.Status.FunctionChain

	l := log.FromContext(ctx)

	var ctl ntthpcv1.ConnectionTypeList
	if err := r.List(ctx, &ctl, &client.ListOptions{Namespace: fc.Spec.ConnectionTypeNamespace}); err != nil {
		l.Error(err, "Failed to fetch connectionTypes")
		return retCTs, err
	}

	retCTs = make([]*ntthpcv1.ConnectionType, len(ctl.Items))
	for i, ct := range ctl.Items {
		getCTs := ct
		retCTs[i] = &getCTs
	}

	var compCTs *[]*ntthpcv1.ConnectionType
	compCTs = &retCTs

	if filter == nil {
		return retCTs, nil
	}

	filterExecuter := connectionTypeFilter{Filter: filter}

	if filter.ConnectionType != nil {
		filterExecuter.FilterConnectionTypes(ctx, compCTs)
	}

	retCTs = *compCTs
	return retCTs, nil
}

// Finalize filter
// Update SchedulingData data and release resouces
func (r *FilterTemplate) Finalize(ctx context.Context, sd *ntthpcv1.SchedulingData) (ctrl.Result, error) {

	if err := r.updateResource(ctx, sd); err != nil {
		return ctrl.Result{}, err
	}

	r.dispose(sd)

	if sd.Status.Status == FinishStatus {
		return ctrl.Result{}, nil
	} else {
		return ctrl.Result{RequeueAfter: time.Second * 1}, nil
	}

}

// Abort filter.
// Release resources and change SchedulingData status to "Failed".
func (r *FilterTemplate) Abort(ctx context.Context, sd *ntthpcv1.SchedulingData, err error) (ctrl.Result, error) {

	r.dispose(sd)

	l := log.FromContext(ctx)
	sd.Status.Status = FailedStatus

	if localErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		return r.Status().Update(ctx, sd)
	}); localErr != nil {
		l.Error(localErr, "unable to update SchedulingData spec")
		return ctrl.Result{}, localErr
	}

	return ctrl.Result{}, err
}

func (r *FilterTemplate) Query(
	ctx context.Context,
	query string,
	sd *ntthpcv1.SchedulingData,
) ([]map[string]interface{}, error) {
	var (
		err error
		res []map[string]interface{}
	)

	l := log.FromContext(ctx)

	idm := r.exclusiveRead(sd)
	if !idm.IsDBConnected {

		var DBInfoConfigMapName string = "db-connection-info"
		var DBInfoConfigMapNamespaceCandidates []string = []string{sd.ObjectMeta.Namespace, "default"}
		var roleName string = "pub-filter"

		// get infomation for connection to DB from ConfigMap/Secret
		dbConnectionInfo, err := fetchDB_ConnectionInfo(
			r, ctx,
			DBInfoConfigMapName,
			DBInfoConfigMapNamespaceCandidates,
			roleName,
			DBInfoConfigMapNamespaceCandidates,
		)
		if err != nil {
			l.Error(err, "Failed to get DB connectionInfo")
			return res, err
		}

		// connect to DB
		err = dbConnectionInfo.Open()

		if err != nil {
			l.Error(err, "Failed to connect DB")
			return res, err
		}

		idm.DBConnection = dbConnectionInfo
		idm.IsDBConnected = true

		r.exclusiveWrite(sd, idm)

	}

	res, err = idm.DBConnection.Query(query)

	if err != nil {
		l.Error(err, "Failed to Query")
		return res, err
	}

	return res, nil
}

func (r *FilterTemplate) DoesUserRequirementHaveKey(sd *ntthpcv1.SchedulingData, key string) bool {
	elm := r.exclusiveRead(sd)
	return doesMapHaveKey(elm.UserRequirement, key)
}

func (r *FilterTemplate) DoesStrategyHaveKey(sd *ntthpcv1.SchedulingData, key string) bool {
	elm := r.exclusiveRead(sd)
	return doesMapHaveKey(elm.Strategy, key)
}

func (r *FilterTemplate) LoadUserRequirementParameter(sd *ntthpcv1.SchedulingData, key string, dest interface{}) {
	elm := r.exclusiveRead(sd)
	parseYAML(elm.UserRequirement[key], dest)
}

func (r *FilterTemplate) LoadStrategyParameter(sd *ntthpcv1.SchedulingData, key string, dest interface{}) {
	elm := r.exclusiveRead(sd)
	parseYAML(elm.Strategy[key], dest)
}

func (r *FilterTemplate) UserRequirementLen(sd *ntthpcv1.SchedulingData) int {
	elm := r.exclusiveRead(sd)
	return len(elm.UserRequirement)
}

func (r *FilterTemplate) StrategyLen(sd *ntthpcv1.SchedulingData) int {
	elm := r.exclusiveRead(sd)
	return len(elm.Strategy)
}

// Private functions
// Check status of SchedulingData in Filtering
func (r *FilterTemplate) checkSchedulingDataStatus(sd *ntthpcv1.SchedulingData) bool {
	isFiltering := sd.Status.Status == FilteringStatus
	var idxIsValid bool
	if sd.Status.CurrentFilterIndex != nil {
		idxIsValid = *sd.Status.CurrentFilterIndex < int32(len(sd.Spec.FilterPipeline))
	}
	return isFiltering && idxIsValid
}

// Fetch strategy and return map which is inputted strategy parameters
func (r *FilterTemplate) fetchStrategy(
	ctx context.Context,
	df *ntthpcv1.DataFlow,
	filterIndex int32,
	name string,
	filterName string,
) (map[string]string, error) {
	return r.recursiveFetchConfigMap(
		ctx,
		name,
		filterName,
		[]string{df.Namespace, "default"},
		ReferenceParameterKey,
		filterIndex,
		[]string{df.Namespace, "default"},
	)
}

// Fetch UserRequest and return map which is inputted UserRequest parameters
func (r *FilterTemplate) fetchUserRequirement(
	ctx context.Context,
	df *ntthpcv1.DataFlow,
) (map[string]string, error) {

	retMap := make(map[string]string)
	namespaceCandidates := []string{df.Namespace, "default"}

	cm, err := tryFetchConfigMapFromSeveralNameSpaceCandidates(
		r, ctx, GetDereferencedValueOrZeroValue(df.Spec.UserRequirement),
		namespaceCandidates)
	if err != nil {
		return retMap, ErrFirstConfigMapNotExist
	}

	for k, v := range cm.Data {
		retMap[k] = v
	}

	return retMap, nil
}

// Update SchedulingData data & status.
// If this filter is last of filter pipeline, status be changed "Finish".
func (r *FilterTemplate) updateResource(ctx context.Context, sd *ntthpcv1.SchedulingData) error {
	// Update Process
	// status update
	l := log.FromContext(ctx)

	*sd.Status.CurrentFilterIndex = *sd.Status.CurrentFilterIndex + 1

	if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		return r.Status().Update(ctx, sd)
	}); err != nil {
		l.Error(err, "unable to update SchedulingData status")
		return err
	}

	if *sd.Status.CurrentFilterIndex >= int32(len(sd.Spec.FilterPipeline)) {
		sd.Status.Status = FinishStatus
		if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			return r.Status().Update(ctx, sd)
		}); err != nil {
			l.Error(err, "unable to update SchedulingData status")
			return err
		}

	}
	return nil
}

func (r *FilterTemplate) fetchDataflow(ctx context.Context, sd *ntthpcv1.SchedulingData) (*ntthpcv1.DataFlow, error) {
	l := log.FromContext(ctx)
	var df *ntthpcv1.DataFlow = &ntthpcv1.DataFlow{}
	if err := r.Get(ctx, client.ObjectKey{Namespace: sd.ObjectMeta.Namespace, Name: sd.ObjectMeta.Name}, df); err != nil {
		l.Error(err, "unable to fetch dataflow")
		return df, err
	}
	return df, nil
}

func (r *FilterTemplate) getSchedulingData(ctx context.Context, req ctrl.Request) (*ntthpcv1.SchedulingData, error) {
	l := log.FromContext(ctx)
	var sd ntthpcv1.SchedulingData
	if err := r.Get(ctx, req.NamespacedName, &sd); err != nil {
		l.Error(err, "Filter: Reconciler unable to fetch SchedulingData")
		return nil, err
	}
	return &sd, nil
}

func (r *FilterTemplate) fetchUserParameters(
	ctx context.Context,
	df *ntthpcv1.DataFlow,
	filterIndex int32,
	filterName string,
) (
	strategy map[string]string,
	usrreq map[string]string,
	err error,
) {
	l := log.FromContext(ctx)

	usrreq, err = r.fetchUserRequirement(ctx, df)
	if err != nil {
		l.Error(err, "Failed to fetch user request parameters")
		switch err {
		case ErrFirstConfigMapNotExist, ErrSecondConfigMapNotExist:
			l.Info("default value is attached for UserRequest")
			usrreq = r.attachDefaultValueForUserRequest()
			err = nil
		default:
			return strategy, usrreq, err
		}
	}

	if name, ok := usrreq["strategy"]; ok {
		strategy, err = r.fetchStrategy(ctx, df, filterIndex, name, filterName)
		if err != nil {
			switch err {
			case ErrFirstConfigMapNotExist:
				l.Info("Failed to fetch first configmap")
				l.Info("default value is attached for strategy")
				strategy = r.attachDefaultValueForStrategy()
				err = nil
			case ErrSecondConfigMapNotExist:
				l.Info("Failed to fetch second configmap")
				l.Info("default value is attached for strategy")
				strategy = r.attachDefaultValueForStrategy()
				err = nil
			default:
				return strategy, usrreq, err
			}
		}
	} else {
		strategy = r.attachDefaultValueForStrategy()
	}

	return strategy, usrreq, err
}

func (r *FilterTemplate) attachDefaultValueForStrategy() map[string]string {
	ret := make(map[string]string)
	return ret
}

func (r *FilterTemplate) attachDefaultValueForUserRequest() map[string]string {
	return make(map[string]string)
}

// Refers to the referenced configMap (second ConfigMap) described in ConfigMap (first ConfiuMap) in two steps,
// inputs the parameters described in them, and returns it.
// priority : "filterNumber"-key > "filterName"-key > key > key (on second ConfigMap)
func (r *FilterTemplate) recursiveFetchConfigMap(
	ctx context.Context,
	configMapName string,
	filterName string,
	nameSpaceCands []string,
	refKey string,
	refIdx int32,
	refConfigMapNameSpaceCands []string,
) (map[string]string, error) {

	l := log.FromContext(ctx)
	var ret map[string]string

	Cm, err := tryFetchConfigMapFromSeveralNameSpaceCandidates(
		r, ctx, configMapName,
		nameSpaceCands)
	if err != nil {
		l.Error(err, fmt.Sprintf("Failed to fetch first configmap : %v", configMapName))
		return ret, ErrFirstConfigMapNotExist
	}

	params := make([]parseConfigMap, 0)

	const (
		nonPrefix = iota
		named
		indexed
		prefixKindNum
	)

	referenceConfigMapNames := make([]*string, prefixKindNum)
	for i, _ := range referenceConfigMapNames {
		referenceConfigMapNames[i] = nil
	}

	prefixChecker := prefixCheckerStruct{
		FilterName:         filterName,
		StringCurrentIndex: strconv.Itoa(int(refIdx)),
	}

	for k, v := range Cm.Data {

		p := parseConfigMap{
			Prefix: "",
		}

		if strings.Contains(k, ".") {
			splitedKey := strings.Split(k, ".")
			// .(dot) which is not filter-index nor filter-name, is not allowed in key.
			if len(splitedKey) != 2 {
				return nil, errInvalidConfigMapParamter
			}
			p.Prefix = splitedKey[0]
			p.Key = splitedKey[1]
			p.Value = v
		} else {
			p.Key = k
			p.Value = v
		}

		if p.Key != refKey {
			params = append(params, p)
		} else {

			if prefixChecker.isNonePrefix(&p) {
				referenceConfigMapNames[nonPrefix] = &p.Value
			} else if prefixChecker.isNamed(&p) {
				referenceConfigMapNames[named] = &p.Value
			} else if prefixChecker.isIndexed(&p) {
				referenceConfigMapNames[indexed] = &p.Value
			}

		}
	}

	ret, err = r.fetchReferenceConfigMap(
		ctx,
		nameSpaceCands,
		filterName,
		refKey,
		refIdx,
		refConfigMapNameSpaceCands,
		params,
		referenceConfigMapNames,
	)

	if err != nil {
		return nil, err
	}

	substituteConfigMapParameters(ret, params, prefixChecker.isNonePrefix)
	substituteConfigMapParameters(ret, params, prefixChecker.isNamed)
	substituteConfigMapParameters(ret, params, prefixChecker.isIndexed)

	return ret, nil
}

func substituteConfigMapParameters(m map[string]string, params []parseConfigMap, check func(*parseConfigMap) bool) {

	for _, p := range params {
		if check(&p) {
			m[p.Key] = p.Value
		}
	}

}

func (r *FilterTemplate) fetchReferenceConfigMap(
	ctx context.Context,
	nameSpaceCands []string,
	filterName string,
	refKey string,
	refIdx int32,
	refConfigMapNameSpaceCands []string,
	parameter []parseConfigMap,
	configMapNames []*string,
) (map[string]string, error) {

	ret := make(map[string]string)

	for _, nameAdr := range configMapNames {

		if nameAdr == nil {
			continue
		}

		name := *nameAdr

		refParams, err := r.recursiveFetchConfigMap(
			ctx,
			name,
			filterName,
			nameSpaceCands,
			refKey,
			refIdx,
			refConfigMapNameSpaceCands,
		)

		if err != nil {
			return nil, err
		}

		overwriteMap(ret, refParams)
	}

	return ret, nil
}

func (r *FilterTemplate) dispose(sd *ntthpcv1.SchedulingData) {

	idm := r.exclusiveRead(sd)
	if idm.IsDBConnected {
		idm.DBConnection.Close()
	}
	r.exclusiveDelete(sd)

}

func (r *FilterTemplate) generateKey(sd *ntthpcv1.SchedulingData) string {
	return sd.ObjectMeta.Namespace + sd.ObjectMeta.Name
}

func (r *FilterTemplate) exclusiveRead(sd *ntthpcv1.SchedulingData) *independentMembers {
	key := r.generateKey(sd)
	r.independentMtx.Lock()
	ret := r.independent[key]
	r.independentMtx.Unlock()
	return &ret
}

func (r *FilterTemplate) exclusiveWrite(sd *ntthpcv1.SchedulingData, val *independentMembers) {
	key := r.generateKey(sd)
	r.independentMtx.Lock()
	r.independent[key] = *val
	r.independentMtx.Unlock()
}

func (r *FilterTemplate) exclusiveDelete(sd *ntthpcv1.SchedulingData) {
	key := r.generateKey(sd)
	r.independentMtx.Lock()
	delete(r.independent, key)
	r.independentMtx.Unlock()
}

type RequestType int

const (
	excludeRequest RequestType = iota
	hopeRequest
)

func checkRequestType(requests []string) (RequestType, error) {
	for _, r := range requests {
		if r[0] != '-' {
			return hopeRequest, nil
		}
	}
	return excludeRequest, nil
}

func filterElements[T *ntthpcv1.FunctionTarget |
	*ntthpcv1.ConnectionTarget |
	*ntthpcv1.ConnectionType |
	*ntthpcv1.TopologyInfo |
	string](
	ctx context.Context,
	cands *[]T,
	reqs []string,
	compFunc func(T, string) bool,
	resourceName string,
	filterName string,
) {

	l := log.FromContext(ctx)

	filteredCands := make([]T, 0)
	reqType, _ := checkRequestType(reqs)
	found := false

	switch reqType {
	case hopeRequest:
		for _, req := range reqs {
			if req[0] == '-' {
				continue
			}
			for _, cand := range *cands {
				if compFunc(cand, req) {
					filteredCands = append(filteredCands, cand)
					found = true
				}
			}
		}
	case excludeRequest:
		rmIdx := make(map[int]bool)

		for _, req := range reqs {
			req = req[1:] // remove "-"
			for i, cand := range *cands {
				if compFunc(cand, req) {
					rmIdx[i] = true
				}
			}
		}

		for i, cand := range *cands {
			if _, ok := rmIdx[i]; !ok {
				filteredCands = append(filteredCands, cand)
				found = true
			}
		}
	}

	if !found {
		l.Info("no " + resourceName + " found that passed " + filterName + "Filter: " + fmt.Sprintf("%v", reqs))
	}

	*cands = filteredCands

}

func (f *functionTargetFilter) FilterFunctionTarget(ctx context.Context, cands *[]*ntthpcv1.FunctionTarget) {
	checkFunc := func(ft *ntthpcv1.FunctionTarget, c string) bool {
		return ft.Name == c
	}
	filterElements(ctx, cands, *f.Filter.FunctionTargets, checkFunc, "FunctionTarget", "FunctionTarget")
}

func (f *functionTargetFilter) FilterRegionTypes(ctx context.Context, cands *[]*ntthpcv1.FunctionTarget) {
	checkFunc := func(ft *ntthpcv1.FunctionTarget, c string) bool {
		return ft.Status.RegionType == c
	}
	filterElements(ctx, cands, *f.Filter.RegionTypes, checkFunc, "FunctionTarget", "RegionType")
}

func (f *functionTargetFilter) FilterNodeNames(ctx context.Context, cands *[]*ntthpcv1.FunctionTarget) {
	checkFunc := func(ft *ntthpcv1.FunctionTarget, c string) bool {
		return ft.Status.NodeName == c
	}
	filterElements(ctx, cands, *f.Filter.NodeNames, checkFunc, "FunctionTarget", "NodeName")
}

func (f *functionTargetFilter) FilterDeviceTypes(ctx context.Context, cands *[]*ntthpcv1.FunctionTarget) {
	checkFunc := func(ft *ntthpcv1.FunctionTarget, c string) bool {
		return ft.Status.DeviceType == c
	}
	filterElements(ctx, cands, *f.Filter.DeviceTypes, checkFunc, "FunctionTarget", "DeviceType")
}

func (f *functionTargetFilter) FilterAvailableFunctionTargets(cands *[]*ntthpcv1.FunctionTarget) {

	retCands := make([]*ntthpcv1.FunctionTarget, 0)

	for _, cand := range *cands {
		if cand.Status.Available {
			retCands = append(retCands, cand)
		}
	}

	*cands = retCands
}

func (f *functionTargetFilter) FilterAvailableFunctionIndexStructs(cands *[]FunctionIndexStruct) {

	retCands := make([]FunctionIndexStruct, 0)

	for _, fis := range *cands {
		if fis.FunctionIndex == nil {
			retCands = append(retCands, fis)
		} else {
			for _, function := range fis.FunctionTarget.Status.Functions {
				if *fis.FunctionIndex == function.FunctionIndex && function.Available {
					retCands = append(retCands, fis)
				}
			}
		}
	}

	*cands = retCands
}

func (f *functionTargetFilter) FilterRegionNames(ctx context.Context, cands *[]*ntthpcv1.FunctionTarget) {
	checkFunc := func(ft *ntthpcv1.FunctionTarget, c string) bool {
		return ft.Status.RegionName == c
	}
	filterElements(ctx, cands, *f.Filter.RegionNames, checkFunc, "FunctionTarget", "RegionName")
}

func (f *functionTargetFilter) FilterFunctionNames(ctx context.Context, cands *[]*ntthpcv1.FunctionTarget) {

	l := log.FromContext(ctx)

	funcNames := *f.Filter.FunctionNames
	retCands := make([]*ntthpcv1.FunctionTarget, 0)

	found := false
	for _, cand := range *cands {
		addFlag := false
		for _, f := range cand.Status.Functions {
			for _, fName := range funcNames {
				if fName == f.FunctionName {
					retCands = append(retCands, cand)
					addFlag = true
					found = true
					break
				}
			}
			if addFlag {
				break
			}
		}
	}

	if !found {
		l.Info("no function found that passed FunctionNameFilter: " + fmt.Sprintf("%v", funcNames))
	}

	*cands = retCands
}

func (f *functionTargetFilter) FilterFunctionIndexes(ctx context.Context, cands *[]FunctionIndexStruct) {

	l := log.FromContext(ctx)

	reqs := *f.Filter.FunctionIndexes
	filteredCands := make([]FunctionIndexStruct, 0)
	reqType, _ := checkRequestType(reqs)
	found := false

	switch reqType {
	case hopeRequest:
		for _, req := range reqs {
			if req[0] == '-' {
				continue
			}
			targetFunctionIndex, _ := strconv.Atoi(req)
			for _, cand := range *cands {
				if cand.FunctionIndex != nil && int32(targetFunctionIndex) == *cand.FunctionIndex {
					filteredCands = append(filteredCands, cand)
					found = true
				}
			}
		}
	case excludeRequest:
		rmIdx := make(map[int]bool)
		for _, req := range reqs {
			req = req[1:] // remove "-"
			excludeTargetFunctionIndex, _ := strconv.Atoi(req)
			for i, cand := range *cands {
				if cand.FunctionIndex != nil && int32(excludeTargetFunctionIndex) == *cand.FunctionIndex {
					rmIdx[i] = true
				}
			}
		}
		for i, cand := range *cands {
			if _, ok := rmIdx[i]; !ok {
				filteredCands = append(filteredCands, cand)
				found = true
			}
		}
	}

	if !found {
		l.Info("no function found that passed FunctionIndexFilter: " + fmt.Sprintf("%v", reqs))
	}

	*cands = filteredCands
}

func (f *connectionTargetFilter) FilterNodeNames(ctx context.Context, cands *[]*ntthpcv1.ConnectionTarget) {
	checkFunc := func(ct *ntthpcv1.ConnectionTarget, c string) bool {
		return ct.Status.NodeName == c
	}
	filterElements(ctx, cands, *f.Filter.NodeNames, checkFunc, "ConnectionTarget", "NodeName")
}

func (f *connectionTargetFilter) FilterDeviceTypes(ctx context.Context, cands *[]*ntthpcv1.ConnectionTarget) {
	checkFunc := func(ct *ntthpcv1.ConnectionTarget, c string) bool {
		return ct.Status.DeviceType == c
	}
	filterElements(ctx, cands, *f.Filter.DeviceTypes, checkFunc, "ConnectionTarget", "DeviceType")
}

func (f *connectionTargetFilter) FilterConnectionTarget(ctx context.Context, cands *[]*ntthpcv1.ConnectionTarget) {
	checkFunc := func(ct *ntthpcv1.ConnectionTarget, c string) bool {
		return ct.Name == c
	}
	filterElements(ctx, cands, *f.Filter.ConnectionTarget, checkFunc, "ConnectionTarget", "ConnectionTarget")
}

func (f *deviceTypeFilter) FilterDeviceTypes(ctx context.Context, cands *[]string) {
	checkFunc := func(v string, c string) bool {
		return v == c
	}
	filterElements(ctx, cands, *f.Filter.DeviceTypes, checkFunc, "DeviceType", "DeviceType")
}

func (f *connectionTypeFilter) FilterConnectionTypes(ctx context.Context, cands *[]*ntthpcv1.ConnectionType) {
	checkFunc := func(v *ntthpcv1.ConnectionType, c string) bool {
		return v.Name == c
	}
	filterElements(ctx, cands, *f.Filter.ConnectionType, checkFunc, "ConnectionType", "ConnectionType")
}

func (r *prefixCheckerStruct) isNonePrefix(p *parseConfigMap) bool {
	return p.Prefix == ""
}

func (r *prefixCheckerStruct) isNamed(p *parseConfigMap) bool {
	return p.Prefix == r.FilterName
}

func (r *prefixCheckerStruct) isIndexed(p *parseConfigMap) bool {
	return p.Prefix == r.StringCurrentIndex
}

func valToAddr[T any](in T) *T {
	return &in
}

func (r *FilterTemplate) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ntthpcv1.SchedulingData{}).
		Complete(r)
}

func (r *FilterTemplate) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

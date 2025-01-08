package combination_filters

import (
	"context"
	"fmt"
	"strings"

	v1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	"github.com/compsysg/whitebox-k8s-flowctrl/lib/filter_template"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Controller Name
const controllerName = "combinationfilter"

// Logger Setting
const (
	loggerKeyController      = "combinationfilter"
	loggerKeyControllerGroup = "example.com"
	loggerKeyControllerKind  = "CombinationFilter"
)

const (
	// The connection method set in Type of ConnectionScheduleInfo is host-mem
	connMethodHostMem string = "host-mem"
	// The connection method set in Type of ConnectionScheduleInfo is host-100gether
	connMethodHost100g string = "host-100gether"
)

const (
	// Default value of TopologyInfo's name
	defaultTiName = "topologyinfo"
	// Default value of TopologyInfo namespace
	defaultTiNameSpace = "topologyinfo"
	// Key name when reading TopologyInfo name from userRequirement
	topologyInfoNameKey = "topologyInfoName"
	// Key name when reading TopologyInfo namespace from userRequirement
	topologyInfoNameSpaceKey = "topologyInfoNameSpace"
)

type CombinationFilters struct {
	filter_template.FilterTemplate
}

const (
	generateCombinationsFilter  = "GenerateCombinations"
	generateRouteFilter         = "GenerateRoute"
	targetResourceFitFilter     = "TargetResourceFit"
	connectionResourceFitFilter = "ConnectionResourceFit"
)

//+kubebuilder:rbac:groups=example.com,resources=schedulingdata,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=schedulingdata/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=schedulingdata/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=dataflows,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=dataflows/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=dataflows/finalizers,verbs=update
//+kubebuilder:rbac:groups=example.com,resources=functiontargets,verbs=get;list;watch
//+kubebuilder:rbac:groups=example.com,resources=functiontargets/status,verbs=get
//+kubebuilder:rbac:groups=example.com,resources=topologyinfos,verbs=get;list;watch
//+kubebuilder:rbac:groups=example.com,resources=topologyinfos/status,verbs=get
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list

func (r *CombinationFilters) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	filters := []string{
		generateCombinationsFilter,
		targetResourceFitFilter,
		connectionResourceFitFilter,
		generateRouteFilter,
	}

	sd, df, err := r.Initialize(ctx, req, filters...)
	if err != nil {
		return ctrl.Result{}, nil
	}

	switch r.GetName(sd) {
	case generateCombinationsFilter:
		return r.generateCombinationsFilter(ctx, req, sd, df)
	case targetResourceFitFilter:
		return r.targetResourceFitFilter(ctx, req, sd, df)
	case connectionResourceFitFilter:
		return r.connectionResourceFitFilter(ctx, req, sd, df)
	case generateRouteFilter:
		return r.generateRouteFilter(ctx, req, sd, df)
	}

	return ctrl.Result{}, nil
}

func (r *CombinationFilters) SetupWithManager(mgr ctrl.Manager) error {

	p := predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return true
		},
	}

	return ctrl.NewControllerManagedBy(mgr).
		Named(controllerName).
		For(&v1.SchedulingData{}).
		WithLogConstructor(func(req *reconcile.Request) logr.Logger {
			return mgr.GetLogger().WithValues("controller", loggerKeyController, "controllerGroup",
				loggerKeyControllerGroup, "controllerKind", loggerKeyControllerKind)
		}).
		WithEventFilter(p).
		Complete(r)
}

// Get TopologyInfo
func (r *CombinationFilters) getTopologyInfo(ctx context.Context, sd *v1.SchedulingData) (*v1.TopologyInfo, error) {

	l := log.FromContext(ctx)

	// Get the default namespace value from a constant
	tiName := defaultTiName
	tiNameSpace := defaultTiNameSpace

	// Check whether TopologyInfo name is specified and its value in userRequirement
	if r.DoesUserRequirementHaveKey(sd, topologyInfoNameKey) {
		r.LoadUserRequirementParameter(sd, topologyInfoNameKey, &tiName)
	}

	// Check whether the namespace of TopologyInfo is specified in userRequirement and its value
	if r.DoesUserRequirementHaveKey(sd, topologyInfoNameSpaceKey) {
		r.LoadUserRequirementParameter(sd, topologyInfoNameSpaceKey, &tiNameSpace)
	}

	var ti v1.TopologyInfo
	if err := r.Get(ctx, types.NamespacedName{Name: tiName, Namespace: tiNameSpace}, &ti); err != nil {
		l.Error(err, "unable to fetch TopologyInfo")
		return nil, err
	}

	return &ti, nil
}

// Eliminate combinations that do not satisfy the conditions
// If only the data with index number [0] is to be excluded, the processing is not performed correctly, so the processing content is corrected on the FJ side.
func exclude[T any](combsAdr *[]T, indexes []int32) {
	combs := *combsAdr

	if len(indexes) == 0 {
		return
	}

	ret := make([]T, 0, len(combs)-len(indexes))
	lastIndex := len(combs) - 1

	// Add the part before the first index
	if indexes[0] > 0 {
		ret = append(ret, combs[:indexes[0]]...)
	}

	// Add the part between the middle indices
	for i := 1; i < len(indexes); i++ {
		from := indexes[i-1] + 1
		to := indexes[i]
		if from <= to && to <= int32(lastIndex) {
			ret = append(ret, combs[from:to]...)
		}
	}

	// Add the part after the last index
	if indexes[len(indexes)-1]+1 <= int32(lastIndex) {
		ret = append(ret, combs[indexes[len(indexes)-1]+1:]...)
	}

	*combsAdr = ret
}

// Get the requested capacity
func (r *CombinationFilters) getRequireCapacity(df *v1.DataFlow) int32 {
	// Returns df.Spec.Requirements.All.Capacity
	// (It seems there is a field that allows you to set the capacity for each function, but I haven't used the FJ version so I'll use this for now)
	return df.Spec.Requirements.All.Capacity
}

// Get the FunctionTargetName from FunctionScheduleInfo.
func getFunctionTargetName(fsi v1.FunctionScheduleInfo) string {
	return fmt.Sprintf("%s.%s-%d.%s", fsi.NodeName, fsi.DeviceType, fsi.DeviceIndex, fsi.RegionName)
}

// get strings indicating route information from ConnectionPath
func (r *CombinationFilters) getRouteStr(wbcRoute []v1.WBConnectionPath) string {
	routeStrArry := []string{}
	for _, route := range wbcRoute {
		routeStrArry = append(routeStrArry, route.EntityID)
	}
	return strings.Join(routeStrArry, " -> ")
}

package score_filters

import (
	"context"
	"fmt"
	"sort"

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
const controllerName = "scorefilter"

// Logger Setting
const (
	loggerKeyController      = "scorefilter"
	loggerKeyControllerGroup = "example.com"
	loggerKeyControllerKind  = "ScoreFilter"
)

const (
	targetResourceFitScore = "TargetResourceFitScore"
	routeScore             = "RouteScore"
)

const (
	selectTopKey = "selectTop"
)

//nolint:unused // FIXME: remove this type
type requirement struct {
	Capacity int
	Resource int
}

type ScoreFilters struct {
	filter_template.FilterTemplate
}

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

func (r *ScoreFilters) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	l.Info("start Scheduler " + req.NamespacedName.Namespace)

	filters := []string{
		targetResourceFitScore,
		routeScore,
	}

	sd, df, err := r.Initialize(ctx, req, filters...)
	if err != nil {
		return ctrl.Result{}, nil
	}

	err = nil
	switch r.GetName(sd) {
	case targetResourceFitScore:
		err = r.targetResourceFitScore(ctx, req, sd, df)
	case routeScore:
		err = r.routeScore(ctx, req, sd, df)
	}

	if err != nil {
		r.Abort(ctx, sd, err)
	}

	r.selectTop(sd)

	return r.Finalize(ctx, sd)
}

func (r *ScoreFilters) SetupWithManager(mgr ctrl.Manager) error {

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

func (r *ScoreFilters) selectTop(sd *v1.SchedulingData) {

	// Check if there is a setting for top n
	if r.DoesStrategyHaveKey(sd, selectTopKey) {
		var selectTopNum int
		r.LoadStrategyParameter(sd, selectTopKey, &selectTopNum)

		sort.SliceStable(sd.Status.TargetCombinations,
			func(i, j int) bool {
				return *sd.Status.TargetCombinations[i].Score > *sd.Status.TargetCombinations[j].Score
			})

		if len(sd.Status.TargetCombinations) > selectTopNum {
			sd.Status.TargetCombinations = sd.Status.TargetCombinations[:selectTopNum]
		}

	}

}

//nolint:unused // FIXME: remove this function
func (r *ScoreFilters) getTargetCombinations(sd *v1.SchedulingData) []v1.TargetCombinationStruct {
	return sd.Status.TargetCombinations
}

// Get the FunctionTargetName from FunctionScheduleInfo.
func getFunctionTargetName(scheduledFunctions v1.FunctionScheduleInfo) string {
	return fmt.Sprintf("%s.%s-%d.%s", scheduledFunctions.NodeName, scheduledFunctions.DeviceType,
		scheduledFunctions.DeviceIndex, scheduledFunctions.RegionName)
}

// Get the requested capacity
func (r *ScoreFilters) getRequireCapacity(df *v1.DataFlow) int32 {
	// Returns df.Spec.Requirements.All.Capacity
	// (It seems there is a field that allows you to set the capacity for each function, but I haven't used the FJ version so I'll use this for now)
	return df.Spec.Requirements.All.Capacity
}

// Get TopologyInfo
func (r *ScoreFilters) getTopologyInfo(ctx context.Context, sd *v1.SchedulingData) (*v1.TopologyInfo, error) {

	// Read from TopologyInfo default values
	// If you specify it separately, read it from the strategy.
	tiName := "topologyinfo"
	tiNameSpace := "topologyinfo"

	var ti v1.TopologyInfo
	if err := r.Get(ctx, types.NamespacedName{Name: tiName, Namespace: tiNameSpace}, &ti); err != nil {
		return nil, err
	}

	return &ti, nil
}

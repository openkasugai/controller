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

package main

import (
	"flag"
	"os"
	"strconv"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.

	zaplog "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	examplecomv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	"github.com/compsysg/whitebox-k8s-flowctrl/internal/controller"

	wbschedulercontroller "github.com/compsysg/whitebox-k8s-flowctrl/internal/controller/wbscheduler_controller"

	"github.com/compsysg/whitebox-k8s-flowctrl/internal/controller/combination_filters"
	"github.com/compsysg/whitebox-k8s-flowctrl/internal/controller/score_filters"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(examplecomv1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var waitTimeSec int
	var requeueTimeSec int
	var enableConnectionSupport bool
	var topologyinfoName string
	var topologyinfoNamespace string
	var topologydataName string
	var topologydataNamespace string
	var defaultFilterPipeline string
	var logLevel int
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.IntVar(&waitTimeSec, "wait-time-sec", 30, "Wait time after scheduling [seconds]")
	flag.IntVar(&requeueTimeSec, "requeue-time-sec", 10, "Requeue interval [seconds]")
	flag.StringVar(&topologyinfoName, "topologyinfo-name", "topologyinfo", "topologyinfo name")
	flag.StringVar(&topologyinfoNamespace, "topologyinfo-namespace", "topologyinfo", "topologyinfo namespace")
	flag.StringVar(&topologydataName, "topologydata-name", "topologydata", "topologydata name")
	flag.StringVar(&topologydataNamespace, "topologydata-namespace", "topologyinfo", "topologydata namespace")
	flag.StringVar(&defaultFilterPipeline, "default-filter-pipeline",
		"GenerateCombinations,TargetResourceFit,TargetResourceFitScore",
		"default filter pipeline")
	flag.IntVar(&logLevel, "log-level", -1, "zapcore log level. default=DEBUG(-1)")

	opts := zap.Options{
		TimeEncoder: zapcore.ISO8601TimeEncoder,
		Development: true,
		ZapOpts:     []zaplog.Option{zaplog.AddCaller()},
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts), zap.Level(zapcore.Level(logLevel))))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metricsserver.Options{BindAddress: metricsAddr},
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "b32462f3.example.com",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	setupLog.Info("enableConnectionSupport=" + strconv.FormatBool(enableConnectionSupport))

	if err = (&controller.DataFlowReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("dataflow-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "DataFlow")
		os.Exit(1)
	}
	if err = (&controller.FunctionChainReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("functionchain-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "FunctionChain")
		os.Exit(1)
	}
	if err = (&controller.FunctionTypeReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "FunctionType")
		os.Exit(1)
	}

	if err = (&wbschedulercontroller.WBschedulerReconciler{
		Client:                mgr.GetClient(),
		Scheme:                mgr.GetScheme(),
		Recorder:              mgr.GetEventRecorderFor("wbscheduler-controller"),
		DefaultFilterPipeline: defaultFilterPipeline,
		WaitTimeSec:           waitTimeSec,
		RequeueTimeSec:        requeueTimeSec,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "WBscheduler")
		os.Exit(1)
	}

	var combFilter combination_filters.CombinationFilters
	combFilter.Client = mgr.GetClient()
	combFilter.Scheme = mgr.GetScheme()
	combFilter.Recorder = mgr.GetEventRecorderFor("combinationfilter-controller")
	if err = combFilter.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CombinationFilter")
		os.Exit(1)
	}

	var scoreFilter score_filters.ScoreFilters
	scoreFilter.Client = mgr.GetClient()
	scoreFilter.Scheme = mgr.GetScheme()
	scoreFilter.Recorder = mgr.GetEventRecorderFor("scorefilter-controller")
	if err = scoreFilter.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ScoreFilter")
		os.Exit(1)
	}

	if err = (&controller.FunctionTargetReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "FunctionTarget")
		os.Exit(1)
	}
	if err = (&controller.TopologyInfoReconciler{
		Client:                     mgr.GetClient(),
		Scheme:                     mgr.GetScheme(),
		TopologyinfoNamespacedName: types.NamespacedName{Namespace: topologyinfoNamespace, Name: topologyinfoName},
		TopologydataNamespacedName: types.NamespacedName{Namespace: topologydataNamespace, Name: topologydataName},
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "TopologyInfo")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

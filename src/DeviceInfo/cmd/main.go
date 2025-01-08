/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED
*/

package main

import (
	"flag"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	examplecomv1 "DeviceInfo/api/v1"
	"DeviceInfo/internal/controller"
	//+kubebuilder:scaffold:imports
	/* Additional files */
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	coveragelog "log"
	"time"
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

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	/* metric server port */
	ms_port := os.Getenv("K8S_DI_MS_PORT")
	if ms_port != "" {
		metricsAddr = ":" + ms_port
	}
	coveragelog.Println("ms_port =", ms_port)

	/* health check server port */
	hc_port := os.Getenv("K8S_DI_HC_PORT")
	if hc_port != "" {
		probeAddr = ":" + hc_port
	}
	coveragelog.Println("hc_port =", hc_port)
	coveragelog.Println("metricsAddr =", metricsAddr)
	coveragelog.Println("probeAddr =", probeAddr)

	ctrl.SetLogger(zapr.NewLogger(SettingZapLogger()))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metricsserver.Options{BindAddress: metricsAddr},
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "56592b87.example.com",
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

	if err = (&controller.DeviceInfoReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("deviceinfo-ctrl"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "DeviceInfo")
		os.Exit(1)
	}
	/*
		if err = (&controller.FPGAReconciler{
			Client:   mgr.GetClient(),
			Scheme:   mgr.GetScheme(),
			Recorder: mgr.GetEventRecorderFor("fpga-ctrl"),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "FPGA")
			os.Exit(1)
		}
	*/
	//+kubebuilder:scaffold:builder
	if err := controller.StartupProccessing(&controller.DeviceInfoReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("deviceinfo-ctrl"),
	}, mgr); err != nil {
		setupLog.Error(err, "controllers.StartupProccessing() error")
		os.Exit(1)
	}
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

// Change the date and time display of the log
func JSTTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(time.RFC3339Nano))
}

// Log output settings
func SettingZapLogger() *zap.Logger {
	zc := zap.NewDevelopmentConfig()
	zc.OutputPaths = []string{"stdout"}
	zc.ErrorOutputPaths = zc.OutputPaths
	zc.DisableStacktrace = true
	zc.DisableCaller = false
	zc.EncoderConfig.EncodeTime = JSTTimeEncoder
	z, _ := zc.Build()
	return z
}

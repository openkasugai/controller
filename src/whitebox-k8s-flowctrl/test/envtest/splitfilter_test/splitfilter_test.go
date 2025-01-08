/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED
*/

package splitfilter_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	v1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme" //nolint:stylecheck // ST1019: intentional import as another name
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	. "github.com/compsysg/whitebox-k8s-flowctrl/internal/controller/combination_filters"
	. "github.com/compsysg/whitebox-k8s-flowctrl/internal/controller/score_filters"
	. "github.com/compsysg/whitebox-k8s-flowctrl/internal/controller/wbscheduler_controller"
	. "github.com/compsysg/whitebox-k8s-flowctrl/test/lib/testutils"
	test_cases "github.com/compsysg/whitebox-k8s-flowctrl/test/test_cases/splitfilter"

	clientgoscheme "k8s.io/client-go/kubernetes/scheme" //nolint:stylecheck // ST1019: intentional import as another name
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var testScheme *runtime.Scheme

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	SetDefaultEventuallyTimeout(10 * time.Second)
	SetDefaultEventuallyPollingInterval(100 * time.Millisecond)
	SetDefaultConsistentlyDuration(10 * time.Second)
	SetDefaultConsistentlyPollingInterval(100 * time.Millisecond)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	testScheme = runtime.NewScheme()

	err = v1.AddToScheme(testScheme)
	Expect(err).NotTo(HaveOccurred())

	err = clientgoscheme.AddToScheme(testScheme)
	Expect(err).NotTo(HaveOccurred())

	err = v1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: testScheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	ns := &corev1.Namespace{}
	ns.Name = "topologyinfo"
	_ = k8sClient.Create(context.Background(), ns) // FIXME: need to handle return value (error)

})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("SplitFilter", func() {

	ctx := context.Background()
	var (
		mgr      manager.Manager
		stopFunc func()
	)

	BeforeEach(func() {

		Expect(DeleteAllOf(ctx, k8sClient, "default",
			&v1.SchedulingData{},
			&v1.DataFlow{},
			&v1.FunctionTarget{},
			&corev1.ConfigMap{},
		)).Should(Succeed())

		Expect(DeleteAllOf(ctx, k8sClient, "topologyinfo",
			&v1.TopologyInfo{},
		)).Should(Succeed())

		Eventually(func(g Gomega) {
			var err error
			mgr, err = newMgr()
			Expect(err).NotTo(HaveOccurred())
		}).Should(Succeed())

		var combFilters CombinationFilters
		combFilters.Client = k8sClient
		combFilters.Scheme = testScheme
		combFilters.Recorder = mgr.GetEventRecorderFor("combinationfilter-controller")
		Expect(combFilters.SetupWithManager(mgr)).Should(Succeed())

		var scoreFilters ScoreFilters
		scoreFilters.Client = k8sClient
		scoreFilters.Scheme = testScheme
		scoreFilters.Recorder = mgr.GetEventRecorderFor("scorefilter-controller")
		Expect(scoreFilters.SetupWithManager(mgr)).Should(Succeed())

	})

	AfterEach(func() {
		Expect(DeleteAllOf(ctx, k8sClient, "default",
			&v1.SchedulingData{},
			&v1.DataFlow{},
			&v1.FunctionTarget{},
			&corev1.ConfigMap{},
		)).Should(Succeed())

		Expect(DeleteAllOf(ctx, k8sClient, "topologyinfo",
			&v1.TopologyInfo{},
		)).Should(Succeed())

		stopFunc()
		time.Sleep(100 * time.Millisecond)
	})

	// Separate the Context depending on whether or not you use the Scheduler.
	Context("WithScheduler", func() {

		BeforeEach(func() {
			// Deploy scheduler
			scheduler := WBschedulerReconciler{
				Client:                k8sClient,
				Scheme:                testScheme,
				Recorder:              mgr.GetEventRecorderFor("wbscheduler-controller"),
				RequeueTimeSec:        1,
				DefaultFilterPipeline: "GenerateCombinations,TargetResourceFit,TargetResourceFitScore",
			}
			Expect(scheduler.SetupWithManager(mgr)).Should(Succeed())

			stopFunc = startMgr(ctx, mgr)
		})

		// Normal TopologyInfo not found
		It("1_Nominal", func() {
			test_cases.Test1_Nominal(ctx, k8sClient)
		})

		// Normal TopologyInfo available
		It("2_Connection", func() {
			test_cases.Test2_Connection(ctx, k8sClient)
		})

		// Normal TopologyInfo not found, various DataFlow
		It("3_1_Nominal_DataFlows_FPGADecode", func() {
			test_cases.Test4_Nominal_DataFlows(ctx, k8sClient, "fpga_decode")
		})

		It("3_2_Nominal_DataFlows_CPUDecode", func() {
			test_cases.Test4_Nominal_DataFlows(ctx, k8sClient, "cpu_decode")
		})

		It("3_3_Nominal_DataFlows_CopyBranch", func() {
			test_cases.Test4_Nominal_DataFlows(ctx, k8sClient, "copy_branch")
		})

		It("3_4_Nominal_DataFlows_GlueFDMAtoTCP", func() {
			test_cases.Test4_Nominal_DataFlows(ctx, k8sClient, "glue_fdma_to_tcp")
		})

		// Normal TopologyInfo available, various DataFlows
		It("4_1_Connection_DataFlows_FPGADecode", func() {
			test_cases.Test5_Connection_DataFlows(ctx, k8sClient, "fpga_decode")
		})

		It("4_2_Connection_DataFlows_CPUDecode", func() {
			test_cases.Test5_Connection_DataFlows(ctx, k8sClient, "cpu_decode")
		})

		It("4_3_Connection_DataFlows_CopyBranch", func() {
			test_cases.Test5_Connection_DataFlows(ctx, k8sClient, "copy_branch")
		})

		It("4_4_Connection_DataFlows_GlueFDMAtoTCP", func() {
			test_cases.Test5_Connection_DataFlows(ctx, k8sClient, "glue_fdma_to_tcp")
		})

		// Dynamic reconfiguration
		It("5_1_Dynamic_Reconfiguraion_NoTopology", func() {
			test_cases.Test6_Dynamic_Reconfiguration(ctx, k8sClient, "no_topology")
		})

		It("5_2_Dynamic_Reconfiguraion_WithTopology", func() {
			test_cases.Test6_Dynamic_Reconfiguration(ctx, k8sClient, "with_topology")
		})

		// Scheduling retry
		It("6_Schduling_Retry_UpdateFt", func() {
			test_cases.Test7_Scheduling_Retry(ctx, k8sClient, "update_ft")
		})
	})

	Context("Unit", func() {
		BeforeEach(func() {
			stopFunc = startMgr(ctx, mgr)
		})

		// Normal TopologyInfo not found
		It("1_Nominal", func() {
			test_cases.UnitTest1_Nominal(ctx, k8sClient)
		})

		// Normal TopologyInfo available
		It("2_Connection", func() {
			test_cases.UnitTest2_Connection(ctx, k8sClient)
		})

		// Normal TopologyInfo not found, various DataFlow
		It("3_1_Nominal_DataFlows_FPGADecode", func() {
			test_cases.UnitTest3_Nominal_DataFlows(ctx, k8sClient, "fpga_decode")
		})

		It("3_2_Nominal_DataFlows_CPUDecode", func() {
			test_cases.UnitTest3_Nominal_DataFlows(ctx, k8sClient, "cpu_decode")
		})

		It("3_3_Nominal_DataFlows_CopyBranch", func() {
			test_cases.UnitTest3_Nominal_DataFlows(ctx, k8sClient, "copy_branch")
		})

		It("3_4_Nominal_DataFlows_GlueFDMAtoTCP", func() {
			test_cases.UnitTest3_Nominal_DataFlows(ctx, k8sClient, "glue_fdma_to_tcp")
		})

		// Normal TopologyInfo available, various DataFlows
		It("4_1_Connection_DataFlows_FPGADecode", func() {
			test_cases.UnitTest4_Connection_DataFlows(ctx, k8sClient, "fpga_decode")
		})

		// Normal TopologyInfo available, various DataFlows
		It("4_2_Connection_DataFlows_CPUDecode", func() {
			test_cases.UnitTest4_Connection_DataFlows(ctx, k8sClient, "cpu_decode")
		})

		// Normal TopologyInfo available, various DataFlows
		It("4_3_Connection_DataFlows_CopyBranch", func() {
			test_cases.UnitTest4_Connection_DataFlows(ctx, k8sClient, "copy_branch")
		})

		// Normal TopologyInfo available, various DataFlows
		It("4_4_Connection_DataFlows_GlueFDMAtoTCP", func() {
			test_cases.UnitTest4_Connection_DataFlows(ctx, k8sClient, "glue_fdma_to_tcp")
		})

		// Dynamic reconfiguration
		It("5_1_Dynamic_Reconfiguraion_NoTopology", func() {
			test_cases.UnitTest5_Dynamic_Reconfiguration(ctx, k8sClient, "no_topology")
		})

		It("5_2_Dynamic_Reconfiguraion_WithTopology", func() {
			test_cases.UnitTest5_Dynamic_Reconfiguration(ctx, k8sClient, "with_topology")
		})
	})
})

func newMgr() (ctrl.Manager, error) {
	return ctrl.NewManager(cfg,
		ctrl.Options{
			Scheme:         testScheme,
			LeaderElection: false,
			Metrics: metricsserver.Options{
				BindAddress: ":50001",
			},
		},
	)
}

func startMgr(ctx context.Context, mgr manager.Manager) func() {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		err := mgr.Start(ctx)
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
	return cancel
}

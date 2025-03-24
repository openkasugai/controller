package combination_filters_test

import (
	"bytes"
	"context"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	examplecomv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1" //nolint:stylecheck // ST1019: intentional import as another name
	ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"     //nolint:stylecheck // ST1019: intentional import as another name

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1" //nolint:stylecheck // ST1019: intentional import as another name
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest" //nolint:stylecheck // ST1019: intentional import as another name
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	. "github.com/compsysg/whitebox-k8s-flowctrl/internal/controller/combination_filters"
	. "github.com/compsysg/whitebox-k8s-flowctrl/lib/filter_template"
	logf "sigs.k8s.io/controller-runtime/pkg/log" //nolint:stylecheck // ST1019: intentional import as another name

	. "github.com/compsysg/whitebox-k8s-flowctrl/test/lib/testutils"
	"github.com/go-logr/logr"
	"k8s.io/client-go/kubernetes/scheme"                //nolint:stylecheck // ST1019: intentional import as another name
	clientgoscheme "k8s.io/client-go/kubernetes/scheme" //nolint:stylecheck // ST1019: intentional import as another name
	//+kubebuilder:scaffold:import
)

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var testScheme *runtime.Scheme
var addr int = 13000
var buf1, buf2, buf3, buf4 *bytes.Buffer
var logger1, logger2, logger3, logger4 logr.Logger

func getMgr(mgr ctrl.Manager) (ctrl.Manager, error) {
	if mgr == nil {
		return ctrl.NewManager(cfg, ctrl.Options{
			Scheme: testScheme,
			Metrics: metricsserver.Options{
				BindAddress: strconv.Itoa(addr),
			},
		})
	}
	return mgr, nil
}

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	SetDefaultEventuallyTimeout(10 * time.Second)
	SetDefaultEventuallyPollingInterval(100 * time.Millisecond)
	SetDefaultConsistentlyDuration(10 * time.Second)
	SetDefaultConsistentlyPollingInterval(100 * time.Millisecond)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
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

	Expect(ntthpcv1.AddToScheme(scheme.Scheme)).Should(Succeed())
	Expect(examplecomv1.AddToScheme(scheme.Scheme)).Should(Succeed())

	Expect(clientgoscheme.AddToScheme(testScheme)).Should(Succeed())
	Expect(ntthpcv1.AddToScheme(testScheme)).Should(Succeed())
	Expect(examplecomv1.AddToScheme(testScheme)).Should(Succeed())

	//+kubebuilder:scaffold:scheme
	k8sClient, err = client.New(cfg, client.Options{Scheme: testScheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	buf1 = &bytes.Buffer{}
	logger1 = zap.New(zap.WriteTo(buf1), zap.UseDevMode(true))

	buf2 = &bytes.Buffer{}
	logger2 = zap.New(zap.WriteTo(buf2), zap.UseDevMode(true))

	buf3 = &bytes.Buffer{}
	logger3 = zap.New(zap.WriteTo(buf3), zap.UseDevMode(true))

	buf4 = &bytes.Buffer{}
	logger4 = zap.New(zap.WriteTo(buf4), zap.UseDevMode(true))

	ns := &corev1.Namespace{}
	ns.Name = "topologyinfo"
	_ = k8sClient.Create(context.Background(), ns) // FIXME: need to handle return value (error)
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("CombinationFilter", func() {

	Context("Test for generateCombinationFilter scheduling condition", func() {
		var mgr ctrl.Manager
		var combFilter CombinationFilters
		var ctx1 context.Context

		BeforeEach(func() {
			ctx1 = logf.IntoContext(context.Background(), logger1)
			commonBeforeEach(ctx1, mgr, &combFilter.FilterTemplate)
		})

		AfterEach(func() {
			ctx1 = context.Background()
			commonAfterEach(ctx1, mgr, &combFilter.FilterTemplate)
		})

		It("filterTargetResourceFit", func() {

			_, _, err := commonSetup(ctx1, "generateCombinationsFilter1")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = combFilter.Reconcile(
				ctx1,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf1.String()).ShouldNot(BeEmpty())
			Expect(buf1.String()).NotTo(ContainSubstring("no FunctionTarget found for RegionTypeCandidates:[alveo]"))
			Expect(buf1.String()).To(ContainSubstring("no FunctionTarget found for RegionTypeCandidates:[a100-x]"))
		})
	})

	Context("Test for generateCombinationFilter capacity check", func() {
		var mgr ctrl.Manager
		var combFilter CombinationFilters
		var ctx2 context.Context

		BeforeEach(func() {
			ctx2 = logf.IntoContext(context.Background(), logger2)
			commonBeforeEach(ctx2, mgr, &combFilter.FilterTemplate)
		})

		AfterEach(func() {
			ctx2 = context.Background()
			commonAfterEach(ctx2, mgr, &combFilter.FilterTemplate)
		})

		It("filterTargetResourceFit", func() {

			_, _, err := commonSetup(ctx2, "generateCombinationsFilter2")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = combFilter.Reconcile(
				ctx2,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf2.String()).ShouldNot(BeEmpty())
			Expect(buf2.String()).To(ContainSubstring("function target MaxFunctions is full. FunctionTarget=node1.a100-0.gpu"))
			Expect(buf2.String()).To(ContainSubstring("function target MaxCapacity will inevitably result in capacity over. FunctionTarget=node1.a100-1.gpu"))
			Expect(buf2.String()).NotTo(ContainSubstring("node1.a100-2.gpu"))
			Expect(buf2.String()).To(ContainSubstring("function MaxDataFlows is full. FunctionTarget=node1.alveou250-0.lane0 FunctionIndex=1"))
			Expect(buf2.String()).NotTo(ContainSubstring("node1.alveou250-0.lane1"))
			Expect(buf2.String()).To(ContainSubstring("function MaxCapacity will inevitably result in capacity over. FunctionTarget=node1.alveou250-1.lane0 FunctionIndex=1"))
			Expect(buf2.String()).NotTo(ContainSubstring("node1.alveou250-1.lane1"))
			Expect(buf2.String()).NotTo(ContainSubstring("node1.cpu-0.cpu"))
			Expect(buf2.String()).To(ContainSubstring("function target MaxCapacity will inevitably result in capacity over. FunctionTarget=node1.a100-1.gpu FunctionIndex=1"))
			Expect(buf2.String()).To(ContainSubstring("function target MaxCapacity will inevitably result in capacity over. FunctionTarget=node1.a100-1.gpu FunctionIndex=2"))
		})
	})

	Context("Test for targetResourceFitFilter", func() {
		var mgr ctrl.Manager
		var combFilter CombinationFilters
		var ctx3 context.Context

		BeforeEach(func() {
			ctx3 = logf.IntoContext(context.Background(), logger3)
			commonBeforeEach(ctx3, mgr, &combFilter.FilterTemplate)
		})

		AfterEach(func() {
			ctx3 = context.Background()
			commonAfterEach(ctx3, mgr, &combFilter.FilterTemplate)
		})

		It("targetResourceFitFilter", func() {

			_, _, err := commonSetup(ctx3, "targetResourceFitFilter")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = combFilter.Reconcile(
				ctx3,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf3.String()).ShouldNot(BeEmpty())
			Expect(buf3.String()).To(ContainSubstring("function target MaxCapacity capacity over. FunctionTarget=node1.a100-0.gpu MaxCapacity=100 postDeployCapacity=105"))
			Expect(buf3.String()).NotTo(ContainSubstring("node1.a100-1.gpu"))
			Expect(buf3.String()).To(ContainSubstring("function target MaxFunctions capacity over. FunctionTarget=node1.a100-2.gpu MaxFunctions=120 postDeployFunctions=121"))
			Expect(buf3.String()).NotTo(ContainSubstring("node1.alveou250-0.lane0"))
			Expect(buf3.String()).NotTo(ContainSubstring("node1.alveou250-0.lane1"))
			Expect(buf3.String()).To(ContainSubstring("function MaxCapacity capacity over. FunctionTarget=node1.alveou250-1.lane0 FunctionIndex=1 MaxCapacity=30 postDeployCapacity=45"))
			Expect(buf3.String()).To(ContainSubstring("function MaxDataflows capacity over. FunctionTarget=node1.alveou250-1.lane1 FunctionIndex=1 MaxDataflows=8 postDeployDataflows=9"))
			Expect(buf3.String()).NotTo(ContainSubstring("node1.cpu-0.cpu"))
		})
	})

	Context("Test for connectionResourceFitFilter", func() {
		var mgr ctrl.Manager
		var combFilter CombinationFilters
		var ctx4 context.Context

		BeforeEach(func() {
			ctx4 = logf.IntoContext(context.Background(), logger4)
			commonBeforeEach(ctx4, mgr, &combFilter.FilterTemplate)
		})

		AfterEach(func() {
			ctx4 = context.Background()
			commonAfterEach(ctx4, mgr, &combFilter.FilterTemplate)
		})

		It("connectionResourceFitFilter", func() {

			_, _, err := commonSetup(ctx4, "connectionResourceFitFilter")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = combFilter.Reconcile(
				ctx4,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf4.String()).ShouldNot(BeEmpty())
			Expect(buf4.String()).NotTo(ContainSubstring("capacity over. EntityId=global.ether-network-0"))
			Expect(buf4.String()).NotTo(ContainSubstring("capacity over. EntityId=node0.host100gether-0"))
			Expect(buf4.String()).NotTo(ContainSubstring("capacity over. EntityId=node1.alveou250-0.dev25gether-0"))
			Expect(buf4.String()).NotTo(ContainSubstring("capacity over. EntityId=node1.alveou250-0.dev25gether-1"))
			Expect(buf4.String()).NotTo(ContainSubstring("capacity over. EntityId=node1.alveou250-1.pcie-0"))
			Expect(buf4.String()).NotTo(ContainSubstring("capacity over. EntityId=node1.alveou250-1.dev25gether-0"))
			Expect(buf4.String()).NotTo(ContainSubstring("capacity over. EntityId=node1.alveou250-1.dev25gether-1"))
			Expect(buf4.String()).To(ContainSubstring("device interface outgoing capacity over. EntityId=node1.alveou250-0.pcie-0 Route=node1.alveou250-0.lane1 -> node1.alveou250-0.pcie-0 -> node1.pcie-network-0 -> node1.cpu-0.pcie-0 -> node1.cpu-0 -> node1.cpu-0.pcie-1 -> node1.pcie-network-1 -> node1.a100-1.pcie-0 -> node1.a100-1"))
			Expect(buf4.String()).To(ContainSubstring("device interface outgoing capacity over. EntityId=node1.alveou250-0.pcie-0 Route=node1.alveou250-0.lane1 -> node1.alveou250-0.pcie-0 -> node1.pcie-network-0 -> node1.cpu-0.pcie-0 -> node1.cpu-0 -> node1.cpu-0.pcie-1 -> node1.pcie-network-1 -> node1.a100-1.pcie-0 -> node1.a100-1"))
			Expect(buf4.String()).To(ContainSubstring("network incoming capacity over. EntityId=node1.pcie-network-0 Route=node1.alveou250-0.lane1 -> node1.alveou250-0.pcie-0 -> node1.pcie-network-0 -> node1.cpu-0.pcie-0 -> node1.cpu-0 -> node1.cpu-0.pcie-1 -> node1.pcie-network-1 -> node1.a100-1.pcie-0 -> node1.a100-1"))
			Expect(buf4.String()).To(ContainSubstring("network outgoing capacity over. EntityId=node1.pcie-network-1 Route=node1.alveou250-0.lane1 -> node1.alveou250-0.pcie-0 -> node1.pcie-network-0 -> node1.cpu-0.pcie-0 -> node1.cpu-0 -> node1.cpu-0.pcie-1 -> node1.pcie-network-1 -> node1.a100-1.pcie-0 -> node1.a100-1"))
		})
	})
})

func commonBeforeEach(ctx context.Context, mgr ctrl.Manager, filterTemplate *FilterTemplate) {

	Eventually(func(g Gomega) {
		var err error
		mgr, err = getMgr(mgr)
		addr++
		Expect(err).ShouldNot(HaveOccurred())
	}).Should(Succeed())

	// If the following is not executed, r.Get, r.List, etc. cannot be executed.
	filterTemplate.Client = k8sClient
	filterTemplate.Scheme = testScheme
	filterTemplate.Recorder = mgr.GetEventRecorderFor("combinationfilter-controller")
	Expect(filterTemplate.SetupWithManager(mgr)).Should(Succeed())
}

func commonSetup(ctx context.Context, filterName string) (*ntthpcv1.SchedulingData, *ntthpcv1.DataFlow, error) {
	commonDir := filepath.Join("..", "..", "resources", "combination_filters_test", filterName)
	yamls := []string{
		filepath.Join(commonDir, "functionTargets"),
		filepath.Join(commonDir, "functioninfo"),
		filepath.Join(commonDir, "topology_customresource.yaml"),
		filepath.Join(commonDir, "dataflow.yaml"),
		filepath.Join(commonDir, "scheduling_data.yaml"),
		filepath.Join(commonDir, "strategy.yaml"),
		filepath.Join(commonDir, "user_requirement.yaml"),
	}
	objs, err := Deploy(ctx, k8sClient, yamls...)

	if err != nil {
		return nil, nil, err
	}

	sd := objs[filepath.Join(commonDir, "scheduling_data.yaml")].(*ntthpcv1.SchedulingData)
	df := objs[filepath.Join(commonDir, "dataflow.yaml")].(*ntthpcv1.DataFlow)

	return sd, df, nil

}

func commonAfterEach(ctx context.Context, mgr ctrl.Manager, filterTemplate *FilterTemplate) {

	Eventually(func(g Gomega) {
		var err error
		mgr, err = getMgr(mgr)
		addr++
		Expect(err).ShouldNot(HaveOccurred())
	}).Should(Succeed())

	// If the following is not executed, r.Get, r.List, etc. cannot be executed.
	filterTemplate.Client = k8sClient
	filterTemplate.Scheme = testScheme
	filterTemplate.Recorder = mgr.GetEventRecorderFor("combinationfilter-controller")

	Expect(DeleteAllOf(ctx, k8sClient, "default",
		&ntthpcv1.DataFlow{},
		&ntthpcv1.SchedulingData{},
		&ntthpcv1.TopologyInfo{},
		&ntthpcv1.FunctionTarget{},
		&corev1.ConfigMap{},
	)).Should(Succeed())

	Expect(DeleteAllOf(ctx, k8sClient, "topologyinfo",
		&ntthpcv1.TopologyInfo{},
	)).Should(Succeed())
}

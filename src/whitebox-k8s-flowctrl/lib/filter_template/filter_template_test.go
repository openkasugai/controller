package filter_template_test

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	examplecomv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1" //nolint:stylecheck // ST1019: intentional import as another name
	ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"     //nolint:stylecheck // ST1019: intentional import as another name

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log" //nolint:stylecheck // ST1019: intentional import as another name
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

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
var buf1, buf2, buf3 *bytes.Buffer
var logger1, logger2, logger3 logr.Logger

type FunctionTargetNameAndIndex struct {
	Name          string
	FunctionIndex *int32
}

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
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
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
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("FilterTemplate", func() {

	Context("Test for FetchDeviceType", func() {
		var mgr ctrl.Manager
		var filterTemplate FilterTemplate
		ctx := context.Background()

		BeforeEach(func() {
			commonBeforeEach(ctx, mgr, &filterTemplate)
		})

		It("Filter", func() {

			sd, df, err := commonSetup(ctx)
			Expect(err).ShouldNot(HaveOccurred())

			got, err := filterTemplate.FetchDeviceTypes(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&DeviceTypeFilter{
					DeviceTypes: &[]string{"alveou250"},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			expect := []string{"alveou250"}
			Expect(got).Should(Equal(expect))

		})

	})

	Context("Test for FetchFunctionTarget", func() {
		var mgr ctrl.Manager
		var filterTemplate FilterTemplate
		ctx := context.Background()

		BeforeEach(func() {
			commonBeforeEach(ctx, mgr, &filterTemplate)
		})

		It("Filter NodeNames", func() {

			sd, df, err := commonSetup(ctx)
			Expect(err).ShouldNot(HaveOccurred())

			got, err := filterTemplate.FetchFunctionTargets(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&FunctionTargetFilter{
					NodeNames: &[]string{"node2"},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			exp := []string{
				"node2.a100-0.gpu",
				"node2.a100-1.gpu",
				"node2.a100-2.gpu",
				"node2.alveou250-0.lane0",
				"node2.alveou250-0.lane1",
				"node2.alveou250-1.lane0",
				"node2.alveou250-1.lane1",
				"node2.cpu-0.cpu",
			}

			act := resourceToNameSlice(toAny[*ntthpcv1.FunctionTarget](got))
			l := log.FromContext(context.Background())
			l.Info(fmt.Sprintf("act : %v ", act))
			Expect(sameStringSlice(exp, act)).Should(BeTrue())

		})

		It("Filter DeviceTypes", func() {

			sd, df, err := commonSetup(ctx)
			Expect(err).ShouldNot(HaveOccurred())

			got, err := filterTemplate.FetchFunctionTargets(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&FunctionTargetFilter{
					DeviceTypes: &[]string{"alveou250"},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			exp := []string{
				"node1.alveou250-0.lane0",
				"node1.alveou250-0.lane1",
				"node1.alveou250-1.lane0",
				"node1.alveou250-1.lane1",
				"node2.alveou250-0.lane0",
				"node2.alveou250-0.lane1",
				"node2.alveou250-1.lane0",
				"node2.alveou250-1.lane1",
			}

			act := resourceToNameSlice(toAny[*ntthpcv1.FunctionTarget](got))
			l := log.FromContext(context.Background())
			l.Info(fmt.Sprintf("act : %v ", act))
			Expect(sameStringSlice(exp, act)).Should(BeTrue())

		})

		It("Filter FunctionTargets", func() {

			sd, df, err := commonSetup(ctx)
			Expect(err).ShouldNot(HaveOccurred())

			got, err := filterTemplate.FetchFunctionTargets(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&FunctionTargetFilter{
					FunctionTargets: &[]string{"node1.a100-2.gpu"},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			exp := []string{
				"node1.a100-2.gpu",
			}

			act := resourceToNameSlice(toAny[*ntthpcv1.FunctionTarget](got))
			l := log.FromContext(context.Background())
			l.Info(fmt.Sprintf("act : %v ", act))
			Expect(sameStringSlice(exp, act)).Should(BeTrue())

		})

		It("Filter FunctionNames", func() {

			sd, df, err := commonSetup(ctx)
			Expect(err).ShouldNot(HaveOccurred())

			got, err := filterTemplate.FetchFunctionTargets(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&FunctionTargetFilter{
					FunctionNames: &[]string{"filter-resize"},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			exp := []string{
				"node1.alveou250-0.lane0",
				"node1.alveou250-0.lane1",
				"node2.alveou250-0.lane0",
				"node2.alveou250-0.lane1",
			}

			act := resourceToNameSlice(toAny[*ntthpcv1.FunctionTarget](got))
			l := log.FromContext(context.Background())
			l.Info(fmt.Sprintf("act : %v ", act))
			Expect(sameStringSlice(exp, act)).Should(BeTrue())

		})

		It("Filter RegionNames", func() {

			sd, df, err := commonSetup(ctx)
			Expect(err).ShouldNot(HaveOccurred())

			got, err := filterTemplate.FetchFunctionTargets(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&FunctionTargetFilter{
					RegionNames: &[]string{"lane0"},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			exp := []string{
				"node1.alveou250-0.lane0",
				"node1.alveou250-1.lane0",
				"node2.alveou250-0.lane0",
				"node2.alveou250-1.lane0",
			}

			act := resourceToNameSlice(toAny[*ntthpcv1.FunctionTarget](got))
			l := log.FromContext(context.Background())
			l.Info(fmt.Sprintf("act : %v ", act))
			Expect(sameStringSlice(exp, act)).Should(BeTrue())

		})

		It("Filter RegionTypes", func() {

			sd, df, err := commonSetup(ctx)
			Expect(err).ShouldNot(HaveOccurred())

			got, err := filterTemplate.FetchFunctionTargets(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&FunctionTargetFilter{
					RegionTypes: &[]string{"alveo"},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			exp := []string{
				"node1.alveou250-0.lane0",
				"node1.alveou250-1.lane0",
				"node1.alveou250-0.lane1",
				"node1.alveou250-1.lane1",
				"node2.alveou250-0.lane0",
				"node2.alveou250-1.lane0",
				"node2.alveou250-0.lane1",
				"node2.alveou250-1.lane1",
			}

			act := resourceToNameSlice(toAny[*ntthpcv1.FunctionTarget](got))
			l := log.FromContext(context.Background())
			l.Info(fmt.Sprintf("got : %v ", act))
			l.Info(fmt.Sprintf("act : %v ", act))
			Expect(sameStringSlice(exp, act)).Should(BeTrue())

		})

		It("Filter for AvailableFunctionTargets", func() {

			sd, df, err := commonSetup(ctx)
			Expect(err).ShouldNot(HaveOccurred())

			got, err := filterTemplate.FetchFunctionTargets(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&FunctionTargetFilter{
					IncludesNotAvailable: ValToAddr[bool](false),
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			exp := []string{
				"node1.a100-0.gpu",
				"node1.a100-1.gpu",
				"node1.a100-2.gpu",
				"node1.alveou250-0.lane0",
				"node1.alveou250-1.lane0",
				"node1.alveou250-0.lane1",
				"node1.alveou250-1.lane1",
				"node1.cpu-0.cpu",
				"node2.a100-0.gpu",
				"node2.a100-1.gpu",
				"node2.a100-2.gpu",
				"node2.alveou250-0.lane0",
				"node2.alveou250-1.lane0",
				"node2.alveou250-0.lane1",
				"node2.alveou250-1.lane1",
				"node2.cpu-0.cpu",
			}

			act := resourceToNameSlice(toAny[*ntthpcv1.FunctionTarget](got))
			l := log.FromContext(context.Background())
			l.Info(fmt.Sprintf("got : %v ", act))
			l.Info(fmt.Sprintf("act : %v ", act))
			Expect(exp).To(ConsistOf(act))

		})

		It("Filter for IncludeNotAvailableFunctionTargets", func() {

			sd, df, err := commonSetup(ctx)
			Expect(err).ShouldNot(HaveOccurred())

			got, err := filterTemplate.FetchFunctionTargets(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&FunctionTargetFilter{
					IncludesNotAvailable: ValToAddr[bool](true),
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			exp := []string{
				"node1.a100-0.gpu",
				"node1.a100-1.gpu",
				"node1.a100-2.gpu",
				"node1.alveou250-0.lane0",
				"node1.alveou250-1.lane0",
				"node1.alveou250-2.lane0",
				"node1.alveou250-0.lane1",
				"node1.alveou250-1.lane1",
				"node1.alveou250-2.lane1",
				"node1.cpu-0.cpu",
				"node2.a100-0.gpu",
				"node2.a100-1.gpu",
				"node2.a100-2.gpu",
				"node2.alveou250-0.lane0",
				"node2.alveou250-1.lane0",
				"node2.alveou250-2.lane0",
				"node2.alveou250-0.lane1",
				"node2.alveou250-1.lane1",
				"node2.alveou250-2.lane1",
				"node2.cpu-0.cpu",
			}

			act := resourceToNameSlice(toAny[*ntthpcv1.FunctionTarget](got))
			l := log.FromContext(context.Background())
			l.Info(fmt.Sprintf("got : %v ", act))
			l.Info(fmt.Sprintf("act : %v ", act))
			Expect(exp).To(ConsistOf(act))

		})

		It("No Filter for FunctionTargets", func() {

			sd, df, err := commonSetup(ctx)
			Expect(err).ShouldNot(HaveOccurred())

			got, err := filterTemplate.FetchFunctionTargets(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				nil,
			)
			Expect(err).ShouldNot(HaveOccurred())

			exp := []string{
				"node1.a100-0.gpu",
				"node1.a100-1.gpu",
				"node1.a100-2.gpu",
				"node1.alveou250-0.lane0",
				"node1.alveou250-1.lane0",
				"node1.alveou250-0.lane1",
				"node1.alveou250-1.lane1",
				"node1.cpu-0.cpu",
				"node2.a100-0.gpu",
				"node2.a100-1.gpu",
				"node2.a100-2.gpu",
				"node2.alveou250-0.lane0",
				"node2.alveou250-1.lane0",
				"node2.alveou250-0.lane1",
				"node2.alveou250-1.lane1",
				"node2.cpu-0.cpu",
			}

			act := resourceToNameSlice(toAny[*ntthpcv1.FunctionTarget](got))
			l := log.FromContext(context.Background())
			l.Info(fmt.Sprintf("got : %v ", act))
			l.Info(fmt.Sprintf("act : %v ", act))
			Expect(exp).To(ConsistOf(act))
		})

	})

	Context("Test for FetchFunctionIndexStructs", func() {
		var mgr ctrl.Manager
		var filterTemplate FilterTemplate
		ctx := context.Background()

		BeforeEach(func() {
			commonBeforeEach(ctx, mgr, &filterTemplate)
		})

		It("Filter FunctionIndexes", func() {

			sd, df, err := commonSetup(ctx)
			Expect(err).ShouldNot(HaveOccurred())

			fil := FunctionTargetFilter{
				FunctionIndexes: &[]string{"2"},
			}

			fts, err := filterTemplate.FetchFunctionTargets(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&fil,
			)

			got := filterTemplate.FetchFunctionIndexStructs(
				ctx,
				&fts,
				&fil,
			)

			Expect(err).ShouldNot(HaveOccurred())

			exp := []FunctionTargetNameAndIndex{
				{Name: "node1.alveou250-0.lane0", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node1.alveou250-1.lane0", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node1.a100-0.gpu", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node1.a100-2.gpu", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.alveou250-0.lane0", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.alveou250-1.lane0", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.a100-0.gpu", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.a100-2.gpu", FunctionIndex: ValToAddr[int32](2)},
			}

			act := functionIndexStructsToNameAndIndexSlice(got)
			l := log.FromContext(context.Background())
			l.Info(fmt.Sprintf("got : %v ", got))
			l.Info(fmt.Sprintf("act : %v ", act))
			Expect(exp).To(ConsistOf(act))

		})

		It("Filter AvailableFunctionIndexStructs", func() {

			sd, df, err := commonSetup(ctx)
			Expect(err).ShouldNot(HaveOccurred())

			fil := FunctionTargetFilter{
				IncludesNotAvailable: ValToAddr[bool](false),
			}

			fts, err := filterTemplate.FetchFunctionTargets(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&fil,
			)

			got := filterTemplate.FetchFunctionIndexStructs(
				ctx,
				&fts,
				&fil,
			)

			Expect(err).ShouldNot(HaveOccurred())

			exp := []FunctionTargetNameAndIndex{
				{Name: "node1.alveou250-0.lane0", FunctionIndex: nil},
				{Name: "node1.alveou250-0.lane0", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node1.alveou250-0.lane1", FunctionIndex: nil},
				{Name: "node1.alveou250-0.lane1", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node1.alveou250-1.lane0", FunctionIndex: nil},
				{Name: "node1.alveou250-1.lane0", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node1.alveou250-1.lane1", FunctionIndex: nil},
				{Name: "node1.alveou250-1.lane1", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node1.a100-0.gpu", FunctionIndex: nil},
				{Name: "node1.a100-0.gpu", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node1.a100-1.gpu", FunctionIndex: nil},
				{Name: "node1.a100-1.gpu", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node1.a100-2.gpu", FunctionIndex: nil},
				{Name: "node1.a100-2.gpu", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node1.cpu-0.cpu", FunctionIndex: nil},
				{Name: "node2.alveou250-0.lane0", FunctionIndex: nil},
				{Name: "node2.alveou250-0.lane0", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.alveou250-0.lane1", FunctionIndex: nil},
				{Name: "node2.alveou250-0.lane1", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node2.alveou250-1.lane0", FunctionIndex: nil},
				{Name: "node2.alveou250-1.lane0", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.alveou250-1.lane1", FunctionIndex: nil},
				{Name: "node2.alveou250-1.lane1", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node2.a100-0.gpu", FunctionIndex: nil},
				{Name: "node2.a100-0.gpu", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.a100-1.gpu", FunctionIndex: nil},
				{Name: "node2.a100-1.gpu", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node2.a100-2.gpu", FunctionIndex: nil},
				{Name: "node2.a100-2.gpu", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.cpu-0.cpu", FunctionIndex: nil},
			}

			act := functionIndexStructsToNameAndIndexSlice(got)
			l := log.FromContext(context.Background())
			l.Info(fmt.Sprintf("exp : %v ", exp))
			l.Info(fmt.Sprintf("got : %v ", got))
			Expect(exp).To(ConsistOf(act))

		})

		It("Filter IncludeNotAvailableFunctionIndexStructs", func() {

			sd, df, err := commonSetup(ctx)
			Expect(err).ShouldNot(HaveOccurred())

			fil := FunctionTargetFilter{
				IncludesNotAvailable: ValToAddr[bool](true),
			}

			fts, err := filterTemplate.FetchFunctionTargets(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&fil,
			)

			got := filterTemplate.FetchFunctionIndexStructs(
				ctx,
				&fts,
				&fil,
			)

			Expect(err).ShouldNot(HaveOccurred())

			exp := []FunctionTargetNameAndIndex{
				{Name: "node1.alveou250-0.lane0", FunctionIndex: nil},
				{Name: "node1.alveou250-0.lane0", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node1.alveou250-0.lane0", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node1.alveou250-0.lane1", FunctionIndex: nil},
				{Name: "node1.alveou250-0.lane1", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node1.alveou250-0.lane1", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node1.alveou250-1.lane0", FunctionIndex: nil},
				{Name: "node1.alveou250-1.lane0", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node1.alveou250-1.lane0", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node1.alveou250-1.lane1", FunctionIndex: nil},
				{Name: "node1.alveou250-1.lane1", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node1.alveou250-1.lane1", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node1.alveou250-2.lane0", FunctionIndex: nil},
				{Name: "node1.alveou250-2.lane0", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node1.alveou250-2.lane0", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node1.alveou250-2.lane1", FunctionIndex: nil},
				{Name: "node1.alveou250-2.lane1", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node1.alveou250-2.lane1", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node1.a100-0.gpu", FunctionIndex: nil},
				{Name: "node1.a100-0.gpu", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node1.a100-0.gpu", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node1.a100-1.gpu", FunctionIndex: nil},
				{Name: "node1.a100-1.gpu", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node1.a100-1.gpu", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node1.a100-2.gpu", FunctionIndex: nil},
				{Name: "node1.a100-2.gpu", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node1.a100-2.gpu", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node1.cpu-0.cpu", FunctionIndex: nil},
				{Name: "node2.alveou250-0.lane0", FunctionIndex: nil},
				{Name: "node2.alveou250-0.lane0", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node2.alveou250-0.lane0", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.alveou250-0.lane1", FunctionIndex: nil},
				{Name: "node2.alveou250-0.lane1", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node2.alveou250-0.lane1", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.alveou250-1.lane0", FunctionIndex: nil},
				{Name: "node2.alveou250-1.lane0", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node2.alveou250-1.lane0", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.alveou250-1.lane1", FunctionIndex: nil},
				{Name: "node2.alveou250-1.lane1", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node2.alveou250-1.lane1", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.alveou250-2.lane0", FunctionIndex: nil},
				{Name: "node2.alveou250-2.lane0", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node2.alveou250-2.lane0", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.alveou250-2.lane1", FunctionIndex: nil},
				{Name: "node2.alveou250-2.lane1", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node2.alveou250-2.lane1", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.a100-0.gpu", FunctionIndex: nil},
				{Name: "node2.a100-0.gpu", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node2.a100-0.gpu", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.a100-1.gpu", FunctionIndex: nil},
				{Name: "node2.a100-1.gpu", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node2.a100-1.gpu", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.a100-2.gpu", FunctionIndex: nil},
				{Name: "node2.a100-2.gpu", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node2.a100-2.gpu", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.cpu-0.cpu", FunctionIndex: nil},
			}

			act := functionIndexStructsToNameAndIndexSlice(got)
			l := log.FromContext(context.Background())
			l.Info(fmt.Sprintf("exp : %v ", exp))
			l.Info(fmt.Sprintf("got : %v ", got))
			Expect(exp).To(ConsistOf(act))

		})

		It("No Filter for FunctionIndexStructs", func() {

			sd, df, err := commonSetup(ctx)
			Expect(err).ShouldNot(HaveOccurred())

			fts, err := filterTemplate.FetchFunctionTargets(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				nil,
			)

			got := filterTemplate.FetchFunctionIndexStructs(
				ctx,
				&fts,
				nil,
			)

			Expect(err).ShouldNot(HaveOccurred())

			exp := []FunctionTargetNameAndIndex{
				{Name: "node1.alveou250-0.lane0", FunctionIndex: nil},
				{Name: "node1.alveou250-0.lane0", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node1.alveou250-0.lane1", FunctionIndex: nil},
				{Name: "node1.alveou250-0.lane1", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node1.alveou250-1.lane0", FunctionIndex: nil},
				{Name: "node1.alveou250-1.lane0", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node1.alveou250-1.lane1", FunctionIndex: nil},
				{Name: "node1.alveou250-1.lane1", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node1.a100-0.gpu", FunctionIndex: nil},
				{Name: "node1.a100-0.gpu", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node1.a100-1.gpu", FunctionIndex: nil},
				{Name: "node1.a100-1.gpu", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node1.a100-2.gpu", FunctionIndex: nil},
				{Name: "node1.a100-2.gpu", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node1.cpu-0.cpu", FunctionIndex: nil},
				{Name: "node2.alveou250-0.lane0", FunctionIndex: nil},
				{Name: "node2.alveou250-0.lane0", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.alveou250-0.lane1", FunctionIndex: nil},
				{Name: "node2.alveou250-0.lane1", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node2.alveou250-1.lane0", FunctionIndex: nil},
				{Name: "node2.alveou250-1.lane0", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.alveou250-1.lane1", FunctionIndex: nil},
				{Name: "node2.alveou250-1.lane1", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node2.a100-0.gpu", FunctionIndex: nil},
				{Name: "node2.a100-0.gpu", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.a100-1.gpu", FunctionIndex: nil},
				{Name: "node2.a100-1.gpu", FunctionIndex: ValToAddr[int32](1)},
				{Name: "node2.a100-2.gpu", FunctionIndex: nil},
				{Name: "node2.a100-2.gpu", FunctionIndex: ValToAddr[int32](2)},
				{Name: "node2.cpu-0.cpu", FunctionIndex: nil},
			}

			act := functionIndexStructsToNameAndIndexSlice(got)
			l := log.FromContext(context.Background())
			l.Info(fmt.Sprintf("exp : %v ", exp))
			l.Info(fmt.Sprintf("got : %v ", got))
			Expect(exp).To(ConsistOf(act))

		})

	})

	Context("Test for FetchFunctionTarget and FunctionIndexStructs", func() {
		var mgr ctrl.Manager
		var filterTemplate FilterTemplate
		ctx := context.Background()

		BeforeEach(func() {
			commonBeforeEach(ctx, mgr, &filterTemplate)
		})

		It("Filter FunctionIndexes", func() {

			sd, df, err := commonSetup(ctx)
			Expect(err).ShouldNot(HaveOccurred())

			fil := FunctionTargetFilter{
				FunctionTargets: &[]string{"node1.a100-2.gpu"},
				FunctionIndexes: &[]string{"2"},
			}

			fts, err := filterTemplate.FetchFunctionTargets(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&fil,
			)

			got := filterTemplate.FetchFunctionIndexStructs(
				ctx,
				&fts,
				&fil,
			)

			Expect(err).ShouldNot(HaveOccurred())

			exp := []FunctionTargetNameAndIndex{
				{Name: "node1.a100-2.gpu", FunctionIndex: ValToAddr[int32](2)},
			}

			act := functionIndexStructsToNameAndIndexSlice(got)
			l := log.FromContext(context.Background())
			l.Info(fmt.Sprintf("got : %v ", got))
			l.Info(fmt.Sprintf("act : %v ", act))
			Expect(exp).To(Equal(act))
		})
	})

	Context("Test for no FunctionTarget that meets the conditions", func() {
		var mgr ctrl.Manager
		var filterTemplate FilterTemplate
		var ctx1 context.Context

		BeforeEach(func() {
			ctx1 = logf.IntoContext(context.Background(), logger1)
			commonBeforeEach(ctx1, mgr, &filterTemplate)
		})

		AfterEach(func() {
			ctx1 = context.Background()
			loggingAfterEach(ctx1, mgr, &filterTemplate)
		})

		It("Filter FunctionTarget", func() {

			fil := FunctionTargetFilter{
				FunctionTargets: &[]string{"node1.a100-3.gpu"},
			}

			sd, df, err := commonSetup(ctx1)
			Expect(err).ShouldNot(HaveOccurred())

			_, err = filterTemplate.FetchFunctionTargets(
				ctx1,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&fil,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf1.String()).ShouldNot(BeEmpty())
			Expect(buf1.String()).To(ContainSubstring(
				"no FunctionTarget found that passed FunctionTargetFilter: [node1.a100-3.gpu]"))
		})

		It("Filter FunctionNames", func() {

			sd, df, err := commonSetup(ctx1)
			Expect(err).ShouldNot(HaveOccurred())

			_, err = filterTemplate.FetchFunctionTargets(
				ctx1,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&FunctionTargetFilter{
					FunctionNames: &[]string{"filter-resize-x"},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf1.String()).ShouldNot(BeEmpty())
			Expect(buf1.String()).To(ContainSubstring(
				"no function found that passed FunctionNameFilter: [filter-resize-x]"))
		})

		It("Filter RegionType", func() {

			fil := FunctionTargetFilter{
				RegionTypes: &[]string{"zynq"},
			}

			sd, df, err := commonSetup(ctx1)
			Expect(err).ShouldNot(HaveOccurred())

			_, err = filterTemplate.FetchFunctionTargets(
				ctx1,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&fil,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf1.String()).ShouldNot(BeEmpty())
			Expect(buf1.String()).To(ContainSubstring(
				"no FunctionTarget found that passed RegionTypeFilter: [zynq]"))
		})

		It("Filter NodeName", func() {

			fil := FunctionTargetFilter{
				NodeNames: &[]string{"node99"},
			}

			sd, df, err := commonSetup(ctx1)
			Expect(err).ShouldNot(HaveOccurred())

			_, err = filterTemplate.FetchFunctionTargets(
				ctx1,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&fil,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf1.String()).ShouldNot(BeEmpty())
			Expect(buf1.String()).To(ContainSubstring(
				"no FunctionTarget found that passed NodeNameFilter: [node99]"))
		})

		It("Filter DeviceType", func() {

			fil := FunctionTargetFilter{
				DeviceTypes: &[]string{"alveou280"},
			}

			sd, df, err := commonSetup(ctx1)
			Expect(err).ShouldNot(HaveOccurred())

			_, err = filterTemplate.FetchFunctionTargets(
				ctx1,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&fil,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf1.String()).ShouldNot(BeEmpty())
			Expect(buf1.String()).To(ContainSubstring(
				"no FunctionTarget found that passed DeviceTypeFilter: [alveou280]"))
		})

		It("Filter RegionName", func() {

			fil := FunctionTargetFilter{
				RegionNames: &[]string{"lane99"},
			}

			sd, df, err := commonSetup(ctx1)
			Expect(err).ShouldNot(HaveOccurred())

			_, err = filterTemplate.FetchFunctionTargets(
				ctx1,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&fil,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf1.String()).ShouldNot(BeEmpty())
			Expect(buf1.String()).To(ContainSubstring(
				"no FunctionTarget found that passed RegionNameFilter: [lane99]"))
		})
	})

	Context("Test for no Function that meets the conditions", func() {
		var mgr ctrl.Manager
		var filterTemplate FilterTemplate
		var ctx2 context.Context

		BeforeEach(func() {
			ctx2 = logf.IntoContext(context.Background(), logger2)
			loggingBeforeEach(ctx2, mgr, &filterTemplate)
		})

		AfterEach(func() {
			ctx2 = context.Background()
			loggingAfterEach(ctx2, mgr, &filterTemplate)
		})

		It("Filter FunctionIndex", func() {

			fil := FunctionTargetFilter{
				FunctionIndexes: &[]string{"3"},
			}

			sd, df, err := commonSetup(ctx2)
			Expect(err).ShouldNot(HaveOccurred())

			fts, err := filterTemplate.FetchFunctionTargets(
				ctx2,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&fil,
			)

			_ = filterTemplate.FetchFunctionIndexStructs(
				ctx2,
				&fts,
				&fil,
			)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf2.String()).ShouldNot(BeEmpty())
			Expect(buf2.String()).To(ContainSubstring(
				"no function found that passed FunctionIndexFilter: [3]"))
		})
	})

	Context("Test for FunctionTarget that meets the conditions", func() {
		var mgr ctrl.Manager
		var filterTemplate FilterTemplate
		var ctx3 context.Context

		BeforeEach(func() {
			ctx3 = logf.IntoContext(context.Background(), logger3)
			commonBeforeEach(ctx3, mgr, &filterTemplate)
		})

		AfterEach(func() {
			ctx3 = context.Background()
			loggingAfterEach(ctx3, mgr, &filterTemplate)
		})

		It("Filter all", func() {

			fil := FunctionTargetFilter{
				FunctionTargets: &[]string{"node1.alveou250-0.lane0"},
				FunctionNames:   &[]string{"filter-resize"},
				NodeNames:       &[]string{"node1"},
				DeviceTypes:     &[]string{"alveou250"},
				RegionNames:     &[]string{"lane0"},
				FunctionIndexes: &[]string{"2"},
			}

			sd, df, err := commonSetup(ctx3)
			Expect(err).ShouldNot(HaveOccurred())

			fts, err := filterTemplate.FetchFunctionTargets(
				ctx3,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "test",
					},
				},
				sd,
				df,
				&fil,
			)

			got := filterTemplate.FetchFunctionIndexStructs(
				ctx3,
				&fts,
				&fil,
			)

			exp := []FunctionTargetNameAndIndex{
				{Name: "node1.alveou250-0.lane0", FunctionIndex: ValToAddr[int32](2)},
			}

			act := functionIndexStructsToNameAndIndexSlice(got)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(exp).To(Equal(act))
			Expect(buf3.String()).NotTo(ContainSubstring("no FunctionTarget found"))
			Expect(buf3.String()).NotTo(ContainSubstring("no function found"))
		})
	})
})

func sameStringSlice(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	// create a map of string -> int
	diff := make(map[string]int, len(x))
	for _, _x := range x {
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++
	}
	for _, _y := range y {
		// If the string _y is not in diff bail out early
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y]--
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	return len(diff) == 0
}

type gettableName interface {
	GetName() string
}

func toAny[T any](in []T) []any {
	ret := make([]any, len(in))
	for i, v := range in {
		ret[i] = v
	}
	return ret
}

func resourceToNameSlice(resources []any) []string {
	names := make([]string, len(resources))
	for i, resource := range resources {
		names[i] = resource.(gettableName).GetName()
	}
	return names
}

func functionIndexStructsToNameAndIndexSlice(
	functionIndexStructs []FunctionIndexStruct,
) []FunctionTargetNameAndIndex {
	ret := make([]FunctionTargetNameAndIndex, len(functionIndexStructs))
	for i, fis := range functionIndexStructs {
		ret[i] = FunctionTargetNameAndIndex{
			Name:          fis.FunctionTarget.GetName(),
			FunctionIndex: fis.FunctionIndex,
		}
	}
	return ret
}

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
	filterTemplate.Recorder = mgr.GetEventRecorderFor("wbscheduler-controller")
	Expect(filterTemplate.SetupWithManager(mgr)).Should(Succeed())

	Expect(DeleteAllOf(ctx, k8sClient, "default",
		&ntthpcv1.DataFlow{},
		&ntthpcv1.SchedulingData{},
		&ntthpcv1.FunctionTarget{},
	)).Should(Succeed())
}

func commonSetup(ctx context.Context) (*ntthpcv1.SchedulingData, *ntthpcv1.DataFlow, error) {
	commonDir := filepath.Join("resources", "common")
	yamls := []string{
		filepath.Join(commonDir, "functionTargets"),
		filepath.Join(commonDir, "dataflow.yaml"),
		filepath.Join(commonDir, "scheduling_data.yaml"),
	}
	objs, err := Deploy(ctx, k8sClient, yamls...)

	if err != nil {
		return nil, nil, err
	}

	sd := objs[filepath.Join(commonDir, "scheduling_data.yaml")].(*ntthpcv1.SchedulingData)
	df := objs[filepath.Join(commonDir, "dataflow.yaml")].(*ntthpcv1.DataFlow)

	return sd, df, nil

}

func loggingBeforeEach(ctx context.Context, mgr ctrl.Manager, filterTemplate *FilterTemplate) {

	Eventually(func(g Gomega) {
		var err error
		mgr, err = getMgr(mgr)
		addr++
		Expect(err).ShouldNot(HaveOccurred())
	}).Should(Succeed())

	// If the following is not executed, r.Get, r.List, etc. cannot be executed.
	filterTemplate.Client = k8sClient
	filterTemplate.Scheme = testScheme
	filterTemplate.Recorder = mgr.GetEventRecorderFor("wbscheduler-controller")
	Expect(filterTemplate.SetupWithManager(mgr)).Should(Succeed())
}

func loggingAfterEach(ctx context.Context, mgr ctrl.Manager, filterTemplate *FilterTemplate) {

	Eventually(func(g Gomega) {
		var err error
		mgr, err = getMgr(mgr)
		addr++
		Expect(err).ShouldNot(HaveOccurred())
	}).Should(Succeed())

	// If the following is not executed, r.Get, r.List, etc. cannot be executed.
	filterTemplate.Client = k8sClient
	filterTemplate.Scheme = testScheme
	filterTemplate.Recorder = mgr.GetEventRecorderFor("wbscheduler-controller")
	Expect(filterTemplate.SetupWithManager(mgr)).Should(Succeed())

	Expect(DeleteAllOf(ctx, k8sClient, "default",
		&ntthpcv1.DataFlow{},
		&ntthpcv1.SchedulingData{},
		&ntthpcv1.FunctionTarget{},
	)).Should(Succeed())
}

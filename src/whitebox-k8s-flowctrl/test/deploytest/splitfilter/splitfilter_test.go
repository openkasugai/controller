/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED
*/

package splitfilter_test

import (
	"context"
	"testing"
	"time"

	v1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	. "github.com/compsysg/whitebox-k8s-flowctrl/test/lib/testutils"
	test_cases "github.com/compsysg/whitebox-k8s-flowctrl/test/test_cases/splitfilter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var k8sClient client.Client

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	SetDefaultEventuallyTimeout(10 * time.Second)
	SetDefaultEventuallyPollingInterval(1 * time.Millisecond)
	SetDefaultConsistentlyDuration(10 * time.Second)
	SetDefaultConsistentlyPollingInterval(100 * time.Millisecond)

	RunSpecs(t, "Controller Suite")
}

var checkNamespaces []string = []string{
	"default",
	"test01",
	"wbfunc-imgproc",
	"chain-imgproc",
	"cluster01",
	"topologyinfo",
}

var _ = BeforeSuite(func() {

	By("bootstrapping test environment")
	var err error

	err = v1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	cfg, err := config.GetConfig()
	Expect(err).NotTo(HaveOccurred())

	k8sClient, err = client.New(cfg, client.Options{})
	Expect(err).NotTo(HaveOccurred())

	namespaces := &corev1.NamespaceList{}
	Expect(k8sClient.List(context.Background(), namespaces)).Should(Succeed())

	for _, name := range checkNamespaces {
		found := false
		for _, ns := range namespaces.Items {
			if ns.Name == name {
				found = true
				break
			}
		}
		if !found {
			ns := &corev1.Namespace{}
			ns.Name = name
			Expect(k8sClient.Create(context.Background(), ns)).Should(Succeed())
		}
	}

})

var _ = AfterSuite(func() {
	deleteTestResource()
})

var _ = Describe("SplitFilterDeploy", func() {

	ctx := context.Background()

	BeforeEach(func() {
		deleteTestResource()
	})

	AfterEach(func() {
		time.Sleep(100 * time.Millisecond)
	})

	// Separate the Context depending on whether or not you use the Scheduler.
	Context("WithScheduler", func() {

		// Normal TopologyInfo not found
		It("1_nominal", func() {
			test_cases.Test1_Nominal(ctx, k8sClient)
		})

		// Normal TopologyInfo available
		It("2_Connection", func() {
			test_cases.Test2_Connection(ctx, k8sClient)
		})

		// OnePass
		It("3_OnePass", func() {
			test_cases.Test3_OnePass(ctx, k8sClient)
		})

		// Normal TopologyInfo not found, various DataFlows
		It("4_1_Nominal_DataFlows_FPGADecode", func() {
			test_cases.Test4_Nominal_DataFlows(ctx, k8sClient, "fpga_decode")
		})

		It("4_2_Nominal_DataFlows_CPUDecode", func() {
			test_cases.Test4_Nominal_DataFlows(ctx, k8sClient, "cpu_decode")
		})

		It("4_3_Nominal_DataFlows_CopyBranch", func() {
			test_cases.Test4_Nominal_DataFlows(ctx, k8sClient, "copy_branch")
		})

		It("4_4_Nominal_DataFlows_GlueFDMAtoTCP", func() {
			test_cases.Test4_Nominal_DataFlows(ctx, k8sClient, "glue_fdma_to_tcp")
		})

		// Normal TopologyInfo available, various DataFlows
		It("5_1_Connection_DataFlows_FPGADecode", func() {
			test_cases.Test5_Connection_DataFlows(ctx, k8sClient, "fpga_decode")
		})

		It("5_2_Connection_DataFlows_CPUDecode", func() {
			test_cases.Test5_Connection_DataFlows(ctx, k8sClient, "cpu_decode")
		})

		It("5_3_Connection_DataFlows_CopyBranch", func() {
			test_cases.Test5_Connection_DataFlows(ctx, k8sClient, "copy_branch")
		})

		It("5_4_Connection_DataFlows_GlueFDMAtoTCP", func() {
			test_cases.Test5_Connection_DataFlows(ctx, k8sClient, "glue_fdma_to_tcp")
		})

		It("5_4_Connection_DataFlows_GlueFDMAtoTCP", func() {
			test_cases.Test5_Connection_DataFlows(ctx, k8sClient, "glue_fdma_to_tcp")
		})

		// Dynamic reconfiguration
		It("6_1_Dynamic_Reconfiguraion_NoTopology", func() {
			test_cases.Test6_Dynamic_Reconfiguration(ctx, k8sClient, "no_topology")
		})

		It("6_2_Dynamic_Reconfiguraion_WithTopology", func() {
			test_cases.Test6_Dynamic_Reconfiguration(ctx, k8sClient, "with_topology")
		})

		// Scheduling retry
		It("7_Schduling_Retry_UpdateFt", func() {
			test_cases.Test7_Scheduling_Retry(ctx, k8sClient, "update_ft")
		})
	})
})

//nolint:unused // createExpYaml is unused. FIXME: remove this function
func createExpYaml(ctx context.Context, dest, name, nameSpace, apiVersion, expStatus string, resource any) {

	switch v := resource.(type) {
	case v1.DataFlow:
		Eventually(func(g Gomega) {
			g.Expect(k8sClient.Get(
				ctx, client.ObjectKey{Namespace: nameSpace, Name: name}, &v)).To(Succeed())
			g.Expect(v.Status.Status).Should(Equal(expStatus))
			g.Expect(GenerateExpectYaml(&v, apiVersion, "DataFlow", dest)).Should(Succeed())
		}).Should(Succeed())
	case v1.SchedulingData:
		Eventually(func(g Gomega) {
			g.Expect(k8sClient.Get(
				ctx, client.ObjectKey{Namespace: nameSpace, Name: name}, &v)).To(Succeed())
			g.Expect(v.Status.Status).Should(Equal(expStatus))
			g.Expect(GenerateExpectYaml(&v, apiVersion, "SchedulingData", dest)).Should(Succeed())
		}).Should(Succeed())
	}

}

func deleteTestResource() {

	ctx := context.Background()

	// In the TopologyInfo Controller specification, if you delete TopologyInfo first,
	// ConfigMap cannot be deleted
	// Delete for adding deleteStamp
	Expect(DeleteAllOf(ctx, k8sClient, "topologyinfo"), &corev1.ConfigMap{}).Should(Succeed())
	time.Sleep(time.Second * 1)

	// Delete the resource used in the previous test.
	for _, ns := range checkNamespaces {

		Expect(DeleteAllOf(ctx, k8sClient, ns,
			&corev1.ConfigMap{},
			&v1.SchedulingData{},
			&v1.DataFlow{},
			&v1.FunctionTarget{},
			&v1.ComputeResource{},
			&v1.FunctionType{},
			&v1.FunctionChain{},
		)).Should(Succeed())
	}

}

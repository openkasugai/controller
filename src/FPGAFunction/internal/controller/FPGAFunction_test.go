package controller

import (
	"bytes"
	"context"
	"fmt"
	"os"

	ctrl "sigs.k8s.io/controller-runtime"
	//	"sigs.k8s.io/controller-runtime/pkg/client"

	examplecomv1 "FPGAFunction/api/v1"
	controllertestcpu "FPGAFunction/internal/controller/test/type/CPU"
	controllertestethernet "FPGAFunction/internal/controller/test/type/Ethernet"
	controllertestgpu "FPGAFunction/internal/controller/test/type/GPU"
	controllertestpcie "FPGAFunction/internal/controller/test/type/PCIe"

	// ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	//"k8s.io/apimachinery/pkg/runtime"
)

func getMgr(mgr ctrl.Manager) (ctrl.Manager, error) {
	if mgr == nil {
		return ctrl.NewManager(cfg, ctrl.Options{
			Scheme: testScheme,
		})
	}
	return mgr, nil
}

// Create FPGAFunctionCR
func createFPGAFunction(ctx context.Context, fpgafcr examplecomv1.FPGAFunction) error {
	tmp := &examplecomv1.FPGAFunction{}
	*tmp = fpgafcr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create EthernetConnectionCR
func createEthernetConnection(ctx context.Context, ethercr controllertestethernet.EthernetConnection) error {
	tmp := &controllertestethernet.EthernetConnection{}
	*tmp = ethercr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create PCIeConnectionCR
func createPCIeConnection(ctx context.Context, pciecr controllertestpcie.PCIeConnection) error {
	tmp := &controllertestpcie.PCIeConnection{}
	*tmp = pciecr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create CPUFunctionCR
func createCPUFunction(ctx context.Context, cpufcr controllertestcpu.CPUFunction) error {
	tmp := &controllertestcpu.CPUFunction{}
	*tmp = cpufcr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create GPUFunctionCR
func createGPUFunction(ctx context.Context, gpufcr controllertestgpu.GPUFunction) error {
	tmp := &controllertestgpu.GPUFunction{}
	*tmp = gpufcr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create ChildBitstreamCR
func createChildBitstream(ctx context.Context, childbscr examplecomv1.ChildBs) error {
	tmp := &examplecomv1.ChildBs{}
	*tmp = childbscr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Update ChildBitstreamCR
func updateChildBitstream(ctx context.Context, childbscr examplecomv1.ChildBs) error {
	tmp := &examplecomv1.ChildBs{}
	*tmp = childbscr
	err := k8sClient.Update(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create FPGACR
func createFPGA(ctx context.Context, fpgacr examplecomv1.FPGA) error {
	tmp := &examplecomv1.FPGA{}
	*tmp = fpgacr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Update FPGACR
/*
func updateFPGA(ctx context.Context, fpgacr examplecomv1.FPGA) error {
	tmp := &examplecomv1.FPGA{}
	*tmp = fpgacr
	err := k8sClient.Update(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}
*/
// Create ComfigMap
func config(ctx context.Context, conf corev1.ConfigMap) error {
	tmp := &corev1.ConfigMap{}
	*tmp = conf
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

var _ = Describe("FPGAFunctionController", func() {
	var mgr ctrl.Manager
	var err error
	ctx := context.Background()
	var fakerecorder = record.NewFakeRecorder(10)
	var writer = bytes.Buffer{}

	//	var reconciler FPGAFunctionReconciler
	Context("Test for FPGAFunctionReconciler", Ordered, func() {
		var reconciler FPGAFunctionReconciler
		BeforeAll(func() {

			opts := zap.Options{
				TimeEncoder: zapcore.ISO8601TimeEncoder,
				Development: true,
				DestWriter:  &writer,
			}
			ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

			// set manager
			mgr, err = getMgr(mgr)
			Expect(err).NotTo(HaveOccurred())

		})

		BeforeEach(func() {
			writer.Reset()

			reconciler = FPGAFunctionReconciler{
				Client:   k8sClient,
				Scheme:   testScheme,
				Recorder: fakerecorder,
			}
			err = reconciler.SetupWithManager(mgr)
			if err != nil {
				fmt.Println("Error in SetupWithManager")
			}
			Expect(err).NotTo(HaveOccurred())
			os.Setenv("K8S_NODENAME", "test01")
			StartupProccessing(mgr)
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestethernet.EthernetConnection{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestpcie.PCIeConnection{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestgpu.GPUFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestcpu.CPUFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGAFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGA{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.ChildBs{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
		})

		It("FPGAFunctionTest 3", func() {
			By("Test Start")

			// Create ConfigMap
			err = config(ctx, fpgafuncconfig_fr_high_infer_2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection3)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR : ", err)
			}

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection4)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR : ", err)
			}

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction2)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR : ", err)
			}

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction2)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR : ", err)
			}

			// Create ChildBitstreamCR
			err = createChildBitstream(ctx, ChildBitstream1)
			if err != nil {
				fmt.Println("There is a problem in createing ChildBitstream ", err)
			}

			// Create FPGACR
			err = createFPGA(ctx, FPGA2[0])
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR ", err)
			}

			// Create FPGAFunctionCR
			err = createFPGAFunction(ctx, FPGAFunction2)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}

			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night02-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}
			fmt.Println("test3:" + writer.String())
			Expect(err).NotTo(HaveOccurred())
		})
		It("FPGAFunctionTest 4", func() {
			By("Test Start")

			// Create ConfigMap
			err = config(ctx, fpgafuncconfig_fr_high_infer_3)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, servicerMgmtConfig2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, deployinfo_configdata2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, regionUniqueInfoConfig2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, functionUniqueInfoConfig2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, fr_childbs_Config2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection3)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR : ", err)
			}

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection4)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR : ", err)
			}

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction3)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR : ", err)
			}

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction3)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR : ", err)
			}
			// Create FPGACR
			err = createFPGA(ctx, FPGA3[0])
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR ", err)
			}

			// Create FPGAFunctionCR
			err = createFPGAFunction(ctx, FPGAFunction3)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}

			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night03-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}
			fmt.Println("test4:" + writer.String())
			Expect(err).NotTo(HaveOccurred())
		})

		It("FPGAFunctionTest 5", func() {
			By("Test Start")

			// Create ConfigMap
			err = config(ctx, fpgafuncconfig_fr_high_infer_4)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection5)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR : ", err)
			}

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection6)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR : ", err)
			}

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction4)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR : ", err)
			}

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction4)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR : ", err)
			}

			// Create ChildBitstreamCR
			err = createChildBitstream(ctx, ChildBitstream2)
			if err != nil {
				fmt.Println("There is a problem in createing ChildBitstream ", err)
			}
			err = createFPGA(ctx, FPGA4[0])
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR ", err)
			}
			err = createFPGAFunction(ctx, FPGAFunction4)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}
			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night04-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}
			fmt.Println("test5:" + writer.String())
			Expect(err).NotTo(HaveOccurred())
		})
		It("FPGAFunctionTest 6", func() {
			By("Test Start")

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection7)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR : ", err)
			}

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection8)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR : ", err)
			}

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection5)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR : ", err)
			}

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection6)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR : ", err)
			}

			// Create ConfigMap
			err = config(ctx, fpgafuncconfig_fr_high_infer_2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction2)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR : ", err)
			}

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction2)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR : ", err)
			}

			// Create ChildBitstreamCR
			err = createChildBitstream(ctx, ChildBitstream1)
			if err != nil {
				fmt.Println("There is a problem in createing ChildBitstream ", err)
			}

			// Create FPGACR
			err = createFPGA(ctx, FPGA2[0])
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR ", err)
			}

			// Create FPGAFunctionCR
			err = createFPGAFunction(ctx, FPGAFunction5)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR  A", err)
			}

			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night05-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			err = createFPGAFunction(ctx, FPGAFunction6)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night06-wbfunction-filter-resize-high-infer-main",
			}})
			fmt.Println("test6:" + writer.String())
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			By("Test End")
			writer.Reset()
			//		fmt.Println("Test End")
		})
		AfterAll(func() {
			By("Test End")
			writer.Reset()
			//		fmt.Println("Test End")
		})
	})

	var fakerecorder2 = record.NewFakeRecorder(10)
	var writer2 = bytes.Buffer{}

	//	var reconciler FPGAFunctionReconciler
	Context("Test for FPGAFunctionReconciler2", Ordered, func() {
		var reconciler2 FPGAFunctionReconciler
		BeforeAll(func() {

			opts := zap.Options{
				TimeEncoder: zapcore.ISO8601TimeEncoder,
				Development: true,
				DestWriter:  &writer2,
			}
			ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

		})

		BeforeEach(func() {
			writer2.Reset()

			// set manager
			mgr, err = getMgr(mgr)
			Expect(err).NotTo(HaveOccurred())

			reconciler2 = FPGAFunctionReconciler{
				Client:   k8sClient,
				Scheme:   testScheme,
				Recorder: fakerecorder2,
			}
			err = reconciler2.SetupWithManager(mgr)
			if err != nil {
				fmt.Println("Error in SetupWithManager")
			}
			Expect(err).NotTo(HaveOccurred())
			os.Setenv("K8S_NODENAME", "test01")
			StartupProccessing(mgr)
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestethernet.EthernetConnection{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestpcie.PCIeConnection{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestgpu.GPUFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestcpu.CPUFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGAFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGA{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.ChildBs{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
		})

		It("FPGAFunctionTest 2", func() {
			By("Test Start")

			// Create ConfigMap
			err = config(ctx, fpgafuncconfig_fr_high_infer_1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, servicerMgmtConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, deployinfo_configdata1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, regionUniqueInfoConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, functionUniqueInfoConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, fr_childbs_Config1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection1)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR : ", err)
			}

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection2)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR : ", err)
			}

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction1)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR : ", err)
			}

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction1)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR : ", err)
			}

			// Create ChildBitstreamCR
			err = createChildBitstream(ctx, ChildBitstream1)
			if err != nil {
				fmt.Println("There is a problem in createing ChildBitstream ", err)
			}

			// Create FPGACR
			err = createFPGA(ctx, FPGA1[0])
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR ", err)
			}

			// Create FPGAFunctionCR
			err = createFPGAFunction(ctx, FPGAFunction1)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}

			_, err := reconciler2.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			var childBsCR examplecomv1.ChildBs
			// Get ChildBs CR
			_ = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00dtest0100000001",
			},
				&childBsCR)

			childBsCR.Status.Regions = childBsCR.Spec.Regions
			childBsCR.Status.Status = examplecomv1.ChildBsStatusPreparing
			childBsCR.Status.State = examplecomv1.ChildBsReady

			// update ChildBs CR
			err = updateChildBitstream(ctx, childBsCR)
			if err != nil {
				fmt.Println("There is a problem in updateing ChildBsCR ", err)
			}

			_, err = reconciler2.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}
			// Get ChildBs CR
			_ = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00dtest0100000001",
			},
				&childBsCR)

			childBsCR.Status.Regions = childBsCR.Spec.Regions
			childBsCR.Status.Status = examplecomv1.ChildBsStatusReady
			childBsCR.Status.State = examplecomv1.ChildBsReady

			// update ChildBs CR
			err = updateChildBitstream(ctx, childBsCR)
			if err != nil {
				fmt.Println("There is a problem in updateing ChildBsCR ", err)
			}

			_, err = reconciler2.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}
			fmt.Println("test2:" + writer2.String())

			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			By("Test End")
			writer2.Reset()
		})
	})
})

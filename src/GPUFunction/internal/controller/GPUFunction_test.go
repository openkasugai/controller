package controller

import (
	"context"
	"fmt"
	"os"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	examplecomv1 "GPUFunction/api/v1"

	// ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// Additional files
	"k8s.io/apimachinery/pkg/types"
	//"sigs.k8s.io/controller-runtime/pkg/client"
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

// Create GPUFunction CR
func createGPUFunction(ctx context.Context, gpufcr examplecomv1.GPUFunction) error {
	tmp := &examplecomv1.GPUFunction{}
	*tmp = gpufcr
	fmt.Println("gpufunc tmp")
	fmt.Println(*tmp)
	tmp.TypeMeta = gpufcr.TypeMeta
	err := k8sClient.Create(ctx, tmp)
	//	err := k8sClient.Create(context.Background(), &GPUFunction)
	if err != nil {
		return err
	}
	// tmp.Status = gpufcr.Status
	err = k8sClient.Status().Update(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create PCIeConnection CR

func createPCIeConnection(ctx context.Context, pcieccr PCIeConnection) error {
	tmp := &PCIeConnection{}
	*tmp = pcieccr
	fmt.Println("pcie tmp")
	fmt.Println(*tmp)
	tmp.Kind = pcieccr.Kind
	tmp.APIVersion = pcieccr.APIVersion
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}

	// tmp.Status = pcieccr.Status
	// err = k8sClient.Status().Update(ctx, tmp)
	// if err != nil {
	// 	return err
	// }
	return nil
}

/*
// Create EthernetConnection CR

	func createEthernetConnection(ctx context.Context, ethernetccr EthernetConnection) error {
		tmp := &EthernetConnection{}
		*tmp = ethernetccr
		tmp.Kind = ethernetccr.Kind
		tmp.APIVersion = ethernetccr.APIVersion
		err := k8sClient.Create(ctx, tmp)
		if err != nil {
			return err
		}
		return nil
	}
*/
func createFPGAFunction(ctx context.Context, fpgafcr FPGAFunction) error {
	tmp := &FPGAFunction{}
	*tmp = fpgafcr
	fmt.Println("fpga tmp:")
	fmt.Println(*tmp)
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	// tmp.Status = fpgafcr.Status
	// err = k8sClient.Status().Update(ctx, tmp)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func createGPUFuncConfig(ctx context.Context, gpuconf corev1.ConfigMap) error {
	tmp := &corev1.ConfigMap{}
	*tmp = gpuconf
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

func createNetworkAttachmentDefinition(ctx context.Context, nad NetworkAttachmentDefinition) error {
	tmp := &NetworkAttachmentDefinition{}
	*tmp = nad
	// tmp.Kind = NetworkAttachmentDefinition1.Kind
	// tmp.APIVersion = NetworkAttachmentDefinition1.APIVersion
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

var _ = Describe("GPUFunctionController", func() {
	var mgr ctrl.Manager
	var err error
	ctx := context.Background()
	//	var reconciler GPUFunctionReconciler
	Context("Test for GPUFunctionReconciler", Ordered, func() {
		var reconciler GPUFunctionReconciler
		BeforeAll(func() {
			// set manager
			mgr, err = getMgr(mgr)
			Expect(err).NotTo(HaveOccurred())

			reconciler = GPUFunctionReconciler{
				Client:   k8sClient,
				Scheme:   testScheme,
				Recorder: mgr.GetEventRecorderFor("gpufunction-controller"),
			}
			err = reconciler.SetupWithManager(mgr)
			if err != nil {
				fmt.Println("Error in SetupWithManager")
			}
			Expect(err).NotTo(HaveOccurred())
			namespace := &corev1.Namespace{}
			/*
				namespace := &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "gpufunctiontest",
						Namespace: "default",
						},
					}
			*/
			namespace.Name = "test01"
			err := k8sClient.Create(context.Background(), namespace)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
		})

		BeforeEach(func() {
			os.Setenv("K8S_NODENAME", "worker1")
			os.Setenv("K8S_POD_ANNOTATION", "")
			mgr, err = getMgr(mgr)
			Expect(err).NotTo(HaveOccurred())
			/*				reconciler = GPUFunctionReconciler{
								Client:   k8sClient,
								Scheme:   testScheme,
								Recorder: mgr.GetEventRecorderFor("gpufunction-controller"),
							}
			*/
		})
		It("Reconcile Test", func() {
			By("Test Start")

			// Create GPUFuncConfig
			err = createGPUFuncConfig(ctx, gpuconfigdecode)
			if err != nil {
				fmt.Println("There is a problem in createing GPUConfig ", err)
				fmt.Printf("%T\n", err)
				fmt.Println(err)
			}

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection1)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}

			// Create FPGAFunctionCR
			err = createFPGAFunction(ctx, FPGAFunction1)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}

			// Create NetworkAttachmentDefinition
			err = createNetworkAttachmentDefinition(ctx, NetworkAttachmentDefinition1)
			if err != nil {
				fmt.Println("There is a problem in createing NetworkAttachmentDefinition")
			}

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction1)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR")
			}

			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}
			// Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var gpuCR examplecomv1.GPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "df-night01-wbfunction-high-infer-main",
				Namespace: "default",
			}, &gpuCR)
			if err != nil {
				// Error route
				fmt.Println("Cannot get GPUFunctionCR:", gpuCR, err)
			}
			// Expect(err).NotTo(HaveOccurred())

			Expect(*gpuCR.Status.FunctionIndex).To(Equal(int32(0)))

			// Delete GPUFunctionCR
			var deltime metav1.Time
			deltime = metav1.Now()
			gpuCR.DeletionTimestamp = &deltime
			//			err = deleteGPUFunction(ctx,gpuCR)

			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction1)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR")
			}

			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}
			// Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var gpuCR2 examplecomv1.GPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "df-night01-wbfunction-high-infer-main",
				Namespace: "default",
			}, &gpuCR2)
			if err != nil {
				// Error route
				fmt.Println("Cannot get GPUFunctionCR:", gpuCR2, err)
			}
			// Expect(err).NotTo(HaveOccurred())

			Expect(*gpuCR2.Status.FunctionIndex).To(Equal(int32(0)))
			gpuCR2.DeletionTimestamp = &deltime
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}
		})
		It("Reconcile Test2", func() {
			By("Test Start")

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection2)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction2)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR")
			}

			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night02-wbfunction-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}
			// Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var gpuCR3 examplecomv1.GPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "df-night02-wbfunction-high-infer-main",
				Namespace: "default",
			}, &gpuCR3)
			if err != nil {
				// Error route
				fmt.Println("Cannot get GPUFunctionCR:", gpuCR3, err)
			}
			// Expect(err).NotTo(HaveOccurred())
			fmt.Println(gpuCR3.Spec)
			Expect(*gpuCR3.Status.FunctionIndex).To(Equal(int32(99)))
		})
		AfterEach(func() {
			By("Test End")
			//		fmt.Println("Test End")
		})
	})
})

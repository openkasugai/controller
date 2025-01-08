/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED
*/

package controller

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap/zapcore"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	examplecomv1 "PCIeConnection/api/v1"
	controllertestcpu "PCIeConnection/internal/controller/test/type/CPU"
	controllertestfpga "PCIeConnection/internal/controller/test/type/FPGA"
	controllertestgpu "PCIeConnection/internal/controller/test/type/GPU"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	//	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"

	"k8s.io/apimachinery/pkg/types"
)

func getMgr(mgr ctrl.Manager) (ctrl.Manager, error) {
	if mgr == nil {
		// return ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		return ctrl.NewManager(cfg, ctrl.Options{
			// Scheme: testScheme,
			Scheme: k8sClient.Scheme(),
		})
	}
	return mgr, nil
	// return ctrl.NewManager(cfg, ctrl.Options{
	// 	// Scheme: testScheme,
	// 	Scheme: k8sClient.Scheme(),
	// })
}

// Create CPUFunction CR
func createCPUFunction(ctx context.Context, cpufcr controllertestcpu.CPUFunction) error {
	tmp := &controllertestcpu.CPUFunction{}
	*tmp = cpufcr
	tmp.TypeMeta = cpufcr.TypeMeta
	err := k8sClient.Create(ctx, tmp)
	// err := k8sClient.Create(context.Background(), &CPUFunction1)
	if err != nil {
		return err
	}

	return nil
}

// Create PCIeConnection CR
func createPCIeConnection(ctx context.Context, pcieccr examplecomv1.PCIeConnection) error {
	tmp := &examplecomv1.PCIeConnection{}
	*tmp = pcieccr
	tmp.TypeMeta.Kind = "PCIeConnection"
	tmp.TypeMeta.APIVersion = "example.com/v1"
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}

	return nil
}

func createFPGAFunction(ctx context.Context, fpgafcr controllertestfpga.FPGAFunction) error {
	tmp := &controllertestfpga.FPGAFunction{}
	*tmp = fpgafcr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}

	return nil
}

// Create CPUFunction CR
func createGPUFunction(ctx context.Context, gpufcr controllertestgpu.GPUFunction) error {
	tmp := &controllertestgpu.GPUFunction{}
	*tmp = gpufcr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}

	return nil
}

/*
func createConfig(ctx context.Context, cm corev1.ConfigMap) error {
	tmp := &corev1.ConfigMap{}
	*tmp = cm
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}
*/
// Create FPGACR
func createFPGACR(ctx context.Context, fpgacr examplecomv1.FPGA) error {
	tmp := &examplecomv1.FPGA{}
	*tmp = fpgacr
	tmp.TypeMeta = fpgacr.TypeMeta
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}

	return nil
}

// To describe test cases in CPUFunctionController
var _ = Describe("PCIeConnectionController", func() {
	var mgr ctrl.Manager
	var err error
	ctx := context.Background()
	// ctx := context.WithValue(context.Background(), log.Logger)
	var fakerecorder = record.NewFakeRecorder(10)
	var writer = bytes.Buffer{}
	// var stopFunc func()

	// This test case is for reconciler
	Context("Test for PCIeConnectionReconciler", Ordered, func() {
		var reconciler PCIeConnectionReconciler

		BeforeAll(func() {

			opts := zap.Options{
				TimeEncoder: zapcore.ISO8601TimeEncoder,
				Development: true,
				DestWriter:  &writer,
			}
			ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

			// os.Setenv("PKG_CONFIG_PATH", "./fpga-software/lib/DPDK/dpdk/lib/x86_64-linux-gnu/pkgconfig:./fpga-software/lib/build/pkgconfig")
			// set manager
			mgr, err = getMgr(mgr)
			Expect(err).NotTo(HaveOccurred())

		})
		// Before Context runs, BeforeAll is executed once.
		BeforeEach(func() {

			// loger initialized
			writer.Reset()

			// recorder initialized
			fakerecorder = record.NewFakeRecorder(10)
			reconciler = PCIeConnectionReconciler{
				Client: k8sClient,
				// Client: mgr.GetClient(),
				Scheme: testScheme,
				// Scheme:   k8sClient.Scheme(),
				Recorder: fakerecorder,
			}

			err = reconciler.SetupWithManager(mgr)
			if err != nil {
				fmt.Println("Error in SetupWithManager", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// set environmental variable
			os.Setenv("K8S_NODENAME", "node01")
			os.Setenv("K8S_DPDK_LOG_FLAG", "1")
			// os.Setenv("K8S_CLUSTERNAME", "default")
			// os.Setenv("K8S_GPU_MS_PORT", "8082")
			// os.Setenv("K8S_GPU_HC_PORT", "8092")
			// os.Setenv("K8S_CPU_MS_PORT", "8083")
			// os.Setenv("K8S_CPU_HC_PORT", "8093")
			// os.Setenv("K8S_GATE_MS_PORT", "8084")
			// os.Setenv("K8S_GATE_HC_PORT", "8094")
			// os.Setenv("K8S_DI_MS_PORT", "8085")
			// os.Setenv("K8S_DI_HC_PORT", "8095")
			// os.Setenv("K8S_WBF_MS_PORT", "8086")
			// os.Setenv("K8S_WBF_HC_PORT", "8096")
			// os.Setenv("K8S_FPGA_MS_PORT", "8087")
			// os.Setenv("K8S_FPGA_HC_PORT", "8097")

			// stopFunc = startMgr(ctx, mgr)

			// To delete testdata
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.PCIeConnection{}, client.InNamespace(TESTNAMESPACE))
			if err != nil {
				fmt.Println("Can not delete PCIeConnectionCR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &controllertestfpga.FPGAFunction{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &controllertestcpu.CPUFunction{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &controllertestgpu.GPUFunction{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			// err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace(TESTNAMESPACE))
			// Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

		})

		AfterAll(func() {
			writer.Reset()
			// stop manager

			// stopFunc()
			// time.Sleep(100 * time.Millisecond)
		})

		//Test for GetFunc
		/*		It("Test_1-1-0_FPGAInit", func() {
					// Create ConfigMap of fpgalist
					err = createConfig(ctx, fpgalist)
					if err != nil {
						fmt.Println("There is a problem in createing configmap ", err)
					}

					PCIeConnectionFPGAInit(mgr)

					expected := []string{
						"/dev/xpcie_21330621T04L",
						"/dev/xpcie_21330621T01J",
						"/dev/xpcie_21330621T00Y",
						"/dev/xpcie_21330621T00D",
					}
					Expect(FpgaDevList).To(Equal(expected))

					// check logs
					Expect(writer.String()).To(ContainSubstring("fpgaDataPh3.NodeName :" + "node01"))
					Expect(writer.String()).To(ContainSubstring("fpga_shmem_controller_init() ret = "))
				})
		*/
		It("Test_1-1-1_FPGA-dma-FPGA", func() {

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection1)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}

			// Create FPGAFunctionCR decode
			err = createFPGAFunction(ctx, FPGAFunctiondecode)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR-decode ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGAFunctionCR filter/resize
			err = createFPGAFunction(ctx, FPGAFunctionfilter)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR-filter/resize ", err)
			}

			// Create FPGACR
			err = createFPGACR(ctx, FPGA1[0])
			Expect(err).NotTo(HaveOccurred())
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			err = createFPGACR(ctx, FPGA1[1])
			Expect(err).NotTo(HaveOccurred())
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			var fpgacr examplecomv1.FPGAList
			err = k8sClient.List(ctx, &fpgacr)
			fmt.Println("fpgacr--------------status[0]")
			fmt.Println(fpgacr.Items[0].Status)
			fmt.Println("fpgacr--------------status[1]")
			fmt.Println(fpgacr.Items[1].Status)

			PCIeConnectionFPGAInit(mgr)
			// Start Reconcile
			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "pcieconnectiontest1-wbconnection-decode-main-filter-resize-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())
			// Get pcieconnectionCR
			var pcieCR examplecomv1.PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest1-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)

			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status
			expectedStatus := examplecomv1.PCIeConnectionStatus{
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "pcieconnectiontest1",
					Namespace: "default",
				},
				Status: "Running",
				From: examplecomv1.PCIeFunctionStatus{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "pcieconnectiontest1-wbfunction-decode-main",
						Namespace: "default",
					},
					Status: "OK",
				},
				To: examplecomv1.PCIeFunctionStatus{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "pcieconnectiontest1-wbfunction-filter-resize-high-infer-main",
						Namespace: "default",
					},
					Status: "OK",
				},
			}

			Expect(pcieCR.Status.DataFlowRef).To(Equal(expectedStatus.DataFlowRef))
			Expect(pcieCR.Status.Status).To(Equal(expectedStatus.Status))
			Expect(pcieCR.Status.From).To(Equal(expectedStatus.From))
			Expect(pcieCR.Status.To).To(Equal(expectedStatus.To))

			// test for updating StartTime
			err = checkTime(testTime.Time, pcieCR.Status.StartTime.Time)
			Expect(err).NotTo(HaveOccurred())

			// test for creating finalizer
			expectedFinalizer := []string{
				"pcieconnection.finalizers.example.com.v1",
			}
			Expect(pcieCR.Finalizers).To(Equal(expectedFinalizer))

			// check logs
			Expect(writer.String()).To(ContainSubstring("fpga_shmem_enable() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_shmem_init() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_shmem_aligned_alloc() OK"))
			Expect(writer.String()).To(ContainSubstring("fpga_chain_connect_egress() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_chain_connect_ingress() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_buf_connect() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("Status Update."))

			// confirmation of events
			events := make([]string, 0)

			for i := 0; i < 2; i++ {
				msg := <-fakerecorder.Events
				events = append(events, msg)

			}

			Expect(events).To(ConsistOf(
				"Normal Create Create Start",
				"Normal Create Create End",
			))

		})
		It("Test_1-1-2_CPU-dma-FPGA", func() {

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection2)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR decode
			err = createCPUFunction(ctx, CPUFunctiondecode)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR-decode ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGAFunctionCR filter/resize
			err = createFPGAFunction(ctx, FPGAFunctionfilter2)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR-filter/resize ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "pcieconnectiontest2-wbconnection-decode-main-filter-resize-low-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var pcieCR examplecomv1.PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest2-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)

			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status
			expectedStatus := examplecomv1.PCIeConnectionStatus{
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "pcieconnectiontest2",
					Namespace: "default",
				},
				Status: "Running",
				From: examplecomv1.PCIeFunctionStatus{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "pcieconnectiontest2-wbfunction-decode-main",
						Namespace: "default",
					},
					Status: "OK",
				},
				To: examplecomv1.PCIeFunctionStatus{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "pcieconnectiontest2-wbfunction-filter-resize-low-infer-main",
						Namespace: "default",
					},
					Status: "OK",
				},
			}

			Expect(pcieCR.Status.DataFlowRef).To(Equal(expectedStatus.DataFlowRef))
			Expect(pcieCR.Status.Status).To(Equal(expectedStatus.Status))
			Expect(pcieCR.Status.From).To(Equal(expectedStatus.From))
			Expect(pcieCR.Status.To).To(Equal(expectedStatus.To))

			// test for updating StartTime
			err = checkTime(testTime.Time, pcieCR.Status.StartTime.Time)
			Expect(err).NotTo(HaveOccurred())

			// test for creating finalizer
			expectedFinalizer := []string{
				"pcieconnection.finalizers.example.com.v1",
			}
			Expect(pcieCR.Finalizers).To(Equal(expectedFinalizer))

			// check logs
			Expect(writer.String()).To(ContainSubstring("nothing to do"))
			Expect(writer.String()).To(ContainSubstring("debug prefix = " + CPUFunctiondecode.Spec.SharedMemory.FilePrefix))
			Expect(writer.String()).To(ContainSubstring("fpga_shmem_enable() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_chain_connect_ingress() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_lldma_init() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("Status Update."))

			// confirmation of events
			events := make([]string, 0)

			for i := 0; i < 2; i++ {
				msg := <-fakerecorder.Events
				events = append(events, msg)

			}

			Expect(events).To(ConsistOf(
				"Normal Create Create Start",
				"Normal Create Create End",
			))
		})

		It("Test_1-1-3_FPGA-dma-GPU", func() {

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection3)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGAFunctionCR filter/resize
			err = createFPGAFunction(ctx, FPGAFunctionfilter3)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR-filter/resize ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create GPUFunctionCR decode
			err = createGPUFunction(ctx, GPUFunctionhighinfer)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "pcieconnectiontest3-wbconnection-filter-resize-high-infer-main-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get pcieconnectionCR
			var pcieCR examplecomv1.PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest3-wbconnection-filter-resize-high-infer-main-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)

			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status
			expectedStatus := examplecomv1.PCIeConnectionStatus{
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "pcieconnectiontest3",
					Namespace: "default",
				},
				Status: "Running",
				From: examplecomv1.PCIeFunctionStatus{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "pcieconnectiontest3-wbfunction-filter-resize-high-infer-main",
						Namespace: "default",
					},
					Status: "OK",
				},
				To: examplecomv1.PCIeFunctionStatus{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "pcieconnectiontest3-wbfunction-high-infer-main",
						Namespace: "default",
					},
					Status: "OK",
				},
			}

			Expect(pcieCR.Status.DataFlowRef).To(Equal(expectedStatus.DataFlowRef))
			Expect(pcieCR.Status.Status).To(Equal(expectedStatus.Status))
			Expect(pcieCR.Status.From).To(Equal(expectedStatus.From))
			Expect(pcieCR.Status.To).To(Equal(expectedStatus.To))

			// test for updating StartTime
			err = checkTime(testTime.Time, pcieCR.Status.StartTime.Time)
			Expect(err).NotTo(HaveOccurred())

			// test for creating finalizer
			expectedFinalizer := []string{
				"pcieconnection.finalizers.example.com.v1",
			}
			Expect(pcieCR.Finalizers).To(Equal(expectedFinalizer))

			// check logs
			Expect(writer.String()).To(ContainSubstring("debug prefix = " + GPUFunctionhighinfer.Spec.SharedMemory.FilePrefix))
			Expect(writer.String()).To(ContainSubstring("fpga_shmem_enable() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_chain_connect_egress() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_lldma_init() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("nothing to do"))
			Expect(writer.String()).To(ContainSubstring("Status Update."))

			// confirmation of events
			events := make([]string, 0)

			for i := 0; i < 2; i++ {
				msg := <-fakerecorder.Events
				events = append(events, msg)

			}

			Expect(events).To(ConsistOf(
				"Normal Create Create Start",
				"Normal Create Create End",
			))
		})

		It("Test_1-1-4_cpu-dma-cpu", func() {

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection4)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR decode
			err = createCPUFunction(ctx, CPUFunctionDecode)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR filter/resize
			err = createCPUFunction(ctx, CPUFunctionFilterResize)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "pcieconnectiontest4-wbconnection-decode-main-filter-resize-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get pcieconnectionCR
			var pcieCR examplecomv1.PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest4-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)

			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status
			expectedStatus := examplecomv1.PCIeConnectionStatus{
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "pcieconnectiontest4",
					Namespace: "default",
				},
				Status: "Running",
				From: examplecomv1.PCIeFunctionStatus{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "pcieconnectiontest4-wbfunction-decode-main",
						Namespace: "default",
					},
					Status: "OK",
				},
				To: examplecomv1.PCIeFunctionStatus{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "pcieconnectiontest4-wbfunction-filter-resize-high-infer-main",
						Namespace: "default",
					},
					Status: "OK",
				},
			}

			Expect(pcieCR.Status.DataFlowRef).To(Equal(expectedStatus.DataFlowRef))
			Expect(pcieCR.Status.Status).To(Equal(expectedStatus.Status))
			Expect(pcieCR.Status.From).To(Equal(expectedStatus.From))
			Expect(pcieCR.Status.To).To(Equal(expectedStatus.To))

			// test for updating StartTime
			err = checkTime(testTime.Time, pcieCR.Status.StartTime.Time)
			Expect(err).NotTo(HaveOccurred())

			// test for creating finalizer
			expectedFinalizer := []string{
				"pcieconnection.finalizers.example.com.v1",
			}
			Expect(pcieCR.Finalizers).To(Equal(expectedFinalizer))

			fmt.Println("test : " + writer.String())

			// check logs
			Expect(writer.String()).To(ContainSubstring("debug prefix = " + CPUFunctionDecode.Spec.SharedMemory.FilePrefix))
			Expect(writer.String()).To(ContainSubstring("fpga_shmem_enable() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("nothing to do"))
			Expect(writer.String()).To(ContainSubstring("Status Update."))

			// confirmation of events
			events := make([]string, 0)

			for i := 0; i < 2; i++ {
				msg := <-fakerecorder.Events
				events = append(events, msg)

			}

			Expect(events).To(ConsistOf(
				"Normal Create Create Start",
				"Normal Create Create End",
			))
		})

		It("Test_2-1_UPDATE", func() {

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, pcieconnectiontestUPDATE)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "pcieconnectiontestupdate-wbconnection-decode-main-filter-resize-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())
			// confirmation of events
			events := make([]string, 0)

			for i := 0; i < 2; i++ {
				msg := <-fakerecorder.Events
				events = append(events, msg)

			}

			Expect(events).To(ConsistOf(
				"Normal Update Update Start",
				"Normal Update Update End",
			))
		})
		It("Test_2-2_DELETE", func() {

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, pcieconnectiontestDELETE)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			var pcieCR examplecomv1.PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontestdelete-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)

			err = k8sClient.Delete(ctx, &pcieCR)

			var pciebCR examplecomv1.PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontestdelete-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pciebCR)

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "pcieconnectiontestdelete-wbconnection-decode-main-filter-resize-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// confirmation of events
			events := make([]string, 0)

			for i := 0; i < 2; i++ {
				msg := <-fakerecorder.Events
				events = append(events, msg)

			}

			Expect(events).To(ConsistOf(
				"Normal Delete Delete Start",
				"Normal Delete Delete End",
			))

			var pcieaCR examplecomv1.PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontestdelete-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieaCR)

			Expect(err).To(MatchError(ContainSubstring("not found")))

		})
	})
})

// func startMgr(ctx context.Context, mgr manager.Manager) func() {
// 	ctx, cancel := context.WithCancel(ctx)
// 	go func() {
// 		err := mgr.Start(ctx)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}()
// 	time.Sleep(100 * time.Millisecond)
// 	return cancel
// }

// Expects that t1 is before t2
func checkTime(t1 time.Time, t2 time.Time) error {
	if t1.Before(t2) {
		return nil
	} else {
		return fmt.Errorf("CR.Status may not be updated. t1 (%s) is after t2 (%s)", t1.Format(time.RFC3339), t2.Format(time.RFC3339))
	}

}

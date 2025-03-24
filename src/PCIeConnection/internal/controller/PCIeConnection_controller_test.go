/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED
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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	examplecomv1 "PCIeConnection/api/v1"
	controllertestcpu "PCIeConnection/internal/controller/test/type/CPU"
	controllertestfpga "PCIeConnection/internal/controller/test/type/FPGA"
	controllertestgpu "PCIeConnection/internal/controller/test/type/GPU"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func getMgr(mgr ctrl.Manager) (ctrl.Manager, error) {
	if mgr == nil {
		return ctrl.NewManager(cfg, ctrl.Options{
			Scheme: k8sClient.Scheme(),
		})
	}
	return mgr, nil
}

// Create CPUFunction CR
func createCPUFunction(ctx context.Context, cpufcr controllertestcpu.CPUFunction) error {
	tmp := &controllertestcpu.CPUFunction{}
	*tmp = cpufcr
	tmp.TypeMeta = cpufcr.TypeMeta
	err := k8sClient.Create(ctx, tmp)
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

// Create GPUFunction CR
func createGPUFunction(ctx context.Context, gpufcr controllertestgpu.GPUFunction) error {
	tmp := &controllertestgpu.GPUFunction{}
	*tmp = gpufcr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}

	return nil
}

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

// Create CPUPod
func createCPUPod(ctx context.Context, cpuPod corev1.Pod) error {
	tmp := &corev1.Pod{}
	*tmp = cpuPod
	tmp.TypeMeta = cpuPod.TypeMeta
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}

	return nil
}

// Create GPUPod
func createGPUPod(ctx context.Context, gpuPod corev1.Pod) error {
	tmp := &corev1.Pod{}
	*tmp = gpuPod
	tmp.TypeMeta = gpuPod.TypeMeta
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}

	return nil
}

// To describe test cases in PCIeConnectionController
var _ = Describe("PCIeConnectionController", func() {
	var mgr ctrl.Manager
	var err error
	ctx := context.Background()
	var fakerecorder = record.NewFakeRecorder(10)
	var writer = bytes.Buffer{}

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
				Client:   k8sClient,
				Scheme:   k8sClient.Scheme(),
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

			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGA{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

		})

		AfterAll(func() {
			writer.Reset()
		})

		/* D2D is not supported
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
		*/
		It("Test_1-1-2_CPU-dma-FPGA", func() {

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection2)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR decode
			err = createCPUFunction(ctx, CPUFunctionDecode2)
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
			Expect(writer.String()).To(ContainSubstring("debug prefix = " + CPUFunctionDecode2.Spec.SharedMemory.FilePrefix))
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

			controllerutil.RemoveFinalizer(&pcieCR, "pcieconnection.finalizers.example.com.v1")
			// Update RemoveFinalizer to fpgafunc CR
			err = k8sClient.Update(ctx, &pcieCR)
			if err != nil {
				fmt.Println("error update RemoveFinalizer to pcieCR.", err)
			}
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

			controllerutil.RemoveFinalizer(&pcieCR, "pcieconnection.finalizers.example.com.v1")
			// Update RemoveFinalizer to fpgafunc CR
			err = k8sClient.Update(ctx, &pcieCR)
			if err != nil {
				fmt.Println("error update RemoveFinalizer to pcieCR.", err)
			}
		})

		It("Test_1-1-4_cpu-dma-cpu", func() {

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection4)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR decode
			err = createCPUFunction(ctx, CPUFunctionDecode4)
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
			Expect(writer.String()).To(ContainSubstring("debug prefix = " + CPUFunctionDecode4.Spec.SharedMemory.FilePrefix))
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

		It("Test_3-1-6_Terminating", func() {

			// Create CPUFunctionCR decode
			err = createCPUFunction(ctx, CPUFunctionDecode2)
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

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection5)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Get pcieconnectionCR
			var pcieCR examplecomv1.PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest5-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)

			// Delete pcieconnectionCR
			err = k8sClient.Delete(ctx, &pcieCR)

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "pcieconnectiontest5-wbconnection-decode-main-filter-resize-low-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))
			Expect(err).NotTo(HaveOccurred())

			// Get pcieconnectionCR
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest5-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)
			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status
			expectedStatus := examplecomv1.PCIeConnectionStatus{
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "pcieconnectiontest5",
					Namespace: "default",
				},
				Status: "Terminating",
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

			// check logs
			Expect(writer.String()).To(ContainSubstring("Status Update."))
			Expect(writer.String()).To(ContainSubstring("Update PCIeConnection to Terminating."))
		})

		/* D2D is not supported
		It("Test_3-1-2_FPGA-dma-FPGA", func() {

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection6)
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
			Expect(err).NotTo(HaveOccurred())

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

			// Get pcieconnectionCR
			var pcieCR examplecomv1.PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest6-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)

			// Delete pcieconnectionCR
			err = k8sClient.Delete(ctx, &pcieCR)

			// Start Reconcile
			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "pcieconnectiontest6-wbconnection-decode-main-filter-resize-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get pcieconnectionCR
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest6-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)
			Expect(err).To(HaveOccurred())

			// check logs
			Expect(writer.String()).To(ContainSubstring("nothing to do"))
			Expect(writer.String()).To(ContainSubstring("Update PCIeConnection to Released."))

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

		})
		*/

		It("Test_3-1-3_CPU-dma-FPGA", func() {

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection7)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR decode
			err = createCPUFunction(ctx, CPUFunctionDecode2)
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

			// Create CPUPodCR
			err = createCPUPod(ctx, CPUPod1)
			if err != nil {
				fmt.Println("There is a problem in createing CPUPod ", err)
			}
			Expect(err).NotTo(HaveOccurred())

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

			// Get pcieconnectionCR
			var pcieCR examplecomv1.PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest7-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)
			Expect(err).NotTo(HaveOccurred())

			// Delete pcieconnectionCR
			err = k8sClient.Delete(ctx, &pcieCR)

			ShmemEnable[CPUFunctionDecode2.Spec.SharedMemory.FilePrefix] = true

			// Start Reconcile
			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "pcieconnectiontest7-wbconnection-decode-main-filter-resize-low-infer-main",
					Namespace: TESTNAMESPACE,
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))
			Expect(err).NotTo(HaveOccurred())

			// Get pcieconnectionCR
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest7-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)
			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status
			expectedStatus := examplecomv1.PCIeConnectionStatus{
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "pcieconnectiontest7",
					Namespace: "default",
				},
				Status: "Terminating",
				From: examplecomv1.PCIeFunctionStatus{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "pcieconnectiontest2-wbfunction-decode-main",
						Namespace: "default",
					},
					Status: "PODDELETING",
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

			// check logs
			Expect(writer.String()).To(ContainSubstring("PodCR Deletion Start."))
			Expect(writer.String()).To(ContainSubstring("CR Status is not Released"))

			// confirmation of events
			events := make([]string, 0)

			for i := 0; i < 1; i++ {
				msg := <-fakerecorder.Events
				events = append(events, msg)

			}

			Expect(events).To(ConsistOf(
				"Normal Delete Delete Start",
			))

			// Start Reconcile
			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "pcieconnectiontest7-wbconnection-decode-main-filter-resize-low-infer-main",
					Namespace: TESTNAMESPACE,
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: true, RequeueAfter: 0}))
			Expect(err).NotTo(HaveOccurred())

			// check logs
			Expect(writer.String()).To(ContainSubstring("Deleting PodCR."))
			Expect(writer.String()).To(ContainSubstring("CR Status is not Released"))

			// confirmation of events
			events = make([]string, 0)

			for i := 0; i < 1; i++ {
				msg := <-fakerecorder.Events
				events = append(events, msg)

			}

			Expect(events).To(ConsistOf(
				"Normal Delete Delete Start",
			))

			// Get podCR
			var podCR corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest1-wbfunction-decode-main-cpu-pod",
				Namespace: TESTNAMESPACE,
			},
				&podCR)
			Expect(err).NotTo(HaveOccurred())

			// Delete podCR
			controllerutil.RemoveFinalizer(&podCR, "kubernetes")
			err = k8sClient.Update(ctx, &podCR)
			Expect(err).NotTo(HaveOccurred())

			// Start Reconcile
			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "pcieconnectiontest7-wbconnection-decode-main-filter-resize-low-infer-main",
					Namespace: TESTNAMESPACE,
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))
			Expect(err).NotTo(HaveOccurred())
			_, exist := ShmemEnable[CPUFunctionDecode2.Status.SharedMemory.FilePrefix]
			Expect(exist).To(BeFalse())

			// check logs
			Expect(writer.String()).To(ContainSubstring("Success to delete PodCR."))
			Expect(writer.String()).To(ContainSubstring("nothing to do"))
			Expect(writer.String()).To(ContainSubstring("fpga_chain_disconnect_ingress() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_chain_wait_stat_egr_free() OK is_success = 1"))
			Expect(writer.String()).To(ContainSubstring("fpga_lldma_finish() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_shmem_disable_with_check() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("Update PCIeConnection to Released."))

			// confirmation of events
			events = make([]string, 0)

			for i := 0; i < 2; i++ {
				msg := <-fakerecorder.Events
				events = append(events, msg)

			}

			Expect(events).To(ConsistOf(
				"Normal Delete Delete Start",
				"Normal Delete Delete End",
			))
		})

		It("Test_3-1-4_FPGA-dma-GPU", func() {

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection8)
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

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunctionhighinfer)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create GPUPodCR
			err = createGPUPod(ctx, GPUPod1)
			if err != nil {
				fmt.Println("There is a problem in createing CPUPod ", err)
			}
			Expect(err).NotTo(HaveOccurred())

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

			// Get pcieconnectionCR
			var pcieCR examplecomv1.PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest8-wbconnection-filter-resize-high-infer-main-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)

			// Delete pcieconnectionCR
			err = k8sClient.Delete(ctx, &pcieCR)

			ShmemEnable[GPUFunctionhighinfer.Spec.SharedMemory.FilePrefix] = true

			// Start Reconcile
			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "pcieconnectiontest8-wbconnection-filter-resize-high-infer-main-high-infer-main",
					Namespace: TESTNAMESPACE,
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))
			Expect(err).NotTo(HaveOccurred())

			// Get pcieconnectionCR
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest8-wbconnection-filter-resize-high-infer-main-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)
			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status
			expectedStatus := examplecomv1.PCIeConnectionStatus{
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "pcieconnectiontest8",
					Namespace: "default",
				},
				Status: "Terminating",
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
					Status: "PODDELETING",
				},
			}

			Expect(pcieCR.Status.DataFlowRef).To(Equal(expectedStatus.DataFlowRef))
			Expect(pcieCR.Status.Status).To(Equal(expectedStatus.Status))
			Expect(pcieCR.Status.From).To(Equal(expectedStatus.From))
			Expect(pcieCR.Status.To).To(Equal(expectedStatus.To))

			// check logs
			Expect(writer.String()).To(ContainSubstring("PodCR Deletion Start."))
			Expect(writer.String()).To(ContainSubstring("CR Status is not Released"))

			// confirmation of events
			events := make([]string, 0)

			for i := 0; i < 1; i++ {
				msg := <-fakerecorder.Events
				events = append(events, msg)

			}

			Expect(events).To(ConsistOf(
				"Normal Delete Delete Start",
			))

			// Start Reconcile
			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "pcieconnectiontest8-wbconnection-filter-resize-high-infer-main-high-infer-main",
					Namespace: TESTNAMESPACE,
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: true, RequeueAfter: 0}))
			Expect(err).NotTo(HaveOccurred())

			// check logs
			Expect(writer.String()).To(ContainSubstring("Deleting PodCR."))
			Expect(writer.String()).To(ContainSubstring("CR Status is not Released"))

			// confirmation of events
			events = make([]string, 0)

			for i := 0; i < 1; i++ {
				msg := <-fakerecorder.Events
				events = append(events, msg)

			}

			Expect(events).To(ConsistOf(
				"Normal Delete Delete Start",
			))

			// Get podCR
			var podCR corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest1-wbfunction-high-infer-main-gpu-pod",
				Namespace: TESTNAMESPACE,
			},
				&podCR)
			Expect(err).NotTo(HaveOccurred())

			// Delete podCR
			controllerutil.RemoveFinalizer(&podCR, "kubernetes")
			err = k8sClient.Update(ctx, &podCR)
			Expect(err).NotTo(HaveOccurred())

			// Start Reconcile
			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "pcieconnectiontest8-wbconnection-filter-resize-high-infer-main-high-infer-main",
					Namespace: TESTNAMESPACE,
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))
			Expect(err).NotTo(HaveOccurred())

			// check logs
			Expect(writer.String()).To(ContainSubstring("Success to delete PodCR."))
			Expect(writer.String()).To(ContainSubstring("CR Status is not Released"))

			// confirmation of events
			events = make([]string, 0)

			for i := 0; i < 1; i++ {
				msg := <-fakerecorder.Events
				events = append(events, msg)

			}

			Expect(events).To(ConsistOf(
				"Normal Delete Delete Start",
			))

			// Start Reconcile
			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "pcieconnectiontest8-wbconnection-filter-resize-high-infer-main-high-infer-main",
					Namespace: TESTNAMESPACE,
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))
			Expect(err).NotTo(HaveOccurred())
			_, exist := ShmemEnable[GPUFunctionhighinfer.Status.SharedMemory.FilePrefix]
			Expect(exist).To(BeFalse())

			// Get pcieconnectionCR
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest8-wbconnection-filter-resize-high-infer-main-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)

			Expect(err).To(HaveOccurred())

			// check logs
			Expect(writer.String()).To(ContainSubstring("nothing to do"))
			Expect(writer.String()).To(ContainSubstring("fpga_chain_wait_disconnection_ingress() OK is_success = 1"))
			Expect(writer.String()).To(ContainSubstring("fpga_chain_wait_stat_egr_free() OK is_success = 1"))
			Expect(writer.String()).To(ContainSubstring("fpga_chain_disconnect_egress() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_lldma_finish() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_shmem_disable_with_check() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("nothing to do"))
			Expect(writer.String()).To(ContainSubstring("Update PCIeConnection to Released."))

			// confirmation of events
			events = make([]string, 0)

			for i := 0; i < 2; i++ {
				msg := <-fakerecorder.Events
				events = append(events, msg)
			}

			Expect(events).To(ConsistOf(
				"Normal Delete Delete Start",
				"Normal Delete Delete End",
			))
		})

		It("Test_3-1-5_CPU-dma-CPU", func() {

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection9)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR decode
			err = createCPUFunction(ctx, CPUFunctionDecode4)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR-decode", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR filter/resize
			err = createCPUFunction(ctx, CPUFunctionFilterResize)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR-filter/resize", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUPodCR
			err = createCPUPod(ctx, CPUPod1)
			if err != nil {
				fmt.Println("There is a problem in createing CPUPod ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			err = createCPUPod(ctx, CPUPod2)
			if err != nil {
				fmt.Println("There is a problem in createing CPUPod ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Get pcieconnectionCR
			var pcieCR examplecomv1.PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest9-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)

			// Delete pcieconnectionCR
			err = k8sClient.Delete(ctx, &pcieCR)

			ShmemEnable[CPUFunctionDecode4.Spec.SharedMemory.FilePrefix] = true

			// Start Reconcile
			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "pcieconnectiontest9-wbconnection-decode-main-filter-resize-high-infer-main",
					Namespace: TESTNAMESPACE,
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))
			Expect(err).NotTo(HaveOccurred())

			// Get pcieconnectionCR
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest9-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)

			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status
			expectedStatus := examplecomv1.PCIeConnectionStatus{
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "pcieconnectiontest9",
					Namespace: "default",
				},
				Status: "Terminating",
				From: examplecomv1.PCIeFunctionStatus{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "pcieconnectiontest4-wbfunction-decode-main",
						Namespace: "default",
					},
					Status: "PODDELETING",
				},
				To: examplecomv1.PCIeFunctionStatus{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "pcieconnectiontest4-wbfunction-filter-resize-high-infer-main",
						Namespace: "default",
					},
					Status: "PODDELETING",
				},
			}

			Expect(pcieCR.Status.DataFlowRef).To(Equal(expectedStatus.DataFlowRef))
			Expect(pcieCR.Status.Status).To(Equal(expectedStatus.Status))
			Expect(pcieCR.Status.From).To(Equal(expectedStatus.From))
			Expect(pcieCR.Status.To).To(Equal(expectedStatus.To))

			// check logs
			Expect(writer.String()).To(ContainSubstring("PodCR Deletion Start."))
			Expect(writer.String()).To(ContainSubstring("PodCR Deletion Start."))
			Expect(writer.String()).To(ContainSubstring("CR Status is not Released"))

			// confirmation of events
			events := make([]string, 0)

			for i := 0; i < 1; i++ {
				msg := <-fakerecorder.Events
				events = append(events, msg)

			}

			Expect(events).To(ConsistOf(
				"Normal Delete Delete Start",
			))

			// Start Reconcile
			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "pcieconnectiontest9-wbconnection-decode-main-filter-resize-high-infer-main",
					Namespace: TESTNAMESPACE,
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: true, RequeueAfter: 0}))
			Expect(err).NotTo(HaveOccurred())

			// check logs
			Expect(writer.String()).To(ContainSubstring("Deleting PodCR."))
			Expect(writer.String()).To(ContainSubstring("Deleting PodCR."))

			// confirmation of events
			events = make([]string, 0)

			for i := 0; i < 1; i++ {
				msg := <-fakerecorder.Events
				events = append(events, msg)

			}

			Expect(events).To(ConsistOf(
				"Normal Delete Delete Start",
			))

			// Get podCR
			var podCR1 corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest1-wbfunction-decode-main-cpu-pod",
				Namespace: TESTNAMESPACE,
			},
				&podCR1)
			Expect(err).NotTo(HaveOccurred())

			var podCR2 corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest2-wbfunction-decode-main-cpu-pod",
				Namespace: TESTNAMESPACE,
			},
				&podCR2)
			Expect(err).NotTo(HaveOccurred())

			// Delete podCR
			controllerutil.RemoveFinalizer(&podCR1, "kubernetes")
			err = k8sClient.Update(ctx, &podCR1)
			Expect(err).NotTo(HaveOccurred())

			controllerutil.RemoveFinalizer(&podCR2, "kubernetes")
			err = k8sClient.Update(ctx, &podCR2)
			Expect(err).NotTo(HaveOccurred())

			// Start Reconcile
			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "pcieconnectiontest9-wbconnection-decode-main-filter-resize-high-infer-main",
					Namespace: TESTNAMESPACE,
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))
			Expect(err).NotTo(HaveOccurred())

			// check logs
			Expect(writer.String()).To(ContainSubstring("Success to delete PodCR."))
			Expect(writer.String()).To(ContainSubstring("Success to delete PodCR."))
			Expect(writer.String()).To(ContainSubstring("CR Status is not Released"))

			// confirmation of events
			events = make([]string, 0)

			for i := 0; i < 1; i++ {
				msg := <-fakerecorder.Events
				events = append(events, msg)

			}

			Expect(events).To(ConsistOf(
				"Normal Delete Delete Start",
			))

			// Start Reconcile
			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "pcieconnectiontest9-wbconnection-decode-main-filter-resize-high-infer-main",
					Namespace: TESTNAMESPACE,
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))
			Expect(err).NotTo(HaveOccurred())
			_, exist := ShmemEnable[CPUFunctionDecode4.Status.SharedMemory.FilePrefix]
			Expect(exist).To(BeFalse())

			// check logs
			Expect(writer.String()).To(ContainSubstring("Update PCIeConnection to Released."))

			// confirmation of events
			events = make([]string, 0)

			for i := 0; i < 2; i++ {
				msg := <-fakerecorder.Events
				events = append(events, msg)

			}

			Expect(events).To(ConsistOf(
				"Normal Delete Delete Start",
				"Normal Delete Delete End",
			))
		})

		It("Test_7-2-1_FPGA-dma-GPU", func() {
			// generated on the basis of Test_1-1-3_FPGA-dma-GPU

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection3)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGAFunctionCR filter/resize
			crData := &unstructured.Unstructured{}
			crData.SetGroupVersionKind(schema.GroupVersionKind{
				Version: FPGAFunctionfilter3.TypeMeta.APIVersion,
				Kind:    FPGAFunctionfilter3.TypeMeta.Kind,
			})
			crData.SetName(FPGAFunctionfilter3.ObjectMeta.Name)
			crData.SetNamespace(FPGAFunctionfilter3.ObjectMeta.Namespace)
			crData.UnstructuredContent()["spec"] = FPGAFunctionfilter3.Spec
			err = k8sClient.Create(ctx, crData)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR-filter/resize ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunctionhighinfer)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			PCIeConnectionFPGAInit(mgr)

			// Start Reconcile
			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "pcieconnectiontest3-wbconnection-filter-resize-high-infer-main-high-infer-main",
					Namespace: TESTNAMESPACE,
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: true, RequeueAfter: 0}))
			Expect(err).NotTo(HaveOccurred())

			var pcieCR examplecomv1.PCIeConnection
			// Get pcieconnectionCR
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest3-wbconnection-filter-resize-high-infer-main-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)
			Expect(err).NotTo(HaveOccurred())

			Expect(pcieCR.Status.DataFlowRef).To(Equal(PCIeConnection3.Status.DataFlowRef))
			Expect(pcieCR.Status.Status).To(Equal(PCIeConnection3.Status.Status))
			Expect(pcieCR.Status.From).To(Equal(PCIeConnection3.Status.From))
			Expect(pcieCR.Status.To).To(Equal(PCIeConnection3.Status.To))

			// test for updating StartTime
			err = checkTime(testTime.Time, pcieCR.Status.StartTime.Time)
			Expect(err).To(HaveOccurred())

			Expect(pcieCR.Finalizers).Should(BeNil())

			// check logs
			Expect(writer.String()).To(ContainSubstring("Requeue because source FPGA FunctionKernelID is not determined."))

			// Start Reconcile
			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "pcieconnectiontest3-wbconnection-filter-resize-high-infer-main-high-infer-main",
					Namespace: TESTNAMESPACE,
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: true, RequeueAfter: 0}))
			Expect(err).NotTo(HaveOccurred())

			// Get pcieconnectionCR
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest3-wbconnection-filter-resize-high-infer-main-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)
			Expect(err).NotTo(HaveOccurred())

			Expect(pcieCR.Status.DataFlowRef).To(Equal(PCIeConnection3.Status.DataFlowRef))
			Expect(pcieCR.Status.Status).To(Equal(PCIeConnection3.Status.Status))
			Expect(pcieCR.Status.From).To(Equal(PCIeConnection3.Status.From))
			Expect(pcieCR.Status.To).To(Equal(PCIeConnection3.Status.To))

			// test for updating StartTime
			err = checkTime(testTime.Time, pcieCR.Status.StartTime.Time)
			Expect(err).To(HaveOccurred())

			Expect(pcieCR.Finalizers).Should(BeNil())

			// check logs
			Expect(writer.String()).To(ContainSubstring("Requeue because source FPGA FunctionKernelID is not determined."))

			// Start Reconcile
			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "pcieconnectiontest3-wbconnection-filter-resize-high-infer-main-high-infer-main",
					Namespace: TESTNAMESPACE,
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: true, RequeueAfter: 0}))
			Expect(err).NotTo(HaveOccurred())

			// Get pcieconnectionCR
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest3-wbconnection-filter-resize-high-infer-main-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)
			Expect(err).NotTo(HaveOccurred())

			Expect(pcieCR.Status.DataFlowRef).To(Equal(PCIeConnection3.Status.DataFlowRef))
			Expect(pcieCR.Status.Status).To(Equal(PCIeConnection3.Status.Status))
			Expect(pcieCR.Status.From).To(Equal(PCIeConnection3.Status.From))
			Expect(pcieCR.Status.To).To(Equal(PCIeConnection3.Status.To))

			// test for updating StartTime
			err = checkTime(testTime.Time, pcieCR.Status.StartTime.Time)
			Expect(err).To(HaveOccurred())

			Expect(pcieCR.Finalizers).Should(BeNil())

			// check logs
			Expect(writer.String()).To(ContainSubstring("Requeue because source FPGA FunctionKernelID is not determined."))
		})
		It("Test_7-2-2_CPU-dma-FPGA", func() {
			// generated on the basis of Test_1-1-2_FPGA-dma-GPU

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection2)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR decode
			err = createCPUFunction(ctx, CPUFunctionDecode2)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR-decode ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGAFunctionCR filter/resize
			crData := &unstructured.Unstructured{}
			crData.SetGroupVersionKind(schema.GroupVersionKind{
				Version: FPGAFunctionfilter2.TypeMeta.APIVersion,
				Kind:    FPGAFunctionfilter2.TypeMeta.Kind,
			})
			crData.SetName(FPGAFunctionfilter2.ObjectMeta.Name)
			crData.SetNamespace(FPGAFunctionfilter2.ObjectMeta.Namespace)
			crData.UnstructuredContent()["spec"] = FPGAFunctionfilter2.Spec
			err = k8sClient.Create(ctx, crData)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR-filter/resize ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			PCIeConnectionFPGAInit(mgr)

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "pcieconnectiontest2-wbconnection-decode-main-filter-resize-low-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: true, RequeueAfter: 0}))
			Expect(err).NotTo(HaveOccurred())

			var pcieCR examplecomv1.PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest2-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)

			Expect(err).NotTo(HaveOccurred())

			Expect(pcieCR.Status.DataFlowRef).To(Equal(PCIeConnection2.Status.DataFlowRef))
			Expect(pcieCR.Status.Status).To(Equal(PCIeConnection2.Status.Status))
			Expect(pcieCR.Status.From).To(Equal(PCIeConnection2.Status.From))
			Expect(pcieCR.Status.To).To(Equal(PCIeConnection2.Status.To))

			// test for updating StartTime
			err = checkTime(testTime.Time, pcieCR.Status.StartTime.Time)
			Expect(err).To(HaveOccurred())

			Expect(pcieCR.Finalizers).Should(BeNil())

			// check logs
			Expect(writer.String()).To(ContainSubstring("Requeue because destination FPGA FunctionKernelID is not determined."))

			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "pcieconnectiontest2-wbconnection-decode-main-filter-resize-low-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: true, RequeueAfter: 0}))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest2-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)

			Expect(err).NotTo(HaveOccurred())

			Expect(pcieCR.Status.DataFlowRef).To(Equal(PCIeConnection2.Status.DataFlowRef))
			Expect(pcieCR.Status.Status).To(Equal(PCIeConnection2.Status.Status))
			Expect(pcieCR.Status.From).To(Equal(PCIeConnection2.Status.From))
			Expect(pcieCR.Status.To).To(Equal(PCIeConnection2.Status.To))

			// test for updating StartTime
			err = checkTime(testTime.Time, pcieCR.Status.StartTime.Time)
			Expect(err).To(HaveOccurred())

			Expect(pcieCR.Finalizers).Should(BeNil())

			// check logs
			Expect(writer.String()).To(ContainSubstring("Requeue because destination FPGA FunctionKernelID is not determined."))

			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "pcieconnectiontest2-wbconnection-decode-main-filter-resize-low-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: true, RequeueAfter: 0}))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest2-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)

			Expect(err).NotTo(HaveOccurred())

			Expect(pcieCR.Status.DataFlowRef).To(Equal(PCIeConnection2.Status.DataFlowRef))
			Expect(pcieCR.Status.Status).To(Equal(PCIeConnection2.Status.Status))
			Expect(pcieCR.Status.From).To(Equal(PCIeConnection2.Status.From))
			Expect(pcieCR.Status.To).To(Equal(PCIeConnection2.Status.To))

			// test for updating StartTime
			err = checkTime(testTime.Time, pcieCR.Status.StartTime.Time)
			Expect(err).To(HaveOccurred())

			Expect(pcieCR.Finalizers).Should(BeNil())

			// check logs
			Expect(writer.String()).To(ContainSubstring("Requeue because destination FPGA FunctionKernelID is not determined."))

		})
		It("Test_7-2-3_Status Not Change route", func() {

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection723)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR decode
			err = createCPUFunction(ctx, CPUFunctionDecode4)
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
					Name:      "pcieconnectiontest723-wbconnection-decode-main-filter-resize-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: true}))
			Expect(err).NotTo(HaveOccurred())

			// Get pcieconnectionCR
			var pcieCR examplecomv1.PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "pcieconnectiontest723-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)

			Expect(err).NotTo(HaveOccurred())

			Expect(pcieCR.Status.DataFlowRef).To(Equal(PCIeConnection723.Status.DataFlowRef))
			Expect(pcieCR.Status.Status).To(Equal(PCIeConnection723.Status.Status))
			Expect(pcieCR.Status.From).To(Equal(PCIeConnection723.Status.From))
			Expect(pcieCR.Status.To).To(Equal(PCIeConnection723.Status.To))

			// test for updating StartTime
			err = checkTime(testTime.Time, pcieCR.Status.StartTime.Time)
			Expect(err).To(HaveOccurred())

			// test for creating finalizer
			var nilfzr []string
			nilfzr = nil
			Expect(pcieCR.Finalizers).To(Equal(nilfzr))

			fmt.Println("test : " + writer.String())

			// check logs
			Expect(writer.String()).To(ContainSubstring("debug prefix = " + CPUFunctionDecode4.Spec.SharedMemory.FilePrefix))
			Expect(writer.String()).To(ContainSubstring("CR Status is not Changed"))

			// confirmation of events
			events := make([]string, 0)

			for i := 0; i < 1; i++ {
				msg := <-fakerecorder.Events
				events = append(events, msg)
			}

			Expect(events).To(ConsistOf(
				"Normal Create Create Start",
			))
		})
	})
})

// Expects that t1 is before t2
func checkTime(t1 time.Time, t2 time.Time) error {
	if t1.Before(t2) {
		return nil
	} else {
		return fmt.Errorf("CR.Status may not be updated. t1 (%s) is after t2 (%s)", t1.Format(time.RFC3339), t2.Format(time.RFC3339))
	}

}

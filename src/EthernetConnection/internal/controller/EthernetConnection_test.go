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
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	examplecomv1 "EthernetConnection/api/v1"
	controllertestcpu "EthernetConnection/internal/controller/test/type/CPU"
	controllertestfpga "EthernetConnection/internal/controller/test/type/FPGA"
	controllertestgpu "EthernetConnection/internal/controller/test/type/GPU"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"

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

// Create Pod CR
func createPod(ctx context.Context, podcr corev1.Pod) error {
	tmp := &corev1.Pod{}
	*tmp = podcr

	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}

	return nil
}

// Create EthernetConnection CR
func createEthernetConnection(ctx context.Context, ethernetccr examplecomv1.EthernetConnection) error {
	tmp := &examplecomv1.EthernetConnection{}
	*tmp = ethernetccr
	tmp.TypeMeta.Kind = "EthernetConnection"
	tmp.TypeMeta.APIVersion = "example.com/v1"
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}

	return nil
}

// Create FPGAFunction CR
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

// Create ChildbsCR
func createChildbs(ctx context.Context, childbscr examplecomv1.ChildBs) error {
	tmp := &examplecomv1.ChildBs{}
	*tmp = childbscr
	tmp.TypeMeta = childbscr.TypeMeta
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}

	return nil
}

// Update ChildbsCR
func updateChildbs(ctx context.Context, childbscr examplecomv1.ChildBs) error {
	tmp := &examplecomv1.ChildBs{}
	*tmp = childbscr
	tmp.TypeMeta = childbscr.TypeMeta
	err := k8sClient.Update(ctx, tmp)
	if err != nil {
		return err
	}

	return nil
}

// Delete ChildbsCR
func deleteChildbs(ctx context.Context, childbscr examplecomv1.ChildBs) error {
	tmp := &examplecomv1.ChildBs{}
	*tmp = childbscr
	tmp.TypeMeta = childbscr.TypeMeta
	err := k8sClient.Delete(ctx, tmp)
	if err != nil {
		return err
	}

	return nil
}

// To describe test cases in CPUFunctionController
var _ = Describe("EthernetConnectionController", func() {
	var mgr ctrl.Manager
	var err error
	ctx := context.Background()
	var fakerecorder = record.NewFakeRecorder(10)
	var writer = bytes.Buffer{}

	// This test case is for reconciler
	Context("Test for EthernetConnectionReconciler", Ordered, func() {
		var reconciler EthernetConnectionReconciler

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
			reconciler = EthernetConnectionReconciler{
				Client:   k8sClient,
				Scheme:   testScheme,
				Recorder: fakerecorder,
			}

			err = reconciler.SetupWithManager(mgr)
			if err != nil {
				fmt.Println("setupwithManagerでerror発生", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// set environmental variable
			os.Setenv("K8S_NODENAME", "node01")

			// To delete testdata
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.EthernetConnection{}, client.InNamespace(TESTNAMESPACE))
			if err != nil {
				fmt.Println("Can not delete EthernetConnectionCR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &controllertestfpga.FPGAFunction{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &controllertestcpu.CPUFunction{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &controllertestgpu.GPUFunction{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

		})

		AfterAll(func() {
			writer.Reset()
		})

		It("Test_1-1-1_FPGA-dma-FPGA", func() {

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection1)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

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

			// Create FPGACR
			err = createFPGACR(ctx, FPGA1[1])
			Expect(err).NotTo(HaveOccurred())
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			CRCStartPtuInitSet(mgr)
			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "ethernetconnectiontest1-wbconnection-decode-main-filter-resize-high-infer-main",
				}})
			if err != nil {
				fmt.Println("err:")
				fmt.Println(err)
			} else {
				fmt.Println("not err:")
				fmt.Println(got)
			}
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get ethernetconnectionCR
			var ethernetCR examplecomv1.EthernetConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "ethernetconnectiontest1-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&ethernetCR)

			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status
			expectedStatus := examplecomv1.EthernetConnectionStatus{
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "ethernetconnectiontest1",
					Namespace: "default",
				},
				Status: "Running",
				From: examplecomv1.EthernetFunctionStatus{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "ethernetconnectiontest1-wbfunction-decode-main",
						Namespace: "default",
					},
					Status: "OK",
				},
				To: examplecomv1.EthernetFunctionStatus{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "ethernetconnectiontest1-wbfunction-filter-resize-high-infer-main",
						Namespace: "default",
					},
					Status: "OK",
				},
			}
			Expect(ethernetCR.Status.DataFlowRef).To(Equal(expectedStatus.DataFlowRef))
			Expect(ethernetCR.Status.Status).To(Equal(expectedStatus.Status))
			Expect(ethernetCR.Status.From).To(Equal(expectedStatus.From))
			Expect(ethernetCR.Status.To).To(Equal(expectedStatus.To))

			// test for updating StartTime
			err = checkTime(testTime.Time, ethernetCR.Status.StartTime.Time)
			Expect(err).NotTo(HaveOccurred())

			// test for creating finalizer
			expectedFinalizer := []string{
				"ethernetconnection.finalizers.example.com.v1",
			}
			Expect(ethernetCR.Finalizers).To(Equal(expectedFinalizer))

			// check logs
			Expect(writer.String()).To(ContainSubstring("EthernetConnectionFPGAAccept() OK"))
			Expect(writer.String()).To(ContainSubstring("chenge srcStatus = STATUS_OK"))
			Expect(writer.String()).To(ContainSubstring("myNodeName node01"))
			Expect(writer.String()).To(ContainSubstring("src_node node01"))
			Expect(writer.String()).To(ContainSubstring("srcStatus OK"))
			Expect(writer.String()).To(ContainSubstring("dst_node node01"))
			Expect(writer.String()).To(ContainSubstring("dstStatus OK"))
			Expect(writer.String()).To(ContainSubstring("fpga_ptu_listen() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_ptu_accept() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_chain_connect_ingress() OK ret = 0"))
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

		It("Test_1-1-2_FPGA-FPGA_tcp-tcp", func() {
			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection2)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGAFunctionCR decode
			err = createFPGAFunction(ctx, FPGAFunctiondecode_tcp)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR-decode ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGAFunctionCR filter/resize
			err = createFPGAFunction(ctx, FPGAFunctionfilter_tcp)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR-filter/resize ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGACR
			err = createFPGACR(ctx, FPGA2[0])
			Expect(err).NotTo(HaveOccurred())
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGACR
			err = createFPGACR(ctx, FPGA2[1])
			Expect(err).NotTo(HaveOccurred())
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR", err)
			}
			Expect(err).NotTo(HaveOccurred())
			CRCStartPtuInitSet(mgr)
			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "ethernetconnectiontest2-wbconnection-decode-main-filter-resize-high-infer-main",
				}})
			if err != nil {
				fmt.Println("err:")
				fmt.Println(err)
			} else {
				fmt.Println("not err:")
				fmt.Println(got)
			}
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())
			// Get EthernetConnectionCR
			var ethernetCR examplecomv1.EthernetConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "ethernetconnectiontest2-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&ethernetCR)

			Expect(err).NotTo(HaveOccurred())

			// test for creating finalizer
			expectedFinalizer := []string{
				"ethernetconnection.finalizers.example.com.v1",
			}
			Expect(ethernetCR.Finalizers).To(Equal(expectedFinalizer))

			// test for updating CR.Status
			expectedStatus := examplecomv1.EthernetConnectionStatus{
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "ethernetconnectiontest2",
					Namespace: "default",
				},
				Status: "Running",
				From: examplecomv1.EthernetFunctionStatus{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "ethernetconnectiontest2-wbfunction-decode-main",
						Namespace: "default",
					},
					Status: "OK",
				},
				To: examplecomv1.EthernetFunctionStatus{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "ethernetconnectiontest2-wbfunction-filter-resize-high-infer-main",
						Namespace: "default",
					},
					Status: "OK",
				},
			}
			Expect(ethernetCR.Status.DataFlowRef).To(Equal(expectedStatus.DataFlowRef))
			Expect(ethernetCR.Status.Status).To(Equal(expectedStatus.Status))
			Expect(ethernetCR.Status.From).To(Equal(expectedStatus.From))
			Expect(ethernetCR.Status.To).To(Equal(expectedStatus.To))

			// test for updating StartTime
			err = checkTime(testTime.Time, ethernetCR.Status.StartTime.Time)
			Expect(err).NotTo(HaveOccurred())

			// check logs
			Expect(writer.String()).To(ContainSubstring("EthernetConnectionFPGAAccept() OK"))
			Expect(writer.String()).To(ContainSubstring("chenge srcStatus = STATUS_OK"))
			Expect(writer.String()).To(ContainSubstring("myNodeName node01"))
			Expect(writer.String()).To(ContainSubstring("src_node node01"))
			Expect(writer.String()).To(ContainSubstring("srcStatus OK"))
			Expect(writer.String()).To(ContainSubstring("dst_node node01"))
			Expect(writer.String()).To(ContainSubstring("dstStatus OK"))
			Expect(writer.String()).To(ContainSubstring("fpga_ptu_listen() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_ptu_connect() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_chain_connect_egress() OK ret = 0"))
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

		It("Test_1-1-3_FPGA_no_networkstatus", func() {
			// Create FPGACR
			err = createFPGACR(ctx, FPGA3[0])
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create ChildbsCR
			err = createChildbs(ctx, Childbs)
			if err != nil {
				fmt.Println("There is a problem in createing ChildbsCR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			CRCStartPtuInitSet(mgr)
			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "testchildbs",
				}})
			if err != nil {
				fmt.Println("err:")
				fmt.Println(err)
			} else {
				fmt.Println("not err:")
				fmt.Println(got)
			}

			// Get ChildbsCR
			var cr examplecomv1.ChildBs
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "testchildbs",
				Namespace: "default",
			},
				&cr)
			Expect(err).NotTo(HaveOccurred())

			// test for creating finalizer
			expectedFinalizer := []string{
				"EthernetConnection.finalizers.example.com.v1",
			}
			Expect(cr.Finalizers).To(Equal(expectedFinalizer))

			expectedStatus := examplecomv1.ChildBsStatus{
				Regions: []examplecomv1.ChildBsRegion{},
				Status:  examplecomv1.ChildBsStatusPreparing,
				State:   examplecomv1.ChildBsReady,
			}
			Expect(cr.Status.State).To(Equal(expectedStatus.State))

			// check logs
			Expect(writer.String()).To(ContainSubstring("ChildBsCR Update Success"))
			Expect(writer.String()).To(ContainSubstring("fpga_ptu_init() ret = 0"))
			Expect(writer.String()).To(ContainSubstring("ChildBsCR Update Success"))
		})

		It("Test_1-1-4_CPU-GPU", func() {
			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection3)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create PodCR
			err = createPod(ctx, podCPU)
			if err != nil {
				fmt.Println("There is a problem in createing CPUPodCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create PodCR
			err = createPod(ctx, podGPU)
			if err != nil {
				fmt.Println("There is a problem in createing GPUPodCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunctiondecode)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunctionhighinfer)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			CRCStartPtuInitSet(mgr)
			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "ethernetconnectiontest3-wbconnection-decode-main-filter-resize-high-infer-main",
				}})
			if err != nil {
				fmt.Println("err:")
				fmt.Println(err)
			} else {
				fmt.Println("not err:")
				fmt.Println(got)
			}
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get EthernetConnectionCR
			var ethernetCR examplecomv1.EthernetConnection

			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "ethernetconnectiontest3-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&ethernetCR)
			Expect(err).NotTo(HaveOccurred())

			expectedStatus := examplecomv1.EthernetConnectionStatus{
				Status: "Running",
			}
			Expect(ethernetCR.Status.Status).To(Equal(expectedStatus.Status))

			// check logs
			Expect(writer.String()).To(ContainSubstring("myNodeName node01"))
			Expect(writer.String()).To(ContainSubstring("src_node node01"))
			Expect(writer.String()).To(ContainSubstring("srcStatus OK"))
			Expect(writer.String()).To(ContainSubstring("dst_node node01"))
			Expect(writer.String()).To(ContainSubstring("dstStatus OK"))
			Expect(writer.String()).To(ContainSubstring("Status Update."))
			Expect(writer.String()).NotTo(ContainSubstring("fpga_ptu_init() ret = 0"))
			Expect(writer.String()).NotTo(ContainSubstring("fpga_ptu_init() ret = 1"))
		})

		It("Test_1-1-5_Delete_FPGA_childbsCR", func() {
			// Create EthernetCR
			err = createEthernetConnection(ctx, EthernetConnection4)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGAFunctionCR decode
			err = createFPGAFunction(ctx, FPGAFunctiondecode_tcp2)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR-decode ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGAFunctionCR filter/resize
			err = createFPGAFunction(ctx, FPGAFunctionfilter_tcp2)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR-filter/resize ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGACR
			err = createFPGACR(ctx, FPGA4[0])
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create ChildbsCR
			err = createChildbs(ctx, Childbs2)
			if err != nil {
				fmt.Println("There is a problem in createing ChildbsCR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			CRCStartPtuInitSet(mgr)
			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "testchildbs2",
				}})
			if err != nil {
				fmt.Println("err:")
				fmt.Println(err)
			} else {
				fmt.Println("not err:")
				fmt.Println(got)
			}
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get ChildBsCR
			var cr examplecomv1.ChildBs
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "testchildbs2",
				Namespace: TESTNAMESPACE,
			},
				&cr)

			Expect(err).NotTo(HaveOccurred())

			// test for creating finalizer
			expectedFinalizer := []string{
				"deviceinfo.finalizers.example.com.v1",
				"EthernetConnection.finalizers.example.com.v1",
			}
			Expect(cr.Finalizers).To(Equal(expectedFinalizer))

			// Delete ChildbsCR
			err = deleteChildbs(ctx, Childbs2)
			if err != nil {
				fmt.Println("There is a problem in deleting ChildbsCR", err)
			}
			Expect(err).NotTo(HaveOccurred())
			CRCStartPtuInitSet(mgr)

			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "testchildbs2",
					Namespace: TESTNAMESPACE,
				}})
			if err != nil {
				fmt.Println("err:")
				fmt.Println(err)
			} else {
				fmt.Println("not err:")
				fmt.Println(got)
			}
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get ChildBsCR
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "testchildbs2",
				Namespace: TESTNAMESPACE,
			},
				&cr)

			Expect(err).NotTo(HaveOccurred())

			// test for deleting finalizer
			expectedFinalizer = []string{
				"deviceinfo.finalizers.example.com.v1",
			}
			Expect(cr.Finalizers).To(Equal(expectedFinalizer))

			CRCStartPtuInitSet(mgr)
			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "ethernetconnectiontest4-wbconnection-decode-main-filter-resize-high-infer-main",
					Namespace: TESTNAMESPACE,
				}})
			if err != nil {
				fmt.Println("err:")
				fmt.Println(err)
			} else {
				fmt.Println("not err:")
				fmt.Println(got)
			}
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get EthernetConnectionCR
			var ethernetCR examplecomv1.EthernetConnection

			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "ethernetconnectiontest4-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&ethernetCR)
			Expect(err).NotTo(HaveOccurred())

			// test for Status is Runnning
			expectedStatus := examplecomv1.EthernetConnectionStatus{
				Status: "Running",
			}
			Expect(ethernetCR.Status.Status).To(Equal(expectedStatus.Status))

			// check logs
			Expect(writer.String()).To(ContainSubstring("EthernetConnectionFPGAAccept() OK"))
			Expect(writer.String()).To(ContainSubstring("chenge srcStatus = STATUS_OK"))
			Expect(writer.String()).To(ContainSubstring("myNodeName node01"))
			Expect(writer.String()).To(ContainSubstring("src_node node01"))
			Expect(writer.String()).To(ContainSubstring("srcStatus OK"))
			Expect(writer.String()).To(ContainSubstring("dst_node node01"))
			Expect(writer.String()).To(ContainSubstring("dstStatus OK"))
			Expect(writer.String()).To(ContainSubstring("fpga_ptu_listen() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_ptu_accept() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_chain_connect_ingress() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("Status Update."))
		})

		It("Test_1-1-6_Update_FPGA_childbsCR__NotNetworkModule", func() {
			// Create EthernetCR
			err = createEthernetConnection(ctx, EthernetConnection5)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGAFunctionCR decode
			err = createFPGAFunction(ctx, FPGAFunctiondecode_tcp3)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR-decode ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGAFunctionCR filter/resize
			err = createFPGAFunction(ctx, FPGAFunctionfilter_tcp3)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR-filter/resize ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGACR
			err = createFPGACR(ctx, FPGA5[0])
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create ChildbsCR
			err = createChildbs(ctx, Childbs3)
			if err != nil {
				fmt.Println("There is a problem in createing ChildbsCR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			CRCStartPtuInitSet(mgr)
			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "ethernetconnectiontest5-wbconnection-decode-main-filter-resize-high-infer-main",
					Namespace: TESTNAMESPACE,
				}})
			if err != nil {
				fmt.Println("err:")
				fmt.Println(err)
			} else {
				fmt.Println("not err:")
				fmt.Println(got)
			}
			Expect(err).NotTo(HaveOccurred())

			CRCStartPtuInitSet(mgr)
			// Start Reconcile
			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "testchildbs3",
					Namespace: TESTNAMESPACE,
				}})
			if err != nil {
				fmt.Println("err:")
				fmt.Println(err)
			} else {
				fmt.Println("not err:")
				fmt.Println(got)
			}

			// Get ChildBsCR
			var cr examplecomv1.ChildBs
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "testchildbs3",
				Namespace: TESTNAMESPACE,
			},
				&cr)

			Expect(err).NotTo(HaveOccurred())

			//  update ChildBsCR.Status.State NotStopNetworkModule
			cr.Status.State = examplecomv1.ChildBsNotStopNetworkModule
			// Update ChildbsCR
			err = updateChildbs(ctx, cr)
			if err != nil {
				fmt.Println("There is a problem in updating ChildbsCR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			CRCStartPtuInitSet(mgr)
			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "testchildbs3",
					Namespace: TESTNAMESPACE,
				}})
			if err != nil {
				fmt.Println("err:")
				fmt.Println(err)
			} else {
				fmt.Println("not err:")
				fmt.Println(got)
			}

			// Get ChildBsCR
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "testchildbs3",
				Namespace: TESTNAMESPACE,
			},
				&cr)

			Expect(err).NotTo(HaveOccurred())

			// test for ChildBsCR.Status.State is NotWriteBitstreamFile
			expectedStatus := examplecomv1.ChildBsStatus{
				State: "NotWriteBitstreamFile",
			}
			Expect(cr.Status.State).To(Equal(expectedStatus.State))

			// check logs
			Expect(writer.String()).To(ContainSubstring("EthernetConnectionFPGAAccept() OK"))
			Expect(writer.String()).To(ContainSubstring("chenge srcStatus = STATUS_OK"))
			Expect(writer.String()).To(ContainSubstring("myNodeName node01"))
			Expect(writer.String()).To(ContainSubstring("src_node node01"))
			Expect(writer.String()).To(ContainSubstring("srcStatus OK"))
			Expect(writer.String()).To(ContainSubstring("dst_node node01"))
			Expect(writer.String()).To(ContainSubstring("dstStatus OK"))
			Expect(writer.String()).To(ContainSubstring("fpga_ptu_listen() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_ptu_accept() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_chain_connect_ingress() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("Status Update."))
			Expect(writer.String()).To(ContainSubstring("fpga_ptu_exit() ret = 0"))
			Expect(writer.String()).To(ContainSubstring("ChildBsCR Update Success"))
		})

		It("Test_1-1-7_Update_FPGA_ChildBs_Not_NotNetworkModule", func() {
			// Create EthernetCR
			err = createEthernetConnection(ctx, EthernetConnection6)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGAFunctionCR decode
			err = createFPGAFunction(ctx, FPGAFunctiondecode_tcp4)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR-decode ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGAFunctionCR filter/resize
			err = createFPGAFunction(ctx, FPGAFunctionfilter_tcp4)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR-filter/resize ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGACR
			err = createFPGACR(ctx, FPGA6[0])
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create ChildbsCR
			err = createChildbs(ctx, Childbs4)
			if err != nil {
				fmt.Println("There is a problem in createing ChildbsCR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			CRCStartPtuInitSet(mgr)
			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "ethernetconnectiontest6-wbconnection-decode-main-filter-resize-high-infer-main",
					Namespace: TESTNAMESPACE,
				}})
			if err != nil {
				fmt.Println("err:")
				fmt.Println(err)
			} else {
				fmt.Println("not err:")
				fmt.Println(got)
			}
			Expect(err).NotTo(HaveOccurred())

			CRCStartPtuInitSet(mgr)
			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "testchildbs4",
					Namespace: TESTNAMESPACE,
				}})
			if err != nil {
				fmt.Println("err:")
				fmt.Println(err)
			} else {
				fmt.Println("not err:")
				fmt.Println(got)
			}

			// Get ChildBsCR
			var cr examplecomv1.ChildBs
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "testchildbs4",
				Namespace: TESTNAMESPACE,
			},
				&cr)

			Expect(err).NotTo(HaveOccurred())

			// Update ChildbsCR
			err = updateChildbs(ctx, cr)
			if err != nil {
				fmt.Println("There is a problem in updating ChildbsCR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			CRCStartPtuInitSet(mgr)
			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Name:      "testchildbs4",
					Namespace: TESTNAMESPACE,
				}})
			if err != nil {
				fmt.Println("err:")
				fmt.Println(err)
			} else {
				fmt.Println("not err:")
				fmt.Println(got)
			}

			// Get ChildBsCR
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "testchildbs4",
				Namespace: TESTNAMESPACE,
			},
				&cr)

			Expect(err).NotTo(HaveOccurred())

			// test for ChildBsCR.Status.State is not NotWriteBitstreamFile
			expectedStatus := examplecomv1.ChildBsStatus{
				State: "NotWriteBitstreamFile",
			}
			Expect(cr.Status.State).NotTo(Equal(expectedStatus.State))

			// check logs
			Expect(writer.String()).To(ContainSubstring("EthernetConnectionFPGAAccept() OK"))
			Expect(writer.String()).To(ContainSubstring("chenge srcStatus = STATUS_OK"))
			Expect(writer.String()).To(ContainSubstring("myNodeName node01"))
			Expect(writer.String()).To(ContainSubstring("src_node node01"))
			Expect(writer.String()).To(ContainSubstring("srcStatus OK"))
			Expect(writer.String()).To(ContainSubstring("dst_node node01"))
			Expect(writer.String()).To(ContainSubstring("dstStatus OK"))
			Expect(writer.String()).To(ContainSubstring("fpga_ptu_listen() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_ptu_accept() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("fpga_chain_connect_ingress() OK ret = 0"))
			Expect(writer.String()).To(ContainSubstring("Status Update."))
			Expect(writer.String()).NotTo(ContainSubstring("fpga_ptu_exit() ret = 0"))
			Expect(writer.String()).NotTo(ContainSubstring("fpga_ptu_exit() ret = 1"))
		})
		It("Test_1-1-8_CPU-CPU_tcp", func() {
			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection118)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR decode
			err = createCPUFunction(ctx, CPUFunctionDecode118)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR-filter/resize ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR filter-resize
			err = createCPUFunction(ctx, CPUFunctionFilterResize118)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR-filter/resize ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			err = reconciler.Create(ctx, &CPUfilterresizePOD)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "ethernetcontest118-wbconnection-decode-main-filter-resize-high-infer-main",
				}})

			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get ethernetconnectionCR
			var ethernetCR examplecomv1.EthernetConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "ethernetcontest118-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&ethernetCR)

			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status
			expectedStatus := examplecomv1.EthernetConnectionStatus{
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "ethernetcontest118",
					Namespace: "default",
				},
				Status: "Running",
				From: examplecomv1.EthernetFunctionStatus{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "ethernetcontest118-decode-main",
						Namespace: "default",
					},
					Status: "OK",
				},
				To: examplecomv1.EthernetFunctionStatus{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "ethernetcontest118-filter-resize-main",
						Namespace: "default",
					},
					Status: "OK",
				},
			}
			Expect(ethernetCR.Status.DataFlowRef).To(Equal(expectedStatus.DataFlowRef))
			Expect(ethernetCR.Status.Status).To(Equal(expectedStatus.Status))
			Expect(ethernetCR.Status.From).To(Equal(expectedStatus.From))
			Expect(ethernetCR.Status.To).To(Equal(expectedStatus.To))

			// test for updating StartTime
			err = checkTime(testTime.Time, ethernetCR.Status.StartTime.Time)
			Expect(err).NotTo(HaveOccurred())

			// test for creating finalizer
			expectedFinalizer := []string{
				"ethernetconnection.finalizers.example.com.v1",
			}
			Expect(ethernetCR.Finalizers).To(Equal(expectedFinalizer))

			// check logs
			Expect(writer.String()).To(ContainSubstring("myNodeName node01"))
			Expect(writer.String()).To(ContainSubstring("src_node node01"))
			Expect(writer.String()).To(ContainSubstring("srcStatus OK"))
			Expect(writer.String()).To(ContainSubstring("dst_node node01"))
			Expect(writer.String()).To(ContainSubstring("dstStatus OK"))
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
			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, ethernetConnectionUPDATE)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "ethernetcontestupdate-wbconnection-main",
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
			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, ethernetConnectionDELETE)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			var ethernetccrCR examplecomv1.EthernetConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "ethernetcontestdelete-wbconnection-main",
				Namespace: TESTNAMESPACE,
			},
				&ethernetccrCR)

			err = k8sClient.Delete(ctx, &ethernetccrCR)

			var ethernetccrCR2 examplecomv1.EthernetConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "ethernetcontestdelete-wbconnection-main",
				Namespace: TESTNAMESPACE,
			},
				&ethernetccrCR2)

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "ethernetcontestdelete-wbconnection-main",
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

			var ethernetccrCR3 examplecomv1.EthernetConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "ethernetcontestdelete-wbconnection-main",
				Namespace: TESTNAMESPACE,
			},
				&ethernetccrCR3)

			Expect(err).To(MatchError(ContainSubstring("not found")))

		})

		It("Test_6-2-1_Status Not Change route ", func() {

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection621)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

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
			err = createFPGACR(ctx, FPGA621[0])
			Expect(err).NotTo(HaveOccurred())
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGACR
			err = createFPGACR(ctx, FPGA621[1])
			Expect(err).NotTo(HaveOccurred())
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			CRCStartPtuInitSet(mgr)
			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "ethernetconnectiontest621-wbconnection-decode-main-filter-resize-high-infer-main",
				}})
			if err != nil {
				fmt.Println("err:")
				fmt.Println(err)
			} else {
				fmt.Println("not err:")
				fmt.Println(got)
			}
			Expect(got).To(Equal(ctrl.Result{Requeue: true}))
			Expect(err).NotTo(HaveOccurred())

			// Get ethernetconnectionCR
			var ethernetCR examplecomv1.EthernetConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "ethernetconnectiontest621-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&ethernetCR)

			Expect(err).NotTo(HaveOccurred())

			Expect(ethernetCR.Status.DataFlowRef).To(Equal(EthernetConnection621.Status.DataFlowRef))
			Expect(ethernetCR.Status.Status).To(Equal(EthernetConnection621.Status.Status))
			Expect(ethernetCR.Status.From).To(Equal(EthernetConnection621.Status.From))
			Expect(ethernetCR.Status.To).To(Equal(EthernetConnection621.Status.To))

			// test for updating StartTime
			err = checkTime(testTime.Time, ethernetCR.Status.StartTime.Time)
			Expect(err).To(HaveOccurred())

			var nilfzr []string
			nilfzr = nil
			Expect(ethernetCR.Finalizers).To(Equal(nilfzr))

			// check logs
			Expect(writer.String()).To(ContainSubstring("myNodeName node01"))
			Expect(writer.String()).To(ContainSubstring("src_node node01"))
			Expect(writer.String()).To(ContainSubstring("srcStatus OK"))
			Expect(writer.String()).To(ContainSubstring("dst_node node01"))
			Expect(writer.String()).To(ContainSubstring("dstStatus OK"))
			Expect(writer.String()).To(ContainSubstring("CR Status is not Changed."))

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

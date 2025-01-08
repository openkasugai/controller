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

	examplecomv1 "EthernetConnection/api/v1"
	controllertestcpu "EthernetConnection/internal/controller/test/type/CPU"
	controllertestfpga "EthernetConnection/internal/controller/test/type/FPGA"
	controllertestgpu "EthernetConnection/internal/controller/test/type/GPU"

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
/*
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
*/
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
/*
func createGPUFunction(ctx context.Context, gpufcr controllertestgpu.GPUFunction) error {
	tmp := &controllertestgpu.GPUFunction{}
	*tmp = gpufcr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}

	return nil
}

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

// ChildbsCR
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

/*
// Create FPGAList
func createFPGAList(ctx context.Context, fpgalist examplecomv1.FPGAList) error {
	tmp := &examplecomv1.FPGAList{}
	*tmp = fpgalist
	tmp.TypeMeta = fpgalist.TypeMeta
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}

	return nil
}
*/
// To describe test cases in CPUFunctionController
var _ = Describe("EthernetConnectionController", func() {
	var mgr ctrl.Manager
	var err error
	ctx := context.Background()
	// ctx := context.WithValue(context.Background(), log.Logger)
	var fakerecorder = record.NewFakeRecorder(10)
	var writer = bytes.Buffer{}
	// var stopFunc func()

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
			reconciler = EthernetConnectionReconciler{
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

					EthernetConnectionFPGAInit(mgr)

					expected := []string{
						"/dev/xethernet_21330621T04L",
						"/dev/xethernet_21330621T01J",
						"/dev/xethernet_21330621T00Y",
						"/dev/xethernet_21330621T00D",
					}
					Expect(FpgaDevList).To(Equal(expected))

					// check logs
					Expect(writer.String()).To(ContainSubstring("fpgaDataPh3.NodeName :" + "node01"))
					Expect(writer.String()).To(ContainSubstring("fpga_shmem_controller_init() ret = "))
				})
		*/

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
			// Get ethernetconnectionCR
			var ethernetCR examplecomv1.EthernetConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "ethernetconnectiontest2-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&ethernetCR)

			Expect(err).NotTo(HaveOccurred())

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
		/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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
			/*
				var cr examplecomv1.ChildBs
				err = k8sClient.Get(ctx, client.ObjectKey{
					Name: "testchildbs",
					Namespace: "default",
					},
				&cr)
				fmt.Println("here ChildBsCRdata")
				fmt.Println((*cr.Spec.Regions[0].Modules.Ptu.Parameters)["MacAddress"].StrVal)
				fmt.Println("------------------------")
			*/
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
			Expect(err).NotTo(HaveOccurred())
			var cr examplecomv1.ChildBs
			_ = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "testchildbs",
				Namespace: "default",
			},
				&cr)

			expectedStatus := examplecomv1.ChildBsStatus{
				Regions: []examplecomv1.ChildBsRegion{},
				Status:  examplecomv1.ChildBsStatusPreparing,
				State:   examplecomv1.ChildBsReady,
			}
			Expect(cr.Status.State).To(Equal(expectedStatus.State))
			Expect(writer.String()).To(ContainSubstring("ChildBsCR Update Success"))
			Expect(writer.String()).To(ContainSubstring("fpga_ptu_init() ret = 0"))
			Expect(writer.String()).To(ContainSubstring("ChildBsCR Update Success"))
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

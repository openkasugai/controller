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

	examplecomv1 "WBConnection/api/v1"
	controllertestethernet "WBConnection/internal/controller/test/type/ethernet"
	controllertestpcie "WBConnection/internal/controller/test/type/pcie"

	// ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"

	// Additional files

	"k8s.io/apimachinery/pkg/types"
)

func getMgr(mgr ctrl.Manager) (ctrl.Manager, error) {
	if mgr == nil {
		return ctrl.NewManager(cfg, ctrl.Options{
			// Scheme: testScheme,
			Scheme: k8sClient.Scheme(),
		})
	}
	return mgr, nil
}

// Create WBConnection CR
func createWBConnection(ctx context.Context, wbccr examplecomv1.WBConnection) error {
	tmp := &examplecomv1.WBConnection{}
	*tmp = wbccr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}

	return nil
}

// Create PCIeConnection CR
func createPCIeConnection(ctx context.Context, pcieccr controllertestpcie.PCIeConnection) error {
	tmp := &controllertestpcie.PCIeConnection{}
	*tmp = pcieccr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}

	return nil
}

// Create EthernetConnection CR
func createEthernetConnection(ctx context.Context, ethernetccr controllertestethernet.EthernetConnection) error {
	tmp := &controllertestethernet.EthernetConnection{}
	*tmp = ethernetccr
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

// To describe test cases in CPUFunctionController
var _ = Describe("PCIeConnectionController", func() {
	var mgr ctrl.Manager
	var err error
	ctx := context.Background()
	var fakerecorder = record.NewFakeRecorder(10)
	var writer = bytes.Buffer{}
	// var stopFunc func()

	// This test case is for reconciler
	Context("Test for PCIeConnectionReconciler", Ordered, func() {
		var reconciler WBConnectionReconciler

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

			err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

		})
		// Before Context runs, BeforeAll is executed once.
		BeforeEach(func() {
			// loger initialized
			writer.Reset()

			// recorder initialized
			fakerecorder = record.NewFakeRecorder(10)
			reconciler = WBConnectionReconciler{
				Client: k8sClient,
				// Client: mgr.GetClient(),
				Scheme: testScheme,
				// Scheme:   k8sClient.Scheme(),
				Recorder: fakerecorder,
			}

			err = reconciler.SetupWithManager(mgr)
			if err != nil {
				fmt.Println("An error occur during setupwithManager: ", err)
			}
			Expect(err).NotTo(HaveOccurred())
			os.Setenv("K8S_NODENAME", "node01")

			// stopFunc = startMgr(ctx, mgr)

			// To delete crdata It
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.WBConnection{}, client.InNamespace(TESTNAMESPACE))
			if err != nil {
				fmt.Println("Can not delete PCIeConnectionCR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &controllertestpcie.PCIeConnection{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &controllertestethernet.EthernetConnection{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			// err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace(TESTNAMESPACE))
			// Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

		})

		AfterAll(func() {
			writer.Reset()
			// 	// stop manager
			// 	stopFunc()
			// 	time.Sleep(100 * time.Millisecond)
		})

		//Test for GetFunc
		It("Test_1-0-1_Setup-configMap", func() {
			err = createConfig(ctx, connectionkindmap)
			Expect(err).NotTo(HaveOccurred())

			err = LoadConfigMap(&reconciler)
			Expect(err).NotTo(HaveOccurred())

			expected := []ConnectionKindMap{
				{
					ConnectionMethod: "host-mem",
					ConnectionCRKind: "PCIeConnection",
				},
				{
					ConnectionMethod: "host-100gether",
					ConnectionCRKind: "EthernetConnection",
				},
			}
			Expect(gConnectionKindmap).To(Equal(expected))

		})
		//Test for GetFunc
		It("Test_1-1-1_Start-of-chain", func() {
			// Create WBConnectionCR
			err = createWBConnection(ctx, wbconnection1Start)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntest1-wbconnection-wb-start-of-chain-decode-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var wbcCR examplecomv1.WBConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest1-wbconnection-wb-start-of-chain-decode-main",
				Namespace: TESTNAMESPACE,
			},
				&wbcCR)

			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status
			expectedStatus := examplecomv1.WBConnectionStatus{
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "wbconntest1",
					Namespace: "default",
				},
				Status: "Deployed",
				From: examplecomv1.FromToWBFunction{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "wb-start-of-chain",
						Namespace: "default",
					},
				},
				To: examplecomv1.FromToWBFunction{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "wbconntest1-wbfunction-decode-main",
						Namespace: "default",
					},
				},
				ConnectionMethod: "host-100gether",
			}

			Expect(wbcCR.Status.DataFlowRef).To(Equal(expectedStatus.DataFlowRef))
			Expect(wbcCR.Status.Status).To(Equal(expectedStatus.Status))
			Expect(wbcCR.Status.From).To(Equal(expectedStatus.From))
			Expect(wbcCR.Status.To).To(Equal(expectedStatus.To))
			Expect(wbcCR.Status).To(Equal(expectedStatus))

			// test for creating finalizer
			expectedFinalizer := []string{
				"WBConnection.finalizers.example.com.v1",
			}
			Expect(wbcCR.Finalizers).To(Equal(expectedFinalizer))

			// check logs
			Expect(writer.String()).To(ContainSubstring("Reconcile start."))
			Expect(writer.String()).To(ContainSubstring("Finalizername=WBConnection.finalizers.example.com.v1"))
			Expect(writer.String()).To(ContainSubstring("Status Information Change start."))
			Expect(writer.String()).To(ContainSubstring("Finalizername=WBConnection.finalizers.example.com.v1"))
			Expect(writer.String()).To(ContainSubstring("Status Update."))
			Expect(writer.String()).To(ContainSubstring("Status Information Change end."))

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
		It("Test_1-1-2_Ethernet", func() {
			// Create WBConnectionCR
			err = createWBConnection(ctx, wbconnection2Ether)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntest2-wbconnection-decode-main-filter-resize-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: true}))
			Expect(err).NotTo(HaveOccurred())

			var wbcCR examplecomv1.WBConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest2-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&wbcCR)

			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status
			expectedStatus := examplecomv1.WBConnectionStatus{
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "wbconntest2",
					Namespace: "default",
				},
				Status: "Waiting",
				From: examplecomv1.FromToWBFunction{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "wbconntest2-wbfunction-decode-main",
						Namespace: "default",
					},
				},
				To: examplecomv1.FromToWBFunction{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "wbconntest2-wbfunction-filter-resize-high-infer-main",
						Namespace: "default",
					},
				},
				ConnectionMethod: "host-100gether",
			}

			Expect(wbcCR.Status.DataFlowRef).To(Equal(expectedStatus.DataFlowRef))
			Expect(wbcCR.Status.Status).To(Equal(expectedStatus.Status))
			Expect(wbcCR.Status.From).To(Equal(expectedStatus.From))
			Expect(wbcCR.Status.To).To(Equal(expectedStatus.To))
			Expect(wbcCR.Status).To(Equal(expectedStatus))

			// test for creating finalizer
			expectedFinalizer := []string{
				"WBConnection.finalizers.example.com.v1",
			}
			Expect(wbcCR.Finalizers).To(Equal(expectedFinalizer))

			// check logs
			Expect(writer.String()).To(ContainSubstring("Reconcile start."))
			Expect(writer.String()).To(ContainSubstring("Finalizername=WBConnection.finalizers.example.com.v1"))
			Expect(writer.String()).To(ContainSubstring("Maked Connection does not exist."))
			Expect(writer.String()).To(ContainSubstring("CustomResource Create."))
			Expect(writer.String()).To(ContainSubstring("kind :EthernetConnection"))
			Expect(writer.String()).To(ContainSubstring("apiVersion :example.com/v1"))
			Expect(writer.String()).To(ContainSubstring("name :wbconntest2-wbconnection-decode-main-filter-resize-high-infer-main"))
			Expect(writer.String()).To(ContainSubstring("namespace :default"))
			Expect(writer.String()).To(ContainSubstring("Status Information Change start."))
			Expect(writer.String()).To(ContainSubstring("Finalizername=WBConnection.finalizers.example.com.v1"))
			Expect(writer.String()).To(ContainSubstring("Status Update."))
			Expect(writer.String()).To(ContainSubstring("Status Information Change end."))

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

			// check EtherntCR
			var ethernetCR controllertestethernet.EthernetConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest2-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&ethernetCR)

			Expect(err).NotTo(HaveOccurred())

			expectedCR := controllertestethernet.EthernetConnection{
				// TypeMeta: metav1.TypeMeta{
				// 	Kind:       "EthernetConnection",
				// 	APIVersion: "example.com/v1",
				// },
				// ObjectMeta: metav1.ObjectMeta{
				// 	Name:      "wbconntest2-wbconnection-decode-main-filter-resize-high-infer-main",
				// 	Namespace: "default",
				// },
				Spec: controllertestethernet.EthernetConnectionSpec{
					DataFlowRef: controllertestethernet.WBNamespacedName{
						Name:      "wbconntest2",
						Namespace: "default",
					},
					From: controllertestethernet.EthernetFunctionSpec{
						WBFunctionRef: controllertestethernet.WBNamespacedName{
							Name:      "wbconntest2-wbfunction-decode-main",
							Namespace: "default",
						},
					},
					To: controllertestethernet.EthernetFunctionSpec{
						WBFunctionRef: controllertestethernet.WBNamespacedName{
							Name:      "wbconntest2-wbfunction-filter-resize-high-infer-main",
							Namespace: "default",
						},
					},
				},
			}

			Expect(ethernetCR.Spec).To(Equal(expectedCR.Spec))

		})
		It("Test_1-1-3_PCIe", func() {
			// Create WBConnectionCR
			err = createWBConnection(ctx, wbconnection3PCIe)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntest3-wbconnection-decode-main-filter-resize-low-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: true}))
			Expect(err).NotTo(HaveOccurred())

			var wbcCR examplecomv1.WBConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest3-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&wbcCR)

			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status
			expectedStatus := examplecomv1.WBConnectionStatus{
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "wbconntest3",
					Namespace: "default",
				},
				Status: "Waiting",
				From: examplecomv1.FromToWBFunction{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "wbconntest3-wbfunction-decode-main",
						Namespace: "default",
					},
				},
				To: examplecomv1.FromToWBFunction{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "wbconntest3-wbfunction-filter-resize-low-infer-main",
						Namespace: "default",
					},
				},
				ConnectionMethod: "host-mem",
			}

			Expect(wbcCR.Status.DataFlowRef).To(Equal(expectedStatus.DataFlowRef))
			Expect(wbcCR.Status.Status).To(Equal(expectedStatus.Status))
			Expect(wbcCR.Status.From).To(Equal(expectedStatus.From))
			Expect(wbcCR.Status.To).To(Equal(expectedStatus.To))
			Expect(wbcCR.Status).To(Equal(expectedStatus))

			// test for creating finalizer
			expectedFinalizer := []string{
				"WBConnection.finalizers.example.com.v1",
			}
			Expect(wbcCR.Finalizers).To(Equal(expectedFinalizer))

			// check logs
			Expect(writer.String()).To(ContainSubstring("Reconcile start."))
			Expect(writer.String()).To(ContainSubstring("Finalizername=WBConnection.finalizers.example.com.v1"))
			Expect(writer.String()).To(ContainSubstring("Maked Connection does not exist."))
			Expect(writer.String()).To(ContainSubstring("CustomResource Create."))
			Expect(writer.String()).To(ContainSubstring("kind :PCIeConnection"))
			Expect(writer.String()).To(ContainSubstring("apiVersion :example.com/v1"))
			Expect(writer.String()).To(ContainSubstring("name :wbconntest3-wbconnection-decode-main-filter-resize-low-infer-main"))
			Expect(writer.String()).To(ContainSubstring("namespace :default"))
			Expect(writer.String()).To(ContainSubstring("Status Information Change start."))
			Expect(writer.String()).To(ContainSubstring("Finalizername=WBConnection.finalizers.example.com.v1"))
			Expect(writer.String()).To(ContainSubstring("Status Update."))
			Expect(writer.String()).To(ContainSubstring("Status Information Change end."))

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

			// check EtherntCR
			var pcieCR controllertestpcie.PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest3-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)

			Expect(err).NotTo(HaveOccurred())

			expectedCR := controllertestpcie.PCIeConnection{
				// TypeMeta: metav1.TypeMeta{
				// 	Kind:       "EthernetConnection",
				// 	APIVersion: "example.com/v1",
				// },
				// ObjectMeta: metav1.ObjectMeta{
				// 	Name:      "wbconntest2-wbconnection-decode-main-filter-resize-high-infer-main",
				// 	Namespace: "default",
				// },
				Spec: controllertestpcie.PCIeConnectionSpec{
					DataFlowRef: controllertestpcie.WBNamespacedName{
						Name:      "wbconntest3",
						Namespace: "default",
					},
					From: controllertestpcie.PCIeFunctionSpec{
						WBFunctionRef: controllertestpcie.WBNamespacedName{
							Name:      "wbconntest3-wbfunction-decode-main",
							Namespace: "default",
						},
					},
					To: controllertestpcie.PCIeFunctionSpec{
						WBFunctionRef: controllertestpcie.WBNamespacedName{
							Name:      "wbconntest3-wbfunction-filter-resize-low-infer-main",
							Namespace: "default",
						},
					},
				},
			}

			Expect(pcieCR.Spec).To(Equal(expectedCR.Spec))
		})
		It("Test_1-1-4_End-of-chain", func() {
			// Create WBConnectionCR
			err = createWBConnection(ctx, wbconnection4End)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntest4-wbconnection-low-infer-main-wb-end-of-chain",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var wbcCR examplecomv1.WBConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest4-wbconnection-low-infer-main-wb-end-of-chain",
				Namespace: TESTNAMESPACE,
			},
				&wbcCR)

			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status
			expectedStatus := examplecomv1.WBConnectionStatus{
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "wbconntest4",
					Namespace: "default",
				},
				Status: "Deployed",
				From: examplecomv1.FromToWBFunction{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "wbconntest4-wbfuncction-low-infer-main",
						Namespace: "default",
					},
				},
				To: examplecomv1.FromToWBFunction{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "wb-end-of-chain",
						Namespace: "default",
					},
				},
				ConnectionMethod: "host-100gether",
			}

			Expect(wbcCR.Status.DataFlowRef).To(Equal(expectedStatus.DataFlowRef))
			Expect(wbcCR.Status.Status).To(Equal(expectedStatus.Status))
			Expect(wbcCR.Status.From).To(Equal(expectedStatus.From))
			Expect(wbcCR.Status.To).To(Equal(expectedStatus.To))
			Expect(wbcCR.Status).To(Equal(expectedStatus))

			// test for creating finalizer
			expectedFinalizer := []string{
				"WBConnection.finalizers.example.com.v1",
			}
			Expect(wbcCR.Finalizers).To(Equal(expectedFinalizer))

			// check logs
			Expect(writer.String()).To(ContainSubstring("Reconcile start."))
			Expect(writer.String()).To(ContainSubstring("Finalizername=WBConnection.finalizers.example.com.v1"))
			Expect(writer.String()).To(ContainSubstring("Status Information Change start."))
			Expect(writer.String()).To(ContainSubstring("Finalizername=WBConnection.finalizers.example.com.v1"))
			Expect(writer.String()).To(ContainSubstring("Status Update."))
			Expect(writer.String()).To(ContainSubstring("Status Information Change end."))

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
		It("Test_1-2-2_Update-Ethernet", func() {
			// Create EternetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnectionUpdate)
			Expect(err).NotTo(HaveOccurred())

			// Create WBConnectionCR
			err = createWBConnection(ctx, wbconnectionUpdate2Ether)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntest2upd-wbconnection-decode-main-filter-resize-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var wbcCR examplecomv1.WBConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest2upd-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&wbcCR)

			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status
			expectedStatus := examplecomv1.WBConnectionStatus{
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "wbconntest2upd",
					Namespace: "default",
				},
				Status: "Deployed",
				From: examplecomv1.FromToWBFunction{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "wbconntest2upd-wbfunction-decode-main",
						Namespace: "default",
					},
				},
				To: examplecomv1.FromToWBFunction{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "wbconntest2upd-wbfunction-filter-resize-high-infer-main",
						Namespace: "default",
					},
				},
				ConnectionMethod: "host-100gether",
			}

			Expect(wbcCR.Status.DataFlowRef).To(Equal(expectedStatus.DataFlowRef))
			Expect(wbcCR.Status.Status).To(Equal(expectedStatus.Status))
			Expect(wbcCR.Status.From).To(Equal(expectedStatus.From))
			Expect(wbcCR.Status.To).To(Equal(expectedStatus.To))
			Expect(wbcCR.Status).To(Equal(expectedStatus))

			// check logs
			Expect(writer.String()).To(ContainSubstring("Reconcile start."))
			Expect(writer.String()).To(ContainSubstring("Finalizername=WBConnection.finalizers.example.com.v1"))
			Expect(writer.String()).To(ContainSubstring("eventConnectionKind=1"))
			Expect(writer.String()).To(ContainSubstring("crData.Status.Status=Waiting"))
			Expect(writer.String()).To(ContainSubstring("crConnectionStatus.Status=Running"))
			Expect(writer.String()).To(ContainSubstring("Status Running Change start."))
			Expect(writer.String()).To(ContainSubstring("Status Update."))
			Expect(writer.String()).To(ContainSubstring("Status Running Change end."))

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
		It("Test_1-2-3_Update-PCIe", func() {
			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnectionUpdate)

			// Create WBConnectionCR
			err = createWBConnection(ctx, wbconnectionUpdate3PCIe)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntest3upd-wbconnection-decode-main-filter-resize-low-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var wbcCR examplecomv1.WBConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest3upd-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&wbcCR)

			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status
			expectedStatus := examplecomv1.WBConnectionStatus{
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "wbconntest3upd",
					Namespace: "default",
				},
				Status: "Deployed",
				From: examplecomv1.FromToWBFunction{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "wbconntest3upd-wbfunction-decode-main",
						Namespace: "default",
					},
				},
				To: examplecomv1.FromToWBFunction{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "wbconntest3upd-wbfunction-filter-resize-low-infer-main",
						Namespace: "default",
					},
				},
				ConnectionMethod: "host-mem",
			}

			Expect(wbcCR.Status.DataFlowRef).To(Equal(expectedStatus.DataFlowRef))
			Expect(wbcCR.Status.Status).To(Equal(expectedStatus.Status))
			Expect(wbcCR.Status.From).To(Equal(expectedStatus.From))
			Expect(wbcCR.Status.To).To(Equal(expectedStatus.To))
			Expect(wbcCR.Status).To(Equal(expectedStatus))

			// check logs
			Expect(writer.String()).To(ContainSubstring("Reconcile start."))
			Expect(writer.String()).To(ContainSubstring("Finalizername=WBConnection.finalizers.example.com.v1"))
			Expect(writer.String()).To(ContainSubstring("eventConnectionKind=1"))
			Expect(writer.String()).To(ContainSubstring("crData.Status.Status=Waiting"))
			Expect(writer.String()).To(ContainSubstring("crConnectionStatus.Status=Running"))
			Expect(writer.String()).To(ContainSubstring("Status Running Change start."))
			Expect(writer.String()).To(ContainSubstring("Status Update."))
			Expect(writer.String()).To(ContainSubstring("Status Running Change end."))

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
		It("Test_1-2-4_DELETE", func() {

			// Create WBConnectionCR
			err = createWBConnection(ctx, wbconnectionDELETE)
			Expect(err).NotTo(HaveOccurred())

			var wbcCR examplecomv1.WBConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntestdel-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&wbcCR)

			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.Delete(ctx, &wbcCR)

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntestdel-wbconnection-decode-main-filter-resize-low-infer-main",
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

			var wbcaCR examplecomv1.WBConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntestdel-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&wbcaCR)

			Expect(err).To(MatchError(ContainSubstring("not found")))

			// check logs
			Expect(writer.String()).To(ContainSubstring("Reconcile start."))
			Expect(writer.String()).To(ContainSubstring("Finalizername=WBConnection.finalizers.example.com.v1"))
			Expect(writer.String()).To(ContainSubstring("RemoveFinalizer Update."))

		})
	})

})

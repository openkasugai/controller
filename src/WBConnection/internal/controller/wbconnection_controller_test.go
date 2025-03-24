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
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	examplecomv1 "WBConnection/api/v1"
	controllertestethernet "WBConnection/internal/controller/test/type/ethernet"
	controllertestpcie "WBConnection/internal/controller/test/type/pcie"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// Additional files

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
)

func getMgr(mgr ctrl.Manager) (ctrl.Manager, error) {
	if mgr == nil {
		return ctrl.NewManager(cfg, ctrl.Options{
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

// Delete WBConnection CR
func deleteWBConnection(ctx context.Context, wbccr examplecomv1.WBConnection) error {
	tmp := &examplecomv1.WBConnection{}
	*tmp = wbccr
	err := k8sClient.Delete(ctx, tmp)
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

// Update PCIeConnection CR
func updatePCIeConnection(ctx context.Context, pcieccr controllertestpcie.PCIeConnection) error {
	tmp := &controllertestpcie.PCIeConnection{}
	*tmp = pcieccr
	err := k8sClient.Update(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Delete PCIeConnection CR
func deletePCIeConnection(ctx context.Context, pcieccr controllertestpcie.PCIeConnection) error {
	tmp := &controllertestpcie.PCIeConnection{}
	*tmp = pcieccr
	err := k8sClient.Delete(ctx, tmp)
	if err != nil {
		return err
	}

	return nil
}

// Update EthernetConnection CR
func updateEthernetConnection(ctx context.Context, ethernetccr controllertestethernet.EthernetConnection) error {
	tmp := &controllertestethernet.EthernetConnection{}
	*tmp = ethernetccr
	err := k8sClient.Update(ctx, tmp)
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
var _ = Describe("WBConnectionController", func() {
	var mgr ctrl.Manager
	var err error
	ctx := context.Background()
	var fakerecorder = record.NewFakeRecorder(10)
	var writer = bytes.Buffer{}

	// This test case is for reconciler
	Context("Test for WBConnectionReconciler", Ordered, func() {
		var reconciler WBConnectionReconciler

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
				Client:   k8sClient,
				Scheme:   k8sClient.Scheme(),
				Recorder: fakerecorder,
			}

			err = reconciler.SetupWithManager(mgr)
			if err != nil {
				fmt.Println("An error occur during setupwithManager: ", err)
			}
			Expect(err).NotTo(HaveOccurred())
			os.Setenv("K8S_NODENAME", "node01")

			// To delete crdata It
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.WBConnection{}, client.InNamespace(TESTNAMESPACE))
			if err != nil {
				fmt.Println("Can not delete PCIeConnectionCR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &controllertestpcie.PCIeConnection{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &controllertestethernet.EthernetConnection{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			err = createConfig(ctx, connectionkindmap)
			Expect(err).NotTo(HaveOccurred())

			err = LoadConfigMap(&reconciler)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterAll(func() {
			writer.Reset()
			time.Sleep(500 * time.Millisecond)
		})
		It("Test_5-1-1_Start-of-chain", func() {
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

		It("Test_5-1-2_Ethernet", func() {
			// Create WBConnectionCR
			err = createWBConnection(ctx, wbconnection2Ether)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntest2-wbconnection-decode-main-filter-resize-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: false}))
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

			var strtime metav1.Time
			strtime = metav1.Now()
			ethernetCR.Status.StartTime = strtime
			ethernetCR.Status.Status = "Running"
			err = updateEthernetConnection(ctx, ethernetCR)
			Expect(err).NotTo(HaveOccurred())

			_, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntest2-wbconnection-decode-main-filter-resize-high-infer-main",
				}})
			Expect(err).NotTo(HaveOccurred())

			var wbcCR2 examplecomv1.WBConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest2-wbconnection-decode-main-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&wbcCR2)
			Expect(err).NotTo(HaveOccurred())
			Expect(wbcCR2.Status.Status).To(Equal(examplecomv1.WBDeployStatusDeployed))
		})

		It("Test_5-1-3_PCIe", func() {
			// Create WBConnectionCR
			err = createWBConnection(ctx, wbconnection3PCIe)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntest3-wbconnection-decode-main-filter-resize-low-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: false}))
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

			// check PCIeCR
			var pcieCR controllertestpcie.PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest3-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)

			Expect(err).NotTo(HaveOccurred())

			expectedCR := controllertestpcie.PCIeConnection{
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

			var strtime metav1.Time
			strtime = metav1.Now()
			pcieCR.Status.StartTime = strtime
			pcieCR.Status.Status = "Running"
			err = updatePCIeConnection(ctx, pcieCR)
			Expect(err).NotTo(HaveOccurred())

			_, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntest3-wbconnection-decode-main-filter-resize-low-infer-main",
				}})
			Expect(err).NotTo(HaveOccurred())

			var wbcCR2 examplecomv1.WBConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest3-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&wbcCR2)
			Expect(err).NotTo(HaveOccurred())
			Expect(wbcCR2.Status.Status).To(Equal(examplecomv1.WBDeployStatusDeployed))
		})

		It("Test_5-1-4_End-of-chain", func() {
			// Create WBConnectionCR
			err = createWBConnection(ctx, wbconnection4End)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntest4-wbconnection-high-infer-main-wb-end-of-chain",
				}},
			)
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var wbcCR examplecomv1.WBConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest4-wbconnection-high-infer-main-wb-end-of-chain",
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
						Name:      "wbconntest4-wbfuncction-high-infer-main",
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

		It("Test_5-1-5_Delete_WBFunction", func() {
			// Create WBConnectionCR
			err = createWBConnection(ctx, wbconnection5PCIe)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntest5-wbconnection-decode-main-filter-resize-low-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: false}))
			Expect(err).NotTo(HaveOccurred())

			var wbcCR examplecomv1.WBConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest5-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&wbcCR)

			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status
			expectedStatus := examplecomv1.WBConnectionStatus{
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "wbconntest5",
					Namespace: "default",
				},
				Status: "Waiting",
				From: examplecomv1.FromToWBFunction{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "wbconntest5-wbfunction-decode-main",
						Namespace: "default",
					},
				},
				To: examplecomv1.FromToWBFunction{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "wbconntest5-wbfunction-filter-resize-low-infer-main",
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
			Expect(writer.String()).To(ContainSubstring("name :wbconntest5-wbconnection-decode-main-filter-resize-low-infer-main"))
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

			// check PCIeCR
			var pcieCR controllertestpcie.PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest5-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)

			Expect(err).NotTo(HaveOccurred())

			expectedCR := controllertestpcie.PCIeConnection{
				Spec: controllertestpcie.PCIeConnectionSpec{
					DataFlowRef: controllertestpcie.WBNamespacedName{
						Name:      "wbconntest5",
						Namespace: "default",
					},
					From: controllertestpcie.PCIeFunctionSpec{
						WBFunctionRef: controllertestpcie.WBNamespacedName{
							Name:      "wbconntest5-wbfunction-decode-main",
							Namespace: "default",
						},
					},
					To: controllertestpcie.PCIeFunctionSpec{
						WBFunctionRef: controllertestpcie.WBNamespacedName{
							Name:      "wbconntest5-wbfunction-filter-resize-low-infer-main",
							Namespace: "default",
						},
					},
				},
			}
			Expect(pcieCR.Spec).To(Equal(expectedCR.Spec))

			var strtime metav1.Time
			strtime = metav1.Now()
			pcieCR.Status.StartTime = strtime
			pcieCR.Status.Status = "Running"
			err = updatePCIeConnection(ctx, pcieCR)
			Expect(err).NotTo(HaveOccurred())

			_, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntest5-wbconnection-decode-main-filter-resize-low-infer-main",
				}})
			Expect(err).NotTo(HaveOccurred())

			var wbcCR2 examplecomv1.WBConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest5-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&wbcCR2)
			Expect(err).NotTo(HaveOccurred())
			Expect(wbcCR2.Status.Status).To(Equal(examplecomv1.WBDeployStatusDeployed))

			err = deleteWBConnection(ctx, wbcCR2)
			Expect(err).NotTo(HaveOccurred())

			_, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntest5-wbconnection-decode-main-filter-resize-low-infer-main",
				}})
			Expect(err).NotTo(HaveOccurred())

			var pcieCR2 controllertestpcie.PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest5-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR2)

			if errors.IsNotFound(err) {
				var wbcCR3 examplecomv1.WBConnection
				err = k8sClient.Get(ctx, client.ObjectKey{
					Name:      "wbconntest5-wbconnection-decode-main-filter-resize-low-infer-main",
					Namespace: TESTNAMESPACE,
				},
					&wbcCR3)
				Expect(err).NotTo(HaveOccurred())
				Expect(writer.String()).To(ContainSubstring("WBConnectionCR state transition success."))
				if wbcCR3.Status.Status == examplecomv1.WBDeployStatusTerminating {
					_, err = reconciler.Reconcile(ctx,
						ctrl.Request{NamespacedName: types.NamespacedName{
							Namespace: TESTNAMESPACE,
							Name:      "wbconntest5-wbconnection-decode-main-filter-resize-low-infer-main",
						}})
					Expect(err).NotTo(HaveOccurred())
					Expect(writer.String()).To(ContainSubstring("WBConnectionCR state transition success."))
				} else {
					fmt.Println("WBConnection is not Terminating.")
				}
			} else {
				fmt.Println("Failed to delete PCIeCR.")
			}
		})

		It("Test_5-1-6_Delete_WBFunction_pat③", func() {
			// Create WBConnectionCR
			err = createWBConnection(ctx, wbconnection6PCIe)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntest6-wbconnection-decode-main-filter-resize-low-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{Requeue: false}))
			Expect(err).NotTo(HaveOccurred())

			var wbcCR examplecomv1.WBConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest6-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&wbcCR)

			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status
			expectedStatus := examplecomv1.WBConnectionStatus{
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "wbconntest6",
					Namespace: "default",
				},
				Status: "Waiting",
				From: examplecomv1.FromToWBFunction{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "wbconntest6-wbfunction-decode-main",
						Namespace: "default",
					},
				},
				To: examplecomv1.FromToWBFunction{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "wbconntest6-wbfunction-filter-resize-low-infer-main",
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
			Expect(writer.String()).To(ContainSubstring("name :wbconntest6-wbconnection-decode-main-filter-resize-low-infer-main"))
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

			// check PCIeCR
			var pcieCR controllertestpcie.PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest6-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR)

			Expect(err).NotTo(HaveOccurred())

			expectedCR := controllertestpcie.PCIeConnection{
				Spec: controllertestpcie.PCIeConnectionSpec{
					DataFlowRef: controllertestpcie.WBNamespacedName{
						Name:      "wbconntest6",
						Namespace: "default",
					},
					From: controllertestpcie.PCIeFunctionSpec{
						WBFunctionRef: controllertestpcie.WBNamespacedName{
							Name:      "wbconntest6-wbfunction-decode-main",
							Namespace: "default",
						},
					},
					To: controllertestpcie.PCIeFunctionSpec{
						WBFunctionRef: controllertestpcie.WBNamespacedName{
							Name:      "wbconntest6-wbfunction-filter-resize-low-infer-main",
							Namespace: "default",
						},
					},
				},
			}
			Expect(pcieCR.Spec).To(Equal(expectedCR.Spec))

			var strtime metav1.Time
			strtime = metav1.Now()
			pcieCR.Status.StartTime = strtime
			pcieCR.Status.Status = "Running"
			err = updatePCIeConnection(ctx, pcieCR)
			Expect(err).NotTo(HaveOccurred())

			_, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntest6-wbconnection-decode-main-filter-resize-low-infer-main",
				}})
			Expect(err).NotTo(HaveOccurred())

			var wbcCR2 examplecomv1.WBConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest6-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&wbcCR2)
			Expect(err).NotTo(HaveOccurred())
			Expect(wbcCR2.Status.Status).To(Equal(examplecomv1.WBDeployStatusDeployed))

			var pcieCR2 controllertestpcie.PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest6-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&pcieCR2)
			Expect(err).NotTo(HaveOccurred())

			controllerutil.AddFinalizer(&pcieCR2, "finalizer")
			err = updatePCIeConnection(ctx, pcieCR2)
			Expect(err).NotTo(HaveOccurred())
			err = deletePCIeConnection(ctx, pcieCR2)
			Expect(err).NotTo(HaveOccurred())
			err = deleteWBConnection(ctx, wbcCR2)
			Expect(err).NotTo(HaveOccurred())

			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntest6-wbconnection-decode-main-filter-resize-low-infer-main",
				}})
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))

			var wbcCR3 examplecomv1.WBConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest6-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&wbcCR3)
			Expect(err).NotTo(HaveOccurred())
			Expect(wbcCR3.Status.Status).To(Equal(examplecomv1.WBDeployStatusTerminating))
		})

		It("Test_5-2-1_PCIe", func() {
			// Create WBConnectionCR
			err := createWBConnection(ctx, wbconnection7PCIe)
			Expect(err).NotTo(HaveOccurred())

			createPCIeConnection(ctx, pcieconnection)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntest7-wbconnection-decode-main-filter-resize-low-infer-main",
				}})
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(ctrl.Result{Requeue: true}))

			var wbcCR examplecomv1.WBConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest7-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&wbcCR)

			//check logs
			Expect(writer.String()).To(ContainSubstring("Reconcile start."))
			Expect(writer.String()).To(ContainSubstring("Finalizername=WBConnection.finalizers.example.com.v1"))

			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(ctrl.Result{Requeue: true}))
		})

		It("Test_5-2-2_PCIe", func() {
			// Create WBConnectionCR
			err := createWBConnection(ctx, wbconnection8PCIe)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntest8-wbconnection-decode-main-filter-resize-low-infer-main",
				}})
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(ctrl.Result{Requeue: false}))

			var wbcCR examplecomv1.WBConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest8-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&wbcCR)

			Expect(err).NotTo(HaveOccurred())
			// Expect(got).To(Equal(ctrl.Result{Requeue: false}))

			// test for CR.Status
			expectedStatus := examplecomv1.WBConnectionStatus{
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "wbconntest8",
					Namespace: "default",
				},
				Status: "Waiting",
				From: examplecomv1.FromToWBFunction{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "wbconntest8-wbfunction-decode-main",
						Namespace: "default",
					},
				},
				To: examplecomv1.FromToWBFunction{
					WBFunctionRef: examplecomv1.WBNamespacedName{
						Name:      "wbconntest8-wbfunction-filter-resize-low-infer-main",
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

			err = deleteWBConnection(ctx, wbcCR)
			Expect(err).NotTo(HaveOccurred())

			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "wbconntest8-wbconnection-decode-main-filter-resize-low-infer-main",
				}})

			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(ctrl.Result{Requeue: true}))

			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "wbconntest8-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&wbcCR)

			Expect(err).NotTo(HaveOccurred())

			// check logs
			Expect(writer.String()).To(ContainSubstring("Reconcile start."))
			Expect(writer.String()).To(ContainSubstring("Finalizername=WBConnection.finalizers.example.com.v1"))
			Expect(writer.String()).To(ContainSubstring("Maked Connection does not exist."))
			Expect(writer.String()).To(ContainSubstring("CustomResource Create."))
			Expect(writer.String()).To(ContainSubstring("kind :PCIeConnection"))
			Expect(writer.String()).To(ContainSubstring("apiVersion :example.com/v1"))
			Expect(writer.String()).To(ContainSubstring("name :wbconntest8-wbconnection-decode-main-filter-resize-low-infer-main"))
			Expect(writer.String()).To(ContainSubstring("namespace :default"))
			Expect(writer.String()).To(ContainSubstring("Status Information Change start."))
			Expect(writer.String()).To(ContainSubstring("Finalizername=WBConnection.finalizers.example.com.v1"))
			Expect(writer.String()).To(ContainSubstring("Status Update."))
			Expect(writer.String()).To(ContainSubstring("Status Information Change end."))
		})
	})
})

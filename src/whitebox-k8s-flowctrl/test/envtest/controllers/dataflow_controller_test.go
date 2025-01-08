/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers_test

import (
	"context"
	"time"

	ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	. "github.com/compsysg/whitebox-k8s-flowctrl/internal/controller"
	. "github.com/compsysg/whitebox-k8s-flowctrl/test/lib/testutils"
)

var ctx = context.Background()

var _ = Describe("Function Check (Reconcile)", func() {
	var stopFunc context.CancelFunc
	BeforeEach(func() {
		DataflowDeleteAll()

		// TODO: Register FT and FunctionInfo
		RegisterDataFlowWithScheduleInfo()

		StartReconciler(&stopFunc)
	})
	AfterEach(func() { stopFunc() })

	Context("Verification of Registered Values", func() {
		BeforeEach(func() {
			// // TODO: Register FT and FunctionInfo
			// Register DataFlow with ScheduleInfo set()
		})

		It("Get WBConnection", func() {
			tmp := ntthpcv1.DataFlow{}
			Eventually(func() error {
				return k8sClient.Get(ctx, client.ObjectKey{Namespace: "default", Name: "test"}, &tmp)
			}).Should(Succeed())
			Expect(tmp.Namespace).To(Equal("default"))
			// GinkgoWriter.Println(tmp)

			time.Sleep(time.Second * 3)

			// TODO: FT and FunctionInfo need to be registered so that DF Controller can create WBConnection
			// wbc := ntthpcv1.WBConnection{}
			// Eventually(func() error {
			// 	return k8sClient.Get(ctx, client.ObjectKey{Namespace: "default", Name: "test-wbconnection-decode-main-filter-resize-high-infer-main"}, &wbc)
			// }).Should(Succeed())
			// Expect(wbc.Status.Status).To(Equal("Deployed"))

		})

	})
})

func getManager(mgr ctrl.Manager) (ctrl.Manager, error) {
	if mgr == nil {
		return ctrl.NewManager(cfg, ctrl.Options{
			Scheme:                 testScheme,
			HealthProbeBindAddress: ":9999",
		})
	}
	return mgr, nil
}

// DF controller + WBC controller (for testing)
func StartReconciler(stopFunc *context.CancelFunc) {
	time.Sleep(100 * time.Millisecond)

	var mgr ctrl.Manager
	var err error

	mgr, err = getManager(mgr)
	Expect(err).ToNot(HaveOccurred())

	reconciler := DataFlowReconciler{
		Client:   k8sClient,
		Scheme:   testScheme,
		Recorder: mgr.GetEventRecorderFor("dataflow-controller-test"),
	}
	err = reconciler.SetupWithManager(mgr)
	Expect(err).NotTo(HaveOccurred())

	wbcreconciler := WBConnectionReconciler{
		Client:   k8sClient,
		Scheme:   testScheme,
		Recorder: mgr.GetEventRecorderFor("wbconnection-controller-test"),
	}
	err = wbcreconciler.SetupWithManager(mgr)
	Expect(err).NotTo(HaveOccurred())

	ctx, cancel := context.WithCancel(ctx)
	*stopFunc = cancel
	go func() {
		err := mgr.Start(ctx)
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
}

func RegisterDataFlowWithScheduleInfo() {
	err := CreateDataFlow(ctx, df5, k8sClient)
	Expect(err).NotTo(HaveOccurred())
}

func DataFlowDelete(df_name string) {
	tmp := ntthpcv1.DataFlow{
		ObjectMeta: metav1.ObjectMeta{
			Name:      df_name,
			Namespace: "default",
		},
	}
	err := k8sClient.Delete(ctx, &tmp)
	Expect(err).NotTo(HaveOccurred())
}

func DataflowDeleteAll() {
	err := k8sClient.DeleteAllOf(ctx, &ntthpcv1.DataFlow{}, client.InNamespace("default"))
	Expect(err).NotTo(HaveOccurred())
}

// WBConnectionReconciler reconciles a DataFlow object
type WBConnectionReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

func (r *WBConnectionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	var wbconnection ntthpcv1.WBConnection

	l.Info("fetching WBConnection Resource")
	if err := r.Get(ctx, req.NamespacedName, &wbconnection); err != nil {
		l.Error(err, "unable to fetch WBConnection")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// TODO: Deploy if both fromFunction and toFunction are present

	wbconnection.Status.Status = "Deployed"
	l.Info("Update WBConnection status to Deploy")
	if err := r.Status().Update(ctx, &wbconnection); err != nil {
		l.Error(err, "unable to update wbconnection status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *WBConnectionReconciler) SetupWithManager(mgr ctrl.Manager) error {

	return ctrl.NewControllerManagedBy(mgr).
		For(&ntthpcv1.WBConnection{}).
		Complete(r)
}

// WBFunctionReconciler reconciles a DataFlow object
type WBFunctionReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

func (r *WBFunctionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	var wbfunction ntthpcv1.WBFunction

	l.Info("fetching WBFunction Resource")
	if err := r.Get(ctx, req.NamespacedName, &wbfunction); err != nil {
		l.Error(err, "unable to fetch WBFunction")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	wbfunction.Status.Status = "Deployed"
	l.Info("Update WBFunction status to Deploy")
	if err := r.Status().Update(ctx, &wbfunction); err != nil {
		l.Error(err, "unable to update wbfunction status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *WBFunctionReconciler) SetupWithManager(mgr ctrl.Manager) error {

	return ctrl.NewControllerManagedBy(mgr).
		For(&ntthpcv1.WBFunction{}).
		Complete(r)
}

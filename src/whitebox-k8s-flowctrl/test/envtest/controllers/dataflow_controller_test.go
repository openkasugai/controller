/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED

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
	"bytes"
	"context"
	"fmt"
	"path/filepath"

	ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1" //nolint:stylecheck // ST1019: intentional import as another name
	corev1 "k8s.io/api/core/v1"                                 //nolint:stylecheck // ST1019: intentional import as another name

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime" //nolint:stylecheck // ST1019: intentional import as another name
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	. "github.com/compsysg/whitebox-k8s-flowctrl/internal/controller"
	logf "sigs.k8s.io/controller-runtime/pkg/log" //nolint:stylecheck // ST1019: intentional import as another name

	. "github.com/compsysg/whitebox-k8s-flowctrl/test/lib/testutils"
	"github.com/go-logr/logr"
)

var _ = Describe("DataflowController", func() {

	Context("Test for add finalier", func() {
		var mgr ctrl.Manager
		var reconciler DataFlowReconciler
		var ctx context.Context
		var buf *bytes.Buffer
		var logger logr.Logger
		var act *ntthpcv1.DataFlow = &ntthpcv1.DataFlow{}

		BeforeEach(func() {
			buf = &bytes.Buffer{}
			logger = zap.New(zap.WriteTo(buf), zap.UseDevMode(true))
			ctx = logf.IntoContext(context.Background(), logger)
			dfCommonBeforeEach(ctx, mgr, &reconciler)
		})

		AfterEach(func() {
			ctx = context.Background()
			dfCommonAfterEach(ctx, mgr, &reconciler)
		})

		It("add finalizer", func() {

			_, err := dfCommonSetup(ctx, "no", "new", "")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = reconciler.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "sample-flow1",
						Namespace: "default",
					},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).ShouldNot(BeEmpty())
			Expect(buf.String()).To(ContainSubstring("update DataFlow to add Finalizer"))
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(
					ctx, client.ObjectKey{Name: "sample-flow1", Namespace: "default"}, act)).To(Succeed())
				g.Expect(act.ObjectMeta.Finalizers).To(ContainElement(ContainSubstring("dataflow.finalizers.example.com.v1")))
			}).Should(Succeed())
		})

		It("don't add finalizer", func() {

			_, err := dfCommonSetup(ctx, "no", "deployed", "")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = reconciler.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "sample-flow1",
						Namespace: "default",
					},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).ShouldNot(BeEmpty())
			Expect(buf.String()).NotTo(ContainSubstring("update DataFlow to add Finalizer"))
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(
					ctx, client.ObjectKey{Name: "sample-flow1", Namespace: "default"}, act)).To(Succeed())
				g.Expect(act.ObjectMeta.Finalizers).To(ContainElement(ContainSubstring("dataflow.finalizers.example.com.v1")))
			}).Should(Succeed())
		})
	})

	Context("Test for update DataFlow and WBResource Creation", func() {
		var mgr ctrl.Manager
		var reconciler DataFlowReconciler
		var ctx context.Context
		var buf *bytes.Buffer
		var logger logr.Logger
		var act *ntthpcv1.DataFlow = &ntthpcv1.DataFlow{}
		var wbFunc *ntthpcv1.WBFunction = &ntthpcv1.WBFunction{}
		var wbConn *ntthpcv1.WBConnection = &ntthpcv1.WBConnection{}

		BeforeEach(func() {
			buf = &bytes.Buffer{}
			logger = zap.New(zap.WriteTo(buf), zap.UseDevMode(true))
			ctx = logf.IntoContext(context.Background(), logger)
			dfCommonBeforeEach(ctx, mgr, &reconciler)
		})

		AfterEach(func() {
			ctx = context.Background()
			dfCommonAfterEach(ctx, mgr, &reconciler)
		})

		It("DataFlow.Status.Status is blank", func() {

			exp, err := dfCommonSetup(ctx, "no", "new", "scheduling")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = reconciler.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "sample-flow1",
						Namespace: "default",
					},
				},
			)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).ShouldNot(BeEmpty())
			Expect(buf.String()).To(ContainSubstring("Update DataFlow status to Scheduling in progress"))
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1", Namespace: "default"}, act)).To(Succeed())
				g.Expect(exp.Spec).To(Equal(act.Spec))
				g.Expect(exp.Status.FunctionChain).To(Equal(act.Status.FunctionChain))
				g.Expect(exp.Status.FunctionType).To(ConsistOf(act.Status.FunctionType))
			}).Should(Succeed())
		})

		It("DataFlow.Status.Status is WBFunction/WBConnection creation in progress", func() {

			exp, err := dfCommonSetup(ctx, "no", "creating", "created")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = reconciler.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "sample-flow1",
						Namespace: "default",
					},
				},
			)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).ShouldNot(BeEmpty())
			Expect(buf.String()).To(ContainSubstring("Update DataFlow status to WBFunction/WBConnection created"))
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1", Namespace: "default"}, act)).To(Succeed())
				g.Expect(Check(exp, act, ".status", ".spec")).Should(Succeed())
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1-wbfunction-decode-main", Namespace: "default"}, wbFunc)).To(Succeed())
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1-wbfunction-filter-resize-main", Namespace: "default"}, wbFunc)).To(Succeed())
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1-wbfunction-copy-branch-main", Namespace: "default"}, wbFunc)).To(Succeed())
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1-wbfunction-high-infer-main-1", Namespace: "default"}, wbFunc)).To(Succeed())
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1-wbfunction-high-infer-main-2", Namespace: "default"}, wbFunc)).To(Succeed())
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1-wbconnection-wb-start-of-chain-decode-main", Namespace: "default"}, wbConn)).To(Succeed())
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1-wbconnection-decode-main-filter-resize-main", Namespace: "default"}, wbConn)).To(Succeed())
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1-wbconnection-filter-resize-main-copy-branch-main", Namespace: "default"}, wbConn)).To(Succeed())
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1-wbconnection-copy-branch-main-high-infer-main-1", Namespace: "default"}, wbConn)).To(Succeed())
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1-wbconnection-copy-branch-main-high-infer-main-2", Namespace: "default"}, wbConn)).To(Succeed())
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1-wbconnection-high-infer-main-1-wb-end-of-chain-1", Namespace: "default"}, wbConn)).To(Succeed())
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1-wbconnection-high-infer-main-2-wb-end-of-chain-2", Namespace: "default"}, wbConn)).To(Succeed())
			}).Should(Succeed())
		})

		It("DataFlow.Status.Status is WBFunction/WBConnection created", func() {

			exp, err := dfCommonSetup(ctx, "deployed", "created", "deployed")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = reconciler.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "sample-flow1",
						Namespace: "default",
					},
				},
			)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).ShouldNot(BeEmpty())
			Expect(buf.String()).To(ContainSubstring("Update DataFlow status to Deployed"))
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1", Namespace: "default"}, act)).To(Succeed())
				g.Expect(Check(exp, act, ".status", ".spec")).Should(Succeed())
			}).Should(Succeed())
		})

		It("DataFlow.Status.Status is other", func() {

			_, err := dfCommonSetup(ctx, "no", "scheduling", "")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = reconciler.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "sample-flow1",
						Namespace: "default",
					},
				},
			)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).NotTo(ContainSubstring("fetching FunctionChain Resource:"))
			Expect(buf.String()).NotTo(ContainSubstring("fetching FunctionTargetList"))
			Expect(buf.String()).NotTo(ContainSubstring("check WBfunction/WBConnection"))
		})
	})

	Context("Test for deleteDataFlow function execution conditon", func() {
		var mgr ctrl.Manager
		var reconciler DataFlowReconciler
		var ctx context.Context
		var buf *bytes.Buffer
		var logger logr.Logger
		var act *ntthpcv1.DataFlow = &ntthpcv1.DataFlow{}

		BeforeEach(func() {
			buf = &bytes.Buffer{}
			logger = zap.New(zap.WriteTo(buf), zap.UseDevMode(true))
			ctx = logf.IntoContext(context.Background(), logger)
			dfCommonBeforeEach(ctx, mgr, &reconciler)
		})

		AfterEach(func() {
			ctx = context.Background()
			dfCommonAfterEach(ctx, mgr, &reconciler)
		})

		It("execute deleteDataFlow function", func() {

			_, err := dfCommonSetup(ctx, "no", "creating", "")
			Expect(err).ShouldNot(HaveOccurred())

			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1", Namespace: "default"}, act)).To(Succeed())
				g.Expect(k8sClient.Delete(ctx, act)).To(Succeed())
			}).Should(Succeed())

			_, err = reconciler.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "sample-flow1",
						Namespace: "default",
					},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).ShouldNot(BeEmpty())
			Expect(buf.String()).To(ContainSubstring("execute deleteDataFlow function"))
		})

		It("don't execute deleteDataFlow function", func() {

			_, err := dfCommonSetup(ctx, "no", "deployed", "")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = reconciler.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "sample-flow1",
						Namespace: "default",
					},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).ShouldNot(BeEmpty())
			Expect(buf.String()).NotTo(ContainSubstring("execute deleteDataFlow function"))
		})
	})

	Context("Test for deleteDataFlow function", func() {
		var mgr ctrl.Manager
		var reconciler DataFlowReconciler
		var ctx context.Context
		var buf *bytes.Buffer
		var logger logr.Logger
		var result reconcile.Result
		var act *ntthpcv1.DataFlow = &ntthpcv1.DataFlow{}

		BeforeEach(func() {
			buf = &bytes.Buffer{}
			logger = zap.New(zap.WriteTo(buf), zap.UseDevMode(true))
			ctx = logf.IntoContext(context.Background(), logger)
			dfCommonBeforeEach(ctx, mgr, &reconciler)
		})

		AfterEach(func() {
			ctx = context.Background()
			dfCommonAfterEach(ctx, mgr, &reconciler)
		})

		It("without WBConnections and WBFunctions", func() {

			_, err := dfCommonSetup(ctx, "no", "scheduling", "")
			Expect(err).ShouldNot(HaveOccurred())

			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1", Namespace: "default"}, act)).To(Succeed())
				g.Expect(k8sClient.Delete(ctx, act)).To(Succeed())
			}).Should(Succeed())

			_, err = reconciler.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "sample-flow1",
						Namespace: "default",
					},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).ShouldNot(BeEmpty())
			Expect(buf.String()).To(ContainSubstring("execute deleteDataFlow function"))
			Expect(buf.String()).NotTo(ContainSubstring("delete request for WBConnections"))
			Expect(buf.String()).NotTo(ContainSubstring("delete request for WBFunctions"))
			Expect(buf.String()).To(ContainSubstring("update DataFlow to remove Finalizer"))
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1", Namespace: "default"}, act)).NotTo(Succeed())
			}).Should(Succeed())
		})

		It("delete requests for WBConnections and WBFunctions", func() {

			_, err := dfCommonSetup(ctx, "with", "deployed", "")
			Expect(err).ShouldNot(HaveOccurred())

			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1", Namespace: "default"}, act)).To(Succeed())
				g.Expect(k8sClient.Delete(ctx, act)).To(Succeed())
			}).Should(Succeed())

			_, err = reconciler.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "sample-flow1",
						Namespace: "default",
					},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).ShouldNot(BeEmpty())
			Expect(buf.String()).To(ContainSubstring("execute deleteDataFlow function"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBConnections"))
			Expect(buf.String()).To(MatchRegexp("(?s)delete request for WBConnection sample-flow1-wbconnection-wb-start-of-chain-decode-main.*delete request for WBConnection sample-flow1-wbconnection-decode-main-filter-resize-main.*delete request for WBConnection sample-flow1-wbconnection-filter-resize-main-copy-branch-main.*delete request for WBConnection sample-flow1-wbconnection-copy-branch-main-high-infer-main-1.*delete request for WBConnection sample-flow1-wbconnection-copy-branch-main-high-infer-main-2.*delete request for WBConnection sample-flow1-wbconnection-high-infer-main-1-wb-end-of-chain-1.*delete request for WBConnection sample-flow1-wbconnection-high-infer-main-2-wb-end-of-chain-2"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBFunctions"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBFunction sample-flow1-wbfunction-decode-main"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBFunction sample-flow1-wbfunction-filter-resize-main"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBFunction sample-flow1-wbfunction-copy-branch-main"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBFunction sample-flow1-wbfunction-high-infer-main-2"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBFunction sample-flow1-wbfunction-high-infer-main-1"))
			Expect(buf.String()).To(ContainSubstring("update DataFlow to remove Finalizer"))
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1", Namespace: "default"}, act)).NotTo(Succeed())
			}).Should(Succeed())
		})

		It("check requeue execution for WBConnection deletion request", func() {

			_, err := dfCommonSetup(ctx, "finalizer_conn", "deployed", "")
			Expect(err).ShouldNot(HaveOccurred())

			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1", Namespace: "default"}, act)).To(Succeed())
				g.Expect(k8sClient.Delete(ctx, act)).To(Succeed())
			}).Should(Succeed())

			req := ctrl.Request{NamespacedName: types.NamespacedName{Name: "sample-flow1", Namespace: "default"}}

			Eventually(func(g Gomega) {
				result, err = reconciler.Reconcile(ctx, req)
				g.Expect(RemoveFinalizers(ctx, k8sClient, "default", &ntthpcv1.WBConnection{})).Should(Succeed())
				g.Expect(result.Requeue).To(BeFalse())
			}).Should(Succeed())

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).ShouldNot(BeEmpty())
			Expect(buf.String()).To(ContainSubstring("execute deleteDataFlow function"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBConnections"))
			Expect(buf.String()).To(ContainSubstring("waiting for WBConnection sample-flow1-wbconnection-copy-branch-main-high-infer-main-1 to be deleted"))
			Expect(buf.String()).To(MatchRegexp("(?s)delete request for WBConnection sample-flow1-wbconnection-wb-start-of-chain-decode-main.*delete request for WBConnection sample-flow1-wbconnection-decode-main-filter-resize-main.*delete request for WBConnection sample-flow1-wbconnection-filter-resize-main-copy-branch-main.*delete request for WBConnection sample-flow1-wbconnection-copy-branch-main-high-infer-main-1.*delete request for WBConnection sample-flow1-wbconnection-copy-branch-main-high-infer-main-2.*delete request for WBConnection sample-flow1-wbconnection-high-infer-main-1-wb-end-of-chain-1.*delete request for WBConnection sample-flow1-wbconnection-high-infer-main-2-wb-end-of-chain-2"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBFunctions"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBFunction sample-flow1-wbfunction-decode-main"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBFunction sample-flow1-wbfunction-filter-resize-main"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBFunction sample-flow1-wbfunction-copy-branch-main"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBFunction sample-flow1-wbfunction-high-infer-main-2"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBFunction sample-flow1-wbfunction-high-infer-main-1"))
			Expect(buf.String()).To(ContainSubstring("update DataFlow to remove Finalizer"))
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1", Namespace: "default"}, act)).NotTo(Succeed())
			}).Should(Succeed())
		})

		It("check requeue execution for WBFunction deletion request", func() {

			_, err := dfCommonSetup(ctx, "finalizer_func", "deployed", "")
			Expect(err).ShouldNot(HaveOccurred())

			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1", Namespace: "default"}, act)).To(Succeed())
				g.Expect(k8sClient.Delete(ctx, act)).To(Succeed())
			}).Should(Succeed())

			req := ctrl.Request{NamespacedName: types.NamespacedName{Name: "sample-flow1", Namespace: "default"}}

			Eventually(func(g Gomega) {
				result, err = reconciler.Reconcile(ctx, req)
				g.Expect(RemoveFinalizers(ctx, k8sClient, "default", &ntthpcv1.WBFunction{})).Should(Succeed())
				g.Expect(result.Requeue).To(BeFalse())
			}).Should(Succeed())

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).ShouldNot(BeEmpty())
			Expect(buf.String()).To(ContainSubstring("execute deleteDataFlow function"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBConnections"))
			Expect(buf.String()).To(ContainSubstring("waiting for WBFunction sample-flow1-wbfunction-copy-branch-main to be deleted"))
			Expect(buf.String()).To(MatchRegexp("(?s)delete request for WBConnection sample-flow1-wbconnection-wb-start-of-chain-decode-main.*delete request for WBConnection sample-flow1-wbconnection-decode-main-filter-resize-main.*delete request for WBConnection sample-flow1-wbconnection-filter-resize-main-copy-branch-main.*delete request for WBConnection sample-flow1-wbconnection-copy-branch-main-high-infer-main-1.*delete request for WBConnection sample-flow1-wbconnection-copy-branch-main-high-infer-main-2.*delete request for WBConnection sample-flow1-wbconnection-high-infer-main-1-wb-end-of-chain-1.*delete request for WBConnection sample-flow1-wbconnection-high-infer-main-2-wb-end-of-chain-2"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBFunctions"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBFunction sample-flow1-wbfunction-decode-main"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBFunction sample-flow1-wbfunction-filter-resize-main"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBFunction sample-flow1-wbfunction-copy-branch-main"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBFunction sample-flow1-wbfunction-high-infer-main-2"))
			Expect(buf.String()).To(ContainSubstring("delete request for WBFunction sample-flow1-wbfunction-high-infer-main-1"))
			Expect(buf.String()).To(ContainSubstring("update DataFlow to remove Finalizer"))
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1", Namespace: "default"}, act)).NotTo(Succeed())
			}).Should(Succeed())
		})
	})
})

func dfCommonBeforeEach(ctx context.Context, mgr ctrl.Manager, reconciler *DataFlowReconciler) {

	Eventually(func(g Gomega) {
		var err error
		mgr, err = getMgr(mgr)
		addr++
		Expect(err).ShouldNot(HaveOccurred())
	}).Should(Succeed())

	// If the following is not executed, r.Get, r.List, etc. cannot be executed.
	reconciler.Client = k8sClient
	reconciler.Scheme = testScheme
	reconciler.Recorder = mgr.GetEventRecorderFor("dataflow-controller")
	Expect(reconciler.SetupWithManager(mgr)).Should(Succeed())
}

func dfCommonSetup(ctx context.Context, wbRscFlg string, dataflowName string, expType string) (*ntthpcv1.DataFlow, error) {

	var yamls []string
	var exp *ntthpcv1.DataFlow

	commonDir := filepath.Join("..", "..", "resources", "controllers_test", "dataflow", fmt.Sprintf("%s_wbresource", wbRscFlg))

	if wbRscFlg == "no" {
		yamls = []string{
			filepath.Join(commonDir, fmt.Sprintf("functioninfos")),
			filepath.Join(commonDir, fmt.Sprintf("functiontypes")),
			filepath.Join(commonDir, fmt.Sprintf("functionchain.yaml")),
			filepath.Join(commonDir, fmt.Sprintf("dataflow_%s.yaml", dataflowName)),
		}
	} else {
		yamls = []string{
			filepath.Join(commonDir, fmt.Sprintf("wbfunctions")),
			filepath.Join(commonDir, fmt.Sprintf("wbconnections")),
			filepath.Join(commonDir, fmt.Sprintf("dataflow_%s.yaml", dataflowName)),
		}
	}
	_, err := DeployWithObjectMeta(ctx, k8sClient, yamls...)

	if err != nil {
		return nil, err
	}

	if expType != "" {
		expYaml := filepath.Join(commonDir, fmt.Sprintf("dataflow_%s.yaml", expType))
		exp, err = GetResourceFromYaml[*ntthpcv1.DataFlow](expYaml)
		if err != nil {
			return nil, err
		}
	}

	return exp, nil

}

func dfCommonAfterEach(ctx context.Context, mgr ctrl.Manager, reconciler *DataFlowReconciler) {

	Eventually(func(g Gomega) {
		var err error
		mgr, err = getMgr(mgr)
		addr++
		Expect(err).ShouldNot(HaveOccurred())
	}).Should(Succeed())

	// If the following is not executed, r.Get, r.List, etc. cannot be executed.
	reconciler.Client = k8sClient
	reconciler.Scheme = testScheme
	reconciler.Recorder = mgr.GetEventRecorderFor("dataflow-controller")

	Expect(RemoveFinalizers(ctx, k8sClient, "default",
		&corev1.ConfigMap{},
		&ntthpcv1.DataFlow{},
		&ntthpcv1.FunctionType{},
		&ntthpcv1.FunctionChain{},
		&ntthpcv1.WBFunction{},
		&ntthpcv1.WBConnection{},
	)).Should(Succeed())

	Expect(DeleteAllOf(ctx, k8sClient, "default",
		&corev1.ConfigMap{},
		&ntthpcv1.DataFlow{},
		&ntthpcv1.FunctionType{},
		&ntthpcv1.FunctionChain{},
		&ntthpcv1.WBFunction{},
		&ntthpcv1.WBConnection{},
	)).Should(Succeed())
}

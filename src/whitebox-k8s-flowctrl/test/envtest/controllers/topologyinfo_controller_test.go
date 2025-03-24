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
	"github.com/go-logr/logr"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1" //nolint:stylecheck // ST1019: intentional import as another name
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime" //nolint:stylecheck // ST1019: intentional import as another name
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	. "github.com/compsysg/whitebox-k8s-flowctrl/internal/controller"
	logf "sigs.k8s.io/controller-runtime/pkg/log" //nolint:stylecheck // ST1019: intentional import as another name

	. "github.com/compsysg/whitebox-k8s-flowctrl/test/lib/testutils"
)

var _ = Describe("TopologyInfoController", func() {

	Context("Test for TopologyInfo Resource", func() {
		var mgr ctrl.Manager
		var reconciler TopologyInfoReconciler
		var ctx context.Context
		var buf *bytes.Buffer
		var logger logr.Logger
		var actTpCr *ntthpcv1.TopologyInfo = &ntthpcv1.TopologyInfo{}
		var actTpCm *corev1.ConfigMap = &corev1.ConfigMap{}

		BeforeEach(func() {
			buf = &bytes.Buffer{}
			logger = zap.New(zap.WriteTo(buf), zap.UseDevMode(true))
			ctx = logf.IntoContext(context.Background(), logger)
			tpCommonBeforeEach(ctx, mgr, &reconciler)
		})

		AfterEach(func() {
			ctx = context.Background()
			tpCommonAfterEach(ctx, mgr, &reconciler)
		})

		It("create TopologyInfo Resource", func() {

			expTpCr, err := tpCommonSetup1(ctx, "create", "created")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = reconciler.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "topologydata",
						Namespace: "default",
					},
				},
			)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).ShouldNot(BeEmpty())
			Expect(buf.String()).To(ContainSubstring(("execute createOrUpdateTopologyInfo function")))
			Expect(buf.String()).To(ContainSubstring(("TopologyInfo does not exists yet, so create new")))
			Expect(buf.String()).To(ContainSubstring(("update TopologyInfo status")))
			Expect(buf.String()).To(ContainSubstring(("update TopologyData ConfigMap to add Finalizer")))
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "topologyinfo", Namespace: "default"}, actTpCr)).To(Succeed())
				g.Expect(Check(expTpCr, actTpCr, ".status", ".spec")).Should(Succeed())
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "topologydata", Namespace: "default"}, actTpCm)).To(Succeed())
				g.Expect(actTpCm.ObjectMeta.Finalizers).To(ContainElement(ContainSubstring("example.com.v1/topologyinfo.finalizers")))
			}).Should(Succeed())
		})

		It("update TopologyInfo Resource", func() {

			expTpCr, err := tpCommonSetup1(ctx, "update", "updated")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = reconciler.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "topologydata",
						Namespace: "default",
					},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).ShouldNot(BeEmpty())
			Expect(buf.String()).To(ContainSubstring(("execute createOrUpdateTopologyInfo function")))
			Expect(buf.String()).To(ContainSubstring(("TopologyInfo resource already exists, so update it")))
			Expect(buf.String()).To(ContainSubstring(("update TopologyInfo status")))
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "topologyinfo", Namespace: "default"}, actTpCr)).To(Succeed())
				g.Expect(expTpCr.Status.Entities).To(ConsistOf(actTpCr.Status.Entities))
				g.Expect(expTpCr.Status.Relations).To(ConsistOf(actTpCr.Status.Relations))
			}).Should(Succeed())
		})

		It("delete TopologyInfo Resource", func() {

			_, err := tpCommonSetup1(ctx, "delete", "")
			Expect(err).ShouldNot(HaveOccurred())

			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "topologydata", Namespace: "default"}, actTpCm)).To(Succeed())
				g.Expect(k8sClient.Delete(ctx, actTpCm)).To(Succeed())
			}).Should(Succeed())

			_, err = reconciler.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "topologydata",
						Namespace: "default",
					},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).ShouldNot(BeEmpty())
			Expect(buf.String()).To(ContainSubstring(("delete TopologyInfo resource")))
			Expect(buf.String()).To(ContainSubstring(("update TopologyInfo ConfigMap to remove Finalizer")))
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "topologydata", Namespace: "default"}, actTpCm)).NotTo(Succeed())
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "topologyinfo", Namespace: "default"}, actTpCr)).NotTo(Succeed())
			}).Should(Succeed())
		})
	})

	Context("Test for checkWbConnEvent", func() {
		var mgr ctrl.Manager
		var reconciler TopologyInfoReconciler
		var ctx context.Context
		var buf *bytes.Buffer
		var logger logr.Logger
		var act *ntthpcv1.WBConnection = &ntthpcv1.WBConnection{}

		BeforeEach(func() {
			buf = &bytes.Buffer{}
			logger = zap.New(zap.WriteTo(buf), zap.UseDevMode(true))
			ctx = logf.IntoContext(context.Background(), logger)
			tpCommonBeforeEach(ctx, mgr, &reconciler)
		})

		AfterEach(func() {
			ctx = context.Background()
			tpCommonAfterEach(ctx, mgr, &reconciler)
		})

		It("deletionTimestamp is not zero", func() {

			_, err := tpCommonSetup2(ctx, "processed", "")
			Expect(err).ShouldNot(HaveOccurred())

			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1-wbconnection-decode-main-filter-resize-main", Namespace: "default"}, act)).To(Succeed())
				g.Expect(k8sClient.Delete(ctx, act)).To(Succeed())
			}).Should(Succeed())

			_, err = reconciler.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "sample-flow1-wbconnection-decode-main-filter-resize-main",
						Namespace: "default",
					},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).ShouldNot(BeEmpty())
			Expect(buf.String()).To(ContainSubstring("WBConnection is now deleting, substract TopologyInfo capacity used"))
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1-wbconnection-decode-main-filter-resize-main", Namespace: "default"}, act)).To(Succeed())
				g.Expect(act.ObjectMeta.Finalizers).NotTo(ContainElement(ContainSubstring("topologyinfo.finalizers.example.com.v1")))
			}).Should(Succeed())
		})

		It("deletionTimestamp is zero with ContainsFinalizer", func() {

			_, err := tpCommonSetup2(ctx, "processed", "")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = reconciler.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "sample-flow1-wbconnection-decode-main-filter-resize-main",
						Namespace: "default",
					},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).ShouldNot(BeEmpty())
			Expect(buf.String()).To(ContainSubstring("WBConnection capacity has already added to TopologyInfo, abort Reconcile"))
		})

		It("deletionTimestamp is zero with no ContainsFinalizer and just created", func() {

			_, err := tpCommonSetup2(ctx, "created", "")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = reconciler.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "sample-flow1-wbconnection-decode-main-filter-resize-main",
						Namespace: "default",
					},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).ShouldNot(BeEmpty())
			Expect(buf.String()).To(ContainSubstring("WBConnection is not yet in the Deployed state, abort Reconcile"))
		})

		It("deletionTimestamp is zero with no ContainsFinalizer and deployed", func() {

			_, err := tpCommonSetup2(ctx, "deployed", "")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = reconciler.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "sample-flow1-wbconnection-decode-main-filter-resize-main",
						Namespace: "default",
					},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).ShouldNot(BeEmpty())
			Expect(buf.String()).To(ContainSubstring("WBConnection is now in the Deployed state, add TopologyInfo capacity used"))
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1-wbconnection-decode-main-filter-resize-main", Namespace: "default"}, act)).To(Succeed())
				g.Expect(act.ObjectMeta.Finalizers).To(ContainElement(ContainSubstring("topologyinfo.finalizers.example.com.v1")))
			}).Should(Succeed())
		})
	})

	Context("Test for updateCapacityUsed", func() {
		var mgr ctrl.Manager
		var reconciler TopologyInfoReconciler
		var ctx context.Context
		var buf *bytes.Buffer
		var logger logr.Logger
		var actWbConn *ntthpcv1.WBConnection = &ntthpcv1.WBConnection{}
		var actTopo *ntthpcv1.TopologyInfo = &ntthpcv1.TopologyInfo{}

		BeforeEach(func() {
			buf = &bytes.Buffer{}
			logger = zap.New(zap.WriteTo(buf), zap.UseDevMode(true))
			ctx = logf.IntoContext(context.Background(), logger)
			tpCommonBeforeEach(ctx, mgr, &reconciler)
		})

		AfterEach(func() {
			ctx = context.Background()
			tpCommonAfterEach(ctx, mgr, &reconciler)
		})

		It("add TopologyInfo capacity used", func() {

			exp, err := tpCommonSetup2(ctx, "deployed", "add")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = reconciler.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "sample-flow1-wbconnection-decode-main-filter-resize-main",
						Namespace: "default",
					},
				},
			)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).To(ContainSubstring("WBConnection is now in the Deployed state, add TopologyInfo capacity used"))
			Expect(buf.String()).To(ContainSubstring("update WBConnection to add Finalizer"))
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1-wbconnection-decode-main-filter-resize-main", Namespace: "default"}, actWbConn)).To(Succeed())
				g.Expect(actWbConn.ObjectMeta.Finalizers).To(ContainElement(ContainSubstring("topologyinfo.finalizers.example.com.v1")))
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "topologyinfo", Namespace: "default"}, actTopo)).To(Succeed())
				g.Expect(Check(exp, actTopo, ".status", ".spec")).Should(Succeed())
			}).Should(Succeed())
		})

		It("substract TopologyInfo capacity used", func() {

			exp, err := tpCommonSetup2(ctx, "deleting", "substract")
			Expect(err).ShouldNot(HaveOccurred())

			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1-wbconnection-decode-main-filter-resize-main", Namespace: "default"}, actWbConn)).To(Succeed())
				g.Expect(k8sClient.Delete(ctx, actWbConn)).To(Succeed())
			}).Should(Succeed())

			_, err = reconciler.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "sample-flow1-wbconnection-decode-main-filter-resize-main",
						Namespace: "default",
					},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).To(ContainSubstring("WBConnection is now deleting, substract TopologyInfo capacity used"))
			Expect(buf.String()).To(ContainSubstring("update WBConnection to remove Finalizer"))
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1-wbconnection-decode-main-filter-resize-main", Namespace: "default"}, actWbConn)).To(Succeed())
				g.Expect(actWbConn.ObjectMeta.Finalizers).NotTo(ContainElement(ContainSubstring("topologyinfo.finalizers.example.com.v1")))
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "topologyinfo", Namespace: "default"}, actTopo)).To(Succeed())
				g.Expect(Check(exp, actTopo, ".status", ".spec")).Should(Succeed())
			}).Should(Succeed())
		})

		It("no connectionpath", func() {

			exp, err := tpCommonSetup2(ctx, "no_connectionpath", "not_add")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = reconciler.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "sample-flow1-wbconnection-decode-main-filter-resize-main",
						Namespace: "default",
					},
				},
			)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buf.String()).To(ContainSubstring("WBConnection is now in the Deployed state, add TopologyInfo capacity used"))
			Expect(buf.String()).To(ContainSubstring("WBConnection has no connectionpath information, abort reconcile"))
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "sample-flow1-wbconnection-decode-main-filter-resize-main", Namespace: "default"}, actWbConn)).To(Succeed())
				g.Expect(actWbConn.ObjectMeta.Finalizers).NotTo(ContainElement(ContainSubstring("topologyinfo.finalizers.example.com.v1")))
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "topologyinfo", Namespace: "default"}, actTopo)).To(Succeed())
				g.Expect(Check(exp, actTopo, ".status", ".spec")).Should(Succeed())
			}).Should(Succeed())
		})
	})
})

func tpCommonBeforeEach(ctx context.Context, mgr ctrl.Manager, reconciler *TopologyInfoReconciler) {

	Eventually(func(g Gomega) {
		var err error
		mgr, err = getMgr(mgr)
		addr++
		Expect(err).ShouldNot(HaveOccurred())
	}).Should(Succeed())

	// If the following is not executed, r.Get, r.List, etc. cannot be executed.
	reconciler.Client = k8sClient
	reconciler.Scheme = testScheme
	reconciler.TopologyinfoNamespacedName = types.NamespacedName{
		Name:      "topologyinfo",
		Namespace: "default",
	}
	reconciler.TopologydataNamespacedName = types.NamespacedName{
		Name:      "topologydata",
		Namespace: "default",
	}
	Expect(reconciler.SetupWithManager(mgr)).Should(Succeed())
}

func tpCommonSetup1(ctx context.Context, testType, expType string) (*ntthpcv1.TopologyInfo, error) {

	var err error
	var yamls []string
	var exp *ntthpcv1.TopologyInfo

	commonDir := filepath.Join("..", "..", "resources", "controllers_test", "topologyinfo")

	switch testType {
	case "create":
		yamls = []string{
			filepath.Join(commonDir, fmt.Sprintf("topology_configmap_%s.yaml", testType)),
		}
	case "update":
		yamls = []string{
			filepath.Join(commonDir, fmt.Sprintf("topology_configmap_%s.yaml", testType)),
			filepath.Join(commonDir, "topology_customresource_created.yaml"),
		}
	case "delete":
		yamls = []string{
			filepath.Join(commonDir, fmt.Sprintf("topology_configmap_%s.yaml", testType)),
			filepath.Join(commonDir, "topology_customresource_created.yaml"),
		}
	default:
		return nil, fmt.Errorf("unsupported testType: %s", testType)
	}

	_, err = DeployWithObjectMeta(ctx, k8sClient, yamls...)
	if err != nil {
		return nil, err
	}

	if expType != "" {
		expYaml := filepath.Join(commonDir, fmt.Sprintf("exp_%s.yaml", expType))
		exp, err = GetResourceFromYaml[*ntthpcv1.TopologyInfo](expYaml)
		if err != nil {
			return nil, err
		}
	}

	return exp, nil

}

func tpCommonSetup2(ctx context.Context, testType, expType string) (*ntthpcv1.TopologyInfo, error) {

	var err error
	var exp *ntthpcv1.TopologyInfo

	commonDir := filepath.Join("..", "..", "resources", "controllers_test", "topologyinfo")

	yamls := []string{
		filepath.Join(commonDir, "topology_customresource_normal.yaml"),
		filepath.Join(commonDir, fmt.Sprintf("wbconnection_%s.yaml", testType)),
	}

	_, err = DeployWithObjectMeta(ctx, k8sClient, yamls...)
	if err != nil {
		return nil, err
	}

	if expType != "" {
		expYaml := filepath.Join(commonDir, fmt.Sprintf("exp_%s.yaml", expType))
		exp, err = GetResourceFromYaml[*ntthpcv1.TopologyInfo](expYaml)
		if err != nil {
			return nil, err
		}
	}

	return exp, nil

}

func tpCommonAfterEach(ctx context.Context, mgr ctrl.Manager, reconciler *TopologyInfoReconciler) {

	Eventually(func(g Gomega) {
		var err error
		mgr, err = getMgr(mgr)
		addr++
		Expect(err).ShouldNot(HaveOccurred())
	}).Should(Succeed())

	// If the following is not executed, r.Get, r.List, etc. cannot be executed.
	reconciler.Client = k8sClient
	reconciler.Scheme = testScheme
	reconciler.TopologyinfoNamespacedName = types.NamespacedName{
		Name:      "topologyinfo",
		Namespace: "default",
	}
	reconciler.TopologydataNamespacedName = types.NamespacedName{
		Name:      "topologydata",
		Namespace: "default",
	}

	Expect(RemoveFinalizers(ctx, k8sClient, "default",
		&corev1.ConfigMap{},
		&ntthpcv1.TopologyInfo{},
		&ntthpcv1.WBConnection{},
	)).Should(Succeed())

	Expect(DeleteAllOf(ctx, k8sClient, "default",
		&corev1.ConfigMap{},
		&ntthpcv1.TopologyInfo{},
		&ntthpcv1.WBConnection{},
	)).Should(Succeed())
}

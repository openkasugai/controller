/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED
*/

package controllers_test

import (
	"bytes"
	"context"
	"time"

	ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1" //nolint:stylecheck // ST1019: intentional import as another name
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1" //nolint:stylecheck // ST1019: intentional import as another name
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime" //nolint:stylecheck // ST1019: intentional import as another name
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	. "github.com/compsysg/whitebox-k8s-flowctrl/internal/controller"

	"github.com/go-logr/logr"
	logf "sigs.k8s.io/controller-runtime/pkg/log" //nolint:stylecheck // ST1019: intentional import as another name
)

var _ = Describe("FunctionTypeController", func() {
	var mgr ctrl.Manager
	var err error
	var buf *bytes.Buffer
	var logger logr.Logger
	ctx := context.Background()

	Context("Test for FunctionTypeController", Ordered, func() {
		var reconciler FunctionTypeReconciler

		BeforeAll(func() {

			buf = &bytes.Buffer{}
			logger = zap.New(zap.WriteTo(buf), zap.UseDevMode(true))
			ctx = logf.IntoContext(context.Background(), logger)

			mgr, err = ctrl.NewManager(cfg, ctrl.Options{
				Scheme: testScheme,
			})
			Expect(err).NotTo(HaveOccurred())

			reconciler = FunctionTypeReconciler{
				Client: k8sClient,
				Scheme: testScheme,
			}
			err = reconciler.SetupWithManager(mgr)
			Expect(err).NotTo(HaveOccurred())
		})

		BeforeEach(func() {
			buf.Reset()
			err = k8sClient.DeleteAllOf(ctx, &ntthpcv1.FunctionType{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)
		})

		AfterAll(func() {
			buf.Reset()
		})

		It("success to get FuctionType resource", func() {

			err = createFunctionInfo(ctx, functionInfo9)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionType(ctx, functionType9)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functiontest9"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			Expect(buf.String()).To(ContainSubstring("fetching FunctionType Resource"))

		})

		It("fail to get FuctionType resource because there's no FuctionType applying", func() {

			err = createFunctionInfo(ctx, functionInfo9)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functiontest9"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).ToNot(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			Expect(buf.String()).To(ContainSubstring("unable to fetch FunctionType"))
		})

		It("success to get FuctionInfoCM resource", func() {

			err = createFunctionInfo(ctx, functionInfo9)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionType(ctx, functionType9)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functiontest9"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			Expect(buf.String()).To(ContainSubstring("Finding configmap with name functiontest9 in namespace default"))

		})

		It("fail to get FuctionInfoCM resource because there's no FuctionInfo applying", func() {

			err = createFunctionType(ctx, functionType9)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functiontest9"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).ToNot(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			Expect(buf.String()).To(ContainSubstring("unable to fetch configmap"))

		})

		It("fail to get FuctionInfoCM because namespace is invalid", func() {

			err = createFunctionInfo(ctx, functionInfo9)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionType(ctx, functionType16)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functiontest16"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).ToNot(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			Expect(buf.String()).To(ContainSubstring("unable to fetch configmap"))

		})

		It("status updating when there is one deployableItems", func() {

			err = createFunctionInfo(ctx, functionInfo9)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionType(ctx, functionType9)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functiontest9"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			expected := ntthpcv1.FunctionTypeStatus{
				Status:               "Ready",
				RegionTypeCandidates: []string{"cpu"},
			}

			var ft ntthpcv1.FunctionType
			err = k8sClient.Get(ctx, client.ObjectKey{Name: "functiontest9", Namespace: "default"}, &ft)
			Expect(err).NotTo(HaveOccurred())

			Expect(ft.Status).To(Equal(expected))
		})

		It("status updating when there are multi deployableItems with same regionType", func() {

			err = createFunctionInfo(ctx, functionInfo10)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionType(ctx, functionType10)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functiontest10"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			expected := ntthpcv1.FunctionTypeStatus{
				Status:               "Ready",
				RegionTypeCandidates: []string{"cpu"},
			}

			var ft ntthpcv1.FunctionType
			err = k8sClient.Get(ctx, client.ObjectKey{Name: "functiontest10", Namespace: "default"}, &ft)
			Expect(err).NotTo(HaveOccurred())

			Expect(ft.Status).To(Equal(expected))
		})

		It("status updating when there are deployableItems with different regionType", func() {

			err = createFunctionInfo(ctx, functionInfo11)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionType(ctx, functionType11)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functiontest11"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			expected := ntthpcv1.FunctionTypeStatus{
				Status:               "Ready",
				RegionTypeCandidates: []string{"cpu", "alveo"},
			}

			var ft ntthpcv1.FunctionType
			err = k8sClient.Get(ctx, client.ObjectKey{Name: "functiontest11", Namespace: "default"}, &ft)
			Expect(err).NotTo(HaveOccurred())

			Expect(ft.Status).To(Equal(expected))
		})

		It("status updating when there is a deployableItems but no regionType", func() {

			err = createFunctionInfo(ctx, functionInfo12)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionType(ctx, functionType12)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functiontest12"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			expected := ntthpcv1.FunctionTypeStatus{
				Status:               "Error",
				RegionTypeCandidates: nil,
			}

			var ft ntthpcv1.FunctionType
			err = k8sClient.Get(ctx, client.ObjectKey{Name: "functiontest12", Namespace: "default"}, &ft)
			Expect(err).NotTo(HaveOccurred())

			Expect(ft.Status).To(Equal(expected))
		})

		It("status updating when there is no deployableItems", func() {

			err = createFunctionInfo(ctx, functionInfo13)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionType(ctx, functionType13)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functiontest13"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			expected := ntthpcv1.FunctionTypeStatus{
				Status:               "Error",
				RegionTypeCandidates: nil,
			}

			var ft ntthpcv1.FunctionType
			err = k8sClient.Get(ctx, client.ObjectKey{Name: "functiontest13", Namespace: "default"}, &ft)
			Expect(err).NotTo(HaveOccurred())

			Expect(ft.Status).To(Equal(expected))
		})

		It("status updating when there is a recommened key", func() {

			err = createFunctionInfo(ctx, functionInfo14)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionType(ctx, functionType14)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functiontest14"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			expected := ntthpcv1.FunctionTypeStatus{
				Status:               "Ready",
				RegionTypeCandidates: []string{"cpu"},
			}

			var ft ntthpcv1.FunctionType
			err = k8sClient.Get(ctx, client.ObjectKey{Name: "functiontest14", Namespace: "default"}, &ft)
			Expect(err).NotTo(HaveOccurred())

			Expect(ft.Status).To(Equal(expected))
		})

		It("status updating when there is a spec key", func() {

			err = createFunctionInfo(ctx, functionInfo9)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionType(ctx, functionType9)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functiontest9"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			Expect(buf.String()).To(ContainSubstring("spec key found"))

		})

		It("status updating when there is a default key", func() {

			err = createFunctionInfo(ctx, functionInfo15)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionType(ctx, functionType15)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functiontest15"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			Expect(buf.String()).To(ContainSubstring("unsupported key found"))
		})
	})
})

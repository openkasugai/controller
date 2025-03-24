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

	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime" //nolint:stylecheck // ST1019: intentional import as another name
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	. "github.com/compsysg/whitebox-k8s-flowctrl/internal/controller"

	"github.com/go-logr/logr"
	logf "sigs.k8s.io/controller-runtime/pkg/log" //nolint:stylecheck // ST1019: intentional import as another name
)

func createFunctionTarget(ctx context.Context, ft ntthpcv1.FunctionTarget) error {
	tmp := &ntthpcv1.FunctionTarget{}
	*tmp = ft
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	tmp.Status = ft.Status
	err = k8sClient.Status().Update(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

func createComputeResource(ctx context.Context, cr ntthpcv1.ComputeResource) error {
	tmp := &ntthpcv1.ComputeResource{}
	*tmp = cr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

var _ = Describe("FunctionTargetController", func() {
	var mgr ctrl.Manager
	var err error
	var buf *bytes.Buffer
	var logger logr.Logger
	ctx := context.Background()

	Context("Test for FunctionTargetController", Ordered, func() {
		var reconciler FunctionTargetReconciler

		BeforeAll(func() {

			buf = &bytes.Buffer{}
			logger = zap.New(zap.WriteTo(buf), zap.UseDevMode(true))
			ctx = logf.IntoContext(context.Background(), logger)

			mgr, err = ctrl.NewManager(cfg, ctrl.Options{
				Scheme: testScheme,
			})
			Expect(err).NotTo(HaveOccurred())

			reconciler = FunctionTargetReconciler{
				Client: k8sClient,
				Scheme: testScheme,
			}
			err = reconciler.SetupWithManager(mgr)
			Expect(err).NotTo(HaveOccurred())
		})

		BeforeEach(func() {
			buf.Reset()
			err = k8sClient.DeleteAllOf(ctx, &ntthpcv1.FunctionTarget{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &ntthpcv1.ComputeResource{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)
		})

		AfterAll(func() {
			buf.Reset()
		})

		It("success to get ComputeResource", func() {

			err = createComputeResource(ctx, computeResource1)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functiontargettest1"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			Expect(buf.String()).To(ContainSubstring("fetching ComputeResource Resource"))

		})

		It("fail to get ComputeResource", func() {

			err = createFunctionTarget(ctx, functionTarget1)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functiontargettest1"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			Expect(buf.String()).To(ContainSubstring("unable to fetch ComputeResource"))
		})

		It("success to create or update functiontarget resource", func() {

			err = createComputeResource(ctx, computeResource1)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionTarget(ctx, functionTarget1)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functiontargettest1"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			Expect(buf.String()).To(ContainSubstring("CreateOrUpdate for FunctionTarget Resouce"))
		})

		It("fail to create or update functiontarget resource", func() {

			err = createComputeResource(ctx, computeResource2)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionTarget(ctx, functionTarget2)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functiontargettest2"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).To(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			Expect(buf.String()).To(ContainSubstring("unable to set ownerReference from ComputeResource to FunctionTarget"))
		})

		It("err msg appears because functiontarget is incorrect", func() {

			err = createComputeResource(ctx, computeResource2)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionTarget(ctx, functionTarget2)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functiontargettest2"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).To(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			Expect(buf.String()).To(ContainSubstring("unable to ensure functiontarget is correct"))
		})

		It("status updating content check", func() {

			err = createComputeResource(ctx, computeResource1)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionTarget(ctx, functionTarget1)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functiontargettest1"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			expected := ntthpcv1.FunctionTargetStatus{
				RegionName:  "lane1",
				RegionType:  "alveou250-0100001c-2lanes-0nics",
				NodeName:    "functiontargettest1",
				DeviceType:  "alveou250",
				DeviceIndex: 3,
				Available:   true,
				Status:      "NotReady",
				Functions: []ntthpcv1.FunctionCapStruct{
					{
						FunctionIndex:    2,
						FunctionName:     "filter-resize",
						Available:        true,
						MaxDataFlows:     func(i int32) *int32 { return &i }(8),
						CurrentDataFlows: func(i int32) *int32 { return &i }(1),
						MaxCapacity:      func(i int32) *int32 { return &i }(30),
						CurrentCapacity:  func(i int32) *int32 { return &i }(8),
					},
				},
			}

			var ft ntthpcv1.FunctionTarget
			err = k8sClient.Get(ctx, client.ObjectKey{Name: "functiontargettest1.alveou250-3.lane1", Namespace: "default"}, &ft)
			Expect(err).NotTo(HaveOccurred())

			Expect(ft.Status).To(Equal(expected))

		})

		It("status updating log check", func() {

			err = createComputeResource(ctx, computeResource1)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionTarget(ctx, functionTarget1)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functiontargettest1"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			Expect(buf.String()).To(ContainSubstring("Update FunctionTarget status"))
		})
	})
})

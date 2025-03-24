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

func createFunctionChain(ctx context.Context, fc ntthpcv1.FunctionChain) error {
	tmp := &ntthpcv1.FunctionChain{}
	*tmp = fc
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	tmp.Status = fc.Status
	err = k8sClient.Status().Update(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

func createFunctionType(ctx context.Context, ft ntthpcv1.FunctionType) error {
	tmp := &ntthpcv1.FunctionType{}
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

func createFunctionInfo(ctx context.Context, funcInfo corev1.ConfigMap) error {
	tmp := &corev1.ConfigMap{}
	*tmp = funcInfo
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

var _ = Describe("FunctionChainController", func() {
	var mgr ctrl.Manager
	var err error
	var buf *bytes.Buffer
	var logger logr.Logger
	ctx := context.Background()

	Context("Test for FunctionChainController", Ordered, func() {
		var reconciler FunctionChainReconciler

		BeforeAll(func() {
			buf = &bytes.Buffer{}
			logger = zap.New(zap.WriteTo(buf), zap.UseDevMode(true))
			ctx = logf.IntoContext(context.Background(), logger)

			mgr, err = ctrl.NewManager(cfg, ctrl.Options{
				Scheme: testScheme,
			})
			Expect(err).NotTo(HaveOccurred())

			reconciler = FunctionChainReconciler{
				Client: k8sClient,
				Scheme: testScheme,
			}
			err = reconciler.SetupWithManager(mgr)
			Expect(err).NotTo(HaveOccurred())
		})

		BeforeEach(func() {
			buf.Reset()
			err = k8sClient.DeleteAllOf(ctx, &ntthpcv1.FunctionChain{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &ntthpcv1.FunctionType{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			time.Sleep(100 * time.Millisecond)
		})

		AfterEach(func() {
			buf.Reset()
		})

		It("check there is no err log", func() {
			err := createFunctionChain(ctx, functionChain1)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionType(ctx, functionType1)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType2)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType3)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType4)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType5)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType6)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType7)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType8)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionInfo(ctx, functionInfo1)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo2)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo3)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo4)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo5)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo6)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo7)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo8)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functionchaintest1"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			Expect(buf.String()).NotTo(ContainSubstring("unable to fetch FunctionChain"))
			Expect(buf.String()).NotTo(ContainSubstring("unable to fetch FunctionTypeList"))
			Expect(buf.String()).NotTo(ContainSubstring("unable to fetch configmap"))
			Expect(buf.String()).NotTo(ContainSubstring("Function was found but not Ready state"))
			Expect(buf.String()).NotTo(ContainSubstring("Function was not found"))
			Expect(buf.String()).NotTo(ContainSubstring("functionchain validation error"))
		})

		It("status updating content check", func() {
			err := createFunctionChain(ctx, functionChain1)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionType(ctx, functionType1)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType2)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType3)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType4)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType5)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType6)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType7)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType8)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionInfo(ctx, functionInfo1)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo2)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo3)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo4)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo5)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo6)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo7)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo8)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functionchaintest1"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			expected := ntthpcv1.FunctionChainStatus{
				Status: "Ready",
			}

			var fc ntthpcv1.FunctionChain
			err = k8sClient.Get(ctx, client.ObjectKey{Name: "functionchaintest1", Namespace: "default"}, &fc)
			Expect(err).NotTo(HaveOccurred())

			Expect(fc.Status).To(Equal(expected))
		})

		It("status updating log check", func() {
			err := createFunctionChain(ctx, functionChain1)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionType(ctx, functionType1)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType2)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType3)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType4)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType5)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType6)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType7)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionType(ctx, functionType8)
			Expect(err).NotTo(HaveOccurred())

			err = createFunctionInfo(ctx, functionInfo1)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo2)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo3)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo4)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo5)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo6)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo7)
			Expect(err).NotTo(HaveOccurred())
			err = createFunctionInfo(ctx, functionInfo8)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default", Name: "functionchaintest1"}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

			Expect(buf.String()).To(ContainSubstring("Setting FunctionChain Status"))
			Expect(buf.String()).NotTo(ContainSubstring("failed to update FunctionChain resource status"))
		})
	})
})

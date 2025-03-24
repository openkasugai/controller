/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED
*/

package controller

import (
	"bytes"
	"context"
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"

	examplecomv1 "WBFunction/api/v1"
	controllertestcpu "WBFunction/internal/controller/test/type/CPU"
	controllertestfpga "WBFunction/internal/controller/test/type/FPGA"
	controllertestgpu "WBFunction/internal/controller/test/type/GPU"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func getMgr(mgr ctrl.Manager) (ctrl.Manager, error) {
	if mgr == nil {
		return ctrl.NewManager(cfg, ctrl.Options{
			Scheme: testScheme,
		})
	}
	return mgr, nil
}

// Create WBFunction CR
func createWBFunction(ctx context.Context, wbfcr examplecomv1.WBFunction) error {
	tmp := &examplecomv1.WBFunction{}
	*tmp = wbfcr
	tmp.TypeMeta = wbfcr.TypeMeta
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Update WBFunction CR
func updateWBFunction(ctx context.Context, wbfcr examplecomv1.WBFunction) error {
	tmp := &examplecomv1.WBFunction{}
	*tmp = wbfcr
	tmp.TypeMeta = wbfcr.TypeMeta
	err := k8sClient.Update(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Delete WBFunction CR
func deleteWBFunctionCR(ctx context.Context, wbfcr examplecomv1.WBFunction) error {
	tmp := &examplecomv1.WBFunction{}
	*tmp = wbfcr
	tmp.TypeMeta = wbfcr.TypeMeta
	err := k8sClient.Delete(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Update CPUFunction CR
func updateCPUFunctionCR(ctx context.Context, cpufunccr controllertestcpu.CPUFunction) error {
	tmp := &controllertestcpu.CPUFunction{}
	*tmp = cpufunccr
	tmp.TypeMeta = cpufunccr.TypeMeta
	err := k8sClient.Update(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Update GPUFunction CR
func updateGPUFunctionCR(ctx context.Context, gpufunccr controllertestgpu.GPUFunction) error {
	tmp := &controllertestgpu.GPUFunction{}
	*tmp = gpufunccr
	tmp.TypeMeta = gpufunccr.TypeMeta
	err := k8sClient.Update(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Update FPGAFunction CR
func updateFPGAFunctionCR(ctx context.Context, fpgafunccr controllertestfpga.FPGAFunction) error {
	tmp := &controllertestfpga.FPGAFunction{}
	*tmp = fpgafunccr
	tmp.TypeMeta = fpgafunccr.TypeMeta
	err := k8sClient.Update(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Update DeviceInfo CR
func updateDeviceInfoCR(ctx context.Context, deviceinfocr examplecomv1.DeviceInfo) error {
	tmp := &examplecomv1.DeviceInfo{}
	*tmp = deviceinfocr
	tmp.TypeMeta = deviceinfocr.TypeMeta
	err := k8sClient.Update(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

func createConfig(ctx context.Context, wbconf corev1.ConfigMap) error {
	tmp := &corev1.ConfigMap{}
	*tmp = wbconf
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create DeviceInfo CR
func createDeviceInfoCR(ctx context.Context, deviceinfocr examplecomv1.DeviceInfo) error {
	tmp := &examplecomv1.DeviceInfo{}
	*tmp = deviceinfocr
	tmp.TypeMeta = deviceinfocr.TypeMeta
	tmp.ObjectMeta = deviceinfocr.ObjectMeta
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Delete DeviceInfo CR
func deleteDeviceInfoCR(ctx context.Context, deviceinfocr examplecomv1.DeviceInfo) error {
	tmp := &examplecomv1.DeviceInfo{}
	*tmp = deviceinfocr
	tmp.TypeMeta = deviceinfocr.TypeMeta
	tmp.ObjectMeta = deviceinfocr.ObjectMeta
	err := k8sClient.Delete(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

var _ = Describe("WBFunctionController", func() {
	var mgr ctrl.Manager
	var err error
	ctx := context.Background()

	var fakerecorder = record.NewFakeRecorder(10)
	var writer = bytes.Buffer{}

	Context("Test for WBFunctionReconciler_CPU", Ordered, func() {
		var reconciler WBFunctionReconciler
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

		BeforeEach(func() {
			// loger initialized
			writer.Reset()

			err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())

			// recorder initialized
			fakerecorder = record.NewFakeRecorder(10)
			reconciler = WBFunctionReconciler{
				Client:   k8sClient,
				Scheme:   testScheme,
				Recorder: fakerecorder,
			}
			err = reconciler.SetupWithManager(mgr)
			if err != nil {
				fmt.Println("Error in SetupWithManager")
			}
			Expect(err).NotTo(HaveOccurred())

			// Create functionkindmapConfig
			err = createConfig(ctx, configdata)
			if err != nil {
				fmt.Println("There is a problem in createing functionkindmap Config", err)
				fmt.Println(err)
			}
			// Create infrastructureinfoConfig
			err = createConfig(ctx, infrainfo_configdata)
			if err != nil {
				fmt.Println("There is a problem in createing infrastructureinfo Config", err)
				fmt.Println(err)
			}
			err = LoadConfigMap(&reconciler)
			Expect(err).NotTo(HaveOccurred())

		})
		It("1-1-1 cpudec", func() {
			By("Test Start")
			// Create functionkindmapConfig
			err = createConfig(ctx, cpuconfig)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFunction Config", err)
				fmt.Println(err)
			}
			// Create WBFunctionCR
			err = createWBFunction(ctx, WBFunction_cpu_decode)
			if err != nil {
				fmt.Println("There is a problem in createing WBFuncCR_cpu ", err)
				fmt.Println(err)
			}
			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-decode-main",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Get CPUFunctionCR
			var cpufunctionCR = controllertestcpu.CPUFunction{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "CPUFunction",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night01-wbfunction-decode-main",
			},
				&cpufunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting CPUFunction CR", err)
				fmt.Println(err)
			}
			cpufunctionCR.Status = CPUFunction1
			err = updateCPUFunctionCR(ctx, cpufunctionCR)
			if err != nil {
				fmt.Println("There is a problem in updateing CPUFunction CR", err)
				fmt.Println(err)
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-decode-main",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}

			// Get DeviceInfoCR
			var deviceInfoCR = examplecomv1.DeviceInfo{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "DeviceInfo",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "deviceinfo-df-night01-wbfunction-decode-main",
			},
				&deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in getting DeviceInfo CR", err)
				fmt.Println(err)
			}
			deviceInfoCR.Status = DeviceInfoRetCPU
			err = updateDeviceInfoCR(ctx, deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in updateing DeviceInfo CR", err)
				fmt.Println(err)
			}
			got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-decode-main",
			}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get WBFunctionCR
			var wbfunctionCR examplecomv1.WBFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night01-wbfunction-decode-main",
			},
				&wbfunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting WBFunction CR", err)
				fmt.Println(err)
			}
			Expect(wbfunctionCR.Status.Status).To(Equal(examplecomv1.WBDeployStatusDeployed))
			Expect(writer.String()).To(ContainSubstring("CPUFunction does not exist."))
			Expect(writer.String()).To(ContainSubstring("CustomResource Create."))
			Expect(writer.String()).To(ContainSubstring("DeviceInfo does not exist. Namespaces/Name=default/df-night01-wbfunction-decode-main"))
			Expect(writer.String()).To(ContainSubstring("Reconcile start."))
			Expect(writer.String()).To(ContainSubstring("Status Information Change end."))
			Expect(writer.String()).To(ContainSubstring("Status Information Change start."))
			Expect(writer.String()).To(ContainSubstring("Status Running Change start."))
			Expect(writer.String()).To(ContainSubstring("Status Update."))
			Expect(writer.String()).To(ContainSubstring("Success to create RequestCR."))
			Expect(writer.String()).To(ContainSubstring("Success to delete RequestCR."))
			Expect(writer.String()).To(ContainSubstring("apiVersion :example.com/v1"))
			Expect(writer.String()).To(ContainSubstring("crData.Status.Status=Waiting"))
			Expect(writer.String()).To(ContainSubstring("crFunc.status.Status=Running"))
			Expect(writer.String()).To(ContainSubstring("eventFunctionKind=1"))
			Expect(writer.String()).To(ContainSubstring("kind :CPUFunction"))
			Expect(writer.String()).To(ContainSubstring("name :df-night01-wbfunction-decode-main"))
			Expect(writer.String()).To(ContainSubstring("namespace :default"))
			Expect(writer.String()).To(ContainSubstring("Status Running Change end."))
		})

		It("1-1-2 gpu high-infer", func() {
			By("Test Start")
			// Create functionkindmapConfig
			err = createConfig(ctx, gpuconfig)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFunction Config", err)
				fmt.Println(err)
			}
			// Create WBFunctionCR
			err = createWBFunction(ctx, WBFunction_gpu)
			if err != nil {
				fmt.Println("There is a problem in createing WBFuncCR_gpu ", err)
				fmt.Println(err)
			}
			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Get GPUFunctionCR
			var gpufunctionCR = controllertestgpu.GPUFunction{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "GPUFunction",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night01-wbfunction-high-infer-main",
			},
				&gpufunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting GPUFunction CR", err)
				fmt.Println(err)
			}
			gpufunctionCR.Status = GPUFunction1
			err = updateGPUFunctionCR(ctx, gpufunctionCR)
			if err != nil {
				fmt.Println("There is a problem in updateing GPUFunction CR", err)
				fmt.Println(err)
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}

			// Get DeviceInfoCR
			var deviceInfoCR = examplecomv1.DeviceInfo{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "DeviceInfo",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "deviceinfo-df-night01-wbfunction-high-infer-main",
			},
				&deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in getting DeviceInfo CR", err)
				fmt.Println(err)
			}
			deviceInfoCR.Status = DeviceInfoRetGPU
			err = updateDeviceInfoCR(ctx, deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in updateing DeviceInfo CR", err)
				fmt.Println(err)
			}
			got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-high-infer-main",
			}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get WBFunctionCR
			var wbfunctionCR examplecomv1.WBFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night01-wbfunction-high-infer-main",
			},
				&wbfunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting WBFunction CR", err)
				fmt.Println(err)
			}

			fmt.Println("test GPU : " + writer.String())
			Expect(wbfunctionCR.Status.Status).To(Equal(examplecomv1.WBDeployStatusDeployed))
			Expect(writer.String()).To(ContainSubstring("GPUFunction does not exist."))
			Expect(writer.String()).To(ContainSubstring("CustomResource Create."))
			Expect(writer.String()).To(ContainSubstring("DeviceInfo does not exist. Namespaces/Name=default/df-night01-wbfunction-high-infer-main"))
			Expect(writer.String()).To(ContainSubstring("Reconcile start."))
			Expect(writer.String()).To(ContainSubstring("Status Information Change end."))
			Expect(writer.String()).To(ContainSubstring("Status Information Change start."))
			Expect(writer.String()).To(ContainSubstring("Status Running Change start."))
			Expect(writer.String()).To(ContainSubstring("Status Update."))
			Expect(writer.String()).To(ContainSubstring("Success to create RequestCR."))
			Expect(writer.String()).To(ContainSubstring("Success to delete RequestCR."))
			Expect(writer.String()).To(ContainSubstring("apiVersion :example.com/v1"))
			Expect(writer.String()).To(ContainSubstring("crData.Status.Status=Waiting"))
			Expect(writer.String()).To(ContainSubstring("crFunc.status.Status=Running"))
			Expect(writer.String()).To(ContainSubstring("eventFunctionKind=1"))
			Expect(writer.String()).To(ContainSubstring("kind :GPUFunction"))
			Expect(writer.String()).To(ContainSubstring("name :df-night01-wbfunction-high-infer-main"))
			Expect(writer.String()).To(ContainSubstring("namespace :default"))
			Expect(writer.String()).To(ContainSubstring("Status Running Change end."))
		})
		It("1-1-3 fpga", func() {
			By("Test Start")
			// Create functionkindmapConfig
			err = createConfig(ctx, fpgafuncconfig_fr_high_infer)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFunction Config", err)
				fmt.Println(err)
			}
			// Create WBFunctionCR
			err = createWBFunction(ctx, WBFunction_fpga_fr)
			if err != nil {
				fmt.Println("There is a problem in createing WBFuncCR_fpga ", err)
				fmt.Println(err)
			}

			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())
			// Get FPGAFunctionCR
			var fpgafunctionCR = controllertestfpga.FPGAFunction{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "FPGAFunction",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night01-wbfunction-filter-resize-high-infer",
			},
				&fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting FPGAFunction CR", err)
				fmt.Println(err)
			}

			var rtx *controllertestfpga.RxTxData
			var nil1 *int32

			Expect(fpgafunctionCR.Spec.FunctionChannelID).To(Equal(nil1))
			Expect(fpgafunctionCR.Spec.FunctionKernelID).To(Equal(nil1))
			Expect(fpgafunctionCR.Spec.PtuKernelID).To(Equal(nil1))
			Expect(fpgafunctionCR.Spec.FrameworkKernelID).To(Equal(nil1))
			Expect(fpgafunctionCR.Spec.Rx).To(Equal(rtx))
			Expect(fpgafunctionCR.Spec.Tx).To(Equal(rtx))
			fpgafunctionCR.Status = FPGAFunction1

			err = updateFPGAFunctionCR(ctx, fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in updateing FPGAFunction CR", err)
				fmt.Println(err)
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}

			// Get DeviceInfoCR
			var deviceInfoCR = examplecomv1.DeviceInfo{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "DeviceInfo",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "deviceinfo-df-night01-wbfunction-filter-resize-high-infer",
			},
				&deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in getting DeviceInfo CR", err)
				fmt.Println(err)
			}
			deviceInfoCR.Status = DeviceInfoRetFPGA
			err = updateDeviceInfoCR(ctx, deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in updateing DeviceInfo CR", err)
				fmt.Println(err)
			}
			got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-filter-resize-high-infer",
			}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get WBFunctionCR
			var wbfunctionCR examplecomv1.WBFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night01-wbfunction-filter-resize-high-infer",
			},
				&wbfunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting WBFunction CR", err)
				fmt.Println(err)
			}

			Expect(wbfunctionCR.Status.Status).To(Equal(examplecomv1.WBDeployStatusDeployed))
			Expect(writer.String()).To(ContainSubstring("FPGAFunction does not exist."))
			Expect(writer.String()).To(ContainSubstring("CustomResource Create."))
			Expect(writer.String()).To(ContainSubstring("DeviceInfo does not exist. Namespaces/Name=default/df-night01-wbfunction-filter-resize-high-infer"))
			Expect(writer.String()).To(ContainSubstring("Reconcile start."))
			Expect(writer.String()).To(ContainSubstring("Status Information Change end."))
			Expect(writer.String()).To(ContainSubstring("Status Information Change start."))
			Expect(writer.String()).To(ContainSubstring("Status Running Change start."))
			Expect(writer.String()).To(ContainSubstring("Status Update."))
			Expect(writer.String()).To(ContainSubstring("Success to create RequestCR."))
			Expect(writer.String()).To(ContainSubstring("Success to delete RequestCR."))
			Expect(writer.String()).To(ContainSubstring("apiVersion :example.com/v1"))
			Expect(writer.String()).To(ContainSubstring("crData.Status.Status=Waiting"))
			Expect(writer.String()).To(ContainSubstring("crFunc.status.Status=Running"))
			Expect(writer.String()).To(ContainSubstring("eventFunctionKind=1"))
			Expect(writer.String()).To(ContainSubstring("kind :FPGAFunction"))
			Expect(writer.String()).To(ContainSubstring("name :df-night01-wbfunction-filter-resize-high-infer"))
			Expect(writer.String()).To(ContainSubstring("namespace :default"))
			Expect(writer.String()).To(ContainSubstring("Status Running Change end."))
		})

		It("1-1-4 fpga_write_childbs", func() {
			By("Test Start")
			// Create functionkindmapConfig
			err = createConfig(ctx, fpgafuncconfig_fr_high_infer)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFunction Config", err)
			}
			// Create WBFunctionCR
			err = createWBFunction(ctx, WBFunction_fpga_fr2)
			if err != nil {
				fmt.Println("There is a problem in createing WBFuncCR_fpga ", err)
			}
			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night02-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())
			// Get FPGAFunctionCR
			var fpgafunctionCR = controllertestfpga.FPGAFunction{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "FPGAFunction",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night02-wbfunction-filter-resize-high-infer",
			},
				&fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting FPGAFunction CR", err)
			}
			var rtx *controllertestfpga.RxTxData
			var nil1 *int32

			Expect(fpgafunctionCR.Spec.FunctionChannelID).To(Equal(nil1))
			Expect(fpgafunctionCR.Spec.FunctionKernelID).To(Equal(nil1))
			Expect(fpgafunctionCR.Spec.PtuKernelID).To(Equal(nil1))
			Expect(fpgafunctionCR.Spec.FrameworkKernelID).To(Equal(nil1))
			Expect(fpgafunctionCR.Spec.Rx).To(Equal(rtx))
			Expect(fpgafunctionCR.Spec.Tx).To(Equal(rtx))

			fpgafunctionCR.Status = FPGAFunction1
			err = updateFPGAFunctionCR(ctx, fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in updateing FPGAFunction CR", err)
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night02-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}

			// Get DeviceInfoCR
			var deviceInfoCR = examplecomv1.DeviceInfo{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "DeviceInfo",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "deviceinfo-df-night02-wbfunction-filter-resize-high-infer",
			},
				&deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in getting DeviceInfo CR", err)
			}
			deviceInfoCR.Status = DeviceInfoRetFPGA
			err = updateDeviceInfoCR(ctx, deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in updateing DeviceInfo CR", err)
			}
			got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night02-wbfunction-filter-resize-high-infer",
			}})

			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get WBFunctionCR
			var wbfunctionCR examplecomv1.WBFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night02-wbfunction-filter-resize-high-infer",
			},
				&wbfunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting WBFunction CR", err)
			}
			Expect(wbfunctionCR.Status.Status).To(Equal(examplecomv1.WBDeployStatusDeployed))

		})

		It("1-1-5 delete_df", func() {
			By("Test Start")
			// Create functionkindmapConfig
			err = createConfig(ctx, fpgafuncconfig_fr_high_infer)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFunction Config", err)
			}
			// Create WBFunctionCR
			err = createWBFunction(ctx, WBFunction_fpga_fr3)
			if err != nil {
				fmt.Println("There is a problem in createing WBFuncCR_fpga ", err)
			}
			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night03-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())
			// Get FPGAFunctionCR
			var fpgafunctionCR = controllertestfpga.FPGAFunction{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "FPGAFunction",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night03-wbfunction-filter-resize-high-infer",
			},
				&fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting FPGAFunction CR", err)
			}
			fpgafunctionCR.Status = FPGAFunction1
			err = updateFPGAFunctionCR(ctx, fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in updateing FPGAFunction CR", err)
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night03-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}

			// Get DeviceInfoCR
			var deviceInfoCR = examplecomv1.DeviceInfo{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "DeviceInfo",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "deviceinfo-df-night03-wbfunction-filter-resize-high-infer",
			},
				&deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in getting DeviceInfo CR", err)
			}
			deviceInfoCR.Status = DeviceInfoRetFPGA
			err = updateDeviceInfoCR(ctx, deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in updateing DeviceInfo CR", err)
			}
			got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night03-wbfunction-filter-resize-high-infer",
			}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get WBFunctionCR
			var wbfunctionCR examplecomv1.WBFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night03-wbfunction-filter-resize-high-infer",
			},
				&wbfunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting WBFunction CR", err)
			}
			if wbfunctionCR.Status.Status == examplecomv1.WBDeployStatusDeployed {
				fmt.Println("DELETE event start.")

				err = deleteWBFunctionCR(ctx, wbfunctionCR)
				if err != nil {
					fmt.Println("There is a problem in delete WBFunction CR", err)
				}
				_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "df-night03-wbfunction-filter-resize-high-infer",
				}})
				Expect(err).NotTo(HaveOccurred())

				var wbfunctionCR2 examplecomv1.WBFunction
				err = k8sClient.Get(ctx, client.ObjectKey{
					Namespace: "default",
					Name:      "df-night03-wbfunction-filter-resize-high-infer",
				},
					&wbfunctionCR2)
				Expect(err).NotTo(HaveOccurred())
				Expect(writer.String()).To(ContainSubstring("Update WBFunction to Terminating."))
				Expect(writer.String()).To(ContainSubstring("Delete FunctionCR was end."))
				Expect(writer.String()).To(ContainSubstring("WBFunction change to Terminating was end."))

				if wbfunctionCR2.Status.Status == examplecomv1.WBDeployStatusTerminating {
					_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "df-night03-wbfunction-filter-resize-high-infer",
					}})
					Expect(err).NotTo(HaveOccurred())
					Expect(writer.String()).To(ContainSubstring("Create DeviceInfoCR."))
					var deviceInfoCR2 = examplecomv1.DeviceInfo{
						TypeMeta: metav1.TypeMeta{
							APIVersion: "example.com/v1",
							Kind:       "DeviceInfo",
						},
					}
					err = k8sClient.Get(ctx, client.ObjectKey{
						Namespace: "default",
						Name:      "deviceinfo-df-night03-wbfunction-filter-resize-high-infer",
					},
						&deviceInfoCR2)
					if err != nil {
						fmt.Println("Failed to get DeviceInfoCR:", err)
					}
					if deviceInfoCR2.Spec.Request.RequestType == examplecomv1.RequestUndeploy {

						deviceInfoCR2.Status.Response.Status = examplecomv1.ResponceUndeployed
						err = updateDeviceInfoCR(ctx, deviceInfoCR2)
						if err != nil {
							fmt.Println("Failed Update to DeviceInfoCR:", err)
						}

						_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
							Namespace: "default",
							Name:      "df-night03-wbfunction-filter-resize-high-infer",
						}})
						Expect(err).NotTo(HaveOccurred())
						Expect(writer.String()).To(ContainSubstring("Delete DeviceInfo start."))
						Expect(writer.String()).To(ContainSubstring("Delete DeviceInfo was end."))
						Expect(writer.String()).To(ContainSubstring("Update WBFunctions to Released."))
						Expect(writer.String()).To(ContainSubstring("WBFunction change to Released was end."))

						var wbfunctionCR3 examplecomv1.WBFunction
						err = k8sClient.Get(ctx, client.ObjectKey{
							Namespace: "default",
							Name:      "df-night03-wbfunction-filter-resize-high-infer",
						},
							&wbfunctionCR3)
						Expect(err).NotTo(HaveOccurred())

						if wbfunctionCR3.Status.Status == examplecomv1.WBDeployStatusReleased {
							var deviceInfoCR3 = examplecomv1.DeviceInfo{
								TypeMeta: metav1.TypeMeta{
									APIVersion: "example.com/v1",
									Kind:       "DeviceInfo",
								},
							}
							err = k8sClient.Get(ctx, client.ObjectKey{
								Namespace: "default",
								Name:      "deviceinfo-df-night03-wbfunction-filter-resize-high-infer",
							},
								&deviceInfoCR3)
							if errors.IsNotFound(err) {
								_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
									Namespace: "default",
									Name:      "df-night03-wbfunction-filter-resize-high-infer",
								}})
								Expect(err).NotTo(HaveOccurred())
								Expect(writer.String()).To(ContainSubstring("Delete WBFunction."))
								Expect(writer.String()).To(ContainSubstring("Delete WBFunctionCR was end."))

								var wbfunctionCR4 examplecomv1.WBFunction
								err = k8sClient.Get(ctx, client.ObjectKey{
									Namespace: "default",
									Name:      "df-night03-wbfunction-filter-resize-high-infer",
								},
									&wbfunctionCR4)
								if err != nil {
									if errors.IsNotFound(err) {
										fmt.Println("Sccess to delete WBFunction.")
									} else {
										fmt.Println("Error:", err)
									}
								} else {
									fmt.Println("Failed to detele WBFunction.")
								}
							} else {
								fmt.Println("Failed to delete DeviceInfo.")
							}
						} else {
							fmt.Println("WBFunctionStatus is not Released:", wbfunctionCR3.Status.Status)
						}
					} else {
						fmt.Println("DeviceInfoCR RequestType is not UnDeploy:", deviceInfoCR2.Spec.Request.RequestType)
					}
				} else {
					fmt.Println("WBFunctionStatus is not Terminating:", wbfunctionCR2.Status.Status)
				}
			} else {
				fmt.Println("WBFunctionStatus is not Deployed:", wbfunctionCR.Status.Status)
			}
		})

		It("1-1-6 delete_df_pat⑤", func() {
			By("Test Start")
			// Create functionkindmapConfig
			err = createConfig(ctx, fpgafuncconfig_fr_high_infer)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFunction Config", err)
			}
			// Create WBFunctionCR
			err = createWBFunction(ctx, WBFunction_fpga_fr4)
			if err != nil {
				fmt.Println("There is a problem in createing WBFuncCR_fpga ", err)
			}
			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night04-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Get FPGAFunctionCR
			var fpgafunctionCR = controllertestfpga.FPGAFunction{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "FPGAFunction",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night04-wbfunction-filter-resize-high-infer",
			},
				&fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting FPGAFunction CR", err)
			}
			fpgafunctionCR.Status = FPGAFunction1
			err = updateFPGAFunctionCR(ctx, fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in updateing FPGAFunction CR", err)
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night04-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}

			// Get DeviceInfoCR
			var deviceInfoCR = examplecomv1.DeviceInfo{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "DeviceInfo",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "deviceinfo-df-night04-wbfunction-filter-resize-high-infer",
			},
				&deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in getting DeviceInfo CR", err)
			}
			deviceInfoCR.Status = DeviceInfoRetFPGA
			err = updateDeviceInfoCR(ctx, deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in updateing DeviceInfo CR", err)
			}
			got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night04-wbfunction-filter-resize-high-infer",
			}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get WBFunctionCR
			var wbfunctionCR examplecomv1.WBFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night04-wbfunction-filter-resize-high-infer",
			},
				&wbfunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting WBFunction CR", err)
			}

			if wbfunctionCR.Status.Status == examplecomv1.WBDeployStatusDeployed {
				fmt.Println("DELETE event start.")

				err = deleteWBFunctionCR(ctx, wbfunctionCR)
				if err != nil {
					fmt.Println("There is a problem in delete WBFunction CR", err)
				}
				Expect(err).NotTo(HaveOccurred())

				err = createDeviceInfoCR(ctx, DeviceInfo1)
				Expect(err).NotTo(HaveOccurred())

				got, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "df-night04-wbfunction-filter-resize-high-infer",
				}})
				Expect(err).NotTo(HaveOccurred())

				var wbfunctionCR2 examplecomv1.WBFunction
				err = k8sClient.Get(ctx, client.ObjectKey{
					Namespace: "default",
					Name:      "df-night04-wbfunction-filter-resize-high-infer",
				},
					&wbfunctionCR2)
				Expect(err).NotTo(HaveOccurred())
				Expect(writer.String()).To(ContainSubstring("Update WBFunction to Terminating."))
				Expect(writer.String()).To(ContainSubstring("Delete FunctionCR was end."))
				Expect(writer.String()).To(ContainSubstring("WBFunction change to Terminating was end."))
				Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))

				got, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "df-night04-wbfunction-filter-resize-high-infer",
				}})
				Expect(err).NotTo(HaveOccurred())
				Expect(got).To(Equal(ctrl.Result{Requeue: true, RequeueAfter: 0}))
			} else {
				fmt.Println("WBFunctionStatus is not Deployed")
			}
		})

		It("1-1-7 delete_df_pat⑥", func() {
			By("Test Start")
			// Create functionkindmapConfig
			err = createConfig(ctx, fpgafuncconfig_fr_high_infer)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFunction Config", err)
			}
			// Create WBFunctionCR
			err = createWBFunction(ctx, WBFunction_fpga_fr5)
			if err != nil {
				fmt.Println("There is a problem in createing WBFuncCR_fpga ", err)
			}
			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night05-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())
			// Get FPGAFunctionCR
			var fpgafunctionCR = controllertestfpga.FPGAFunction{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "FPGAFunction",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night05-wbfunction-filter-resize-high-infer",
			},
				&fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting FPGAFunction CR", err)
			}
			fpgafunctionCR.Status = FPGAFunction1
			err = updateFPGAFunctionCR(ctx, fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in updateing FPGAFunction CR", err)
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night05-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}

			// Get DeviceInfoCR
			var deviceInfoCR = examplecomv1.DeviceInfo{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "DeviceInfo",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "deviceinfo-df-night05-wbfunction-filter-resize-high-infer",
			},
				&deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in getting DeviceInfo CR", err)
			}
			deviceInfoCR.Status = DeviceInfoRetFPGA
			err = updateDeviceInfoCR(ctx, deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in updateing DeviceInfo CR", err)
			}
			got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night05-wbfunction-filter-resize-high-infer",
			}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get WBFunctionCR
			var wbfunctionCR examplecomv1.WBFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night05-wbfunction-filter-resize-high-infer",
			},
				&wbfunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting WBFunction CR", err)
			}
			if wbfunctionCR.Status.Status == examplecomv1.WBDeployStatusDeployed {
				fmt.Println("DELETE event start.")

				err = deleteWBFunctionCR(ctx, wbfunctionCR)
				if err != nil {
					fmt.Println("There is a problem in delete WBFunction CR", err)
				}

				got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "df-night05-wbfunction-filter-resize-high-infer",
				}})
				Expect(err).NotTo(HaveOccurred())
				Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))

				var wbfunctionCR2 examplecomv1.WBFunction
				err = k8sClient.Get(ctx, client.ObjectKey{
					Namespace: "default",
					Name:      "df-night05-wbfunction-filter-resize-high-infer",
				},
					&wbfunctionCR2)
				Expect(err).NotTo(HaveOccurred())
				Expect(writer.String()).To(ContainSubstring("Update WBFunction to Terminating."))
				Expect(writer.String()).To(ContainSubstring("Delete FunctionCR was end."))
				Expect(writer.String()).To(ContainSubstring("WBFunction change to Terminating was end."))

				if wbfunctionCR2.Status.Status == examplecomv1.WBDeployStatusTerminating {

					err = createDeviceInfoCR(ctx, DeviceInfo2)
					Expect(err).NotTo(HaveOccurred())

					got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "df-night05-wbfunction-filter-resize-high-infer",
					}})
					Expect(err).NotTo(HaveOccurred())
					Expect(got).To(Equal(ctrl.Result{Requeue: true, RequeueAfter: 0}))
					Expect(writer.String()).To(ContainSubstring("Delete FunctionCR was end."))
				} else {
					fmt.Println("WBFunctionStatus is not Terminating")
				}
			} else {
				fmt.Println("WBFunctionStatus is not Deployed")
			}
		})

		It("1-1-8 delete_df_pat⑦", func() {
			By("Test Start")
			// Create functionkindmapConfig
			err = createConfig(ctx, fpgafuncconfig_fr_high_infer)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFunction Config", err)
			}
			// Create WBFunctionCR
			err = createWBFunction(ctx, WBFunction_fpga_fr6)
			if err != nil {
				fmt.Println("There is a problem in createing WBFuncCR_fpga ", err)
			}
			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night06-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())
			// Get FPGAFunctionCR
			var fpgafunctionCR = controllertestfpga.FPGAFunction{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "FPGAFunction",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night06-wbfunction-filter-resize-high-infer",
			},
				&fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting FPGAFunction CR", err)
			}
			fpgafunctionCR.Status = FPGAFunction1
			err = updateFPGAFunctionCR(ctx, fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in updateing FPGAFunction CR", err)
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night06-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}

			// Get DeviceInfoCR
			var deviceInfoCR = examplecomv1.DeviceInfo{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "DeviceInfo",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "deviceinfo-df-night06-wbfunction-filter-resize-high-infer",
			},
				&deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in getting DeviceInfo CR", err)
			}
			deviceInfoCR.Status = DeviceInfoRetFPGA
			err = updateDeviceInfoCR(ctx, deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in updateing DeviceInfo CR", err)
			}
			got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night06-wbfunction-filter-resize-high-infer",
			}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get WBFunctionCR
			var wbfunctionCR examplecomv1.WBFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night06-wbfunction-filter-resize-high-infer",
			},
				&wbfunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting WBFunction CR", err)
			}
			if wbfunctionCR.Status.Status == examplecomv1.WBDeployStatusDeployed {
				fmt.Println("DELETE event start.")

				err = deleteWBFunctionCR(ctx, wbfunctionCR)
				if err != nil {
					fmt.Println("There is a problem in delete WBFunction CR", err)
				}

				_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "df-night06-wbfunction-filter-resize-high-infer",
				}})
				Expect(err).NotTo(HaveOccurred())

				var wbfunctionCR2 examplecomv1.WBFunction
				err = k8sClient.Get(ctx, client.ObjectKey{
					Namespace: "default",
					Name:      "df-night06-wbfunction-filter-resize-high-infer",
				},
					&wbfunctionCR2)
				Expect(err).NotTo(HaveOccurred())
				Expect(writer.String()).To(ContainSubstring("Update WBFunction to Terminating."))
				Expect(writer.String()).To(ContainSubstring("Delete FunctionCR was end."))
				Expect(writer.String()).To(ContainSubstring("WBFunction change to Terminating was end."))

				if wbfunctionCR2.Status.Status == examplecomv1.WBDeployStatusTerminating {

					err = createDeviceInfoCR(ctx, DeviceInfo3)
					Expect(err).NotTo(HaveOccurred())

					var deviceInfoCR2 = examplecomv1.DeviceInfo{
						TypeMeta: metav1.TypeMeta{
							APIVersion: "example.com/v1",
							Kind:       "DeviceInfo",
						},
					}

					err = k8sClient.Get(ctx, client.ObjectKey{
						Namespace: "default",
						Name:      "deviceinfo-df-night06-wbfunction-filter-resize-high-infer",
					},
						&deviceInfoCR2)
					Expect(err).NotTo(HaveOccurred())

					controllerutil.AddFinalizer(&deviceInfoCR2, "finalizer")
					err = updateDeviceInfoCR(ctx, deviceInfoCR2)
					Expect(err).NotTo(HaveOccurred())
					err = deleteDeviceInfoCR(ctx, deviceInfoCR2)
					Expect(err).NotTo(HaveOccurred())

					got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "df-night06-wbfunction-filter-resize-high-infer",
					}})
					Expect(err).NotTo(HaveOccurred())
					Expect(got).To(Equal(ctrl.Result{Requeue: true, RequeueAfter: 0}))

				} else {
					fmt.Println("WBFunctionStatus is not Terminating")
				}
			} else {
				fmt.Println("WBFunctionStatus is not Deployed")
			}
		})

		It("1-1-9 delete_df_pat⑨", func() {
			By("Test Start")
			// Create functionkindmapConfig
			err = createConfig(ctx, fpgafuncconfig_fr_high_infer)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFunction Config", err)
			}
			// Create WBFunctionCR
			err = createWBFunction(ctx, WBFunction_fpga_fr7)
			if err != nil {
				fmt.Println("There is a problem in createing WBFuncCR_fpga ", err)
			}
			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night07-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())
			// Get FPGAFunctionCR
			var fpgafunctionCR = controllertestfpga.FPGAFunction{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "FPGAFunction",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night07-wbfunction-filter-resize-high-infer",
			},
				&fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting FPGAFunction CR", err)
			}
			fpgafunctionCR.Status = FPGAFunction1
			err = updateFPGAFunctionCR(ctx, fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in updateing FPGAFunction CR", err)
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night07-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}

			// Get DeviceInfoCR
			var deviceInfoCR = examplecomv1.DeviceInfo{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "DeviceInfo",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "deviceinfo-df-night07-wbfunction-filter-resize-high-infer",
			},
				&deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in getting DeviceInfo CR", err)
			}
			deviceInfoCR.Status = DeviceInfoRetFPGA
			err = updateDeviceInfoCR(ctx, deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in updateing DeviceInfo CR", err)
			}
			got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night07-wbfunction-filter-resize-high-infer",
			}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get WBFunctionCR
			var wbfunctionCR examplecomv1.WBFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night07-wbfunction-filter-resize-high-infer",
			},
				&wbfunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting WBFunction CR", err)
			}
			if wbfunctionCR.Status.Status == examplecomv1.WBDeployStatusDeployed {
				fmt.Println("DELETE event start.")

				err = deleteWBFunctionCR(ctx, wbfunctionCR)
				if err != nil {
					fmt.Println("There is a problem in delete WBFunction CR", err)
				}

				_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "df-night07-wbfunction-filter-resize-high-infer",
				}})
				Expect(err).NotTo(HaveOccurred())

				var wbfunctionCR2 examplecomv1.WBFunction
				err = k8sClient.Get(ctx, client.ObjectKey{
					Namespace: "default",
					Name:      "df-night07-wbfunction-filter-resize-high-infer",
				},
					&wbfunctionCR2)
				Expect(err).NotTo(HaveOccurred())
				Expect(writer.String()).To(ContainSubstring("Update WBFunction to Terminating."))
				Expect(writer.String()).To(ContainSubstring("Delete FunctionCR was end."))
				Expect(writer.String()).To(ContainSubstring("WBFunction change to Terminating was end."))

				if wbfunctionCR2.Status.Status == examplecomv1.WBDeployStatusTerminating {

					err = createDeviceInfoCR(ctx, DeviceInfo4)
					Expect(err).NotTo(HaveOccurred())

					got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "df-night07-wbfunction-filter-resize-high-infer",
					}})
					Expect(err).NotTo(HaveOccurred())
					Expect(got).To(Equal(ctrl.Result{Requeue: true, RequeueAfter: 0}))

				} else {
					fmt.Println("WBFunctionStatus is not Terminating")
				}
			} else {
				fmt.Println("WBFunctionStatus is not Deployed")
			}
		})

		It("1-1-10 delete_df_pat⑩-1", func() {
			By("Test Start")
			// Create functionkindmapConfig
			err = createConfig(ctx, fpgafuncconfig_fr_high_infer)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFunction Config", err)
			}
			// Create WBFunctionCR
			err = createWBFunction(ctx, WBFunction_fpga_fr8)
			if err != nil {
				fmt.Println("There is a problem in createing WBFuncCR_fpga ", err)
			}
			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night08-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())
			// Get FPGAFunctionCR
			var fpgafunctionCR = controllertestfpga.FPGAFunction{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "FPGAFunction",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night08-wbfunction-filter-resize-high-infer",
			},
				&fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting FPGAFunction CR", err)
			}
			fpgafunctionCR.Status = FPGAFunction1
			err = updateFPGAFunctionCR(ctx, fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in updateing FPGAFunction CR", err)
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night08-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}

			// Get DeviceInfoCR
			var deviceInfoCR = examplecomv1.DeviceInfo{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "DeviceInfo",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "deviceinfo-df-night08-wbfunction-filter-resize-high-infer",
			},
				&deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in getting DeviceInfo CR", err)
			}
			deviceInfoCR.Status = DeviceInfoRetFPGA
			err = updateDeviceInfoCR(ctx, deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in updateing DeviceInfo CR", err)
			}
			got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night08-wbfunction-filter-resize-high-infer",
			}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get WBFunctionCR
			var wbfunctionCR examplecomv1.WBFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night08-wbfunction-filter-resize-high-infer",
			},
				&wbfunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting WBFunction CR", err)
			}
			if wbfunctionCR.Status.Status == examplecomv1.WBDeployStatusDeployed {
				fmt.Println("DELETE event start.")

				err = deleteWBFunctionCR(ctx, wbfunctionCR)
				if err != nil {
					fmt.Println("There is a problem in delete WBFunction CR", err)
				}

				got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "df-night08-wbfunction-filter-resize-high-infer",
				}})
				Expect(err).NotTo(HaveOccurred())

				var wbfunctionCR2 examplecomv1.WBFunction
				err = k8sClient.Get(ctx, client.ObjectKey{
					Namespace: "default",
					Name:      "df-night08-wbfunction-filter-resize-high-infer",
				},
					&wbfunctionCR2)
				Expect(err).NotTo(HaveOccurred())
				Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))
				Expect(writer.String()).To(ContainSubstring("Update WBFunction to Terminating."))
				Expect(writer.String()).To(ContainSubstring("Delete FunctionCR was end."))
				Expect(writer.String()).To(ContainSubstring("WBFunction change to Terminating was end."))

				if wbfunctionCR2.Status.Status == examplecomv1.WBDeployStatusTerminating {

					err = createDeviceInfoCR(ctx, DeviceInfo5)
					Expect(err).NotTo(HaveOccurred())

					got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "df-night08-wbfunction-filter-resize-high-infer",
					}})
					Expect(err).NotTo(HaveOccurred())
					Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))

					var wbfunctionCR3 examplecomv1.WBFunction
					err = k8sClient.Get(ctx, client.ObjectKey{
						Namespace: "default",
						Name:      "df-night08-wbfunction-filter-resize-high-infer",
					},
						&wbfunctionCR3)
					Expect(err).NotTo(HaveOccurred())
					Expect(writer.String()).To(ContainSubstring("Delete DeviceInfo was end."))
					Expect(writer.String()).To(ContainSubstring("Update WBFunctions to Released."))
					Expect(wbfunctionCR3.Status.Status).To(Equal(examplecomv1.WBDeployStatusReleased))
				} else {
					fmt.Println("WBFunctionStatus is not Terminating")
				}
			} else {
				fmt.Println("WBFunctionStatus is not Deployed")
			}
		})
		It("1-1-11 delete_df_pat⑩-2", func() {
			By("Test Start")
			// Create functionkindmapConfig
			err = createConfig(ctx, fpgafuncconfig_fr_high_infer)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFunction Config", err)
			}
			// Create WBFunctionCR
			err = createWBFunction(ctx, WBFunction_fpga_fr9)
			if err != nil {
				fmt.Println("There is a problem in createing WBFuncCR_fpga ", err)
			}
			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night09-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())
			// Get FPGAFunctionCR
			var fpgafunctionCR = controllertestfpga.FPGAFunction{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "FPGAFunction",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night09-wbfunction-filter-resize-high-infer",
			},
				&fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting FPGAFunction CR", err)
			}
			fpgafunctionCR.Status = FPGAFunction1
			err = updateFPGAFunctionCR(ctx, fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in updateing FPGAFunction CR", err)
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night09-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}

			// Get DeviceInfoCR
			var deviceInfoCR = examplecomv1.DeviceInfo{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "DeviceInfo",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "deviceinfo-df-night09-wbfunction-filter-resize-high-infer",
			},
				&deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in getting DeviceInfo CR", err)
			}
			deviceInfoCR.Status = DeviceInfoRetFPGA
			err = updateDeviceInfoCR(ctx, deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in updateing DeviceInfo CR", err)
			}
			got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night09-wbfunction-filter-resize-high-infer",
			}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get WBFunctionCR
			var wbfunctionCR examplecomv1.WBFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night09-wbfunction-filter-resize-high-infer",
			},
				&wbfunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting WBFunction CR", err)
			}
			if wbfunctionCR.Status.Status == examplecomv1.WBDeployStatusDeployed {
				fmt.Println("DELETE event start.")

				err = deleteWBFunctionCR(ctx, wbfunctionCR)
				if err != nil {
					fmt.Println("There is a problem in delete WBFunction CR", err)
				}

				got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "df-night09-wbfunction-filter-resize-high-infer",
				}})
				Expect(err).NotTo(HaveOccurred())
				Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))

				var wbfunctionCR2 examplecomv1.WBFunction
				err = k8sClient.Get(ctx, client.ObjectKey{
					Namespace: "default",
					Name:      "df-night09-wbfunction-filter-resize-high-infer",
				},
					&wbfunctionCR2)
				Expect(err).NotTo(HaveOccurred())
				Expect(writer.String()).To(ContainSubstring("Update WBFunction to Terminating."))
				Expect(writer.String()).To(ContainSubstring("Delete FunctionCR was end."))
				Expect(writer.String()).To(ContainSubstring("WBFunction change to Terminating was end."))

				if wbfunctionCR2.Status.Status == examplecomv1.WBDeployStatusTerminating {

					err = createDeviceInfoCR(ctx, DeviceInfo6)
					Expect(err).NotTo(HaveOccurred())

					got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "df-night09-wbfunction-filter-resize-high-infer",
					}})
					Expect(err).NotTo(HaveOccurred())
					Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))

					var wbfunctionCR3 examplecomv1.WBFunction
					err = k8sClient.Get(ctx, client.ObjectKey{
						Namespace: "default",
						Name:      "df-night09-wbfunction-filter-resize-high-infer",
					},
						&wbfunctionCR3)
					Expect(err).NotTo(HaveOccurred())
					Expect(writer.String()).To(ContainSubstring("Update WBFunctions to Failed."))
					Expect(wbfunctionCR3.Status.Status).To(Equal(examplecomv1.WBDeployStatusFailed))
				} else {
					fmt.Println("WBFunctionStatus is not Terminating")
				}
			} else {
				fmt.Println("WBFunctionStatus is not Deployed")
			}
		})
		It("1-1-12 delete_df_pat⑪", func() {
			By("Test Start")
			// Create functionkindmapConfig
			err = createConfig(ctx, fpgafuncconfig_fr_high_infer)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFunction Config", err)
			}
			// Create WBFunctionCR
			err = createWBFunction(ctx, WBFunction_fpga_fr10)
			if err != nil {
				fmt.Println("There is a problem in createing WBFuncCR_fpga ", err)
			}
			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night10-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())
			// Get FPGAFunctionCR
			var fpgafunctionCR = controllertestfpga.FPGAFunction{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "FPGAFunction",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night10-wbfunction-filter-resize-high-infer",
			},
				&fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting FPGAFunction CR", err)
			}
			fpgafunctionCR.Status = FPGAFunction1
			err = updateFPGAFunctionCR(ctx, fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in updateing FPGAFunction CR", err)
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night10-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}

			// Get DeviceInfoCR
			var deviceInfoCR = examplecomv1.DeviceInfo{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "DeviceInfo",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "deviceinfo-df-night10-wbfunction-filter-resize-high-infer",
			},
				&deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in getting DeviceInfo CR", err)
			}
			deviceInfoCR.Status = DeviceInfoRetFPGA
			err = updateDeviceInfoCR(ctx, deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in updateing DeviceInfo CR", err)
			}
			got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night10-wbfunction-filter-resize-high-infer",
			}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get WBFunctionCR
			var wbfunctionCR examplecomv1.WBFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night10-wbfunction-filter-resize-high-infer",
			},
				&wbfunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting WBFunction CR", err)
			}
			if wbfunctionCR.Status.Status == examplecomv1.WBDeployStatusDeployed {
				fmt.Println("DELETE event start.")

				err = deleteWBFunctionCR(ctx, wbfunctionCR)
				if err != nil {
					fmt.Println("There is a problem in delete WBFunction CR", err)
				}

				got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "df-night10-wbfunction-filter-resize-high-infer",
				}})
				Expect(err).NotTo(HaveOccurred())

				var wbfunctionCR2 examplecomv1.WBFunction
				err = k8sClient.Get(ctx, client.ObjectKey{
					Namespace: "default",
					Name:      "df-night10-wbfunction-filter-resize-high-infer",
				},
					&wbfunctionCR2)
				Expect(err).NotTo(HaveOccurred())
				Expect(writer.String()).To(ContainSubstring("Update WBFunction to Terminating."))
				Expect(writer.String()).To(ContainSubstring("Delete FunctionCR was end."))
				Expect(writer.String()).To(ContainSubstring("WBFunction change to Terminating was end."))
				Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))

				if wbfunctionCR2.Status.Status == examplecomv1.WBDeployStatusTerminating {
					err = createDeviceInfoCR(ctx, DeviceInfo7)
					Expect(err).NotTo(HaveOccurred())

					var deviceInfoCR2 = examplecomv1.DeviceInfo{
						TypeMeta: metav1.TypeMeta{
							APIVersion: "example.com/v1",
							Kind:       "DeviceInfo",
						},
					}

					err = k8sClient.Get(ctx, client.ObjectKey{
						Namespace: "default",
						Name:      "deviceinfo-df-night10-wbfunction-filter-resize-high-infer",
					},
						&deviceInfoCR2)
					Expect(err).NotTo(HaveOccurred())

					wbfunctionCR2.Status.Status = examplecomv1.WBDeployStatusReleased
					err = updateWBFunction(ctx, wbfunctionCR2)
					Expect(err).NotTo(HaveOccurred())

					controllerutil.AddFinalizer(&deviceInfoCR2, "finalizer")
					err = updateDeviceInfoCR(ctx, deviceInfoCR2)
					Expect(err).NotTo(HaveOccurred())
					err = deleteDeviceInfoCR(ctx, deviceInfoCR2)
					Expect(err).NotTo(HaveOccurred())

					got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "df-night10-wbfunction-filter-resize-high-infer",
					}})
					Expect(err).NotTo(HaveOccurred())
					Expect(got).To(Equal(ctrl.Result{Requeue: true, RequeueAfter: 0}))
				} else {
					fmt.Println("WBFunctionStatus is not Terminating")
				}
			} else {
				fmt.Println("WBFunctionStatus is not Deployed")
			}
		})

		It("1-2-1 Unexpected route_None DeletionTimeStamp", func() {
			By("Test Start")
			// Create functionkindmapConfig
			err = createConfig(ctx, fpgafuncconfig_fr_high_infer)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFunction Config", err)
			}
			// Create WBFunctionCR
			err = createWBFunction(ctx, WBFunction_fpga_fr11)
			if err != nil {
				fmt.Println("There is a problem in createing WBFuncCR_fpga ", err)
			}
			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night11-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())
			// Get FPGAFunctionCR
			var fpgafunctionCR = controllertestfpga.FPGAFunction{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "FPGAFunction",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night11-wbfunction-filter-resize-high-infer",
			},
				&fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting FPGAFunction CR", err)
			}
			fpgafunctionCR.Status = FPGAFunction1
			err = updateFPGAFunctionCR(ctx, fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in updateing FPGAFunction CR", err)
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night11-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}

			// Get DeviceInfoCR
			var deviceInfoCR = examplecomv1.DeviceInfo{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "DeviceInfo",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "deviceinfo-df-night11-wbfunction-filter-resize-high-infer",
			},
				&deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in getting DeviceInfo CR", err)
			}
			deviceInfoCR.Status = DeviceInfoRetFPGA
			err = updateDeviceInfoCR(ctx, deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in updateing DeviceInfo CR", err)
			}
			got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night11-wbfunction-filter-resize-high-infer",
			}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get WBFunctionCR
			var wbfunctionCR examplecomv1.WBFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night11-wbfunction-filter-resize-high-infer",
			},
				&wbfunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting WBFunction CR", err)
			}
			if wbfunctionCR.Status.Status == examplecomv1.WBDeployStatusDeployed {
				fmt.Println("DELETE event start.")

				err = deleteWBFunctionCR(ctx, wbfunctionCR)
				if err != nil {
					fmt.Println("There is a problem in delete WBFunction CR", err)
				}

				got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "df-night11-wbfunction-filter-resize-high-infer",
				}})
				Expect(err).NotTo(HaveOccurred())

				var wbfunctionCR2 examplecomv1.WBFunction
				err = k8sClient.Get(ctx, client.ObjectKey{
					Namespace: "default",
					Name:      "df-night11-wbfunction-filter-resize-high-infer",
				},
					&wbfunctionCR2)
				Expect(err).NotTo(HaveOccurred())
				Expect(writer.String()).To(ContainSubstring("Update WBFunction to Terminating."))
				Expect(writer.String()).To(ContainSubstring("Delete FunctionCR was end."))
				Expect(writer.String()).To(ContainSubstring("WBFunction change to Terminating was end."))
				Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))

				if wbfunctionCR2.Status.Status == examplecomv1.WBDeployStatusTerminating {
					err = createDeviceInfoCR(ctx, DeviceInfo8)
					Expect(err).NotTo(HaveOccurred())

					var deviceInfoCR2 = examplecomv1.DeviceInfo{
						TypeMeta: metav1.TypeMeta{
							APIVersion: "example.com/v1",
							Kind:       "DeviceInfo",
						},
					}

					err = k8sClient.Get(ctx, client.ObjectKey{
						Namespace: "default",
						Name:      "deviceinfo-df-night11-wbfunction-filter-resize-high-infer",
					},
						&deviceInfoCR2)
					Expect(err).NotTo(HaveOccurred())

					wbfunctionCR2.Status.Status = examplecomv1.WBDeployStatusReleased
					err = updateWBFunction(ctx, wbfunctionCR2)
					Expect(err).NotTo(HaveOccurred())

					controllerutil.AddFinalizer(&deviceInfoCR2, "finalizer")
					err = updateDeviceInfoCR(ctx, deviceInfoCR2)
					Expect(err).NotTo(HaveOccurred())

					got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "df-night11-wbfunction-filter-resize-high-infer",
					}})
					Expect(err).NotTo(HaveOccurred())
					Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))
					Expect(writer.String()).To(ContainSubstring("Unexpected route: WBFunctionCRStatus: Released, DeviceInfoCRRequestType: Undeploy, DeletionTimeStamp: DeletionTimeStamp is nil,  DeviceInfoCRResponceStatus: Undeployed"))
				} else {
					fmt.Println("WBFunctionStatus is not Terminating")
				}
			} else {
				fmt.Println("WBFunctionStatus is not Deployed")
			}
		})

		It("1-2-2 Unexpected route_DeletionTimeStamp exists", func() {
			By("Test Start")
			// Create functionkindmapConfig
			err = createConfig(ctx, fpgafuncconfig_fr_high_infer)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFunction Config", err)
			}
			// Create WBFunctionCR
			err = createWBFunction(ctx, WBFunction_fpga_fr12)
			if err != nil {
				fmt.Println("There is a problem in createing WBFuncCR_fpga ", err)
			}
			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night12-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())
			// Get FPGAFunctionCR
			var fpgafunctionCR = controllertestfpga.FPGAFunction{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "FPGAFunction",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night12-wbfunction-filter-resize-high-infer",
			},
				&fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting FPGAFunction CR", err)
			}
			fpgafunctionCR.Status = FPGAFunction1
			err = updateFPGAFunctionCR(ctx, fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in updateing FPGAFunction CR", err)
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night12-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}

			// Get DeviceInfoCR
			var deviceInfoCR = examplecomv1.DeviceInfo{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "DeviceInfo",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "deviceinfo-df-night12-wbfunction-filter-resize-high-infer",
			},
				&deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in getting DeviceInfo CR", err)
			}
			deviceInfoCR.Status = DeviceInfoRetFPGA
			err = updateDeviceInfoCR(ctx, deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in updateing DeviceInfo CR", err)
			}
			got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night12-wbfunction-filter-resize-high-infer",
			}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get WBFunctionCR
			var wbfunctionCR examplecomv1.WBFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night12-wbfunction-filter-resize-high-infer",
			},
				&wbfunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting WBFunction CR", err)
			}
			if wbfunctionCR.Status.Status == examplecomv1.WBDeployStatusDeployed {
				fmt.Println("DELETE event start.")

				err = deleteWBFunctionCR(ctx, wbfunctionCR)
				if err != nil {
					fmt.Println("There is a problem in delete WBFunction CR", err)
				}

				got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "df-night12-wbfunction-filter-resize-high-infer",
				}})
				Expect(err).NotTo(HaveOccurred())

				var wbfunctionCR2 examplecomv1.WBFunction
				err = k8sClient.Get(ctx, client.ObjectKey{
					Namespace: "default",
					Name:      "df-night12-wbfunction-filter-resize-high-infer",
				},
					&wbfunctionCR2)
				Expect(err).NotTo(HaveOccurred())
				Expect(writer.String()).To(ContainSubstring("Update WBFunction to Terminating."))
				Expect(writer.String()).To(ContainSubstring("Delete FunctionCR was end."))
				Expect(writer.String()).To(ContainSubstring("WBFunction change to Terminating was end."))
				Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))

				if wbfunctionCR2.Status.Status == examplecomv1.WBDeployStatusTerminating {
					err = createDeviceInfoCR(ctx, DeviceInfo9)
					Expect(err).NotTo(HaveOccurred())

					var deviceInfoCR2 = examplecomv1.DeviceInfo{
						TypeMeta: metav1.TypeMeta{
							APIVersion: "example.com/v1",
							Kind:       "DeviceInfo",
						},
					}

					err = k8sClient.Get(ctx, client.ObjectKey{
						Namespace: "default",
						Name:      "deviceinfo-df-night12-wbfunction-filter-resize-high-infer",
					},
						&deviceInfoCR2)
					Expect(err).NotTo(HaveOccurred())

					wbfunctionCR2.Status.Status = examplecomv1.WBDeployStatusReleased
					err = updateWBFunction(ctx, wbfunctionCR2)
					Expect(err).NotTo(HaveOccurred())

					controllerutil.AddFinalizer(&deviceInfoCR2, "finalizer")
					deviceInfoCR2.Status.Response.Status = examplecomv1.ResponceInitial
					err = updateDeviceInfoCR(ctx, deviceInfoCR2)
					Expect(err).NotTo(HaveOccurred())

					err = deleteDeviceInfoCR(ctx, deviceInfoCR2)
					Expect(err).NotTo(HaveOccurred())

					got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "df-night12-wbfunction-filter-resize-high-infer",
					}})
					Expect(err).NotTo(HaveOccurred())
					Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))
					Expect(writer.String()).To(ContainSubstring("Unexpected route: WBFunctionCRStatus: Released, DeletionTimeStamp:"))
				} else {
					fmt.Println("WBFunctionStatus is not Terminating")
				}
			} else {
				fmt.Println("WBFunctionStatus is not Deployed")
			}
		})
		It("1-3-1 Terminating_and_deletionstamp1", func() {
			By("Test Start")
			// Create functionkindmapConfig
			err = createConfig(ctx, fpgafuncconfig_fr_high_infer)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFunction Config", err)
			}
			// Create WBFunctionCR
			err = createWBFunction(ctx, WBFunction_fpga_fr13)
			if err != nil {
				fmt.Println("There is a problem in createing WBFuncCR_fpga ", err)
			}
			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night13-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())
			// Get FPGAFunctionCR
			var fpgafunctionCR = controllertestfpga.FPGAFunction{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "FPGAFunction",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night13-wbfunction-filter-resize-high-infer",
			},
				&fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting FPGAFunction CR", err)
			}
			fpgafunctionCR.Status = FPGAFunction1
			err = updateFPGAFunctionCR(ctx, fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in updateing FPGAFunction CR", err)
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night13-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}

			// Get DeviceInfoCR
			var deviceInfoCR = examplecomv1.DeviceInfo{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "DeviceInfo",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "deviceinfo-df-night13-wbfunction-filter-resize-high-infer",
			},
				&deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in getting DeviceInfo CR", err)
			}
			deviceInfoCR.Status = DeviceInfoRetFPGA
			err = updateDeviceInfoCR(ctx, deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in updateing DeviceInfo CR", err)
			}
			got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night13-wbfunction-filter-resize-high-infer",
			}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get WBFunctionCR
			var wbfunctionCR examplecomv1.WBFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night13-wbfunction-filter-resize-high-infer",
			},
				&wbfunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting WBFunction CR", err)
			}
			if wbfunctionCR.Status.Status == examplecomv1.WBDeployStatusDeployed {
				fmt.Println("DELETE event start.")

				err = deleteWBFunctionCR(ctx, wbfunctionCR)
				if err != nil {
					fmt.Println("There is a problem in delete WBFunction CR", err)
				}

				got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "df-night13-wbfunction-filter-resize-high-infer",
				}})
				Expect(err).NotTo(HaveOccurred())

				var wbfunctionCR2 examplecomv1.WBFunction
				err = k8sClient.Get(ctx, client.ObjectKey{
					Namespace: "default",
					Name:      "df-night13-wbfunction-filter-resize-high-infer",
				},
					&wbfunctionCR2)
				Expect(err).NotTo(HaveOccurred())
				Expect(writer.String()).To(ContainSubstring("Update WBFunction to Terminating."))
				Expect(writer.String()).To(ContainSubstring("Delete FunctionCR was end."))
				Expect(writer.String()).To(ContainSubstring("WBFunction change to Terminating was end."))
				Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))

				if wbfunctionCR2.Status.Status == examplecomv1.WBDeployStatusTerminating {
					err = createDeviceInfoCR(ctx, DeviceInfo10)
					Expect(err).NotTo(HaveOccurred())

					var deviceInfoCR2 = examplecomv1.DeviceInfo{
						TypeMeta: metav1.TypeMeta{
							APIVersion: "example.com/v1",
							Kind:       "DeviceInfo",
						},
					}

					err = k8sClient.Get(ctx, client.ObjectKey{
						Namespace: "default",
						Name:      "deviceinfo-df-night13-wbfunction-filter-resize-high-infer",
					},
						&deviceInfoCR2)
					Expect(err).NotTo(HaveOccurred())

					controllerutil.AddFinalizer(&deviceInfoCR2, "finalizer")
					deviceInfoCR2.Status.Response.Status = examplecomv1.ResponceUndeployed
					err = updateDeviceInfoCR(ctx, deviceInfoCR2)
					Expect(err).NotTo(HaveOccurred())

					err = deleteDeviceInfoCR(ctx, deviceInfoCR2)
					Expect(err).NotTo(HaveOccurred())

					got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "df-night13-wbfunction-filter-resize-high-infer",
					}})
					Expect(err).NotTo(HaveOccurred())
					Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))
					Expect(writer.String()).To(ContainSubstring("WBFunction change to Released was end."))
					var wbfunctionCR3 examplecomv1.WBFunction
					k8sClient.Get(ctx, client.ObjectKey{
						Namespace: "default",
						Name:      "df-night13-wbfunction-filter-resize-high-infer",
					},
						&wbfunctionCR3)
					if wbfunctionCR3.Status.Status == examplecomv1.WBDeployStatusReleased {
					} else {
						fmt.Println("WBFunctionStatus is not Released.")
					}
				} else {
					fmt.Println("WBFunctionStatus is not Terminating")
				}
			} else {
				fmt.Println("WBFunctionStatus is not Deployed")
			}
		})
		It("1-3-2 Terminating_and_deletionstamp2", func() {
			By("Test Start")
			// Create functionkindmapConfig
			err = createConfig(ctx, fpgafuncconfig_fr_high_infer)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFunction Config", err)
			}
			// Create WBFunctionCR
			err = createWBFunction(ctx, WBFunction_fpga_fr14)
			if err != nil {
				fmt.Println("There is a problem in createing WBFuncCR_fpga ", err)
			}
			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night14-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())
			// Get FPGAFunctionCR
			var fpgafunctionCR = controllertestfpga.FPGAFunction{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "FPGAFunction",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night14-wbfunction-filter-resize-high-infer",
			},
				&fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting FPGAFunction CR", err)
			}
			fpgafunctionCR.Status = FPGAFunction1
			err = updateFPGAFunctionCR(ctx, fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in updateing FPGAFunction CR", err)
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night14-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}

			// Get DeviceInfoCR
			var deviceInfoCR = examplecomv1.DeviceInfo{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "DeviceInfo",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "deviceinfo-df-night14-wbfunction-filter-resize-high-infer",
			},
				&deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in getting DeviceInfo CR", err)
			}
			deviceInfoCR.Status = DeviceInfoRetFPGA
			err = updateDeviceInfoCR(ctx, deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in updateing DeviceInfo CR", err)
			}
			got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night14-wbfunction-filter-resize-high-infer",
			}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get WBFunctionCR
			var wbfunctionCR examplecomv1.WBFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night14-wbfunction-filter-resize-high-infer",
			},
				&wbfunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting WBFunction CR", err)
			}
			if wbfunctionCR.Status.Status == examplecomv1.WBDeployStatusDeployed {
				fmt.Println("DELETE event start.")

				err = deleteWBFunctionCR(ctx, wbfunctionCR)
				if err != nil {
					fmt.Println("There is a problem in delete WBFunction CR", err)
				}

				got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "df-night14-wbfunction-filter-resize-high-infer",
				}})
				Expect(err).NotTo(HaveOccurred())

				var wbfunctionCR2 examplecomv1.WBFunction
				err = k8sClient.Get(ctx, client.ObjectKey{
					Namespace: "default",
					Name:      "df-night14-wbfunction-filter-resize-high-infer",
				},
					&wbfunctionCR2)
				Expect(err).NotTo(HaveOccurred())
				Expect(writer.String()).To(ContainSubstring("Update WBFunction to Terminating."))
				Expect(writer.String()).To(ContainSubstring("Delete FunctionCR was end."))
				Expect(writer.String()).To(ContainSubstring("WBFunction change to Terminating was end."))
				Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))

				if wbfunctionCR2.Status.Status == examplecomv1.WBDeployStatusTerminating {
					err = createDeviceInfoCR(ctx, DeviceInfo11)
					Expect(err).NotTo(HaveOccurred())

					var deviceInfoCR2 = examplecomv1.DeviceInfo{
						TypeMeta: metav1.TypeMeta{
							APIVersion: "example.com/v1",
							Kind:       "DeviceInfo",
						},
					}

					err = k8sClient.Get(ctx, client.ObjectKey{
						Namespace: "default",
						Name:      "deviceinfo-df-night14-wbfunction-filter-resize-high-infer",
					},
						&deviceInfoCR2)
					Expect(err).NotTo(HaveOccurred())

					controllerutil.AddFinalizer(&deviceInfoCR2, "finalizer")
					deviceInfoCR2.Status.Response.Status = ""
					deviceInfoCR2.Spec.Request.RequestType = examplecomv1.RequestUndeploy
					err = updateDeviceInfoCR(ctx, deviceInfoCR2)
					Expect(err).NotTo(HaveOccurred())

					err = deleteDeviceInfoCR(ctx, deviceInfoCR2)
					Expect(err).NotTo(HaveOccurred())

					got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "df-night14-wbfunction-filter-resize-high-infer",
					}})
					Expect(err).NotTo(HaveOccurred())
					Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))
					Expect(writer.String()).To(ContainSubstring("WBFunction change to Released was end."))
					var wbfunctionCR3 examplecomv1.WBFunction
					k8sClient.Get(ctx, client.ObjectKey{
						Namespace: "default",
						Name:      "df-night14-wbfunction-filter-resize-high-infer",
					},
						&wbfunctionCR3)
					if wbfunctionCR3.Status.Status == examplecomv1.WBDeployStatusReleased {
					} else {
						fmt.Println("WBFunctionStatus is not Released.")
					}
				} else {
					fmt.Println("WBFunctionStatus is not Terminating")
				}
			} else {
				fmt.Println("WBFunctionStatus is not Deployed")
			}
		})
		It("1-3-3 Terminating_and_deletionstamp3", func() {
			By("Test Start")
			// Create functionkindmapConfig
			err = createConfig(ctx, fpgafuncconfig_fr_high_infer)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFunction Config", err)
			}
			// Create WBFunctionCR
			err = createWBFunction(ctx, WBFunction_fpga_fr15)
			if err != nil {
				fmt.Println("There is a problem in createing WBFuncCR_fpga ", err)
			}
			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night15-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())
			// Get FPGAFunctionCR
			var fpgafunctionCR = controllertestfpga.FPGAFunction{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "FPGAFunction",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night15-wbfunction-filter-resize-high-infer",
			},
				&fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting FPGAFunction CR", err)
			}
			fpgafunctionCR.Status = FPGAFunction1
			err = updateFPGAFunctionCR(ctx, fpgafunctionCR)
			if err != nil {
				fmt.Println("There is a problem in updateing FPGAFunction CR", err)
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night15-wbfunction-filter-resize-high-infer",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}

			// Get DeviceInfoCR
			var deviceInfoCR = examplecomv1.DeviceInfo{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "example.com/v1",
					Kind:       "DeviceInfo",
				},
			}

			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "deviceinfo-df-night15-wbfunction-filter-resize-high-infer",
			},
				&deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in getting DeviceInfo CR", err)
			}
			deviceInfoCR.Status = DeviceInfoRetFPGA
			err = updateDeviceInfoCR(ctx, deviceInfoCR)
			if err != nil {
				fmt.Println("There is a problem in updateing DeviceInfo CR", err)
			}
			got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night15-wbfunction-filter-resize-high-infer",
			}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Get WBFunctionCR
			var wbfunctionCR examplecomv1.WBFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "df-night15-wbfunction-filter-resize-high-infer",
			},
				&wbfunctionCR)
			if err != nil {
				fmt.Println("There is a problem in getting WBFunction CR", err)
			}
			if wbfunctionCR.Status.Status == examplecomv1.WBDeployStatusDeployed {
				fmt.Println("DELETE event start.")

				err = deleteWBFunctionCR(ctx, wbfunctionCR)
				if err != nil {
					fmt.Println("There is a problem in delete WBFunction CR", err)
				}

				got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "df-night15-wbfunction-filter-resize-high-infer",
				}})
				Expect(err).NotTo(HaveOccurred())

				var wbfunctionCR2 examplecomv1.WBFunction
				err = k8sClient.Get(ctx, client.ObjectKey{
					Namespace: "default",
					Name:      "df-night15-wbfunction-filter-resize-high-infer",
				},
					&wbfunctionCR2)
				Expect(err).NotTo(HaveOccurred())
				Expect(writer.String()).To(ContainSubstring("Update WBFunction to Terminating."))
				Expect(writer.String()).To(ContainSubstring("Delete FunctionCR was end."))
				Expect(writer.String()).To(ContainSubstring("WBFunction change to Terminating was end."))
				Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))

				if wbfunctionCR2.Status.Status == examplecomv1.WBDeployStatusTerminating {
					err = createDeviceInfoCR(ctx, DeviceInfo12)
					Expect(err).NotTo(HaveOccurred())

					var deviceInfoCR2 = examplecomv1.DeviceInfo{
						TypeMeta: metav1.TypeMeta{
							APIVersion: "example.com/v1",
							Kind:       "DeviceInfo",
						},
					}

					err = k8sClient.Get(ctx, client.ObjectKey{
						Namespace: "default",
						Name:      "deviceinfo-df-night15-wbfunction-filter-resize-high-infer",
					},
						&deviceInfoCR2)
					Expect(err).NotTo(HaveOccurred())

					controllerutil.AddFinalizer(&deviceInfoCR2, "finalizer")
					deviceInfoCR2.Status.Response.Status = examplecomv1.ResponceError
					deviceInfoCR2.Spec.Request.RequestType = examplecomv1.RequestDeploy
					err = updateDeviceInfoCR(ctx, deviceInfoCR2)
					Expect(err).NotTo(HaveOccurred())

					err = deleteDeviceInfoCR(ctx, deviceInfoCR2)
					Expect(err).NotTo(HaveOccurred())

					got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "df-night15-wbfunction-filter-resize-high-infer",
					}})
					Expect(err).NotTo(HaveOccurred())
					Expect(got).To(Equal(ctrl.Result{Requeue: false, RequeueAfter: 0}))
					Expect(writer.String()).To(ContainSubstring("WBFunction change to Released was end."))
					var wbfunctionCR3 examplecomv1.WBFunction
					k8sClient.Get(ctx, client.ObjectKey{
						Namespace: "default",
						Name:      "df-night15-wbfunction-filter-resize-high-infer",
					},
						&wbfunctionCR3)
					if wbfunctionCR3.Status.Status == examplecomv1.WBDeployStatusReleased {
					} else {
						fmt.Println("WBFunctionStatus is not Released.")
					}
				} else {
					fmt.Println("WBFunctionStatus is not Terminating")
				}
			} else {
				fmt.Println("WBFunctionStatus is not Deployed")
			}
		})
		AfterEach(func() {
			By("Test End")
			writer.Reset()
		})
	})
})

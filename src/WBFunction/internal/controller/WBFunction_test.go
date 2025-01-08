package controller

import (
	"bytes"
	"context"
	"fmt"
	// "os"

	ctrl "sigs.k8s.io/controller-runtime"
	//	"sigs.k8s.io/controller-runtime/pkg/client"

	examplecomv1 "WBFunction/api/v1"
	controllertestcpu "WBFunction/internal/controller/test/type/CPU"
	controllertestfpga "WBFunction/internal/controller/test/type/FPGA"
	controllertestgpu "WBFunction/internal/controller/test/type/GPU"

	// ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	//"k8s.io/apimachinery/pkg/runtime"
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
	//	err := k8sClient.Create(context.Background(), &wbfcr)
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
	//	err := k8sClient.Create(context.Background(), &wbfcr)
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
	//	err := k8sClient.Create(context.Background(), &wbfcr)
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
	//	err := k8sClient.Create(context.Background(), &wbfcr)
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
	//	err := k8sClient.Create(context.Background(), &wbfcr)
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
			// fmt.Println("test:" + writer.String())
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

		It("1-2-1 gpu high-infer", func() {
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
		It("1-3-1 fpga", func() {
			By("Test Start")
			// Create functionkindmapConfig
			err = createConfig(ctx, fpgafuncconfig_fr_high_infer)
			// err = createConfig(ctx, fpgaconfig)
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

			/* debug
			// Get configdata
			var configdata corev1.ConfigMap
			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "fpgafunc-config-filter-resize-high-infer",
			},
				&configdata)
			if err != nil {
				fmt.Println("There is a problem in getting WBFunction CR", err)
				fmt.Println(err)
			}
			fmt.Println("test config data : ", configdata.Data["fpgafunc-config-filter-resize-high-infer.json"])
			fmt.Println("test config data : ", configdata.Data)
			// Get configdata
			var configdata2 corev1.ConfigMap
			err = k8sClient.Get(ctx, client.ObjectKey{
				Namespace: "default",
				Name:      "infrastructureinfo",
			},
				&configdata2)
			if err != nil {
				fmt.Println("There is a problem in getting WBFunction CR", err)
				fmt.Println(err)
			}
			fmt.Println("test config data : ", configdata2.Data["infrastructureinfo.json"])
			fmt.Println("test config data : ", configdata2.Data)
			*/

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
			// fmt.Println("test:" + writer.String())
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

		AfterEach(func() {
			By("Test End")
			writer.Reset()
		})
	})
})

/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED
*/

package controller

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap/zapcore"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	examplecomv1 "GPUFunction/api/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"

	// Additional files
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func getMgr(mgr ctrl.Manager) (ctrl.Manager, error) {
	if mgr == nil {
		return ctrl.NewManager(cfg, ctrl.Options{
			Scheme: testScheme,
		})
	}
	return mgr, nil
}

// Create GPUFunction CR
func createGPUFunction(ctx context.Context, gpufcr examplecomv1.GPUFunction) error {
	tmp := &examplecomv1.GPUFunction{}
	*tmp = gpufcr
	tmp.TypeMeta = gpufcr.TypeMeta
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create PCIeConnection CR
func createPCIeConnection(ctx context.Context, pcieccr PCIeConnection) error {
	tmp := &PCIeConnection{}
	*tmp = pcieccr
	tmp.Kind = pcieccr.Kind
	tmp.APIVersion = pcieccr.APIVersion
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create EthernetConnection CR
func createEthernetConnection(ctx context.Context, ethernetccr EthernetConnection) error {
	tmp := &EthernetConnection{}
	*tmp = ethernetccr
	tmp.Kind = ethernetccr.Kind
	tmp.APIVersion = ethernetccr.APIVersion
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create FPGAFunction CR
func createFPGAFunction(ctx context.Context, fpgafcr FPGAFunction) error {
	tmp := &FPGAFunction{}
	*tmp = fpgafcr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create CPUFunction CR
func createCPUFunction(ctx context.Context, cpufcr CPUFunction) error {
	tmp := &CPUFunction{}
	*tmp = cpufcr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create ConfigMap
func createConfig(ctx context.Context, conf corev1.ConfigMap) error {
	tmp := &corev1.ConfigMap{}
	*tmp = conf
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

func createNetworkAttachmentDefinition(ctx context.Context, nad NetworkAttachmentDefinition) error {
	tmp := &NetworkAttachmentDefinition{}
	*tmp = nad
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

var _ = Describe("GPUFunctionController", func() {
	var mgr ctrl.Manager
	var err error
	ctx := context.Background()
	var fakerecorder = record.NewFakeRecorder(10)
	var writer = bytes.Buffer{}

	// This test case is for reconciler
	Context("Test for GPUFunctionReconciler", Ordered, func() {
		var reconciler GPUFunctionReconciler
		BeforeAll(func() {
			// set up log
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

			os.Setenv("K8S_NODENAME", "worker1")
			os.Setenv("K8S_POD_ANNOTATION", "")

			mgr, err = getMgr(mgr)
			Expect(err).NotTo(HaveOccurred())
			fakerecorder = record.NewFakeRecorder(10)

			reconciler = GPUFunctionReconciler{
				Client:   k8sClient,
				Scheme:   testScheme,
				Recorder: fakerecorder,
			}

			err = reconciler.SetupWithManager(mgr)
			if err != nil {
				fmt.Println("An error occur during setupwithManager: ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// To delete crdata It
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.GPUFunction{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &PCIeConnection{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &EthernetConnection{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &FPGAFunction{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &CPUFunction{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &NetworkAttachmentDefinition{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)
		})
		AfterAll(func() {
			writer.Reset()
		})

		It("Reconcile Test", func() {
			By("Test Start")

			// Create GPUFuncConfig
			err = createConfig(ctx, gpuconfigdecode)
			if err != nil {
				fmt.Println("There is a problem in createing GPUConfig ", err)
				fmt.Printf("%T\n", err)
				fmt.Println(err)
			}

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection1)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}

			// Create FPGAFunctionCR
			err = createFPGAFunction(ctx, FPGAFunction1)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}

			// Create NetworkAttachmentDefinition
			err = createNetworkAttachmentDefinition(ctx, NetworkAttachmentDefinition1)
			if err != nil {
				fmt.Println("There is a problem in createing NetworkAttachmentDefinition")
			}

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction1)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR")
			}

			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}
			Expect(err).NotTo(HaveOccurred())

			var gpuCR examplecomv1.GPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "df-night01-wbfunction-high-infer-main",
				Namespace: "default",
			}, &gpuCR)
			if err != nil {
				// Error route
				fmt.Println("Cannot get GPUFunctionCR:", gpuCR, err)
			}

			Expect(*gpuCR.Status.FunctionIndex).To(Equal(int32(0)))

			// Delete GPUFunctionCR
			err = k8sClient.Delete(ctx, &gpuCR)
			Expect(err).NotTo(HaveOccurred())

			var gpupod corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "df-night01-wbfunction-high-infer-main-mps-dgpu-gpu-702fb653-43a4-732d-6bc4-7b3487696c90-pod",
				Namespace: "default",
			},
				&gpupod)
			if err != nil {
				// Error route
				fmt.Println("Cannot get Pod:", gpupod, err)
			}
			err = reconciler.Delete(ctx, &gpupod)
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "df-night01-wbfunction-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var gpuaCR examplecomv1.GPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "df-night01-wbfunction-high-infer-main",
				Namespace: "default",
			},
				&gpuaCR)

			Expect(err).To(MatchError(ContainSubstring("not found")))

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction1)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR")
			}
			Expect(err).NotTo(HaveOccurred())

			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}
			Expect(err).NotTo(HaveOccurred())

			var gpuCR2 examplecomv1.GPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "df-night01-wbfunction-high-infer-main",
				Namespace: "default",
			}, &gpuCR2)
			if err != nil {
				// Error route
				fmt.Println("Cannot get GPUFunctionCR:", gpuCR2, err)
			}

			Expect(*gpuCR2.Status.FunctionIndex).To(Equal(int32(0)))
			_ = k8sClient.Delete(ctx, &gpuCR2)
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}
		})
		It("Reconcile Test2", func() {
			By("Test Start")

			// Create GPUFuncConfig
			err = createConfig(ctx, gpuconfigdecode)
			if err != nil {
				fmt.Println("There is a problem in createing GPUConfig ", err)
				fmt.Printf("%T\n", err)
				fmt.Println(err)
			}

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection2)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			// Create FPGAFunctionCR
			err = createFPGAFunction(ctx, FPGAFunction1)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}
			// Create NetworkAttachmentDefinition
			err = createNetworkAttachmentDefinition(ctx, NetworkAttachmentDefinition1)
			if err != nil {
				fmt.Println("There is a problem in createing NetworkAttachmentDefinition")
			}
			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction2)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR")
			}

			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night02-wbfunction-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}
			Expect(err).NotTo(HaveOccurred())
			var gpuCR3 examplecomv1.GPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "df-night02-wbfunction-high-infer-main",
				Namespace: "default",
			}, &gpuCR3)
			if err != nil {
				// Error route
				fmt.Println("Cannot get GPUFunctionCR:", gpuCR3, err)
			}
			Expect(err).NotTo(HaveOccurred())
			uuid := strings.Replace(gpuCR3.Spec.AcceleratorIDs[0].ID, "GPU-", "", -1)
			podname := "df-night02-wbfunction-high-infer-main-mps-dgpu-gpu-" + uuid + "-pod"
			var gpupod corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      podname,
				Namespace: "default",
			},
				&gpupod)
			if err != nil {
				// Error route
				fmt.Println("Cannot get Pod:", gpupod, err)
			}
			Expect(*gpuCR3.Status.FunctionIndex).To(Equal(int32(99)))
			Expect(*gpuCR3.Status.PodName).To(Equal(gpupod.Name))
		})
		//Test for GetFunc
		It("Test_1-1-1_fpga-filter-resize-dma-high-infer", func() {
			// Create GPUFuncConfig
			err = createConfig(ctx, gpuconfighigh)
			if err != nil {
				fmt.Println("There is a problem in createing configmap ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection111)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGAFunctionCR
			err = createFPGAFunction(ctx, FPGAFunction111)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create NetworkAttachmentDifinition
			err = createNetworkAttachmentDefinition(ctx, NetworkAttachmentDefinition1)
			if err != nil {
				fmt.Println("There is a problem in createing NetworkAttachmentDefinition ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction111)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "gpufunctest111-wbfunction-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var gpuCR examplecomv1.GPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "gpufunctest111-wbfunction-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&gpuCR)
			Expect(err).NotTo(HaveOccurred())

			var addTest = true
			var rxTest = "DMA"
			var txTest = "RTP"
			var functionIndex int32 = 0
			var podName = "gpufunctest111-wbfunction-high-infer-main-mps-dgpu-0-pod"
			// test for updating CR.Status
			expectedStatus := examplecomv1.GPUFunctionStatus{
				StartTime: gpuCR.Status.StartTime,
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "gpufunctest111",
					Namespace: "default",
				},
				FunctionName: "high-infer",
				ImageURI:     "localhost/gpu-deepstream-app:3.1.0",
				ConfigName:   "gpufunc-config-high-infer",
				Status:       "Running",
				SharedMemory: &examplecomv1.SharedMemorySpec{
					FilePrefix:      "gpufunctest111-wbfunction-high-infer-main",
					CommandQueueID:  "gpufunctest111-wbfunction-high-infer-main",
					SharedMemoryMiB: 1,
				},
				VirtualNetworkDeviceDriverType: "sriov",
				AdditionalNetwork:              &addTest,
				FunctionIndex:                  &functionIndex,
				RxProtocol:                     &rxTest,
				TxProtocol:                     &txTest,
				PodName:                        &podName,
			}

			// test for updating StartTime
			err = checkTime(testTime.Time, gpuCR.Status.StartTime.Time)
			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status except for Status.StartTime
			Expect(gpuCR.Status).To(Equal(expectedStatus))

			// test for creating finalizer
			expectedFinalizer := []string{
				"gpufunction.finalizers.example.com.v1",
			}
			Expect(gpuCR.Finalizers).To(Equal(expectedFinalizer))

			// check logs
			Expect(writer.String()).To(ContainSubstring("previousConnectionCRName : gpufunctest111-wbconnection-filter-resize-high-infer-main-high-infer-main"))
			Expect(writer.String()).To(ContainSubstring("rxProtocol : DMA"))
			Expect(writer.String()).To(ContainSubstring("txProtocol : RTP"))
			Expect(writer.String()).To(ContainSubstring("Success to create Pod."))
			Expect(writer.String()).To(ContainSubstring("Status Update."))

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

			var gpupod corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "gpufunctest111-wbfunction-high-infer-main-mps-dgpu-0-pod",
				Namespace: TESTNAMESPACE,
			},
				&gpupod)
			if err != nil {
				// error
				fmt.Println("Could not get Pod:", gpupod, err)
			}
			Expect(err).NotTo(HaveOccurred())

			expectedPodMeta := metav1.ObjectMeta{
				Name:      "gpufunctest111-wbfunction-high-infer-main-mps-dgpu-0-pod",
				Namespace: "default",
				Annotations: map[string]string{
					"k8s.v1.cni.cncf.io/networks": "[ {\"name\": \"worker1-config-net-sriov\",\"ips\": [\"192.174.90.141/24\"] } ]",
				},
			}

			Expect(gpupod.ObjectMeta.Name).To(Equal(expectedPodMeta.Name))
			Expect(gpupod.ObjectMeta.Namespace).To(Equal(expectedPodMeta.Namespace))
			Expect(gpupod.ObjectMeta.Annotations["k8s.v1.cni.cncf.io/networks"]).To(Equal(expectedPodMeta.Annotations["k8s.v1.cni.cncf.io/networks"]))

			expectedPodSpecContainerEnv := []corev1.EnvVar{
				{
					Name:  "CONNECTOR_ID",
					Value: "gpufunctest111-wbfunction-high-infer-main",
				},
				{
					Name:  "CUDA_MPS_LOG_DIRECTORY",
					Value: "/tmp/nvidia-mps",
				},
				{
					Name:  "CUDA_MPS_PIPE_DIRECTORY",
					Value: "/tmp/nvidia-mps",
				},
				{
					Name:  "CUDA_VISIBLE_DEVICES",
					Value: "0",
				},
				{
					Name:  "FILE_PREFIX",
					Value: "gpufunctest111-wbfunction-high-infer-main",
				},
				{
					Name:  "FPGA_DEV",
					Value: "/dev/xpcie_21330621T00D",
				},
				{
					Name:  "HEIGHT",
					Value: "1280",
				},
				{
					Name:  "K8S_POD_NAME",
					Value: "",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							APIVersion: "v1",
							FieldPath:  "metadata.name",
						},
						ResourceFieldRef: nil,
						ConfigMapKeyRef:  nil,
						SecretKeyRef:     nil,
					},
				},
				{
					Name:  "K8S_POD_NAMESPACE",
					Value: "",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							APIVersion: "v1",
							FieldPath:  "metadata.namespace",
						},
						ResourceFieldRef: nil,
						ConfigMapKeyRef:  nil,
						SecretKeyRef:     nil,
					},
				},
				{
					Name:  "RECEIVING_SERVER_IP",
					Value: "192.174.90.10",
				},
				{
					Name:  "RECEIVING_SERVER_PORT",
					Value: "2001",
				},
				{
					Name:  "SHMEM_SECONDARY",
					Value: "1",
				},
				{
					Name:  "WIDTH",
					Value: "1280",
				},
			}

			sort.Slice(expectedPodSpecContainerEnv, func(i, j int) bool {
				return expectedPodSpecContainerEnv[i].Name < expectedPodSpecContainerEnv[j].Name
			})
			var gpupodEnv []corev1.EnvVar
			for _, containers := range gpupod.Spec.Containers {
				gpupodEnv = append(gpupodEnv, containers.Env...)
			}

			sort.Slice(gpupodEnv, func(i, j int) bool {
				return gpupodEnv[i].Name < gpupodEnv[j].Name
			})

			Expect(gpupodEnv).To(Equal(expectedPodSpecContainerEnv))

			var podLimits = corev1.ResourceList{
				"hugepages-1Gi":                   resource.Quantity{},
				"intel.com/intel_sriov_netdevice": {},
			}
			podLimits["hugepages-1Gi"] = resource.Quantity(resource.MustParse("1Gi"))
			podLimits["intel.com/intel_sriov_netdevice"] = resource.Quantity(resource.MustParse("1"))

			var podRequests = corev1.ResourceList{
				"hugepages-1Gi":                   resource.Quantity{},
				"intel.com/intel_sriov_netdevice": {},
				"memory":                          {},
			}
			podRequests["hugepages-1Gi"] = resource.Quantity(resource.MustParse("1Gi"))
			podRequests["intel.com/intel_sriov_netdevice"] = resource.Quantity(resource.MustParse("1"))
			podRequests["memory"] = resource.Quantity(resource.MustParse("32Gi"))

			var privileged bool = true
			expectedPodSpecContainer := []corev1.Container{
				{
					Name:  "gpu-container0",
					Image: "localhost/gpu-deepstream-app:3.1.0",
					Command: []string{
						"sh",
						"-c",
					},
					ImagePullPolicy: corev1.PullIfNotPresent,
					Args: []string{
						"cd /opt/DeepStream-Yolo && gst-launch-1.0 -ev fpgasrc ! 'video/x-raw,format=(string)BGR,width=1280,height=1280' ! nvvideoconvert ! 'video/x-raw(memory:NVMM), format=(string)RGBA' ! m.sink_0 nvstreammux name=m nvbuf-memory-type=0 batch-size=1 width=1280 height=1280 ! queue ! nvinfer config-file-path=./config_infer_primary_yoloV4_p6_th020_040.txt batch-size=1 model-engine-file=./model_b1_gpu0_fp16.engine ! queue ! nvdsosd process-mode=1 ! nvvideoconvert ! 'video/x-raw, format=(string)BGR' ! videoconvert ! queue ! perf ! rtpvrawpay ! udpsink host=192.174.90.10 port=2001 sync=true ",
					},
					WorkingDir: "/opt/nvidia/deepstream/deepstream-6.3",
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "hugepage-1gi",
							MountPath: "/dev/hugepages",
						},
						{
							Name:      "host-nvidia-mps",
							MountPath: "/tmp/nvidia-mps",
						},
						{
							Name:      "dpdk",
							MountPath: "/var/run/dpdk",
						},
					},
					Resources: corev1.ResourceRequirements{
						Limits:   podLimits,
						Requests: podRequests,
					},
					SecurityContext: &corev1.SecurityContext{
						Privileged: &privileged,
					},
				},
			}

			Expect(gpupod.Spec.Containers[0].Name).To(Equal(expectedPodSpecContainer[0].Name))
			Expect(gpupod.Spec.Containers[0].Image).To(Equal(expectedPodSpecContainer[0].Image))
			Expect(gpupod.Spec.Containers[0].ImagePullPolicy).To(Equal(expectedPodSpecContainer[0].ImagePullPolicy))
			Expect(gpupod.Spec.Containers[0].Args).To(Equal(expectedPodSpecContainer[0].Args))
			Expect(gpupod.Spec.Containers[0].Command).To(Equal(expectedPodSpecContainer[0].Command))
			Expect(gpupod.Spec.Containers[0].SecurityContext).To(Equal(expectedPodSpecContainer[0].SecurityContext))
			Expect(gpupod.Spec.Containers[0].WorkingDir).To(Equal(expectedPodSpecContainer[0].WorkingDir))
			Expect(gpupod.Spec.Containers[0].VolumeMounts).To(Equal(expectedPodSpecContainer[0].VolumeMounts))
			Expect(gpupod.Spec.Containers[0].Resources).To(Equal(expectedPodSpecContainer[0].Resources))

			volumeType := corev1.HostPathType("")
			expectedPodSpecVolumes := []corev1.Volume{
				{
					Name: "hugepage-1gi",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/dev/hugepages",
							Type: &volumeType,
						},
					},
				},
				{
					Name: "host-nvidia-mps",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/tmp/nvidia-mps",
							Type: &volumeType,
						},
					},
				},
				{
					Name: "dpdk",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/var/run/dpdk",
							Type: &volumeType,
						},
					},
				},
			}
			Expect(gpupod.Spec.Volumes).To(Equal(expectedPodSpecVolumes))

			var expectedPodSpecrResartPolicy corev1.RestartPolicy = "Always"
			expectedPodSpecHostNW := false
			expectedPodSpecHostIPC := true
			Expect(gpupod.Spec.RestartPolicy).To(Equal(expectedPodSpecrResartPolicy))
			Expect(gpupod.Spec.HostNetwork).To(Equal(expectedPodSpecHostNW))
			Expect(gpupod.Spec.HostIPC).To(Equal(expectedPodSpecHostIPC))

			expectedPodSpecAffinity := &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
						NodeSelectorTerms: []corev1.NodeSelectorTerm{
							{
								MatchExpressions: []corev1.NodeSelectorRequirement{
									{
										Key:      "kubernetes.io/hostname",
										Operator: corev1.NodeSelectorOpIn,
										Values: []string{
											"worker1",
										},
									},
								},
							},
						},
					},
				},
			}
			Expect(gpupod.Spec.Affinity).To(Equal(expectedPodSpecAffinity))
		})
		It("Test_1-1-2_cpu-filter-resize-tcp-high-infer", func() {
			// Create GPUFuncConfig
			err = createConfig(ctx, gpuconfighigh)
			if err != nil {
				fmt.Println("There is a problem in createing configmap ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection112)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunctionFilterResize112)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create NetworkAttachmentDifinition
			err = createNetworkAttachmentDefinition(ctx, NetworkAttachmentDefinition1)
			if err != nil {
				fmt.Println("There is a problem in createing NetworkAttachmentDefinition ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction112)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "gpufunctest112-wbfunction-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var gpuCR examplecomv1.GPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "gpufunctest112-wbfunction-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&gpuCR)
			Expect(err).NotTo(HaveOccurred())

			var addTest = true
			var rxTest = "TCP"
			var txTest = "RTP"
			var functionIndex int32 = 1
			var podName = "gpufunctest112-wbfunction-high-infer-main-mps-dgpu-0-pod"
			// test for updating CR.Status
			expectedStatus := examplecomv1.GPUFunctionStatus{
				StartTime: gpuCR.Status.StartTime,
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "gpufunctest112",
					Namespace: "default",
				},
				FunctionName: "high-infer",
				ImageURI:     "localhost/gpu-deepstream-app_tcprcv:3.1.0",
				ConfigName:   "gpufunc-config-high-infer",
				Status:       "Running",
				SharedMemory: &examplecomv1.SharedMemorySpec{
					FilePrefix:      "gpufunctest112-wbfunction-high-infer-main",
					CommandQueueID:  "gpufunctest112-wbfunction-high-infer-main",
					SharedMemoryMiB: 1,
				},
				VirtualNetworkDeviceDriverType: "sriov",
				AdditionalNetwork:              &addTest,
				FunctionIndex:                  &functionIndex,
				RxProtocol:                     &rxTest,
				TxProtocol:                     &txTest,
				PodName:                        &podName,
			}

			// test for updating StartTime
			err = checkTime(testTime.Time, gpuCR.Status.StartTime.Time)
			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status except for Status.StartTime
			Expect(gpuCR.Status).To(Equal(expectedStatus))

			// test for creating finalizer
			expectedFinalizer := []string{
				"gpufunction.finalizers.example.com.v1",
			}
			Expect(gpuCR.Finalizers).To(Equal(expectedFinalizer))

			// check logs
			Expect(writer.String()).To(ContainSubstring("previousConnectionCRName : gpufunctest112-wbconnection-filter-resize-high-infer-main-high-infer-main"))
			Expect(writer.String()).To(ContainSubstring("rxProtocol : TCP"))
			Expect(writer.String()).To(ContainSubstring("txProtocol : RTP"))
			Expect(writer.String()).To(ContainSubstring("Success to create Pod."))
			Expect(writer.String()).To(ContainSubstring("Status Update."))

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

			var gpupod corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "gpufunctest112-wbfunction-high-infer-main-mps-dgpu-0-pod",
				Namespace: TESTNAMESPACE,
			},
				&gpupod)
			if err != nil {
				// error
				fmt.Println("Could not get Pod:", gpupod, err)
			}
			Expect(err).NotTo(HaveOccurred())

			expectedPodMeta := metav1.ObjectMeta{
				Name:      "gpufunctest112-wbfunction-high-infer-main-mps-dgpu-0-pod",
				Namespace: "default",
				Annotations: map[string]string{
					"k8s.v1.cni.cncf.io/networks": "[ {\"name\": \"worker1-config-net-sriov\",\"ips\": [\"192.174.90.142/24\"] } ]",
				},
			}

			Expect(gpupod.ObjectMeta.Name).To(Equal(expectedPodMeta.Name))
			Expect(gpupod.ObjectMeta.Namespace).To(Equal(expectedPodMeta.Namespace))
			Expect(gpupod.ObjectMeta.Annotations["k8s.v1.cni.cncf.io/networks"]).To(Equal(expectedPodMeta.Annotations["k8s.v1.cni.cncf.io/networks"]))

			expectedPodSpecContainerEnv := []corev1.EnvVar{
				{
					Name:  "CUDA_MPS_LOG_DIRECTORY",
					Value: "/tmp/nvidia-mps",
				},
				{
					Name:  "CUDA_MPS_PIPE_DIRECTORY",
					Value: "/tmp/nvidia-mps",
				},
				{
					Name:  "CUDA_VISIBLE_DEVICES",
					Value: "0",
				},
				{
					Name:  "GST_PLUGIN_PATH",
					Value: "/opt/nvidia/deepstream/deepstream-6.3/fpga-software/tools/tcp_plugins/fpga_depayloader",
				},
				{
					Name:  "HEIGHT",
					Value: "1280",
				},
				{
					Name:  "RECEIVING_SERVER_IP",
					Value: "192.174.90.10",
				},
				{
					Name:  "RECEIVING_SERVER_PORT",
					Value: "2011",
				},
				{
					Name:  "WIDTH",
					Value: "1280",
				},
				{
					Name:  "K8S_POD_NAME",
					Value: "",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							APIVersion: "v1",
							FieldPath:  "metadata.name",
						},
						ResourceFieldRef: nil,
						ConfigMapKeyRef:  nil,
						SecretKeyRef:     nil,
					},
				},
				{
					Name:  "K8S_POD_NAMESPACE",
					Value: "",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							APIVersion: "v1",
							FieldPath:  "metadata.namespace",
						},
						ResourceFieldRef: nil,
						ConfigMapKeyRef:  nil,
						SecretKeyRef:     nil,
					},
				},
			}

			sort.Slice(expectedPodSpecContainerEnv, func(i, j int) bool {
				return expectedPodSpecContainerEnv[i].Name < expectedPodSpecContainerEnv[j].Name
			})
			var gpupodEnv []corev1.EnvVar
			for _, containers := range gpupod.Spec.Containers {
				gpupodEnv = append(gpupodEnv, containers.Env...)
			}

			sort.Slice(gpupodEnv, func(i, j int) bool {
				return gpupodEnv[i].Name < gpupodEnv[j].Name
			})

			Expect(gpupodEnv).To(Equal(expectedPodSpecContainerEnv))

			var podLimits = corev1.ResourceList{
				"intel.com/intel_sriov_netdevice": {},
			}
			podLimits["intel.com/intel_sriov_netdevice"] = resource.Quantity(resource.MustParse("1"))

			var podRequests = corev1.ResourceList{
				"intel.com/intel_sriov_netdevice": {},
			}
			podRequests["intel.com/intel_sriov_netdevice"] = resource.Quantity(resource.MustParse("1"))

			var privileged bool = true
			expectedPodSpecContainer := []corev1.Container{
				{
					Name:  "gpu-container0",
					Image: "localhost/gpu-deepstream-app_tcprcv:3.1.0",
					Command: []string{
						"sh",
						"-c",
					},
					ImagePullPolicy: corev1.PullIfNotPresent,
					Args: []string{
						"cd /opt/DeepStream-Yolo && gst-launch-1.0 -ev fpgadepay host=192.174.90.142 port=15001 ! 'video/x-raw,format=(string)BGR,width=1280,height=1280' ! nvvideoconvert ! 'video/x-raw(memory:NVMM), format=(string)RGBA' ! m.sink_0 nvstreammux name=m nvbuf-memory-type=0 batch-size=1 width=1280 height=1280 ! queue ! nvinfer config-file-path=./config_infer_primary_yoloV4_p6_th020_040.txt batch-size=1 model-engine-file=./model_b1_gpu0_fp16.engine ! queue ! nvdsosd process-mode=1 ! nvvideoconvert ! 'video/x-raw, format=(string)BGR' ! videoconvert ! queue ! perf ! rtpvrawpay ! udpsink host=192.174.90.10 port=2011 sync=true ",
					},
					WorkingDir: "/opt/nvidia/deepstream/deepstream-6.3",
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "host-nvidia-mps",
							MountPath: "/tmp/nvidia-mps",
						},
					},
					Resources: corev1.ResourceRequirements{
						Limits:   podLimits,
						Requests: podRequests,
					},
					SecurityContext: &corev1.SecurityContext{
						Privileged: &privileged,
					},
				},
			}

			Expect(gpupod.Spec.Containers[0].Name).To(Equal(expectedPodSpecContainer[0].Name))
			Expect(gpupod.Spec.Containers[0].Image).To(Equal(expectedPodSpecContainer[0].Image))
			Expect(gpupod.Spec.Containers[0].ImagePullPolicy).To(Equal(expectedPodSpecContainer[0].ImagePullPolicy))
			Expect(gpupod.Spec.Containers[0].Args).To(Equal(expectedPodSpecContainer[0].Args))
			Expect(gpupod.Spec.Containers[0].Command).To(Equal(expectedPodSpecContainer[0].Command))
			Expect(gpupod.Spec.Containers[0].SecurityContext).To(Equal(expectedPodSpecContainer[0].SecurityContext))
			Expect(gpupod.Spec.Containers[0].WorkingDir).To(Equal(expectedPodSpecContainer[0].WorkingDir))
			Expect(gpupod.Spec.Containers[0].VolumeMounts).To(Equal(expectedPodSpecContainer[0].VolumeMounts))
			Expect(gpupod.Spec.Containers[0].Resources).To(Equal(expectedPodSpecContainer[0].Resources))

			volumeType := corev1.HostPathType("")
			expectedPodSpecVolumes := []corev1.Volume{
				{
					Name: "host-nvidia-mps",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/tmp/nvidia-mps",
							Type: &volumeType,
						},
					},
				},
			}
			Expect(gpupod.Spec.Volumes).To(Equal(expectedPodSpecVolumes))

			var expectedPodSpecrResartPolicy corev1.RestartPolicy = "Always"
			expectedPodSpecHostNW := false
			expectedPodSpecHostIPC := true
			Expect(gpupod.Spec.RestartPolicy).To(Equal(expectedPodSpecrResartPolicy))
			Expect(gpupod.Spec.HostNetwork).To(Equal(expectedPodSpecHostNW))
			Expect(gpupod.Spec.HostIPC).To(Equal(expectedPodSpecHostIPC))

			expectedPodSpecAffinity := &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
						NodeSelectorTerms: []corev1.NodeSelectorTerm{
							{
								MatchExpressions: []corev1.NodeSelectorRequirement{
									{
										Key:      "kubernetes.io/hostname",
										Operator: corev1.NodeSelectorOpIn,
										Values: []string{
											"worker1",
										},
									},
								},
							},
						},
					},
				},
			}
			Expect(gpupod.Spec.Affinity).To(Equal(expectedPodSpecAffinity))
		})
		//Test for GetFunc
		It("Test_1-2-1_fpga-filter-resize-dma-low-infer", func() {
			// Create GPUFuncConfig
			err = createConfig(ctx, gpuconfiglow)
			if err != nil {
				fmt.Println("There is a problem in createing configmap ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection121)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGAFunctionCR
			err = createFPGAFunction(ctx, FPGAFunction121)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create NetworkAttachmentDifinition
			err = createNetworkAttachmentDefinition(ctx, NetworkAttachmentDefinition1)
			if err != nil {
				fmt.Println("There is a problem in createing NetworkAttachmentDefinition ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction121)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "gpufunctest121-wbfunction-low-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var gpuCR examplecomv1.GPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "gpufunctest121-wbfunction-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&gpuCR)
			Expect(err).NotTo(HaveOccurred())

			var addTest = true
			var rxTest = "DMA"
			var txTest = "RTP"
			var functionIndex int32 = 2
			var podName = "gpufunctest121-wbfunction-low-infer-main-mps-dgpu-0-pod"
			// test for updating CR.Status
			expectedStatus := examplecomv1.GPUFunctionStatus{
				StartTime: gpuCR.Status.StartTime,
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "gpufunctest121",
					Namespace: "default",
				},
				FunctionName: "low-infer",
				ImageURI:     "localhost/gpu-deepstream-app:3.1.0",
				ConfigName:   "gpufunc-config-low-infer",
				Status:       "Running",
				SharedMemory: &examplecomv1.SharedMemorySpec{
					FilePrefix:      "gpufunctest121-wbfunction-low-infer-main",
					CommandQueueID:  "gpufunctest121-wbfunction-low-infer-main",
					SharedMemoryMiB: 1,
				},
				VirtualNetworkDeviceDriverType: "sriov",
				AdditionalNetwork:              &addTest,
				FunctionIndex:                  &functionIndex,
				RxProtocol:                     &rxTest,
				TxProtocol:                     &txTest,
				PodName:                        &podName,
			}

			// test for updating StartTime
			err = checkTime(testTime.Time, gpuCR.Status.StartTime.Time)
			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status except for Status.StartTime
			Expect(gpuCR.Status).To(Equal(expectedStatus))

			// test for creating finalizer
			expectedFinalizer := []string{
				"gpufunction.finalizers.example.com.v1",
			}
			Expect(gpuCR.Finalizers).To(Equal(expectedFinalizer))

			// check logs
			Expect(writer.String()).To(ContainSubstring("previousConnectionCRName : gpufunctest121-wbconnection-filter-resize-low-infer-main-low-infer-main"))
			Expect(writer.String()).To(ContainSubstring("rxProtocol : DMA"))
			Expect(writer.String()).To(ContainSubstring("txProtocol : RTP"))
			Expect(writer.String()).To(ContainSubstring("Success to create Pod."))
			Expect(writer.String()).To(ContainSubstring("Status Update."))

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

			var gpupod corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "gpufunctest121-wbfunction-low-infer-main-mps-dgpu-0-pod",
				Namespace: TESTNAMESPACE,
			},
				&gpupod)
			if err != nil {
				// error
				fmt.Println("Could not get Pod:", gpupod, err)
			}
			Expect(err).NotTo(HaveOccurred())

			expectedPodMeta := metav1.ObjectMeta{
				Name:      "gpufunctest121-wbfunction-low-infer-main-mps-dgpu-0-pod",
				Namespace: "default",
				Annotations: map[string]string{
					"k8s.v1.cni.cncf.io/networks": "[ {\"name\": \"worker1-config-net-sriov\",\"ips\": [\"192.174.90.141/24\"] } ]",
				},
			}

			Expect(gpupod.ObjectMeta.Name).To(Equal(expectedPodMeta.Name))
			Expect(gpupod.ObjectMeta.Namespace).To(Equal(expectedPodMeta.Namespace))
			Expect(gpupod.ObjectMeta.Annotations["k8s.v1.cni.cncf.io/networks"]).To(Equal(expectedPodMeta.Annotations["k8s.v1.cni.cncf.io/networks"]))

			expectedPodSpecContainerEnv := []corev1.EnvVar{
				{
					Name:  "CONNECTOR_ID",
					Value: "gpufunctest121-wbfunction-low-infer-main",
				},
				{
					Name:  "CUDA_MPS_LOG_DIRECTORY",
					Value: "/tmp/nvidia-mps",
				},
				{
					Name:  "CUDA_MPS_PIPE_DIRECTORY",
					Value: "/tmp/nvidia-mps",
				},
				{
					Name:  "CUDA_VISIBLE_DEVICES",
					Value: "0",
				},
				{
					Name:  "FILE_PREFIX",
					Value: "gpufunctest121-wbfunction-low-infer-main",
				},
				{
					Name:  "FPGA_DEV",
					Value: "/dev/xpcie_21330621T00D",
				},
				{
					Name:  "HEIGHT",
					Value: "416",
				},
				{
					Name:  "K8S_POD_NAME",
					Value: "",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							APIVersion: "v1",
							FieldPath:  "metadata.name",
						},
						ResourceFieldRef: nil,
						ConfigMapKeyRef:  nil,
						SecretKeyRef:     nil,
					},
				},
				{
					Name:  "K8S_POD_NAMESPACE",
					Value: "",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							APIVersion: "v1",
							FieldPath:  "metadata.namespace",
						},
						ResourceFieldRef: nil,
						ConfigMapKeyRef:  nil,
						SecretKeyRef:     nil,
					},
				},
				{
					Name:  "RECEIVING_SERVER_IP",
					Value: "192.174.90.10",
				},
				{
					Name:  "RECEIVING_SERVER_PORT",
					Value: "2001",
				},
				{
					Name:  "SHMEM_SECONDARY",
					Value: "1",
				},
				{
					Name:  "WIDTH",
					Value: "416",
				},
			}

			sort.Slice(expectedPodSpecContainerEnv, func(i, j int) bool {
				return expectedPodSpecContainerEnv[i].Name < expectedPodSpecContainerEnv[j].Name
			})
			var gpupodEnv []corev1.EnvVar
			for _, containers := range gpupod.Spec.Containers {
				gpupodEnv = append(gpupodEnv, containers.Env...)
			}

			sort.Slice(gpupodEnv, func(i, j int) bool {
				return gpupodEnv[i].Name < gpupodEnv[j].Name
			})

			Expect(gpupodEnv).To(Equal(expectedPodSpecContainerEnv))

			var podLimits = corev1.ResourceList{
				"hugepages-1Gi":                   resource.Quantity{},
				"intel.com/intel_sriov_netdevice": {},
			}
			podLimits["hugepages-1Gi"] = resource.Quantity(resource.MustParse("1Gi"))
			podLimits["intel.com/intel_sriov_netdevice"] = resource.Quantity(resource.MustParse("1"))

			var podRequests = corev1.ResourceList{
				"hugepages-1Gi":                   resource.Quantity{},
				"intel.com/intel_sriov_netdevice": {},
				"memory":                          {},
			}
			podRequests["hugepages-1Gi"] = resource.Quantity(resource.MustParse("1Gi"))
			podRequests["intel.com/intel_sriov_netdevice"] = resource.Quantity(resource.MustParse("1"))
			podRequests["memory"] = resource.Quantity(resource.MustParse("32Gi"))

			var privileged bool = true
			expectedPodSpecContainer := []corev1.Container{
				{
					Name:  "gpu-container0",
					Image: "localhost/gpu-deepstream-app:3.1.0",
					Command: []string{
						"sh",
						"-c",
					},
					ImagePullPolicy: corev1.PullIfNotPresent,
					Args: []string{
						"cd /opt/nvidia/deepstream/deepstream-6.3/sources/objectDetector_Yolo/ && gst-launch-1.0 -ev fpgasrc ! 'video/x-raw,format=(string)BGR,width=416,height=416' ! nvvideoconvert ! 'video/x-raw(memory:NVMM), format=(string)RGBA' ! m.sink_0 nvstreammux name=m nvbuf-memory-type=0 batch-size=1  width=416 height=416 ! queue ! nvinfer config-file-path=./config_infer_primary_yoloV3_tiny.txt batch-size=1 model-engine-file=./model_b1_gpu0_int8.engine ! queue ! nvvideoconvert ! 'video/x-raw, format=(string)BGR' ! videoconvert ! queue ! perf ! rtpvrawpay ! udpsink host=192.174.90.10 port=2001 sync=true ",
					},
					WorkingDir: "/opt/nvidia/deepstream/deepstream-6.3",
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "hugepage-1gi",
							MountPath: "/dev/hugepages",
						},
						{
							Name:      "host-nvidia-mps",
							MountPath: "/tmp/nvidia-mps",
						},
						{
							Name:      "dpdk",
							MountPath: "/var/run/dpdk",
						},
					},
					Resources: corev1.ResourceRequirements{
						Limits:   podLimits,
						Requests: podRequests,
					},
					SecurityContext: &corev1.SecurityContext{
						Privileged: &privileged,
					},
				},
			}

			Expect(gpupod.Spec.Containers[0].Name).To(Equal(expectedPodSpecContainer[0].Name))
			Expect(gpupod.Spec.Containers[0].Image).To(Equal(expectedPodSpecContainer[0].Image))
			Expect(gpupod.Spec.Containers[0].ImagePullPolicy).To(Equal(expectedPodSpecContainer[0].ImagePullPolicy))
			Expect(gpupod.Spec.Containers[0].Args).To(Equal(expectedPodSpecContainer[0].Args))
			Expect(gpupod.Spec.Containers[0].Command).To(Equal(expectedPodSpecContainer[0].Command))
			Expect(gpupod.Spec.Containers[0].SecurityContext).To(Equal(expectedPodSpecContainer[0].SecurityContext))
			Expect(gpupod.Spec.Containers[0].WorkingDir).To(Equal(expectedPodSpecContainer[0].WorkingDir))
			Expect(gpupod.Spec.Containers[0].VolumeMounts).To(Equal(expectedPodSpecContainer[0].VolumeMounts))
			Expect(gpupod.Spec.Containers[0].Resources).To(Equal(expectedPodSpecContainer[0].Resources))

			volumeType := corev1.HostPathType("")
			expectedPodSpecVolumes := []corev1.Volume{
				{
					Name: "hugepage-1gi",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/dev/hugepages",
							Type: &volumeType,
						},
					},
				},
				{
					Name: "host-nvidia-mps",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/tmp/nvidia-mps",
							Type: &volumeType,
						},
					},
				},
				{
					Name: "dpdk",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/var/run/dpdk",
							Type: &volumeType,
						},
					},
				},
			}
			Expect(gpupod.Spec.Volumes).To(Equal(expectedPodSpecVolumes))

			var expectedPodSpecrResartPolicy corev1.RestartPolicy = "Always"
			expectedPodSpecHostNW := false
			expectedPodSpecHostIPC := true
			Expect(gpupod.Spec.RestartPolicy).To(Equal(expectedPodSpecrResartPolicy))
			Expect(gpupod.Spec.HostNetwork).To(Equal(expectedPodSpecHostNW))
			Expect(gpupod.Spec.HostIPC).To(Equal(expectedPodSpecHostIPC))

			expectedPodSpecAffinity := &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
						NodeSelectorTerms: []corev1.NodeSelectorTerm{
							{
								MatchExpressions: []corev1.NodeSelectorRequirement{
									{
										Key:      "kubernetes.io/hostname",
										Operator: corev1.NodeSelectorOpIn,
										Values: []string{
											"worker1",
										},
									},
								},
							},
						},
					},
				},
			}
			Expect(gpupod.Spec.Affinity).To(Equal(expectedPodSpecAffinity))
		})
		//Test for GetFunc
		It("Test_1-2-2_cpu-filter-resize-tcp-low-infer", func() {
			// Create GPUFuncConfig
			err = createConfig(ctx, gpuconfiglow)
			if err != nil {
				fmt.Println("There is a problem in createing configmap ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection122)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunctionFilterResize122)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create NetworkAttachmentDifinition
			err = createNetworkAttachmentDefinition(ctx, NetworkAttachmentDefinition1)
			if err != nil {
				fmt.Println("There is a problem in createing NetworkAttachmentDefinition ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction122)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "gpufunctest122-wbfunction-low-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var gpuCR examplecomv1.GPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "gpufunctest122-wbfunction-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&gpuCR)
			Expect(err).NotTo(HaveOccurred())

			var addTest = true
			var rxTest = "TCP"
			var txTest = "RTP"
			var functionIndex int32 = 3
			var podName = "gpufunctest122-wbfunction-low-infer-main-mps-dgpu-0-pod"
			// test for updating CR.Status
			expectedStatus := examplecomv1.GPUFunctionStatus{
				StartTime: gpuCR.Status.StartTime,
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "gpufunctest122",
					Namespace: "default",
				},
				FunctionName: "low-infer",
				ImageURI:     "localhost/gpu-deepstream-app_tcprcv:3.1.0",
				ConfigName:   "gpufunc-config-low-infer",
				Status:       "Running",
				SharedMemory: &examplecomv1.SharedMemorySpec{
					FilePrefix:      "gpufunctest122-wbfunction-low-infer-main",
					CommandQueueID:  "gpufunctest122-wbfunction-low-infer-main",
					SharedMemoryMiB: 1,
				},
				VirtualNetworkDeviceDriverType: "sriov",
				AdditionalNetwork:              &addTest,
				FunctionIndex:                  &functionIndex,
				RxProtocol:                     &rxTest,
				TxProtocol:                     &txTest,
				PodName:                        &podName,
			}

			// test for updating StartTime
			err = checkTime(testTime.Time, gpuCR.Status.StartTime.Time)
			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status except for Status.StartTime
			Expect(gpuCR.Status).To(Equal(expectedStatus))

			// test for creating finalizer
			expectedFinalizer := []string{
				"gpufunction.finalizers.example.com.v1",
			}
			Expect(gpuCR.Finalizers).To(Equal(expectedFinalizer))

			// check logs
			Expect(writer.String()).To(ContainSubstring("previousConnectionCRName : gpufunctest122-wbconnection-filter-resize-low-infer-main-low-infer-main"))
			Expect(writer.String()).To(ContainSubstring("rxProtocol : TCP"))
			Expect(writer.String()).To(ContainSubstring("txProtocol : RTP"))
			Expect(writer.String()).To(ContainSubstring("Success to create Pod."))
			Expect(writer.String()).To(ContainSubstring("Status Update."))

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

			var gpupod corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "gpufunctest122-wbfunction-low-infer-main-mps-dgpu-0-pod",
				Namespace: TESTNAMESPACE,
			},
				&gpupod)
			if err != nil {
				// error
				fmt.Println("Could not get Pod:", gpupod, err)
			}
			Expect(err).NotTo(HaveOccurred())

			expectedPodMeta := metav1.ObjectMeta{
				Name:      "gpufunctest122-wbfunction-low-infer-main-mps-dgpu-0-pod",
				Namespace: "default",
				Annotations: map[string]string{
					"k8s.v1.cni.cncf.io/networks": "[ {\"name\": \"worker1-config-net-sriov\",\"ips\": [\"192.174.90.142/24\"] } ]",
				},
			}

			Expect(gpupod.ObjectMeta.Name).To(Equal(expectedPodMeta.Name))
			Expect(gpupod.ObjectMeta.Namespace).To(Equal(expectedPodMeta.Namespace))
			Expect(gpupod.ObjectMeta.Annotations["k8s.v1.cni.cncf.io/networks"]).To(Equal(expectedPodMeta.Annotations["k8s.v1.cni.cncf.io/networks"]))

			expectedPodSpecContainerEnv := []corev1.EnvVar{
				{
					Name:  "CUDA_MPS_LOG_DIRECTORY",
					Value: "/tmp/nvidia-mps",
				},
				{
					Name:  "CUDA_MPS_PIPE_DIRECTORY",
					Value: "/tmp/nvidia-mps",
				},
				{
					Name:  "CUDA_VISIBLE_DEVICES",
					Value: "0",
				},
				{
					Name:  "GST_PLUGIN_PATH",
					Value: "/opt/nvidia/deepstream/deepstream-6.3/fpga-software/tools/tcp_plugins/fpga_depayloader",
				},
				{
					Name:  "HEIGHT",
					Value: "416",
				},
				{
					Name:  "RECEIVING_SERVER_IP",
					Value: "192.174.90.10",
				},
				{
					Name:  "RECEIVING_SERVER_PORT",
					Value: "2011",
				},
				{
					Name:  "WIDTH",
					Value: "416",
				},
				{
					Name:  "K8S_POD_NAME",
					Value: "",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							APIVersion: "v1",
							FieldPath:  "metadata.name",
						},
						ResourceFieldRef: nil,
						ConfigMapKeyRef:  nil,
						SecretKeyRef:     nil,
					},
				},
				{
					Name:  "K8S_POD_NAMESPACE",
					Value: "",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							APIVersion: "v1",
							FieldPath:  "metadata.namespace",
						},
						ResourceFieldRef: nil,
						ConfigMapKeyRef:  nil,
						SecretKeyRef:     nil,
					},
				},
			}

			sort.Slice(expectedPodSpecContainerEnv, func(i, j int) bool {
				return expectedPodSpecContainerEnv[i].Name < expectedPodSpecContainerEnv[j].Name
			})
			var gpupodEnv []corev1.EnvVar
			for _, containers := range gpupod.Spec.Containers {
				gpupodEnv = append(gpupodEnv, containers.Env...)
			}

			sort.Slice(gpupodEnv, func(i, j int) bool {
				return gpupodEnv[i].Name < gpupodEnv[j].Name
			})

			Expect(gpupodEnv).To(Equal(expectedPodSpecContainerEnv))

			var podLimits = corev1.ResourceList{
				"intel.com/intel_sriov_netdevice": {},
			}
			podLimits["intel.com/intel_sriov_netdevice"] = resource.Quantity(resource.MustParse("1"))

			var podRequests = corev1.ResourceList{
				"intel.com/intel_sriov_netdevice": {},
			}
			podRequests["intel.com/intel_sriov_netdevice"] = resource.Quantity(resource.MustParse("1"))

			var privileged bool = true
			expectedPodSpecContainer := []corev1.Container{
				{
					Name:  "gpu-container0",
					Image: "localhost/gpu-deepstream-app_tcprcv:3.1.0",
					Command: []string{
						"sh",
						"-c",
					},
					ImagePullPolicy: corev1.PullIfNotPresent,
					Args: []string{
						"cd /opt/nvidia/deepstream/deepstream-6.3/sources/objectDetector_Yolo/ && gst-launch-1.0 -ev fpgadepay host=192.174.90.142 port=15001 ! 'video/x-raw,format=(string)BGR,width=416,height=416' ' ! nvvideoconvert ! 'video/x-raw(memory:NVMM), format=(string)RGBA' ! m.sink_0 nvstreammux name=m nvbuf-memory-type=0 batch-size=1  width=416 height=416 ! queue ! nvinfer config-file-path=./config_infer_primary_yoloV3_tiny.txt batch-size=1 model-engine-file=./model_b1_gpu0_int8.engine ! queue ! nvvideoconvert ! 'video/x-raw, format=(string)BGR' ! videoconvert ! queue ! perf ! rtpvrawpay ! udpsink host=192.174.90.10 port=2011 sync=true ",
					},
					WorkingDir: "/opt/nvidia/deepstream/deepstream-6.3",
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "host-nvidia-mps",
							MountPath: "/tmp/nvidia-mps",
						},
					},
					Resources: corev1.ResourceRequirements{
						Limits:   podLimits,
						Requests: podRequests,
					},
					SecurityContext: &corev1.SecurityContext{
						Privileged: &privileged,
					},
				},
			}

			Expect(gpupod.Spec.Containers[0].Name).To(Equal(expectedPodSpecContainer[0].Name))
			Expect(gpupod.Spec.Containers[0].Image).To(Equal(expectedPodSpecContainer[0].Image))
			Expect(gpupod.Spec.Containers[0].ImagePullPolicy).To(Equal(expectedPodSpecContainer[0].ImagePullPolicy))
			Expect(gpupod.Spec.Containers[0].Args).To(Equal(expectedPodSpecContainer[0].Args))
			Expect(gpupod.Spec.Containers[0].Command).To(Equal(expectedPodSpecContainer[0].Command))
			Expect(gpupod.Spec.Containers[0].SecurityContext).To(Equal(expectedPodSpecContainer[0].SecurityContext))
			Expect(gpupod.Spec.Containers[0].WorkingDir).To(Equal(expectedPodSpecContainer[0].WorkingDir))
			Expect(gpupod.Spec.Containers[0].VolumeMounts).To(Equal(expectedPodSpecContainer[0].VolumeMounts))
			Expect(gpupod.Spec.Containers[0].Resources).To(Equal(expectedPodSpecContainer[0].Resources))

			volumeType := corev1.HostPathType("")
			expectedPodSpecVolumes := []corev1.Volume{
				{
					Name: "host-nvidia-mps",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/tmp/nvidia-mps",
							Type: &volumeType,
						},
					},
				},
			}
			Expect(gpupod.Spec.Volumes).To(Equal(expectedPodSpecVolumes))

			var expectedPodSpecrResartPolicy corev1.RestartPolicy = "Always"
			expectedPodSpecHostNW := false
			expectedPodSpecHostIPC := true
			Expect(gpupod.Spec.RestartPolicy).To(Equal(expectedPodSpecrResartPolicy))
			Expect(gpupod.Spec.HostNetwork).To(Equal(expectedPodSpecHostNW))
			Expect(gpupod.Spec.HostIPC).To(Equal(expectedPodSpecHostIPC))

			expectedPodSpecAffinity := &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
						NodeSelectorTerms: []corev1.NodeSelectorTerm{
							{
								MatchExpressions: []corev1.NodeSelectorRequirement{
									{
										Key:      "kubernetes.io/hostname",
										Operator: corev1.NodeSelectorOpIn,
										Values: []string{
											"worker1",
										},
									},
								},
							},
						},
					},
				},
			}
			Expect(gpupod.Spec.Affinity).To(Equal(expectedPodSpecAffinity))
		})

		It("Test_2-1_UPDATE", func() {

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, gpufunctestUPDATE)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFunctionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "gpufunctestupdate-wbfunction-high-infer-main",
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
				"Normal Update Update Start",
				"Normal Update Update End",
			))
		})
		It("Test_2-2_DELETE", func() {

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, gpufunctestDELETE)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFunctionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			var gpuCR examplecomv1.GPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "gpufunctestdelete-wbfunction-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&gpuCR)

			err = k8sClient.Delete(ctx, &gpuCR)

			var gpubCR examplecomv1.GPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "gpufunctestdelete-wbfunction-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&gpubCR)

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "gpufunctestdelete-wbfunction-high-infer-main",
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

			var gpuaCR examplecomv1.GPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "gpufunctestdelete-wbfunction-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&gpuaCR)

			Expect(err).To(MatchError(ContainSubstring("not found")))

		})

		It("3-1-4 Config in LifeCycle", func() {
			By("Test Start")

			// Create GPUFuncConfig
			err = createConfig(ctx, gpuconfigdecode314)
			if err != nil {
				fmt.Println("There is a problem in createing GPUConfig ", err)
				fmt.Printf("%T\n", err)
				fmt.Println(err)
			}

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection1)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}

			// Create FPGAFunctionCR
			err = createFPGAFunction(ctx, FPGAFunction1)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}

			// Create NetworkAttachmentDefinition
			err = createNetworkAttachmentDefinition(ctx, NetworkAttachmentDefinition1)
			if err != nil {
				fmt.Println("There is a problem in createing NetworkAttachmentDefinition")
			}

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction314)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR")
			}

			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night314-wbfunction-high-infer-main",
			}})

			if err != nil {
				By("Reconcile Error")
			}
			Expect(err).NotTo(HaveOccurred())

			var gpuCR examplecomv1.GPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "df-night314-wbfunction-high-infer-main",
				Namespace: "default",
			}, &gpuCR)

			if err != nil {
				// Error route
				fmt.Println("Cannot get GPUFunctionCR:", gpuCR, err)
			}

			var podCR corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      *gpuCR.Status.PodName,
				Namespace: "default",
			}, &podCR)
			if err != nil {
				// Error route
				fmt.Println("Cannot get PodCR:", podCR, err)
			}

			var licynil corev1.Lifecycle
			Expect(podCR.Spec.Containers[0].Lifecycle).NotTo(HaveValue(Equal(licynil)))
		})
		AfterEach(func() {
			By("Test End")
		})
	})
})

func checkTime(t1 time.Time, t2 time.Time) error {
	if t1.Before(t2) {
		return nil
	} else {
		return fmt.Errorf("CR.Status may not be updated. t1 (%s) is after t2 (%s)", t1.Format(time.RFC3339), t2.Format(time.RFC3339))
	}
}

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
	"time"

	"go.uber.org/zap/zapcore"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	examplecomv1 "CPUFunction/api/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"

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

// Create CPUFunction CR
func createCPUFunction(ctx context.Context, cpufcr examplecomv1.CPUFunction) error {
	tmp := &examplecomv1.CPUFunction{}
	*tmp = cpufcr
	tmp.TypeMeta = cpufcr.TypeMeta
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Delete CPUFunction CR
func deleteCPUFunction(ctx context.Context, cpufcr examplecomv1.CPUFunction) error {
	tmp := &examplecomv1.CPUFunction{}
	*tmp = cpufcr
	tmp.TypeMeta = cpufcr.TypeMeta
	err := k8sClient.Delete(ctx, tmp)
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
func createConfig(ctx context.Context, cpuconf corev1.ConfigMap) error {
	tmp := &corev1.ConfigMap{}
	*tmp = cpuconf
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

// To describe test cases in CPUFunctionController
var _ = Describe("CPUFunctionController", func() {
	var mgr ctrl.Manager
	var err error
	ctx := context.Background()
	var fakerecorder = record.NewFakeRecorder(10)
	var writer = bytes.Buffer{}

	// This test case is for reconciler
	Context("Test for CPUFunctionReconciler", Ordered, func() {
		var reconciler CPUFunctionReconciler

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

		// Each time It runs, BeforeEach is executed.
		BeforeEach(func() {

			// loger initialized
			writer.Reset()

			// set environmental variable
			os.Setenv("K8S_NODENAME", "node01")

			// recorder initialized
			fakerecorder = record.NewFakeRecorder(10)

			reconciler = CPUFunctionReconciler{
				Client:   k8sClient,
				Scheme:   testScheme,
				Recorder: fakerecorder,
			}

			err = reconciler.SetupWithManager(mgr)
			if err != nil {
				fmt.Println("An error occur during setupwithManager: ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// To delete crdata
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.CPUFunction{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &corev1.Pod{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &PCIeConnection{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &EthernetConnection{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &FPGAFunction{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &NetworkAttachmentDefinition{}, client.InNamespace(TESTNAMESPACE))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

		})

		AfterAll(func() {
			writer.Reset()
		})

		//Test for GetFunc
		It("Test_1-1-1_decode-dma-FPGA", func() {
			// Create CPUFuncConfig
			err = createConfig(ctx, cpuconfigdecode)
			if err != nil {
				fmt.Println("There is a problem in createing configmap ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection1)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGAFunctionCR
			err = createFPGAFunction(ctx, FPGAFunction1)
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

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction1)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "cpufunctest1-wbfunction-decode-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var cpuCR examplecomv1.CPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctest1-wbfunction-decode-main",
				Namespace: TESTNAMESPACE,
			},
				&cpuCR)
			Expect(err).NotTo(HaveOccurred())

			var addTest = true
			var rxTest = "RTP"
			var txTest = "DMA"
			var functionIndexTest int32 = 0
			var podName = "cpufunctest1-wbfunction-decode-main-cpu-pod"
			// test for updating CR.Status
			expectedStatus := examplecomv1.CPUFunctionStatus{
				StartTime: cpuCR.Status.StartTime,
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctest1",
					Namespace: "default",
				},
				FunctionName: "cpu-decode",
				ImageURI:     "localhost/host_decode:3.1.0",
				ConfigName:   "cpufunc-config-decode",
				Status:       "Running",
				SharedMemory: &examplecomv1.SharedMemorySpec{
					FilePrefix:      "test01-cpufunctest1-wbfunction-decode-main",
					CommandQueueID:  "test01-cpufunctest1-wbfunction-decode-main",
					SharedMemoryMiB: 1,
				},
				VirtualNetworkDeviceDriverType: "sriov",
				AdditionalNetwork:              &addTest,
				FunctionIndex:                  &functionIndexTest,
				RxProtocol:                     &rxTest,
				TxProtocol:                     &txTest,
				PodName:                        &podName,
			}

			// test for updating StartTime
			err = checkTime(testTime.Time, cpuCR.Status.StartTime.Time)
			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status except for Status.StartTime
			Expect(cpuCR.Status).To(Equal(expectedStatus))

			// test for creating finalizer
			expectedFinalizer := []string{
				"cpufunction.finalizers.example.com.v1",
			}
			Expect(cpuCR.Finalizers).To(Equal(expectedFinalizer))

			// check logs
			Expect(writer.String()).To(ContainSubstring("nextConnectionCRName : cpufunctest1-wbconnection-decode-main-filter-resize-low-infer-main"))
			Expect(writer.String()).To(ContainSubstring("rxProtocol : RTP"))
			Expect(writer.String()).To(ContainSubstring("txProtocol : DMA"))
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

			var cpupod corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctest1-wbfunction-decode-main-cpu-pod",
				Namespace: TESTNAMESPACE,
			},
				&cpupod)
			if err != nil {
				// error handling
				fmt.Println("Could not get Pod:", cpupod, err)
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(*cpuCR.Status.PodName).To(Equal(cpupod.Name))

			expectedPodMeta := metav1.ObjectMeta{
				Name:      "cpufunctest1-wbfunction-decode-main-cpu-pod",
				Namespace: "default",
				Annotations: map[string]string{
					"k8s.v1.cni.cncf.io/networks": "[ {\"name\": \"node01-config-net-sriov\",\"ips\": [\"192.174.90.102/24\"] } ]",
				},
			}

			Expect(cpupod.ObjectMeta.Name).To(Equal(expectedPodMeta.Name))
			Expect(cpupod.ObjectMeta.Namespace).To(Equal(expectedPodMeta.Namespace))
			Expect(cpupod.ObjectMeta.Annotations["k8s.v1.cni.cncf.io/networks"]).To(Equal(expectedPodMeta.Annotations["k8s.v1.cni.cncf.io/networks"]))

			expectedPodSpecContainerEnv := []corev1.EnvVar{
				{
					Name:      "DECENV_APPLOG_LEVEL",
					Value:     "6",
					ValueFrom: nil,
				},
				{
					Name:      "DECENV_FRAME_WIDTH",
					Value:     "3840",
					ValueFrom: nil,
				},
				{
					Name:      "DECENV_FRAME_HEIGHT",
					Value:     "2160",
					ValueFrom: nil,
				},
				{
					Name:      "DECENV_VIDEO_CONNECT_LIMIT",
					Value:     "0",
					ValueFrom: nil,
				},
				{
					Name:      "DECENV_VIDEOSRC_PROTOCOL",
					Value:     "RTP",
					ValueFrom: nil,
				},
				{
					Name:      "DECENV_OUTDST_PROTOCOL",
					Value:     "DMA",
					ValueFrom: nil,
				},
				{
					Name:  "DECENV_VIDEOSRC_PORT",
					Value: "8556",
				},
				{
					Name:      "DECENV_FRAME_FPS",
					Value:     "5",
					ValueFrom: nil,
				},
				{
					Name:      "DECENV_FPGA_DEV_NAME",
					Value:     "/dev/xpcie_21330621T00D",
					ValueFrom: nil,
				},
				{
					Name:      "DECENV_DPDK_FILE_PREFIX",
					Value:     "test01-cpufunctest1-wbfunction-decode-main",
					ValueFrom: nil,
				},
				{
					Name:      "DECENV_FPGA_DMA_HOST_TO_DEV_CONNECTOR_ID",
					Value:     "test01-cpufunctest1-wbfunction-decode-main",
					ValueFrom: nil,
				},
				{
					Name:      "DECENV_VIDEOSRC_IPA",
					Value:     "192.174.90.102",
					ValueFrom: nil,
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
			var cpupodEnv []corev1.EnvVar
			for _, containers := range cpupod.Spec.Containers {
				cpupodEnv = append(cpupodEnv, containers.Env...)
			}

			sort.Slice(cpupodEnv, func(i, j int) bool {
				return cpupodEnv[i].Name < cpupodEnv[j].Name
			})

			Expect(cpupodEnv).To(Equal(expectedPodSpecContainerEnv))

			var cpupodLimits = corev1.ResourceList{
				"hugepages-1Gi":                   resource.Quantity{},
				"intel.com/intel_sriov_netdevice": {},
			}
			cpupodLimits["hugepages-1Gi"] = resource.Quantity(resource.MustParse("1Gi"))
			cpupodLimits["intel.com/intel_sriov_netdevice"] = resource.Quantity(resource.MustParse("1"))

			var cpupodRequests = corev1.ResourceList{
				"hugepages-1Gi":                   resource.Quantity{},
				"intel.com/intel_sriov_netdevice": {},
				"memory":                          {},
			}
			cpupodRequests["hugepages-1Gi"] = resource.Quantity(resource.MustParse("1Gi"))
			cpupodRequests["intel.com/intel_sriov_netdevice"] = resource.Quantity(resource.MustParse("1"))
			cpupodRequests["memory"] = resource.Quantity(resource.MustParse("32Gi"))

			var privileged bool = true
			expectedPodSpecContainer := []corev1.Container{
				{
					Name:  "cpu-container0",
					Image: "localhost/host_decode:3.1.0",
					Command: []string{
						"sh",
						"-c",
					},
					ImagePullPolicy: corev1.PullIfNotPresent,
					Args: []string{
						"./tools/host_decode/build/host_decode-shared ",
					},
					WorkingDir: "",
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "hugepage-1gi",
							MountPath: "/dev/hugepages",
						},
						{
							Name:      "dpdk",
							MountPath: "/var/run/dpdk",
						},
					},
					Resources: corev1.ResourceRequirements{
						Limits:   cpupodLimits,
						Requests: cpupodRequests,
					},
					SecurityContext: &corev1.SecurityContext{
						Privileged: &privileged,
					},
				},
			}

			Expect(cpupod.Spec.Containers[0].Name).To(Equal(expectedPodSpecContainer[0].Name))
			Expect(cpupod.Spec.Containers[0].Image).To(Equal(expectedPodSpecContainer[0].Image))
			Expect(cpupod.Spec.Containers[0].ImagePullPolicy).To(Equal(expectedPodSpecContainer[0].ImagePullPolicy))
			Expect(cpupod.Spec.Containers[0].Args).To(Equal(expectedPodSpecContainer[0].Args))
			Expect(cpupod.Spec.Containers[0].Command).To(Equal(expectedPodSpecContainer[0].Command))
			Expect(cpupod.Spec.Containers[0].SecurityContext).To(Equal(expectedPodSpecContainer[0].SecurityContext))
			Expect(cpupod.Spec.Containers[0].WorkingDir).To(Equal(expectedPodSpecContainer[0].WorkingDir))
			Expect(cpupod.Spec.Containers[0].VolumeMounts).To(Equal(expectedPodSpecContainer[0].VolumeMounts))
			Expect(cpupod.Spec.Containers[0].Resources).To(Equal(expectedPodSpecContainer[0].Resources))

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
					Name: "dpdk",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/var/run/dpdk",
							Type: &volumeType,
						},
					},
				},
			}
			Expect(cpupod.Spec.Volumes).To(Equal(expectedPodSpecVolumes))

			var expectedPodSpecrResartPolicy corev1.RestartPolicy = "Always"
			expectedPodSpecHostNW := false
			expectedPodSpecHostIPC := true
			Expect(cpupod.Spec.RestartPolicy).To(Equal(expectedPodSpecrResartPolicy))
			Expect(cpupod.Spec.HostNetwork).To(Equal(expectedPodSpecHostNW))
			Expect(cpupod.Spec.HostIPC).To(Equal(expectedPodSpecHostIPC))

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
											"node01",
										},
									},
								},
							},
						},
					},
				},
			}
			Expect(cpupod.Spec.Affinity).To(Equal(expectedPodSpecAffinity))
		})
		It("Test_1-1-2_decode-tcp-CPU", func() {
			// Create CPUFuncConfig
			err = createConfig(ctx, cpuconfigdecode)
			if err != nil {
				fmt.Println("There is a problem in createing configmap ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection2)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create NetworkAttachmentDifinition
			err = createNetworkAttachmentDefinition(ctx, NetworkAttachmentDefinition1)
			if err != nil {
				fmt.Println("There is a problem in createing NetworkAttachmentDefinition ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR for cpu-filter-resize
			err = createCPUFunction(ctx, CPUFunction3)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFunc-filter-resizeCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())
			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction2)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFunc-decodeCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "cpufunctest2-wbfunction-decode-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var cpuCR examplecomv1.CPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctest2-wbfunction-decode-main",
				Namespace: TESTNAMESPACE,
			},
				&cpuCR)
			Expect(err).NotTo(HaveOccurred())

			var addTest = true
			var functionIndex int32 = 1
			var rxTest = "RTP"
			var txTest = "TCP"
			var podName = "cpufunctest2-wbfunction-decode-main-cpu-pod"
			// test for updating CR.Status
			expectedStatus := examplecomv1.CPUFunctionStatus{
				StartTime: cpuCR.Status.StartTime,
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctest2",
					Namespace: "default",
				},
				FunctionName: "cpu-decode",
				ImageURI:     "localhost/host_decode:3.1.0",
				ConfigName:   "cpufunc-config-decode",
				Status:       "Running",
				SharedMemory: &examplecomv1.SharedMemorySpec{
					FilePrefix:      "test01-cpufunctest2-wbfunction-decode-main",
					CommandQueueID:  "test01-cpufunctest2-wbfunction-decode-main",
					SharedMemoryMiB: 256,
				},
				VirtualNetworkDeviceDriverType: "sriov",
				AdditionalNetwork:              &addTest,
				FunctionIndex:                  &functionIndex,
				RxProtocol:                     &rxTest,
				TxProtocol:                     &txTest,
				PodName:                        &podName,
			}

			// test for updating StartTime
			err = checkTime(testTime.Time, cpuCR.Status.StartTime.Time)
			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status except for Status.StartTime
			Expect(cpuCR.Status).To(Equal(expectedStatus))

			// test for creating finalizer
			expectedFinalizer := []string{
				"cpufunction.finalizers.example.com.v1",
			}
			Expect(cpuCR.Finalizers).To(Equal(expectedFinalizer))

			// check logs
			Expect(writer.String()).To(ContainSubstring("nextConnectionCRName : cpufunctest2-wbconnection-decode-main-filter-resize-high-infer-main"))
			Expect(writer.String()).To(ContainSubstring("rxProtocol : RTP"))
			Expect(writer.String()).To(ContainSubstring("txProtocol : TCP"))
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

			var cpupod corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctest2-wbfunction-decode-main-cpu-pod",
				Namespace: TESTNAMESPACE,
			},
				&cpupod)
			if err != nil {
				// error handling
				fmt.Println("Could not get Pod:", cpupod, err)
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(*cpuCR.Status.PodName).To(Equal(cpupod.Name))

			expectedPodMeta := metav1.ObjectMeta{
				Name:      "cpufunctest2-wbfunction-decode-main-cpu-pod",
				Namespace: "default",
				Annotations: map[string]string{
					"k8s.v1.cni.cncf.io/networks": "[ {\"name\": \"node01-config-net-sriov\",\"ips\": [\"192.174.90.111/24\"] } ]",
				},
			}

			Expect(cpupod.ObjectMeta.Name).To(Equal(expectedPodMeta.Name))
			Expect(cpupod.ObjectMeta.Namespace).To(Equal(expectedPodMeta.Namespace))
			Expect(cpupod.ObjectMeta.Annotations["k8s.v1.cni.cncf.io/networks"]).To(Equal(expectedPodMeta.Annotations["k8s.v1.cni.cncf.io/networks"]))

			expectedPodSpecContainerEnv := []corev1.EnvVar{
				{
					Name:  "DECENV_APPLOG_LEVEL",
					Value: "6",
				},
				{
					Name:  "DECENV_FRAME_WIDTH",
					Value: "3840",
				},
				{
					Name:  "DECENV_FRAME_HEIGHT",
					Value: "2160",
				},
				{
					Name:  "DECENV_VIDEO_CONNECT_LIMIT",
					Value: "0",
				},
				{
					Name:  "DECENV_VIDEOSRC_PROTOCOL",
					Value: "RTP",
				},
				{
					Name:  "DECENV_OUTDST_PROTOCOL",
					Value: "TCP",
				},
				{
					Name:  "DECENV_VIDEOSRC_PORT",
					Value: "5004",
				},
				{
					Name:  "DECENV_FRAME_FPS",
					Value: "15",
				},
				{
					Name:  "DECENV_OUTDST_IPA",
					Value: "192.168.90.131",
				},
				{
					Name:  "DECENV_OUTDST_PORT",
					Value: "15000",
				},
				{
					Name:      "DECENV_VIDEOSRC_IPA",
					Value:     "192.174.90.111",
					ValueFrom: nil,
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
			var cpupodEnv []corev1.EnvVar
			for _, containers := range cpupod.Spec.Containers {
				cpupodEnv = append(cpupodEnv, containers.Env...)
			}

			sort.Slice(cpupodEnv, func(i, j int) bool {
				return cpupodEnv[i].Name < cpupodEnv[j].Name
			})

			Expect(cpupodEnv).To(Equal(expectedPodSpecContainerEnv))

			var cpupodLimits = corev1.ResourceList{
				"intel.com/intel_sriov_netdevice": {},
			}
			cpupodLimits["intel.com/intel_sriov_netdevice"] = resource.Quantity(resource.MustParse("1"))

			var cpupodRequests = corev1.ResourceList{
				"intel.com/intel_sriov_netdevice": {},
			}
			cpupodRequests["intel.com/intel_sriov_netdevice"] = resource.Quantity(resource.MustParse("1"))

			var privileged bool = true
			expectedPodSpecContainer := []corev1.Container{
				{
					Name:  "cpu-container0",
					Image: "localhost/host_decode:3.1.0",
					Command: []string{
						"sh",
						"-c",
					},
					ImagePullPolicy: corev1.PullIfNotPresent,
					Args: []string{
						"./tools/host_decode/build/host_decode-shared ",
					},
					WorkingDir:   "",
					VolumeMounts: nil,
					Resources: corev1.ResourceRequirements{
						Limits:   cpupodLimits,
						Requests: cpupodRequests,
					},
					SecurityContext: &corev1.SecurityContext{
						Privileged: &privileged,
					},
				},
			}

			Expect(cpupod.Spec.Containers[0].Name).To(Equal(expectedPodSpecContainer[0].Name))
			Expect(cpupod.Spec.Containers[0].Image).To(Equal(expectedPodSpecContainer[0].Image))
			Expect(cpupod.Spec.Containers[0].ImagePullPolicy).To(Equal(expectedPodSpecContainer[0].ImagePullPolicy))
			Expect(cpupod.Spec.Containers[0].Args).To(Equal(expectedPodSpecContainer[0].Args))
			Expect(cpupod.Spec.Containers[0].Command).To(Equal(expectedPodSpecContainer[0].Command))
			Expect(cpupod.Spec.Containers[0].SecurityContext).To(Equal(expectedPodSpecContainer[0].SecurityContext))
			Expect(cpupod.Spec.Containers[0].WorkingDir).To(Equal(expectedPodSpecContainer[0].WorkingDir))
			Expect(cpupod.Spec.Containers[0].VolumeMounts).To(Equal(expectedPodSpecContainer[0].VolumeMounts))
			Expect(cpupod.Spec.Containers[0].Resources).To(Equal(expectedPodSpecContainer[0].Resources))

			var expectedPodSpecVolumes []corev1.Volume = nil
			Expect(cpupod.Spec.Volumes).To(Equal(expectedPodSpecVolumes))

			var expectedPodSpecrResartPolicy corev1.RestartPolicy = "Always"
			expectedPodSpecHostNW := false
			expectedPodSpecHostIPC := true
			Expect(cpupod.Spec.RestartPolicy).To(Equal(expectedPodSpecrResartPolicy))
			Expect(cpupod.Spec.HostNetwork).To(Equal(expectedPodSpecHostNW))
			Expect(cpupod.Spec.HostIPC).To(Equal(expectedPodSpecHostIPC))

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
											"node01",
										},
									},
								},
							},
						},
					},
				},
			}
			Expect(cpupod.Spec.Affinity).To(Equal(expectedPodSpecAffinity))
		})
		It("Test_1-1-3_filter-resize-high-infer", func() {
			// Create CPUFuncConfig
			err = createConfig(ctx, cpuconfigfrhigh)
			if err != nil {
				fmt.Println("There is a problem in createing configmap ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnectionfrhigh)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection3frhigh)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create NetworkAttachmentDifinition
			err = createNetworkAttachmentDefinition(ctx, NetworkAttachmentDefinition1)
			if err != nil {
				fmt.Println("There is a problem in createing NetworkAttachmentDefinition ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction3)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFunc-filter-resizeCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "cpufunctest3-wbfunction-filter-resize-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var cpuCR examplecomv1.CPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctest3-wbfunction-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&cpuCR)
			Expect(err).NotTo(HaveOccurred())

			var addTest = true
			var functionIndex int32 = 2
			var rxTest = "TCP"
			var txTest = "TCP"
			var podName = "cpufunctest3-wbfunction-filter-resize-high-infer-main-cpu-pod"
			// test for updating CR.Status
			expectedStatus := examplecomv1.CPUFunctionStatus{
				StartTime: cpuCR.Status.StartTime,
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctest3",
					Namespace: "default",
				},
				FunctionName: "cpu-filter-resize-high-infer",
				ImageURI:     "localhost/cpu-filterresize-app:3.1.0",
				ConfigName:   "cpufunc-config-filter-resize-high-infer",
				Status:       "Running",
				SharedMemory: &examplecomv1.SharedMemorySpec{
					FilePrefix:      "test01-cpufunctest3-wbfunction-filter-resize-high-infer-main",
					CommandQueueID:  "test01-cpufunctest3-wbfunction-filter-resize-high-infer-main",
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
			err = checkTime(testTime.Time, cpuCR.Status.StartTime.Time)
			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status except for Status.StartTime
			Expect(cpuCR.Status).To(Equal(expectedStatus))

			// test for creating finalizer
			expectedFinalizer := []string{
				"cpufunction.finalizers.example.com.v1",
			}
			Expect(cpuCR.Finalizers).To(Equal(expectedFinalizer))

			// check logs
			Expect(writer.String()).To(ContainSubstring("previousConnectionCRName : cpufunctest3-wbconnection-decode-main-filter-resize-high-infer-main"))
			Expect(writer.String()).To(ContainSubstring("nextConnectionCRName : cpufunctest3-wbconnection-filter-resize-high-infer-main-high-infer-main"))
			Expect(writer.String()).To(ContainSubstring("rxProtocol : TCP"))
			Expect(writer.String()).To(ContainSubstring("txProtocol : TCP"))
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

			var cpupod corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctest3-wbfunction-filter-resize-high-infer-main-cpu-pod",
				Namespace: TESTNAMESPACE,
			},
				&cpupod)

			Expect(err).NotTo(HaveOccurred())
			Expect(*cpuCR.Status.PodName).To(Equal(cpupod.Name))

			expectedPodMeta := metav1.ObjectMeta{
				Name:      "cpufunctest3-wbfunction-filter-resize-high-infer-main-cpu-pod",
				Namespace: "default",
				Annotations: map[string]string{
					"k8s.v1.cni.cncf.io/networks": "[ {\"name\": \"node01-config-net-sriov\",\"ips\": [\"192.168.122.50/24\"] } ]",
				},
			}

			Expect(cpupod.ObjectMeta.Name).To(Equal(expectedPodMeta.Name))
			Expect(cpupod.ObjectMeta.Namespace).To(Equal(expectedPodMeta.Namespace))
			Expect(cpupod.ObjectMeta.Annotations["k8s.v1.cni.cncf.io/networks"]).To(Equal(expectedPodMeta.Annotations["k8s.v1.cni.cncf.io/networks"]))

			expectedPodSpecContainerEnv := []corev1.EnvVar{
				{
					Name:  "FRENV_APPLOG_LEVEL",
					Value: "INFO",
				},
				{
					Name:  "FRENV_INPUT_HEIGHT",
					Value: "2160",
				},
				{
					Name:  "FRENV_INPUT_PORT",
					Value: "15000",
				},
				{
					Name:  "FRENV_INPUT_WIDTH",
					Value: "3840",
				},
				{
					Name:  "FRENV_OUTPUT_HEIGHT",
					Value: "1280",
				},
				{
					Name:  "FRENV_OUTPUT_IP",
					Value: "192.168.122.100",
				},
				{
					Name:  "FRENV_OUTPUT_PORT",
					Value: "16000",
				},
				{
					Name:  "FRENV_OUTPUT_WIDTH",
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
			var cpupodEnv []corev1.EnvVar
			for _, containers := range cpupod.Spec.Containers {
				cpupodEnv = append(cpupodEnv, containers.Env...)
			}

			sort.Slice(cpupodEnv, func(i, j int) bool {
				return cpupodEnv[i].Name < cpupodEnv[j].Name
			})

			Expect(cpupodEnv).To(Equal(expectedPodSpecContainerEnv))

			var cpupodLimits = corev1.ResourceList{
				"intel.com/intel_sriov_netdevice": {},
			}
			cpupodLimits["intel.com/intel_sriov_netdevice"] = resource.Quantity(resource.MustParse("1"))

			var cpupodRequests = corev1.ResourceList{
				"intel.com/intel_sriov_netdevice": {},
			}
			cpupodRequests["intel.com/intel_sriov_netdevice"] = resource.Quantity(resource.MustParse("1"))

			var privileged bool = true
			expectedPodSpecContainer := []corev1.Container{
				{
					Name:  "cpu-container0",
					Image: "localhost/cpu-filterresize-app:3.1.0",
					Command: []string{
						"python",
						"fr.py",
						"--in_port=$(FRENV_INPUT_PORT)",
						"--out_addr=$(FRENV_OUTPUT_IP)",
						"--out_port=$(FRENV_OUTPUT_PORT)",
						"--in_width=$(FRENV_INPUT_WIDTH)",
						"--in_height=$(FRENV_INPUT_HEIGHT)",
						"--out_width=$(FRENV_OUTPUT_WIDTH)",
						"--out_height=$(FRENV_OUTPUT_HEIGHT)",
						"--loglevel=$(FRENV_APPLOG_LEVEL)",
					},
					ImagePullPolicy: corev1.PullIfNotPresent,
					WorkingDir:      "",
					VolumeMounts:    nil,
					Resources: corev1.ResourceRequirements{
						Limits:   cpupodLimits,
						Requests: cpupodRequests,
					},
					SecurityContext: &corev1.SecurityContext{
						Privileged: &privileged,
					},
				},
			}

			Expect(cpupod.Spec.Containers[0].Name).To(Equal(expectedPodSpecContainer[0].Name))
			Expect(cpupod.Spec.Containers[0].Image).To(Equal(expectedPodSpecContainer[0].Image))
			Expect(cpupod.Spec.Containers[0].ImagePullPolicy).To(Equal(expectedPodSpecContainer[0].ImagePullPolicy))
			Expect(cpupod.Spec.Containers[0].Args).To(Equal(expectedPodSpecContainer[0].Args))
			Expect(cpupod.Spec.Containers[0].Command).To(Equal(expectedPodSpecContainer[0].Command))
			Expect(cpupod.Spec.Containers[0].SecurityContext).To(Equal(expectedPodSpecContainer[0].SecurityContext))
			Expect(cpupod.Spec.Containers[0].WorkingDir).To(Equal(expectedPodSpecContainer[0].WorkingDir))
			Expect(cpupod.Spec.Containers[0].VolumeMounts).To(Equal(expectedPodSpecContainer[0].VolumeMounts))
			Expect(cpupod.Spec.Containers[0].Resources).To(Equal(expectedPodSpecContainer[0].Resources))

			var expectedPodSpecVolumes []corev1.Volume = nil
			Expect(cpupod.Spec.Volumes).To(Equal(expectedPodSpecVolumes))

			var expectedPodSpecrResartPolicy corev1.RestartPolicy = "Always"
			expectedPodSpecHostNW := false
			expectedPodSpecHostIPC := true
			Expect(cpupod.Spec.RestartPolicy).To(Equal(expectedPodSpecrResartPolicy))
			Expect(cpupod.Spec.HostNetwork).To(Equal(expectedPodSpecHostNW))
			Expect(cpupod.Spec.HostIPC).To(Equal(expectedPodSpecHostIPC))

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
											"node01",
										},
									},
								},
							},
						},
					},
				},
			}
			Expect(cpupod.Spec.Affinity).To(Equal(expectedPodSpecAffinity))

		})
		It("Test_1-1-4_filter-resize-low-infer", func() {
			// Create CPUFuncConfig
			err = createConfig(ctx, cpuconfigfrlow)
			if err != nil {
				fmt.Println("There is a problem in createing configmap ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnectionfrlow)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection4frlow)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create NetworkAttachmentDifinition
			err = createNetworkAttachmentDefinition(ctx, NetworkAttachmentDefinition1)
			if err != nil {
				fmt.Println("There is a problem in createing NetworkAttachmentDefinition ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction4frlow)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFunc-filter-resizeCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "cpufunctest4-wbfunction-filter-resize-low-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var cpuCR examplecomv1.CPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctest4-wbfunction-filter-resize-low-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&cpuCR)
			Expect(err).NotTo(HaveOccurred())

			var addTest = true
			var functionIndex int32 = 3
			var rxTest = "TCP"
			var txTest = "TCP"
			var podName = "cpufunctest4-wbfunction-filter-resize-low-infer-main-cpu-pod"
			// test for updating CR.Status
			expectedStatus := examplecomv1.CPUFunctionStatus{
				StartTime: cpuCR.Status.StartTime,
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctest4",
					Namespace: "default",
				},
				FunctionName: "cpu-filter-resize-low-infer",
				ImageURI:     "localhost/cpu-filterresize-app:3.1.0",
				ConfigName:   "cpufunc-config-filter-resize-low-infer",
				Status:       "Running",
				SharedMemory: &examplecomv1.SharedMemorySpec{
					FilePrefix:      "test01-cpufunctest4-wbfunction-filter-resize-low-infer-main",
					CommandQueueID:  "test01-cpufunctest4-wbfunction-filter-resize-low-infer-main",
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
			err = checkTime(testTime.Time, cpuCR.Status.StartTime.Time)
			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status except for Status.StartTime
			Expect(cpuCR.Status).To(Equal(expectedStatus))

			// test for creating finalizer
			expectedFinalizer := []string{
				"cpufunction.finalizers.example.com.v1",
			}
			Expect(cpuCR.Finalizers).To(Equal(expectedFinalizer))

			// check logs
			Expect(writer.String()).To(ContainSubstring("previousConnectionCRName : cpufunctest4-wbconnection-decode-main-filter-resize-low-infer-main"))
			Expect(writer.String()).To(ContainSubstring("nextConnectionCRName : cpufunctest4-wbconnection-filter-resize-low-infer-main-low-infer-main"))
			Expect(writer.String()).To(ContainSubstring("rxProtocol : TCP"))
			Expect(writer.String()).To(ContainSubstring("txProtocol : TCP"))
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

			var cpupod corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctest4-wbfunction-filter-resize-low-infer-main-cpu-pod",
				Namespace: TESTNAMESPACE,
			},
				&cpupod)

			Expect(err).NotTo(HaveOccurred())

			expectedPodMeta := metav1.ObjectMeta{
				Name:      "cpufunctest4-wbfunction-filter-resize-low-infer-main-cpu-pod",
				Namespace: "default",
				Annotations: map[string]string{
					"k8s.v1.cni.cncf.io/networks": "[ {\"name\": \"node01-config-net-sriov\",\"ips\": [\"192.168.122.150/24\"] } ]",
				},
			}

			Expect(cpupod.ObjectMeta.Name).To(Equal(expectedPodMeta.Name))
			Expect(cpupod.ObjectMeta.Namespace).To(Equal(expectedPodMeta.Namespace))
			Expect(cpupod.ObjectMeta.Annotations["k8s.v1.cni.cncf.io/networks"]).To(Equal(expectedPodMeta.Annotations["k8s.v1.cni.cncf.io/networks"]))

			expectedPodSpecContainerEnv := []corev1.EnvVar{
				{
					Name:  "FRENV_APPLOG_LEVEL",
					Value: "INFO",
				},
				{
					Name:  "FRENV_INPUT_HEIGHT",
					Value: "2160",
				},
				{
					Name:  "FRENV_INPUT_PORT",
					Value: "15000",
				},
				{
					Name:  "FRENV_INPUT_WIDTH",
					Value: "3840",
				},
				{
					Name:  "FRENV_OUTPUT_HEIGHT",
					Value: "416",
				},
				{
					Name:  "FRENV_OUTPUT_IP",
					Value: "192.168.122.121",
				},
				{
					Name:  "FRENV_OUTPUT_PORT",
					Value: "16000",
				},
				{
					Name:  "FRENV_OUTPUT_WIDTH",
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
				{Name: "test", Value: "testvalue", ValueFrom: nil},
			}

			sort.Slice(expectedPodSpecContainerEnv, func(i, j int) bool {
				return expectedPodSpecContainerEnv[i].Name < expectedPodSpecContainerEnv[j].Name
			})
			var cpupodEnv []corev1.EnvVar
			for _, containers := range cpupod.Spec.Containers {
				cpupodEnv = append(cpupodEnv, containers.Env...)
			}

			sort.Slice(cpupodEnv, func(i, j int) bool {
				return cpupodEnv[i].Name < cpupodEnv[j].Name
			})

			Expect(cpupodEnv).To(Equal(expectedPodSpecContainerEnv))

			var cpupodLimits = corev1.ResourceList{
				"intel.com/intel_sriov_netdevice": {},
			}
			cpupodLimits["intel.com/intel_sriov_netdevice"] = resource.Quantity(resource.MustParse("1"))

			var cpupodRequests = corev1.ResourceList{
				"intel.com/intel_sriov_netdevice": {},
			}
			cpupodRequests["intel.com/intel_sriov_netdevice"] = resource.Quantity(resource.MustParse("1"))

			var privileged bool = true
			expectedPodSpecContainer := []corev1.Container{
				{
					Name:  "cpu-container0",
					Image: "localhost/cpu-filterresize-app:3.1.0",
					Command: []string{
						"python",
						"fr.py",
						"--in_port=$(FRENV_INPUT_PORT)",
						"--out_addr=$(FRENV_OUTPUT_IP)",
						"--out_port=$(FRENV_OUTPUT_PORT)",
						"--in_width=$(FRENV_INPUT_WIDTH)",
						"--in_height=$(FRENV_INPUT_HEIGHT)",
						"--out_width=$(FRENV_OUTPUT_WIDTH)",
						"--out_height=$(FRENV_OUTPUT_HEIGHT)",
						"--loglevel=$(FRENV_APPLOG_LEVEL)",
					},
					ImagePullPolicy: corev1.PullIfNotPresent,
					WorkingDir:      "",
					VolumeMounts:    nil,
					Resources: corev1.ResourceRequirements{
						Limits:   cpupodLimits,
						Requests: cpupodRequests,
					},
					SecurityContext: &corev1.SecurityContext{
						Privileged: &privileged,
					},
				},
			}

			Expect(cpupod.Spec.Containers[0].Name).To(Equal(expectedPodSpecContainer[0].Name))
			Expect(cpupod.Spec.Containers[0].Image).To(Equal(expectedPodSpecContainer[0].Image))
			Expect(cpupod.Spec.Containers[0].ImagePullPolicy).To(Equal(expectedPodSpecContainer[0].ImagePullPolicy))
			Expect(cpupod.Spec.Containers[0].Args).To(Equal(expectedPodSpecContainer[0].Args))
			Expect(cpupod.Spec.Containers[0].Command).To(Equal(expectedPodSpecContainer[0].Command))
			Expect(cpupod.Spec.Containers[0].SecurityContext).To(Equal(expectedPodSpecContainer[0].SecurityContext))
			Expect(cpupod.Spec.Containers[0].WorkingDir).To(Equal(expectedPodSpecContainer[0].WorkingDir))
			Expect(cpupod.Spec.Containers[0].VolumeMounts).To(Equal(expectedPodSpecContainer[0].VolumeMounts))
			Expect(cpupod.Spec.Containers[0].Resources).To(Equal(expectedPodSpecContainer[0].Resources))

			var expectedPodSpecVolumes []corev1.Volume = nil
			Expect(cpupod.Spec.Volumes).To(Equal(expectedPodSpecVolumes))

			var expectedPodSpecrResartPolicy corev1.RestartPolicy = "Always"
			expectedPodSpecHostNW := false
			expectedPodSpecHostIPC := true
			Expect(cpupod.Spec.RestartPolicy).To(Equal(expectedPodSpecrResartPolicy))
			Expect(cpupod.Spec.HostNetwork).To(Equal(expectedPodSpecHostNW))
			Expect(cpupod.Spec.HostIPC).To(Equal(expectedPodSpecHostIPC))

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
											"node01",
										},
									},
								},
							},
						},
					},
				},
			}
			Expect(cpupod.Spec.Affinity).To(Equal(expectedPodSpecAffinity))

		})
		It("Test_1-1-5_copy-branch", func() {
			// Create CPUFuncConfig
			err = createConfig(ctx, cpuconfigcopybranch)
			if err != nil {
				fmt.Println("There is a problem in createing configmap ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection4)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection5)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection6)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create NetworkAttachmentDifinition
			err = createNetworkAttachmentDefinition(ctx, NetworkAttachmentDefinition1)
			if err != nil {
				fmt.Println("There is a problem in createing NetworkAttachmentDefinition ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction5copy)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR-copy-branch", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "cpufunctest5-wbfunction-copy-branch-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var cpuCR examplecomv1.CPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctest5-wbfunction-copy-branch-main",
				Namespace: TESTNAMESPACE,
			},
				&cpuCR)
			Expect(err).NotTo(HaveOccurred())

			var addTest = true
			var functionIndex int32 = 4
			var rxTest = "TCP"
			var txTest = "TCP"
			var podName = "cpufunctest5-wbfunction-copy-branch-main-cpu-pod"
			// test for updating CR.Status
			expectedStatus := examplecomv1.CPUFunctionStatus{
				StartTime: cpuCR.Status.StartTime,
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctest5",
					Namespace: "default",
				},
				FunctionName: "copy-branch",
				ImageURI:     "localhost/cpu-copybranch-app:3.1.0",
				ConfigName:   "cpufunc-config-copy-branch",
				Status:       "Running",
				SharedMemory: &examplecomv1.SharedMemorySpec{
					FilePrefix:      "test01-cpufunctest5-wbfunction-copy-branch-main",
					CommandQueueID:  "test01-cpufunctest5-wbfunction-copy-branch-main",
					SharedMemoryMiB: 0,
				},
				VirtualNetworkDeviceDriverType: "sriov",
				AdditionalNetwork:              &addTest,
				FunctionIndex:                  &functionIndex,
				RxProtocol:                     &rxTest,
				TxProtocol:                     &txTest,
				PodName:                        &podName,
			}

			// test for updating StartTime
			err = checkTime(testTime.Time, cpuCR.Status.StartTime.Time)
			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status except for Status.StartTime
			Expect(cpuCR.Status).To(Equal(expectedStatus))

			// test for creating finalizer
			expectedFinalizer := []string{
				"cpufunction.finalizers.example.com.v1",
			}
			Expect(cpuCR.Finalizers).To(Equal(expectedFinalizer))

			// check logs
			Expect(writer.String()).To(ContainSubstring("previousConnectionCRName : cpufunctest5-wbconnection-filter-resize-low-infer-main-copy-branch-main"))
			Expect(writer.String()).To(ContainSubstring("nextConnectionCRName : cpufunctest5-wbconnection-copy-branch-main-infer-"))
			Expect(writer.String()).To(ContainSubstring("rxProtocol : TCP"))
			Expect(writer.String()).To(ContainSubstring("txProtocol : TCP"))
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

			var cpupod corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctest5-wbfunction-copy-branch-main-cpu-pod",
				Namespace: TESTNAMESPACE,
			},
				&cpupod)

			Expect(err).NotTo(HaveOccurred())
			Expect(*cpuCR.Status.PodName).To(Equal(cpupod.Name))

			expectedPodMeta := metav1.ObjectMeta{
				Name:      "cpufunctest5-wbfunction-copy-branch-main-cpu-pod",
				Namespace: "default",
				Annotations: map[string]string{
					"k8s.v1.cni.cncf.io/networks": "[ {\"name\": \"node01-config-net-sriov\",\"ips\": [\"192.168.122.121/24\"] } ]",
				},
			}

			Expect(cpupod.ObjectMeta.Name).To(Equal(expectedPodMeta.Name))
			Expect(cpupod.ObjectMeta.Namespace).To(Equal(expectedPodMeta.Namespace))
			Expect(cpupod.ObjectMeta.Annotations["k8s.v1.cni.cncf.io/networks"]).To(Equal(expectedPodMeta.Annotations["k8s.v1.cni.cncf.io/networks"]))

			var cpupodLimits = corev1.ResourceList{
				"intel.com/intel_sriov_netdevice": {},
			}
			cpupodLimits["intel.com/intel_sriov_netdevice"] = resource.Quantity(resource.MustParse("1"))

			var cpupodRequests = corev1.ResourceList{
				"intel.com/intel_sriov_netdevice": {},
			}
			cpupodRequests["intel.com/intel_sriov_netdevice"] = resource.Quantity(resource.MustParse("1"))

			var privileged bool = true
			expectedPodSpecContainer := []corev1.Container{
				{
					Args: []string{
						"./copy_branch 192.168.122.121:16000 2 192.168.90.141:17000,192.168.90.142:18000 1024 ",
					},
					Name:  "cpu-container0",
					Image: "localhost/cpu-copybranch-app:3.1.0",
					Command: []string{
						"sh",
						"-c",
					},
					ImagePullPolicy: corev1.PullIfNotPresent,
					WorkingDir:      "/opt/fpga-software/tools/copy_branch",
					VolumeMounts:    nil,
					Resources: corev1.ResourceRequirements{
						Limits:   cpupodLimits,
						Requests: cpupodRequests,
					},
					SecurityContext: &corev1.SecurityContext{
						Privileged: &privileged,
					},
				},
			}

			Expect(cpupod.Spec.Containers[0].Name).To(Equal(expectedPodSpecContainer[0].Name))
			Expect(cpupod.Spec.Containers[0].Image).To(Equal(expectedPodSpecContainer[0].Image))
			Expect(cpupod.Spec.Containers[0].ImagePullPolicy).To(Equal(expectedPodSpecContainer[0].ImagePullPolicy))
			Expect(cpupod.Spec.Containers[0].Args).To(Equal(expectedPodSpecContainer[0].Args))
			Expect(cpupod.Spec.Containers[0].Command).To(Equal(expectedPodSpecContainer[0].Command))
			Expect(cpupod.Spec.Containers[0].SecurityContext).To(Equal(expectedPodSpecContainer[0].SecurityContext))
			Expect(cpupod.Spec.Containers[0].WorkingDir).To(Equal(expectedPodSpecContainer[0].WorkingDir))
			Expect(cpupod.Spec.Containers[0].VolumeMounts).To(Equal(expectedPodSpecContainer[0].VolumeMounts))
			Expect(cpupod.Spec.Containers[0].Resources).To(Equal(expectedPodSpecContainer[0].Resources))

			var expectedPodSpecVolumes []corev1.Volume = nil
			Expect(cpupod.Spec.Volumes).To(Equal(expectedPodSpecVolumes))

			var expectedPodSpecrResartPolicy corev1.RestartPolicy = "Always"
			expectedPodSpecHostNW := false
			expectedPodSpecHostIPC := true
			Expect(cpupod.Spec.RestartPolicy).To(Equal(expectedPodSpecrResartPolicy))
			Expect(cpupod.Spec.HostNetwork).To(Equal(expectedPodSpecHostNW))
			Expect(cpupod.Spec.HostIPC).To(Equal(expectedPodSpecHostIPC))

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
											"node01",
										},
									},
								},
							},
						},
					},
				},
			}
			Expect(cpupod.Spec.Affinity).To(Equal(expectedPodSpecAffinity))

		})
		It("Test_1-1-6_glue-dma-to-tcp", func() {
			// Create CPUFuncConfig
			err = createConfig(ctx, cpuconfigglue)
			if err != nil {
				fmt.Println("There is a problem in createing configmap ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection5)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGAFunctionCR
			err = createFPGAFunction(ctx, FPGAFunction5)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction6glue)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection6glue)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create NetworkAttachmentDifinition
			err = createNetworkAttachmentDefinition(ctx, NetworkAttachmentDefinition1)
			if err != nil {
				fmt.Println("There is a problem in createing NetworkAttachmentDefinition ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "cpufunctest6-wbfunction-glue-fdma-to-tcp-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var cpuCR examplecomv1.CPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctest6-wbfunction-glue-fdma-to-tcp-main",
				Namespace: TESTNAMESPACE,
			},
				&cpuCR)
			Expect(err).NotTo(HaveOccurred())

			var addTest = true
			var functionIndex int32 = 5
			var rxTest = "DMA"
			var txTest = "TCP"
			var podName = "cpufunctest6-wbfunction-glue-fdma-to-tcp-main-cpu-pod"
			// test for updating CR.Status
			expectedStatus := examplecomv1.CPUFunctionStatus{
				StartTime: cpuCR.Status.StartTime,
				DataFlowRef: examplecomv1.WBNamespacedName{
					Name:      "cpufunctest6",
					Namespace: "default",
				},
				FunctionName: "glue-fdma-to-tcp",
				ImageURI:     "localhost/cpu-glue-app:3.1.0",
				ConfigName:   "cpufunc-config-glue-fdma-to-tcp",
				Status:       "Running",
				SharedMemory: &examplecomv1.SharedMemorySpec{
					FilePrefix:      "test01-cpufunctest6-wbfunction-glue-fdma-to-tcp-main",
					CommandQueueID:  "test01-cpufunctest6-wbfunction-glue-fdma-to-tcp-main",
					SharedMemoryMiB: 256,
				},
				VirtualNetworkDeviceDriverType: "sriov",
				AdditionalNetwork:              &addTest,
				FunctionIndex:                  &functionIndex,
				RxProtocol:                     &rxTest,
				TxProtocol:                     &txTest,
				PodName:                        &podName,
			}

			// test for updating StartTime
			err = checkTime(testTime.Time, cpuCR.Status.StartTime.Time)
			Expect(err).NotTo(HaveOccurred())

			// test for updating CR.Status except for Status.StartTime
			Expect(cpuCR.Status).To(Equal(expectedStatus))

			// test for creating finalizer
			expectedFinalizer := []string{
				"cpufunction.finalizers.example.com.v1",
			}
			Expect(cpuCR.Finalizers).To(Equal(expectedFinalizer))

			// check logs
			Expect(writer.String()).To(ContainSubstring("previousConnectionCRName : cpufunctest6-wbconnection-filter-resize-high-infer-main-glue-fdma-to-tcp-main"))
			Expect(writer.String()).To(ContainSubstring("nextConnectionCRName : cpufunctest6-wbconnection-glue-fdma-to-tcp-main-high-infer-main"))
			Expect(writer.String()).To(ContainSubstring("rxProtocol : DMA"))
			Expect(writer.String()).To(ContainSubstring("txProtocol : TCP"))
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

			var cpupod corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctest6-wbfunction-glue-fdma-to-tcp-main-cpu-pod",
				Namespace: TESTNAMESPACE,
			},
				&cpupod)

			Expect(err).NotTo(HaveOccurred())
			Expect(*cpuCR.Status.PodName).To(Equal(cpupod.Name))

			expectedPodMeta := metav1.ObjectMeta{
				Name:      "cpufunctest6-wbfunction-glue-fdma-to-tcp-main-cpu-pod",
				Namespace: "default",
				Annotations: map[string]string{
					"k8s.v1.cni.cncf.io/networks": "[ {\"name\": \"node01-config-net-sriov\",\"ips\": [\"192.174.122.131/24\"] } ]",
				},
			}

			Expect(cpupod.ObjectMeta.Name).To(Equal(expectedPodMeta.Name))
			Expect(cpupod.ObjectMeta.Namespace).To(Equal(expectedPodMeta.Namespace))
			Expect(cpupod.ObjectMeta.Annotations["k8s.v1.cni.cncf.io/networks"]).To(Equal(expectedPodMeta.Annotations["k8s.v1.cni.cncf.io/networks"]))

			expectedPodSpecContainerEnv := []corev1.EnvVar{
				{
					Name:  "GLUEENV_DPDK_FILE_PREFIX",
					Value: "test01-cpufunctest6-wbfunction-glue-fdma-to-tcp-main",
				},
				{
					Name:  "GLUEENV_FPGA_DEV_NAME",
					Value: "/dev/xpcie_21330621T00D",
				},
				{
					Name:  "GLUEENV_FPGA_DMA_DEV_TO_HOST_CONNECTOR_ID",
					Value: "test01-cpufunctest6-wbfunction-glue-fdma-to-tcp-main",
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
			var cpupodEnv []corev1.EnvVar
			for _, containers := range cpupod.Spec.Containers {
				cpupodEnv = append(cpupodEnv, containers.Env...)
			}

			sort.Slice(cpupodEnv, func(i, j int) bool {
				return cpupodEnv[i].Name < cpupodEnv[j].Name
			})

			Expect(cpupodEnv).To(Equal(expectedPodSpecContainerEnv))

			var cpupodLimits = corev1.ResourceList{
				"hugepages-1Gi":                   resource.Quantity{},
				"intel.com/intel_sriov_netdevice": {},
			}
			cpupodLimits["hugepages-1Gi"] = resource.Quantity(resource.MustParse("1Gi"))
			cpupodLimits["intel.com/intel_sriov_netdevice"] = resource.Quantity(resource.MustParse("1"))

			var cpupodRequests = corev1.ResourceList{
				"hugepages-1Gi":                   resource.Quantity{},
				"intel.com/intel_sriov_netdevice": {},
				"memory":                          {},
			}
			cpupodRequests["hugepages-1Gi"] = resource.Quantity(resource.MustParse("1Gi"))
			cpupodRequests["intel.com/intel_sriov_netdevice"] = resource.Quantity(resource.MustParse("1"))
			cpupodRequests["memory"] = resource.Quantity(resource.MustParse("32Gi"))

			var privileged bool = true
			expectedPodSpecContainer := []corev1.Container{
				{
					Args: []string{
						"./build/glue 192.174.90.141:16000 1280 1280 ",
					},
					Name:  "cpu-container0",
					Image: "localhost/cpu-glue-app:3.1.0",
					Command: []string{
						"sh",
						"-c",
					},
					ImagePullPolicy: corev1.PullIfNotPresent,
					WorkingDir:      "/opt/fpga-software/tools/glue_fdma_tcp",
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "hugepage-1gi",
							MountPath: "/dev/hugepages",
						},
						{
							Name:      "dpdk",
							MountPath: "/var/run/dpdk",
						},
					},
					Resources: corev1.ResourceRequirements{
						Limits:   cpupodLimits,
						Requests: cpupodRequests,
					},
					SecurityContext: &corev1.SecurityContext{
						Privileged: &privileged,
					},
				},
			}

			Expect(cpupod.Spec.Containers[0].Name).To(Equal(expectedPodSpecContainer[0].Name))
			Expect(cpupod.Spec.Containers[0].Image).To(Equal(expectedPodSpecContainer[0].Image))
			Expect(cpupod.Spec.Containers[0].ImagePullPolicy).To(Equal(expectedPodSpecContainer[0].ImagePullPolicy))
			Expect(cpupod.Spec.Containers[0].Args).To(Equal(expectedPodSpecContainer[0].Args))
			Expect(cpupod.Spec.Containers[0].Command).To(Equal(expectedPodSpecContainer[0].Command))
			Expect(cpupod.Spec.Containers[0].SecurityContext).To(Equal(expectedPodSpecContainer[0].SecurityContext))
			Expect(cpupod.Spec.Containers[0].WorkingDir).To(Equal(expectedPodSpecContainer[0].WorkingDir))
			Expect(cpupod.Spec.Containers[0].VolumeMounts).To(Equal(expectedPodSpecContainer[0].VolumeMounts))
			Expect(cpupod.Spec.Containers[0].Resources).To(Equal(expectedPodSpecContainer[0].Resources))

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
					Name: "dpdk",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/var/run/dpdk",
							Type: &volumeType,
						},
					},
				},
			}
			Expect(cpupod.Spec.Volumes).To(Equal(expectedPodSpecVolumes))

			var expectedPodSpecrResartPolicy corev1.RestartPolicy = "Always"
			expectedPodSpecHostNW := false
			expectedPodSpecHostIPC := true
			Expect(cpupod.Spec.RestartPolicy).To(Equal(expectedPodSpecrResartPolicy))
			Expect(cpupod.Spec.HostNetwork).To(Equal(expectedPodSpecHostNW))
			Expect(cpupod.Spec.HostIPC).To(Equal(expectedPodSpecHostIPC))

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
											"node01",
										},
									},
								},
							},
						},
					},
				},
			}
			Expect(cpupod.Spec.Affinity).To(Equal(expectedPodSpecAffinity))

		})
		It("Test_1-2-1_redeployment", func() {
			// Create CPUFuncConfig
			err = createConfig(ctx, cpuconfigfrhigh)
			if err != nil {
				fmt.Println("There is a problem in createing configmap ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection1)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection3)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction12)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFunc-filter-resizeCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create NetworkAttachmentDifinition
			err = createNetworkAttachmentDefinition(ctx, NetworkAttachmentDefinition1)
			if err != nil {
				fmt.Println("There is a problem in createing NetworkAttachmentDefinition ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "dftest-wbfunction-filter-resize-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var cpuCR examplecomv1.CPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "dftest-wbfunction-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&cpuCR)
			if err != nil {
				// Error route
				fmt.Println("Cannot get CPUFunctionCR:", cpuCR, err)
			}
			Expect(err).NotTo(HaveOccurred())

			Expect(*cpuCR.Status.FunctionIndex).To(Equal(int32(6)))

			// delete CPUFunctionCR
			err = deleteCPUFunction(ctx, CPUFunction12)
			if err != nil {
				fmt.Println("There is a problem in deleteing CPUFunc-filter-resizeCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "dftest-wbfunction-filter-resize-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// delete Pod
			var cpupod1 corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "dftest-wbfunction-filter-resize-high-infer-main-cpu-pod",
				Namespace: "default",
			},
				&cpupod1)
			if err != nil {
				// Error route
				fmt.Println("Cannot get Pod:", cpupod1, err)
			}
			err = reconciler.Delete(ctx, &cpupod1)
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction12)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFunc-filter-resizeCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "dftest-wbfunction-filter-resize-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var cpuCRredeploy examplecomv1.CPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "dftest-wbfunction-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&cpuCRredeploy)
			if err != nil {
				// Error route
				fmt.Println("Cannot get CPUFunctionCR:", cpuCRredeploy, err)
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(*cpuCR.Status.FunctionIndex).To(Equal(int32(6)))
			var cpupod corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "dftest-wbfunction-filter-resize-high-infer-main-cpu-pod",
				Namespace: TESTNAMESPACE,
			},
				&cpupod)
			if err != nil {
				// Error route
				fmt.Println("Cannot get Pod:", cpupod, err)
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(*cpuCR.Status.PodName).To(Equal(cpupod.Name))

			// delete CPUFunctionCR
			err = deleteCPUFunction(ctx, CPUFunction12)
			if err != nil {
				fmt.Println("There is a problem in deleteing CPUFunc-filter-resizeCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			_, _ = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "dftest-wbfunction-filter-resize-high-infer-main",
				}})
		})
		It("Test_1-2-2_next-deployment", func() {
			// Create CPUFuncConfig
			err = createConfig(ctx, cpuconfigfrhigh)
			if err != nil {
				fmt.Println("There is a problem in createing configmap ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection1)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection3)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction122)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFunc-filter-resizeCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create NetworkAttachmentDifinition
			err = createNetworkAttachmentDefinition(ctx, NetworkAttachmentDefinition1)
			if err != nil {
				fmt.Println("There is a problem in createing NetworkAttachmentDefinition ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "dftest-wbfunction-filter-resize-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var cpuCR examplecomv1.CPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "dftest-wbfunction-filter-resize-high-infer-main",
				Namespace: TESTNAMESPACE,
			},
				&cpuCR)
			if err != nil {
				// Error route
				fmt.Println("Cannot get CPUFunctionCR:", cpuCR, err)
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(*cpuCR.Status.FunctionIndex).To(Equal(int32(99)))
			var cpupod corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "dftest-wbfunction-filter-resize-high-infer-main-cpu-pod",
				Namespace: TESTNAMESPACE,
			},
				&cpupod)
			if err != nil {
				// Error route
				fmt.Println("Cannot get Pod:", cpupod, err)
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(*cpuCR.Status.PodName).To(Equal(cpupod.Name))
		})
		It("Test_2-1_UPDATE", func() {

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, cpufunctestUPDATE)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "cpufunctestupdate-wbfunction-decode-main",
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

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, cpufunctestDELETE)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			var cpuCR examplecomv1.CPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctestdelete-wbfunction-decode-main",
				Namespace: TESTNAMESPACE,
			},
				&cpuCR)

			err = k8sClient.Delete(ctx, &cpuCR)

			var cpubCR examplecomv1.CPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctestdelete-wbfunction-decode-main",
				Namespace: TESTNAMESPACE,
			},
				&cpubCR)

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "cpufunctestdelete-wbfunction-decode-main",
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

			var cpuaCR examplecomv1.CPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctestdelete-wbfunction-decode-main",
				Namespace: TESTNAMESPACE,
			},
				&cpuaCR)

			Expect(err).To(MatchError(ContainSubstring("not found")))

		})
		It("Test_2-1-4_Lifecycle", func() {
			// Create CPUFuncConfig
			err = createConfig(ctx, cpuconfigdecode214)
			if err != nil {
				fmt.Println("There is a problem in createing configmap ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection214)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create FPGAFunctionCR
			err = createFPGAFunction(ctx, FPGAFunction214)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create NetworkAttachmentDifinition
			err = createNetworkAttachmentDefinition(ctx, NetworkAttachmentDefinition214)
			if err != nil {
				fmt.Println("There is a problem in createing NetworkAttachmentDefinition ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction214)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: TESTNAMESPACE,
					Name:      "cpufunctest214-wbfunction-decode-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var cpuCR examplecomv1.CPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctest214-wbfunction-decode-main",
				Namespace: TESTNAMESPACE,
			},
				&cpuCR)
			Expect(err).NotTo(HaveOccurred())

			var cpupod corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctest214-wbfunction-decode-main-cpu-pod",
				Namespace: TESTNAMESPACE,
			},
				&cpupod)
			if err != nil {
				// error handling
				fmt.Println("Could not get Pod:", cpupod, err)
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(*cpuCR.Status.PodName).To(Equal(cpupod.Name))

			// check Lifecycle is not nil
			var licynil corev1.Lifecycle
			Expect(cpupod.Spec.Containers[0].Lifecycle).NotTo(HaveValue(Equal(licynil)))

			// check logs
			Expect(writer.String()).To(ContainSubstring("nextConnectionCRName : cpufunctest214-wbconnection-decode-main-filter-resize-low-infer-main"))
			Expect(writer.String()).To(ContainSubstring("rxProtocol : RTP"))
			Expect(writer.String()).To(ContainSubstring("txProtocol : DMA"))
			Expect(writer.String()).To(ContainSubstring("Success to create Pod."))
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

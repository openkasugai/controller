package controller

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	examplecomv1 "CPUFunction/api/v1"

	// ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
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
	// tmp.Kind = cpufcr.Kind
	// tmp.APIVersion = cpufcr.APIVersion
	err := k8sClient.Create(ctx, tmp)
	// err := k8sClient.Create(context.Background(), &CPUFunction1)
	if err != nil {
		return err
	}
	// tmp.Status = cpufcr.Status
	// err = k8sClient.Status().Update(ctx, tmp)
	// if err != nil {
	// 	return err
	// }
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

	// tmp.Status = pcieccr.Status
	// err = k8sClient.Status().Update(ctx, tmp)
	// if err != nil {
	// 	return err
	// }
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

func createFPGAFunction(ctx context.Context, fpgafcr FPGAFunction) error {
	tmp := &FPGAFunction{}
	*tmp = fpgafcr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	// tmp.Status = fpgafcr.Status
	// err = k8sClient.Status().Update(ctx, tmp)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func createCPUFuncConfig(ctx context.Context, cpuconf corev1.ConfigMap) error {
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
	// tmp.Kind = NetworkAttachmentDefinition1.Kind
	// tmp.APIVersion = NetworkAttachmentDefinition1.APIVersion
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
	// var stopFunc func()

	// This test case is for reconciler
	Context("Test for CPUFunctionReconciler", Ordered, func() {
		var reconciler CPUFunctionReconciler

		BeforeAll(func() {
			// set manager
			mgr, err = getMgr(mgr)
			Expect(err).NotTo(HaveOccurred())

			reconciler = CPUFunctionReconciler{
				Client:   k8sClient,
				Scheme:   testScheme,
				Recorder: mgr.GetEventRecorderFor("cpufunction-controller"),
			}

			err = reconciler.SetupWithManager(mgr)
			if err != nil {
				fmt.Println("Error in SetupWithManager", err)
			}
			Expect(err).NotTo(HaveOccurred())

			namespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					//Namespace: "test01",
					Name: "test01",
				},
			}
			err := k8sClient.Create(context.Background(), namespace)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
		})
		// Before Context runs, BeforeAll is executed once.
		BeforeEach(func() {

			os.Setenv("K8S_NODENAME", "node01")
			os.Setenv("K8S_CLUSTERNAME", "default")
			os.Setenv("K8S_GPU_MS_PORT", "8082")
			os.Setenv("K8S_GPU_HC_PORT", "8092")
			os.Setenv("K8S_CPU_MS_PORT", "8083")
			os.Setenv("K8S_CPU_HC_PORT", "8093")
			os.Setenv("K8S_GATE_MS_PORT", "8084")
			os.Setenv("K8S_GATE_HC_PORT", "8094")
			os.Setenv("K8S_DI_MS_PORT", "8085")
			os.Setenv("K8S_DI_HC_PORT", "8095")
			os.Setenv("K8S_WBF_MS_PORT", "8086")
			os.Setenv("K8S_WBF_HC_PORT", "8096")
			os.Setenv("K8S_FPGA_MS_PORT", "8087")
			os.Setenv("K8S_FPGA_HC_PORT", "8097")

			// stopFunc = startMgr(ctx, mgr)

		})

		// Each time It runs, BeforeEach is executed.
		BeforeEach(func() {

			// To delete crdata It
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.CPUFunction{}, client.InNamespace("test01"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &corev1.Pod{}, client.InNamespace("test01"))
			Expect(err).NotTo(HaveOccurred())
			// var yamlname string = "example.com_pcieconnections.yaml"
			// //pcieyaml := strings.Split(string(cm.Data[yamlName]), SPLIT)
			// pcieyaml, errors := ReadFile(yamlname)
			// if errors != nil {
			// 	fmt.Println("Error :", err)
			// 	return
			// }
			// var PCIeConnectionCRD unstructured.Unstructured
			// err = yaml.Unmarshal([]byte(pcieyaml), &PCIeConnectionCRD)
			// _, err = ctrl.CreateOrUpdate(ctx, r.Client, &PCIeConnectionCRD, func() error {
			// 				return nil
			// })

			err = k8sClient.DeleteAllOf(ctx, &PCIeConnection{}, client.InNamespace("test01"))
			if err != nil {
				fmt.Println("Can not delete PCIeConnectionCR", err)
			}
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &EthernetConnection{}, client.InNamespace("test01"))
			if err != nil {
				fmt.Println("Can not delete EthernetConnectionCR", err)
			}
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &FPGAFunction{}, client.InNamespace("test01"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(100 * time.Millisecond)

		})

		// AfterAll(func() {
		// 	// stop manager
		// 	stopFunc()
		// 	time.Sleep(100 * time.Millisecond)
		// })
		/*
			It("Test_1-2-1_redeployment", func() {
				// Create CPUFuncConfig
				err = createCPUFuncConfig(ctx, cpuconfigfrhigh)
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
				err = createCPUFunction(ctx, CPUFunction3)
				if err != nil {
					fmt.Println("There is a problem in createing CPUFunc-filter-resizeCR ", err)
				}
				Expect(err).NotTo(HaveOccurred())

				got, err := reconciler.Reconcile(ctx,
					ctrl.Request{NamespacedName: types.NamespacedName{
						Namespace: "test01",
						Name:      "dftest-wbfunction-filter-resize-high-infer-main",
					}})
				Expect(got).To(Equal(ctrl.Result{}))
				Expect(err).NotTo(HaveOccurred())

				var cpuCR examplecomv1.CPUFunction
				err = k8sClient.Get(ctx, client.ObjectKey{
					Name:      "dftest-wbfunction-filter-resize-high-infer-main",
					Namespace: "test01",
				},
					&cpuCR)
				if err != nil {
					// Error route
					fmt.Println("Cannot get CPUFunctionCR:", cpuCR, err)
				}
				Expect(err).NotTo(HaveOccurred())

				Expect(*cpuCR.Status.FunctionIndex).To(Equal(int32(0)))

				// delete CPUFunctionCR
				err = deleteCPUFunction(ctx, CPUFunction3)
				if err != nil {
					fmt.Println("There is a problem in deleteing CPUFunc-filter-resizeCR ", err)
				}
				Expect(err).NotTo(HaveOccurred())

				got, err = reconciler.Reconcile(ctx,
					ctrl.Request{NamespacedName: types.NamespacedName{
						Namespace: "test01",
						Name:      "dftest-wbfunction-filter-resize-high-infer-main",
					}})
				Expect(got).To(Equal(ctrl.Result{}))
				Expect(err).NotTo(HaveOccurred())

				// Create CPUFunctionCR
				err = createCPUFunction(ctx, CPUFunction3)
				if err != nil {
					fmt.Println("There is a problem in createing CPUFunc-filter-resizeCR ", err)
				}
				Expect(err).NotTo(HaveOccurred())

				got, err = reconciler.Reconcile(ctx,
					ctrl.Request{NamespacedName: types.NamespacedName{
						Namespace: "test01",
						Name:      "dftest-wbfunction-filter-resize-high-infer-main",
					}})
				Expect(got).To(Equal(ctrl.Result{}))
				Expect(err).NotTo(HaveOccurred())

				var cpuCRredeploy examplecomv1.CPUFunction
				err = k8sClient.Get(ctx, client.ObjectKey{
					Name:      "dftest-wbfunction-filter-resize-high-infer-main",
					Namespace: "test01",
				},
					&cpuCRredeploy)
				if err != nil {
					// Error route
					fmt.Println("Cannot get CPUFunctionCR:", cpuCRredeploy, err)
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(*cpuCR.Status.FunctionIndex).To(Equal(int32(0)))
			})
			It("Test_1-2-2_next-deployment", func() {
				// Create CPUFuncConfig
				err = createCPUFuncConfig(ctx, cpuconfigfrhigh)
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

				// Create EthernetConnectionCR
				err = createEthernetConnection(ctx, EthernetConnection1221)
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

				got, err := reconciler.Reconcile(ctx,
					ctrl.Request{NamespacedName: types.NamespacedName{
						Namespace: "test01",
						Name:      "dftest122-wbfunction-filter-resize-high-infer-main",
					}})
				Expect(got).To(Equal(ctrl.Result{}))
				Expect(err).NotTo(HaveOccurred())

				var cpuCR examplecomv1.CPUFunction
				err = k8sClient.Get(ctx, client.ObjectKey{
					Name:      "dftest-wbfunction-filter-resize-high-infer-main",
					Namespace: "test01",
				},
					&cpuCR)
				if err != nil {
					// Error route
					fmt.Println("Cannot get CPUFunctionCR:", cpuCR, err)
				}
				Expect(err).NotTo(HaveOccurred())

				Expect(*cpuCR.Status.FunctionIndex).To(Equal(int32(99)))
			})
		*/
		//Test for GetFunc
		It("Test_1-1-1_decode-dma", func() {
			// Create CPUFuncConfig
			err = createCPUFuncConfig(ctx, cpuconfigdecode)
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
					Namespace: "test01",
					Name:      "cpufunctiontest1-wbfunction-decode-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())
			fmt.Println("Reconcile Errors:", err)
			Expect(err).NotTo(HaveOccurred())

			var cpuCR examplecomv1.CPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctiontest1-wbfunction-decode-main",
				Namespace: "test01",
			},
				&cpuCR)
			Expect(err).NotTo(HaveOccurred())
			var pcieCR PCIeConnection
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctiontest1-wbconnection-decode-main-filter-resize-low-infer-main",
				Namespace: "test01",
			},
				&pcieCR)
			fmt.Println("getpcieCR Error", err)
			fmt.Println("Kind:", pcieCR.Kind)
			fmt.Println("Spec:", pcieCR.Spec)
			// yamlData2, err := yaml.Marshal(CPUFunction1)
			// fmt.Println("Input CPUFunctionCR:", string(yamlData2))
			// fmt.Println("K8s.Client CPU CR:", cpuCR)
			// yamlData, err := yaml.Marshal(cpuCR)
			// fmt.Println("Pod output in YAML format:", string(yamlData))

			Expect(err).NotTo(HaveOccurred())

			// time.Sleep(100 * time.Millisecond)

			var cpupod corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctiontest1-wbfunction-decode-main-cpu-pod",
				Namespace: "test01",
			},
				&cpupod)
			if err != nil {
				// Error route
				fmt.Println("Could not get Pod:", cpupod, err)
				// return
			}
			fmt.Println("Pod creation complete:", cpupod)
			Expect(err).NotTo(HaveOccurred())

			// encoder := yaml.NewEncoder(os.Stdout)
			// err = encoder.Encode(cpupod)
			// if err != nil {
			// 	fmt.Println("Could not output pod in YAML format:", err)
			// 	return
			// }

			Expect(err).NotTo(HaveOccurred())
			var envVarFPS string
			var envVarPORT string
			var envVarIPA string
			for _, containers := range cpupod.Spec.Containers {

				for _, envVar := range containers.Env {
					fmt.Println(envVar.Name, "：", envVar.Value)
					if envVar.Name == "DECENV_FRAME_FPS" {
						envVarFPS = envVar.Value
						fmt.Println("The value of DECENV_FRAME_FPS：", envVarFPS)
					}
					if envVar.Name == "DECENV_VIDEOSRC_PORT" {
						envVarPORT = envVar.Value
						fmt.Println("The value of DECENV_VIDEOSRC_PORT：", envVarPORT)
					}
					if envVar.Name == "DECENV_VIDEOSRC_IPA" {
						envVarIPA = envVar.Value
						fmt.Println("The value of DECENV_VIDEOSRC_IPA：", envVarIPA)
					}
				}
			}
			fmt.Println(envVarFPS)
			// confirmation Pod Yaml
			expectedEnvFPSVar := corev1.EnvVar{
				Name:  "DECENV_FRAME_FPS",
				Value: strconv.Itoa(int(CPUFunction1.Spec.Params["decEnvFrameFPS"].IntVal)),
			}

			Expect(err).NotTo(HaveOccurred())
			Expect(envVarFPS).To(Equal(expectedEnvFPSVar.Value))
			fmt.Println(envVarFPS)
			// confirmation Pod Yaml
			expectedEnvPORTVar := corev1.EnvVar{
				Name:  "DECENV_VIDEOSRC_PORT",
				Value: strconv.Itoa(int(CPUFunction1.Spec.Params["inputPort"].IntVal)),
			}

			Expect(err).NotTo(HaveOccurred())
			Expect(envVarPORT).To(Equal(expectedEnvPORTVar.Value))
			expectedEnvIPAVar := corev1.EnvVar{
				Name:  "DECENV_VIDEOSRC_IPA",
				Value: CPUFunction1.Spec.Params["inputIPAddress"].StrVal,
			}
			Expect(envVarIPA).To(Equal(expectedEnvIPAVar.Value))
			// Expect(envVarValue).To(Equal("DECENV_FRAME_FPS"))
			//actualEnvVar := findEnvVarByName(cpupod.Spec.Containers.Env, "DECENV_FRAME_FPS")

			//	Expect(envVar).To(Equal(expectedEnvVar))

			// var cpupod2 corev1.Pod
			// err = k8sClient.Get(ctx, client.ObjectKey{
			// 	Name:      "cpufunctiontest1-wbfunction-decode-main-cpu-pod",
			// 	Namespace: "default",
			// },
			// 	&cpupod2)
			// Expect(err).NotTo(HaveOccurred())
			// if err != nil {
			// 	// Error route
			// 	fmt.Println("Could not get Pod:", cpupod2, err)
			// 	return
			// }

		})
		It("Test_1-1-2_decode-tcp", func() {
			// Create CPUFuncConfig
			err = createCPUFuncConfig(ctx, cpuconfigdecode)
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
					Namespace: "test01",
					Name:      "cpufunctiontest2-wbfunction-decode-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var cpupod corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctiontest2-wbfunction-decode-main-cpu-pod",
				Namespace: "test01",
			},
				&cpupod)
			if err != nil {
				// Error route
				fmt.Println("Could not get Pod:", cpupod, err)
			}
			fmt.Println("Complete to create Pod:", cpupod)
			Expect(err).NotTo(HaveOccurred())

			var envVarFPS string
			var envVarPORT string
			var envVarIPA string
			for _, containers := range cpupod.Spec.Containers {

				for _, envVar := range containers.Env {
					fmt.Println(envVar.Name, "：", envVar.Value)
					if envVar.Name == "DECENV_FRAME_FPS" {
						envVarFPS = envVar.Value
						fmt.Println("The value of DECENV_FRAME_FPS：", envVarFPS)
					}
					if envVar.Name == "DECENV_VIDEOSRC_PORT" {
						envVarPORT = envVar.Value
						fmt.Println("The value of DECENV_VIDEOSRC_PORT：", envVarPORT)
					}
					if envVar.Name == "DECENV_VIDEOSRC_IPA" {
						envVarIPA = envVar.Value
						fmt.Println("The value of DECENV_VIDEOSRC_IPA：", envVarIPA)
					}
				}
			}
			fmt.Println(envVarFPS)
			// confirmation Pod Yaml
			expectedEnvFPSVar := corev1.EnvVar{
				Name:  "DECENV_FRAME_FPS",
				Value: strconv.Itoa(int(CPUFunction2.Spec.Params["decEnvFrameFPS"].IntVal)),
			}

			Expect(err).NotTo(HaveOccurred())
			Expect(envVarFPS).To(Equal(expectedEnvFPSVar.Value))
			fmt.Println(envVarFPS)
			// confirmation Pod Yaml
			expectedEnvPORTVar := corev1.EnvVar{
				Name:  "DECENV_VIDEOSRC_PORT",
				Value: strconv.Itoa(int(CPUFunction2.Spec.Params["inputPort"].IntVal)),
			}

			Expect(err).NotTo(HaveOccurred())
			Expect(envVarPORT).To(Equal(expectedEnvPORTVar.Value))
			expectedEnvIPAVar := corev1.EnvVar{
				Name:  "DECENV_VIDEOSRC_IPA",
				Value: CPUFunction2.Spec.Params["inputIPAddress"].StrVal,
			}
			Expect(envVarIPA).To(Equal(expectedEnvIPAVar.Value))
		})
		It("Test_1-1-3_filter-resize-high-infer", func() {
			// Create CPUFuncConfig
			err = createCPUFuncConfig(ctx, cpuconfigfrhigh)
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
			err = createCPUFunction(ctx, CPUFunction3)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFunc-filter-resizeCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// // Create CPUFunctionCR
			// err = createCPUFunction(ctx, CPUFunction1)
			// if err != nil {
			// 	fmt.Println("There is a problem in createing CPUFunc-decodeCR ", err)
			// }
			// Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "test01",
					Name:      "dftest-wbfunction-filter-resize-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var cpupod corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "dftest-wbfunction-filter-resize-high-infer-main-cpu-pod",
				Namespace: "test01",
			},
				&cpupod)
			if err != nil {
				// Error route
				fmt.Println("Cannot get Pod:", cpupod, err)
			}
			fmt.Println("Complete to create Pod:", cpupod)
			Expect(err).NotTo(HaveOccurred())

			var envVarInputPort string
			var envVarOutputIP string
			var envVarOutputPort string
			for _, containers := range cpupod.Spec.Containers {

				for _, envVar := range containers.Env {
					fmt.Println(envVar.Name, "：", envVar.Value)
					if envVar.Name == "FRENV_INPUT_PORT" {
						envVarInputPort = envVar.Value
						fmt.Println("The value of FRENV_INPUT_PORT：", envVarInputPort)
					}
					if envVar.Name == "FRENV_OUTPUT_IP" {
						envVarOutputIP = envVar.Value
						fmt.Println("The value of FRENV_OUTPUT_IP：", envVarOutputIP)
					}
					if envVar.Name == "FRENV_OUTPUT_PORT" {
						envVarOutputPort = envVar.Value
						fmt.Println("The value of FRENV_OUTPUT_PORT：", envVarOutputPort)
					}
				}
			}
			fmt.Println(envVarInputPort)
			// confirmation Pod Yaml
			expectedEnvInputPortVar := corev1.EnvVar{
				Name:  "FRENV_INPUT_PORT",
				Value: strconv.Itoa(int(CPUFunction3.Spec.Params["inputPort"].IntVal)),
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(envVarInputPort).To(Equal(expectedEnvInputPortVar.Value))

			fmt.Println(envVarInputPort)

			// confirmation Pod Yaml
			expectedEnvOutputIPVar := corev1.EnvVar{
				Name:  "FRENV_OUTPUT_IP",
				Value: CPUFunction3.Spec.Params["outputIPAddress"].StrVal,
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(envVarOutputIP).To(Equal(expectedEnvOutputIPVar.Value))

			expectedEnvOutputPortVar := corev1.EnvVar{
				Name:  "FRENV_OUTPUT_PORT",
				Value: strconv.Itoa(int(CPUFunction3.Spec.Params["outputPort"].IntVal)),
			}
			Expect(envVarOutputPort).To(Equal(expectedEnvOutputPortVar.Value))

		})
		// // It("Test_1-1-4_filter-resize-low-infer", func() {
		// 	// Create CPUFuncConfig
		// 	err = createCPUFuncConfig(ctx, cpuconfigfrlow)
		// 	if err != nil {
		// 		fmt.Println("There is a problem in createing configmap ", err)
		// 	}
		// 	Expect(err).NotTo(HaveOccurred())

		// 	// Create EthernetConnectionCR
		// 	err = createEthernetConnection(ctx, EthernetConnection1)
		// 	if err != nil {
		// 		fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
		// 	}
		// 	Expect(err).NotTo(HaveOccurred())

		// 	// Create EthernetConnectionCR
		// 	err = createEthernetConnection(ctx, EthernetConnection3)
		// 	if err != nil {
		// 		fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
		// 	}
		// 	Expect(err).NotTo(HaveOccurred())

		// 	// Create CPUFunctionCR
		// 	err = createCPUFunction(ctx, CPUFunction3)
		// 	if err != nil {
		// 		fmt.Println("There is a problem in createing CPUFunc-filter-resizeCR ", err)
		// 	}
		// 	Expect(err).NotTo(HaveOccurred())

		// 	// // Create CPUFunctionCR
		// 	// err = createCPUFunction(ctx, CPUFunction1)
		// 	// if err != nil {
		// 	// 	fmt.Println("There is a problem in createing CPUFunc-decodeCR ", err)
		// 	// }
		// 	// Expect(err).NotTo(HaveOccurred())

		// 	got, err := reconciler.Reconcile(ctx,
		// 		ctrl.Request{NamespacedName: types.NamespacedName{
		// 			Namespace: "test01",
		// 			Name:      "dftest-wbfunction-filter-resize-high-infer-main",
		// 		}})
		// 	Expect(got).To(Equal(ctrl.Result{}))
		// 	Expect(err).NotTo(HaveOccurred())

		// 	var cpupod corev1.Pod
		// 	err = k8sClient.Get(ctx, client.ObjectKey{
		// 		Name:      "dftest-wbfunction-filter-resize-high-infer-main-cpu-pod",
		// 		Namespace: "test01",
		// 	},
		// 		&cpupod)
		// 	if err != nil {
		// 		// Erro route
		// 		fmt.Println("Cannot get Pod:", cpupod, err)
		// 	}
		// 	fmt.Println("Complete to create Pod:", cpupod)
		// 	Expect(err).NotTo(HaveOccurred())

		// 	var envVarInputPort string
		// 	var envVarOutputIP string
		// 	var envVarOutputPort string
		// 	for _, containers := range cpupod.Spec.Containers {

		// 		for _, envVar := range containers.Env {
		// 			fmt.Println(envVar.Name, "：", envVar.Value)
		// 			if envVar.Name == "FRENV_INPUT_PORT" {
		// 				envVarInputPort = envVar.Value
		// 				fmt.Println("The value of FRENV_INPUT_PORT：", envVarInputPort)
		// 			}
		// 			if envVar.Name == "FRENV_OUTPUT_IP" {
		// 				envVarOutputIP = envVar.Value
		// 				fmt.Println("The value of FRENV_OUTPUT_IP：", envVarOutputIP)
		// 			}
		// 			if envVar.Name == "FRENV_OUTPUT_PORT" {
		// 				envVarOutputPort = envVar.Value
		// 				fmt.Println("The value of FRENV_OUTPUT_PORT：", envVarOutputPort)
		// 			}
		// 		}
		// 	}
		// 	fmt.Println(envVarInputPort)
		// 	// confirmation Pod Yaml
		// 	expectedEnvInputPortVar := corev1.EnvVar{
		// 		Name:  "FRENV_INPUT_PORT",
		// 		Value: strconv.Itoa(int(CPUFunction3.Spec.Params["inputPort"].IntVal)),
		// 	}
		// 	Expect(err).NotTo(HaveOccurred())
		// 	Expect(envVarInputPort).To(Equal(expectedEnvInputPortVar.Value))

		// 	fmt.Println(envVarInputPort)

		// 	// confirmation Pod Yaml
		// 	expectedEnvOutputIPVar := corev1.EnvVar{
		// 		Name:  "FRENV_OUTPUT_IP",
		// 		Value: CPUFunction3.Spec.Params["outputIPAddress"].StrVal,
		// 	}
		// 	Expect(err).NotTo(HaveOccurred())
		// 	Expect(envVarOutputIP).To(Equal(expectedEnvOutputIPVar.Value))

		// 	expectedEnvOutputPortVar := corev1.EnvVar{
		// 		Name:  "FRENV_OUTPUT_PORT",
		// 		Value: strconv.Itoa(int(CPUFunction3.Spec.Params["outputPort"].IntVal)),
		// 	}
		// 	Expect(envVarOutputPort).To(Equal(expectedEnvOutputPortVar.Value))

		// })
		It("Test_1-1-5_copy-branch", func() {
			// Create CPUFuncConfig
			err = createCPUFuncConfig(ctx, cpuconfigcopybranch)
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

			// // Create CPUFunctionCR
			// err = createCPUFunction(ctx, CPUFunction4frlow)
			// if err != nil {
			// 	fmt.Println("There is a problem in createing CPUFuncCR-filter-resize-low-infer ", err)
			// }
			// Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction4)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR-copy-branch", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "test01",
					Name:      "cpufunctiontest4-wbfunction-copy-branch-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var cpupod corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctiontest4-wbfunction-copy-branch-main-cpu-pod",
				Namespace: "test01",
			},
				&cpupod)
			if err != nil {
				// Error route
				fmt.Println("Cannot get Pod:", cpupod, err)
			}
			fmt.Println("Complete to create Pod:", cpupod)
			Expect(err).NotTo(HaveOccurred())

		})
		It("Test_1-1-6_glue-tcp-to-dma", func() {
			// Create CPUFuncConfig
			err = createCPUFuncConfig(ctx, cpuconfigglue)
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
			err = createCPUFunction(ctx, CPUFunction5)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection5glue)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "test01",
					Name:      "cpufunctiontest5-wbfunction-glue-fdma-to-tcp-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var cpupod corev1.Pod
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "cpufunctiontest5-wbfunction-glue-fdma-to-tcp-main-cpu-pod",
				Namespace: "test01",
			},
				&cpupod)
			if err != nil {
				// Error route
				fmt.Println("Cannot get Pod:", cpupod, err)
			}
			fmt.Println("Complete to create Pod:", cpupod)
			Expect(err).NotTo(HaveOccurred())

		})
		It("Test_1-2-1_redeployment", func() {
			// Create CPUFuncConfig
			err = createCPUFuncConfig(ctx, cpuconfigfrhigh)
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

			/*
				// Create CPUFunctionCR
				err = createCPUFunction(ctx, CPUFunction3)
				if err != nil {
					fmt.Println("There is a problem in createing CPUFunc-filter-resizeCR ", err)
				}
				Expect(err).NotTo(HaveOccurred())

				got, err := reconciler.Reconcile(ctx,
					ctrl.Request{NamespacedName: types.NamespacedName{
						Namespace: "test01",
						Name:      "dftest-wbfunction-filter-resize-high-infer-main",
					}})
				Expect(got).To(Equal(ctrl.Result{}))
				Expect(err).NotTo(HaveOccurred())
			*/
			var cpuCR examplecomv1.CPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "dftest-wbfunction-filter-resize-high-infer-main",
				Namespace: "test01",
			},
				&cpuCR)
			if err != nil {
				// Error route
				fmt.Println("Cannot get CPUFunctionCR:", cpuCR, err)
			}
			Expect(err).NotTo(HaveOccurred())

			Expect(*cpuCR.Status.FunctionIndex).To(Equal(int32(2)))

			// delete CPUFunctionCR
			err = deleteCPUFunction(ctx, CPUFunction3)
			if err != nil {
				fmt.Println("There is a problem in deleteing CPUFunc-filter-resizeCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "test01",
					Name:      "dftest-wbfunction-filter-resize-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction3)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFunc-filter-resizeCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "test01",
					Name:      "dftest-wbfunction-filter-resize-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var cpuCRredeploy examplecomv1.CPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "dftest-wbfunction-filter-resize-high-infer-main",
				Namespace: "test01",
			},
				&cpuCRredeploy)
			if err != nil {
				// Error route
				fmt.Println("Cannot get CPUFunctionCR:", cpuCRredeploy, err)
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(*cpuCR.Status.FunctionIndex).To(Equal(int32(2)))

			// delete CPUFunctionCR
			err = deleteCPUFunction(ctx, CPUFunction3)
			if err != nil {
				fmt.Println("There is a problem in deleteing CPUFunc-filter-resizeCR ", err)
			}
			Expect(err).NotTo(HaveOccurred())

			_, _ = reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "test01",
					Name:      "dftest-wbfunction-filter-resize-high-infer-main",
				}})
		})
		It("Test_1-2-2_next-deployment", func() {
			// Create CPUFuncConfig
			err = createCPUFuncConfig(ctx, cpuconfigfrhigh)
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

			got, err := reconciler.Reconcile(ctx,
				ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "test01",
					Name:      "dftest-wbfunction-filter-resize-high-infer-main",
				}})
			Expect(got).To(Equal(ctrl.Result{}))
			Expect(err).NotTo(HaveOccurred())

			var cpuCR examplecomv1.CPUFunction
			err = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "dftest-wbfunction-filter-resize-high-infer-main",
				Namespace: "test01",
			},
				&cpuCR)
			if err != nil {
				// Error route
				fmt.Println("Cannot get CPUFunctionCR:", cpuCR, err)
			}
			Expect(err).NotTo(HaveOccurred())

			Expect(*cpuCR.Status.FunctionIndex).To(Equal(int32(99)))
		})
	})
})

// func startMgr(ctx context.Context, mgr manager.Manager) func() {
// 	ctx, cancel := context.WithCancel(ctx)
// 	go func() {
// 		err := mgr.Start(ctx)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}()
// 	time.Sleep(100 * time.Millisecond)
// 	return cancel
// }

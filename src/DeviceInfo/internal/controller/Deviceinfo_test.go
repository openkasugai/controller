package controller

import (
	"bytes"
	"context"
	"fmt"
	"os"

	ctrl "sigs.k8s.io/controller-runtime"
	//  "sigs.k8s.io/controller-runtime/pkg/client"

	examplecomv1 "DeviceInfo/api/v1"

	// ntthpcv1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"errors"
	// v1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

var (
	ErrElementsNotEqual      = errors.New("Elements doesn't equal")
	ErrInvalidPathIsInputeed = errors.New("Invalid path is inputted")
	ErrFailedToConvert       = errors.New("Failed to convert")
)

func getMgr(mgr ctrl.Manager) (ctrl.Manager, error) {
	if mgr == nil {
		return ctrl.NewManager(cfg, ctrl.Options{
			Scheme: testScheme,
		})
	}
	return mgr, nil
}

func decodeManifest(fileName string) (obj runtime.Object, err error) {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return obj, err
	}
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err = decode(bytes, nil, nil)
	return obj, err
}

func GetResourceFromYaml[T any](path string) (T, error) {

	obj, err := decodeManifest(path)
	if err != nil {
		var tmp T
		return tmp, err
	}
	ret, ok := obj.(T)
	if !ok {
		var tmp T
		return tmp, ErrFailedToConvert
	} else {
		return ret, nil
	}
}

// Create ComputeResourceCR
func createCompureResource(ctx context.Context, comres examplecomv1.ComputeResource) error {
	tmp := &examplecomv1.ComputeResource{}
	*tmp = comres
	tmp.TypeMeta = comres.TypeMeta
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	/*
	   err = k8sClient.Status().Update(ctx, tmp)
	   if err != nil {
	       return err
	   }
	*/
	return nil
}

func createFPGACR(ctx context.Context, fpgaCR examplecomv1.FPGA) error {
	tmp := &examplecomv1.FPGA{}
	*tmp = fpgaCR
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

func createChildBsCR(ctx context.Context, childBsCR examplecomv1.ChildBs) error {
	tmp := &examplecomv1.ChildBs{}
	*tmp = childBsCR
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

func createInfraInfoConfig(ctx context.Context, infraConfig corev1.ConfigMap) error {
	tmp := &corev1.ConfigMap{}
	*tmp = infraConfig
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

func createDeployInfoConfig(ctx context.Context, deployInfoConfig corev1.ConfigMap) error {
	tmp := &corev1.ConfigMap{}
	*tmp = deployInfoConfig
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

func updateDeployInfoConfig(ctx context.Context, deployInfoConfig corev1.ConfigMap) error {
	tmp := &corev1.ConfigMap{}
	*tmp = deployInfoConfig
	err := k8sClient.Update(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create DeviceInfoCR
func createDeviceInfo(ctx context.Context, devinfo examplecomv1.DeviceInfo) error {
	tmp := &examplecomv1.DeviceInfo{}
	*tmp = devinfo
	tmp.TypeMeta = devinfo.TypeMeta
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	/*
		err = k8sClient.Status().Update(ctx, tmp)
		if err != nil {
			return err
		}
	*/
	return nil
}

/*
	func updateDeviceInfo(ctx context.Context, devinfo examplecomv1.DeviceInfo) error {
		tmp := &examplecomv1.DeviceInfo{}
		*tmp = devinfo
		tmp.TypeMeta = devinfo.TypeMeta
		err := k8sClient.Update(ctx, tmp)
		if err != nil {
			return err
		}
		return nil
	}
*/
func deleteDeviceInfo(ctx context.Context, devinfo examplecomv1.DeviceInfo) error {
	tmp := &examplecomv1.DeviceInfo{}
	*tmp = devinfo
	tmp.TypeMeta = devinfo.TypeMeta
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

	Context("Test for DeviceInfoCR", Ordered, func() {
		var reconciler DeviceInfoReconciler
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

			reconciler = DeviceInfoReconciler{
				Client:   k8sClient,
				Scheme:   testScheme,
				Recorder: fakerecorder,
				// Recorder: mgr.GetEventRecorderFor("deviceinfo-controller"),
			}
			err = reconciler.SetupWithManager(mgr)
			if err != nil {
				fmt.Println("Error in SetupWithManager")
			}
			Expect(err).NotTo(HaveOccurred())
			os.Setenv("K8S_NODENAME", "test01")
			os.Setenv("K8S_CLUSTERNAME", "default")
			namespace := &corev1.Namespace{}
			namespace.Name = "test01"
			err := k8sClient.Create(context.Background(), namespace)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
		})

		BeforeEach(func() {
			writer.Reset()
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.DeviceInfo{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())

			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "deviceinfo-df-night01-wbfunction-decode-main",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}

			err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.ComputeResource{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGA{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.ChildBs{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
		})

		It("7-1", func() {
			By("Test Start")
			// Create ComputeResourceCR
			err = createInfraInfoConfig(ctx, infrainfo_configdata)
			if err != nil {
				fmt.Println("There is a problem in createing InfraInfo Config ", err)
				fmt.Println(err)
			}
			err = createDeployInfoConfig(ctx, deployinfo_configdata)
			if err != nil {
				fmt.Println("There is a problem in createing DeployInfo Config", err)
				fmt.Println(err)
			}
			err = StartupProccessing(&reconciler, mgr)
			if err != nil {
				fmt.Println("Error in StartupProccessing")
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(writer.String()).To(ContainSubstring("Startup Create Success"))
		})

		It("7-2", func() {
			By("Test Start")
			// Create ComputeResourceCR
			err = createInfraInfoConfig(ctx, infrainfo_configdata)
			if err != nil {
				fmt.Println("There is a problem in createing InfraInfo Config ", err)
				fmt.Println(err)
			}
			err = createDeployInfoConfig(ctx, deployinfo_configdata)
			if err != nil {
				fmt.Println("There is a problem in createing DeployInfo Config", err)
				fmt.Println(err)
			}
			err = StartupProccessing(&reconciler, mgr)
			if err != nil {
				fmt.Println("Error in StartupProccessing")
			}
			err = updateDeployInfoConfig(ctx, deployinfo_configdata2)
			if err != nil {
				fmt.Println("There is a problem in updateing DeployInfo Config", err)
				fmt.Println(err)
			}
			err = createFPGACR(ctx, fpgaCRdata)
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR ", err)
				fmt.Println(err)
			}
			err = createChildBsCR(ctx, childBsCRdata1)
			if err != nil {
				fmt.Println("There is a problem in createing ChildBsCR ", err)
				fmt.Println(err)
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "childbs1",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(writer.String()).To(ContainSubstring("Startup Create Success"))
			Expect(writer.String()).To(ContainSubstring("Reconcile start."))
			Expect(writer.String()).To(ContainSubstring("ComputeResource Update Success"))
			Expect(writer.String()).To(ContainSubstring("Reconcile end."))
		})

		It("7-3", func() {
			By("Test Start")
			// Create ComputeResourceCR
			err = createInfraInfoConfig(ctx, infrainfo_configdata)
			if err != nil {
				fmt.Println("There is a problem in createing InfraInfo Config ", err)
				fmt.Println(err)
			}
			err = createDeployInfoConfig(ctx, deployinfo_configdata)
			if err != nil {
				fmt.Println("There is a problem in createing DeployInfo Config", err)
				fmt.Println(err)
			}
			err = StartupProccessing(&reconciler, mgr)
			if err != nil {
				fmt.Println("Error in StartupProccessing")
			}
			err = updateDeployInfoConfig(ctx, deployinfo_configdata2)
			if err != nil {
				fmt.Println("There is a problem in updateing DeployInfo Config", err)
				fmt.Println(err)
			}
			err = createFPGACR(ctx, fpgaCRdata)
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR ", err)
				fmt.Println(err)
			}
			err = createChildBsCR(ctx, childBsCRdata2)
			if err != nil {
				fmt.Println("There is a problem in createing ChildBsCR ", err)
				fmt.Println(err)
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "childbs1",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(writer.String()).To(ContainSubstring("Startup Create Success"))
			Expect(writer.String()).To(ContainSubstring("Reconcile start."))
			Expect(writer.String()).To(ContainSubstring("ComputeResource Update Success"))
			Expect(writer.String()).To(ContainSubstring("Reconcile end."))
		})

		It("7-4", func() {
			By("Test Start")
			// Create ComputeResourceCR
			err = createCompureResource(ctx, comres1)
			if err != nil {
				fmt.Println("There is a problem in createing ComputeresourceCR ", err)
				fmt.Println(err)
			}
			err = createInfraInfoConfig(ctx, infrainfo_configdata)
			if err != nil {
				fmt.Println("There is a problem in createing InfraInfo Config ", err)
				fmt.Println(err)
			}
			err = createDeployInfoConfig(ctx, deployinfo_configdata)
			if err != nil {
				fmt.Println("There is a problem in createing DeployInfo Config", err)
				fmt.Println(err)
			}
			// Create DeviceInfoCR
			err = createDeviceInfo(ctx, DeviceInfo1)
			if err != nil {
				fmt.Println("There is a problem in createing DeviceInfo ", err)
				fmt.Println(err)
			}

			err = StartupProccessing(&reconciler, mgr)
			if err != nil {
				fmt.Println("Error in StartupProccessing")
			}
			Expect(err).NotTo(HaveOccurred())

			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "deviceinfo-df-night01-wbfunction-decode-main",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(writer.String()).To(ContainSubstring("ComputeResourceCR is exist"))
			Expect(writer.String()).To(ContainSubstring("Reconcile start."))
			Expect(writer.String()).To(ContainSubstring("ComputeResource Update Success"))
			Expect(writer.String()).To(ContainSubstring("DeviceInfoCR Update Success"))
			Expect(writer.String()).To(ContainSubstring("Reconcile end."))
		})

		It("7-5", func() {
			By("Test Start")
			// Create ComputeResourceCR
			err = createCompureResource(ctx, comres1)
			if err != nil {
				fmt.Println("There is a problem in createing ComputeresourceCR ", err)
				fmt.Println(err)
			}
			err = createInfraInfoConfig(ctx, infrainfo_configdata)
			if err != nil {
				fmt.Println("There is a problem in createing InfraInfo Config ", err)
				fmt.Println(err)
			}
			err = createDeployInfoConfig(ctx, deployinfo_configdata)
			if err != nil {
				fmt.Println("There is a problem in createing DeployInfo Config", err)
				fmt.Println(err)
			}
			// Create DeviceInfoCR
			err = createDeviceInfo(ctx, DeviceInfo1)
			if err != nil {
				fmt.Println("There is a problem in createing DeviceInfo1 ", err)
				fmt.Println(err)
			}

			err = StartupProccessing(&reconciler, mgr)
			if err != nil {
				fmt.Println("Error in StartupProccessing")
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "deviceinfo-df-night01-wbfunction-decode-main",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())
			var devinfocr examplecomv1.DeviceInfo
			_ = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "deviceinfo-df-night01-wbfunction-decode-main",
				Namespace: "default",
			},
				&devinfocr)
			if devinfocr.Status.Response.Status == "Deployed" {
				fmt.Println("Change Undeploy Start")

				var deltime metav1.Time

				deltime = metav1.Now()

				devinfocr.DeletionTimestamp = &deltime
				err = deleteDeviceInfo(ctx, devinfocr)
				if err != nil {
					fmt.Println("There is a problem in updateing devinfocr ", err)
					fmt.Println(err)
				}

				_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "deviceinfo-df-night01-wbfunction-decode-main",
				}})
				if err != nil {
					By("Reconcile Error")
					fmt.Println(err)
				}

				err = createDeviceInfo(ctx, DeviceInfo2)
				if err != nil {
					fmt.Println("There is a problem in createing DeviceInfo2 ", err)
					fmt.Println(err)
				}
				_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "deviceinfo-df-night01-wbfunction-decode-main",
				}})
				if err != nil {
					By("Reconcile Error")
					fmt.Println(err)
				}
				_ = k8sClient.Get(ctx, client.ObjectKey{
					Name:      "deviceinfo-df-night01-wbfunction-decode-main",
					Namespace: "default",
				},
					&devinfocr)
				if devinfocr.Status.Response.Status == "Undeployed" {
					var cpr examplecomv1.ComputeResource
					_ = k8sClient.Get(ctx, client.ObjectKey{
						Name:      "compute-test01",
						Namespace: "default",
					},
						&cpr)
					for _, val := range cpr.Spec.Regions {
						if val.DeviceType != "cpu" {
							continue
						}
						var ni []examplecomv1.FunctionInfrastruct
						ni = nil

						Expect(val.Available).Should(Equal(chkComRes1.Available))
						Expect(val.CurrentCapacity).Should(Equal(chkComRes1.CurrentCapacity))
						Expect(val.CurrentFunctions).Should(Equal(chkComRes1.CurrentFunctions))
						Expect(val.DeviceFilePath).Should(Equal(chkComRes1.DeviceFilePath))
						Expect(val.DeviceIndex).Should(Equal(chkComRes1.DeviceIndex))
						Expect(val.DeviceType).Should(Equal(chkComRes1.DeviceType))
						Expect(val.DeviceUUID).Should(Equal(chkComRes1.DeviceUUID))
						Expect(val.MaxCapacity).Should(Equal(chkComRes1.MaxCapacity))
						Expect(val.MaxFunctions).Should(Equal(chkComRes1.MaxFunctions))
						Expect(val.Name).Should(Equal(chkComRes1.Name))
						Expect(val.Type).Should(Equal(chkComRes1.Type))
						Expect(val.Functions).Should(Equal(ni))
					}
				} else {
					fmt.Println("Status is NotUndeployed")
				}
			} else {
				fmt.Println("Status is NotDeployed")
			}
		})

		It("7-6", func() {
			By("Test Start")
			// Create ComputeResourceCR
			err = createCompureResource(ctx, comres1)
			if err != nil {
				fmt.Println("There is a problem in createing ComputeresourceCR ", err)
				fmt.Println(err)
			}
			// Create InfraInfoConfig
			err = createInfraInfoConfig(ctx, infrainfo_configdata)
			if err != nil {
				fmt.Println("There is a problem in createing InfraInfo Config ", err)
				fmt.Println(err)
			}
			// Create DeployInfoConfig
			err = createDeployInfoConfig(ctx, deployinfo_configdata)
			if err != nil {
				fmt.Println("There is a problem in createing DeployInfo Config", err)
				fmt.Println(err)
			}
			// Create DeviceInfoCR
			err = createDeviceInfo(ctx, DeviceInfo1)
			if err != nil {
				fmt.Println("There is a problem in createing DeviceInfo1 ", err)
				fmt.Println(err)
			}

			err = StartupProccessing(&reconciler, mgr)
			if err != nil {
				fmt.Println("Error in StartupProccessing")
			}
			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "deviceinfo-df-night01-wbfunction-decode-main",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println(err)
			}
			Expect(err).NotTo(HaveOccurred())
			var devinfocr examplecomv1.DeviceInfo
			_ = k8sClient.Get(ctx, client.ObjectKey{
				Name:      "deviceinfo-df-night01-wbfunction-decode-main",
				Namespace: "default",
			},
				&devinfocr)
			if devinfocr.Status.Response.Status == "Deployed" {
				fmt.Println("Change Undeploy Start")

				var deltime metav1.Time

				deltime = metav1.Now()

				devinfocr.DeletionTimestamp = &deltime
				err = deleteDeviceInfo(ctx, devinfocr)
				if err != nil {
					fmt.Println("There is a problem in updateing devinfocr ", err)
					fmt.Println(err)
				}

				_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "deviceinfo-df-night01-wbfunction-decode-main",
				}})
				if err != nil {
					By("Reconcile Error")
					fmt.Println(err)
				}
				err = createDeviceInfo(ctx, DeviceInfo3)
				if err != nil {
					fmt.Println("There is a problem in createing DeviceInfo3 ", err)
					fmt.Println(err)
				}
				_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "deviceinfo-df-night02-wbfunction-decode-main",
				}})
				if err != nil {
					By("Reconcile Error")
					fmt.Println(err)
				}
				var cpr examplecomv1.ComputeResource
				_ = k8sClient.Get(ctx, client.ObjectKey{
					Name:      "compute-test01",
					Namespace: "default",
				},
					&cpr)
				for _, val := range cpr.Spec.Regions {
					if val.DeviceType != "cpu" {
						continue
					}
					if len(val.Functions) == 2 {
						var devinfocr2 examplecomv1.DeviceInfo
						_ = k8sClient.Get(ctx, client.ObjectKey{
							Name:      "deviceinfo-df-night02-wbfunction-decode-main",
							Namespace: "default",
						},
							&devinfocr2)
						if devinfocr2.Status.Response.Status == "Deployed" {
							devinfocr2.DeletionTimestamp = &deltime
							err = deleteDeviceInfo(ctx, devinfocr2)
							if err != nil {
								fmt.Println("There is a problem in updateing devinfocr ", err)
								fmt.Println(err)
							}
							_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
								Namespace: "default",
								Name:      "deviceinfo-df-night02-wbfunction-decode-main",
							}})
							if err != nil {
								By("Reconcile Error")
								fmt.Println(err)
							}
							err = createDeviceInfo(ctx, DeviceInfo2)
							if err != nil {
								fmt.Println("There is a problem in createing DeviceInfo2 ", err)
								fmt.Println(err)
							}
							_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
								Namespace: "default",
								Name:      "deviceinfo-df-night01-wbfunction-decode-main",
							}})
							if err != nil {
								By("Reconcile Error")
								fmt.Println(err)
							}
							var devinfocr3 examplecomv1.DeviceInfo
							_ = k8sClient.Get(ctx, client.ObjectKey{
								Name:      "deviceinfo-df-night01-wbfunction-decode-main",
								Namespace: "default",
							},
								&devinfocr3)
							if devinfocr3.Status.Response.Status == "Undeployed" {
								var cpr examplecomv1.ComputeResource
								_ = k8sClient.Get(ctx, client.ObjectKey{
									Name:      "compute-test01",
									Namespace: "default",
								},
									&cpr)
								for _, val := range cpr.Spec.Regions {
									if val.DeviceType == "cpu" {
										if len(val.Functions) == 1 {
											Expect(val.Available).Should(Equal(chkComRes2.Available))
											Expect(val.CurrentCapacity).Should(Equal(chkComRes2.CurrentCapacity))
											Expect(val.CurrentFunctions).Should(Equal(chkComRes2.CurrentFunctions))
											Expect(val.DeviceFilePath).Should(Equal(chkComRes2.DeviceFilePath))
											Expect(val.DeviceIndex).Should(Equal(chkComRes2.DeviceIndex))
											Expect(val.DeviceType).Should(Equal(chkComRes2.DeviceType))
											Expect(val.DeviceUUID).Should(Equal(chkComRes2.DeviceUUID))
											Expect(val.MaxCapacity).Should(Equal(chkComRes2.MaxCapacity))
											Expect(val.MaxFunctions).Should(Equal(chkComRes2.MaxFunctions))
											Expect(val.Name).Should(Equal(chkComRes2.Name))
											Expect(val.Type).Should(Equal(chkComRes2.Type))
											Expect(val.Functions[0].Available).Should(Equal(chkComRes2.Functions[0].Available))
											Expect(val.Functions[0].CurrentCapacity).Should(Equal(chkComRes2.Functions[0].CurrentCapacity))
											Expect(val.Functions[0].CurrentDataFlows).Should(Equal(chkComRes2.Functions[0].CurrentDataFlows))
											Expect(val.Functions[0].FunctionIndex).Should(Equal(chkComRes2.Functions[0].FunctionIndex))
											Expect(val.Functions[0].FunctionName).Should(Equal(chkComRes2.Functions[0].FunctionName))
											Expect(val.Functions[0].MaxCapacity).Should(Equal(chkComRes2.Functions[0].MaxCapacity))
											Expect(val.Functions[0].MaxDataFlows).Should(Equal(chkComRes2.Functions[0].MaxDataFlows))
											Expect(val.Functions[0].PartitionName).Should(Equal(chkComRes2.Functions[0].PartitionName))
										} else {
											fmt.Println("Failed to delete Functions[0]")
										}
									} else {
										continue
									}
								}
							} else {
								fmt.Println("Status is not Undeploy")
							}
						} else {
							fmt.Println("DF Status is not Deploy 2")
						}
					} else {
						fmt.Println("Failed add 2 dataflows")
					}
				}
			}
		})

		AfterEach(func() {
			By("Test End")
			writer.Reset()
		})
	})
})

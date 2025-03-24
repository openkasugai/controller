/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED
*/

package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	ctrl "sigs.k8s.io/controller-runtime"

	examplecomv1 "FPGAFunction/api/v1"
	controllertestcpu "FPGAFunction/internal/controller/test/type/CPU"
	controllertestethernet "FPGAFunction/internal/controller/test/type/Ethernet"
	controllertestgpu "FPGAFunction/internal/controller/test/type/GPU"
	controllertestpcie "FPGAFunction/internal/controller/test/type/PCIe"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"

	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"

	// "k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// "k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func getMgr(mgr ctrl.Manager) (ctrl.Manager, error) {
	if mgr == nil {
		return ctrl.NewManager(cfg, ctrl.Options{
			Scheme: testScheme,
		})
	}
	return mgr, nil
}

// Create FPGAFunctionCR
func createFPGAFunction(ctx context.Context, fpgafcr examplecomv1.FPGAFunction) error {
	tmp := &examplecomv1.FPGAFunction{}
	*tmp = fpgafcr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Delete FPGAFunctionCR
func deleteFPGAFunction(ctx context.Context, fpgafcr examplecomv1.FPGAFunction) error {
	tmp := &examplecomv1.FPGAFunction{}
	*tmp = fpgafcr
	err := k8sClient.Delete(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create EthernetConnectionCR
func createEthernetConnection(ctx context.Context, ethercr controllertestethernet.EthernetConnection) error {
	tmp := &controllertestethernet.EthernetConnection{}
	*tmp = ethercr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create PCIeConnectionCR
func createPCIeConnection(ctx context.Context, pciecr controllertestpcie.PCIeConnection) error {
	tmp := &controllertestpcie.PCIeConnection{}
	*tmp = pciecr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create CPUFunctionCR
func createCPUFunction(ctx context.Context, cpufcr controllertestcpu.CPUFunction) error {
	tmp := &controllertestcpu.CPUFunction{}
	*tmp = cpufcr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create GPUFunctionCR
func createGPUFunction(ctx context.Context, gpufcr controllertestgpu.GPUFunction) error {
	tmp := &controllertestgpu.GPUFunction{}
	*tmp = gpufcr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create ChildBitstreamCR
func createChildBitstream(ctx context.Context, childbscr examplecomv1.ChildBs) error {
	tmp := &examplecomv1.ChildBs{}
	*tmp = childbscr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Update ChildBitstreamCR
func updateChildBitstream(ctx context.Context, childbscr examplecomv1.ChildBs) error {
	tmp := &examplecomv1.ChildBs{}
	*tmp = childbscr
	err := k8sClient.Update(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create FPGACR
func createFPGA(ctx context.Context, fpgacr examplecomv1.FPGA) error {
	tmp := &examplecomv1.FPGA{}
	*tmp = fpgacr
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Create ComfigMap
func config(ctx context.Context, conf corev1.ConfigMap) error {
	tmp := &corev1.ConfigMap{}
	*tmp = conf
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
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
	return nil
}

// Create FPGAReconfigurationCR
func createFPGAReconfiguration(ctx context.Context, fpgareconf examplecomv1.FPGAReconfiguration) error {
	tmp := &examplecomv1.FPGAReconfiguration{}
	*tmp = fpgareconf
	err := k8sClient.Create(ctx, tmp)
	if err != nil {
		return err
	}
	return nil
}

// Print ChildBitstreamCR
func printChildBitstream(ctx context.Context,
	childBsCR examplecomv1.ChildBs) {

	fmt.Println("ChildBsCR Spec")
	if nil != childBsCR.Spec.ChildBitstreamID {
		fmt.Println(" ChildBitstreamID=" + *childBsCR.Spec.ChildBitstreamID)
	}
	if nil != childBsCR.Spec.ChildBitstreamFile {
		fmt.Println(" ChildBitstreamFile=" + *childBsCR.Spec.ChildBitstreamFile)
	}
	for regionidx := 0; regionidx < len(childBsCR.Spec.Regions); regionidx++ {
		regioninfo := childBsCR.Status.Regions[regionidx]
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Name=" + *regioninfo.Name)
		if nil == regioninfo.MaxFunctions {
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".MaxFunctions=nil")
		} else {
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".MaxFunctions=" + strconv.Itoa(int(*regioninfo.MaxFunctions)))
		}
		if nil == regioninfo.MaxCapacity {
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".MaxCapacity=nil")
		} else {
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".MaxCapacity=" + strconv.Itoa(int(*regioninfo.MaxCapacity)))
		}
		if nil == regioninfo.Modules.Ptu {
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".Modules.Ptu=nil")
		} else {
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".Modules.Ptu.Cids=" + *regioninfo.Modules.Ptu.Cids)
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".Modules.Ptu.ID=" + strconv.Itoa(int(*regioninfo.Modules.Ptu.ID)))
			if nil == regioninfo.Modules.Ptu.Parameters {
				fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
					".Modules.Ptu.Parameters=nil")
			} else {
				for index, _ := range *regioninfo.Modules.Ptu.Parameters {
					fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
						".Modules.Ptu.Parameters[" + index +
						"].Type=" + strconv.Itoa(int((*regioninfo.Modules.Ptu.Parameters)[index].Type)))
					fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
						".Modules.Ptu.Parameters[" + index +
						"].IntVal=" + strconv.Itoa(int((*regioninfo.Modules.Ptu.Parameters)[index].IntVal)))
					fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
						".Modules.Ptu.Parameters[" + index +
						"].StrVal=" + (*regioninfo.Modules.Ptu.Parameters)[index].StrVal)
				}
			}
		}
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.LLDMA.Cids=" + *regioninfo.Modules.LLDMA.Cids)
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.LLDMA.ID=" + strconv.Itoa(int(*regioninfo.Modules.LLDMA.ID)))
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.Chain.ID=" + strconv.Itoa(int(*regioninfo.Modules.Chain.ID)))
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.Chain.Identifier=" + *regioninfo.Modules.Chain.Identifier)
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.Chain.Type=" + *regioninfo.Modules.Chain.Type)
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.Chain.Version=" + *regioninfo.Modules.Chain.Version)
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.Directtrans.ID=" + strconv.Itoa(int(*regioninfo.Modules.Directtrans.ID)))
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.Directtrans.Identifier=" + *regioninfo.Modules.Directtrans.Identifier)
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.Directtrans.Type=" + *regioninfo.Modules.Directtrans.Type)
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.Directtrans.Version=" + *regioninfo.Modules.Directtrans.Version)
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.Conversion.ID=" + strconv.Itoa(int(*regioninfo.Modules.Conversion.ID)))
		for convidx := 0; convidx < len(*regioninfo.Modules.Conversion.Module); convidx++ {
			convinfo := (*regioninfo.Modules.Conversion.Module)[convidx]
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".Modules.Conversion.Module[" + strconv.Itoa(convidx) +
				".Identifier=" + *convinfo.Identifier)
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".Modules.Conversion.Module[" + strconv.Itoa(convidx) +
				".Type=" + *convinfo.Type)
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".Modules.Conversion.Module[" + strconv.Itoa(convidx) +
				".Version=" + *convinfo.Version)
		}
		if nil == regioninfo.Modules.Functions {
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".Modules.Functions=nil")
		} else {
			for funcidx := 0; funcidx < len(*regioninfo.Modules.Functions); funcidx++ {
				funcinfo := (*regioninfo.Modules.Functions)[funcidx]
				fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
					".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
					".ID=" + strconv.Itoa(int(*funcinfo.ID)))
				if nil != funcinfo.FunctionName {
					fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
						".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
						".FunctionName=" + *funcinfo.FunctionName)
				} else {
					fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
						".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
						".FunctionName=nil")
				}
				for funcmdidx := 0; funcmdidx < len(*funcinfo.Module); funcmdidx++ {
					moduleinfo := (*funcinfo.Module)[funcmdidx]
					if nil != moduleinfo.FunctionChannelIDs {
						fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
							".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
							".Module[" + strconv.Itoa(funcmdidx) + "]" +
							".FunctionChannelIDs=" + *moduleinfo.FunctionChannelIDs)
					} else {
						fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
							".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
							".Module[" + strconv.Itoa(funcmdidx) + "]" +
							".FunctionChannelIDs=nil")
					}
					fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
						".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
						".Module[" + strconv.Itoa(funcmdidx) + "]" +
						".Identifier=" + *moduleinfo.Identifier)
					fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
						".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
						".Module[" + strconv.Itoa(funcmdidx) + "]" +
						".Type=" + *moduleinfo.Type)
					fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
						".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
						".Module[" + strconv.Itoa(funcmdidx) + "]" +
						".Version=" + *moduleinfo.Version)
				}
				if nil != funcinfo.Parameters {
					for mapidx, _ := range *funcinfo.Parameters {
						fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
							".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
							".Parameters[" + mapidx + "]" +
							".Type=" + strconv.Itoa(int((*funcinfo.Parameters)[mapidx].Type)))
						fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
							".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
							".Parameters[" + mapidx + "]" +
							".IntVal=" + strconv.Itoa(int((*funcinfo.Parameters)[mapidx].IntVal)))
						fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
							".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
							".Parameters[" + mapidx + "]" +
							".StrVal=" + (*funcinfo.Parameters)[mapidx].StrVal)
					}
				} else {
					fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
						".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
						".Parameters=nil")
				}
				for mapidx, param := range *funcinfo.IntraResourceMgmtMap {
					fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
						".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
						".IntraResourceMgmtMap[" + mapidx + "]" +
						".Available=" + strconv.FormatBool(*param.Available))
					if nil != param.FunctionCRName {
						fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
							".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
							".IntraResourceMgmtMap[" + mapidx + "]" +
							".FunctionCRName=" + *param.FunctionCRName)
					} else {
						fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
							".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
							".IntraResourceMgmtMap[" + mapidx + "]" +
							".FunctionCRName=nil")
					}
					if nil != param.Rx && nil != param.Rx.Protocol {
						for rxtxidx, rxtxparam := range *param.Rx.Protocol {
							fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
								".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
								".IntraResourceMgmtMap[" + mapidx + "]" +
								".Rx.Protocol[" + rxtxidx + "]" +
								".Port=" + strconv.Itoa(int(*rxtxparam.Port)))
							if nil != rxtxparam.DMAChannelID {
								fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
									".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
									".IntraResourceMgmtMap[" + mapidx + "]" +
									".Rx.Protocol[" + rxtxidx + "]" +
									".DMAChannelID=" + strconv.Itoa(int(*rxtxparam.DMAChannelID)))
							} else {
								fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
									".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
									".IntraResourceMgmtMap[" + mapidx + "]" +
									".Rx.Protocol[" + rxtxidx + "]" +
									".DMAChannelID=nil")
							}
							if nil != rxtxparam.LLDMAConnectorID {
								fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
									".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
									".IntraResourceMgmtMap[" + mapidx + "]" +
									".Rx.Protocol[" + rxtxidx + "]" +
									".LLDMAConnectorID=" + strconv.Itoa(int(*rxtxparam.LLDMAConnectorID)))
							} else {
								fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
									".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
									".IntraResourceMgmtMap[" + mapidx + "]" +
									".Rx.Protocol[" + rxtxidx + "]" +
									".LLDMAConnectorID=nil")
							}
						}
					}
					if nil != param.Tx && nil != param.Tx.Protocol {
						for rxtxidx, rxtxparam := range *param.Tx.Protocol {
							fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
								".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
								".IntraResourceMgmtMap[" + mapidx + "]" +
								".Tx.Protocol[" + rxtxidx + "]" +
								".Port=" + strconv.Itoa(int(*rxtxparam.Port)))
							if nil != rxtxparam.DMAChannelID {
								fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
									".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
									".IntraResourceMgmtMap[" + mapidx + "]" +
									".Tx.Protocol[" + rxtxidx + "]" +
									".DMAChannelID=" + strconv.Itoa(int(*rxtxparam.DMAChannelID)))
							} else {
								fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
									".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
									".IntraResourceMgmtMap[" + mapidx + "]" +
									".Rx.Protocol[" + rxtxidx + "]" +
									".DMAChannelID=nil")
							}
							if nil != rxtxparam.LLDMAConnectorID {
								fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
									".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
									".IntraResourceMgmtMap[" + mapidx + "]" +
									".Tx.Protocol[" + rxtxidx + "]" +
									".LLDMAConnectorID=" + strconv.Itoa(int(*rxtxparam.LLDMAConnectorID)))
							} else {
								fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
									".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
									".IntraResourceMgmtMap[" + mapidx + "]" +
									".Rx.Protocol[" + rxtxidx + "]" +
									".LLDMAConnectorID=nil")
							}
						}
					}
				}
				fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
					".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
					".DeploySpec.MaxCapacity=" + strconv.Itoa(int(*funcinfo.DeploySpec.MaxCapacity)))
				fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
					".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
					".DeploySpec.MaxDataFlows=" + strconv.Itoa(int(*funcinfo.DeploySpec.MaxDataFlows)))
			}
		}
	}
	fmt.Println("ChildBsCR Status")
	fmt.Println(" State=" + childBsCR.Status.State)
	fmt.Println(" Status=" + childBsCR.Status.Status)
	if nil != childBsCR.Status.ChildBitstreamID {
		fmt.Println(" ChildBitstreamID=" + *childBsCR.Status.ChildBitstreamID)
	}
	if nil != childBsCR.Status.ChildBitstreamFile {
		fmt.Println(" ChildBitstreamFile=" + *childBsCR.Status.ChildBitstreamFile)
	}
	for regionidx := 0; regionidx < len(childBsCR.Status.Regions); regionidx++ {
		regioninfo := childBsCR.Status.Regions[regionidx]
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Name=" + *regioninfo.Name)
		if nil == regioninfo.MaxFunctions {
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".MaxFunctions=nil")
		} else {
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".MaxFunctions=" + strconv.Itoa(int(*regioninfo.MaxFunctions)))
		}
		if nil == regioninfo.MaxCapacity {
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".MaxCapacity=nil")
		} else {
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".MaxCapacity=" + strconv.Itoa(int(*regioninfo.MaxCapacity)))
		}
		if nil == regioninfo.Modules.Ptu {
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".Modules.Ptu=nil")
		} else {
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".Modules.Ptu.Cids=" + *regioninfo.Modules.Ptu.Cids)
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".Modules.Ptu.ID=" + strconv.Itoa(int(*regioninfo.Modules.Ptu.ID)))
			if nil == regioninfo.Modules.Ptu.Parameters {
				fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
					".Modules.Ptu.Parameters=nil")
			} else {
				for index, _ := range *regioninfo.Modules.Ptu.Parameters {
					fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
						".Modules.Ptu.Parameters[" + index +
						"].Type=" + strconv.Itoa(int((*regioninfo.Modules.Ptu.Parameters)[index].Type)))
					fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
						".Modules.Ptu.Parameters[" + index +
						"].IntVal=" + strconv.Itoa(int((*regioninfo.Modules.Ptu.Parameters)[index].IntVal)))
					fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
						".Modules.Ptu.Parameters[" + index +
						"].StrVal=" + (*regioninfo.Modules.Ptu.Parameters)[index].StrVal)
				}
			}
		}
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.LLDMA.Cids=" + *regioninfo.Modules.LLDMA.Cids)
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.LLDMA.ID=" + strconv.Itoa(int(*regioninfo.Modules.LLDMA.ID)))
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.Chain.ID=" + strconv.Itoa(int(*regioninfo.Modules.Chain.ID)))
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.Chain.Identifier=" + *regioninfo.Modules.Chain.Identifier)
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.Chain.Type=" + *regioninfo.Modules.Chain.Type)
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.Chain.Version=" + *regioninfo.Modules.Chain.Version)
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.Directtrans.ID=" + strconv.Itoa(int(*regioninfo.Modules.Directtrans.ID)))
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.Directtrans.Identifier=" + *regioninfo.Modules.Directtrans.Identifier)
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.Directtrans.Type=" + *regioninfo.Modules.Directtrans.Type)
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.Directtrans.Version=" + *regioninfo.Modules.Directtrans.Version)
		fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
			".Modules.Conversion.ID=" + strconv.Itoa(int(*regioninfo.Modules.Conversion.ID)))
		for convidx := 0; convidx < len(*regioninfo.Modules.Conversion.Module); convidx++ {
			convinfo := (*regioninfo.Modules.Conversion.Module)[convidx]
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".Modules.Conversion.Module[" + strconv.Itoa(convidx) + "]" +
				".Identifier=" + *convinfo.Identifier)
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".Modules.Conversion.Module[" + strconv.Itoa(convidx) + "]" +
				".Type=" + *convinfo.Type)
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".Modules.Conversion.Module[" + strconv.Itoa(convidx) + "]" +
				".Version=" + *convinfo.Version)
		}
		if nil == regioninfo.Modules.Functions {
			fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
				".Modules.Functions=nil")
		} else {
			for funcidx := 0; funcidx < len(*regioninfo.Modules.Functions); funcidx++ {
				funcinfo := (*regioninfo.Modules.Functions)[funcidx]
				fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
					".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
					".ID=" + strconv.Itoa(int(*funcinfo.ID)))
				if nil != funcinfo.FunctionName {
					fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
						".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
						".FunctionName=" + *funcinfo.FunctionName)
				} else {
					fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
						".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
						".FunctionName=nil")
				}
				for funcmdidx := 0; funcmdidx < len(*funcinfo.Module); funcmdidx++ {
					moduleinfo := (*funcinfo.Module)[funcmdidx]
					if nil != moduleinfo.FunctionChannelIDs {
						fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
							".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
							".Module[" + strconv.Itoa(funcmdidx) + "]" +
							".FunctionChannelIDs=" + *moduleinfo.FunctionChannelIDs)
					} else {
						fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
							".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
							".Module[" + strconv.Itoa(funcmdidx) + "]" +
							".FunctionChannelIDs=nil")
					}
					fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
						".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
						".Module[" + strconv.Itoa(funcmdidx) + "]" +
						".Identifier=" + *moduleinfo.Identifier)
					fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
						".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
						".Module[" + strconv.Itoa(funcmdidx) + "]" +
						".Type=" + *moduleinfo.Type)
					fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
						".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
						".Module[" + strconv.Itoa(funcmdidx) + "]" +
						".Version=" + *moduleinfo.Version)
				}
				if nil != funcinfo.Parameters {
					for mapidx, _ := range *funcinfo.Parameters {
						fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
							".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
							".Parameters[" + mapidx + "]" +
							".Type=" + strconv.Itoa(int((*funcinfo.Parameters)[mapidx].Type)))
						fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
							".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
							".Parameters[" + mapidx + "]" +
							".IntVal=" + strconv.Itoa(int((*funcinfo.Parameters)[mapidx].IntVal)))
						fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
							".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
							".Parameters[" + mapidx + "]" +
							".StrVal=" + (*funcinfo.Parameters)[mapidx].StrVal)
					}
				} else {
					fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
						".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
						".Parameters=nil")
				}
				for mapidx, param := range *funcinfo.IntraResourceMgmtMap {
					fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
						".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
						".IntraResourceMgmtMap[" + mapidx + "]" +
						".Available=" + strconv.FormatBool(*param.Available))
					if nil != param.FunctionCRName {
						fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
							".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
							".IntraResourceMgmtMap[" + mapidx + "]" +
							".FunctionCRName=" + *param.FunctionCRName)
					} else {
						fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
							".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
							".IntraResourceMgmtMap[" + mapidx + "]" +
							".FunctionCRName=nil")
					}
					if nil != param.Rx && nil != param.Rx.Protocol {
						for rxtxidx, rxtxparam := range *param.Rx.Protocol {
							fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
								".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
								".IntraResourceMgmtMap[" + mapidx + "]" +
								".Rx.Protocol[" + rxtxidx + "]" +
								".Port=" + strconv.Itoa(int(*rxtxparam.Port)))
							if nil != rxtxparam.DMAChannelID {
								fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
									".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
									".IntraResourceMgmtMap[" + mapidx + "]" +
									".Rx.Protocol[" + rxtxidx + "]" +
									".DMAChannelID=" + strconv.Itoa(int(*rxtxparam.DMAChannelID)))
							} else {
								fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
									".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
									".IntraResourceMgmtMap[" + mapidx + "]" +
									".Rx.Protocol[" + rxtxidx + "]" +
									".DMAChannelID=nil")
							}
							if nil != rxtxparam.LLDMAConnectorID {
								fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
									".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
									".IntraResourceMgmtMap[" + mapidx + "]" +
									".Rx.Protocol[" + rxtxidx + "]" +
									".LLDMAConnectorID=" + strconv.Itoa(int(*rxtxparam.LLDMAConnectorID)))
							} else {
								fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
									".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
									".IntraResourceMgmtMap[" + mapidx + "]" +
									".Rx.Protocol[" + rxtxidx + "]" +
									".LLDMAConnectorID=nil")
							}
						}
					}
					if nil != param.Tx && nil != param.Tx.Protocol {
						for rxtxidx, rxtxparam := range *param.Tx.Protocol {
							fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
								".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
								".IntraResourceMgmtMap[" + mapidx + "]" +
								".Tx.Protocol[" + rxtxidx + "]" +
								".Port=" + strconv.Itoa(int(*rxtxparam.Port)))
							if nil != rxtxparam.DMAChannelID {
								fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
									".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
									".IntraResourceMgmtMap[" + mapidx + "]" +
									".Tx.Protocol[" + rxtxidx + "]" +
									".DMAChannelID=" + strconv.Itoa(int(*rxtxparam.DMAChannelID)))
							} else {
								fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
									".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
									".IntraResourceMgmtMap[" + mapidx + "]" +
									".Tx.Protocol[" + rxtxidx + "]" +
									".DMAChannelID=nil")
							}
							if nil != rxtxparam.LLDMAConnectorID {
								fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
									".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
									".IntraResourceMgmtMap[" + mapidx + "]" +
									".Tx.Protocol[" + rxtxidx + "]" +
									".LLDMAConnectorID=" + strconv.Itoa(int(*rxtxparam.LLDMAConnectorID)))
							} else {
								fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
									".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
									".IntraResourceMgmtMap[" + mapidx + "]" +
									".Tx.Protocol[" + rxtxidx + "]" +
									".LLDMAConnectorID=nil")
							}
						}
					}
				}
				fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
					".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
					".DeploySpec.MaxCapacity=" + strconv.Itoa(int(*funcinfo.DeploySpec.MaxCapacity)))
				fmt.Println(" Regions[" + strconv.Itoa(regionidx) + "]" +
					".Modules.Functions[" + strconv.Itoa(funcidx) + "]" +
					".DeploySpec.MaxDataFlows=" + strconv.Itoa(int(*funcinfo.DeploySpec.MaxDataFlows)))
			}
		}
	}
}

// Print FPGACR
func printFPGA(ctx context.Context,
	fpgaCR examplecomv1.FPGA) {
	fmt.Println("fpgaCR Spec")
	if nil != fpgaCR.Spec.ChildBitstreamID {
		fmt.Println(" ChildBitstreamID=" + *fpgaCR.Spec.ChildBitstreamID)
	} else {
		fmt.Println(" ChildBitstreamID=nil")
	}
	fmt.Println("fpgaCR Status")
	if nil != fpgaCR.Status.ChildBitstreamID {
		fmt.Println(" ChildBitstreamID=" + *fpgaCR.Status.ChildBitstreamID)
	} else {
		fmt.Println(" ChildBitstreamID=nil")
	}
	if nil != fpgaCR.Status.ChildBitstreamCRName {
		fmt.Println(" ChildBitstreamCRName=" + *fpgaCR.Status.ChildBitstreamCRName)
	} else {
		fmt.Println(" ChildBitstreamCRName=nil")
	}
}

// Print FPGAFunctionCR
func printFPGAFunction(ctx context.Context,
	fpgaFunctionCR examplecomv1.FPGAFunction) {
	fmt.Println("fpgaFunctionCR Status")
	fmt.Println(" Status=" + fpgaFunctionCR.Status.Status)
	fmt.Println(" FunctionIndex=" + strconv.Itoa(int(fpgaFunctionCR.Status.FunctionIndex)))
	fmt.Println(" Rx.Protocol=" + fpgaFunctionCR.Status.Rx.Protocol)
	if nil != fpgaFunctionCR.Status.Rx.IPAddress {
		fmt.Println(" Rx.IPAddress=" + *fpgaFunctionCR.Status.Rx.IPAddress)
	} else {
		fmt.Println(" Rx.IPAddress=nil")
	}
	if nil != fpgaFunctionCR.Status.Rx.Port {
		fmt.Println(" Rx.Port=" + strconv.Itoa(int(*fpgaFunctionCR.Status.Rx.Port)))
	} else {
		fmt.Println(" Rx.Port=nil")
	}
	if nil != fpgaFunctionCR.Status.Rx.SubnetAddress {
		fmt.Println(" Rx.SubnetAddress=" + *fpgaFunctionCR.Status.Rx.SubnetAddress)
	} else {
		fmt.Println(" Rx.SubnetAddress=nil")
	}
	if nil != fpgaFunctionCR.Status.Rx.GatewayAddress {
		fmt.Println(" Rx.GatewayAddress=" + *fpgaFunctionCR.Status.Rx.GatewayAddress)
	} else {
		fmt.Println(" Rx.GatewayAddress=nil")
	}
	if nil != fpgaFunctionCR.Status.Rx.DMAChannelID {
		fmt.Println(" Rx.DMAChannelID=" + strconv.Itoa(int(*fpgaFunctionCR.Status.Rx.DMAChannelID)))
	} else {
		fmt.Println(" Rx.DMAChannelID=nil")
	}
	if nil != fpgaFunctionCR.Status.Rx.LLDMAConnectorID {
		fmt.Println(" Rx.LLDMAConnectorID=" + strconv.Itoa(int(*fpgaFunctionCR.Status.Rx.LLDMAConnectorID)))
	} else {
		fmt.Println(" Rx.LLDMAConnectorID=nil")
	}
	fmt.Println(" Tx.Protocol=" + fpgaFunctionCR.Status.Rx.Protocol)
	if nil != fpgaFunctionCR.Status.Rx.IPAddress {
		fmt.Println(" Tx.IPAddress=" + *fpgaFunctionCR.Status.Rx.IPAddress)
	} else {
		fmt.Println(" Tx.IPAddress=nil")
	}
	if nil != fpgaFunctionCR.Status.Rx.Port {
		fmt.Println(" Tx.Port=" + strconv.Itoa(int(*fpgaFunctionCR.Status.Rx.Port)))
	} else {
		fmt.Println(" Tx.Port=nil")
	}
	if nil != fpgaFunctionCR.Status.Rx.SubnetAddress {
		fmt.Println(" Tx.SubnetAddress=" + *fpgaFunctionCR.Status.Rx.SubnetAddress)
	} else {
		fmt.Println(" Tx.SubnetAddress=nil")
	}
	if nil != fpgaFunctionCR.Status.Rx.GatewayAddress {
		fmt.Println(" Tx.GatewayAddress=" + *fpgaFunctionCR.Status.Rx.GatewayAddress)
	} else {
		fmt.Println(" Tx.GatewayAddress=nil")
	}
	if nil != fpgaFunctionCR.Status.Rx.DMAChannelID {
		fmt.Println(" Tx.DMAChannelID=" + strconv.Itoa(int(*fpgaFunctionCR.Status.Rx.DMAChannelID)))
	} else {
		fmt.Println(" Tx.DMAChannelID=nil")
	}
	if nil != fpgaFunctionCR.Status.Rx.LLDMAConnectorID {
		fmt.Println(" Tx.LLDMAConnectorID=" + strconv.Itoa(int(*fpgaFunctionCR.Status.Rx.LLDMAConnectorID)))
	} else {
		fmt.Println(" Tx.LLDMAConnectorID=nil")
	}
}

// Print FPGAReconfigurationCR
func printFPGAReconfiguration(ctx context.Context,
	fpgaReconfiguratonCR examplecomv1.FPGAReconfiguration) {
	fmt.Println("fpgaReconfiguratonCR Status")
	fmt.Println(" Status=" + fpgaReconfiguratonCR.Status.Status)
}

// Print DeployInfoCM
func printDeployInfo(ctx context.Context) {
	var err error
	var deployInfoCM map[string][]examplecomv1.DeviceRegionInfo
	// Get DeployInfo CM
	tmpData := &unstructured.Unstructured{}
	tmpData.SetGroupVersionKind(schema.GroupVersionKind{
		Kind:    "ConfigMap",
		Version: "v1",
	})
	err = k8sClient.Get(ctx, types.NamespacedName{
		Namespace: "default",
		Name:      "deployinfo",
	}, tmpData)
	if err != nil {
		fmt.Println("Get deployinfo error.", err)
	} else {
		mapdata, _, _ := unstructured.NestedStringMap(tmpData.Object, "data")
		var configData []byte
		for _, jsonrecord := range mapdata {
			configData = []byte(jsonrecord)
		}
		err = json.Unmarshal(configData, &deployInfoCM)
		if err != nil {
			fmt.Println("Unmarshal deployinfo error.", err)
		}
	}
	if nil == err {
		deployDevices := deployInfoCM["devices"]
		fmt.Println("deployinfoCM")
		for mapidx := 0; mapidx < len(deployDevices); mapidx++ {
			deployInfo := deployDevices[mapidx]
			fmt.Println(" deployinfoCM[" + strconv.Itoa(mapidx) + "]" +
				".NodeName=" + deployInfo.NodeName)
			if nil != deployInfo.DeviceFilePath {
				fmt.Println(" deployinfoCM[" + strconv.Itoa(mapidx) + "]" +
					".DeviceFilePath=" + *deployInfo.DeviceFilePath)
			} else {
				fmt.Println(" deployinfoCM[" + strconv.Itoa(mapidx) + "]" +
					".DeviceFilePath=nil")
			}
			if nil != deployInfo.DeviceFilePath {
				fmt.Println(" deployinfoCM[" + strconv.Itoa(mapidx) + "]" +
					".DeviceUUID=" + *deployInfo.DeviceUUID)
			} else {
				fmt.Println(" deployinfoCM[" + strconv.Itoa(mapidx) + "]" +
					".DeviceUUID=nil")
			}
			for targetidx := 0; targetidx < len(deployInfo.FunctionTargets); targetidx++ {
				functg := deployInfo.FunctionTargets[targetidx]
				fmt.Println(" deployinfoCM[" + strconv.Itoa(mapidx) + "]" +
					".FunctionTargets[" + strconv.Itoa(targetidx) + "]" +
					".RegionType=" + functg.RegionType)
				fmt.Println(" deployinfoCM[" + strconv.Itoa(mapidx) + "]" +
					".FunctionTargets[" + strconv.Itoa(targetidx) + "]" +
					".RegionName=" + functg.RegionName)
				fmt.Println(" deployinfoCM[" + strconv.Itoa(mapidx) + "]" +
					".FunctionTargets[" + strconv.Itoa(targetidx) + "]" +
					".MaxFunctions=" + strconv.Itoa(int(functg.MaxFunctions)))
				fmt.Println(" deployinfoCM[" + strconv.Itoa(mapidx) + "]" +
					".FunctionTargets[" + strconv.Itoa(targetidx) + "]" +
					".MaxCapacity=" + strconv.Itoa(int(functg.MaxCapacity)))
				for funcidx := 0; funcidx < len(functg.Functions); funcidx++ {
					funcs := functg.Functions[funcidx]
					if nil != funcs.FunctionIndex {
						fmt.Println(" deployinfoCM[" + strconv.Itoa(mapidx) + "]" +
							".FunctionTargets[" + strconv.Itoa(targetidx) + "]" +
							".Functions[" + strconv.Itoa(funcidx) + "]" +
							".FunctionIndex=" + strconv.Itoa(int(*funcs.FunctionIndex)))
					} else {
						fmt.Println(" deployinfoCM[" + strconv.Itoa(mapidx) + "]" +
							".FunctionTargets[" + strconv.Itoa(targetidx) + "]" +
							".Functions[" + strconv.Itoa(funcidx) + "]" +
							".FunctionIndex=nil")
					}
					fmt.Println(" deployinfoCM[" + strconv.Itoa(mapidx) + "]" +
						".FunctionTargets[" + strconv.Itoa(targetidx) + "]" +
						".Functions[" + strconv.Itoa(funcidx) + "]" +
						".PartitionName=" + funcs.PartitionName)
					fmt.Println(" deployinfoCM[" + strconv.Itoa(mapidx) + "]" +
						".FunctionTargets[" + strconv.Itoa(targetidx) + "]" +
						".Functions[" + strconv.Itoa(funcidx) + "]" +
						".FunctionName=" + funcs.FunctionName)
					fmt.Println(" deployinfoCM[" + strconv.Itoa(mapidx) + "]" +
						".FunctionTargets[" + strconv.Itoa(targetidx) + "]" +
						".Functions[" + strconv.Itoa(funcidx) + "]" +
						".MaxDataFlows=" + strconv.Itoa(int(funcs.MaxDataFlows)))
					fmt.Println(" deployinfoCM[" + strconv.Itoa(mapidx) + "]" +
						".FunctionTargets[" + strconv.Itoa(targetidx) + "]" +
						".Functions[" + strconv.Itoa(funcidx) + "]" +
						".MaxCapacity=" + strconv.Itoa(int(funcs.MaxCapacity)))
				}
			}
		}
	}
}

var _ = Describe("FPGAFunctionController", func() {
	var mgr ctrl.Manager
	var err error
	ctx := context.Background()

	var fakerecorder = record.NewFakeRecorder(10)
	var writer = bytes.Buffer{}

	Context("Test for FPGAFunctionReconciler", Ordered, func() {
		var reconciler FPGAFunctionReconciler
		BeforeAll(func() {

			opts := zap.Options{
				TimeEncoder: zapcore.ISO8601TimeEncoder,
				Development: true,
				DestWriter:  &writer,
			}
			ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
		})

		BeforeEach(func() {
			writer.Reset()

			// set manager
			mgr, err = getMgr(mgr)
			Expect(err).NotTo(HaveOccurred())

			// recorder initialized
			fakerecorder = record.NewFakeRecorder(10)
			reconciler = FPGAFunctionReconciler{
				Client:   k8sClient,
				Scheme:   testScheme,
				Recorder: fakerecorder,
			}
			err = reconciler.SetupWithManager(mgr)
			if err != nil {
				fmt.Println("Error in SetupWithManager")
			}
			Expect(err).NotTo(HaveOccurred())
			os.Setenv("K8S_NODENAME", "test01")
			os.Setenv("K8S_CLUSTERNAME", "default")
			StartupProccessing(mgr)
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestethernet.EthernetConnection{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestpcie.PCIeConnection{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestgpu.GPUFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestcpu.CPUFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGAFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGA{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.ChildBs{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.ComputeResource{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGAReconfiguration{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
		})
		It("FPGAFunctionTest 3", func() {
			By("Test Start")
			fmt.Println("test3:--------------------------------------------------")

			// Create ConfigMap
			err = config(ctx, fpgafuncconfig_fr_high_infer_2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, servicerMgmtConfig2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, deployinfo_configdata2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, regionUniqueInfoConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, functionUniqueInfoConfig2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, preDeterminedRegionInfoConfig8)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, fr_childbs_Config2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection3)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR : ", err)
			}

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection4)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR : ", err)
			}

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction2)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR : ", err)
			}

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction2)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR : ", err)
			}

			// Create ChildBitstreamCR
			err = createChildBitstream(ctx, ChildBitstream1)
			if err != nil {
				fmt.Println("There is a problem in createing ChildBitstream ", err)
			}

			// Create FPGACR
			err = createFPGA(ctx, FPGA2[0])
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR ", err)
			}

			// Create ComputeResourceCR
			err = createCompureResource(ctx, comres1)
			if err != nil {
				fmt.Println("There is a problem in createing ComputeresourceCR ", err)
			}

			// Create FPGAFunctionCR
			err = createFPGAFunction(ctx, FPGAFunction2)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}

			var childBsCR examplecomv1.ChildBs
			// Get ChildBs CR
			_ = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01-00000002",
			},
				&childBsCR)
			if err != nil {
				fmt.Println("get childBsCR error.", err)
			} else {
				printChildBitstream(ctx, childBsCR)
			}

			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night02-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			// Get ChildBs CR
			_ = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01-00000002",
			},
				&childBsCR)
			if err != nil {
				fmt.Println("get childBsCR error.", err)
			} else {
				printChildBitstream(ctx, childBsCR)
			}

			var fpgaCR examplecomv1.FPGA
			// Get fpga CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01",
			},
				&fpgaCR)
			if err != nil {
				fmt.Println("Get fpgaCR error.", err)
			} else {
				printFPGA(ctx, fpgaCR)
			}

			var fpgafuncCR examplecomv1.FPGAFunction
			// Get fpga CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "df-night02-wbfunction-filter-resize-high-infer-main",
			},
				&fpgafuncCR)
			if err != nil {
				fmt.Println("get fpgafuncCR error.", err)
			} else {
				printFPGAFunction(ctx, fpgafuncCR)
				Expect(strconv.Itoa(int(fpgafuncCR.Status.FunctionIndex))).Should(Equal("1"))
			}

			fmt.Println("test3:^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
			fmt.Println(writer.String())
			Expect(err).NotTo(HaveOccurred())
		})
		It("FPGAFunctionTest 4", func() {
			By("Test Start")
			fmt.Println("test4:--------------------------------------------------")

			// Create ConfigMap
			err = config(ctx, fpgafuncconfig_fr_high_infer_3)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, servicerMgmtConfig2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, deployinfo_configdata2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, regionUniqueInfoConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, functionUniqueInfoConfig2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, preDeterminedRegionInfoConfig8)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, fr_childbs_Config2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection3)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR : ", err)
			}

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection4)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR : ", err)
			}

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction3)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR : ", err)
			}

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction3)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR : ", err)
			}
			// Create FPGACR
			err = createFPGA(ctx, FPGA3[0])
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR ", err)
			}

			// Create ComputeResourceCR
			err = createCompureResource(ctx, comres1)
			if err != nil {
				fmt.Println("There is a problem in createing ComputeresourceCR ", err)
			}

			// Create FPGAFunctionCR
			err = createFPGAFunction(ctx, FPGAFunction3)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}

			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night03-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			var childBsCR examplecomv1.ChildBs
			// Get ChildBs CR
			_ = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01-00000001",
			},
				&childBsCR)

			printChildBitstream(ctx, childBsCR)
			printDeployInfo(ctx)

			var fpgaCR examplecomv1.FPGA
			// Get fpga CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01",
			},
				&fpgaCR)
			if err != nil {
				fmt.Println("Get fpgaCR error.", err)
			} else {
				printFPGA(ctx, fpgaCR)
			}

			var fpgafuncCR examplecomv1.FPGAFunction
			// Get fpga CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "df-night03-wbfunction-filter-resize-high-infer-main",
			},
				&fpgafuncCR)
			if err != nil {
				fmt.Println("get fpgafuncCR error.", err)
			} else {
				printFPGAFunction(ctx, fpgafuncCR)
				Expect(fpgafuncCR.Status.FunctionIndex).Should(BeZero())
			}
			controllerutil.RemoveFinalizer(&fpgafuncCR, "fpgafunction.finalizers.example.com.v1")
			// Update RemoveFinalizer to fpgafunc CR
			err = k8sClient.Update(ctx, &fpgafuncCR)
			if err != nil {
				fmt.Println("error update RemoveFinalizer to fpgafuncCR.", err)
			}

			fmt.Println("test4:^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
			fmt.Println(writer.String())
			Expect(err).NotTo(HaveOccurred())
		})
		It("FPGAFunctionTest 5", func() {
			By("Test Start")
			fmt.Println("test5:--------------------------------------------------")

			// Create ConfigMap
			err = config(ctx, fpgafuncconfig_fr_high_infer_4)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, servicerMgmtConfig2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, deployinfo_configdata2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, regionUniqueInfoConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, functionUniqueInfoConfig2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, preDeterminedRegionInfoConfig8)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, fr_childbs_Config1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection5)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR : ", err)
			}

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection6)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR : ", err)
			}

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction4)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR : ", err)
			}

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction4)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR : ", err)
			}

			// Create ChildBitstreamCR
			err = createChildBitstream(ctx, ChildBitstream2)
			if err != nil {
				fmt.Println("There is a problem in createing ChildBitstream ", err)
			}
			err = createFPGA(ctx, FPGA4[0])
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR ", err)
			}
			err = createFPGAFunction(ctx, FPGAFunction4)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}

			var childBsCR examplecomv1.ChildBs
			// Get ChildBs CR
			_ = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01-4444444444",
			},
				&childBsCR)
			if err != nil {
				fmt.Println("Get childBsCR error.", err)
			} else {
				printChildBitstream(ctx, childBsCR)
			}

			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night04-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			// Get ChildBs CR
			_ = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01-4444444444",
			},
				&childBsCR)
			if err != nil {
				fmt.Println("Get childBsCR error.", err)
			} else {
				printChildBitstream(ctx, childBsCR)
			}

			var fpgaCR examplecomv1.FPGA
			// Get fpga CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01",
			},
				&fpgaCR)
			if err != nil {
				fmt.Println("Get fpgaCR error.", err)
			} else {
				printFPGA(ctx, fpgaCR)
			}

			var fpgafuncCR examplecomv1.FPGAFunction
			// Get fpga CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "df-night04-wbfunction-filter-resize-high-infer-main",
			},
				&fpgafuncCR)
			if err != nil {
				fmt.Println("get fpgafuncCR error.", err)
			} else {
				printFPGAFunction(ctx, fpgafuncCR)
				Expect(strconv.Itoa(int(fpgafuncCR.Status.FunctionIndex))).Should(Equal("3"))
			}

			fmt.Println("test5:^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
			fmt.Println(writer.String())
			Expect(err).NotTo(HaveOccurred())
		})
		It("FPGAFunctionTest 6", func() {
			By("Test Start")
			fmt.Println("test6:--------------------------------------------------")
			// Create ConfigMap
			err = config(ctx, fpgafuncconfig_fr_high_infer_2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, servicerMgmtConfig2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, deployinfo_configdata2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, regionUniqueInfoConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, functionUniqueInfoConfig2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, preDeterminedRegionInfoConfig8)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, fr_childbs_Config1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection7)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR : ", err)
			}

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection8)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR : ", err)
			}

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection5)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR : ", err)
			}

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection6)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR : ", err)
			}

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction2)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR : ", err)
			}

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction2)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR : ", err)
			}

			// Create ChildBitstreamCR
			err = createChildBitstream(ctx, ChildBitstream1)
			if err != nil {
				fmt.Println("There is a problem in createing ChildBitstream ", err)
			}

			// Create FPGACR
			err = createFPGA(ctx, FPGA2[0])
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR ", err)
			}

			// Create FPGAFunctionCR
			err = createFPGAFunction(ctx, FPGAFunction5)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR  A", err)
			}

			_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night05-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}
			var childBsCR examplecomv1.ChildBs
			// Get ChildBs CR
			_ = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01-00000002",
			},
				&childBsCR)
			if err != nil {
				fmt.Println("Get childBsCR error.", err)
			} else {
				printChildBitstream(ctx, childBsCR)
			}

			var fpgaCR examplecomv1.FPGA
			// Get fpga CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01",
			},
				&fpgaCR)
			if err != nil {
				fmt.Println("Get fpgaCR error.", err)
			} else {
				printFPGA(ctx, fpgaCR)
			}

			var fpgafuncCR examplecomv1.FPGAFunction
			// Get fpga CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "df-night05-wbfunction-filter-resize-high-infer-main",
			},
				&fpgafuncCR)
			if err != nil {
				fmt.Println("get fpgafuncCR error.", err)
			} else {
				printFPGAFunction(ctx, fpgafuncCR)
			}

			err = createFPGAFunction(ctx, FPGAFunction6)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}

			_, err = reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night06-wbfunction-filter-resize-high-infer-main",
			}})

			// Get ChildBs CR
			_ = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01-00000002",
			},
				&childBsCR)
			if err != nil {
				fmt.Println("Get childBsCR error.", err)
			} else {
				printChildBitstream(ctx, childBsCR)
			}

			// Get fpga CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01",
			},
				&fpgaCR)
			if err != nil {
				fmt.Println("Get fpgaCR error.", err)
			} else {
				printFPGA(ctx, fpgaCR)
			}

			// Get fpga CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "df-night06-wbfunction-filter-resize-high-infer-main",
			},
				&fpgafuncCR)
			if err != nil {
				fmt.Println("get fpgafuncCR error.", err)
			} else {
				printFPGAFunction(ctx, fpgafuncCR)
			}

			fmt.Println("test6:^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
			fmt.Println(writer.String())
			Expect(err).NotTo(HaveOccurred())
		})

		It("FPGAFunctionTest 7-UPDATE", func() {
			By("Test Start")
			// Create FPGAFunctionCR
			err = createFPGAFunction(ctx, fpgafunctionUPDATE)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR", err)
			}
			Expect(err).NotTo(HaveOccurred())

			got, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: TESTNAMESPACE,
				Name:      "fpgafunctestupdate-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}

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

		AfterEach(func() {
			By("Test End")
			writer.Reset()
		})
		AfterAll(func() {
			By("Test End")
			writer.Reset()
		})
	})

	var fakerecorder2 = record.NewFakeRecorder(10)
	var writer2 = bytes.Buffer{}

	Context("Test for FPGAFunctionReconciler2", Ordered, func() {
		var reconciler2 FPGAFunctionReconciler
		BeforeAll(func() {

			opts := zap.Options{
				TimeEncoder: zapcore.ISO8601TimeEncoder,
				Development: true,
				DestWriter:  &writer2,
			}
			ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
		})

		BeforeEach(func() {
			writer2.Reset()

			// set manager
			mgr, err = getMgr(mgr)
			Expect(err).NotTo(HaveOccurred())

			reconciler2 = FPGAFunctionReconciler{
				Client:   k8sClient,
				Scheme:   testScheme,
				Recorder: fakerecorder2,
			}
			err = reconciler2.SetupWithManager(mgr)
			if err != nil {
				fmt.Println("Error in SetupWithManager")
			}
			Expect(err).NotTo(HaveOccurred())
			os.Setenv("K8S_NODENAME", "test01")
			os.Setenv("K8S_CLUSTERNAME", "default")
			StartupProccessing(mgr)
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestethernet.EthernetConnection{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestpcie.PCIeConnection{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestgpu.GPUFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestcpu.CPUFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGAFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGA{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.ChildBs{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.ComputeResource{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGAReconfiguration{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
		})

		It("FPGAFunctionTest 2", func() {
			By("Test Start")
			fmt.Println("test2:--------------------------------------------------")

			// Create ConfigMap
			err = config(ctx, fpgafuncconfig_fr_high_infer_1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, servicerMgmtConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, deployinfo_configdata1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, regionUniqueInfoConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, functionUniqueInfoConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, preDeterminedRegionInfoConfig8)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, fr_childbs_Config1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection1)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR : ", err)
			}

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection2)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR : ", err)
			}

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction1)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR : ", err)
			}

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction1)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR : ", err)
			}

			// Create ChildBitstreamCR
			err = createChildBitstream(ctx, ChildBitstream1)
			if err != nil {
				fmt.Println("There is a problem in createing ChildBitstream ", err)
			}

			// Create FPGACR
			err = createFPGA(ctx, FPGA1[0])
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR ", err)
			}

			// Create ComputeResourceCR
			err = createCompureResource(ctx, comres1)
			if err != nil {
				fmt.Println("There is a problem in createing ComputeresourceCR ", err)
			}

			// Create FPGAFunctionCR
			err = createFPGAFunction(ctx, FPGAFunction1)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}

			printDeployInfo(ctx)

			_, err := reconciler2.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			var childBsCR examplecomv1.ChildBs
			var childBitstreamID string
			// Get ChildBs CR
			_ = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01-00000001",
			},
				&childBsCR)
			if nil != childBsCR.Spec.ChildBitstreamID {
				childBitstreamID = *childBsCR.Spec.ChildBitstreamID
			}
			printChildBitstream(ctx, childBsCR)
			printDeployInfo(ctx)

			var fpgaCR examplecomv1.FPGA
			// Get fpga CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01",
			},
				&fpgaCR)
			if err != nil {
				fmt.Println("Get fpgaCR error.", err)
			} else {
				printFPGA(ctx, fpgaCR)
				Expect(fpgaCR.Status.ChildBitstreamID).ShouldNot(BeNil())
				Expect(fpgaCR.Status.ChildBitstreamCRName).ShouldNot(BeNil())
				Expect(*fpgaCR.Status.ChildBitstreamID).To(Equal(childBitstreamID))
				Expect(*fpgaCR.Status.ChildBitstreamCRName).To(Equal(fpgaCR.Name + "-" + childBitstreamID))
			}

			var fpgafuncCR examplecomv1.FPGAFunction
			// Get fpga CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
			},
				&fpgafuncCR)
			if err != nil {
				fmt.Println("get fpgafuncCR error.", err)
			} else {
				printFPGAFunction(ctx, fpgafuncCR)
			}

			childBsCR.Status.State = examplecomv1.ChildBsReady

			// update ChildBs CR
			err = updateChildBitstream(ctx, childBsCR)
			if err != nil {
				fmt.Println("There is a problem in updateing ChildBsCR ", err)
			}
			printChildBitstream(ctx, childBsCR)
			printDeployInfo(ctx)

			_, err = reconciler2.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}
			// Get ChildBs CR
			_ = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01-00000001",
			},
				&childBsCR)
			printChildBitstream(ctx, childBsCR)
			printDeployInfo(ctx)

			// Get fpga CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
			},
				&fpgafuncCR)
			if err != nil {
				fmt.Println("get fpgafuncCR error.", err)
			} else {
				printFPGAFunction(ctx, fpgafuncCR)
				Expect(fpgafuncCR.Status.FunctionIndex).Should(BeZero())
				Expect(fpgafuncCR.Status.Status).To(Equal("Running"))

				err = k8sClient.Delete(ctx, &fpgafuncCR)
				if err != nil {
					fmt.Println("get fpgafuncCR error.", err)
				} else {
					fmt.Println("delete fpgafuncCR")
				}
			}

			_, err = reconciler2.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			// Delete ChildBs CR
			err = k8sClient.Delete(ctx, &childBsCR)
			if err != nil {
				fmt.Println("Delete childBsCR error.", err)
			}

			fmt.Println("test2:^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
			fmt.Println(writer2.String())

			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			By("Test End")
			writer2.Reset()
		})
		AfterAll(func() {
			By("Test End")
			writer2.Reset()
		})
	})

	var fakerecorder3 = record.NewFakeRecorder(10)
	var writer3 = bytes.Buffer{}

	Context("Test for FPGAFunctionReconciler3", Ordered, func() {
		var reconciler3 FPGAFunctionReconciler
		BeforeAll(func() {

			opts := zap.Options{
				TimeEncoder: zapcore.ISO8601TimeEncoder,
				Development: true,
				DestWriter:  &writer3,
			}
			ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
		})

		BeforeEach(func() {
			writer3.Reset()

			// set manager
			mgr, err = getMgr(mgr)
			Expect(err).NotTo(HaveOccurred())

			reconciler3 = FPGAFunctionReconciler{
				Client:   k8sClient,
				Scheme:   testScheme,
				Recorder: fakerecorder3,
			}
			err = reconciler3.SetupWithManager(mgr)
			if err != nil {
				fmt.Println("Error in SetupWithManager")
			}
			Expect(err).NotTo(HaveOccurred())
			os.Setenv("K8S_NODENAME", "test01")
			os.Setenv("K8S_CLUSTERNAME", "default")
			StartupProccessing(mgr)
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestethernet.EthernetConnection{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestpcie.PCIeConnection{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestgpu.GPUFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestcpu.CPUFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGAFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGA{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.ChildBs{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.ComputeResource{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGAReconfiguration{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
		})
		It("FPGAFunctionTest 7", func() {
			By("Test Start")
			fmt.Println("test7:--------------------------------------------------")

			// Create ConfigMap
			err = config(ctx, fpgafuncconfig_fr_high_infer_1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, servicerMgmtConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, deployinfo_configdata1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, regionUniqueInfoConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, functionUniqueInfoConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, preDeterminedRegionInfoConfig8)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, fr_childbs_Config1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection1)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR : ", err)
			}

			// Create EthernetConnectionCR
			err = createEthernetConnection(ctx, EthernetConnection2)
			if err != nil {
				fmt.Println("There is a problem in createing EthernetConnectionCR : ", err)
			}

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction1)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR : ", err)
			}

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction1)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR : ", err)
			}

			// Create ChildBitstreamCR
			err = createChildBitstream(ctx, ChildBitstream1)
			if err != nil {
				fmt.Println("There is a problem in createing ChildBitstream ", err)
			}

			// Create FPGACR
			err = createFPGA(ctx, FPGA1[0])
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR ", err)
			}

			// Create ComputeResourceCR
			err = createCompureResource(ctx, comres1)
			if err != nil {
				fmt.Println("There is a problem in createing ComputeresourceCR ", err)
			}

			// Create FPGAFunctionCR
			err = createFPGAFunction(ctx, FPGAFunction1)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}

			printDeployInfo(ctx)

			_, err := reconciler3.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			var childBsCR examplecomv1.ChildBs
			// Get ChildBs CR
			_ = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01-00000001",
			},
				&childBsCR)

			printChildBitstream(ctx, childBsCR)
			printDeployInfo(ctx)

			_, err = reconciler3.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			// Get ChildBs CR
			_ = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01-00000001",
			},
				&childBsCR)

			printChildBitstream(ctx, childBsCR)

			childBsCR.Status.State = examplecomv1.ChildBsConfiguringNetwork

			// update ChildBs CR
			err = updateChildBitstream(ctx, childBsCR)
			if err != nil {
				fmt.Println("There is a problem in updateing ChildBsCR ", err)
			}

			printChildBitstream(ctx, childBsCR)
			printDeployInfo(ctx)

			_, err = reconciler3.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			// Get ChildBs CR
			_ = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01-00000001",
			},
				&childBsCR)

			printChildBitstream(ctx, childBsCR)

			childBsCR.Status.State = examplecomv1.ChildBsReady

			// update ChildBs CR
			err = updateChildBitstream(ctx, childBsCR)
			if err != nil {
				fmt.Println("There is a problem in updateing ChildBsCR ", err)
			}

			printChildBitstream(ctx, childBsCR)
			printDeployInfo(ctx)

			_, err = reconciler3.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}
			// Get ChildBs CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01-00000001",
			},
				&childBsCR)
			if err != nil {
				fmt.Println("get ChildBsCR error", err)
			} else {
				printChildBitstream(ctx, childBsCR)
			}
			printDeployInfo(ctx)

			// Delete FPGAFunctionCR
			err = deleteFPGAFunction(ctx, FPGAFunction1)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}
			_, err = reconciler3.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night01-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}
			// Get ChildBs CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01-00000001",
			},
				&childBsCR)
			if err != nil {
				fmt.Println("get ChildBsCR error", err)
			} else {
				printChildBitstream(ctx, childBsCR)
				for regionIndex := 0; regionIndex < len(childBsCR.Spec.Regions); regionIndex++ {
					specRegion := childBsCR.Spec.Regions[regionIndex]
					for functionIndex := 0; functionIndex < len(*specRegion.Modules.Functions); functionIndex++ {
						specFunctionData := (*specRegion.Modules.Functions)[functionIndex]
						for funcmoduleIndex := 0; funcmoduleIndex < len(*specFunctionData.Module); funcmoduleIndex++ {
							functionChannelIDs := (*specFunctionData.Module)[funcmoduleIndex].FunctionChannelIDs
							if nil == functionChannelIDs {
								continue
							}
							idRange := strings.SplitN(*functionChannelIDs, "-", 2)
							idRangeMin, _ := strconv.Atoi(idRange[0])
							idRangeMax, _ := strconv.Atoi(idRange[1])
							for idIndex := idRangeMin; idIndex <= idRangeMax; idIndex++ {
								idIndexString := strconv.Itoa(idIndex)
								intraResourceMgmt := (*specFunctionData.IntraResourceMgmtMap)[idIndexString]
								fmt.Println("IntraResourceMgmtMap[" + idIndexString +
									"].Available=" + strconv.FormatBool(*intraResourceMgmt.Available))
								if nil == intraResourceMgmt.FunctionCRName {
									fmt.Println("IntraResourceMgmtMap[" + idIndexString + "].FunctionCRName=nil")
								} else {
									fmt.Println("IntraResourceMgmtMap[" + idIndexString +
										"].FunctionCRName=" + *intraResourceMgmt.FunctionCRName)
								}
								if false == *intraResourceMgmt.Available {
									Expect(*intraResourceMgmt.Available).Should(BeFalse())
								} else {
									Expect(*intraResourceMgmt.Available).Should(BeTrue())
								}
								if nil == intraResourceMgmt.FunctionCRName {
									Expect(intraResourceMgmt.FunctionCRName).Should(BeNil())
								} else {
									Expect(*intraResourceMgmt.FunctionCRName).To(Equal("df-night01-wbfunction-filter-resize-high-infer-main"))
								}
							}
						}
					}
				}
			}

			fmt.Println("test7:^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
			fmt.Println(writer3.String())

			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			By("Test End")
			writer3.Reset()
		})
	})

	var fakerecorder4 = record.NewFakeRecorder(10)
	var writer4 = bytes.Buffer{}

	Context("Test for FPGAFunctionReconciler4", Ordered, func() {
		var reconciler4 FPGAFunctionReconciler
		BeforeAll(func() {

			opts := zap.Options{
				TimeEncoder: zapcore.ISO8601TimeEncoder,
				Development: true,
				DestWriter:  &writer4,
			}
			ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
		})

		BeforeEach(func() {
			writer4.Reset()

			// set manager
			mgr, err = getMgr(mgr)
			Expect(err).NotTo(HaveOccurred())

			reconciler4 = FPGAFunctionReconciler{
				Client:   k8sClient,
				Scheme:   testScheme,
				Recorder: fakerecorder4,
			}
			err = reconciler4.SetupWithManager(mgr)
			if err != nil {
				fmt.Println("Error in SetupWithManager")
			}
			Expect(err).NotTo(HaveOccurred())
			os.Setenv("K8S_NODENAME", "test01")
			os.Setenv("K8S_CLUSTERNAME", "default")
			StartupProccessing(mgr)
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestethernet.EthernetConnection{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestpcie.PCIeConnection{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestgpu.GPUFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestcpu.CPUFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGAFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGA{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.ChildBs{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.ComputeResource{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGAReconfiguration{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
		})
		It("FPGAFunctionTest 8", func() {
			By("Test Start")
			fmt.Println("test8:--------------------------------------------------")

			// Create ConfigMap
			err = config(ctx, fpgafuncconfig_fr_high_infer_1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, fpgafuncconfig_fr_low_infer_8)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, servicerMgmtConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, deployinfo_configdata8)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, regionUniqueInfoConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, functionUniqueInfoConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, preDeterminedRegionInfoConfig8)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, fr_childbs_Config1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}

			// Create FPGACR
			err = createFPGA(ctx, FPGA8[0])
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR ", err)
			}

			// Create ComputeResourceCR
			err = createCompureResource(ctx, comres8)
			if err != nil {
				fmt.Println("There is a problem in createing ComputeresourceCR ", err)
			}

			// Create FPGAReconfigurationCR
			err = createFPGAReconfiguration(ctx, FPGAReconfiguration8)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAReconfigurationCR ", err)
			}

			_, err := reconciler4.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "manualfpgareconfig-test01-21320621v01d",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			var childBsCR examplecomv1.ChildBs
			var childBitstreamID string
			// Get ChildBs CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v01d-test01-00000001",
			},
				&childBsCR)
			if err != nil {
				fmt.Println("Get ChildBsCR error.", err)
			} else {
				if nil != childBsCR.Spec.ChildBitstreamID {
					childBitstreamID = *childBsCR.Spec.ChildBitstreamID
				}
				printChildBitstream(ctx, childBsCR)
			}

			var fpgaCR examplecomv1.FPGA
			// Get fpga CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v01d-test01",
			},
				&fpgaCR)
			if err != nil {
				fmt.Println("Get fpgaCR error.", err)
			} else {
				printFPGA(ctx, fpgaCR)
				Expect(fpgaCR.Status.ChildBitstreamID).ShouldNot(BeNil())
				Expect(fpgaCR.Status.ChildBitstreamCRName).ShouldNot(BeNil())
				Expect(*fpgaCR.Status.ChildBitstreamID).To(Equal(childBitstreamID))
				Expect(*fpgaCR.Status.ChildBitstreamCRName).To(Equal(fpgaCR.Name + "-" + childBitstreamID))
			}

			var fpgaReconfiguratonCR examplecomv1.FPGAReconfiguration
			// Get fpgaReconfiguraton CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "manualfpgareconfig-test01-21320621v01d",
			},
				&fpgaReconfiguratonCR)
			if err != nil {
				fmt.Println("Get fpgaReconfiguratonCR error.", err)
			} else {
				printFPGAReconfiguration(ctx, fpgaReconfiguratonCR)
				// Expect(fpgaReconfiguratonCR.Status.Status).To(Equal("Succeeded"))
			}

			printDeployInfo(ctx)

			// Delete fpgaReconfiguraton CR
			err = k8sClient.Delete(ctx, &fpgaReconfiguratonCR)
			if err != nil {
				fmt.Println("Delete fpgaReconfiguratonCR error.", err)
			}

			fmt.Println("Reconcile fpgaReconfiguratonCR delete start")
			_, err = reconciler4.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "manualfpgareconfig-test01-21320621v01d",
			}})
			if err != nil {
				By("Reconcile Error")
			}
			fmt.Println("Reconcile fpgaReconfiguratonCR dalete end")

			fmt.Println("test8:^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
			fmt.Println(writer4.String())

			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			By("Test End")
			writer4.Reset()
		})
	})

	var fakerecorder5 = record.NewFakeRecorder(10)
	var writer5 = bytes.Buffer{}

	Context("Test for FPGAFunctionReconciler5", Ordered, func() {
		var reconciler5 FPGAFunctionReconciler
		BeforeAll(func() {

			opts := zap.Options{
				TimeEncoder: zapcore.ISO8601TimeEncoder,
				Development: true,
				DestWriter:  &writer5,
			}
			ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
		})

		BeforeEach(func() {
			writer5.Reset()

			// set manager
			mgr, err = getMgr(mgr)
			Expect(err).NotTo(HaveOccurred())

			reconciler5 = FPGAFunctionReconciler{
				Client:   k8sClient,
				Scheme:   testScheme,
				Recorder: fakerecorder5,
			}
			err = reconciler5.SetupWithManager(mgr)
			if err != nil {
				fmt.Println("Error in SetupWithManager")
			}
			Expect(err).NotTo(HaveOccurred())
			os.Setenv("K8S_NODENAME", "test01")
			os.Setenv("K8S_CLUSTERNAME", "default")
			StartupProccessing(mgr)
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestethernet.EthernetConnection{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestpcie.PCIeConnection{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestgpu.GPUFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestcpu.CPUFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGAFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGA{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.ChildBs{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.ComputeResource{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGAReconfiguration{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
		})
		It("FPGAFunctionTest 9", func() {
			By("Test Start")
			fmt.Println("test9:--------------------------------------------------")

			// Create ConfigMap
			err = config(ctx, fpgafuncconfig_fr_high_infer_1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, fpgafuncconfig_fr_low_infer_8)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, servicerMgmtConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, deployinfo_configdata9)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, regionUniqueInfoConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, functionUniqueInfoConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, preDeterminedRegionInfoConfig8)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, fr_childbs_Config1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}

			// Create FPGACR
			err = createFPGA(ctx, FPGA9[0])
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR ", err)
			}

			// Create ComputeResourceCR
			err = createCompureResource(ctx, comres8)
			if err != nil {
				fmt.Println("There is a problem in createing ComputeresourceCR ", err)
			}

			// Create ChildBitstreamCR
			err = createChildBitstream(ctx, ChildBitstream9)
			if err != nil {
				fmt.Println("There is a problem in createing ChildBitstream ", err)
			}

			var childBsCR examplecomv1.ChildBs
			// Get ChildBs CR
			_ = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v01d-test01-00000002",
			},
				&childBsCR)
			if err != nil {
				fmt.Println("get childBsCR error.", err)
			}

			childBsCR.Status.Regions = childBsCR.Spec.Regions
			childBsCR.Status.Status = examplecomv1.ChildBsStatusReady
			childBsCR.Status.State = examplecomv1.ChildBsReady

			// update ChildBs CR
			err = updateChildBitstream(ctx, childBsCR)
			if err != nil {
				fmt.Println("There is a problem in updateing ChildBsCR ", err)
			}

			printChildBitstream(ctx, childBsCR)

			// Create FPGAReconfigurationCR
			err = createFPGAReconfiguration(ctx, FPGAReconfiguration9)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAReconfigurationCR ", err)
			}

			printDeployInfo(ctx)

			_, err := reconciler5.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "manualfpgareconfig-test01-21320621v01d",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			printDeployInfo(ctx)

			// Get ChildBs CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v01d-test01-00000002",
			},
				&childBsCR)
			if err != nil {
				fmt.Println("get ChildBsCR error.", err)
			} else {
				printChildBitstream(ctx, childBsCR)
			}

			_, err = reconciler5.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "manualfpgareconfig-test01-21320621v01d",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			printDeployInfo(ctx)

			// Get ChildBs CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v01d-test01-00000002",
			},
				&childBsCR)
			if err != nil {
				fmt.Println("get ChildBsCR error.", err)
			} else {
				printChildBitstream(ctx, childBsCR)
			}

			childBsCR.Status.State = examplecomv1.ChildBsStoppingNetworkModule

			// update ChildBs CR
			err = updateChildBitstream(ctx, childBsCR)
			if err != nil {
				fmt.Println("There is a problem in updateing ChildBsCR ", err)
			}

			printChildBitstream(ctx, childBsCR)

			_, err = reconciler5.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "manualfpgareconfig-test01-21320621v01d",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			printDeployInfo(ctx)

			// Get ChildBs CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v01d-test01-00000002",
			},
				&childBsCR)
			if err != nil {
				fmt.Println("get ChildBsCR error.", err)
			} else {
				printChildBitstream(ctx, childBsCR)
			}

			childBsCR.Status.State = examplecomv1.ChildBsNotWriteBsfile

			// update ChildBs CR
			err = updateChildBitstream(ctx, childBsCR)
			if err != nil {
				fmt.Println("There is a problem in updateing ChildBsCR ", err)
			}

			printChildBitstream(ctx, childBsCR)

			_, err = reconciler5.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "manualfpgareconfig-test01-21320621v01d",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			printDeployInfo(ctx)

			// Get ChildBs CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v01d-test01-00000002",
			},
				&childBsCR)
			if err != nil {
				fmt.Println("get ChildBsCR error.", err)
			} else {
				printChildBitstream(ctx, childBsCR)
			}

			var fpgaReconfiguratonCR examplecomv1.FPGAReconfiguration
			// Get fpgaReconfiguraton CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "manualfpgareconfig-test01-21320621v01d",
			},
				&fpgaReconfiguratonCR)
			if err != nil {
				fmt.Println("Get fpgaReconfiguratonCR error.", err)
			} else {
				printFPGAReconfiguration(ctx, fpgaReconfiguratonCR)
				Expect(fpgaReconfiguratonCR.Status.Status).To(Equal("Succeeded"))
			}

			// Delete fpgaReconfiguraton CR
			err = k8sClient.Delete(ctx, &fpgaReconfiguratonCR)
			if err != nil {
				fmt.Println("Delete fpgaReconfiguratonCR error.", err)
			}

			fmt.Println("Reconcile fpgaReconfiguratonCR delete start")
			_, err = reconciler5.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "manualfpgareconfig-test01-21320621v01d",
			}})
			if err != nil {
				By("Reconcile Error")
				fmt.Println("Reconcile fpgaReconfiguratonCR dalete error")
				err = nil
			}
			fmt.Println("Reconcile fpgaReconfiguratonCR dalete end")

			fmt.Println("test9:^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
			fmt.Println(writer5.String())

			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			By("Test End")
			writer5.Reset()
		})
	})

	var fakerecorder6 = record.NewFakeRecorder(10)
	var writer6 = bytes.Buffer{}

	Context("Test for FPGAFunctionReconciler6", Ordered, func() {
		var reconciler6 FPGAFunctionReconciler
		BeforeAll(func() {

			opts := zap.Options{
				TimeEncoder: zapcore.ISO8601TimeEncoder,
				Development: true,
				DestWriter:  &writer6,
			}
			ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
		})

		BeforeEach(func() {
			writer6.Reset()

			// set manager
			mgr, err = getMgr(mgr)
			Expect(err).NotTo(HaveOccurred())

			reconciler6 = FPGAFunctionReconciler{
				Client:   k8sClient,
				Scheme:   testScheme,
				Recorder: fakerecorder6,
			}
			err = reconciler6.SetupWithManager(mgr)
			if err != nil {
				fmt.Println("Error in SetupWithManager")
			}
			Expect(err).NotTo(HaveOccurred())
			os.Setenv("K8S_NODENAME", "test01")
			os.Setenv("K8S_CLUSTERNAME", "default")
			StartupProccessing(mgr)
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestethernet.EthernetConnection{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestpcie.PCIeConnection{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestgpu.GPUFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestcpu.CPUFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGAFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGA{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.ChildBs{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.ComputeResource{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGAReconfiguration{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
		})
		It("FPGAFunctionTest 10", func() {
			By("Test Start")
			fmt.Println("test10:--------------------------------------------------")

			// Create ConfigMap
			err = config(ctx, fpgafuncconfig_fr_high_infer_1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, fpgafuncconfig_fr_low_infer_8)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, servicerMgmtConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, deployinfo_configdata9)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, regionUniqueInfoConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, functionUniqueInfoConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, preDeterminedRegionInfoConfig8)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, fr_childbs_Config1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}

			// Create FPGACR
			err = createFPGA(ctx, FPGA9[0])
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR ", err)
			}

			// Create ComputeResourceCR
			err = createCompureResource(ctx, comres8)
			if err != nil {
				fmt.Println("There is a problem in createing ComputeresourceCR ", err)
			}

			// Create ChildBitstreamCR
			err = createChildBitstream(ctx, ChildBitstream9)
			if err != nil {
				fmt.Println("There is a problem in createing ChildBitstream ", err)
			}

			var childBsCR examplecomv1.ChildBs
			// Get ChildBs CR
			_ = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v01d-test01-00000002",
			},
				&childBsCR)

			childBsCR.Status.Regions = childBsCR.Spec.Regions
			childBsCR.Status.Status = examplecomv1.ChildBsStatusReady
			childBsCR.Status.State = examplecomv1.ChildBsReady

			// update ChildBs CR
			err = updateChildBitstream(ctx, childBsCR)
			if err != nil {
				fmt.Println("There is a problem in updateing ChildBsCR ", err)
			}
			printChildBitstream(ctx, childBsCR)

			// Create FPGAReconfigurationCR
			err = createFPGAReconfiguration(ctx, FPGAReconfiguration10)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAReconfigurationCR ", err)
			}

			printDeployInfo(ctx)

			_, err := reconciler6.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "manualfpgareconfig-test01-21320621v01d",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			printDeployInfo(ctx)

			// Get ChildBs CR
			_ = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v01d-test01-00000002",
			},
				&childBsCR)

			printChildBitstream(ctx, childBsCR)

			_, err = reconciler6.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "manualfpgareconfig-test01-21320621v01d",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			printDeployInfo(ctx)

			// Get ChildBs CR
			_ = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v01d-test01-00000002",
			},
				&childBsCR)

			printChildBitstream(ctx, childBsCR)

			childBsCR.Status.State = examplecomv1.ChildBsStoppingNetworkModule

			// update ChildBs CR
			err = updateChildBitstream(ctx, childBsCR)
			if err != nil {
				fmt.Println("There is a problem in updateing ChildBsCR ", err)
			}
			printChildBitstream(ctx, childBsCR)

			_, err = reconciler6.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "manualfpgareconfig-test01-21320621v01d",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			printDeployInfo(ctx)

			// Get ChildBs CR
			_ = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v01d-test01-00000002",
			},
				&childBsCR)

			printChildBitstream(ctx, childBsCR)

			childBsCR.Status.State = examplecomv1.ChildBsNotWriteBsfile

			// update ChildBs CR
			err = updateChildBitstream(ctx, childBsCR)
			if err != nil {
				fmt.Println("There is a problem in updateing ChildBsCR ", err)
			}
			printChildBitstream(ctx, childBsCR)

			printDeployInfo(ctx)

			_, err = reconciler6.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "manualfpgareconfig-test01-21320621v01d",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			printDeployInfo(ctx)

			// Get ChildBs CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v01d-test01-00000002",
			},
				&childBsCR)
			if errors.IsNotFound(err) {
				fmt.Println("delete ChildBsCR success.")
			} else if err != nil {
				fmt.Println("get ChildBsCR error.", err)
			} else {
				fmt.Println("delete ChildBsCR found error.")
			}

			var fpgaCR examplecomv1.FPGA
			// Get fpga CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v01d-test01",
			},
				&fpgaCR)
			if err != nil {
				fmt.Println("Get fpgaCR error.", err)
			} else {
				printFPGA(ctx, fpgaCR)
				Expect(fpgaCR.Status.ChildBitstreamID).Should(BeNil())
				Expect(fpgaCR.Status.ChildBitstreamCRName).Should(BeNil())
			}

			printDeployInfo(ctx)

			var fpgaReconfiguratonCR examplecomv1.FPGAReconfiguration
			// Get fpgaReconfiguraton CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "manualfpgareconfig-test01-21320621v01d",
			},
				&fpgaReconfiguratonCR)
			if err != nil {
				fmt.Println("Get fpgaReconfiguratonCR error.", err)
			} else {
				printFPGAReconfiguration(ctx, fpgaReconfiguratonCR)
				Expect(fpgaReconfiguratonCR.Status.Status).To(Equal("Succeeded"))
			}
			// Delete fpgaReconfiguraton CR
			err = k8sClient.Delete(ctx, &fpgaReconfiguratonCR)
			if err != nil {
				fmt.Println("Delete fpgaReconfiguratonCR error.", err)
			}

			fmt.Println("Reconcile fpgaReconfiguratonCR delete start")
			_, err = reconciler6.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "manualfpgareconfig-test01-21320621v01d",
			}})
			if err != nil {
				By("Reconcile Error")
			}
			fmt.Println("Reconcile fpgaReconfiguratonCR dalete end")

			fmt.Println("test10:^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
			fmt.Println(writer6.String())

			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			By("Test End")
			writer6.Reset()
		})
	})

	var fakerecorder7 = record.NewFakeRecorder(10)
	var writer7 = bytes.Buffer{}

	Context("Test for FPGAFunctionReconciler7", Ordered, func() {
		var reconciler7 FPGAFunctionReconciler
		BeforeAll(func() {

			opts := zap.Options{
				TimeEncoder: zapcore.ISO8601TimeEncoder,
				Development: true,
				DestWriter:  &writer7,
			}
			ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
		})

		BeforeEach(func() {
			writer7.Reset()

			// set manager
			mgr, err = getMgr(mgr)
			Expect(err).NotTo(HaveOccurred())

			reconciler7 = FPGAFunctionReconciler{
				Client:   k8sClient,
				Scheme:   testScheme,
				Recorder: fakerecorder7,
			}
			err = reconciler7.SetupWithManager(mgr)
			if err != nil {
				fmt.Println("Error in SetupWithManager")
			}
			Expect(err).NotTo(HaveOccurred())
			os.Setenv("K8S_NODENAME", "test01")
			os.Setenv("K8S_CLUSTERNAME", "default")
			StartupProccessing(mgr)
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestethernet.EthernetConnection{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestpcie.PCIeConnection{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestgpu.GPUFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &controllertestcpu.CPUFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGAFunction{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGA{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.ChildBs{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.ComputeResource{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &examplecomv1.FPGAReconfiguration{}, client.InNamespace("default"))
			Expect(err).NotTo(HaveOccurred())
		})
		It("FPGAFunctionTest 2-1", func() {
			By("Test Start")
			fmt.Println("test2-1:--------------------------------------------------")
			// It was made based on Test4.
			// The generation of PCIeConnection3 was deleted.
			// Two Reconcile were added..

			// Create ConfigMap
			err = config(ctx, fpgafuncconfig_fr_high_infer_3)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, servicerMgmtConfig2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, deployinfo_configdata2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, regionUniqueInfoConfig1)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, functionUniqueInfoConfig2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, preDeterminedRegionInfoConfig8)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}
			err = config(ctx, fr_childbs_Config2)
			if err != nil {
				fmt.Println("There is a problem in createing CM : ", err)
			}

			// Create PCIeConnectionCR
			err = createPCIeConnection(ctx, PCIeConnection4)
			if err != nil {
				fmt.Println("There is a problem in createing PCIeConnectionCR : ", err)
			}

			// Create CPUFunctionCR
			err = createCPUFunction(ctx, CPUFunction3)
			if err != nil {
				fmt.Println("There is a problem in createing CPUFuncCR : ", err)
			}

			// Create GPUFunctionCR
			err = createGPUFunction(ctx, GPUFunction3)
			if err != nil {
				fmt.Println("There is a problem in createing GPUFuncCR : ", err)
			}
			// Create FPGACR
			err = createFPGA(ctx, FPGA3[0])
			if err != nil {
				fmt.Println("There is a problem in createing FPGACR ", err)
			}

			// Create ComputeResourceCR
			err = createCompureResource(ctx, comres1)
			if err != nil {
				fmt.Println("There is a problem in createing ComputeresourceCR ", err)
			}

			// Create FPGAFunctionCR
			err = createFPGAFunction(ctx, FPGAFunction3)
			if err != nil {
				fmt.Println("There is a problem in createing FPGAFuncCR ", err)
			}

			_, err := reconciler7.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night03-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			_, err = reconciler7.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night03-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			_, err = reconciler7.Reconcile(ctx, ctrl.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "df-night03-wbfunction-filter-resize-high-infer-main",
			}})
			if err != nil {
				By("Reconcile Error")
			}

			var childBsCR examplecomv1.ChildBs
			// Get ChildBs CR
			_ = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01-00000001",
			},
				&childBsCR)

			printChildBitstream(ctx, childBsCR)
			printDeployInfo(ctx)

			var fpgaCR examplecomv1.FPGA
			// Get fpga CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "fpga-21320621v00d-test01",
			},
				&fpgaCR)
			if err != nil {
				fmt.Println("Get fpgaCR error.", err)
			} else {
				printFPGA(ctx, fpgaCR)
			}

			var fpgafuncCR examplecomv1.FPGAFunction
			// Get fpga CR
			err = k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "df-night03-wbfunction-filter-resize-high-infer-main",
			},
				&fpgafuncCR)
			if err != nil {
				fmt.Println("get fpgafuncCR error.", err)
			} else {
				printFPGAFunction(ctx, fpgafuncCR)
				Expect(fpgafuncCR.Status.Rx.Protocol).To(Equal(""))
				Expect(fpgafuncCR.Status.Tx.Protocol).To(Equal(""))
			}

			controllerutil.RemoveFinalizer(&fpgafuncCR, "fpgafunction.finalizers.example.com.v1")
			// Update RemoveFinalizer to fpgafunc CR
			err = k8sClient.Update(ctx, &fpgafuncCR)
			if err != nil {
				fmt.Println("error update RemoveFinalizer to fpgafuncCR.", err)
			}

			fmt.Println("test2-1:^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
			fmt.Println(writer7.String())
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			By("Test End")
			writer7.Reset()
		})
	})
})

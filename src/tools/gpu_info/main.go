/* Copyright 2024 NTT Corporation , FUJITSU LIMITED */
package main

import (
	"fmt"
	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

func main() {
	creanup, err := dcgm.Init(dcgm.StartHostengine)
	defer creanup()
	if nil != err {
		return
	}
	count, _ := dcgm.GetSupportedDevices()
	for i := 0; i < len(count); i++ {
		device, _ := dcgm.GetDeviceInfo(count[i])
		fmt.Println("=[ device info Start ]============================================== ")
		fmt.Println("device.GPU                :", device.GPU)
		fmt.Println("device.DCGMSupported      :", device.DCGMSupported)
		fmt.Println("device.UUID               :", device.UUID)
		fmt.Println("device.Power              :", device.Power)
		fmt.Println("device.PCI.BusID          :", device.PCI.BusID)
		fmt.Println("device.PCI.BAR1           :", device.PCI.BAR1)
		fmt.Println("device.PCI.FBTotal        :", device.PCI.FBTotal)
		fmt.Println("device.PCI.Bandwidth      :", device.PCI.Bandwidth)
		fmt.Println("device.Identifiers.Brand  :", device.Identifiers.Brand)
		fmt.Println("device.Identifiers.Model  :", device.Identifiers.Model)
		fmt.Println("device.Identifiers.Serial :", device.Identifiers.Serial)
		fmt.Println("device.Identifiers.Vbios  :", device.Identifiers.Vbios)
		fmt.Println("device.Identifiers.InforomImageVersion :", device.Identifiers.InforomImageVersion)
		fmt.Println("device.Identifiers.DriverVersion       :", device.Identifiers.DriverVersion)
		fmt.Println("device.CPUAffinity        :", device.CPUAffinity)
		fmt.Println("=[ device info End   ]============================================== ")
	}
}

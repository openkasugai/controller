/* Copyright 2025 NTT Corporation , FUJITSU LIMITED */

package common

import (
	metav1 "example.com/FPGACheckPerDF/api/v1"
)

// custom resource save data
type CrDataSt struct {
	CrDf       metav1.DataFlow
	CrWBConn   metav1.WBConnection
	CrEthConn  metav1.EthernetConnection
	CrPCIeConn metav1.PCIeConnection
	CrWBFunc   metav1.WBFunction
	CrFPGAFunc metav1.FPGAFunction
	CrCPUFunc  metav1.CPUFunction
	CrGPUFunc  metav1.GPUFunction
}

// custom resource save table
type CrTable struct {
	ChainFlag     bool   `json:"chainflag"`
	FPGAPosition  int    `json:"fpgaposition"`
	APIVersion    string `json:"apiversion"`
	Kind          string `json:"kind"`
	Namespace     string `json:"namespace"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	APIVersionSub string `json:"apiversionsub"`
	KindSub       string `json:"kindsub"`
	StatusSub     string `json:"statussub"`
	CrData        CrDataSt
}

const (
	SHELLFOLDER          = "/tmp/FPGACheckPerDF/"
	SHELLNAME            = "FPGACheckPerDF-chk-"
	SHELLEXTENSION       = ".sh"
	LOGFOLDER            = "/tmp/FPGACheckPerDF/"
	LOGNAME              = "FPGACheckPerDF-chk-"
	LOGEXTENSION         = ".log"
	FPGACOMMANDPATH      = "../fpga-chk-connection/bin/fpga-chk-connection"
	FPGAPOSITIONNONE     = 0
	FPGAPOSITIONPREVIOUS = 1
	FPGAPOSITIONPERSONAL = 2
	FPGAPOSITIONNEXT     = 3
)

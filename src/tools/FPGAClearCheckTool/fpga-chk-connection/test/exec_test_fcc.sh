#! /bin/bash
# Copyright 2025 NTT Corporation , FUJITSU LIMITED

set -eu

# Set environment parameters
INPUT_FPGA_NAME="21320621V01M"	# /dev/xpcie_21320621V01M
INPUT_HOST_ADDR_LANE0="C0A8006F"
INPUT_HOST_ADDR_LANE1="C0A8016F"

# Set FPGA's parameters
# When you want to change this values, you should change source code too.
INPUT_FPGA_ADDR_LANE0="C0A80065"
INPUT_FPGA_GATEWAY_LANE0="C0A80001"
INPUT_FPGA_SUBNET_LANE0="FFFFFF00"
INPUT_FPGA_ADDR_LANE1="C0A80165"
INPUT_FPGA_GATEWAY_LANE1="C0A80101"
INPUT_FPGA_SUBNET_LANE1="FFFFFF00"

usage () {
	echo "usage: $0"
	echo "Note:"
	echo " Please change below parameters in this file:"
	echo "   - INPUT_FPGA_NAME       : FPGA's UUID"
	echo "   - INPUT_HOST_ADDR_LANE0 : Host's address connect to FPGA's lane0"
	echo "                           : See below FPGA's setting."
	echo "   - INPUT_HOST_ADDR_LANE1 : Host's address connect to FPGA's lane1"
	echo "                           : See below FPGA's setting."
	echo " FPGA's IP settings for Lane0:"
	echo "   - Address : $INPUT_FPGA_ADDR_LANE0"
	echo "   - Gateway : $INPUT_FPGA_GATEWAY_LANE0"
	echo "   - Subnet  : $INPUT_FPGA_SUBNET_LANE0"
	echo " FPGA's IP settings for Lane1:"
	echo "   - Address : $INPUT_FPGA_ADDR_LANE1"
	echo "   - Gateway : $INPUT_FPGA_GATEWAY_LANE1"
	echo "   - Subnet  : $INPUT_FPGA_SUBNET_LANE1"
}

main () {
	# include shrc
	. ./shrc/setup_connection_def.shrc
	. ./shrc/setup_connection_case_lldma.shrc
	. ./shrc/setup_connection_case_ptu.shrc
	. ./shrc/setup_connection_case_chain.shrc
	. ./shrc/setup_connection_case_option.shrc

	# Check parameter
	if [ $# -ne 0 ];then
		if [ "$1" = "-h" ];then
			usage
			exit 0
		elif [ "$1" = "--help" ];then
			usage
			exit 0
		fi
	fi
	if [ ! -e "/dev/xpcie_$INPUT_FPGA_NAME" ];then
		echo "** Invalid device file : /dev/xpcie_$INPUT_FPGA_NAME"
		usage
		exit 0
	fi

	# Check binary
	if [ ! -e $CMD_SETUP_CONNECTION ];then
		echo " ! Not found : $CMD_SETUP_CONNECTION"
		exit 1
	fi
	if [ ! -e $CMD_FCC ];then
		echo " ! Not found : $CMD_FCC"
		exit 1
	fi

	echo "# Start : ALL : `date +%H-%M-%S`"

	# Pre process
	exec_daemon

	# Execute test
	exec_test test_lldma
	exec_test test_ptu
	exec_test test_chain
	exec_test test_option

	echo "# End   : ALL : `date +%H-%M-%S`"

	# Post process
	request_quit || true

	mv_logfiles
}

main $* | tee -a main.log

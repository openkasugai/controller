# Test for fpga-chk-connection

## Setup
1. Install libfpga if not installed.
2. Install DPDK if not installed.
3. Change libfpga's include directory path(`LIBFPGA_INC_DIR`) in Makefile if needed.
4. Change libfpga's library directory path(`LIBFPGA_LIB_DIR`) in Makefile if needed.
5. Change DPDK's library directory path(`DPDK_LIB_DIR`) in Makefile if needed.
6. Build `setup-connection` as follows:
```sh
make
```

## Build Commands
|Command|Description|
|-|-|
|make (default)|= make static|
|make all|make clean + make static|
|make static|create `setup-connection`|
|make clean|clean `setup-connection` and bin directory|

## Note
- Please be careful of FPGA's parameters:
	- FPGA's IP settings for Lane0:
		- Address : C0A80065(192.168.0.101)
		- Gateway : **C0A800**01(192.168.0.1)
		- Subnet  : FFFFFF00(255.255.255.0)
		- Port    : 20000-20100
	- FPGA's IP settings for Lane1:
		- Address : C0A80165(192.168.1.101)
		- Gateway : **C0A801**01(192.168.1.1)
		- Subnet  : FFFFFF00(255.255.255.0)
		- Port    : 20000-20100

## Prepare
- Set parameters in exec_test_fcc.sh at first.
	|Parameter|Description|Default|
	|-|-|-|
	|INPUT_FPGA_NAME|FPGA's UUID|-|
	|INPUT_HOST_ADDR_LANE0|Host's address connect to FPGA's lane0<br>port will be used 20000-20100|**C0A800**6F(192.168.0.111 )|
	|INPUT_HOST_ADDR_LANE1|Host's address connect to FPGA's lane1<br>port will be used 20000-20100|**C0A801**6F(192.168.1.111 )|

- Setup IP address for INPUT_HOST_ADDR_LANE0 and INPUT_HOST_ADDR_LANE1 if needed.
	```sh
	$ setup_ip () {
		# 1 : IP_ADDRESS(format:x.x.x.x)
		# 2 : DEVICE
		sudo ip addr add ${1}/24 dev $2
		sudo ip link set dev $2 up
		sleep 1
		ip -4 a show dev $2
	}
	$ setup_ip 192.168.0.111 ens6f0
	$ setup_ip 192.168.1.111 ens6f1
	```
- Load xpcie driver if xpcie driver is not loaded.
	```sh
	$ lsmod | grep xpcie
	$ sudo insmod path/to/xpcie.ko
	$ lsmod | grep xpcie
	xpcie                 110592  0
	```

## Execute
- Command:
	```sh
	$ ./exec_test_fcc.sh [-h|--help]
	```
- Execute exec_test_fcc.sh.
	```sh
	$ ./exec_test_fcc.sh # about 15minutes
	```
- Check if error happened(exist "!" in log file) or not.
	```sh
	$ grep "!" log -rn
	```

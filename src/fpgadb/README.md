# FPGA Additional Database library(libfpgadb)

## Introduction
- `libfpgadb` is a library to get configuration information about FPGAs from json file
- `libfpgadb` is based on `libfpga`.

## Setup
1. Inastall libfpga if not installed.
2. Change libfpga's include directory path(`LIBFPGA_INC_DIR`) in Makefile if needed.
3. Build libfpgadb as follows:
	- Parson is installed automatically via `git clone`.
		Therefore, if Git is not installed, you should install Git first.
```sh
make
```

### Build Commands
|Command|Description|
|-|-|
|make (default)|= make static|
|make all|make clean + make static|
|make static|create libfpgadb.a|
|make shared|create libfpgadb.so(not tested)|
|make clean|clean all object files and libraries and build directory|

## APIs
- See `libfpgadb.h` about the details.
	- fpga_db_get_bitstream_id
	- fpga_db_get_bitstream_id_by_dev_id
	- fpga_db_get_device_config
	- fpga_db_get_device_config_by_dev_id
	- fpga_db_get_device_config_by_bitstream_id
	- fpga_db_get_child_bitstream_ids
	- fpga_db_get_child_bitstream_ids_by_dev_id
	- fpga_db_get_child_bitstream_ids_by_parent
	- fpga_db_free_child_bitstream_ids
	- fpga_db_enable_dummy_bitstream
	- fpga_db_enable_dummy_bitstream_by_dev_id
	- fpga_db_disable_dummy_bitstream
	- fpga_db_disable_dummy_bitstream_by_dev_id

## Remarks
- This library uses only parson.h but does not use parson.c in default,
	because this library uses objects of parson.c.o in `libfpga`.
	- If you store 3rdparty/parson/parson.h in advance,
		`git clone` for parson will be skipped.
- This library can include parson.c.o with setting `ENABLE_INCLUDE_PARSON=1` as Makefile argument if you need.
```
make ENABLE_INCLUDE_PARSON=1
```
```
make ENABLE_INCLUDE_PARSON=1 shared
```


## bitstream_id-config-table.json
- When you want to update the file, please check the below logs or register values of FPGA.
- `parent-bitstream-id`
	- offset : PCI configuration register+0x35C
	- use setpci(device id: 0x903f)
	```sh
	$ sudo setpci -d:903f 35C.l
	00000000
	```
	- use setpci(slot: <BB>:<DD>.<F>)
	```sh
	$ sudo setpci -s 1f:00.0 35C.l
	00000000
	```
	- use dmesg after loading xpcie.ko
	```sh
	$ sudo dmesg | grep ParentBitstream
	[  185.533297] xpcie: ParentBitstream=00000000
	```
- `child-bitstream-id`
	- use dbgreg tool(reg32r)
	- offset : 0x0
	```sh
	$ ./reg32r /dev/xpcie_21320621V01M 0
	offset:        0
	 0000 : 00012114
	```
- `chain`
	- use dbgreg tool(reg32r)
	- `identifier`
		- offset : 0x1010(lane0), 0x2010(lane1)
		```sh
		$ ./reg32r /dev/xpcie_21320621V01M 1010
		offset:        0
		 1010 : 0000f0c0
		```
	- `version`
		- offset : 0x1020(lane0), 0x2020(lane1)
		```sh
		$ ./reg32r /dev/xpcie_21320621V01M 2020
		offset:        0
		 1020 : 24052427
		```
- `directtrans`
	- use dbgreg tool(reg32r)
	- `identifier`
		- offset : 0x5010(lane0), 0x6010(lane1)
		```sh
		$ ./reg32r /dev/xpcie_21320621V01M 5010
		offset:        0
		 5010 : 0000f3c0
		```
	- `version`
		- offset : 0x5020(lane0), 0x6020(lane1)
		```sh
		$ ./reg32r /dev/xpcie_21320621V01M 5020
		offset:        0
		 5020 : 24052406
		```
- `conversion`
	- use dbgreg tool(reg32r)
	- `identifier`
		- offset : 0x24010(lane0), 0x24410(lane1)
			- 0x400 bytes per a lane
		```sh
		$ ./reg32r /dev/xpcie_21320621V01M 24010
		offset:        0
		 24010 : 0000f1c2
		```
	- `version`
		- offset : 0x24020(lane0), 0x24420(lane1)
			- 0x400 bytes per a lane
		```sh
		$ ./reg32r /dev/xpcie_21320621V01M 24020
		offset:        0
		 24020 : 24052407
		```
- `functions`
	- use dbgreg tool(reg32r)
	- `identifier`
		- offset : 0x25010(lane0), 0x26010(lane1)
		```sh
		$ ./reg32r /dev/xpcie_21320621V01M 25010
		offset:        0
		 25010 : 0000f2c2
		```
	- `version`
		- offset : 0x25020(lane0), 0x26020(lane1)
		```sh
		$ ./reg32r /dev/xpcie_21320621V01M 25020
		offset:        0
		 25020 : 24052427
		```

## Test and Example
- See test/README.md.

## Files
```
fpgadb
 +---- Makefile
 +---- README.md
 +---- 3rdparty/
 +---- include/
 |     \---- libfpgadb.h
 +---- src/
 |     \---- libfpgadb.c
 \---- test/
       +---- bitstream_id-config-table.json
       +---- dummy-bitstream_id-config-table.json
       +---- test-bitstream_id-config-table.json
       +---- main.c
       +---- Makefile
       \---- README.md
```

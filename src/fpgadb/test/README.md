# Test for libfpgadb

## Setup
1. Inastall libfpga if not installed.
2. Build libfpgadb if not build.
3. Change Makefile parameters if needed:
	- libfpga's library directory path(`LIBFPGA_LIB_DIR`)
		- default : /usr/local/include/fpgalib
	- libfpga's include directory path(`LIBFPGA_INC_DIR`)
		- default : /usr/local/lib/fpgalib
	- libfpgadb's library directory path(`LIBFPGA_LIB_DIR`)
		- default : ../build
	- libfpgadb's include directory path(`LIBFPGA_INC_DIR`)
		- default : ../build/include
4. Build App.
```sh
make
```

### Build Commands
|Command|Description|
|-|-|
|make (default)|= make all|
|make all|make clean + make build|
|make build|Build App|
|make clean|clean all object files and libraries and build directory|

## Execute options
- Get bitstream id from FPGA and create files with available configuration information.
```sh
./test [options for liblogging]
```

- Get bitstream id from scanf() and create files with available configuration information.
```sh
./test [options for liblogging] -- --without-fpga
```

## Execute
```sh
$ ./test
 * Execute with FPGAs
 *  Use 'bitstream_id-config-table.json'
 * Get FPGAs' serial ids
 * Device[21320621V01M]
 *    Create 'config-device-21320621V01M.json'
 * Parent BSID : 00000000
 * This parent-bitstream-id(00000000) can use the below child-bitstream-id:
 * - 00012114
 *    Create 'config-available-parent(00000000)-child(00012114).json'
```

```sh
$ ./test -- --without-fpga
 * Execute without FPGAs
 *  Use 'dummy-bitstream_id-config-table.json'
 * Input Parent Bitstream ID in hex.
 > ffffffff
 * Parent BSID : ffffffff
 * This parent-bitstream-id(ffffffff) can use the below child-bitstream-id:
 * - fffffff0
 *    Create 'config-available-parent(ffffffff)-child(fffffff0).json'
 * - fffffff1
 *    Create 'config-available-parent(ffffffff)-child(fffffff1).json'
 * - fffffff2
 *    Create 'config-available-parent(ffffffff)-child(fffffff2).json'
```

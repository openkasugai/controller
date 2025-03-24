# FPGA Check Connection tool

## Introduction
- `fpga-chk-connection` is a tool to check FPGA connection settings.
	- LLDMA channel is used or not.
		- Checked by: `Connector_id`
	- PTU session is alive or not.
		- Checked by: `Device file path`, `Lane`, `External interface id`, `Connection_id`
	- Chain setting is set or not.
		- Checked by: `Device file path`, `Lane`, `Function channel id`(, `Direction`)
- `fpga-chk-connection` needs `DPDK`, `libfpga` and `parson`.

## Setup
1. Install libfpga if not installed.
2. Install DPDK if not installed.
3. Change libfpga's include directory path(`LIBFPGA_INC_DIR`) in Makefile if needed.
4. Change libfpga's library directory path(`LIBFPGA_LIB_DIR`) in Makefile if needed.
5. Change DPDK's library directory path(`DPDK_LIB_DIR`) in Makefile if needed.
6. Build `fpga-chk-connection` as follows:
	- Parson will be installed automatically via `git clone`.
		Therefore, if Git is not installed, you should install Git first.
```sh
make
```

### Build Commands
|Command|Description|
|-|-|
|make (default)|= make static|
|make all|make clean + make static|
|make static|create `fpga-chk-connection`|
|make clean|clean `fpga-chk-connection` and bin directory|
|make clean-all|make clean + make remove-parson|
|make remove-parson|remove `parson` directory|

## Options
```sh
fpga-chk-connection : v01000000
usage: fpga-chk-connection [-dlfeckjio <PARAMETER>] [--dump] [-h]
       -d/--device <DEVICE>         : Device file path[/dev/xpcie_<UUID>,<UUID>]
       -l/--lane <LANE>             : Lane number[0-1]
       -f/--fchid <FCHID>           : Function channel id[0-511]
       -e/--extif_id <EXTIF_ID>     : External interface id[lldma,LLDMA,0,ptu,PTU,1]
       -c/--cid <CID>               : Connection id[0-15(LLDMA),1-511(PTU)]
          --dir <DIRECTION>         : Direction[ingress,1,egress,2,both,3]
       -k/--connector_id <KEY>      : Connector_id[String]
       -j/--json_params <JSON>      : JSON format parameters[String]
       -i/--input_json_file <FILE>  : JSON format input file path[String]
       -o/--output_json_file <FILE> : JSON format output file path[String]
                                    : This file will be created only when settings are found.
          --dump                    : Dump all settings
       -h/--help                    : Print this message
```

- When checking **LLDMA** channel settings:
	- Required arguments
		- `connector_id`
- When checking **PTU** sessions:
	- Required arguments
		- `device`
		- `lane`
		- `extif_id`
		- `cid`
- When checking **Chain** settings:
	- Required arguments
		- `device`
		- `lane`
		- `fchid`
	- Optional arguments
		- `dir`
- `External interface id(extif_id)` depends on the FPGA implementation.
	- In the current FPGA implementation: `0`=`LLDMA`, `1`=`PTU`.


## Output example
- Check LLDMA setting : Not found
	```sh
	$ ./bin/fpga-chk-connection -k RX-0
	# Extif Settings
	- Dataflow-session : 0
	        - Parameters
	                - Extif_id : 0(LLDMA)
	                - Connector_id : RX-0
	        - Result : Not found

	# Chain Settings

	# Summary
	- Found Dataflow Index
	        - Dataflow-session : []
	        - Dataflow-chain : []
	```

- Check LLDMA setting : Found
	```sh
	$ ./bin/fpga-chk-connection -k RX-0
	# Extif Settings
	- Dataflow-session : 0
	        - Parameters
	                - Extif_id : 0(LLDMA)
	                - Connector_id : RX-0
	        - Result : Found
	        - Found Dataflow Settings
	                - Device : /dev/xpcie_21320621V01M
	                - Direction : RX(Ingress)
	                - Connection_id : 0

	# Chain Settings

	# Summary
	- Found Dataflow Index
	        - Dataflow-session : [0]
	        - Dataflow-chain : []
	```

- Check PTU setting : Not found
	```sh
	$ ./bin/fpga-chk-connection -d /dev/xpcie_21320621V01M -l 0 -e 1 -c 2
	# Extif Settings
	- Dataflow-session : 0
	        - Parameters
	                - Device : /dev/xpcie_21320621V01M
	                - Lane : 0
	                - Extif_id : 1(PTU)
	                - Connection_id : 2
	        - Result : Not found

	# Chain Settings

	# Summary
	- Found Dataflow Index
	        - Dataflow-session : []
	        - Dataflow-chain : []
	```

- Check PTU setting : Found
	```sh
	$ ./bin/fpga-chk-connection -d /dev/xpcie_21320621V01M -l 0 -e 1 -c 2
	# Extif Settings
	- Dataflow-session : 0
	        - Parameters
	                - Device : /dev/xpcie_21320621V01M
	                - Lane : 0
	                - Extif_id : 1(PTU)
	                - Connection_id : 2
	        - Result : Found

	# Chain Settings

	# Summary
	- Found Dataflow Index
	        - Dataflow-session : [0]
	        - Dataflow-chain : []
	```

- Check Chain setting : Not found
	```sh
	$ ./bin/fpga-chk-connection -d /dev/xpcie_21320621V01M -l 0 -f 0
	# Extif Settings

	# Chain Settings
	- Dataflow-chain : 0
	        - Parameters
	                - Device : /dev/xpcie_21320621V01M
	                - Lane : 0
	                - Function_channel_id : 0
	                - Direction : 1(Ingress)
	        - Result : Not found
	- Dataflow-chain : 1
	        - Parameters
	                - Device : /dev/xpcie_21320621V01M
	                - Lane : 0
	                - Function_channel_id : 0
	                - Direction : 2(Egress)
	        - Result : Not found

	# Summary
	- Found Dataflow Index
	        - Dataflow-session : []
	        - Dataflow-chain : []
	```

- Check Chain setting : Found
	```sh
	$ ./bin/fpga-chk-connection -d /dev/xpcie_21320621V01M -l 0 -f 0
	# Extif Settings

	# Chain Settings
	- Dataflow-chain : 0
	        - Parameters
	                - Device : /dev/xpcie_21320621V01M
	                - Lane : 0
	                - Function_channel_id : 0
	                - Direction : 1(Ingress)
	        - Result : Found
	                - Extif_id : 0(LLDMA)
	                - Connection_id : 0
	- Dataflow-chain : 1
	        - Parameters
	                - Device : /dev/xpcie_21320621V01M
	                - Lane : 0
	                - Function_channel_id : 0
	                - Direction : 2(Egress)
	        - Result : Found
	                - Extif_id : 0(LLDMA)
	                - Connection_id : 0

	# Summary
	- Found Dataflow Index
	        - Dataflow-session : []
	        - Dataflow-chain : [0,1]
	```

- Check LLDMA->Func->LLDMA setting : Found
	```sh
	$ ./bin/fpga-chk-connection -k RX-0 -d /dev/xpcie_21320621V01M -l 0 -f 0 --dir both -k TX-0
	# Extif Settings
	- Dataflow-session : 0
	        - Parameters
	                - Extif_id : 0(LLDMA)
	                - Connector_id : RX-0
	        - Result : Found
	        - Found Dataflow Settings
	                - Device : /dev/xpcie_21320621V01M
	                - Direction : RX(Ingress)
	                - Connection_id : 0
	- Dataflow-session : 1
	        - Parameters
	                - Extif_id : 0(LLDMA)
	                - Connector_id : TX-0
	        - Result : Found
	        - Found Dataflow Settings
	                - Device : /dev/xpcie_21320621V01M
	                - Direction : TX(Egress)
	                - Connection_id : 0

	# Chain Settings
	- Dataflow-chain : 0
	        - Parameters
	                - Device : /dev/xpcie_21320621V01M
	                - Lane : 0
	                - Function_channel_id : 0
	                - Direction : 1(Ingress)
	        - Result : Found
	                - Extif_id : 0(LLDMA)
	                - Connection_id : 0
	- Dataflow-chain : 1
	        - Parameters
	                - Device : /dev/xpcie_21320621V01M
	                - Lane : 0
	                - Function_channel_id : 0
	                - Direction : 2(Egress)
	        - Result : Found
	                - Extif_id : 0(LLDMA)
	                - Connection_id : 0

	# Summary
	- Found Dataflow Index
	        - Dataflow-session : [0,1]
	        - Dataflow-chain : [0,1]
	```

- Dump All setting : Not found
	```sh
	$ ./bin/fpga-chk-connection --dump
	# Extif Settings
	- LLDMA channels(21320621V01M : RX : Avail)  : [0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15]
	- LLDMA channels(21320621V01M : RX : Active) : []
	- LLDMA channels(21320621V01M : TX : Avail)  : [0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15]
	- LLDMA channels(21320621V01M : TX : Active) : []

	# Chain Settings
	```

- Dump All setting : Found
	```sh
	$ ./bin/fpga-chk-connection --dump
	# Extif Settings
	- LLDMA channels(21320621V01M : RX : Avail)  : [0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15]
	- LLDMA channels(21320621V01M : RX : Active) : [0,1,2,3]
	- LLDMA channels(21320621V01M : TX : Avail)  : [0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15]
	- LLDMA channels(21320621V01M : TX : Active) : [0,1,2,3]
	- LLDMA-session : 0
	        - Device : /dev/xpcie_21320621V01M
	        - Lane : 0
	        - Extif_id : 0(LLDMA)
	        - Connection_id : 0
	- LLDMA-session : 1
	        - Device : /dev/xpcie_21320621V01M
	        - Lane : 0
	        - Extif_id : 0(LLDMA)
	        - Connection_id : 1
	- LLDMA-session : 2
	        - Device : /dev/xpcie_21320621V01M
	        - Lane : 0
	        - Extif_id : 0(LLDMA)
	        - Connection_id : 2
	- LLDMA-session : 3
	        - Device : /dev/xpcie_21320621V01M
	        - Lane : 0
	        - Extif_id : 0(LLDMA)
	        - Connection_id : 3

	# Chain Settings
	- Chain-settings : 0
	        - Device : /dev/xpcie_21320621V01M
	        - Lane : 0
	        - Function_channel_id : 0
	        - Direction : 1(Ingress)
	        - Extif_id : 0(LLDMA)
	        - Connection_id : 0
	- Chain-settings : 1
	        - Device : /dev/xpcie_21320621V01M
	        - Lane : 0
	        - Function_channel_id : 0
	        - Direction : 2(Egress)
	        - Extif_id : 0(LLDMA)
	        - Connection_id : 0
	- Chain-settings : 2
	        - Device : /dev/xpcie_21320621V01M
	        - Lane : 0
	        - Function_channel_id : 1
	        - Direction : 1(Ingress)
	        - Extif_id : 0(LLDMA)
	        - Connection_id : 1
	- Chain-settings : 3
	        - Device : /dev/xpcie_21320621V01M
	        - Lane : 0
	        - Function_channel_id : 1
	        - Direction : 2(Egress)
	        - Extif_id : 0(LLDMA)
	        - Connection_id : 1
	- Chain-settings : 4
	        - Device : /dev/xpcie_21320621V01M
	        - Lane : 0
	        - Function_channel_id : 2
	        - Direction : 1(Ingress)
	        - Extif_id : 0(LLDMA)
	        - Connection_id : 2
	- Chain-settings : 5
	        - Device : /dev/xpcie_21320621V01M
	        - Lane : 0
	        - Function_channel_id : 2
	        - Direction : 2(Egress)
	        - Extif_id : 0(LLDMA)
	        - Connection_id : 2
	- Chain-settings : 6
	        - Device : /dev/xpcie_21320621V01M
	        - Lane : 0
	        - Function_channel_id : 3
	        - Direction : 1(Ingress)
	        - Extif_id : 0(LLDMA)
	        - Connection_id : 3
	- Chain-settings : 7
	        - Device : /dev/xpcie_21320621V01M
	        - Lane : 0
	        - Function_channel_id : 3
	        - Direction : 2(Egress)
	        - Extif_id : 0(LLDMA)
	        - Connection_id : 3
	```

## Files
```
fpga-chk-connection
 +---- Makefile
 +---- README.md
 +---- 3rdparty/
 +---- include/
 |     +---- fcc_arg.h
 |     +---- fcc_check.h
 |     +---- fcc_json.h
 |     +---- fcc_log.h
 |     \---- fcc_prm.h
 +---- src/
 |     +---- fcc_arg.c
 |     +---- fcc_check.c
 |     +---- fcc_json.c
 |     +---- fcc_log.c
 |     +---- fcc_prm.c
 |     \---- fcc_main.c
 \---- test/
```

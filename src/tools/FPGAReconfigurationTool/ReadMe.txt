FPGA Reconfiguration Tool

This tool writes bitstream and configuration information to the FPGA device.

Writing FPGA Bitstream (child bs) to initialized FPGA
	Write FPGA Bitstreams (child bs) to Initialized FPGA
FPGA Reset
	Places the FPGA in use in the initialized state.
Child Bitstream Reset
	Reset an FPGA Bitstream (child bs) in use to the state it was just written to.

[Command input format]

FPGAReconfigurationTool nodename devicefilepath [-l0 configname] [-l1 configname] [-reset FPGA|ChildBs]

The first argument specifies the node name of the host where the FPGA device is deployed.
The second argument specifies the file path of the FPGA device.
The third argument can be followed by the following optional arguments.
	The -lx argument specifies the ConfigMap name of the FPGAFunc configuration information to be set in lane.
		*The configuration information for FPGAFunc must be created in advance as a ConfigMap that 
            contains information (parameters to be set for the functions, etc.) 
            about the functions to be executed in each Lane of the FPGA .
		x indicates a lane index. The current target FPGA Bitstream (child bs) has a 2-lane configuration,
        so when writing FPGA Bitstream (child bs) to the initialized FPGA, the -l0 and -l1 lane arguments must be specified.
	The -reset argument can indicate FPGA: FPGA reset, ChildBs: child Bitstream reset.
		Do not specify the -reset argument when writing FPGA Bitstream (child bs) to an initialized FPGA.

[Example of command input]

・FPGA Bitstream (child bs) write to an initialized FPGA (example of writing different function information on lane0 and lane1)
	FPGAReconfigurationTool workernode /dev/xpcie_XXXXXXX -l0 fpgafunc-config-filter-resize-high-infer -l1 fpgafunc-config-filter-resize-low-infer

・FPGA Reset
	FPGAReconfigurationTool workernode /dev/xpcie_XXXXXXX -reset FPGA

・Child Bitstream Reset
	FPGAReconfigurationTool workernode /dev/xpcie_XXXXXXX -reset ChildBs
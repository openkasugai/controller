FPGA Clear Check Tool

This tool is a tool to confirm that various resources in the FPGA used by DataFlow to be deleted when deleting DataFlow have been cleared.
After DataFlow is deployed, obtain FPGA information for the various CR associated with the deployed DataFlow CR.
After DataFlow is deleted, LLDMA information and CHAIN information are acquired for the FPGA device from the FPGA information of the deleted DataFlow CR, and the current clear status is displayed.

[Command input format]

./FPGACheckPerDF get|check [-n Namespace] [-node Nodename]

The first argument specifies get or check.
	When get, FPGA information is obtained from various CR associated with all DataFlow CR deployed on the Namespace.	
	When check, from the FPGA information acquired during get, the LLDMA information and CHAIN information are acquired for the FPGA device 
		using the FPGA information of DataFlow CR (deleted DataFlow CR), which does not exist at this time, and the current clear status is displayed.
The second and subsequent arguments are optional, and the following optional arguments can be specified. 
	The -n argument must be the Namespace where DataFlow CR is deployed.
		*By default, Namespace is set to "default."
	For the -node argument, specify the name of the node on which the FPGA device is installed.
		This is the node name in the k8s cluster, specifically the node name displayed by "kubectl get node."
		*By default, the environment variable "K8S_NODENAME" is acquired.
		 If the environment variable "K8S_NODENAME" is not set, this command displays a message and then terminates.

[Example of command input]

・get
	./FPGACheckPerDF get -n Namespace -node Nodename

	Execution result
		Namespace:test01  Name:df-test-1-1-1  XXXXXXX

	Description of what is displayed
		Namespace: Namespace
		Name	 : Name of DataFlow CR
		XXXXXXX  : State of DataFlow CR
			OK	 : Each CR's information associated with DataFlow CR is valid
			NG	 : Each CR's information associated with DataFlow CR is invalid
					*Information is missing in CR chain up to the wb-end-of-chain
			NotFPGA: Each CR tied to DataFlow CR does not contain an FPGA
			OtherNode: The Nodename of the FPGA associated with DataFlow CR does not match the Nodename specification in the argument.
		*DataFlow CR whose DataFlow CR status is "OK" will be checked during a post-deletion check.
		 DataFlow CR in any other state is excluded from the check after deletion.

・check
	./FPGACheckPerDF check -n Namespace -node Nodename

	Execution result
		Namespace:test01  Name:df-test-1-1-1
		  #LLDMA ingress Command-Result:../fpga-chk-connection/fpga-chk-connection -k test01-df-test-1-1-1-wbfunction-decode-main
		    - Result : XXXXXXXXX
		  #CHAIN ingress Command-Result:../fpga-chk-connection/fpga-chk-connection -d 21320621V00D -l 0 -f 1 --dir ingress
		    - Result : XXXXXXXXX
		  #CHAIN egress  Command-Result:../fpga-chk-connection/fpga-chk-connection -d 21320621V00D -l 0 -f 1 --dir egress
		    - Result : XXXXXXXXX
		  #LLDMA egress  Command-Result:../fpga-chk-connection/fpga-chk-connection -k test01-df-test-1-1-1-wbfunction-low-infer-main
		    - Result : XXXXXXXXX

	Description of what is displayed
		Namespace: Namespace
		Name	 : Name of DataFlow CR
		  #LLDMA ingress Command-Result : Command to get LLDMA (ingress) status
		    - Result : XXXXXXXXX: Result of command to get LLDMA (ingress) status (Result line only)
						Not found: LLDMA (ingress) cleared
						Found	: LLDMA (ingress) not cleared
		  #CHAIN ingress Command-Result : Command to get CHAIN (ingress) status
		    - Result : XXXXXXXXX: Result of command to get CHAIN (ingress) status (Result line only)
						Not found: CHAIN (ingress) cleared
						Found	: CHAIN (ingress) not cleared
		  #CHAIN egress  Command-Result : Command to get CHAIN (egress) status
		    - Result : XXXXXXXXX: Result of command to get CHAIN (egress) status (Result line only)
						Not found: CHAIN (egress) cleared
						Found	: CHAIN (egress) not cleared
		  #LLDMA egress  Command-Result : Command to get LLDMA (egress) status
		    - Result : XXXXXXXXX: Result of command to get LLDMA (egress) status (Result line only)
						Not found: LLDMA (egress) cleared
						Found	: LLDMA (egress) not cleared

		*If all the results (Result) of each status acquisition command are "Not found," 
		 it can be judged that all the resources of the FPGA used by the relevant DataFlow CR have been cleared.


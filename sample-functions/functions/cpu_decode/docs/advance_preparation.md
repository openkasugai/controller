###### Copyright 2024 NTT Corporation, FUJITSU LIMITED

## Advance preparation

### Hugepage

As an operating condition, the server must be configured with Hugepage.
Here is an example of checking the Hugepage configuration on ubuntu:

- How to check

  ```
  $ cat /proc/meminfo
  ...
  HugePages_Total:      16
  HugePages_Free:       16
  HugePages_Rsvd:        0
  HugePages_Surp:        0
  Hugepagesize:    1048576 kB      <- The hugepage size is 1GB.
  ```

- How to set

  1. Edit `/etc/default/grub`

     - Before Editing

       ```
       GRUB_CMDLINE_LINUX_DEFAULT=""
       ```

     - After Editing

       ```
       GRUB_CMDLINE_LINUX_DEFAULT="default_hugepagesz=1G hugepagesz=1G hugepages=16"
       ```
       Add hugepage-related boot options to GRUB_CMDLINE_LINUX_DEFAULT.
       With the above setting, when hugepagesz=1G, 16 pages are allocated at startup.

  2. Update Grub and reboot.

     ```
     $ update-grub && reboot 
     ```


### FPGA driver/library

- Get the FPGA driver/library

  ```
  $ cd {Repository root directory}
  $ git config -f .gitmodules submodule.src/submodules/fpga-software.url https://github.com/openkasugai/hardware-drivers.git
  $ git submodule sync
  $ git submodule update --init --recursive
  ```

- Build and Install the FPGA xpcie Driver

  https://github.com/openkasugai/hardware-drivers/blob/main/driver/README.md

- Build the FPGA library (libfpga)

  https://github.com/openkasugai/hardware-drivers/blob/main/lib/README.md


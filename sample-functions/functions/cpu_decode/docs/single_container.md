###### Copyright 2024 NTT Corporation, FUJITSU LIMITED

# Using the cpu_decode alone to execute FPGA DMA (For debug)

The cpu_decode performs the initial FPGA settings and executes both DMA transmission and reception.
The subsequent FPGA application will become unnecessary.


## Makefile settings

Perform a container build with the following Makefile settings.

```
APPLOG_PRINT := 1

MODULE_FPGA := 1

DEBUG_PRINT := 0

　
DPDK_SECONDARY_PROC_MODE := 0

　
ALLOCATE_SRC_SHMEM_MODE := 1

EXEC_FPGA_ENQUEUE_MODE := 1

EXEC_FPGA_DEQUEUE_MODE := 1

　
CONTROL_FPGA_DEV_INIT := 1

CONTROL_FPGA_FUNC_INIT := 1

CONTROL_FPGA_ENQUEUE_LLDMA_INIT := 1

CONTROL_FPGA_CHAIN_CONNECT := 1

CONTROL_FPGA_FINISH := 1

CONTROL_FPGA_FUNC_FINISH := 1

CONTROL_FPGA_SHMEM_FINISH := 1

　
DEBUG_PRINT_FPGA_OUT_HEADER := 0

DEBUG_FPGA_OUT_IMAGE_TO_MP4 := 1
```


## Save the mp4 video file

When using only the cpu_decode to execute FPGA DMA, the FPGA output frame can be saved as an mp4 video file.

If you build with `DEBUG_FPGA_OUT_IMAGE_TO_MP4 := 1`, save the mp4 video file.

Video file name : fpga_out_frame_[YYYYMMDDhhmmss].mp4


## Header log added by FPGA

FPGA added headers can be checked.

Log file name：recv_header_[YYYYMMDD-hhmmss].log

The header format is different between Non-Module FPGAs and Module FPGAs.

- Non-Module FPGA Header Example

  ```
  --------------------------------------------------------
  FrameHeader CH(0) TASK(1)
                  Receive value | Expected value | compare
    marker      :    0xe0ff10ad |     0xe0ff10ad |   OK
    payload_len :    0x004b0000 |     0x004b0000 |   OK
    payload_type:          0x01 |           0x01 |   OK
    channel_id  :        0x0000 |         0x0000 |   OK
    frame_index :    0x00000000 |              - |   -
    color_space :          0x00 |           0x00 |   OK
    data_type   :          0x00 |           0x00 |   OK
    num_ch      :        0x0000 |         0x0000 |   OK
    width       :        0x0500 |         0x0500 |   OK
    height      :        0x0500 |         0x0500 |   OK
  ```

- Module FPGA Header Example

  ```
  ----------------------------------------------------------------
  FrameHeader CH(0) TASK(1)
                          Receive value | Expected value | compare
    marker         :         0xe0ff10ad |     0xe0ff10ad |   OK
    payload_len    :         0x004b0000 |     0x004b0000 |   OK
    sequence_num   :         0x00000002 |              - |   -
    timestamp      : 0x0000000000000000 |              - |   -
    data_id        :         0x00000000 |     0x00000000 |   OK
    header_checksum:             0x0000 |              - |   -
  ```

Receive value  : Header value of the data output from the FPGA to the destination shared memory.

Expected value : Expected value of the header.

compare        : Result of comparing output and expected values.


## Execution Order

Set DECENV_OUTDST_PROTOCOL=`"DMA"` and execute cpu_decode.

1. Start the cpu_decode.

   It will remain in standby mode until it connects to the video distribution server.

2. Start the video distribution server.

   Establish a connection between the video distribution server and the cpu_decode, and start video distribution.


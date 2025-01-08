###### Copyright 2024 NTT Corporation, FUJITSU LIMITED

# CPU Decode Application [cpu_decode]

## Overview

- Decodes H.264 streaming and sends frames using FPGA DMA or TCP.
  - H.264 Streaming (RTP or RTSP) -> CPU Decode -> Sent using FPGA DMA
  - H.264 Streaming (RTP or RTSP) -> CPU Decode -> Sent using TCP
- The output frame format is header + data frame.

## Requirement

- FPGA library (libfpga)

  https://github.com/openkasugai/hardware-drivers/tree/main/lib

- FPGA xpcie driver

  https://github.com/openkasugai/hardware-drivers/tree/main/driver

- OpenCV 3.4.3 (with Gstreamer)

  [docker/install_opencv_for_container.sh](docker/install_opencv_for_container.sh)


## Advance preparation

- [Advance preparation](docs/advance_preparation.md)


## Build

### Makefile settings

The default build configuration is for running controller.

[docker/Makefile](docker/Makefile)


|  Build options | Target |  Desctiption  |
| :---- | :---- | :---- |
| APPLOG_PRINT                | Common | When set to '1', the application log is output to standard output. |
| MODULE_FPGA                 | Common | When set to '1', it corresponds to the Module FPGA header format. When using FPGA, the FPGA API also corresponds to the Module FPGA. |
| DEBUG_PRINT                 | Common | When set to '1', the execution progress (frame count for each thread) is output to standard output. |
| DPDK_SECONDARY_PROC_MODE    | FPGA DMA | When set to '1', it becomes the DPDK secondary process.<br>When set to '0', it becomes the DPDK primary process. |
| ALLOCATE_SRC_SHMEM_MODE     | FPGA DMA | When set to '1', the cpu_decode application uses DPDK to allocate shared memory for the FPGA SRC.<br>When set to '0', the cpu_decode application will not allocate shared memory for the FPGA SRC, but will instead refer to the shared memory address given by the environment variable parameter. |
| EXEC_FPGA_ENQUEUE_MODE      | FPGA DMA | When set to '1', the cpu_decode application will put an ENQUEUE DMA request to the FPGA.<br>When set to '0', the cpu_decode application will not put ENQUEUE DMA requests to the FPGA. |
| EXEC_FPGA_DEQUEUE_MODE      | FPGA DMA | When set to '1', the cpu_decode application will put an DEQUEUE DMA request to the FPGA. It also allocates shared memory for FPGA DST. (Using DPDK).<br>When set to '0', the cpu_decode application will not put DEQUEUE DMA requests to the FPGA. |
| CONTROL_FPGA_DEV_INIT       | FPGA DMA | When set to '1', the cpu_decode application will perform FPGA device settings.<br>When set to '0', the cpu_decode application will not perform FPGA device settings. |
| CONTROL_FPGA_FUNC_INIT      | FPGA DMA | When set to '1', the cpu_decode application will set the FPGA function.<br>When set to '0', the cpu_decode application will not perform FPGA function settings. |
| CONTROL_FPGA_ENQUEUE_LLDMA_INIT | FPGA DMA | When set to '1', the cpu_decode application will set the FPGA ENQUEUE DMA (DMA_HOST_TO_DEV) channel.<br>When set to '0', the cpu_decode application will not set the FPGA ENQUEUE DMA (DMA_HOST_TO_DEV) channel. |
| CONTROL_FPGA_CHAIN_CONNECT  | FPGA DMA | When set to '1', the cpu_decode application will perform the FPGA function chain connection settings.<br>When set to '0', the cpu_decode application will not perform the FPGA function chain connection settings. |
| CONTROL_FPGA_FINISH         | FPGA DMA | Only valid when CONTROL_FPGA_DEV_INIT=1.<br>When set to '1', the cpu_decode application will terminate the FPGA device.<br>When set to '0', the cpu_decode application does not terminate the FPGA device. |
| CONTROL_FPGA_FUNC_FINISH    | FPGA DMA | Only valid when CONTROL_FPGA_FUNC_INIT=1.<br>When set to '1', the cpu_decode application will terminate the FPGA function.<br>When set to '0', the cpu_decode application will not terminate the FPGA function.|
| CONTROL_FPGA_SHMEM_FINISH   | FPGA DMA | When set to '1', the cpu_decode application will terminate the DPDK process that allocates the FPGA shared memory.<br>When set to '0', the cpu_decode application does not terminate the DPDK process that allocates FPGA shared memory. |
| DEBUG_PRINT_FPGA_OUT_HEADER | FPGA DMA | Only valid when EXEC_FPGA_DEQUEUE_MODE=1.<br>When set to '1', the header information of the FPGA output data is output to standard output. |
| DEBUG_FPGA_OUT_IMAGE_TO_MP4 | FPGA DMA | Only valid when EXEC_FPGA_DEQUEUE_MODE=1.<br>When set to '1', save the FPGA output frame as an mp4 video file (fpga_out_frame_[YYYYMMDDhhmmss].mp4). |

### Building container image

```sh
$ cd docker/
```

- When creating a container image using "buildah bud"

  Execute [buildah_bud.sh](docker/buildah_bud.sh).

  ```sh
  $ ./buildah_bud.sh [image tag version]
  ```

  Ex.)
  ```sh
  $ ./buildah_bud.sh 1.0.0
  ```

- When creating a container image using "docker build"

  Execute [docker_build.sh](docker/docker_build.sh).

  ```sh
  $ ./docker_build.sh [image tag version]
  ```

  Ex.)
  ```sh
  $ ./docker_build.sh 1.0.0
  ```

### Downloading container image

You can also download container image from ghcr.
```sh
$ sudo buildah pull ghcr.io/openkasugai/controller/cpu_decode:1.0.0
$ sudo buildah tag ghcr.io/openkasugai/controller/cpu_decode:1.0.0 localhost/cpu_decode:1.0.0
```

## Parameter settings for the cpu_decode

The cpu_decode execution parameters are set using environment variables.

- Common parameters
  | Environment variable | Target | Required/Optional | Desctiption |
  | :--- | :---: | :--- | :--- |
  | DECENV_APPLOG_LEVEL          | Common | Optional | Set the log level.<br>0 : No log<br>2 : ERROR (ERROR only)<br>3 : WARN (ERROR and WARN)<br>4 : INFO (ERROR,WARN and INFO). This is the default.<br>5 : DEBUG (ERROR,WARN,INFO and DEBUG)<br>6 : ALL (All levels) |
  | DECENV_VIDEOSRC_PROTOCOL     | Common | Required | Set the input streaming video protocol ("RTP" or "RTSP"). |
  | DECENV_VIDEOSRC_PORT         | Common | Required | Set the IPv4 port for input streaming video. |
  | DECENV_FRAME_FPS             | Common | Required | Set the input FPS. |
  | DECENV_FRAME_WIDTH           | Common | Required | Set the width of the input frame size. |
  | DECENV_FRAME_HEIGHT          | Common | Required | Set the height of the input frame size. |
  | DECENV_OUTDST_PROTOCOL       | Common | Required | Set the output protocol ("DMA" or "TCP").<br>DMA : FPGA DMA output.<br>TCP : TCP network output. |
  | DECENV_VIDEO_CONNECT_LIMIT   | Common | Optional | Set a limit on the number of connections to video distribution.<br>'0' means unlimited. This is the default.<br>If it is greater than '0', it will end after the specified number of connections. |

- Parameters for RTSP video streaming servers [Only used when DECENV_VIDEOSRC_PROTOCOL=`"RTSP"` is set.]
  | Environment variable | Target | Required/Optional | Desctiption |
  | :--- | :---: | :--- | :--- |
  | DECENV_VIDEOSRC_IPA          | Common | Required | Set the IPv4 addres for RTSP server. |

- FPGA DMA parameters [Only used when DECENV_OUTDST_PROTOCOL=`"DMA"` is set.]
  | Environment variable | Target | Required/Optional | Desctiption |
  | :--- | :---: | :--- | :--- |
  | DECENV_DPDK_FILE_PREFIX      | FPGA DMA | Optional | Set the DPDK file prefix.<br>Used when DPDK_SECONDARY_PROC_MODE=1. |
  | DECENV_FPGA_DMA_HOST_TO_DEV_CONNECTOR_ID | FPGA DMA | Required | Set the FPGA DMA_HOST_TO_DEV connector ID.<br>If CONTROL_FPGA_ENQUEUE_LLDMA_INIT=1, the cpu_decode app will set the connector ID using this setting name.<br>If CONTROL_FPGA_ENQUEUE_LLDMA_INIT=0, the connector ID set by the external process is used. |
  | DECENV_FPGA_SRC_SHMEM_ADDR   | FPGA DMA | Optional | Set the shared memory address to which the decoded frame is sent. (FPGA SRC shared memory).<br>When ALLOCATE_SRC_SHMEM_MODE=0, this address is referenced. |
  | DECENV_FPGA_DEV_NAME         | FPGA DMA | Optional | Set the FPGA device file path.<br>Used when CONTROL_FPGA_DEV_INIT=1 or EXEC_FPGA_DEQUEUE_MODE=1. |
  | DECENV_FPGA_CH_ID            | FPGA DMA | Optional | Set the Channel ID for the FPGA.<br>Used when CONTROL_FPGA_ENQUEUE_LLDMA_INIT=1 or EXEC_FPGA_DEQUEUE_MODE=1. |
  | DECENV_FPGA_OUT_FRAME_WIDTH  | FPGA DMA | Optional | Set the width of the FPGA output frame size.<br>Used when CONTROL_FPGA_FUNC_INIT=1 or EXEC_FPGA_DEQUEUE_MODE=1. |
  | DECENV_FPGA_OUT_FRAME_HEIGHT | FPGA DMA | Optional | Set the height of the FPGA output frame size.<br>Used when CONTROL_FPGA_FUNC_INIT=1 or EXEC_FPGA_DEQUEUE_MODE=1. |

- TCP output parameters [Only used when DECENV_OUTDST_PROTOCOL=`"TCP"` is set.]
  | Environment variable | Target | Required/Optional | Desctiption |
  | :--- | :---: | :--- | :--- |
  | DECENV_OUTDST_IPA            | TCP output | Required | Set the TCP output destination IP address. |
  | DECENV_OUTDST_PORT           | TCP output | Required | Set the TCP output destination IP port. |


## The cpu_decode processing overview

### When the output is FPGA DMA. [DECENV_OUTDST_PROTOCOL=`"DMA"`]

1. Receiving streaming video distributed from a video server.
2. H.264 decoding.
3. Header is added to the decode frame and stored in the FPGA SRC shared memory.
4. Execution of FPGA DMA transfer request.

### When the output is TCP. [DECENV_OUTDST_PROTOCOL=`"TCP"`]

1. Receiving streaming video distributed from a video server.
2. H.264 decoding.
3. Header is added to the decode frame.
4. TCP send.


## Container execution sample

### The cpu_decode container creation

Sample scripts for "docker run" for each input/output pattern.

- input:RTP, output:FPGA DMA

  [docker/docker_run.rtp_to_fpga_dma.sh](docker/docker_run.rtp_to_fpga_dma.sh)

- input:RTP, output:TCP

  [docker/docker_run.rtp_to_tcp.sh](docker/docker_run.rtp_to_tcp.sh)

- input:RTSP, output:FPGA DMA

  [docker/docker_run.rtsp_to_fpga_dma.sh](docker/docker_run.rtsp_to_fpga_dma.sh)

- input:RTSP, output:FPGA DMA

  [docker/docker_run.rtsp_to_tcp.sh](docker/docker_run.rtsp_to_tcp.sh)

### Execute the cpu_decode on the container

Sample script for "docker exec" to execute the cpu_decode.

[docker/docker_exec.sh](docker/docker_exec.sh)


## Execution Order

### When the output is FPGA DMA. [DECENV_OUTDST_PROTOCOL=`"DMA"`]

1. Start the application for the FPGA DMA destination.

2. Start the cpu_decode.

   It will remain in standby mode until it connects to the video distribution server.

3. Start the video distribution server.

   Establish a connection between the video distribution server and the cpu_decode, and start video distribution.


### When the output is TCP. [DECENV_OUTDST_PROTOCOL=`"TCP"`]

1. Start the TCP receiver application.

   The TCP connection is in a waiting state.

2. Start the cpu_decode.

   Start connecting to the TCP receiver application.
   If the TCP receiver application cannot be connected to, it will retry to connect every second, and if it reaches 60 retries, it will terminate with an error.

   After a connection is established with the TCP receiver, it will wait until it is connected to the video distribution server.
   
3. Start the video distribution server.

   Establish a connection between the video distribution server and the cpu_decode, and start video distribution.


## The cpu_decode log

Log file name ：app_[YYYYMMDD-hhmmss].log

If you build with `APPLOG_PRINT := 1`, standard output is also enabled.

- Sample log when the output is FPGA DMA. [DECENV_OUTDST_PROTOCOL=`"DMA"`]

  ```
  log start...20241119-102843
  loglevel: 4
  10:28:43 Version: 0.6.00
  10:28:43 [info ] parameter VIDEOSRC_PROTOCOL: "RTP"
  10:28:43 [info ] parameter videosrc: "udpsrc port=5001 buffer-size=512000000 caps=application/x-rtp ! rtpjitterbuffer ! rtph264depay ! h264parse ! openh264dec ! queue ! videoconver
  t ! appsink"
  10:28:43 [info ] parameter FRAME_FPS: 15.000000
  10:28:43 [info ] parameter FRAME_WIDTH: 3840
  10:28:43 [info ] parameter FRAME_HEIGHT: 2160
  10:28:43 [info ] parameter VIDEO_CONNECT_LIMIT: 0
  10:28:43 [info ] parameter OUTDST_PROTOCOL: "DMA"
  10:28:43 [info ] parameter DPDK_FILE_PREFIX: "test01-df-test-1-1-1-1-wbfunction-decode-main"
  10:28:43 [info ] parameter FPGA_DMA_HOST_TO_DEV_CONNECTOR_ID: "test01-df-test-1-1-1-1-wbfunction-decode-main"
  10:28:43 [info ] parameter FPGA_DEV_NAME: "/dev/xpcie_2133072BM03P"
  10:28:43 [info ] PID(2) CPU affinity(0xffffffffffffffff)
  10:28:43 [info ] --- fpga_shmem_alloc for srcbuf ---
  10:28:43 [info ] --- enqueue fpga_lldma_queue_setup ---
  10:28:43 [info ] parameter FPGA_CH_ID: 0
  10:28:43 [info ] --- send_fpga_enq_thread start ---
  10:28:43 [info ] --- send_fpga_deq_thread start ---
  10:28:43 [info ] CH(0) DMA RX dma_info: dir(0) chid(0) queue_addr(0x7f7ae67b4000) queue_size(255)
  10:31:45 [info ] VideoCapture opened.
  10:31:45 [info ] --- decode_to_dma_thread start ---
  10:31:45 [info ] CH(0) DMA RX dmacmd_info: enq(0) task_id(1) src_len(24883264) src_addr(0x17e7f8c00)
  10:31:45 [info ] CH(0) DMA RX dmacmd_info: enq(0) result_task_id(1) result_status(0) result_data_len(0)
  10:31:45 [info ] CH(0) DMA RX dmacmd_info: enq(1) task_id(2) src_len(24883264) src_addr(0x17e7f8c00)
  10:31:45 [info ] CH(0) DMA RX dmacmd_info: enq(1) result_task_id(2) result_status(0) result_data_len(0)
  10:31:45 [info ] CH(0) DMA RX dmacmd_info: enq(2) task_id(3) src_len(24883264) src_addr(0x17e7f8c00)
  10:31:45 [info ] CH(0) DMA RX dmacmd_info: enq(2) result_task_id(3) result_status(0) result_data_len(0)
  10:31:45 [info ] CH(0) DMA RX dmacmd_info: enq(3) task_id(4) src_len(24883264) src_addr(0x17e7f8c00)
  10:31:45 [info ] CH(0) DMA RX dmacmd_info: enq(3) result_task_id(4) result_status(0) result_data_len(0)
  10:31:45 [info ] CH(0) DMA RX dmacmd_info: enq(4) task_id(5) src_len(24883264) src_addr(0x17e7f8c00)
  10:31:45 [info ] CH(0) DMA RX dmacmd_info: enq(4) result_task_id(5) result_status(0) result_data_len(0)
  10:31:46 [info ] CH(0) DMA RX dmacmd_info: enq(5) task_id(6) src_len(24883264) src_addr(0x17e7f8c00)
  10:31:46 [info ] CH(0) DMA RX dmacmd_info: enq(5) result_task_id(6) result_status(0) result_data_len(0)
  ```

- Sample log when the output is TCP. [DECENV_OUTDST_PROTOCOL=`"TCP"`]

  ```
  log start...20241120-022446
  loglevel: 6
  02:24:46 Version: 0.6.00
  02:24:46 [info ] parameter VIDEOSRC_PROTOCOL: "RTP"
  02:24:46 [info ] parameter videosrc: "udpsrc port=5004 buffer-size=512000000 caps=application/x-rtp ! rtpjitterbuffer ! rtph264depay ! h264parse ! openh264dec ! queue ! videoconvert ! appsink"
  02:24:46 [info ] parameter FRAME_FPS: 15.000000
  02:24:46 [info ] parameter FRAME_WIDTH: 3840
  02:24:46 [info ] parameter FRAME_HEIGHT: 2160
  02:24:46 [info ] parameter VIDEO_CONNECT_LIMIT: 0
  02:24:46 [info ] parameter OUTDST_PROTOCOL: "TCP"
  02:24:46 [info ] parameter OUTDST_IPA: 192.174.90.111
  02:24:46 [info ] parameter OUTDST_PORT: 15000
  02:24:46 [info ] parameter FPGA_CH_ID: 0
  02:24:46 [info ] TCP connect try 1. (Connection refused)
  02:24:47 [info ] TCP connect try 2. (Success)
  02:24:47 [info ] --- send_tcp_thread start ---
  02:27:27 [info ] VideoCapture opened.
  02:27:27 [info ] --- decode_to_tcp_thread start ---
  02:27:27 [debug] CH(0) TCP send enq(0)
  02:27:27 [debug] CH(0) TCP send enq(1)
  02:27:27 [debug] CH(0) TCP send enq(2)
  02:27:27 [debug] CH(0) TCP send enq(3)
  02:27:27 [debug] CH(0) TCP send enq(4)
  02:27:28 [debug] CH(0) TCP send enq(5)
  ```


## For debug

### Using the cpu_decode alone to execute FPGA DMA

The cpu_decode performs the initial FPGA settings and executes both DMA transmission and reception.
The subsequent FPGA application will become unnecessary.

[How to use the cpu_decode alone to execute FPGA DMA](docs/single_container.md)



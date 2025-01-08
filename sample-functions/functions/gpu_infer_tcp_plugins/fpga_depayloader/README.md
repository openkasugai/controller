###### Copyright 2023 NTT Corporation, FUJITSU LIMITED

# Container image creation method for inference application

1. Get the source of the inference app from the repository
    ```
    $ cd ~/ 
    $ git clone https://github.com/openkasugai/controller.git
    
    // checkout the appropriate branch or tag
    ```

2. Changing Shell Script Permissions
    ```
    $ cd ~/controller/sample-functions/functions/gpu_infer_tcp_plugins/fpga_depayloader
    $ chmod a+x build_app.sh
    $ cd build_docker/gpu-deepstream 
    $ chmod a+x generate_engine_file.sh 
    $ chmod a+x find_gpu.sh 
    $ chmod a+x check_gpus.sh
    ```

3. Run the buildah command (configure the proxy as needed)
    ```
    $ sudo buildah bud --runtime=/usr/bin/nvidia-container-runtime -t gpu_infer_tcp:1.0.0 -f Dockerfile ../../../../../
    ```
Container image is not available on ghcr according to GStreamer licensing terms.

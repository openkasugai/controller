######################################
# INSTALL OPENCV ON UBUNTU OR DEBIAN #
######################################

# -------------------------------------------------------------------- |
#                       SCRIPT OPTIONS                                 |
# ---------------------------------------------------------------------|
OPENCV_VERSION='3.4.3'                        # Version to be installed
OPENCV_INSTALL_PATH='/usr/local/opencv-3.4.3' # PATH to installed
# -------------------------------------------------------------------- |

## 1. KEEP UBUNTU OR DEBIAN UP TO DATE

#apt-get -y update


## 2. INSTALL THE DEPENDENCIES

# Gstreamer for Video I/O:
apt-get install -y libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev \
                   gstreamer1.0-plugins-base gstreamer1.0-plugins-bad gstreamer1.0-plugins-good


## 3. INSTALL THE LIBRARY

wget https://github.com/opencv/opencv/archive/${OPENCV_VERSION}.zip
unzip ${OPENCV_VERSION}.zip && rm ${OPENCV_VERSION}.zip
mv opencv-${OPENCV_VERSION} OpenCV

cd OpenCV && mkdir build && cd build

cmake -DBUILD_IPP_IW=ON \
      -DBUILD_ITT=ON \
      -DBUILD_JAVA=OFF \
      -DBUILD_PACKAGE=OFF \
      -DBUILD_PERF_TESTS=OFF \
      -DBUILD_PROTOBUF=OFF \
      -DBUILD_TESTS=OFF \
      -DBUILD_opencv_apps=OFF \
      -DBUILD_opencv_calib3d=OFF \
      -DBUILD_opencv_dnn=OFF \
      -DBUILD_opencv_features2d=OFF \
      -DBUILD_opencv_flann=OFF \
      -DBUILD_opencv_highgui=OFF \
      -DBUILD_opencv_imgcodecs=ON \
      -DBUILD_opencv_imgproc=ON \
      -DBUILD_opencv_java=OFF \
      -DBUILD_opencv_java_bindings_generator=OFF \
      -DBUILD_opencv_js=OFF \
      -DBUILD_opencv_ml=OFF \
      -DBUILD_opencv_objdetect=OFF \
      -DBUILD_opencv_photo=OFF \
      -DBUILD_opencv_python3=OFF \
      -DBUILD_opencv_python_bindings_generator=OFF \
      -DBUILD_opencv_shape=OFF \
      -DBUILD_opencv_stitching=OFF \
      -DBUILD_opencv_superres=OFF \
      -DBUILD_opencv_ts=OFF \
      -DBUILD_opencv_video=OFF \
      -DBUILD_opencv_videoio=ON \
      -DBUILD_opencv_videostab=OFF \
      -DCMAKE_BUILD_TYPE=Release \
      -DCMAKE_INSTALL_PREFIX=${OPENCV_INSTALL_PATH} \
      -DCV_TRACE=OFF \
      -DENABLE_PRECOMPILED_HEADERS=OFF \
      -DWITH_1394=OFF \
      -DWITH_ARAVIS=OFF \
      -DWITH_CUBLAS=OFF \
      -DWITH_CUFFT=OFF \
      -DWITH_EIGEN=OFF \
      -DWITH_FFMPEG=OFF \
      -DWITH_GSTREAMER=ON \
      -DWITH_GTK=OFF \
      -DWITH_IMGCODEC_HDR=OFF \
      -DWITH_IMGCODEC_PXM=OFF \
      -DWITH_IMGCODEC_SUNRASTER=OFF \
      -DWITH_JASPER=OFF \
      -DWITH_JPEG=OFF \
      -DWITH_LAPACK=OFF \
      -DWITH_MATLAB=OFF \
      -DWITH_NVCUVID=OFF \
      -DWITH_OPENCL=OFF \
      -DWITH_OPENCLAMDBLAS=OFF \
      -DWITH_OPENCLAMDFFT=OFF \
      -DWITH_OPENEXR=OFF \
      -DWITH_PNG=OFF \
      -DWITH_PROTOBUF=OFF \
      -DWITH_PVAPI=OFF \
      -DWITH_TBB=ON \
      -DWITH_TIFF=OFF \
      -DWITH_V4L=OFF \
      -DWITH_VTK=OFF \
      -DWITH_WEBP=OFF \
      ..

make -j8
make install
echo "${OPENCV_INSTALL_PATH}/lib" > /etc/ld.so.conf.d/opencv.conf
ldconfig


## 4. REMOVE THE TEMP DATA

cd ../../
rm -rf ./OpenCV

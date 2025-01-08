#!/bin/bash
# 1. Download GStreamer source file(gstgst-plugins-base).
git clone https://github.com/GStreamer/gst-plugins-base.git gst-plugins-base

# 2. Checkout target version for patching
git -C gst-plugins-base checkout refs/tags/1.19.2

# 3. Move source file.
mv gst-plugins-base/gst-libs/gst/gettext.h . && mv gst-plugins-base/gst-libs/gst/gst-i18n-plugin.h . && \
mv gst-plugins-base/gst/tcp/gsttcpserversrc.c . && mv gst-plugins-base/gst/tcp/gsttcpserversrc.h . 

# 4. Rename source file.
mv gsttcpserversrc.c gstfpgadepayloader.c & mv gsttcpserversrc.h gstfpgadepayloader.h  

# 5. Apply patch file to tcpserversrc.
patch -p4 < gpu_inference_tcp_rcv.patch

# 6. Build GPU inference app
make clean && make

# 7. Delete source file.
rm gettext.h gst-i18n-plugin.h gstfpgadepayloader.*

# 8. Delete GStreamer source file.
rm -rf gst-plugins-base

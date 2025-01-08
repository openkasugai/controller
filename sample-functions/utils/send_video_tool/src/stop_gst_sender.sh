#!/bin/bash

# Copyright 2022 NTT Corporation, FUJITSU LIMITED 


pgrep -f gst_child_process.sh | xargs kill >/dev/null 2>&1
pgrep gst-launch-1.0 | xargs kill >/dev/null 2>&1

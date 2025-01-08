# Copyright 2024 NTT Corporation, FUJITSU LIMITED

import time

start_time = time.time()


class PERF:
    def __init__(self):
        global start_time
        self.start_time = start_time
        self.is_first = True
        self.frame_count = 0

    def update_fps(self):
        end_time = time.time()
        if self.is_first:
            self.start_time = end_time
            self.is_first = False
        else:
            self.frame_count = self.frame_count + 1

    def get_fps(self):
        end_time = time.time()
        stream_fps = float(self.frame_count/(end_time - self.start_time))
        self.frame_count = 0
        self.start_time = end_time
        return '{:.3f}'.format(round(stream_fps, 3))

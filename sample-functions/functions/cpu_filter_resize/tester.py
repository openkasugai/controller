# -*- coding:utf-8 -*-

# Copyright 2024 NTT Corporation, FUJITSU LIMITED

import logging
import os
import select
import socket
import struct
import threading
import time
import cv2 as cv
import fire
import numpy as np
import schedule
from fps import PERF

logger = logging.getLogger(__name__)
ch = logging.StreamHandler()
fmt = logging.Formatter(
    '%(asctime)s %(levelname)-8s %(message)s [%(filename)s:%(lineno)d]')
ch.setFormatter(fmt)
logger.addHandler(ch)

stop_event = threading.Event()

BUFF_SIZE = 1024*1024
FHEADER_SIZE = 48
FHEADER_MARKER = b'\xad\x10\xff\xe0'
PORT = 9999


def run_server(width, height, fps):

    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # Set Receive Buffer Size
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_RCVBUF, 1024*1024)
    rcvbufsize = sock.getsockopt(socket.SOL_SOCKET, socket.SO_RCVBUF)
    logger.info(f'Recv buffer size: {rcvbufsize}')

    sock.bind(('0.0.0.0', PORT))
    sock.listen()
    logger.info("Listening on port: "+str(PORT))

    while not stop_event.is_set():
        client, addr = sock.accept()
        logger.info(f"New client from {addr}")

        plen = 0
        recv_data = bytearray()

        codec = cv.VideoWriter_fourcc(*'mp4v')
        path = os.getcwd()
        out_file = f'{path}/out_{width}x{height}_{fps}fps.mp4'
        writer = cv.VideoWriter(out_file, codec, fps, (width, height))

        while not stop_event.is_set():
            try:
                # Raise recv timeout exception in select every second, giving chance to stop_event check
                readable, _, _ = select.select([client], [], [], 1)
                fd = readable[0]

                if len(recv_data) < FHEADER_SIZE:
                    # When frame header is not received. Support for receiving split frame headers
                    buf = buf = fd.recv(FHEADER_SIZE - len(recv_data))
                else:
                    # When the frame header has been received and image data is to be received
                    # Calculate the length of the remaining unreceived image data in the frame data
                    remaining = plen + FHEADER_SIZE - len(recv_data)
                    if remaining > BUFF_SIZE:
                        buf = fd.recv(BUFF_SIZE)
                    else:
                        buf = fd.recv(remaining)

                if not buf:
                    logger.info(f'buf is none')
                    break

                # Add incoming data (buf) to recv_data
                recv_data.extend(buf)

                if len(recv_data) == FHEADER_SIZE:
                    # Get Frame Header Value
                    fh = struct.unpack('<4s I 4s I 32s',
                                       bytes(recv_data))
                    if FHEADER_MARKER != fh[0]:
                        logger.error(
                            f'frame header is incorrect. first 4 bytes: {fh[0]}')
                    plen = fh[1]
                    seqnum = fh[3]
                    logger.debug(
                        f'sequence_number: {seqnum} , payload_len: {plen}')
                    continue
                elif len(recv_data) < (plen+FHEADER_SIZE):
                    # When all frame data has not been received yet
                    continue

                global perf
                perf.update_fps()

                raw_data = bytes(recv_data[FHEADER_SIZE:])

                # byes -> image data(ndarray)
                arr = np.frombuffer(raw_data, dtype=np.uint8)
                img = arr.reshape(height, width, 3)
                cv.imwrite(f'{path}/out.jpg', img)
                writer.write(img)

                recv_data = bytearray()

            except:
                continue

        client.close()
        logger.info('connection closed')
        writer.release()
        logger.info('writer released')
        stop_event.set()

    sock.close()
    logger.info('server socket closed')


def suchedule_event():
    schedule.every(3).seconds.do(print_fps)
    while not stop_event.is_set():
        schedule.run_pending()
        time.sleep(1)


def print_fps():
    global perf
    logger.info(f'FPS: {perf.get_fps()}')


def main(width=1280, height=1280, fps=15, loglevel='INFO'):
    logger.setLevel(logging.getLevelName(loglevel))

    global perf
    perf = PERF()

    t1 = threading.Thread(target=suchedule_event)
    t1.start()

    t2 = threading.Thread(target=run_server, args=[width, height, fps])
    t2.start()

    try:
        while not stop_event.is_set():
            time.sleep(1)
    except KeyboardInterrupt:
        logger.info('Caught keyboardinterrupt')
        stop_event.set()
        t1.join()
        t2.join()

    logger.info('finish')


if __name__ == '__main__':
    fire.Fire(main)

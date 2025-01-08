# -*- coding:utf-8 -*-

# Copyright 2024 NTT Corporation, FUJITSU LIMITED

import asyncio
from dataclasses import dataclass
import logging
import queue
import threading
import os
import socket
import struct
import time
import numpy as np
import fire
usegpu = os.getenv('USE_GPU', '0')
if usegpu == '1':
    import cvcuda
    import cupy as cp
else:
    import cv2 as cv

logger = logging.getLogger(__name__)
ch = logging.StreamHandler()
fmt = logging.Formatter(
    '%(asctime)s %(levelname)-8s %(message)s [%(filename)s:%(lineno)d]')
ch.setFormatter(fmt)
logger.addHandler(ch)

fr_func = None
stop_event = threading.Event()


BUFF_SIZE = 1024*1024
FHEADER_SIZE = 48
FHEADER_MARKER = b'\xad\x10\xff\xe0'


@dataclass
class FrameData:
    fheader: tuple
    img: np.ndarray


def filter_resize_cvcuda(queue1, queue2, dsize_width, dsize_height):

    logger.info('start filter-resize(cv-cuda) thread')
    while not stop_event.is_set():
        fd: FrameData = queue1.get()
        img = cvcuda.as_image(cp.asarray(fd.img))
        img_tensor = cvcuda.as_tensor(img)

        filtered_tensor = cvcuda.median_blur(img_tensor, (5, 5))
        resized_tensor = cvcuda.resize(
            filtered_tensor, (1, dsize_width, dsize_height, 3), cvcuda.Interp.CUBIC)

        resized_img = cp.asarray(resized_tensor.cuda())
        fd.img = cp.asnumpy(resized_img[0])

        # Passing the frame header and Image to the processing result sending thread
        queue2.put(fd)


def filter_resize(queue1, queue2, dsize_width, dsize_height):
    logger.info('start filter-resize thread')
    while not stop_event.is_set():
        fd: FrameData = queue1.get()
        dst = cv.medianBlur(fd.img, ksize=5)
        fd.img = cv.resize(dst, (dsize_width, dsize_height))

        # Passing the frame header and Image to the processing result sending thread
        queue2.put(fd)


def make_connection(out_addr, out_port, sockbufsize):
    while not stop_event.is_set():
        try:
            sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            sock.setsockopt(socket.SOL_SOCKET, socket.SO_SNDBUF, sockbufsize)
            sendbufsize = sock.getsockopt(
                socket.SOL_SOCKET, socket.SO_SNDBUF)
            sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            sock.connect((out_addr, out_port))
            logger.info(f'Connected host: {out_addr}, port: {out_port}')
            logger.info(f'Send buffer size: {sendbufsize}')
            return sock
        except socket.error as e:
            logger.warning(f'failed to connect, try reconnect. {e}')
            time.sleep(1)


def run_sender(queue2, out_addr, out_port, out_width, out_height, sockbufsize):
    logger.info('start sender thread')
    sock = make_connection(out_addr, out_port, sockbufsize)

    while not stop_event.is_set():
        try:
            fd: FrameData = queue2.get(timeout=0.5)
        except queue.Empty:
            continue

        fh = fd.fheader

        # Update frame header values and combine with filtered resized data
        updated_fh = struct.pack(
            "<4s I 4s I 32s", fh[0],
            out_width*out_height*3,
            fh[2], fh[3], fh[4])

        send_data = bytearray(updated_fh)
        send_data.extend(fd.img.tobytes())

        try:
            sock.send(send_data)
            logger.debug(f"send_data len={len(send_data)}")
        except socket.error as e:
            logger.warning(f'connection lost, try reconnect. {e}')
            sock = make_connection(out_addr, out_port, sockbufsize)
    if sock is not None:
        sock.close()
        logger.info('Close sender connection')


async def handle_client(loop, queue, client, in_width, in_height):
    plen = 0
    recv_data = bytearray()

    while True:
        if len(recv_data) < FHEADER_SIZE:
            # When frame header is not received. Support for receiving split frame headers
            buf = await loop.sock_recv(client, FHEADER_SIZE - len(recv_data))
        else:
            # When the frame header has been received and image data is to be received
            # Calculate the length of the remaining unreceived image data in the frame data
            remaining = plen + FHEADER_SIZE - len(recv_data)
            if remaining > BUFF_SIZE:
                buf = await loop.sock_recv(client, BUFF_SIZE)
            else:
                buf = await loop.sock_recv(client, remaining)

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

        raw_data = bytes(recv_data[FHEADER_SIZE:])

        # byes -> image data(ndarray)
        arr = np.frombuffer(raw_data, dtype=np.uint8)
        img = arr.reshape(in_height, in_width, 3)

        # filter-resize Pass frame header and Image to thread
        fd = FrameData(fh, img)
        queue.put(fd)

        recv_data = bytearray()

    logger.info('Close the connection')
    client.close()


async def run_receiver(queue, in_port, in_width, in_height, sockbufsize):
    logger.info('start receiver')

    server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # Set socket receive buffer size
    server.setsockopt(socket.SOL_SOCKET, socket.SO_RCVBUF, sockbufsize)
    rcvbufsize = server.getsockopt(socket.SOL_SOCKET, socket.SO_RCVBUF)
    logger.info(f'Recv buffer size: {rcvbufsize}')

    server.bind(('0.0.0.0', in_port))
    server.listen()

    logger.info(f'Listening on port: {in_port}')
    server.setblocking(False)
    loop = asyncio.get_event_loop()

    while True:
        client, addr = await loop.sock_accept(server)
        client.setblocking(False)
        logger.info(f"New client from {addr}")
        asyncio.create_task(handle_client(
            loop, queue, client, in_width, in_height))


async def start_threads(in_port, out_addr, out_port,
                        in_width, in_height, out_width, out_height, sockbufsize):
    queue1: 'queue.Queue[FrameData]' = queue.Queue(maxsize=200)
    queue2: 'queue.Queue[FrameData]' = queue.Queue(maxsize=200)

    # fr_func already contains filter resize functions for CPU/GPU modes in main()
    t1 = threading.Thread(
        target=fr_func,
        args=[queue1, queue2, out_width, out_height],
        daemon=True)
    t1.start()

    t2 = threading.Thread(
        target=run_sender,
        args=[queue2, out_addr, out_port, out_width, out_height, sockbufsize],
        daemon=True)
    t2.start()

    task = asyncio.create_task(run_receiver(
        queue1, in_port, in_width, in_height, sockbufsize))
    await task


def main(in_port=8888, out_addr='127.0.0.1', out_port=9999,
         in_width=3840, in_height=2160, out_width=1280, out_height=1280,
         sockbufsize=1024*1024*4, loglevel='INFO'):
    logger.setLevel(logging.getLevelName(loglevel))

    global fr_func
    if usegpu == '1':
        logger.info('GPU(CV-CUDA) mode')
        fr_func = filter_resize_cvcuda
    else:
        logger.info('CPU mode')
        fr_func = filter_resize

    logger.info(f'in_port={in_port}, out_addr={out_addr}, out_port={out_port}')
    logger.info(
        f'in_width={in_width}, in_height={in_height}, out_width={out_width}, out_height={out_height}')

    try:
        asyncio.run(start_threads(in_port, out_addr, out_port,
                                  in_width, in_height, out_width, out_height,
                                  sockbufsize))
    except:
        logger.info(f'Caught exception')
        stop_event.set()
        time.sleep(3)


if __name__ == '__main__':
    fire.Fire(main)

/*************************************************
* Copyright 2025 NTT Corporation , FUJITSU LIMITED
*************************************************/

#include "connection.hpp"

#include <liblogging.h>
#include <libfpgactl.h>
#include <liblldma.h>
#include <libptu.h>
#include <libchain.h>

#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <unistd.h>
#include <errno.h>
#include <fcntl.h>
#include <sys/file.h>
#include <sys/un.h>
#include <sys/types.h>
#include <sys/mman.h>
#include <sys/select.h>
#include <arpa/inet.h>


connection::connection(
  logger *plog)
{
  libfpga_log_set_level(LIBFPGA_LOG_NOTHING);
  if (plog) {
    plog_ = plog;
  } else {
    plog_self_ = new logger(default_filename_);
    if (!plog_self_)
      fprintf(stdout, " ! Failed to new logger(%s)\n", default_filename_);
    plog_ = plog_self_;
  }
}

connection::~connection(
  void)
{
  if (plog_self_)
    delete plog_self_;
};

int connection::fpga_initialize(
  void)
{
  plog_->print(" * Start : %s\n", __func__);
  int ret = fpga_scan_devices();
  if (ret < 0) {
    plog_->print(" ! Failed to fpga_scan_devices : ret(%d)\n", ret);
    return ret;
  }
  if (ret == 0) {
    plog_->print(" ! FPGAs not found...\n");
    plog_->print(" ! Please check if driver is loaded\n");
    return -1;
  }
  ret = fpga_enable_regrw_all();
  if (ret) {
    plog_->print(" ! Failed to fpga_enable_regrw_all() in %d\n", ret);
    return ret;
  }
  plog_->print(" * Succeed to fpga_enable_regrw_all()\n");
  plog_->print(" * End   : %s\n", __func__);
  return 0;
}

int connection::fpga_ptu_initialize(
  void)
{
  plog_->print(" * Start : %s\n", __func__);
  int ret;
  for (uint32_t dev_id = 0; dev_id < (uint32_t)fpga_get_num(); dev_id++) {
    for (uint32_t lane = 0; lane < fpga_lane_max_; lane++) {
      in_addr_t fpga_addr    = fpga_addr_    + 0x100 * (dev_id + lane);
      in_addr_t fpga_gateway = fpga_gateway_ + 0x100 * (dev_id + lane);
      uint8_t fpga_mac[6];
      memcpy(fpga_mac, fpga_mac_, sizeof(fpga_mac));
      fpga_mac[5] = fpga_mac_[5] + dev_id + lane;
      ret = fpga_ptu_init(dev_id, lane, fpga_addr, fpga_subnet_, fpga_gateway, fpga_mac);
      if (ret) {
        plog_->print(" ! Failed to fpga_ptu_init(%u,%u,%08x,%08x,%08x,[%02x,%02x,%02x,%02x,%02x,%02x])\n",
          dev_id, lane, fpga_addr, fpga_subnet_, fpga_gateway,
          fpga_mac[0], fpga_mac[1], fpga_mac[2], fpga_mac[3], fpga_mac[4], fpga_mac[5]);
        fpga_finalize();
        return -1;
      } else {
        plog_->print(" * Succeed to fpga_ptu_init(%u,%u,%08x,%08x,%08x,[%02x,%02x,%02x,%02x,%02x,%02x])\n",
          dev_id, lane, fpga_addr, fpga_subnet_, fpga_gateway,
          fpga_mac[0], fpga_mac[1], fpga_mac[2], fpga_mac[3], fpga_mac[4], fpga_mac[5]);
      }
    }
  }
  plog_->print(" * End   : %s\n", __func__);
  return 0;
}

int connection::fpga_finalize(
  void)
{
  plog_->print(" * Start : %s\n", __func__);
  while(fpga_defmap_.size()) {
    server_delete_connection(fpga_defmap_[0]);
  }
  for (uint32_t dev_id = 0; dev_id < (uint32_t)fpga_get_num(); dev_id++) {
    for (uint32_t lane = 0; lane < fpga_lane_max_; lane++) {
      if (fpga_ptu_exit(dev_id, lane))
        plog_->print(" ! Failed to fpga_ptu_exit(%u,%u)\n", dev_id, lane);
      else
        plog_->print(" * Succeed to fpga_ptu_exit(%u,%u)\n", dev_id, lane);
    }
  }
  fpga_finish();
  plog_->print(" * End   : %s\n", __func__);
  return 0;
}

int connection::common_listen(
  const char *laddr,
  uint32_t lport)
{
  plog_->print(" * Start : %s\n", __func__);

  // socket
  int fd_listen = socket(AF_INET, SOCK_STREAM, 0);
  if (fd_listen < 0) {
    int err = errno;
    plog_->print(" ! Failed to socket : %s\n", strerror(err));
    return fd_listen;
  }
  plog_->print(" * socket\n");

  // bind
  struct sockaddr_in listen_addr;
  listen_addr.sin_family = AF_INET;
  listen_addr.sin_port = htons(lport);
  listen_addr.sin_addr.s_addr = inet_addr(laddr);
  int ret = bind(fd_listen, (struct sockaddr*)&listen_addr, sizeof(listen_addr));
  if (ret < 0) {
    int err = errno;
    plog_->print(" ! bind : %s(%08x),%u : %s\n",
      laddr, listen_addr.sin_addr.s_addr, lport, strerror(err));
    close(fd_listen);
    return ret;
  }
  plog_->print(" * bind : %s(%08x),%u\n",
    laddr, listen_addr.sin_addr.s_addr, lport);

  // listen
  ret = listen(fd_listen, 16);
  if (ret < 0) {
    int err = errno;
    plog_->print(" ! listen : %s\n", strerror(err));
    close(fd_listen);
    return ret;
  }
  plog_->print(" * listen\n");

  plog_->print(" * End   : %s\n", __func__);
  return fd_listen;
}

int connection::common_accept(
  int fd_listen)
{
  plog_->print(" * Start : %s\n", __func__);

  if (fd_listen < 0) {
    plog_->print(" ! Invalid fd(%d)...\n", fd_listen);
    return -1;
  }

  // accept
  struct sockaddr_in accept_addr;
  socklen_t len = sizeof(accept_addr);
  int fd_accept = accept(fd_listen, (struct sockaddr*)&accept_addr, &len);
  if (fd_accept < 0) {
    int err = errno;
    plog_->print(" ! accept : %d : %s\n", fd_listen, strerror(err));
    return fd_accept;
  }
  plog_->print(" * accept : %d->%d\n", fd_listen, fd_accept);

  plog_->print(" * End   : %s\n", __func__);
  return fd_accept;
}


int connection::common_connect(
  const char *raddr,
  uint32_t rport,
  const char *laddr,
  uint32_t lport)
{
  plog_->print(" * Start : %s\n", __func__);

  // socket
  int fd_connect = socket(AF_INET, SOCK_STREAM, 0);
  if (fd_connect < 0) {
    int err = errno;
    plog_->print(" ! Failed to socket : %s\n", strerror(err));
    return fd_connect;
  }
  plog_->print(" * socket\n");

  // bind
  if (laddr) {
    struct sockaddr_in bind_addr;
    bind_addr.sin_family = AF_INET;
    bind_addr.sin_port = htons(lport);
    bind_addr.sin_addr.s_addr = inet_addr(laddr);
    int ret = bind(fd_connect, (struct sockaddr*)&bind_addr, sizeof(bind_addr));
    if (ret < 0) {
      int err = errno;
      plog_->print(" ! bind : %s(%08x),%u : %s\n",
        laddr, bind_addr.sin_addr.s_addr, lport, strerror(err));
      close(fd_connect);
      return ret;
    }
    plog_->print(" * bind : %s(%08x),%u\n", laddr, bind_addr.sin_addr.s_addr, lport);
  }

  // connect
  struct sockaddr_in connect_addr;
  connect_addr.sin_family = AF_INET;
  connect_addr.sin_port = htons(rport);
  connect_addr.sin_addr.s_addr = inet_addr(raddr);
  int ret = connect(fd_connect, (struct sockaddr*)&connect_addr, sizeof(connect_addr));
  if (ret < 0) {
    int err = errno;
    plog_->print(" ! connect : %s(%08x),%u : %s\n",
      raddr, connect_addr.sin_addr.s_addr, rport, strerror(err));
    close(fd_connect);
    return ret;
  }
  plog_->print(" * connect : %s(%08x),%u\n",
    raddr, connect_addr.sin_addr.s_addr, rport);

  plog_->print(" * End   : %s\n", __func__);
  return fd_connect;
}


int connection::common_send(
  int fd,
  const void *data,
  socklen_t len)
{
  plog_->print(" * Start : %s\n", __func__);

  // send
  int send_flag = 0;
  socklen_t send_len_remain = len;
  socklen_t send_len;
  while (send_len_remain) {
    send_len = send(
      fd,
      (char*)data + (len - send_len_remain),
      send_len_remain,
      send_flag);
    if (send_len <= 0) {
      int err = errno;
      plog_->print(" ! send : (%dbyte/%dbyte)[%d] : %s\n",
        len - send_len_remain, len, strerror(err));
      return -1;
    }
    send_len_remain -= send_len;
  }
  plog_->print(" * send : (%dbyte)[%d]\n", len, fd);

  plog_->print(" * End   : %s\n", __func__);
  return 0;
}


int connection::common_recv(
  int fd,
  void *data,
  socklen_t len)
{
  plog_->print(" * Start : %s\n", __func__);

  // recv
  int recv_flag = 0;
  socklen_t recv_len_remain = len;
  socklen_t recv_len;
  while (recv_len_remain) {
    recv_len = recv(
      fd,
      (char*)data + (len - recv_len_remain),
      recv_len_remain,
      recv_flag);
    if (recv_len == 0) {
      plog_->print(" * recv : (connection lost)\n");
      return 1;
    } else if (recv_len < 0) {
      int err = errno;
      plog_->print(" ! recv : (%dbyte/%dbyte)[%d] : %s\n",
        len - recv_len_remain, len, strerror(err));
      return -1;
    }
    recv_len_remain -= recv_len;
  }
  plog_->print(" * recv : (%dbyte)[%d]\n", len, fd);

  plog_->print(" * End   : %s\n", __func__);
  return 0;
}


int connection::common_close(
  int &fd)
{
  plog_->print(" * Start : %s\n", __func__);

  // close
  if (fd < 0) {
    plog_->print(" ! Invalid fd(%d)...\n", fd);
    return -1;
  }
  close(fd);
  fd = -1;

  plog_->print(" * End   : %s\n", __func__);
  return 0;
}


int connection::host_connect_fpga(
  const char *raddr,
  uint32_t rport,
  const char *laddr,
  uint32_t lport)
{
  plog_->print(" * Start : %s\n", __func__);

  int fd_connect = common_connect(raddr, rport, laddr, lport);
  if (fd_connect < 0) {
    plog_->print(" ! Failed to connect to FPGA\n", __func__);
    return fd_connect;
  }

  plog_->print(" * End   : %s\n", __func__);
  return fd_connect;
}


int connection::host_accept_fpga(
  const char *laddr,
  uint32_t lport)
{
  plog_->print(" * Start : %s\n", __func__);

  int fd_listen = common_listen(laddr, lport);
  if (fd_listen < 0) {
    plog_->print(" ! Failed to listen\n", __func__);
    return fd_listen;
  }

  int fd_accept = common_accept(fd_listen);
  common_close(fd_listen);
  if (fd_accept < 0) {
    plog_->print(" ! Failed to listen\n", __func__);
    return fd_accept;
  }

  plog_->print(" * End   : %s\n", __func__);
  return fd_accept;
}


int connection::host_close_fpga(
  int fd)
{
  plog_->print(" * Start : %s\n", __func__);

  common_close(fd);

  plog_->print(" * End   : %s\n", __func__);
  return 0;
}


int connection::host_exec_fcc(
  const char *options)
{
  if (!options) {
    plog_->print(" ! Invalid argument : NULL\n");
    return -1;
  }
  int ret = system(NULL);
  if (!ret) {
    plog_->print(" ! Fatal error : system(NULL) : %d\n", ret);
    return ret;
  }
  char cmd[1024];
  sprintf(cmd, "%s %s", fcc_path_, options);
  ret = system(cmd) >> 8;
  plog_->print(" # Result(%d) : Command(%s)\n", ret, cmd);
  return ret;
}


int connection::server_listen(
  void)
{
  plog_->print(" * Start : %s\n", __func__);

  if (server_fd_listen_ != -1) {
    plog_->print(" ! Already listening\n");
    return -1;
  }

  int fd_listen = common_listen(server_listen_addr_, server_listen_port_);
  while (fd_listen < 0) {
    plog_->print(" ! Failed to listen as server\n");
    sleep(5);
    fd_listen = common_listen(server_listen_addr_, server_listen_port_);
  }

  server_fd_listen_ = fd_listen;

  plog_->print(" * End   : %s\n", __func__);
  return 0;
}

int connection::server_accept(
  void)
{
  plog_->print(" * Start : %s\n", __func__);

  if (server_fd_listen_ == -1) {
    plog_->print(" ! Not listening as server yet\n");
    return -1;
  }

  int fd_accept = common_accept(server_fd_listen_);
  if (fd_accept < 0) {
    plog_->print(" ! Failed to accept as server\n");
    return -1;
  }

  plog_->print(" * End   : %s\n", __func__);
  return fd_accept;
}

int connection::server_close(
  void)
{
  return server_close(server_fd_listen_);
}


int connection::server_close(
  int fd)
{
  plog_->print(" * Start : %s\n", __func__);

  if (fd < 0) {
    plog_->print(" ! Already closed\n");
    return 0;
  }
  common_close(fd);

  plog_->print(" * End   : %s\n", __func__);
  return 0;
}


int connection::server_create_connection(
  const char *device,
  uint32_t dir,
  uint32_t chid,
  const char *connector_id,
  void **pdmainfo)
{
  if (!device || !connector_id || !pdmainfo) {
    plog_->print(" ! Invalid argument : NULL\n");
    return -1;
  }
  uint32_t dev_id;
  int ret = fpga_get_dev_id(device, &dev_id);
  if (ret < 0) {
    plog_->print(" ! Failed to fpga_get_dev_id(%s) in %d\n", device, ret);
    return ret;
  }
  plog_->print(" * Succeed to fpga_get_dev_id(%s) : dev_id(%u)\n", device, dev_id);
  dma_info_t *dmainfo = new dma_info_t;
  ret = fpga_lldma_init(
    dev_id,
    (dma_dir_t)dir,
    chid,
    connector_id,
    dmainfo);
  if (ret < 0) {
    plog_->print(" ! Failed to fpga_lldma_init(%u,%u,%u,%s) in %d\n",
      dev_id, dir, chid, connector_id, ret);
    delete(dmainfo);
    return ret;
  }
  plog_->print(" * Succeed to fpga_lldma_init(%u,%u,%u,%s,%lx)\n",
    dev_id, dir, chid, connector_id, (uintptr_t)dmainfo);
  *pdmainfo = (void*)dmainfo;
  return 0;
}


int connection::server_create_connection(
  const char *device,
  uint32_t lane,
  in_port_t fpga_port,
  in_addr_t host_addr,
  in_port_t host_port,
  bool is_server,
  uint32_t *cid)
{
  if (!device || !cid) {
    plog_->print(" ! Invalid argument : NULL\n");
    return -1;
  }
  uint32_t dev_id;
  int ret = fpga_get_dev_id(device, &dev_id);
  if (ret < 0) {
    plog_->print(" ! Failed to fpga_get_dev_id(%s) in %d\n", device, ret);
    return ret;
  }
  plog_->print(" * Succeed to fpga_get_dev_id(%s) : dev_id(%u)\n", device, dev_id);
  char cid_filename[256];
  sprintf(cid_filename, "cid-%s-%u-%u.log", device, lane, fpga_port);
  struct timeval ptu_timeout;
  ptu_timeout.tv_sec  = 10;
  ptu_timeout.tv_usec = 0;
  struct timeval *pptu_timeout = &ptu_timeout;
  if (is_server) {
    ret = fpga_ptu_listen(dev_id, lane, fpga_port);
    if (ret) {
      plog_->print(" ! Failed to fpga_ptu_listen(%u,%u,%u) in %d\n",
        dev_id, lane, fpga_port, ret);
      return ret;
    }
    plog_->print(" * Succeed to fpga_ptu_listen(%u,%u,%u)\n",
        dev_id, lane, fpga_port);
    ret = fpga_ptu_accept(
      dev_id,
      lane,
      fpga_port,
      host_addr,
      host_port,
      pptu_timeout,
      cid);
    if (ret < 0) {
      plog_->print(" ! Failed to fpga_ptu_accept(%u,%u,%u,%08x,%u,%lx,%lx) in %d\n",
        dev_id, lane, fpga_port, host_addr, host_port,
        (uintptr_t)pptu_timeout, (uintptr_t)cid, ret);
      int _ret = fpga_ptu_listen_close(dev_id, lane, fpga_port);
      if (_ret < 0)
        plog_->print(" ! Failed to fpga_ptu_listen_close(%u,%u,%u) in %d\n",
          dev_id, lane, fpga_port, _ret);
      else
        plog_->print(" * Succeed to fpga_ptu_listen_close(%u,%u,%u)\n",
          dev_id, lane, fpga_port);
      return ret;
    }
    plog_->print(" * Succeed to fpga_ptu_accept(%u,%u,%u,%08x,%u,%lx,%lx) : cid(%u)\n",
      dev_id, lane, fpga_port, host_addr, host_port,
      (uintptr_t)pptu_timeout, (uintptr_t)cid, *cid);
  } else {
    ret = fpga_ptu_connect(
      dev_id,
      lane,
      fpga_port,
      host_addr,
      host_port,
      pptu_timeout,
      cid);
    if (ret < 0) {
      plog_->print(" ! Failed to fpga_ptu_connect(%u,%u,%u,%08x,%u,%lx,%lx) in %d\n",
        dev_id, lane, fpga_port, host_addr, host_port,
        (uintptr_t)pptu_timeout, (uintptr_t)cid, ret);
      return ret;
    }
    plog_->print(" * Succeed to fpga_ptu_connect(%u,%u,%u,%08x,%u,%lx,%lx) : cid(%u)\n",
      dev_id, lane, fpga_port, host_addr, host_port,
      (uintptr_t)pptu_timeout, (uintptr_t)cid, *cid);
  }
  logger cid_logger(cid_filename);
  cid_logger.set_timestamp(NULL);
  cid_logger.print("%u\n", *cid);
  return 0;
}

int connection::server_create_connection(
  const char *device,
  uint32_t lane,
  uint32_t fchid,
  uint32_t extif_id,
  uint32_t cid,
  uint32_t active_flag,
  uint32_t direct_flag,
  uint32_t virtual_flag,
  uint32_t blocking_flag)
{
  if (!device) {
    plog_->print(" ! Invalid argument : NULL\n");
    return -1;
  }
  uint32_t dev_id;
  int ret = fpga_get_dev_id(device, &dev_id);
  if (ret < 0) {
    plog_->print(" ! Failed to fpga_get_dev_id(%s) in %d\n", device, ret);
    return ret;
  }
  plog_->print(" * Succeed to fpga_get_dev_id(%s) : dev_id(%u)\n", device, dev_id);
  if (direct_flag) {
    // ingress
    ret = fpga_chain_connect_ingress(
      dev_id,
      lane,
      fchid,
      extif_id,
      cid,
      active_flag,
      direct_flag);
    if (ret < 0) {
      plog_->print(" ! Failed to fpga_chain_connect_ingress(%u,%u,%u,%u,%u,%u,%u) in %d\n",
        dev_id, lane, fchid, extif_id, cid, active_flag, direct_flag, ret);
      return ret;
    }
    plog_->print(" * Succeed to fpga_chain_connect_ingress(%u,%u,%u,%u,%u,%u,%u)\n",
      dev_id, lane, fchid, extif_id, cid, active_flag, direct_flag);
  } else {
    ret = fpga_chain_connect_egress(
      dev_id,
      lane,
      fchid,
      extif_id,
      cid,
      active_flag,
      virtual_flag,
      blocking_flag);
    if (ret < 0) {
      plog_->print(" ! Failed to fpga_chain_connect_egress(%u,%u,%u,%u,%u,%u,%u) in %d\n",
        dev_id, lane, fchid, extif_id, cid, active_flag, direct_flag, ret);
      return ret;
    }
    plog_->print(" * Succeed to fpga_chain_connect_egress(%u,%u,%u,%u,%u,%u,%u)\n",
      dev_id, lane, fchid, extif_id, cid, active_flag, direct_flag);
  }
  return 0;
}


int connection::server_create_connection(
  datadef *request)
{
  plog_->print(" * Start : %s\n", __func__);
  int ret = -1;
  if (request->lldma.device[0] != '\0') {
    plog_->print(" # Setup for LLDMA\n");
    ret = server_create_connection(
      request->lldma.device,
      request->lldma.dir,
      request->lldma.chid,
      request->lldma.connector_id,
      &request->lldma.dmainfo);
  } else if (request->ptu.device[0] != '\0') {
    plog_->print(" # Setup for PTU\n");
    ret = server_create_connection(
      request->ptu.device,
      request->ptu.lane,
      request->ptu.fpga_port,
      request->ptu.host_addr,
      request->ptu.host_port,
      request->ptu.is_server,
      &request->ptu.cid);
  } else if (request->chain.device[0] != '\0') {
    plog_->print(" # Setup for CHAIN\n");
    ret = server_create_connection(
      request->chain.device,
      request->chain.lane,
      request->chain.fchid,
      request->chain.extif_id,
      request->chain.cid,
      request->chain.active_flag,
      request->chain.direct_flag,
      request->chain.virtual_flag,
      request->chain.blocking_flag);
  }
  if (ret < 0) {
    plog_->print(" ! Failed to create connection in %d\n", ret);
    return ret;
  }
  plog_->print(" * Succeed to create connection\n");
  fpga_defmap_.push_back(request);
  plog_->print(" * End   : %s\n", __func__);
  return 0;
}


int connection::server_delete_connection(
  const char *connector_id)
{
  if (!connector_id) {
    plog_->print(" ! Invalid argument : NULL\n");
    return -1;
  }
  for (auto it = fpga_defmap_.begin(); it != fpga_defmap_.end(); it++) {
    if ((*it)->lldma.device[0] == '\0') continue;
    if (strcmp((*it)->lldma.connector_id, connector_id) != 0) continue;
    int ret = fpga_lldma_finish((dma_info_t*)(*it)->lldma.dmainfo);
    if (ret < 0)
      plog_->print(" ! Failed to fpga_lldma_finish(%lx) in %d\n",
        (uintptr_t)(*it)->lldma.dmainfo, ret);
    else
      plog_->print(" * Succeed to fpga_lldma_finish(%lx)\n",
        (uintptr_t)(*it)->lldma.dmainfo);
    delete((dma_info_t*)((*it)->lldma.dmainfo));
    plog_->print(" * Succeed to delete dmainfo\n");
    delete(*it);
    plog_->print(" * Delete connection data\n");
    fpga_defmap_.erase(it);
    plog_->print(" * Delete connection from map\n");
    return ret;
  }
  plog_->print(" ! Failed to %s : Not found...\n", __func__);
  return -1;  // Not found
}


int connection::server_delete_connection(
  const char *device,
  uint32_t lane,
  in_port_t fpga_port)
{
  if (!device) {
    plog_->print(" ! Invalid argument : NULL\n");
    return -1;
  }
  uint32_t dev_id;
  int ret = fpga_get_dev_id(device, &dev_id);
  if (ret < 0) {
    plog_->print(" ! Failed to fpga_get_dev_id(%s) in %d\n", device, ret);
    return ret;
  }
  plog_->print(" * Succeed to fpga_get_dev_id(%s) : dev_id(%u)\n", device, dev_id);
  for (auto it = fpga_defmap_.begin(); it != fpga_defmap_.end(); it++) {
    char *dev = (*it)->ptu.device;
    if (strcmp(dev, device) != 0) {
      // '/dev/xpcie_'=11word
      if (dev[0] == '/' && strlen(dev) > 11 && device[0] != '/') {
        if (strcmp(&dev[11], device) != 0)
          continue;
      } else if (dev[0] != '/' && strlen(device) > 11 && device[0] == '/') {
        if (strcmp(dev, &device[11]) != 0)
          continue;
      } else {
        continue;
      }
    }
    if ((*it)->ptu.lane != lane) continue;
    if ((*it)->ptu.fpga_port != fpga_port) continue;
    ret = fpga_ptu_disconnect(dev_id, lane, (*it)->ptu.cid);
    if (ret < 0)
      plog_->print(" ! Failed to fpga_ptu_disconnect(%u,%u,%u) in %d\n",
        dev_id, lane, (*it)->ptu.cid, ret);
    else
      plog_->print(" * Succeed to fpga_ptu_disconnect(%u,%u,%u)\n",
        dev_id, lane, (*it)->ptu.cid);
    if ((*it)->ptu.is_server) {
      ret = fpga_ptu_listen_close(dev_id, lane, fpga_port);
      if (ret < 0)
        plog_->print(" ! Failed to fpga_ptu_listen_close(%u,%u,%u) in %d\n",
          dev_id, lane, fpga_port, ret);
      else
        plog_->print(" * Succeed to fpga_ptu_listen_close(%u,%u,%u)\n",
          dev_id, lane, fpga_port);
    }
    delete(*it);
    plog_->print(" * Delete connection data\n");
    fpga_defmap_.erase(it);
    plog_->print(" * Delete connection from map\n");
    return ret;
  }
  plog_->print(" ! Failed to %s : Not found...\n", __func__);
  return -1;  // Not found
}

int connection::server_delete_connection(
  const char *device,
  uint32_t lane,
  uint32_t fchid,
  uint32_t direct_flag)
{
  if (!device) {
    plog_->print(" ! Invalid argument : NULL\n");
    return -1;
  }
  uint32_t dev_id;
  int ret = fpga_get_dev_id(device, &dev_id);
  if (ret < 0) {
    plog_->print(" ! Failed to fpga_get_dev_id(%s) in %d\n", device, ret);
    return ret;
  }
  plog_->print(" * Succeed to fpga_get_dev_id(%s) : dev_id(%u)\n", device, dev_id);
  for (auto it = fpga_defmap_.begin(); it != fpga_defmap_.end(); it++) {
    char *dev = (*it)->chain.device;
    if (strcmp(dev, device) != 0) {
      // '/dev/xpcie_'=11word
      if (dev[0] == '/' && strlen(dev) > 11 && device[0] != '/') {
        if (strcmp(&dev[11], device) != 0)
          continue;
      } else if (dev[0] != '/' && strlen(device) > 11 && device[0] == '/') {
        if (strcmp(dev, &device[11]) != 0)
          continue;
      } else {
        continue;
      }
    }
    if ((*it)->chain.lane != lane) continue;
    if ((*it)->chain.fchid != fchid) continue;
    if ((*it)->chain.direct_flag != direct_flag) continue;
    if (direct_flag) {
      // ingress
      ret = fpga_chain_disconnect_ingress(dev_id, lane, fchid);
      if (ret < 0)
        plog_->print(" ! Failed to fpga_chain_disconnect_ingress(%u,%u,%u) in %d\n",
          dev_id, lane, fchid, ret);
      else
        plog_->print(" * Succeed to fpga_chain_disconnect_ingress(%u,%u,%u)\n",
          dev_id, lane, fchid);
    } else {
      ret = fpga_chain_disconnect_egress(dev_id, lane, fchid);
      if (ret < 0)
        plog_->print(" ! Failed to fpga_chain_disconnect_egress(%u,%u,%u) in %d\n",
          dev_id, lane, fchid, ret);
      else
        plog_->print(" * Succeed to fpga_chain_disconnect_egress(%u,%u,%u)\n",
          dev_id, lane, fchid);
    }
    delete(*it);
    plog_->print(" * Delete connection data\n");
    fpga_defmap_.erase(it);
    plog_->print(" * Delete connection from map\n");
    return ret;
  }
  plog_->print(" ! Failed to %s : Not found...\n", __func__);
  return -1;  // Not found
}


int connection::server_delete_connection(
  datadef *request)
{
  plog_->print(" * Start : %s\n", __func__);
  int ret = -1;
  if (request->lldma.device[0] != '\0') {
    plog_->print(" # Cleanup for LLDMA\n");
    ret = server_delete_connection(
      request->lldma.connector_id);
  } else if (request->ptu.device[0] != '\0') {
    plog_->print(" # Cleanup for PTU\n");
    ret = server_delete_connection(
      request->ptu.device,
      request->ptu.lane,
      request->ptu.fpga_port);
  } else if (request->chain.device[0] != '\0') {
    plog_->print(" # Cleanup for CHAIN\n");
    ret = server_delete_connection(
      request->chain.device,
      request->chain.lane,
      request->chain.fchid,
      request->chain.direct_flag);
  }
  if (ret < 0) {
    plog_->print(" ! Failed to delete connection in %d\n", ret);
    return ret;
  }
  plog_->print(" * Succeed to delete connection\n");
  plog_->print(" * End   : %s\n", __func__);
  return 0;
}


int connection::server_exec_request(
  int fd_accept)
{
  plog_->print(" * Start : %s\n", __func__);

  int ret;
  datadef *request = new datadef;

  // recv
  ret = common_recv(fd_accept, request, sizeof(*request));
  if (ret < 0) {
    plog_->print(" ! Failed to recv request\n");
    delete request;
    return -1;
  }

  int response = 0;
  switch (request->cmd) {
    case FLOWDEF_CMD_QUIT:
      plog_->print(" # Recieve Quit command\n");
      response = 1;
      delete request;
      break;
    case FLOWDEF_CMD_CREATE:
      plog_->print(" # Recieve Create command\n");
      ret = server_create_connection(request);
      if (ret < 0) {
        plog_->print(" ! Failed to create connection\n");
        response = -1;
        delete request;
      } else {
        plog_->print(" * Succeed to create connection\n");
        // request is stored, so should not be deleted
      }
      break;
    case FLOWDEF_CMD_DELETE:
      plog_->print(" # Recieve Delete command\n");
      ret = server_delete_connection(request);
      if (ret < 0) {
        plog_->print(" ! Failed to delete connection\n");
        response = -1;
      } else {
        plog_->print(" * Succeed to delete connection\n");
      }
      delete request;
      break;
    default:
      plog_->print(" ! Failed to recv request\n");
      response = -1;
      delete request;
      break;
  }

  ret = common_send(fd_accept, &response, sizeof(response));
  if (ret < 0) {
    plog_->print(" ! Failed to send response\n");
    return -1;
  }

  plog_->print(" * End   : %s\n", __func__);
  return response == 1 ? 1 : 0;
}


int connection::server_wait_request(
  void)
{
  plog_->print(" * Start : %s\n", __func__);
  do {
    int fd_accept = server_accept();
    if (fd_accept < 0) {
      plog_->print(" ! Failed to accept client\n", __func__);
      continue;
    }

    int ret = server_exec_request(fd_accept);
    server_close(fd_accept);
    if (ret == 1)
      break;
  } while (1);

  plog_->print(" * End   : %s\n", __func__);
  return 0;
}


int connection::server(
  void)
{
  plog_->print(" * Start : %s\n", __func__);
  int ret;

  ret = server_listen();
  if (ret) {
    plog_->print(" ! Faliled to listen as server\n");
    return ret;
  }

  ret = daemon(1, 0);
  if (ret) {
    plog_->print(" ! Faliled to execute daemon\n");
    server_close();
    fpga_finalize();
    return ret;
  }

  ret = fpga_initialize();
  if (ret < 0) {
    plog_->print(" ! Faliled to initialize fpga\n");
    server_close();
    return ret;
  }

  ret = fpga_ptu_initialize();
  if (ret < 0) {
    plog_->print(" ! Faliled to initialize fpga(PTU)\n");
    server_close();
    fpga_finalize();
    return ret;
  }

  server_wait_request();

  server_close();
  fpga_finalize();

  plog_->print(" * End   : %s\n", __func__);
  return 0;
}


int connection::client_connect(
  void)
{
  plog_->print(" * Start : %s\n", __func__);

  if (client_fd_connect_ != -1) {
    plog_->print(" ! Already connecting\n");
    return -1;
  }

  int fd_connect = common_connect(server_listen_addr_, server_listen_port_);
  if (fd_connect < 0) {
    plog_->print(" ! Failed to connect to server\n");
    return -1;
  }

  client_fd_connect_ = fd_connect;

  plog_->print(" * End   : %s\n", __func__);
  return 0;
}


int connection::client_disconnect(
  void)
{
  plog_->print(" * Start : %s\n", __func__);

  if (client_fd_connect_ == -1) {
    plog_->print(" ! Not connecting to server yet\n");
    return -1;
  }

  common_close(client_fd_connect_);

  plog_->print(" * End   : %s\n", __func__);
  return 0;
}


int connection::client_request(
  datadef *request)
{
  plog_->print(" * Start : %s\n", __func__);

  int ret = client_connect();
  if (ret < 0) {
    plog_->print(" ! Failed to connect to server\n");
    return -1;
  }

  // send
  ret = common_send(client_fd_connect_, request, sizeof(*request));
  if (ret < 0) {
    plog_->print(" ! Failed to send request\n");
    return -1;
  }

  // recv
  int response;
  ret = common_recv(client_fd_connect_, &response, sizeof(response));
  if (ret < 0) {
    plog_->print(" ! Failed to recv response\n");
    return -1;
  }
  plog_->print(" # Request finished in %d\n", response);

  client_disconnect();

  plog_->print(" * End   : %s\n", __func__);
  return response;
}

int connection::client_request(
  void)
{
  plog_->print(" * Start : %s\n", __func__);

  datadef *request = new datadef;
  request->cmd = FLOWDEF_CMD_QUIT;

  int ret = client_request(request);
  delete(request);
  
  plog_->print(" * End   : %s\n", __func__);
  return ret;
}


int connection::client_request_quit(
  void)
{
  plog_->print(" # Request %s\n", __func__);
  return client_request();
}


int connection::client_request(
  const char *device,
  const char *connector_id,
  uint32_t dir,
  uint32_t chid,
  FLOWDEF_CMD cmd)
{
  plog_->print(" * Start : %s\n", __func__);

  datadef *request = new datadef;
  request->cmd = cmd;
  request->lldma.dir  = dir;
  request->lldma.chid = chid;
  strcpy(request->lldma.device, device);
  strcpy(request->lldma.connector_id, connector_id);

  int ret = client_request(request);
  delete(request);
  
  plog_->print(" * End   : %s\n", __func__);
  return ret;
}


int connection::client_request_create(
  const char *device,
  const char *connector_id,
  uint32_t dir,
  uint32_t chid)
{
  plog_->print(" # Request %s(LLDMA)\n", __func__);
  return client_request(
    device,
    connector_id,
    dir,
    chid,
    FLOWDEF_CMD_CREATE);
}


int connection::client_request_delete(
  const char *device,
  const char *connector_id,
  uint32_t dir,
  uint32_t chid)
{
  plog_->print(" # Request %s(LLDMA)\n", __func__);
  return client_request(
    device,
    connector_id,
    dir,
    chid,
    FLOWDEF_CMD_DELETE);
}


int connection::client_request(
  const char *device,
  uint32_t lane,
  uint32_t fchid,
  uint32_t extif_id,
  uint32_t cid,
  uint32_t active_flag,
  uint32_t direct_flag,
  uint32_t virtual_flag,
  uint32_t blocking_flag,
  FLOWDEF_CMD cmd)
{
  plog_->print(" * Start : %s\n", __func__);

  datadef *request = new datadef;
  request->cmd                 = cmd;
  request->chain.lane          = lane;
  request->chain.fchid         = fchid;
  request->chain.extif_id      = extif_id;
  request->chain.cid           = cid;
  request->chain.active_flag   = active_flag;
  request->chain.direct_flag   = direct_flag;
  request->chain.virtual_flag  = virtual_flag;
  request->chain.blocking_flag = blocking_flag;
  strcpy(request->chain.device, device);

  int ret = client_request(request);
  delete(request);
  
  plog_->print(" * End   : %s\n", __func__);
  return ret;
}


int connection::client_request_create(
  const char *device,
  uint32_t lane,
  uint32_t fchid,
  uint32_t extif_id,
  uint32_t cid,
  uint32_t active_flag,
  uint32_t direct_flag,
  uint32_t virtual_flag,
  uint32_t blocking_flag)
{
  plog_->print(" # Request %s(CHAIN)\n", __func__);
  return client_request(
    device,
    lane,
    fchid,
    extif_id,
    cid,
    active_flag,
    direct_flag,
    virtual_flag,
    blocking_flag,
    FLOWDEF_CMD_CREATE);
}


int connection::client_request_delete(
  const char *device,
  uint32_t lane,
  uint32_t fchid,
  uint32_t extif_id,
  uint32_t cid,
  uint32_t active_flag,
  uint32_t direct_flag,
  uint32_t virtual_flag,
  uint32_t blocking_flag)
{
  plog_->print(" # Request %s(CHAIN)\n", __func__);
  return client_request(
    device,
    lane,
    fchid,
    extif_id,
    cid,
    active_flag,
    direct_flag,
    virtual_flag,
    blocking_flag,
    FLOWDEF_CMD_DELETE);
}


int connection::client_request(
  const char *device,
  uint32_t lane,
  in_port_t fpga_port,
  in_addr_t host_addr,
  in_port_t host_port,
  bool is_server,
  FLOWDEF_CMD cmd)
{
  plog_->print(" * Start : %s\n", __func__);

  datadef *request = new datadef;
  request->cmd           = cmd;
  request->ptu.lane      = lane;
  request->ptu.fpga_port = fpga_port;
  request->ptu.host_addr = host_addr;
  request->ptu.host_port = host_port;
  request->ptu.is_server = is_server;
  strcpy(request->ptu.device, device);

  int ret = client_request(request);
  delete(request);
  
  plog_->print(" * End   : %s\n", __func__);
  return ret;
}


int connection::client_request_create(
  const char *device,
  uint32_t lane,
  in_port_t fpga_port,
  in_addr_t host_addr,
  in_port_t host_port,
  bool is_server)
{
  plog_->print(" # Request %s(PTU)\n", __func__);
  return client_request(
    device,
    lane,
    fpga_port,
    host_addr,
    host_port,
    is_server,
    FLOWDEF_CMD_CREATE);
}


int connection::client_request_delete(
  const char *device,
  uint32_t lane,
  in_port_t fpga_port,
  in_addr_t host_addr,
  in_port_t host_port,
  bool is_server)
{
  plog_->print(" # Request %s(PTU)\n", __func__);
  return client_request(
    device,
    lane,
    fpga_port,
    host_addr,
    host_port,
    is_server,
    FLOWDEF_CMD_DELETE);
}

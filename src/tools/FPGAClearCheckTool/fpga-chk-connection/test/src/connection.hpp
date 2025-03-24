/*************************************************
* Copyright 2025 NTT Corporation , FUJITSU LIMITED
*************************************************/

#ifndef CONNECTION_HPP__
#define CONNECTION_HPP__


#include <stdint.h>
#include "logger.hpp"

#include <vector>
#include <arpa/inet.h>

#define DEVICE_FILE_LEN  64
#define CONNECTOR_ID_LEN 64


class connection {
 public:
  connection(
          logger *plog = NULL);
  ~connection(
          void);

  int server(
          void);

  int client_request_quit(
          void);
  int client_request_create(
          const char *device,
          const char *connector_id,
          uint32_t dir,
          uint32_t chid);
  int client_request_create(
          const char *device,
          uint32_t lane,
          uint32_t fchid,
          uint32_t extif_id,
          uint32_t cid,
          uint32_t active_flag,
          uint32_t direct_flag,
          uint32_t virtual_flag,
          uint32_t blocking_flag);
  int client_request_create(
          const char *device,
          uint32_t lane,
          in_port_t fpga_port,
          in_addr_t host_addr,
          in_port_t host_port,
          bool is_server);
  int client_request_delete(
          const char *device,
          const char *connector_id,
          uint32_t dir,
          uint32_t chid);
  int client_request_delete(
          const char *device,
          uint32_t lane,
          uint32_t fchid,
          uint32_t extif_id,
          uint32_t cid,
          uint32_t active_flag,
          uint32_t direct_flag,
          uint32_t virtual_flag,
          uint32_t blocking_flag);
  int client_request_delete(
          const char *device,
          uint32_t lane,
          in_port_t fpga_port,
          in_addr_t host_addr,
          in_port_t host_port,
          bool is_server);

  int host_connect_fpga(
          const char *raddr,
          uint32_t rport,
          const char *laddr,
          uint32_t lport);
  int host_accept_fpga(
          const char *laddr,
          uint32_t lport);
  int host_close_fpga(
          int fd);
  int host_exec_fcc(
          const char *options);

 private:
  enum FLOWDEF_CMD {
    FLOWDEF_CMD_QUIT,
    FLOWDEF_CMD_CREATE,
    FLOWDEF_CMD_DELETE,
  };
  struct flowdef_lldma {
    flowdef_lldma(){
      device[0] = '\0';
    }
    char device[DEVICE_FILE_LEN];
    uint32_t dir;
    uint32_t chid;
    char connector_id[CONNECTOR_ID_LEN];
    void *dmainfo;
  };
  struct flowdef_ptu {
    flowdef_ptu(){
      device[0] = '\0';
    }
    char device[DEVICE_FILE_LEN];
    uint32_t lane;
    in_port_t fpga_port;
    in_addr_t host_addr;
    in_port_t host_port;
    bool is_server;
    uint32_t cid;
  };
  struct flowdef_chain {
    flowdef_chain(){
      device[0] = '\0';
    }
    char device[DEVICE_FILE_LEN];
    uint32_t lane;
    uint32_t fchid;
    uint32_t extif_id;
    uint32_t cid;
    uint32_t active_flag;
    uint32_t direct_flag;
    uint32_t virtual_flag;
    uint32_t blocking_flag;
  };
  struct datadef {
    FLOWDEF_CMD   cmd;
    flowdef_lldma lldma;
    flowdef_ptu   ptu;
    flowdef_chain chain;
  };

  int common_listen(
          const char *laddr,
          uint32_t lport);
  int common_accept(
          int fd_listen);
  int common_connect(
          const char *raddr,
          uint32_t rport,
          const char *laddr = NULL,
          uint32_t lport = 0);
  int common_send(
          int fd,
          const void *data,
          socklen_t len);
  int common_recv(
          int fd,
          void *data,
          socklen_t len);
  int common_close(
          int &fd);

  int fpga_initialize(
          void);
  int fpga_ptu_initialize(
          void);
  int fpga_finalize(
          void);

  int server_listen(
          void);
  int server_accept(
          void);
  int server_close(
          void);
  int server_close(
          int fd);
  int server_exec_request(
          int fd);
  int server_wait_request(
          void);
  int server_create_connection(
          const char *device,
          uint32_t dir,
          uint32_t chid,
          const char *connector_id,
          void **pdmainfo);
  int server_create_connection(
          const char *device,
          uint32_t lane,
          in_port_t fpga_port,
          in_addr_t host_addr,
          in_port_t host_port,
          bool is_server,
          uint32_t *cid);
  int server_create_connection(
          const char *device,
          uint32_t lane,
          uint32_t fchid,
          uint32_t extif_id,
          uint32_t cid,
          uint32_t active_flag,
          uint32_t direct_flag,
          uint32_t virtual_flag,
          uint32_t blocking_flag);
  int server_create_connection(
          datadef *request);

  int server_delete_connection(
          const char *connector_id);
  int server_delete_connection(
          const char *device,
          uint32_t lane,
          in_port_t fpga_port);
  int server_delete_connection(
          const char *device,
          uint32_t lane,
          uint32_t fchid,
          uint32_t direct_flag);
  int server_delete_connection(
          datadef *request);

  int client_connect(
          void);
  int client_disconnect(
          void);
  int client_request(
          datadef *request);
  int client_request(
          void);
  int client_request(
          const char *device,
          const char *connector_id,
          uint32_t dir,
          uint32_t chid,
          FLOWDEF_CMD cmd);
  int client_request(
          const char *device,
          uint32_t lane,
          in_port_t fpga_port,
          in_addr_t host_addr,
          in_port_t host_port,
          bool is_server,
          FLOWDEF_CMD cmd);
  int client_request(
          const char *device,
          uint32_t lane,
          uint32_t fchid,
          uint32_t extif_id,
          uint32_t cid,
          uint32_t active_flag,
          uint32_t direct_flag,
          uint32_t virtual_flag,
          uint32_t blocking_flag,
          FLOWDEF_CMD cmd);

  const char *default_filename_ = "connection.log";
  logger *plog_self_ = NULL;
  logger *plog_      = NULL;

  const char *server_listen_addr_ = "127.0.0.1";
  uint32_t    server_listen_port_ = 10000;
  int         server_fd_listen_   = -1;
  int         client_fd_connect_  = -1;

  std::vector<datadef*> fpga_defmap_;

  const uint32_t  fpga_lane_max_ = 2;
  const in_addr_t fpga_addr_     = 0xC0A80065; // 192.168.0.101
  const in_addr_t fpga_subnet_   = 0xffffff00; // 255.255.255.0
  const in_addr_t fpga_gateway_  = 0xC0A80001; // 192.168.0.1
  const uint8_t   fpga_mac_[6]   = {0x00, 0x12, 0x34, 0x56, 0x78, 0x91};

  const char *fcc_path_ = "../bin/fpga-chk-connection";
};

#endif  // CONNECTION_HPP__

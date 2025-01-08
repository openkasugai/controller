#include <arpa/inet.h>
#include <errno.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <time.h>
#include <unistd.h>

#include "util.h"

int main() {

  int ret = 0;
  char *strenv;

  int num_frames = 8;
  int width, height;

  int listen_fd = socket(AF_INET, SOCK_STREAM, 0);

  struct sockaddr_in saddr;
  memset(&saddr, '0', sizeof(saddr));
  saddr.sin_family = AF_INET;
  saddr.sin_addr.s_addr = htonl(INADDR_ANY);
  saddr.sin_port = htons(5000);

  bind(listen_fd, (struct sockaddr*)&saddr, sizeof(saddr));

  listen(listen_fd, 10);

  strenv = getenv("INPUT_HEIGHT");
  if (strenv == NULL) {
    printf("environment variable INPUT_HEIGHT must be specified\n");
    return -1;
  }
  height = atoi(strenv);

  strenv = getenv("INPUT_WIDTH");
  if (strenv == NULL) {
    printf("environment variable INPUT_WIDTH must be specified\n");
    return -1;
  }
  width = atoi(strenv);

  int num_channels = 3;
  size_t payload_size = sizeof(uint8_t) * width * height * num_channels;
  
  int conn_fd = accept(listen_fd, (struct sockaddr*)NULL, NULL);

  char *buffer = (char*)malloc(sizeof(frameheader_t) + payload_size);
  for (int count = 0; count < num_frames; count++) {
    // recv header
    frameheader_t fh;
    int rest = sizeof(frameheader_t);
    while (rest > 0) {
      int n = recv(conn_fd, &fh, sizeof(frameheader_t), 0);
      if (n < 0) {
	printf("recv failed: %d\n", errno);
	exit(1);
      }
      rest -= n;
    }
    uint32_t payload_len = fh.payload_len;
    if (payload_len == payload_size) {
      printf("OK! %u\n", payload_len);
    } else {
      printf("NG! %u\n", payload_len);
    }
    rest = payload_len;
    while (rest > 0) {
      int n = recv(conn_fd, buffer, payload_len, 0);
      if (n < 0) {
	printf("recv failed: %d\n", errno);
	exit(1);
      }
      rest -= n;
    }
  }
  close(conn_fd);
  return 0;
}

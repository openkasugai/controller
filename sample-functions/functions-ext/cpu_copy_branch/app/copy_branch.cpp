//=========================================================================
// copy_branch.cpp Copyright 2024 NTT Corporation , FUJITSU LIMITED
//-------------------------------------------------------------------------
#include <cstdint>
#include <iostream>
#include <mutex>
#include <thread>
#include <iostream> //standard input/output
#include <sys/socket.h> //address domain
#include <sys/types.h> //Socket Type
#include <arpa/inet.h> //Used for byte order conversion
#include <unistd.h> //Used for close()
#include <cstring> //string type
#include <string>
#include <stdio.h>
#include <stdlib.h>
#include <cxxabi.h>
#include <vector>
#include "image_packet.hpp"

#define BUF_NUM			32//1024//32
#define TIMEOUT_MILLISEC	500//Units: milliseconds

struct tcp_server_info_t{
  char src[128];
  char port[128];
};

struct tcp_client_info_t{
  char dst[128];
  char port[128];
  int  id;
};

struct tcp_recv_buf_info_t{
  std::mutex mtx;
  bool	pending;
  int	copy_size;
};

struct tcp_recv_info_t{
  tcp_recv_buf_info_t	*rbuf_info;
};

tcp_client_info_t *tcp_client_info;
tcp_recv_info_t g_tcp_recv_buf_info[BUF_NUM];
bool *g_timeout_skip_flag;

int g_memsize;


//prototype declaration
void tcp_server_thread(int branch_num, tcp_server_info_t tcp_server_info, uint8_t* recv_buffer);
void tcp_client_thread(tcp_client_info_t tcp_client_info, uint8_t* recv_buffer);

int main(int argc, char *argv[])
{
  std::vector<std::thread> threads;
  uint8_t* recv_buf;

  //Obtaining Distribution Numbers
  int branch_num = atoi(argv[2]);
  std::cout << "branch num : " << branch_num << std::endl; //standard output

  //memory resource amount acquisition
  g_memsize = atoi(argv[4]);
  printf("memsize = %d\n",g_memsize);
  fflush(stdout);

  //Argument Check
  if(argc != 5){
    std::cout << "Invalid arguments" << std::endl; //standard output
    exit(0);
  }
  std::cout << "arguments num : " << argc << std::endl; //standard output

  //Allocating memory for storing TCP communication information
  tcp_server_info_t tcp_server_info;
  tcp_client_info = new tcp_client_info_t[branch_num];

  //Receive flag allocation
  for(int m=0;m<BUF_NUM;m++){
    g_tcp_recv_buf_info[m].rbuf_info = new tcp_recv_buf_info_t[branch_num];
  }
	
  //Allocating buffer for TCP reception
  uint8_t *tcp_recv_buf = (uint8_t*)malloc(BUF_NUM * g_memsize);
  
  //Allocate skip flag for timeout
  g_timeout_skip_flag = new bool[branch_num];

  //TCP receive buffer initialization
  memset(tcp_recv_buf,0x00,BUF_NUM * g_memsize);

  //TCP reception flag initialization
  for(int p=0;p<BUF_NUM;p++){
    memset((void*)g_tcp_recv_buf_info[p].rbuf_info,0x00,sizeof(tcp_recv_buf_info_t)*branch_num);
  }

  //timeout skip flag initialization
  memset(g_timeout_skip_flag,0x00,sizeof(bool)*branch_num);

  //TCP communication information storage memory initialization
  memset(tcp_client_info,0x00,sizeof(tcp_client_info_t)*branch_num);

  //TCP Receive Buffer Address Passing
  recv_buf = (uint8_t *)tcp_recv_buf;

  //Acquisition of information for TCP reception
  std::string tcp_src_info = argv[1];
  int colon_pos = tcp_src_info.find(":");
  std::strcpy(tcp_server_info.src,tcp_src_info.substr(0,colon_pos).c_str());
  std::strcpy(tcp_server_info.port,tcp_src_info.substr(colon_pos+1).c_str());

  //TCP transmission information acquisition (for the number of distributions)
  std::string tcp_dst_info = argv[3];
  std::cout << "tcp_dst_info : " << tcp_dst_info << std::endl; //standard output
  int start_pos = 0;
  int comma_pos = 0;
  for(int i=0;i<branch_num;i++){
    colon_pos = tcp_dst_info.find(":",start_pos);//9->24
    comma_pos = tcp_dst_info.find(",",start_pos);//14->29

    std::strcpy(tcp_client_info[i].dst,tcp_dst_info.substr(start_pos,colon_pos-start_pos).c_str());

    if(i!=branch_num-1){
      std::strcpy(tcp_client_info[i].port,tcp_dst_info.substr(colon_pos+1,comma_pos-1).c_str());
    }
    else{
      std::strcpy(tcp_client_info[i].port,tcp_dst_info.substr(colon_pos+1).c_str());
    }
    tcp_client_info[i].id   = i;
    start_pos = comma_pos + 1;

  }

  //TCP receive thread start
  std::thread th_tcp_recv(tcp_server_thread, branch_num, tcp_server_info, recv_buf);
  //TCP send thread start
  for(int id=0;id<branch_num;id++){
    threads.emplace_back(std::thread(tcp_client_thread, tcp_client_info[id], recv_buf));
  }

  //Waiting for Thread Completion
  th_tcp_recv.join();
  for(auto& thread : threads){
    thread.join();
  }

  printf("                                    ↓frame_index\n");
  fflush(stdout);
  uint8_t value_check;
  for(int j=0;j<BUF_NUM;j++){
    printf("BufNo = %d\n",j);
    fflush(stdout);
    for(int k=0;k<g_memsize;k++){
      memcpy(&value_check,recv_buf+j*(g_memsize)+k,sizeof(uint8_t));
      printf("%2x ",(int)value_check);
      fflush(stdout);
    }
    printf("\n");
    fflush(stdout);
  }

  //Freeing dynamically allocated space
  free(tcp_recv_buf);
  delete[] tcp_client_info;
  for(int n=0;n<BUF_NUM;n++){
    delete[] g_tcp_recv_buf_info[n].rbuf_info;
  }
  delete[] g_timeout_skip_flag;

  std::cout << "ALL THREAD END" << std::endl;

  return 0;
}





// TCP receive thread
void tcp_server_thread(int branch_num, tcp_server_info_t tcp_server_info, uint8_t* recv_buffer)
{
  int sockfd;
  int client_sockfd;
  struct sockaddr_in addr;
  socklen_t len = sizeof( struct sockaddr_in );
  struct sockaddr_in from_addr;
  uint8_t* buf;
  int rlen;
  int recvsize=0;
  int offset=0;
  int buf_idx=0;
  uint32_t datasize=0;
  uint8_t frame_header[sizeof(ImagePacketHeader)]={};
  ImagePacketHeader* header;
  int divide_num;
  int32_t last_data_size;
  bool buf_in_use = false;

  std::chrono::system_clock::time_point start_time, current_time;

  printf("[TCP_S] THREAD START : TCP SERVER (ip:%s, port:%d)\n",tcp_server_info.src,atoi(tcp_server_info.port));
  fflush(stdout);

  // socket generation
  if( ( sockfd = socket( AF_INET, SOCK_STREAM, 0 ) ) < 0 ) {
    perror("socket failed");
    fflush(stderr);
  }
	
  // Standby IP/port number setting
  addr.sin_family = AF_INET;
  addr.sin_port = htons( atoi(tcp_server_info.port) );
  addr.sin_addr.s_addr = inet_addr(tcp_server_info.src);

  // Bind
  if( bind( sockfd, (struct sockaddr *)&addr, sizeof( addr ) ) < 0 ) {
    perror("bind failed");
    fflush(stderr);
  }
 
  // Waiting to receive
  if( listen( sockfd, SOMAXCONN ) < 0 ) {
    perror("listen failed");
    fflush(stderr);
  }
 
  // Waiting for connect requests from clients
  if( ( client_sockfd = accept( sockfd, (struct sockaddr *)&from_addr, &len ) ) < 0 ) {
    perror("accept failed");
    fflush(stderr);
  }
	
  //initial buffer pointer
  buf = (uint8_t *)recv_buffer;

  while (true) {

    // Waiting to receive
    recvsize =0;
    //Acquisition of start time for timeout
    start_time = std::chrono::system_clock::now();
    //Wait for the receive buffer to become free
    do{
      //Get current time for timeout
      current_time = std::chrono::system_clock::now();
      //When the timeout period has elapsed from the start of the wait
      if(std::chrono::duration_cast<std::chrono::milliseconds>(current_time-start_time).count()>TIMEOUT_MILLISEC){
	for(int r=0;r<branch_num;r++){
	  //Continue to skip destinations that have timed out
	  if((!g_timeout_skip_flag[r])&&(g_tcp_recv_buf_info[buf_idx].rbuf_info[r].pending)){
	    printf("Branch No.%d : TIMEOUT !!!\n",r);
	    fflush(stdout);
	  }
	  g_timeout_skip_flag[r] = g_tcp_recv_buf_info[buf_idx].rbuf_info[r].pending;
	}
      }
      //Check whether transmission to all distribution destinations has been completed (whether the buffer can be overwritten)
      buf_in_use = false;
      for(int j=0;j<branch_num;j++){
	//Skip when timeout skip flag is 1 (treat as pending= false)
	buf_in_use = buf_in_use | (g_tcp_recv_buf_info[buf_idx].rbuf_info[j].pending & ~g_timeout_skip_flag[j]);
      }
    }
    while(buf_in_use);

    //Receive frame header when receive buffer is free
    rlen = recv( client_sockfd, buf+(buf_idx*g_memsize), sizeof(ImagePacketHeader), 0 );//The return value is the number of bytes received. If 0, the connection was terminated normally.
    recvsize += rlen;
    if ( rlen == 0 ) {
      printf("[TCP_S] TCP connection end\n");
      fflush(stdout);
      break;
    } else if ( rlen == -1 ) {//SOCKET_ERROR
      perror("1st recv for frame header failed");
      fflush(stderr);
    } else {
      //Reception of remaining frame headers, if any
      while(recvsize<sizeof(ImagePacketHeader)){
	rlen = recv( client_sockfd, buf+recvsize+(buf_idx*g_memsize), sizeof(ImagePacketHeader)-recvsize, 0 );
	//printf("Recv Byte Size = %d : \n",rlen);
	if ( rlen == 0 ) {
	  printf("[TCP_S] TCP connection end\n");
	  fflush(stdout);
	  break;
	} else if ( rlen == -1 ) {//SOCKET_ERROR
	  perror("recv for frame header failed");
	  fflush(stderr);
	} else {
	  recvsize += rlen;
	}
      }
      //Stores the specified amount of received data
      memcpy(frame_header, buf+(buf_idx*g_memsize), sizeof(ImagePacketHeader));
      //Received data size display
      //printf("[TCP_S] Copy(Recieve) Byte Size = %d\n",recvsize);
      //Get Frame Header
      header = (ImagePacketHeader*)frame_header;
      //marker display for debugging
      //printf("[TCP_S] marker = %x\n",header->marker);
      //Get payload length from frame header
      datasize = header->payload_len;
      //Debug Payload Length Display
      //printf("[TCP_S] payload_len = %d\n",header->payload_len);
    }

    //Receive non-frame header (payload)
    //loop count (number of buffers)
    //Size calculation when one frame is stored in the last buffer
    if((datasize+sizeof(ImagePacketHeader))%g_memsize==0){
      divide_num = (datasize+sizeof(ImagePacketHeader))/g_memsize;
      last_data_size = g_memsize;
    }else{
      divide_num = (datasize+sizeof(ImagePacketHeader))/g_memsize+1;
      last_data_size = (datasize+sizeof(ImagePacketHeader))%g_memsize;
    }

    //Repeat until one frame is received.
    for(int i=divide_num;i>0;i--){
      //Provide an offset to include the frame header first time
      if(i==divide_num){
	offset = recvsize;
      } else {
	offset = 0;
      }
      recvsize = 0;

      //Acquisition of start time for timeout
      start_time = std::chrono::system_clock::now();
      //Wait for the receive buffer to become free// The first time (when the frame header is received), it will be dropped immediately.
      do{
	//Get current time for timeout
	current_time = std::chrono::system_clock::now();
	//When the timeout period has elapsed from the start of the wait
	if(std::chrono::duration_cast<std::chrono::milliseconds>(current_time-start_time).count()>TIMEOUT_MILLISEC){
	  for(int s=0;s<branch_num;s++){
	    //Continue to skip destinations that have timed out
	    if((!g_timeout_skip_flag[s])&&(g_tcp_recv_buf_info[buf_idx].rbuf_info[s].pending)){
	      printf("Branch No.%d : TIMEOUT !!!\n",s);
	      fflush(stdout);
	    }
	    g_timeout_skip_flag[s] = g_tcp_recv_buf_info[buf_idx].rbuf_info[s].pending;
	  }
	}
	//Check whether transmission to all distribution destinations has been completed (whether the buffer can be overwritten)
	buf_in_use = false;
	for(int l=0;l<branch_num;l++){
	  buf_in_use = buf_in_use | (g_tcp_recv_buf_info[buf_idx].rbuf_info[l].pending & ~g_timeout_skip_flag[l]);
	}
      }
      while(buf_in_use);

      if(i==1){//Final data for one frame
	rlen = recv( client_sockfd, buf+(buf_idx*g_memsize)+offset, last_data_size-offset, 0 );
	recvsize += rlen;
	while(recvsize<last_data_size-offset){
	  rlen = recv( client_sockfd, buf+(buf_idx*g_memsize)+offset+recvsize, last_data_size-offset-recvsize, 0 );
	  if ( rlen == 0 ) {
	    printf("[TCP_S] TCP connection end\n");
	    fflush(stdout);
	    break;
	  } else if ( rlen == -1 ) {//SOCKET_ERROR
	    perror("last recv for payload failed");
	    fflush(stderr);
	  } else {
	    recvsize += rlen;
	  }
	}
      }else{//one frame of non-final data
	rlen = recv( client_sockfd, buf+(buf_idx*g_memsize)+offset, g_memsize-offset, 0 );
	recvsize += rlen;
	while(recvsize<g_memsize-offset){
	  rlen = recv( client_sockfd, buf+(buf_idx*g_memsize)+offset+recvsize, g_memsize-offset-recvsize, 0 );
	  if ( rlen == 0 ) {
	    printf("[TCP_S] TCP connection end\n");
	    fflush(stdout);
	    break;
	  } else if ( rlen == -1 ) {//SOCKET_ERROR
	    perror("recv for payload failed");
	    fflush(stderr);
	  } else {
	    recvsize += rlen;
	  }
	}
      }
      //One buffer reception is complete.
      //Update receive flags in the receive buffer
      for(int k=0;k<branch_num;k++){
	tcp_recv_buf_info_t &rbuf_info_per_branch = g_tcp_recv_buf_info[buf_idx].rbuf_info[k];
	rbuf_info_per_branch.mtx.lock();
	rbuf_info_per_branch.copy_size = recvsize+offset;
	rbuf_info_per_branch.pending = true;
	rbuf_info_per_branch.mtx.unlock();
      }
      //Update buf_idx
      if(buf_idx<(BUF_NUM-1)){
	buf_idx++;
      }
      else{
	buf_idx = 0;
      }			
    }
  }
  // Socket Close
  close( client_sockfd );
  close( sockfd );
  printf("[TCP_S] THREAD END : TCP SERVER\n");
  fflush(stdout);
}


// TCP send thread
void tcp_client_thread(tcp_client_info_t tcp_client_info, uint8_t* recv_buffer)
{
  int sockfd;
  struct sockaddr_in addr;
  int send_ret = 0;
  int send_ret2 = 0;
  int buf_idx = 0;

  printf("[TCP_C] THREAD START : TCP CLIENT (ID:%d, ip:%s, port:%d)\n",tcp_client_info.id,tcp_client_info.dst,atoi(tcp_client_info.port));
  fflush(stdout);

  // Destination address and port number setting
  addr.sin_family = AF_INET;
  addr.sin_port = htons( atoi(tcp_client_info.port) );
  addr.sin_addr.s_addr = inet_addr( tcp_client_info.dst );

  // socket generation
  if( (sockfd = socket( AF_INET, SOCK_STREAM, 0) ) < 0 ) {
    perror("socket failed");
    fflush(stderr);
  }
	
  // server connection
  while (true) {
    // Wait for the connection to complete
    if(!connect( sockfd, (struct sockaddr *)&addr, sizeof( struct sockaddr_in ) )) {
      break;
    }
  }
	
  while (true){
    //If Skip is targeted, break;
    if(g_timeout_skip_flag[tcp_client_info.id]){
      break;
    }
    send_ret = 0;
    send_ret2 = 0;

    tcp_recv_buf_info_t &my_rbuf_info = g_tcp_recv_buf_info[buf_idx].rbuf_info[tcp_client_info.id];
    my_rbuf_info.mtx.lock();
    bool pending = my_rbuf_info.pending;
    int copy_size = my_rbuf_info.copy_size;
    my_rbuf_info.mtx.unlock();

    if(pending){
      //TCP transmission of storage size information
      send_ret = send( sockfd, recv_buffer+(buf_idx*g_memsize), copy_size, 0 );
      if(send_ret  < 0 ) {
	perror("1st send failed");
	fflush(stderr);
      }
      else if (send_ret < copy_size) {
	// If the specified size cannot be sent, send the rest
	do{
	  //If Skip is targeted, break;
	  if(g_timeout_skip_flag[tcp_client_info.id]){
	    break;
	  }
	  send_ret2 = send( sockfd, recv_buffer+(buf_idx*g_memsize)+send_ret, copy_size-send_ret, 0 );
	  if(send_ret2  < 0 ) {
	    perror("send failed");
	    fflush(stderr);
	  } else {
	    send_ret += send_ret2;
	  }
	}
	while(send_ret < copy_size);
      }

      //Notify Receiving Thread (Flag Clear)
      my_rbuf_info.pending = false;//Ring buffer pointer for destination [i] =ID
      //Update buf_idx
      if(buf_idx<(BUF_NUM-1)){
	buf_idx++;
      }
      else{
	buf_idx = 0;
      }
      //debug break
      //break;
    }
  }

  close(sockfd);
  printf("[TCP_C] THREAD END : TCP CLIENT (ID = %d)\n",tcp_client_info.id);
  fflush(stdout);
}

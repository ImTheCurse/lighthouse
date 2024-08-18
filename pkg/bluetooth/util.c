
#include "util.h"
#include <stdlib.h>
#include <string.h>

unsigned char* serverDataBuffer;
int serverKeyflag = KEY_OFF | PASSKEY_OFF;

const int handleConnection(int node,unsigned char* data){
  serverDataBuffer = (unsigned char*)malloc(strlen(data)+1);
  strcpy(serverDataBuffer, "Recieved your message, and stored it!");
  
  printf("Received: %s",data);
  printf("Sending back: %s",serverDataBuffer);
  write_node(node,serverDataBuffer,strlen(serverDataBuffer));
  free(serverDataBuffer);
  serverDataBuffer = (unsigned char*)malloc(strlen(data)+1);
  strcpy(serverDataBuffer, data);
  return(SERVER_CONTINUE);

}











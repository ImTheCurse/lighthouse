#ifndef UTIL_H_SERVER
#define UTIL_H_SERVER

#include <string.h>
#include <stdio.h>
#include "btlib.h"

extern unsigned char* serverDataBuffer;
extern int serverKeyflag;
const int handleConnection(int node,unsigned char* data);

#endif // !UTIL_H

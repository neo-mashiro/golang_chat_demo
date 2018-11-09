/*
g++ -o demo_client demo_client.cpp 
./demo_client 192.168.32.235 8181 name1
*/

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netdb.h>
#include <netinet/in.h>
#include <errno.h>
#include <string.h>
#include <arpa/inet.h>

struct Message {
	int PayloadLen;
	char SenderName[16];
	int MsgDataLen;
	char MsgData[];
};


int main(int argc,char* argv[])
{
    struct sockaddr_in serverAddr;
    int nFd = 0;
    int nRet = 0;
    int nReadLen = 0;
    struct Message msg;


    /* 创建套接字描述符 */
    nFd = socket(AF_INET,SOCK_STREAM,0);
    if (-1 == nFd)
    {
        perror("socket:");
        return -1;
    }

    serverAddr.sin_family = AF_INET;
    serverAddr.sin_addr.s_addr = inet_addr(argv[1]);
    serverAddr.sin_port = htons(atoi(argv[2]));

    /* 和服务器端建立连接 */
    nRet = connect(nFd,(struct sockaddr*)&serverAddr,sizeof(serverAddr));
    if (nRet == -1)
    {
        printf("connect:");
        return -1;
    }

    int i= 1;
    while(i)
    {
        memset(&msg,0,sizeof(struct Message));

        strcpy(msg.SenderName, argv[3]);
        // msg.SenderNameLen = strlen(msg.SenderName);
        sprintf(msg.MsgData, "this is message %d. End.", i);
        msg.MsgDataLen = strlen(msg.MsgData);
        msg.PayloadLen = sizeof(struct Message) - sizeof(int) + msg.MsgDataLen;
        
        ssize_t ss = send(nFd, &msg, sizeof(msg) + msg.MsgDataLen, 0);
        printf("send %ld bytes, count:%d, msg:%s[%d], payloadLen:%d\n", ss, i, msg.MsgData, msg.MsgDataLen, msg.PayloadLen);
        
        ssize_t rs = recv(nFd, &msg.PayloadLen, sizeof(msg.PayloadLen), MSG_DONTWAIT);
        if(rs == sizeof(msg.PayloadLen))
        {
            rs = recv(nFd, msg.SenderName, msg.PayloadLen, 0);
            printf("recv %s sent: %s\n", msg.SenderName, msg.MsgData);
        }
            
        
        i++;
        sleep(5);
    }

}

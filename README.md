## Introduce

server.go and client.go make up a chat demo.  
demo_client.cpp is C code client, it's use for test golang server.  

  
While any client is connecting this server, it would send message to other client,  
meanwhile, it can also receive message from other client.

:relieved:
1. the message protocol use popular C/C++ binary format.
2. no security check.
3. code is Goroutine-safe.
4. Golang code has distinct C++ style. :yum:

Golang Newbie.
debug with Goland 1.11

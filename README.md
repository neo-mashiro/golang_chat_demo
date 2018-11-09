based golang 1.11
server.go and client.go make up a chat demo.
demo_client.cpp is C code client, it's use for test golang server.



++++++++++   tcp  +++++++++++
+ server + ------ + client1 +
++++++++++        +++++++++++
    |     \
    | tcp  \ tcp +++++++++++
    |       \____+ client2 +
+++++++++++      +++++++++++  
+ clientN +   
+++++++++++


While any client connected the server, it would send message to other client,
meanwhile, it can also receive message from other client.


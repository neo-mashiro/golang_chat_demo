// tcp chats demo.
// this is server.
package main

import (
	"fmt"
	"net"
	"sync"
)

func handleConnection(cid int64, conn net.Conn, mc chan *Message) {
	addConnMap(cid, conn)

	defer func() {
		removeConnMap(cid)
		conn.Close()
	}()

	for {
		msg, err := ReadMessage(conn)
		if err != nil {
			fmt.Println("Error to read message because of ", err)
			break
		}
		fmt.Printf("Received msg '%v' from %v, cid:%d\n",
			string(msg.MsgData[:msg.MsgDataLen]),
			string(msg.SenderName),
			cid)

		msg.Cid = cid
		mc <- msg
	}

}

func boardcastMessage(mc chan *Message) {
	for {
		select {
		case msg := <-mc:
			gMapLock.Lock()
			for k, v := range gConnMap {
				if k != msg.Cid { // find others except sender self.
					WriteMessage(v, msg)
				}
			}
			gMapLock.Unlock()
		}
	}
}

func addConnMap(cid int64, conn net.Conn) {
	gMapLock.Lock()
	defer gMapLock.Unlock()
	gConnMap[cid] = conn
	fmt.Printf("add %d\n", cid)
}
func removeConnMap(cid int64) {
	gMapLock.Lock()
	defer gMapLock.Unlock()
	delete(gConnMap, cid)
	fmt.Printf("remove %d\n", cid)
}

var gConnMap = make(map[int64]net.Conn)
var gMapLock sync.Mutex

func main() {
	svr, err := net.Listen("tcp", ":8181")
	if err != nil {
		fmt.Println("Error listening: ", err)
		return
	}
	defer svr.Close()

	var msgChannel = make(chan *Message)
	go boardcastMessage(msgChannel)

	fmt.Println("Server is running, port: 8181")
	var cid int64
	cid = 0
	for {
		conn, err := svr.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
		} else {
			cid++
			fmt.Printf("New client is coming. cid:%d\n", cid)
			go handleConnection(cid, conn, msgChannel)
		}
	}

	return
}

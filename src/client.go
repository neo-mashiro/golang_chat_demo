// tcp chats demo.
// this is client.
package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
)

// 得到一行输入并发送
func scanLineAndSend(conn net.Conn, name []byte, wg *sync.WaitGroup) error {
	defer wg.Done()
	fmt.Printf("Input and press enter to send:\n")
	for {
		var c byte
		var err error
		var rline []byte

		for err == nil {
			_, err = fmt.Scanf("%c", &c)
			if c != '\r' && c != '\n' {
				rline = append(rline, c)
			} else {
				break
			}
		}

		if len(rline) > 0 {
			if string(rline) == "quit" {
				fmt.Println("byebye.")
				break
			}
			fmt.Printf("%s: %s\n", name, rline)

			msg := &Message{}
			msg.SenderName = make([]byte, len(name))
			msg.SenderName = name
			msg.MsgData = make([]byte, len(rline))
			msg.MsgData = rline
			msg.MsgDataLen = int32(len(rline))
			msg.PayloadLen = 16 + 4 + msg.MsgDataLen // SenderName长度+int长度+消息体长度

			_, err = WriteMessage(conn, msg)
			if err != nil {
				return errors.New("write tcp error: " + err.Error())
			}
		}
	}
	return nil
}

func recvMsgAndShow(conn net.Conn, wg *sync.WaitGroup) error {
	defer wg.Done()
	for {
		msg, err := ReadMessage(conn)
		if err != nil {
			fmt.Println("Error to read message because of ", err)
			return err
		}
		fmt.Printf("[%v]: %v\n",
			string(msg.SenderName),
			string(msg.MsgData[:msg.MsgDataLen]))
	}
	return nil
}

func main() {
	argNum := len(os.Args)
	if argNum != 4 {
		fmt.Println("Usage: ./client serverip port aliasname")
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", os.Args[1]+":"+os.Args[2])
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go scanLineAndSend(conn, []byte(os.Args[3]), &wg)
	go recvMsgAndShow(conn, &wg)
	wg.Wait()

}

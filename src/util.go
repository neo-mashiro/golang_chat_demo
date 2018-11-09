package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
)

/*
 * low tcp protocol:
   ++++++++++++++++++++++++++++++++++++++++++++
   +  int32 +  char[16]  +  int32  +  char[]  +
   ++++++++++++++++++++++++++++++++++++++++++++
*/
type Message struct {
	PayloadLen int32 // sendername + msgdatalen + msgdata
	SenderName []byte
	MsgDataLen int32
	MsgData    []byte

	Cid int64
}

/*
 * 从socket读取c-type struct，转换成golang struct Message.
 */
func ReadMessage(conn net.Conn) (*Message, error) {

	// 读取长度dataSize
	buf := make([]byte, 4)
	_, err := io.ReadFull(conn, buf) // read from socket.
	if err != nil {
		return nil, err
	}

	bufReader := bytes.NewReader(buf)
	var dataSize int32
	err = binary.Read(bufReader, binary.LittleEndian, &dataSize)
	if err != nil {
		return nil, err
	}

	// 读取长度为dataSize的bytes流 到 dataBuf
	dataBuf := make([]byte, dataSize)
	_, err = io.ReadFull(conn, dataBuf) // read from socket.
	if err != nil {
		return nil, err
	}

	// dataBuf 2 struct Message
	bufReader2 := bytes.NewReader(dataBuf)
	msg := &Message{}
	nameBuf := make([]byte, 16)
	binary.Read(bufReader2, binary.LittleEndian, &nameBuf)

	binary.Read(bufReader2, binary.LittleEndian, &msg.MsgDataLen)
	msgBuf := make([]byte, msg.MsgDataLen)
	binary.Read(bufReader2, binary.LittleEndian, &msgBuf)

	msg.PayloadLen = dataSize
	// 只获取有效字节数存入golang bytes[]
	msg.SenderName = bytes.TrimRight(nameBuf, "\x00")
	msg.MsgData = msgBuf

	return msg, nil
}

/*
 * 将struct Message 序列化成 c-type struct, 进行发送
 */
func WriteMessage(conn net.Conn, msg *Message) (int, error) {

	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, msg.PayloadLen)
	nameBuf := make([]byte, 16)
	copy(nameBuf, msg.SenderName)
	// nameBuf[len(msg.SenderName)] = '\x00'
	binary.Write(buffer, binary.LittleEndian, nameBuf)
	binary.Write(buffer, binary.LittleEndian, msg.MsgDataLen)
	binary.Write(buffer, binary.LittleEndian, msg.MsgData)

	sl, err := conn.Write(buffer.Bytes())
	if err != nil {
		return sl, err
	} else {
		return sl, nil
	}
}

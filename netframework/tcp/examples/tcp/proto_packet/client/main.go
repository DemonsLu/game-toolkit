package main

import (
	"fmt"
	"github.com/DarthPestilane/easytcp"
	"github.com/sirupsen/logrus"
	"net"
	"netframework/tcp/examples/fixture"
	"netframework/tcp/examples/tcp/proto_packet/common"
	"time"
)

var log *logrus.Logger
var successConn = 0

func init() {
	log = logrus.New()
	log.SetLevel(logrus.DebugLevel)
}

func connAndSend() {
	conn, err := net.Dial("tcp", fixture.ServerAddr)
	if err != nil {
		//fmt.Println(err)
		return
	}

	packer := &common.CustomPacker{}
	codec := &easytcp.ProtobufCodec{}

	go func() {
		for {
			var id = common.ID_FooReqID
			req := &common.FooReq{
				Bar: "bar",
				Buz: 22,
			}
			data, err := codec.Encode(req)
			if err != nil {
				panic(err)
			}
			packedMsg, err := packer.Pack(easytcp.NewMessage(id, data))
			if err != nil {
				panic(err)
			}
			if _, err := conn.Write(packedMsg); err != nil {
				panic(err)
			}
			//log.Debugf("send | id: %d; size: %d; data: %s", id, len(data), req.String())
			time.Sleep(time.Second)
		}
	}()

	successConn++
	fmt.Println("success conn count", successConn)
	for {
		msg, err := packer.Unpack(conn)
		if err != nil {
			panic(err)
		}
		var respData common.FooResp
		if err := codec.Decode(msg.Data(), &respData); err != nil {
			panic(err)
		}
		//log.Infof("recv | id: %d; size: %d; data: %s", msg.ID(), len(msg.Data()), respData.String())
	}
}

func main() {
	ch := make(chan struct{})
	//for i := 0; i < 30000; i++ {
	for {
		time.Sleep(time.Millisecond)
		go connAndSend()
	}
	<-ch
}

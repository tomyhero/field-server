package handler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tomyhero/go-tcp_server/client"
	"github.com/tomyhero/go-tcp_server/context"
	"github.com/tomyhero/go-tcp_server/server"
	"github.com/ugorji/go/codec"
	"reflect"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	port := 8080
	config := &server.ServerConfig{Port: port}
	assert.Equal(t, 1, 1)

	sv := server.NewServer(config)
	handlers := make([]context.IHandler, 1)
	handlers[0] = NewFieldHandler()
	sv.Setup(handlers)
	go sv.Run()
	time.Sleep(100 * time.Millisecond)

	var h = new(codec.MsgpackHandle)
	h.MapType = reflect.TypeOf(map[string]interface{}{})
	h.RawToString = true
	cl := client.Client{
		CDataManager: &context.CDataManager{CodecHandle: h},
	}

	err := cl.Connect(fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("connect Fail", err)
		t.Fail()
		return
	}

	err = cl.Send(&context.CData{
		Header: map[string]interface{}{"CMD": "field_login"},
		Body:   map[string]interface{}{},
	})

	if err != nil {
		fmt.Println("Login Fail", err)
		t.Fail()
		return
	}
	res, err := cl.Receive()
	if err != nil {
		fmt.Println("receive error", err)
		return
	}
	accessToken := res.Body["AUTH_ACCESS_TOKEN"]
	assert.NotNil(t, accessToken)

	go ReceiveHandler(&cl, t)

	err = cl.Send(&context.CData{
		Header: map[string]interface{}{"CMD": "field_UpdateStatus", "AUTH_ACCESS_TOKEN": accessToken},
		Body:   map[string]interface{}{"uid": "12234", "foo": "a"}, // 必須はuidだけ
	})

	//	cl.Disconnect()
	assert.Equal(t, 1, 1)

	time.Sleep(1000 * time.Millisecond)

}

func ReceiveHandler(cl *client.Client, t *testing.T) {
	defer func() {
		cl.Disconnect()
	}()
	for {
		cdata, err := cl.Receive()
		if err != nil {
			fmt.Println("ERROR", err)
			return
		} else {
			fmt.Println("RECEIVE FROM SERVER", cdata)
		}

	}

}

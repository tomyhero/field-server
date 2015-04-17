package handler

import (
	"fmt"
	"github.com/tomyhero/go-tcp_server/authorizer"
	"github.com/tomyhero/go-tcp_server/context"
	"github.com/ugorji/go/codec"
	"net"
	"reflect"
	"time"
)

func (h *FieldHandler) Prefix() string {
	return "field"
}

type FieldHandler struct {
	Authorizer context.IAuthorizer
	quit       chan bool
}

func (h *FieldHandler) GetAuthorizer() context.IAuthorizer {
	return h.Authorizer
}
func NewFieldHandler() *FieldHandler {
	return &FieldHandler{Authorizer: authorizer.AccessToken{}}
}
func (h *FieldHandler) HookInitialize(database map[string]interface{}, conns map[net.Conn]interface{}) {
	h.quit = make(chan bool)

	database["DATA"] = map[string]interface{}{}

	go h.broadcast(database, conns)
}

func (h *FieldHandler) HookDestroy(database map[string]interface{}) {
	close(h.quit)
}
func (h *FieldHandler) HookBeforeExecute(c *context.Context) {
}
func (h *FieldHandler) HookAfterExecute(c *context.Context) {

}

func (h *FieldHandler) HookDisconnect(conn net.Conn, database map[string]interface{}, conns map[net.Conn]interface{}) {
	// hack getting session data
	sessionID := fmt.Sprint(conns[conn].(map[string]interface{})["session_id"])
	session, ok := database[h.Prefix()].(map[string]interface{})[sessionID].(map[string]interface{})
	//	fmt.Println(sessionID)
	if ok {
		// todo remove data
		uid := session["uid"].(string)
		data := database["DATA"].(map[string]interface{})
		delete(data, uid)
	}

}

func (h *FieldHandler) ActionUpdateStatus(c *context.Context) (*context.Context, error) {
	uid := c.Req.Body["uid"].(string)
	c.Session["uid"] = uid
	data := c.Database["DATA"].(map[string]interface{})
	data[uid] = c.Req.Body

	return c, nil
}

func (h *FieldHandler) broadcast(database map[string]interface{}, conns map[net.Conn]interface{}) {

	var hl = new(codec.MsgpackHandle)
	hl.MapType = reflect.TypeOf(map[string]interface{}{})
	hl.RawToString = true
	cm := context.CDataManager{CodecHandle: hl}
	for {
		select {
		case <-h.quit:
			// ちゃんとうごいてないかもまぁいいや
			fmt.Println("quit Command Received")
			return
		default:
		}

		data := database["DATA"].(map[string]interface{})

		tmp := []map[string]interface{}{}
		for _, d := range data {
			tmp = append(tmp, d.(map[string]interface{}))
		}

		cdata := context.CData{Header: map[string]interface{}{"CMD": "broadcast"}, Body: map[string]interface{}{"LIST": tmp}}
		for conn, _ := range conns {
			err := cm.Send(conn, cdata.GetData())
			if err != nil {
				fmt.Println("Fail to send : %s", err)
			}
		}
		// ここ
		time.Sleep(300 * time.Millisecond)

	}
}

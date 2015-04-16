package handler

import (
	"fmt"
	"github.com/tomyhero/go-tcp_server/authorizer"
	"github.com/tomyhero/go-tcp_server/context"
	"net"
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
func (h *FieldHandler) HookInitialize(database map[string]interface{}) {
	h.quit = make(chan bool)
	go h.broadcast(database)
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
	session, ok := database[h.Prefix()].(map[string]interface{})[sessionID]
	fmt.Println(sessionID)
	if ok {
		// todo remove data
		fmt.Println(session)
	}

}

func (h *FieldHandler) ActionUpdateStatus(c *context.Context) (*context.Context, error) {
	fmt.Println(c.Req)
	return c, nil
}

func (h *FieldHandler) broadcast(database map[string]interface{}) {
	for {
		select {
		case <-h.quit:
			// ちゃんとうごいてないかもまぁいいや
			fmt.Println("quit Command Received")
			return
		default:
		}

		fmt.Println(database)
		time.Sleep(100 * time.Millisecond)

		// OWATA : ここからconnsにアクセスするしかたがわからん

	}
}

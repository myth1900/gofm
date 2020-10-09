package gofm

import "github.com/gorilla/websocket"

type Audience interface {
	Close()
	Connect()
}
type audience struct {
	roomID int
	ws     *websocket.Conn
}

func NewAudience(roomID int) *audience {
	adcs := &audience{
		roomID: 0,
		ws:     nil,
	}
	ws := &websocket.Dialer{}
	ws.Dial("")
	return &audience{}
}

func (a *audience) Close() {

}
func (a *audience) Connect() {

}

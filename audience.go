package gofm

import (
	"fmt"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Audience interface {
	Close()
	Connect()
}
type audience struct {
	roomID int
	ws     *websocket.Conn
}

func NewAudience(roomID int) *audience {
	return &audience{
		roomID: roomID,
	}
}

func (a *audience) Close() {

}

type joinAction struct {
	Action string `json:"action"`
	Uuid   string `json:"uuid"`
	Type   string `json:"type"`
	RoomID int    `json:"room_id"`
}

func (a *audience) Connect() {
	const url = "wss://fm.missevan.com:3016/ws"
	ws := &websocket.Dialer{}
	header := make(map[string][]string)
	header["Cookie"] = a.getCookie()
	conn, resp, err := ws.Dial(url, header)
	if err != nil {
		panic(err)
	}
	ioutil.ReadAll(resp.Body)

	// join room
	conn.WriteJSON(&joinAction{
		Action: "join",
		Uuid:   uuid.NewV1().String(),
		Type:   "room",
		RoomID: a.roomID,
	})

	//
	go func() {
		for {
			_, bts, _ := conn.ReadMessage()
			fmt.Println(string(bts))
		}
	}()

	// 维持心跳
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for {
			<-ticker.C
			conn.WriteMessage(1, []byte("❤️"))
		}
	}()
}

// FM_SESS=20201009|2e30kyum2jik391yvwz42rz79; path=/; expires=Mon, 12 Oct 2020 11:06:56 GMT; secure; httponly
func (a *audience) getCookie() []string {
	const url = "https://fm.missevan.com/api/user/info"
	var rst []string
	resp, _ := http.Get(url)
	ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	strs := resp.Header.Values("Set-Cookie")
	if len(strs) != 2 {
		panic("获取 cookie 不符合预期")
	}
	for _, str := range strs {
		ss := strings.Split(str, ";")
		rst = append(rst, ss[0])
	}
	return rst
}

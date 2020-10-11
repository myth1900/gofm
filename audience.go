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

const (
	wsUrl        = "wss://fm.missevan.com:3016/ws"
	userInfoUrl  = "https://fm.missevan.com/api/user/info"
	heartBeatMsg = "❤️"
)

type Audience interface {
	Close()
	Connect()
}
type audience struct {
	roomID int
	conn   *websocket.Conn
}

func NewAudience(roomID int) *audience {
	return &audience{
		roomID: roomID,
	}
}

func (a *audience) Close() {
	if a.conn != nil {
		if err := a.conn.Close(); err != nil {
			fmt.Println(err.Error())
		}
	}

}

type joinAction struct {
	Action string `json:"action"`
	Uuid   string `json:"uuid"`
	Type   string `json:"type"`
	RoomID int    `json:"room_id"`
}

func (a *audience) keepAlive() {
	// 读消息
	go func() {
		for {
			_, _, err := a.conn.ReadMessage()
			if err != nil {
				panic(err)
			}
			//fmt.Println(string(bts))
		}
	}()

	// 维持心跳
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for {
			<-ticker.C
			if err := a.conn.WriteMessage(1, []byte(heartBeatMsg)); err != nil {
				panic(err)
			}
		}
	}()
}

func (a *audience) Connect() {
	ws := &websocket.Dialer{}
	header := make(map[string][]string)
	header["Cookie"] = a.getCookie()
	conn, resp, err := ws.Dial(wsUrl, header)
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
	a.conn = conn

}

// FM_SESS=20201009|2e30kyum2jik391yvwz42rz79; path=/; expires=Mon, 12 Oct 2020 11:06:56 GMT; secure; httponly
func (a *audience) getCookie() []string {
	var rst []string
	resp, _ := http.Get(userInfoUrl)
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

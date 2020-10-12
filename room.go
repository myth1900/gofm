package gofm

import (
	"fmt"
	"time"
)

type Room interface {
	IncAdcs(nums int) (int, error)
	DecAdcs(nums int) (int, error)
}
type room struct {
	roomID          int
	connectAdcs     chan *audience
	waitConnectAdcs chan *audience
	waitClosedAdcs  chan *audience
}

func NewLiveRoom(roomID int) Room {
	r := &room{
		roomID:          roomID,
		connectAdcs:     make(chan *audience, 100),
		waitConnectAdcs: make(chan *audience, 100),
		waitClosedAdcs:  make(chan *audience, 100),
	}
	r.backGround()
	return r
}

// IncAdcs 增加听众，返回增加后的人数
func (r *room) IncAdcs(nums int) (int, error) {
	for i := 0; i < nums; i++ {
		r.waitConnectAdcs <- NewAudience(r.roomID)
	}
	return len(r.connectAdcs) + len(r.waitConnectAdcs), nil
}

// DecAdcs 减少听众，返回减少后的人数
func (r *room) DecAdcs(nums int) (int, error) {

	// 从未连接的队列中直接减少部分人数
	var i int
	for i < nums {
		select {
		case e := <-r.waitConnectAdcs:
			e.Close()
			i++
		default:
			goto J1
		}
	}

J1:
	// 从已连接的队列中减少部分人数
	for i < nums {
		select {
		case e := <-r.connectAdcs:
			r.waitClosedAdcs <- e
			i++
		default:
			goto J2
		}
	}
J2:
	return len(r.connectAdcs) + len(r.waitConnectAdcs), nil
}

func (r *room) backGround() {
	const duration = 10 * time.Second
	// 每隔一段时间，从等待连接的队列中拉取一个进行连接
	go func() {

		ticker := time.NewTicker(duration)
		defer ticker.Stop()
		for {
			<-ticker.C
			select {
			case e := <-r.waitConnectAdcs:
				e.Connect()
				r.connectAdcs <- e
			default:
				break
			}
		}
	}()

	// 每隔一段时间，从等待关闭的队列中拉取一个人进行关闭
	go func() {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()
		for {
			<-ticker.C
			select {
			case e := <-r.waitClosedAdcs:
				e.Close()
			default:
				break
			}
		}
	}()

	// 房间信息
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			r.log()
		}
	}()
}

func (r *room) log() {
	fmt.Printf("connected:\t%d\n", len(r.connectAdcs))
	fmt.Printf("waitConnected:\t%d\n", len(r.waitConnectAdcs))
	fmt.Printf("waitClosed:\t%d\n\n", len(r.waitClosedAdcs))
}

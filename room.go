package gofm

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	MaxAudiences = 100
)

type Room interface {
	// 更新人数
	UpdateAudience(nums int) error
	// 返回人数
	Status() RoomStatus
}

func NewRoom(roomID int) Room {
	r := &room{
		roomID:          roomID,
		mu:              &sync.Mutex{},
		connectedAdcs:   make(chan Audience, MaxAudiences),
		waitConnectAdcs: make(chan Audience, MaxAudiences),
		waitClosedAdcs:  make(chan Audience, MaxAudiences),
	}
	r.init()
	return r
}

type RoomStatus struct {
	Connected   int `json:"connected"`
	WaitConnect int `json:"wait_connected"`
	WaitClosed  int `json:"wait_closed"`
}

type room struct {
	roomID          int
	mu              *sync.Mutex
	adcNums         int
	connectedAdcs   chan Audience
	waitConnectAdcs chan Audience
	waitClosedAdcs  chan Audience
}

func (r *room) Status() RoomStatus {
	return RoomStatus{
		Connected:   len(r.connectedAdcs),
		WaitConnect: len(r.waitConnectAdcs),
		WaitClosed:  len(r.waitClosedAdcs),
	}
}

func (r *room) UpdateAudience(nums int) error {
	if nums > MaxAudiences || nums < 0 {
		return errors.New(fmt.Sprintf("超过允许的最大在线人数 %d", MaxAudiences))
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	switch {
	case r.adcNums < nums:
		r.incAdcs(nums - r.adcNums)
	case r.adcNums > nums:
		r.decAdcs(r.adcNums - nums)
	}
	r.adcNums = nums
	return nil
}

// IncAdcs 增加听众，返回增加后的人数
func (r *room) incAdcs(nums int) {
	for i := 0; i < nums; i++ {
		r.waitConnectAdcs <- NewAudience(r.roomID)
	}
}

// DecAdcs 减少听众，返回减少后的人数
func (r *room) decAdcs(nums int) {

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
		case e := <-r.connectedAdcs:
			r.waitClosedAdcs <- e
			i++
		default:
			goto J2
		}
	}
J2:
}

func (r *room) init() {
	const duration = 10 * time.Second
	// 每隔一段时间，从等待连接的队列中建立一个连接
	go func() {

		ticker := time.NewTicker(duration)
		defer ticker.Stop()
		for {
			<-ticker.C
			select {
			case e := <-r.waitConnectAdcs:
				e.Connect()
				r.connectedAdcs <- e
			default:
				break
			}
		}
	}()

	// 每隔一段时间，从等待关闭的队列中关闭一个连接
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
}

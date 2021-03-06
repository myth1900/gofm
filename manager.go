package gofm

import "sync"

type Manager interface {
	UpdateAudienceWithRoomID(roomID, nums int) error
	Status() []RoomStatus
}

func NewManager() Manager {
	return &manager{
		mu:    &sync.RWMutex{},
		rooms: make(map[int]Room),
	}
}

type manager struct {
	mu    *sync.RWMutex
	rooms map[int]Room
}

func (m *manager) Status() []RoomStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()
	rss := make([]RoomStatus, 0)
	for _, room := range m.rooms {
		rss = append(rss, room.Status())
	}
	return rss
}

func (m *manager) UpdateAudienceWithRoomID(roomID, nums int) error {
	room, err := m.getRoom(roomID)
	if err != nil {
		return err
	}
	return room.UpdateAudience(nums)
}

func (m *manager) getRoom(roomID int) (Room, error) {
	var err error
	m.mu.RLock()
	defer m.mu.RUnlock()
	room, ok := m.rooms[roomID]
	if !ok {
		m.mu.RUnlock()
		defer m.mu.RLock()
		m.mu.Lock()
		defer m.mu.Unlock()
		room, err = NewRoom(roomID)
		if err != nil {
			return nil, err
		}
		m.rooms[roomID] = room
	}
	return room, nil
}

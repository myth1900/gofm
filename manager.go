package gofm

import "sync"

type Manager interface {
	IncreaseAudience(roomID, nums int) (int, error)
	DecreaseAudience(roomID, nums int) (int, error)
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

func (m *manager) IncreaseAudience(roomID, nums int) (int, error) {
	room := m.getRoom(roomID)
	return room.IncAdcs(nums)
}

func (m manager) DecreaseAudience(roomID, nums int) (int, error) {
	room := m.getRoom(roomID)
	return room.DecAdcs(nums)
}

func (m *manager) getRoom(roomID int) Room {
	m.mu.RLock()
	defer m.mu.RUnlock()
	room, ok := m.rooms[roomID]
	if !ok {
		m.mu.RUnlock()
		m.mu.Lock()
		room = NewLiveRoom(roomID)
		m.rooms[roomID] = room
		m.mu.Unlock()
		m.mu.RLock()
	}
	return room
}

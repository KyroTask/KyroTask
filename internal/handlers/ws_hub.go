package handlers

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

// PomodoroHub manages WebSocket connections per user for real-time sync.
type PomodoroHub struct {
	mu    sync.RWMutex
	conns map[uint][]*websocket.Conn // userID → list of connections
}

var Hub = &PomodoroHub{
	conns: make(map[uint][]*websocket.Conn),
}

// Register adds a connection for a user.
func (h *PomodoroHub) Register(userID uint, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.conns[userID] = append(h.conns[userID], conn)
}

// Unregister removes a specific connection for a user.
func (h *PomodoroHub) Unregister(userID uint, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	conns := h.conns[userID]
	for i, c := range conns {
		if c == conn {
			h.conns[userID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	if len(h.conns[userID]) == 0 {
		delete(h.conns, userID)
	}
}

// BroadcastToUser sends a JSON payload to all connections for a given user.
func (h *PomodoroHub) BroadcastToUser(userID uint, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	h.mu.RLock()
	conns := make([]*websocket.Conn, len(h.conns[userID]))
	copy(conns, h.conns[userID])
	h.mu.RUnlock()

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			conn.Close()
			h.Unregister(userID, conn)
		}
	}
}

package httpsse

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/ZhangPengPaul/LightOTA/internal/model"
)

type Client struct {
	deviceID string
	m       *Manager
	writer   http.ResponseWriter
	flusher  http.Flusher
	done     chan struct{}
}

type Manager struct {
	clients map[string]map[*Client]struct{}
	mu       sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients: make(map[string]map[*Client]struct{}),
	}
}

func (m *Manager) Add(deviceID string, w http.ResponseWriter, f http.Flusher) *Client {
	client := &Client{
		deviceID: deviceID,
		m:       m,
		writer:   w,
		flusher:  f,
		done:     make(chan struct{}),
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.clients[deviceID]; !ok {
		m.clients[deviceID] = make(map[*Client]struct{})
	}
	m.clients[deviceID][client] = struct{}{}

	return client
}

func (m *Manager) Remove(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if clients, ok := m.clients[client.deviceID]; ok {
		delete(clients, client)
		close(client.done)
		if len(clients) == 0 {
			delete(m.clients, client.deviceID)
		}
	}
}

func (m *Manager) NotifyUpgrade(deviceID string, firmware *model.Firmware) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if clients, ok := m.clients[deviceID]; ok {
		message := fmt.Sprintf("event: upgrade\n"+
			"data: {\"version\":\"%s\",\"firmwareId\":\"%s\",\"md5\":\"%s\",\"fileSize\":%d}\n\n",
			firmware.Version, firmware.ID, firmware.MD5, firmware.FileSize)

		for client := range clients {
			_, err := client.writer.Write([]byte(message))
			if err != nil {
				continue
			}
			client.flusher.Flush()
		}
	}
}

func (c *Client) Wait(ctx context.Context) {
	select {
	case <-ctx.Done():
		c.m.Remove(c)
	case <-c.done:
	}
}

func (m *Manager) HandleSSE(w http.ResponseWriter, r *http.Request, deviceID string) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	client := m.Add(deviceID, w, flusher)
	defer m.Remove(client)

	ctx := r.Context()
	client.Wait(ctx)
}

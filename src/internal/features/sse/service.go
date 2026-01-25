package sse

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/ccrsxx/api-go/src/internal/api"
	"github.com/ccrsxx/api-go/src/internal/features/jellyfin"
	"github.com/ccrsxx/api-go/src/internal/features/spotify"
	"github.com/ccrsxx/api-go/src/internal/model"
	"github.com/google/uuid"
)

const (
	MaxGlobalClients = 100
	MaxClientsPerIP  = 10
)

type ClientMetadata struct {
	ID          string
	IpAddress   string
	UserAgent   string
	ConnectedAt time.Time
}

type service struct {
	mu       sync.RWMutex
	clients  map[chan string]ClientMetadata
	ipCounts map[string]int
	stopChan chan struct{}
}

var Service = &service{
	clients:  map[chan string]ClientMetadata{},
	ipCounts: map[string]int{},
}

func (s *service) IsConnectionAllowed(ip string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	isGlobalClientLimitReached := len(s.clients) >= MaxGlobalClients

	if isGlobalClientLimitReached {
		return &api.HttpError{
			Message:    "Maximum number of clients reached. Try again later.",
			StatusCode: http.StatusServiceUnavailable,
		}
	}

	isClientIPLimitReached := s.ipCounts[ip] >= MaxClientsPerIP

	if isClientIPLimitReached {
		return &api.HttpError{
			Message:    "Maximum number of clients for your IP reached. Try again later.",
			StatusCode: http.StatusTooManyRequests,
		}
	}

	return nil
}

func (s *service) AddClient(clientChan chan string, r *http.Request, ip string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	meta := ClientMetadata{
		ID:          uuid.New().String(),
		IpAddress:   ip,
		UserAgent:   r.UserAgent(),
		ConnectedAt: time.Now(),
	}

	s.clients[clientChan] = meta
	s.ipCounts[ip]++

	slog.Info("sse client connected",
		"id", meta.ID,
		"ip_address", meta.IpAddress,
		"user_agent", meta.UserAgent,
		"active_clients", len(s.clients),
	)

	welcomeMsg := `data: {"data":{"message":"Connection established. Waiting for updates..."}}` + "\n\n"

	clientChan <- welcomeMsg

	if s.stopChan == nil {
		s.startWorkerLocked()
	}
}

func (s *service) RemoveClient(clientChan chan string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	meta, exists := s.clients[clientChan]

	// Safety guard: If the client doesn't exist, just return
	// Happens if AddClient failed due to limits or server calling RemoveClient multiple times
	if !exists {
		return
	}

	delete(s.clients, clientChan)

	close(clientChan)

	s.ipCounts[meta.IpAddress]--

	if s.ipCounts[meta.IpAddress] <= 0 {
		delete(s.ipCounts, meta.IpAddress)
	}

	slog.Info("sse client disconnected",
		"id", meta.ID,
		"ip_address", meta.IpAddress,
		"user_agent", meta.UserAgent,
		"duration", time.Since(meta.ConnectedAt).String(),
		"active_clients", len(s.clients),
	)

	shouldStopPoller := len(s.clients) == 0 && s.stopChan != nil

	if shouldStopPoller {
		s.stopWorkerLocked()
	}
}

func (s *service) startWorkerLocked() {
	if s.stopChan != nil {
		slog.Warn("Attempted to start SSE poller, but it's already running")
		return
	}

	slog.Info("SSE System woke up: Starting background poller")

	s.stopChan = make(chan struct{})

	go s.pollLoop(s.stopChan)
}

func (s *service) stopWorkerLocked() {
	if s.stopChan == nil {
		slog.Warn("Attempted to stop SSE poller, but it's not running")
		return
	}

	slog.Info("SSE System idle: Stopping background poller")

	close(s.stopChan)

	s.stopChan = nil
}

func (s *service) pollLoop(stopChan chan struct{}) {
	const interval = 1 * time.Second

	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	ctx := context.Background()

	for {
		select {
		case <-stopChan:
			return // Exit goroutine / cancel polling
		case <-ticker.C:
			s.pollAndBroadcast(ctx)
		}
	}
}

func (s *service) pollAndBroadcast(ctx context.Context) {
	var wg sync.WaitGroup

	var spotifyData, jellyfinData model.CurrentlyPlaying

	wg.Add(2)

	go func() {
		defer wg.Done()

		data, err := spotify.Service.GetCurrentlyPlaying(ctx)

		if err != nil {
			spotifyData = model.NewDefaultCurrentlyPlaying(model.PlatformSpotify)
			return
		}

		spotifyData = data
	}()

	go func() {
		defer wg.Done()

		data, err := jellyfin.Service.GetCurrentlyPlaying(ctx)

		if err != nil {
			jellyfinData = model.NewDefaultCurrentlyPlaying(model.PlatformJellyfin)
			return
		}

		jellyfinData = data
	}()

	wg.Wait()

	spotifyJSON, _ := json.Marshal(map[string]model.CurrentlyPlaying{"data": spotifyData})
	jellyfinJSON, _ := json.Marshal(map[string]model.CurrentlyPlaying{"data": jellyfinData})

	msgSpotify := fmt.Sprintf("event: spotify\ndata: %s\n\n", spotifyJSON)
	msgJellyfin := fmt.Sprintf("event: jellyfin\ndata: %s\n\n", jellyfinJSON)

	// Protect map iteration with read lock
	s.mu.RLock()
	defer s.mu.RUnlock()

	for clientChan := range s.clients {
		select {
		case clientChan <- msgSpotify:
		default:
			// If the client channel is full, skip sending to avoid blocking other clients
			// Happens when the client has a slow connection
		}

		select {
		case clientChan <- msgJellyfin:
		default:
			// If the client channel is full, skip sending to avoid blocking other clients
			// Happens when the client has a slow connection
		}
	}
}

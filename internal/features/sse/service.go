package sse

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/model"
	"github.com/ccrsxx/api/internal/observability"
	"github.com/google/uuid"
)

const (
	maxClientsPerIP     = 10
	maxGlobalClients    = 100
	defaultPollInterval = 1 * time.Second
)

type clientMetadata struct {
	ID          string
	IPAddress   string
	UserAgent   string
	ConnectedAt time.Time
}

type dataFetcher interface {
	GetCurrentlyPlaying(ctx context.Context) (model.CurrentlyPlaying, error)
}

type Service struct {
	mu               sync.RWMutex
	clients          map[chan string]clientMetadata
	stopChan         chan struct{}
	appContext       context.Context
	pollInterval     time.Duration
	spotifyService   dataFetcher
	ipAddressCounts  map[string]int
	jellyfinService  dataFetcher
	navidromeService dataFetcher
}

type ServiceConfig struct {
	AppContext       context.Context
	PollInterval     time.Duration
	SpotifyService   dataFetcher
	JellyfinService  dataFetcher
	NavidromeService dataFetcher
}

func NewService(cfg ServiceConfig) *Service {
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = defaultPollInterval
	}

	return &Service{
		clients:         map[chan string]clientMetadata{},
		ipAddressCounts: map[string]int{},

		appContext:       cfg.AppContext,
		pollInterval:     cfg.PollInterval,
		spotifyService:   cfg.SpotifyService,
		jellyfinService:  cfg.JellyfinService,
		navidromeService: cfg.NavidromeService,
	}
}

func (s *Service) IsConnectionAllowed(ip string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	isGlobalClientLimitReached := len(s.clients) >= maxGlobalClients

	if isGlobalClientLimitReached {
		return &api.HTTPError{
			Message:    "Maximum number of clients reached. Try again later.",
			StatusCode: http.StatusServiceUnavailable,
		}
	}

	isClientIPLimitReached := s.ipAddressCounts[ip] >= maxClientsPerIP

	if isClientIPLimitReached {
		return &api.HTTPError{
			Message:    "Maximum number of clients for your IP reached. Try again later.",
			StatusCode: http.StatusTooManyRequests,
		}
	}

	return nil
}

func (s *Service) AddClient(ctx context.Context, clientChan chan string, ipAddress string, userAgent string) {
	sseData := s.getSSEData(ctx)

	if ctx.Err() != nil {
		slog.Warn("sse client cancelled", "ip_address", ipAddress)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	meta := clientMetadata{
		ID:          uuid.New().String(),
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		ConnectedAt: time.Now(),
	}

	s.clients[clientChan] = meta
	s.ipAddressCounts[ipAddress]++

	observability.SSEActiveClients.Set(float64(len(s.clients)))

	slog.Info("sse client connected",
		"id", meta.ID,
		"ip_address", meta.IPAddress,
		"user_agent", meta.UserAgent,
		"active_clients", len(s.clients),
	)

	welcomeMsg := `data: {"data":{"message":"Connection established. Waiting for updates..."}}` + "\n\n"

	// Send initial data immediately upon connection
	clientChan <- welcomeMsg
	clientChan <- sseData.spotify
	clientChan <- sseData.jellyfin
	clientChan <- sseData.navidrome

	if s.stopChan == nil {
		s.startWorkerLocked()
	}
}

func (s *Service) RemoveClient(ctx context.Context, clientChan chan string) {
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

	s.ipAddressCounts[meta.IPAddress]--

	if s.ipAddressCounts[meta.IPAddress] <= 0 {
		delete(s.ipAddressCounts, meta.IPAddress)
	}

	observability.SSEActiveClients.Set(float64(len(s.clients)))

	slog.Info("sse client disconnected",
		"id", meta.ID,
		"ip_address", meta.IPAddress,
		"user_agent", meta.UserAgent,
		"duration", time.Since(meta.ConnectedAt).String(),
		"active_clients", len(s.clients),
	)

	shouldStopPoller := len(s.clients) == 0 && s.stopChan != nil

	if shouldStopPoller {
		s.stopWorkerLocked()
	}
}

func (s *Service) startWorkerLocked() {
	if s.stopChan != nil {
		slog.Warn("sse poller already running")
		return
	}

	s.stopChan = make(chan struct{})

	go s.pollLoop(s.stopChan)

	slog.Info("sse poller started")
}

func (s *Service) stopWorkerLocked() {
	if s.stopChan == nil {
		slog.Warn("sse poller not running")
		return
	}

	close(s.stopChan)

	s.stopChan = nil

	slog.Info("sse poller stopped")
}

func (s *Service) pollLoop(stopChan chan struct{}) {
	interval := s.pollInterval

	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			// Note: intentionally do NOT listen to <-s.appContext.Done() here.
			// Doing so would cause a hard exit and bypass the graceful teardown.
			// Instead, let the controller catch the global shutdown, trigger RemoveClient,
			// and once the last client is removed, stopChan is safely closed here.
			return
		case <-ticker.C:
			// pass appContext down purely to instantly abort any hanging
			// network requests to Spotify/Jellyfin during a shutdown.
			s.pollAndBroadcast(s.appContext)
		}
	}
}

func (s *Service) pollAndBroadcast(ctx context.Context) {
	sseData := s.getSSEData(ctx)

	// Protect map iteration with read lock
	s.mu.RLock()
	defer s.mu.RUnlock()

	for clientChan := range s.clients {
		// Inside default block, if the client channel is full, skip sending to avoid blocking other clients
		// Happens when the client has a slow connection

		select {
		case clientChan <- sseData.spotify:
		default:
		}

		select {
		case clientChan <- sseData.jellyfin:
		default:
		}

		select {
		case clientChan <- sseData.navidrome:
		default:
		}
	}
}

type sseData struct {
	spotify   string
	jellyfin  string
	navidrome string
}

// fetchCurrentlyPlaying fetches from a single upstream, records its latency and
// outcome, and falls back to a default payload on error.
//
// The metric matters because the fallback is silent: without it, an upstream
// outage is indistinguishable from "the user isn't listening to anything".
func fetchCurrentlyPlaying(
	ctx context.Context,
	platform model.Platform,
	fetcher dataFetcher,
) model.CurrentlyPlaying {
	start := time.Now()

	data, err := fetcher.GetCurrentlyPlaying(ctx)

	result := "success"

	if err != nil {
		result = "error"
	}

	observability.UpstreamRequestDuration.
		WithLabelValues(string(platform), result).
		Observe(time.Since(start).Seconds())

	if err != nil {
		slog.WarnContext(ctx, "sse upstream fetch error",
			"platform", platform,
			"error", err,
		)

		return model.NewDefaultCurrentlyPlaying(platform)
	}

	return data
}

func (s *Service) getSSEData(ctx context.Context) sseData {
	var spotifyData, jellyfinData, navidromeData model.CurrentlyPlaying

	var wg sync.WaitGroup

	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)

	defer cancel()

	wg.Go(func() {
		spotifyData = fetchCurrentlyPlaying(timeoutCtx, model.PlatformSpotify, s.spotifyService)
	})

	wg.Go(func() {
		jellyfinData = fetchCurrentlyPlaying(timeoutCtx, model.PlatformJellyfin, s.jellyfinService)
	})

	wg.Go(func() {
		navidromeData = fetchCurrentlyPlaying(timeoutCtx, model.PlatformNavidrome, s.navidromeService)
	})

	wg.Wait()

	spotifyJSON, _ := json.Marshal(map[string]model.CurrentlyPlaying{"data": spotifyData})
	jellyfinJSON, _ := json.Marshal(map[string]model.CurrentlyPlaying{"data": jellyfinData})
	navidromeJSON, _ := json.Marshal(map[string]model.CurrentlyPlaying{"data": navidromeData})

	msgSpotify := fmt.Sprintf("event: spotify\ndata: %s\n\n", spotifyJSON)
	msgJellyfin := fmt.Sprintf("event: jellyfin\ndata: %s\n\n", jellyfinJSON)
	msgNavidrome := fmt.Sprintf("event: navidrome\ndata: %s\n\n", navidromeJSON)

	return sseData{
		spotify:   msgSpotify,
		jellyfin:  msgJellyfin,
		navidrome: msgNavidrome,
	}
}

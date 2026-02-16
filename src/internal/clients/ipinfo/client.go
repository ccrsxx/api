package ipinfo

import (
	"sync"
	"time"

	"github.com/ccrsxx/api/src/internal/config"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/ipinfo/go/v2/ipinfo/cache"
)

var (
	once   sync.Once
	client *ipinfo.Client
)

func New(token string) *ipinfo.Client {
	clientCache := ipinfo.NewCache(cache.NewInMemory().WithExpiration(5 * time.Minute))

	return ipinfo.NewClient(
		nil,
		clientCache,
		token,
	)
}

func DefaultClient() *ipinfo.Client {
	once.Do(func() {
		client = New(config.Env().IpInfoToken)
	})

	return client
}

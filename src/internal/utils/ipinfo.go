package utils

import (
	"sync"
	"time"

	"github.com/ccrsxx/api-go/src/internal/config"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/ipinfo/go/v2/ipinfo/cache"
)

var (
	once   sync.Once
	ipInfo *ipinfo.Client
)

func IPInfo() *ipinfo.Client {
	once.Do(func() {
		token := config.Env().IpInfoToken

		clientCache := ipinfo.NewCache(cache.NewInMemory().WithExpiration(5 * time.Minute))

		ipInfo = ipinfo.NewClient(
			nil,
			clientCache,
			token,
		)
	})

	return ipInfo
}

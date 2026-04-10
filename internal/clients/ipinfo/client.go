package ipinfo

import (
	"time"

	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/ipinfo/go/v2/ipinfo/cache"
)

func NewClient(token string) *ipinfo.Client {
	clientCache := ipinfo.NewCache(cache.NewInMemory().WithExpiration(5 * time.Minute))

	return ipinfo.NewClient(
		nil,
		clientCache,
		token,
	)
}

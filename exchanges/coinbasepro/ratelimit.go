package coinbasepro

import (
	"time"

	"github.com/antonk9021/qocryptotrader/exchanges/request"
)

// Coinbasepro rate limit conts
const (
	coinbaseproRateInterval = time.Second
	coinbaseproAuthRate     = 5
	coinbaseproUnauthRate   = 2
)

// GetRateLimit returns the rate limit for the exchange
func GetRateLimit() request.RateLimitDefinitions {
	return request.RateLimitDefinitions{
		request.Auth:   request.NewRateLimitWithWeight(coinbaseproRateInterval, coinbaseproAuthRate, 1),
		request.UnAuth: request.NewRateLimitWithWeight(coinbaseproRateInterval, coinbaseproUnauthRate, 1),
	}
}

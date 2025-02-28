package kline

import (
	"github.com/shopspring/decimal"
	"github.com/antonk9021/qocryptotrader/backtester/data"
	"github.com/antonk9021/qocryptotrader/backtester/eventtypes/event"
)

// Kline holds kline data and an event to be processed as
// a common.Event type
type Kline struct {
	*event.Base
	Open             decimal.Decimal
	Close            decimal.Decimal
	Low              decimal.Decimal
	High             decimal.Decimal
	Volume           decimal.Decimal
	ValidationIssues string
}

// Event is a kline data event
type Event interface {
	data.Event
	IsKline() bool
}

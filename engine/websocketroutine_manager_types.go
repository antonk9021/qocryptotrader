package engine

import (
	"errors"
	"github.com/antonk9021/qocryptotrader/exchanges/ticker"
	"sync"
	"time"

	"github.com/antonk9021/qocryptotrader/currency"
)

var (
	errNilCurrencyPairSyncer           = errors.New("nil currency pair syncer received")
	errNilCurrencyConfig               = errors.New("nil currency config received")
	errNilCurrencyPairFormat           = errors.New("nil currency pair format received")
	errNilWebsocketDataHandlerFunction = errors.New("websocket data handler function is nil")
	errNilWebsocket                    = errors.New("websocket is nil")
	errRoutineManagerNotStarted        = errors.New("websocket routine manager not started")
	errUseAPointer                     = errors.New("could not process, pass to websocket routine manager as a pointer")
)

const (
	stoppedState int32 = iota
	startingState
	readyState
)

type TickerUpdate struct {
	Exchange    string
	Pair        currency.Pair
	Candle      ticker.Price
	LastUpdated time.Time
}

// WebsocketRoutineManager is used to process websocket updates from a unified location
type WebsocketRoutineManager struct {
	state           int32
	verbose         bool
	exchangeManager iExchangeManager
	orderManager    iOrderManager
	syncer          iCurrencyPairSyncer
	currencyConfig  *currency.Config
	shutdown        chan struct{}
	dataHandlers    []WebsocketDataHandler
	TickerUpdates   chan TickerUpdate
	wg              sync.WaitGroup
	Mu              sync.RWMutex
}

// WebsocketDataHandler defines a function signature for a function that handles
// data coming from websocket connections.
type WebsocketDataHandler func(service string, incoming interface{}) error

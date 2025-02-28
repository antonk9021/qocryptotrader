package engine

import (
	"sync"
	"time"

	"github.com/antonk9021/qocryptotrader/common/key"
	"github.com/antonk9021/qocryptotrader/config"
	"github.com/antonk9021/qocryptotrader/currency"
)

// syncBase stores information
type syncBase struct {
	IsUsingWebsocket bool
	IsUsingREST      bool
	LastUpdated      time.Time
	HaveData         bool
	NumErrors        int
}

// currencyPairSyncAgent stores the sync agent info
type currencyPairSyncAgent struct {
	Key      key.ExchangePairAsset
	Pair     currency.Pair
	Created  time.Time
	trackers []*syncBase
	locks    []sync.Mutex
}

// SyncManager stores the exchange currency pair syncer object
type SyncManager struct {
	initSyncCompleted              int32
	initSyncStarted                int32
	started                        int32
	shutdown                       chan bool
	format                         currency.PairFormat
	initSyncStartTime              time.Time
	fiatDisplayCurrency            currency.Code
	websocketRoutineManagerEnabled bool
	mux                            sync.Mutex
	initSyncWG                     sync.WaitGroup
	inService                      sync.WaitGroup

	currencyPairs            map[key.ExchangePairAsset]*currencyPairSyncAgent
	tickerBatchLastRequested map[key.ExchangeAsset]time.Time

	remoteConfig    *config.RemoteControlConfig
	config          config.SyncManagerConfig
	exchangeManager iExchangeManager
}

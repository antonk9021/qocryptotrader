package engine

import (
	"errors"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/antonk9021/qocryptotrader/backtester/data"
	"github.com/antonk9021/qocryptotrader/backtester/eventhandlers/eventholder"
	"github.com/antonk9021/qocryptotrader/backtester/eventhandlers/exchange"
	"github.com/antonk9021/qocryptotrader/backtester/eventhandlers/portfolio"
	"github.com/antonk9021/qocryptotrader/backtester/eventhandlers/statistics"
	"github.com/antonk9021/qocryptotrader/backtester/eventhandlers/strategies"
	"github.com/antonk9021/qocryptotrader/backtester/funding"
	"github.com/antonk9021/qocryptotrader/backtester/report"
	"github.com/antonk9021/qocryptotrader/engine"
)

var (
	errNilConfig           = errors.New("unable to setup backtester with nil config")
	errAmbiguousDataSource = errors.New("ambiguous settings received. Only one data type can be set")
	errNoDataSource        = errors.New("no data settings set in config")
	errIntervalUnset       = errors.New("candle interval unset")
	errUnhandledDatatype   = errors.New("unhandled datatype")
	errNilData             = errors.New("nil data received")
	errLiveOnly            = errors.New("close all positions is only supported by live data type")
	errNotSetup            = errors.New("backtesting task not setup")
)

// BackTest is the main holder of all backtesting functionality
type BackTest struct {
	m                        sync.Mutex
	wg                       sync.WaitGroup
	verbose                  bool
	hasProcessedAnEvent      bool
	hasShutdown              bool
	shutdown                 chan struct{}
	MetaData                 TaskMetaData
	DataHolder               data.Holder
	LiveDataHandler          Handler
	Strategy                 strategies.Handler
	Portfolio                portfolio.Handler
	Exchange                 exchange.ExecutionHandler
	Statistic                statistics.Handler
	EventQueue               eventholder.EventHolder
	Reports                  report.Handler
	Funding                  funding.IFundingManager
	exchangeManager          *engine.ExchangeManager
	orderManager             *engine.OrderManager
	databaseManager          *engine.DatabaseConnectionManager
	hasProcessedDataAtOffset map[int64]bool
}

// TaskSummary holds details of a BackTest
// rather than passing entire contents around
type TaskSummary struct {
	MetaData TaskMetaData
}

// TaskMetaData contains details about a run such as when it was loaded
type TaskMetaData struct {
	ID                   uuid.UUID
	Strategy             string
	DateLoaded           time.Time
	DateStarted          time.Time
	DateEnded            time.Time
	Closed               bool
	ClosePositionsOnStop bool
	LiveTesting          bool
	RealOrders           bool
}

// TaskManager contains all strategy tasks
type TaskManager struct {
	m     sync.Mutex
	tasks []*BackTest
}

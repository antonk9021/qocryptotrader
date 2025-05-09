package engine

import (
	"errors"
	"sync"
	"time"

	"github.com/antonk9021/qocryptotrader/currency"
	"github.com/antonk9021/qocryptotrader/exchanges/futures"
	"github.com/antonk9021/qocryptotrader/exchanges/order"
)

// OrderManagerName is an exported subsystem name
const OrderManagerName = "orders"

// vars for the fund manager package
var (
	// ErrOrdersAlreadyExists occurs when the order already exists in the manager
	ErrOrdersAlreadyExists = errors.New("order already exists")
	// ErrOrderIDCannotBeEmpty occurs when an order does not have an ID
	ErrOrderIDCannotBeEmpty = errors.New("orderID cannot be empty")
	// ErrOrderNotFound occurs when an order is not found in the orderstore
	ErrOrderNotFound = errors.New("order does not exist")

	errNilCommunicationsManager = errors.New("cannot start with nil communications manager")
	errNilOrder                 = errors.New("nil order received")
	errFuturesTrackingDisabled  = errors.New("tracking futures positions disabled. enable it via config under orderManager activelyTrackFuturesPositions")
	orderManagerInterval        = time.Second * 10
	defaultOrderSeekTime        = -time.Hour * 24 * 365
)

type orderManagerConfig struct {
	EnforceLimitConfig     bool
	AllowMarketOrders      bool
	CancelOrdersOnShutdown bool
	LimitAmount            float64
	AllowedPairs           currency.Pairs
	AllowedExchanges       []string
	OrderSubmissionRetries int64
}

// OrderManager processes and stores orders across enabled exchanges
type OrderManager struct {
	started                       int32
	processingOrders              int32
	shutdown                      chan struct{}
	orderStore                    store
	cfg                           orderManagerConfig
	verbose                       bool
	activelyTrackFuturesPositions bool
	futuresPositionSeekDuration   time.Duration
	respectOrderHistoryLimits     bool
}

// store holds all orders by exchange
type store struct {
	m                         sync.RWMutex
	Orders                    map[string][]*order.Detail
	commsManager              iCommsManager
	exchangeManager           iExchangeManager
	wg                        *sync.WaitGroup
	futuresPositionController futures.PositionController
}

// OrderSubmitResponse contains the order response along with an internal order ID
type OrderSubmitResponse struct {
	*order.Detail
	InternalOrderID string
}

// OrderUpsertResponse contains a copy of the resulting order details and a bool
// indicating if the order details were inserted (true) or updated (false)
type OrderUpsertResponse struct {
	OrderDetails order.Detail
	IsNewOrder   bool
}

package eventholder

import (
	"github.com/antonk9021/qocryptotrader/backtester/common"
)

// Holder contains the event queue for backtester processing
type Holder struct {
	Queue []common.Event
}

// EventHolder interface details what is expected of an event holder to perform
type EventHolder interface {
	Reset() error
	AppendEvent(common.Event)
	NextEvent() common.Event
}

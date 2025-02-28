package order

import (
	"github.com/shopspring/decimal"
	"github.com/antonk9021/qocryptotrader/backtester/common"
	"github.com/antonk9021/qocryptotrader/backtester/eventtypes/event"
	"github.com/antonk9021/qocryptotrader/backtester/eventtypes/signal"
	"github.com/antonk9021/qocryptotrader/exchanges/order"
)

// Order contains all details for an order event
type Order struct {
	*event.Base
	ID                  string
	Direction           order.Side
	Status              order.Status
	ClosePrice          decimal.Decimal
	Amount              decimal.Decimal
	OrderType           order.Type
	Leverage            decimal.Decimal
	AllocatedFunds      decimal.Decimal
	BuyLimit            decimal.Decimal
	SellLimit           decimal.Decimal
	FillDependentEvent  signal.Event
	ClosingPosition     bool
	LiquidatingPosition bool
}

// Event inherits common event interfaces along with extra functions related to handling orders
type Event interface {
	common.Event
	common.Directioner
	GetClosePrice() decimal.Decimal
	GetBuyLimit() decimal.Decimal
	GetSellLimit() decimal.Decimal
	SetAmount(decimal.Decimal)
	GetAmount() decimal.Decimal
	IsOrder() bool
	GetStatus() order.Status
	SetID(id string)
	GetID() string
	IsLeveraged() bool
	GetAllocatedFunds() decimal.Decimal
	GetFillDependentEvent() signal.Event
	IsClosingPosition() bool
	IsLiquidating() bool
}

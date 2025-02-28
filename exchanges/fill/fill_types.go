package fill

import (
	"time"

	"github.com/antonk9021/qocryptotrader/currency"
	"github.com/antonk9021/qocryptotrader/exchanges/asset"
	"github.com/antonk9021/qocryptotrader/exchanges/order"
)

// Fills is used to hold data and methods related to fill dissemination
type Fills struct {
	dataHandler      chan interface{}
	fillsFeedEnabled bool
}

// Data defines fill data
type Data struct {
	ID            string
	Timestamp     time.Time
	Exchange      string
	AssetType     asset.Item
	CurrencyPair  currency.Pair
	Side          order.Side
	OrderID       string
	ClientOrderID string
	TradeID       string
	Price         float64
	Amount        float64
}

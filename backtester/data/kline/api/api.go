package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/antonk9021/qocryptotrader/backtester/common"
	"github.com/antonk9021/qocryptotrader/currency"
	exchange "github.com/antonk9021/qocryptotrader/exchanges"
	"github.com/antonk9021/qocryptotrader/exchanges/asset"
	"github.com/antonk9021/qocryptotrader/exchanges/kline"
	"github.com/antonk9021/qocryptotrader/exchanges/trade"
)

// LoadData retrieves data from a GoCryptoTrader exchange wrapper which calls the exchange's API
func LoadData(ctx context.Context, dataType int64, startDate, endDate time.Time, interval time.Duration, exch exchange.IBotExchange, fPair currency.Pair, a asset.Item) (*kline.Item, error) {
	var candles *kline.Item
	var err error
	switch dataType {
	case common.DataCandle:
		candles, err = exch.GetHistoricCandlesExtended(ctx,
			fPair,
			a,
			kline.Interval(interval),
			startDate,
			endDate)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve candle data for %v %v %v, %v", exch.GetName(), a, fPair, err)
		}
	case common.DataTrade:
		var trades []trade.Data
		trades, err = exch.GetHistoricTrades(ctx,
			fPair,
			a,
			startDate,
			endDate)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve trade data for %v %v %v, %v", exch.GetName(), a, fPair, err)
		}

		candles, err = trade.ConvertTradesToCandles(kline.Interval(interval), trades...)
		if err != nil {
			return nil, fmt.Errorf("could not convert trade data to candles for %v %v %v, %v", exch.GetName(), a, fPair, err)
		}
	default:
		return nil, fmt.Errorf("could not retrieve data for %v %v %v, %w", exch.GetName(), a, fPair, common.ErrInvalidDataType)
	}
	candles.Exchange = strings.ToLower(candles.Exchange)
	return candles, nil
}

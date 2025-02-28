package database

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/antonk9021/qocryptotrader/backtester/common"
	"github.com/antonk9021/qocryptotrader/backtester/data/kline"
	"github.com/antonk9021/qocryptotrader/currency"
	"github.com/antonk9021/qocryptotrader/exchanges/asset"
	gctkline "github.com/antonk9021/qocryptotrader/exchanges/kline"
	"github.com/antonk9021/qocryptotrader/exchanges/trade"
	"github.com/antonk9021/qocryptotrader/log"
)

var errNoUSDData = errors.New("could not retrieve USD database candle data")

// LoadData retrieves data from an existing database using GoCryptoTrader's database handling implementation
func LoadData(startDate, endDate time.Time, interval time.Duration, exchangeName string, dataType int64, fPair currency.Pair, a asset.Item, isUSDTrackingPair bool) (*kline.DataFromKline, error) {
	resp := kline.NewDataFromKline()
	switch dataType {
	case common.DataCandle:
		klineItem, err := getCandleDatabaseData(
			startDate,
			endDate,
			interval,
			exchangeName,
			fPair,
			a)
		if err != nil {
			if isUSDTrackingPair {
				return nil, fmt.Errorf("%w for %v %v %v. Please save USD candle pair data to the database or set `disable-usd-tracking` to `true` in your config. %v", errNoUSDData, exchangeName, a, fPair, err)
			}
			return nil, fmt.Errorf("could not retrieve database candle data for %v %v %v, %v", exchangeName, a, fPair, err)
		}
		resp.Item = klineItem
		for i := range klineItem.Candles {
			if klineItem.Candles[i].ValidationIssues != "" {
				log.Warnf(common.Data, "Candle validation issue for %v %v %v: %v", klineItem.Exchange, klineItem.Asset, klineItem.Pair, klineItem.Candles[i].ValidationIssues)
			}
		}
	case common.DataTrade:
		trades, err := trade.GetTradesInRange(
			exchangeName,
			a.String(),
			fPair.Base.String(),
			fPair.Quote.String(),
			startDate,
			endDate)
		if err != nil {
			return nil, err
		}
		klineItem, err := trade.ConvertTradesToCandles(
			gctkline.Interval(interval),
			trades...)
		if err != nil {
			if isUSDTrackingPair {
				return nil, fmt.Errorf("%w for %v %v %v. Please save USD pair trade data to the database or set `disable-usd-tracking` to `true` in your config. %v", errNoUSDData, exchangeName, a, fPair, err)
			}
			return nil, fmt.Errorf("could not retrieve database trade data for %v %v %v, %v", exchangeName, a, fPair, err)
		}
		resp.Item = klineItem
	default:
		if isUSDTrackingPair {
			return nil, fmt.Errorf("%w for %v %v %v. Please add USD pair data to your CSV or set `disable-usd-tracking` to `true` in your config", errNoUSDData, exchangeName, a, fPair)
		}
		return nil, fmt.Errorf("could not retrieve database data for %v %v %v, %w", exchangeName, a, fPair, common.ErrInvalidDataType)
	}
	resp.Item.Exchange = strings.ToLower(resp.Item.Exchange)

	return resp, nil
}

func getCandleDatabaseData(startDate, endDate time.Time, interval time.Duration, exchangeName string, fPair currency.Pair, a asset.Item) (*gctkline.Item, error) {
	return gctkline.LoadFromDatabase(
		exchangeName,
		fPair,
		a,
		gctkline.Interval(interval),
		startDate,
		endDate)
}

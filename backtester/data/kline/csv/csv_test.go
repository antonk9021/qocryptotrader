package csv

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/antonk9021/qocryptotrader/backtester/common"
	"github.com/antonk9021/qocryptotrader/currency"
	"github.com/antonk9021/qocryptotrader/exchanges/asset"
	gctkline "github.com/antonk9021/qocryptotrader/exchanges/kline"
)

const testExchange = "binance"

func TestLoadDataCandles(t *testing.T) {
	exch := testExchange
	a := asset.Spot
	p := currency.NewPair(currency.BTC, currency.USDT)
	_, err := LoadData(
		common.DataCandle,
		filepath.Join("..", "..", "..", "..", "testdata", "binance_BTCUSDT_24h_2019_01_01_2020_01_01.csv"),
		exch,
		gctkline.FifteenMin.Duration(),
		p,
		a,
		false)
	if !errors.Is(err, nil) {
		t.Errorf("received: %v, expected: %v", err, nil)
	}
}

func TestLoadDataTrades(t *testing.T) {
	exch := testExchange
	a := asset.Spot
	p := currency.NewPair(currency.BTC, currency.USDT)
	_, err := LoadData(
		common.DataTrade,
		filepath.Join("..", "..", "..", "..", "testdata", "binance_BTCUSDT_24h-trades_2020_11_16.csv"),
		exch,
		gctkline.FifteenMin.Duration(),
		p,
		a,
		false)
	if !errors.Is(err, nil) {
		t.Errorf("received: %v, expected: %v", err, nil)
	}
}

func TestLoadDataInvalid(t *testing.T) {
	exch := testExchange
	a := asset.Spot
	p := currency.NewPair(currency.BTC, currency.USDT)
	_, err := LoadData(
		-1,
		filepath.Join("..", "..", "..", "..", "testdata", "binance_BTCUSDT_24h-trades_2020_11_16.csv"),
		exch,
		gctkline.FifteenMin.Duration(),
		p,
		a,
		false)
	if !errors.Is(err, common.ErrInvalidDataType) {
		t.Errorf("received: %v, expected: %v", err, common.ErrInvalidDataType)
	}

	_, err = LoadData(
		-1,
		filepath.Join("..", "..", "..", "..", "testdata", "binance_BTCUSDT_24h-trades_2020_11_16.csv"),
		exch,
		gctkline.FifteenMin.Duration(),
		p,
		a,
		true)
	if !errors.Is(err, errNoUSDData) {
		t.Errorf("received: %v, expected: %v", err, errNoUSDData)
	}
}

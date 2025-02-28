package live

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/antonk9021/qocryptotrader/backtester/common"
	"github.com/antonk9021/qocryptotrader/common/convert"
	"github.com/antonk9021/qocryptotrader/currency"
	"github.com/antonk9021/qocryptotrader/engine"
	"github.com/antonk9021/qocryptotrader/exchanges/asset"
	gctkline "github.com/antonk9021/qocryptotrader/exchanges/kline"
)

const testExchange = "okx"

func TestLoadCandles(t *testing.T) {
	t.Parallel()
	interval := gctkline.OneHour
	cp := currency.NewPair(currency.BTC, currency.USDT)
	a := asset.Spot
	em := engine.NewExchangeManager()
	exch, err := em.NewExchangeByName(testExchange)
	if err != nil {
		t.Fatal(err)
	}
	pFormat := &currency.PairFormat{Uppercase: true}
	b := exch.GetBase()
	exch.SetDefaults()
	b.CurrencyPairs.Pairs = make(map[asset.Item]*currency.PairStore)
	b.CurrencyPairs.Pairs[asset.Spot] = &currency.PairStore{
		Available:     currency.Pairs{cp},
		Enabled:       currency.Pairs{cp},
		AssetEnabled:  convert.BoolPtr(true),
		RequestFormat: pFormat,
		ConfigFormat:  pFormat,
	}
	var data *gctkline.Item
	data, err = LoadData(context.Background(), time.Now().Add(-interval.Duration()*10), exch, common.DataCandle, interval.Duration(), cp, currency.EMPTYPAIR, a, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(data.Candles) == 0 {
		t.Error("expected candles")
	}
	_, err = LoadData(context.Background(), time.Now(), exch, -1, interval.Duration(), cp, currency.EMPTYPAIR, a, true)
	if !errors.Is(err, common.ErrInvalidDataType) {
		t.Errorf("received: %v, expected: %v", err, common.ErrInvalidDataType)
	}
}

func TestLoadTrades(t *testing.T) {
	t.Parallel()
	interval := gctkline.OneMin
	cp := currency.NewPair(currency.BTC, currency.USDT)
	a := asset.Spot
	em := engine.NewExchangeManager()
	exch, err := em.NewExchangeByName(testExchange)
	if err != nil {
		t.Fatal(err)
	}
	pFormat := &currency.PairFormat{Uppercase: true}
	b := exch.GetBase()
	exch.SetDefaults()
	b.CurrencyPairs.Pairs = make(map[asset.Item]*currency.PairStore)
	b.CurrencyPairs.Pairs[asset.Spot] = &currency.PairStore{
		Available:     currency.Pairs{cp},
		Enabled:       currency.Pairs{cp},
		AssetEnabled:  convert.BoolPtr(true),
		RequestFormat: pFormat,
		ConfigFormat:  pFormat,
	}
	var data *gctkline.Item
	data, err = LoadData(context.Background(), time.Now().Add(-interval.Duration()*60), exch, common.DataTrade, interval.Duration(), cp, currency.EMPTYPAIR, a, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(data.Candles) == 0 {
		t.Error("expected candles")
	}
}

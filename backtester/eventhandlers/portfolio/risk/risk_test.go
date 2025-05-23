package risk

import (
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/antonk9021/qocryptotrader/backtester/eventhandlers/portfolio/compliance"
	"github.com/antonk9021/qocryptotrader/backtester/eventhandlers/portfolio/holdings"
	"github.com/antonk9021/qocryptotrader/backtester/eventtypes/event"
	"github.com/antonk9021/qocryptotrader/backtester/eventtypes/order"
	gctcommon "github.com/antonk9021/qocryptotrader/common"
	"github.com/antonk9021/qocryptotrader/common/key"
	"github.com/antonk9021/qocryptotrader/currency"
	"github.com/antonk9021/qocryptotrader/exchanges/asset"
	gctorder "github.com/antonk9021/qocryptotrader/exchanges/order"
)

func TestAssessHoldingsRatio(t *testing.T) {
	t.Parallel()
	ratio := assessHoldingsRatio(currency.NewPair(currency.BTC, currency.USDT), []holdings.Holding{
		{
			Pair:      currency.NewPair(currency.BTC, currency.USDT),
			BaseValue: decimal.NewFromInt(2),
		},
		{
			Pair:      currency.NewPair(currency.LTC, currency.USDT),
			BaseValue: decimal.NewFromInt(2),
		},
	})
	if !ratio.Equal(decimal.NewFromFloat(0.5)) {
		t.Errorf("expected %v received %v", 0.5, ratio)
	}

	ratio = assessHoldingsRatio(currency.NewPair(currency.BTC, currency.USDT), []holdings.Holding{
		{
			Pair:      currency.NewPair(currency.BTC, currency.USDT),
			BaseValue: decimal.NewFromInt(1),
		},
		{
			Pair:      currency.NewPair(currency.LTC, currency.USDT),
			BaseValue: decimal.NewFromInt(2),
		},
		{
			Pair:      currency.NewPair(currency.DOGE, currency.USDT),
			BaseValue: decimal.NewFromInt(1),
		},
	})
	if !ratio.Equal(decimal.NewFromFloat(0.25)) {
		t.Errorf("expected %v received %v", 0.25, ratio)
	}
}

func TestEvaluateOrder(t *testing.T) {
	t.Parallel()
	r := Risk{}
	_, err := r.EvaluateOrder(nil, nil, compliance.Snapshot{})
	if !errors.Is(err, gctcommon.ErrNilPointer) {
		t.Error(err)
	}
	p := currency.NewPair(currency.BTC, currency.USDT)
	e := "binance"
	a := asset.Spot
	o := &order.Order{
		Base: &event.Base{
			Exchange:     e,
			AssetType:    a,
			CurrencyPair: p,
		},
	}
	h := []holdings.Holding{}
	r.CurrencySettings = make(map[key.ExchangePairAsset]*CurrencySettings)
	_, err = r.EvaluateOrder(o, h, compliance.Snapshot{})
	if !errors.Is(err, errNoCurrencySettings) {
		t.Error(err)
	}

	r.CurrencySettings[key.ExchangePairAsset{
		Exchange: e,
		Base:     p.Base.Item,
		Quote:    p.Quote.Item,
		Asset:    a,
	}] = &CurrencySettings{
		MaximumOrdersWithLeverageRatio: decimal.NewFromFloat(0.3),
		MaxLeverageRate:                decimal.NewFromFloat(0.3),
		MaximumHoldingRatio:            decimal.NewFromFloat(0.3),
	}

	h = append(h, holdings.Holding{
		Pair:     p,
		BaseSize: decimal.NewFromInt(1),
	})
	_, err = r.EvaluateOrder(o, h, compliance.Snapshot{})
	if !errors.Is(err, nil) {
		t.Errorf("received: %v, expected: %v", err, nil)
	}

	h = append(h, holdings.Holding{
		Pair: currency.NewPair(currency.DOGE, currency.USDT),
	})
	o.Leverage = decimal.NewFromFloat(1.1)
	r.CurrencySettings[key.ExchangePairAsset{
		Exchange: e,
		Base:     p.Base.Item,
		Quote:    p.Quote.Item,
		Asset:    a,
	}].MaximumHoldingRatio = decimal.Zero
	_, err = r.EvaluateOrder(o, h, compliance.Snapshot{})
	if !errors.Is(err, errLeverageNotAllowed) {
		t.Error(err)
	}
	r.CanUseLeverage = true
	_, err = r.EvaluateOrder(o, h, compliance.Snapshot{})
	if !errors.Is(err, errCannotPlaceLeverageOrder) {
		t.Error(err)
	}

	r.MaximumLeverage = decimal.NewFromInt(33)
	r.CurrencySettings[key.ExchangePairAsset{
		Exchange: e,
		Base:     p.Base.Item,
		Quote:    p.Quote.Item,
		Asset:    a,
	}].MaxLeverageRate = decimal.NewFromInt(33)
	_, err = r.EvaluateOrder(o, h, compliance.Snapshot{})
	if !errors.Is(err, nil) {
		t.Errorf("received: %v, expected: %v", err, nil)
	}

	r.MaximumLeverage = decimal.NewFromInt(33)
	r.CurrencySettings[key.ExchangePairAsset{
		Exchange: e,
		Base:     p.Base.Item,
		Quote:    p.Quote.Item,
		Asset:    a,
	}].MaxLeverageRate = decimal.NewFromInt(33)

	_, err = r.EvaluateOrder(o, h, compliance.Snapshot{
		Orders: []compliance.SnapshotOrder{
			{
				Order: &gctorder.Detail{
					Leverage: 3,
				},
			},
		},
	})
	if !errors.Is(err, errCannotPlaceLeverageOrder) {
		t.Error(err)
	}

	h = append(h, holdings.Holding{Pair: p, BaseValue: decimal.NewFromInt(1337)}, holdings.Holding{Pair: p, BaseValue: decimal.NewFromFloat(1337.42)})
	r.CurrencySettings[key.ExchangePairAsset{
		Exchange: e,
		Base:     p.Base.Item,
		Quote:    p.Quote.Item,
		Asset:    a,
	}].MaximumHoldingRatio = decimal.NewFromFloat(0.1)
	_, err = r.EvaluateOrder(o, h, compliance.Snapshot{})
	if !errors.Is(err, nil) {
		t.Errorf("received: %v, expected: %v", err, nil)
	}

	h = append(h, holdings.Holding{Pair: currency.NewPair(currency.DOGE, currency.LTC), BaseValue: decimal.NewFromInt(1337)})
	_, err = r.EvaluateOrder(o, h, compliance.Snapshot{})
	if !errors.Is(err, errCannotPlaceLeverageOrder) {
		t.Error(err)
	}
}

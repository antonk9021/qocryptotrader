package portfolio

import (
	"strings"

	"github.com/shopspring/decimal"
	"github.com/antonk9021/qocryptotrader/backtester/eventhandlers/exchange"
	"github.com/antonk9021/qocryptotrader/backtester/eventhandlers/portfolio/holdings"
	"github.com/antonk9021/qocryptotrader/backtester/eventhandlers/portfolio/risk"
	gctcommon "github.com/antonk9021/qocryptotrader/common"
	"github.com/antonk9021/qocryptotrader/common/key"
	"github.com/antonk9021/qocryptotrader/exchanges/asset"
	"github.com/antonk9021/qocryptotrader/exchanges/futures"
)

// Setup creates a portfolio manager instance and sets private fields
func Setup(sh SizeHandler, r risk.Handler, riskFreeRate decimal.Decimal) (*Portfolio, error) {
	if sh == nil {
		return nil, errSizeManagerUnset
	}
	if riskFreeRate.IsNegative() {
		return nil, errNegativeRiskFreeRate
	}
	if r == nil {
		return nil, errRiskManagerUnset
	}
	p := &Portfolio{}
	p.sizeManager = sh
	p.riskManager = r
	p.riskFreeRate = riskFreeRate

	return p, nil
}

// Reset returns the portfolio manager to its default state
func (p *Portfolio) Reset() error {
	if p == nil {
		return gctcommon.ErrNilPointer
	}
	p.exchangeAssetPairPortfolioSettings = make(map[key.ExchangePairAsset]*Settings)
	p.riskFreeRate = decimal.Zero
	p.sizeManager = nil
	p.riskManager = nil
	return nil
}

// SetCurrencySettingsMap ensures a map is created and no panics happen
func (p *Portfolio) SetCurrencySettingsMap(setup *exchange.Settings) error {
	if setup == nil {
		return errNoPortfolioSettings
	}
	if setup.Exchange == nil {
		return errExchangeUnset
	}
	if setup.Asset == asset.Empty {
		return errAssetUnset
	}
	if setup.Pair.IsEmpty() {
		return errCurrencyPairUnset
	}

	if p.exchangeAssetPairPortfolioSettings == nil {
		p.exchangeAssetPairPortfolioSettings = make(map[key.ExchangePairAsset]*Settings)
	}
	name := strings.ToLower(setup.Exchange.GetName())

	settings := &Settings{
		Exchange:          setup.Exchange,
		exchangeName:      name,
		assetType:         setup.Asset,
		pair:              setup.Pair,
		BuySideSizing:     setup.BuySide,
		SellSideSizing:    setup.SellSide,
		Leverage:          setup.Leverage,
		HoldingsSnapshots: make(map[int64]*holdings.Holding),
	}
	if setup.Asset.IsFutures() {
		collateralCurrency, _, err := setup.Exchange.GetCollateralCurrencyForContract(setup.Asset, setup.Pair)
		if err != nil {
			return err
		}
		futureTrackerSetup := &futures.MultiPositionTrackerSetup{
			Exchange:                  name,
			Asset:                     setup.Asset,
			Pair:                      setup.Pair,
			Underlying:                setup.Pair.Base,
			OfflineCalculation:        true,
			UseExchangePNLCalculation: setup.UseExchangePNLCalculation,
			CollateralCurrency:        collateralCurrency,
		}
		if setup.UseExchangePNLCalculation {
			futureTrackerSetup.ExchangePNLCalculation = setup.Exchange
		}
		var tracker *futures.MultiPositionTracker
		tracker, err = futures.SetupMultiPositionTracker(futureTrackerSetup)
		if err != nil {
			return err
		}
		settings.FuturesTracker = tracker
	}
	p.exchangeAssetPairPortfolioSettings[key.ExchangePairAsset{
		Exchange: name,
		Base:     setup.Pair.Base.Item,
		Quote:    setup.Pair.Quote.Item,
		Asset:    setup.Asset,
	}] = settings
	return nil
}

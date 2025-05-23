package strategies

import (
	"errors"
	"testing"

	"github.com/antonk9021/qocryptotrader/backtester/data"
	"github.com/antonk9021/qocryptotrader/backtester/eventhandlers/portfolio"
	"github.com/antonk9021/qocryptotrader/backtester/eventhandlers/strategies"
	"github.com/antonk9021/qocryptotrader/backtester/eventhandlers/strategies/base"
	"github.com/antonk9021/qocryptotrader/backtester/eventhandlers/strategies/dollarcostaverage"
	"github.com/antonk9021/qocryptotrader/backtester/eventtypes/signal"
	"github.com/antonk9021/qocryptotrader/backtester/funding"
)

func TestAddStrategies(t *testing.T) {
	t.Parallel()
	err := addStrategies(nil)
	if !errors.Is(err, errNoStrategies) {
		t.Error(err)
	}

	err = addStrategies([]strategies.Handler{&dollarcostaverage.Strategy{}})
	if !errors.Is(err, strategies.ErrStrategyAlreadyExists) {
		t.Error(err)
	}

	err = addStrategies([]strategies.Handler{&CustomStrategy{}})
	if !errors.Is(err, nil) {
		t.Error(err)
	}
}

type CustomStrategy struct {
	base.Strategy
}

func (s *CustomStrategy) Name() string {
	return "custom-strategy"
}

func (s *CustomStrategy) Description() string {
	return "this is a demonstration of loading strategies via custom plugins"
}

func (s *CustomStrategy) SupportsSimultaneousProcessing() bool {
	return true
}

func (s *CustomStrategy) OnSignal(d data.Handler, _ funding.IFundingTransferer, _ portfolio.Handler) (signal.Event, error) {
	return s.createSignal(d)
}
func (s *CustomStrategy) OnSimultaneousSignals(_ []data.Handler, _ funding.IFundingTransferer, _ portfolio.Handler) ([]signal.Event, error) {
	return nil, nil
}

func (s *CustomStrategy) createSignal(_ data.Handler) (*signal.Signal, error) {
	return nil, nil
}

func (s *CustomStrategy) SetCustomSettings(map[string]interface{}) error {
	return nil
}

func (s *CustomStrategy) SetDefaults() {}

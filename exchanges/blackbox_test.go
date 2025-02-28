package exchange_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	exchange "github.com/antonk9021/qocryptotrader/exchanges"
	shared "github.com/antonk9021/qocryptotrader/exchanges/sharedtestvalues"
)

type mockEx struct {
	shared.CustomEx
	flow chan int
}

func (m *mockEx) UpdateTradablePairs(_ context.Context, _ bool) error {
	m.flow <- 42
	return nil
}

func TestBootstrap(t *testing.T) {
	m := &mockEx{
		shared.CustomEx{},
		make(chan int, 1),
	}
	m.Features.Enabled.AutoPairUpdates = true
	err := exchange.Bootstrap(context.TODO(), m)
	assert.NoError(t, err, "Bootstrap should not error")
	assert.Equal(t, 42, <-m.flow, "UpdateTradablePairs should be called on the exchange")
}

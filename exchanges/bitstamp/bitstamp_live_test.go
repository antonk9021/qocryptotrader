//go:build mock_test_off

// This will build if build tag mock_test_off is parsed and will do live testing
// using all tests in (exchange)_test.go
package bitstamp

import (
	"log"
	"os"
	"testing"

	"github.com/antonk9021/qocryptotrader/config"
	"github.com/antonk9021/qocryptotrader/exchanges/sharedtestvalues"
)

var mockTests = false

func TestMain(m *testing.M) {
	cfg := config.GetConfig()
	err := cfg.LoadConfig("../../testdata/configtest.json", true)
	if err != nil {
		log.Fatal("Bitstamp load config error", err)
	}
	bitstampConfig, err := cfg.GetExchangeConfig("Bitstamp")
	if err != nil {
		log.Fatal("Bitstamp Setup() init error", err)
	}
	bitstampConfig.API.AuthenticatedSupport = true
	if apiKey != "" {
		bitstampConfig.API.Credentials.Key = apiKey
	}
	if apiSecret != "" {
		bitstampConfig.API.Credentials.Secret = apiSecret
	}
	if customerID != "" {
		bitstampConfig.API.Credentials.ClientID = customerID
	}
	b.SetDefaults()
	b.Websocket = sharedtestvalues.NewTestWebsocket()
	err = b.Setup(bitstampConfig)
	if err != nil {
		log.Fatal("Bitstamp setup error", err)
	}
	log.Printf(sharedtestvalues.LiveTesting, b.Name)
	os.Exit(m.Run())
}

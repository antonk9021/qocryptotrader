//go:build mock_test_off

// This will build if build tag mock_test_off is parsed and will do live testing
// using all tests in (exchange)_test.go
package gemini

import (
	"log"
	"os"
	"testing"

	"github.com/antonk9021/qocryptotrader/config"
	exchange "github.com/antonk9021/qocryptotrader/exchanges"
	"github.com/antonk9021/qocryptotrader/exchanges/sharedtestvalues"
)

var mockTests = false

func TestMain(m *testing.M) {
	cfg := config.GetConfig()
	err := cfg.LoadConfig("../../testdata/configtest.json", true)
	if err != nil {
		log.Fatal("Gemini load config error", err)
	}
	geminiConfig, err := cfg.GetExchangeConfig("Gemini")
	if err != nil {
		log.Fatal("Gemini Setup() init error", err)
	}
	geminiConfig.API.AuthenticatedSupport = true
	geminiConfig.API.Credentials.Key = apiKey
	geminiConfig.API.Credentials.Secret = apiSecret
	g.SetDefaults()
	g.Websocket = sharedtestvalues.NewTestWebsocket()
	err = g.Setup(geminiConfig)
	if err != nil {
		log.Fatal("Gemini setup error", err)
	}
	err = g.API.Endpoints.SetRunning(exchange.RestSpot.String(), geminiAPIURL)
	if err != nil {
		log.Fatalf("endpoint setting failed. key: %s, val: %s", exchange.RestSpot.String(), geminiAPIURL)
	}
	log.Printf(sharedtestvalues.LiveTesting, g.Name)
	os.Exit(m.Run())
}

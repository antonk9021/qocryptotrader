package main

import (
	"testing"

	"github.com/antonk9021/qocryptotrader/currency"
)

func TestDisruptFormatting(t *testing.T) {
	_, err := disruptFormatting(currency.EMPTYPAIR)
	if err == nil {
		t.Fatal("error cannot be nil")
	}

	_, err = disruptFormatting(currency.Pair{Quote: currency.BTC})
	if err == nil {
		t.Fatal("error cannot be nil")
	}

	p := currency.NewPair(currency.BTC, currency.USDT)

	badPair, err := disruptFormatting(p)
	if err != nil {
		t.Fatal(err)
	}

	if badPair.String() != "BTC-TEST-DELIM-usdt" {
		t.Fatal("incorrect disrupted pair")
	}
}

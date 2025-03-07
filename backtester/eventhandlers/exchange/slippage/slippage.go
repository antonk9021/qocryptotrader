package slippage

import (
	"math/rand"

	"github.com/shopspring/decimal"
	gctorder "github.com/antonk9021/qocryptotrader/exchanges/order"
	"github.com/antonk9021/qocryptotrader/exchanges/orderbook"
)

// EstimateSlippagePercentage takes in an int range of numbers
// turns it into a percentage
func EstimateSlippagePercentage(maximumSlippageRate, minimumSlippageRate decimal.Decimal) decimal.Decimal {
	if minimumSlippageRate.LessThan(decimal.NewFromInt(1)) || minimumSlippageRate.GreaterThan(decimal.NewFromInt(100)) {
		return decimal.NewFromInt(1)
	}
	if maximumSlippageRate.LessThan(decimal.NewFromInt(1)) || maximumSlippageRate.GreaterThan(decimal.NewFromInt(100)) {
		return decimal.NewFromInt(1)
	}

	// the language here is confusing. The maximum slippage rate is the lower bounds of the number,
	// eg 80 means for every dollar, keep 80%
	randSeed := int(minimumSlippageRate.IntPart()) - int(maximumSlippageRate.IntPart())
	if randSeed > 0 {
		result := int64(rand.Intn(randSeed)) //nolint:gosec // basic number generation required, no need for crypto/rand

		return maximumSlippageRate.Add(decimal.NewFromInt(result)).Div(decimal.NewFromInt(100))
	}
	return decimal.NewFromInt(1)
}

// CalculateSlippageByOrderbook will analyse a provided orderbook and return the result of attempting to
// place the order on there
func CalculateSlippageByOrderbook(ob *orderbook.Base, side gctorder.Side, allocatedFunds, feeRate decimal.Decimal) (price, amount decimal.Decimal, err error) {
	var result *orderbook.WhaleBombResult
	result, err = ob.SimulateOrder(allocatedFunds.InexactFloat64(), side == gctorder.Buy)
	if err != nil {
		return
	}
	rate := (result.MinimumPrice - result.MaximumPrice) / result.MaximumPrice
	price = decimal.NewFromFloat(result.MinimumPrice * (rate + 1))
	amount = decimal.NewFromFloat(result.Amount * (1 - feeRate.InexactFloat64()))
	return
}

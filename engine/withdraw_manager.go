package engine

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/antonk9021/qocryptotrader/common"
	dbwithdraw "github.com/antonk9021/qocryptotrader/database/repository/withdraw"
	"github.com/antonk9021/qocryptotrader/exchanges/asset"
	"github.com/antonk9021/qocryptotrader/exchanges/currencystate"
	"github.com/antonk9021/qocryptotrader/log"
	"github.com/antonk9021/qocryptotrader/portfolio/withdraw"
)

// SetupWithdrawManager creates a new withdraw manager
func SetupWithdrawManager(em iExchangeManager, pm iPortfolioManager, isDryRun bool) (*WithdrawManager, error) {
	if em == nil {
		return nil, errors.New("nil manager")
	}
	return &WithdrawManager{
		exchangeManager:  em,
		portfolioManager: pm,
		isDryRun:         isDryRun,
	}, nil
}

// SubmitWithdrawal performs validation and submits a new withdraw request to
// exchange
func (m *WithdrawManager) SubmitWithdrawal(ctx context.Context, req *withdraw.Request) (*withdraw.Response, error) {
	if m == nil {
		return nil, ErrNilSubsystem
	}
	if req == nil {
		return nil, withdraw.ErrRequestCannotBeNil
	}

	exch, err := m.exchangeManager.GetExchangeByName(req.Exchange)
	if err != nil {
		return nil, err
	}

	resp := &withdraw.Response{
		Exchange: withdraw.ExchangeResponse{
			Name: req.Exchange,
		},
		RequestDetails: *req,
	}

	// Determines if the currency can be withdrawn from the exchange
	errF := exch.CanWithdraw(req.Currency, asset.Spot)
	if errF != nil && !errors.Is(errF, currencystate.ErrCurrencyStateNotFound) { // Suppress not found error
		return nil, errF
	}

	if m.isDryRun {
		log.Warnln(log.Global, "Dry run enabled, no withdrawal request will be submitted or have an event created")
		resp.ID = withdraw.DryRunID
		resp.Exchange.Status = "dryrun"
		resp.Exchange.ID = withdraw.DryRunID.String()
	} else {
		var ret *withdraw.ExchangeResponse
		if req.Type == withdraw.Crypto {
			if !m.portfolioManager.IsWhiteListed(req.Crypto.Address) {
				return nil, withdraw.ErrStrAddressNotWhiteListed
			}
			if !m.portfolioManager.IsExchangeSupported(req.Exchange, req.Crypto.Address) {
				return nil, withdraw.ErrStrExchangeNotSupportedByAddress
			}
		}
		if req.Type == withdraw.Fiat {
			ret, err = exch.WithdrawFiatFunds(ctx, req)
			if err != nil {
				resp.Exchange.Status = err.Error()
			} else {
				resp.Exchange.Status = ret.Status
				resp.Exchange.ID = ret.ID
			}
		} else if req.Type == withdraw.Crypto {
			ret, err = exch.WithdrawCryptocurrencyFunds(ctx, req)
			if err != nil {
				resp.Exchange.Status = err.Error()
			} else {
				resp.Exchange.Status = ret.Status
				resp.Exchange.ID = ret.ID
			}
		}
	}
	dbwithdraw.Event(resp)
	if err == nil {
		withdraw.Cache.Add(resp.ID, resp)
	}
	return resp, err
}

// WithdrawalEventByID returns a withdrawal request by ID
func (m *WithdrawManager) WithdrawalEventByID(id string) (*withdraw.Response, error) {
	if m == nil {
		return nil, ErrNilSubsystem
	}
	if v := withdraw.Cache.Get(id); v != nil {
		wdResp, ok := v.(*withdraw.Response)
		if !ok {
			return nil, common.GetTypeAssertError("*withdraw.Response", v)
		}
		return wdResp, nil
	}

	l, err := dbwithdraw.GetEventByUUID(id)
	if err != nil {
		return nil, fmt.Errorf("%w %v", ErrWithdrawRequestNotFound, id)
	}
	withdraw.Cache.Add(id, l)
	return l, nil
}

// WithdrawalEventByExchange returns a withdrawal request by ID
func (m *WithdrawManager) WithdrawalEventByExchange(exchange string, limit int) ([]*withdraw.Response, error) {
	if m == nil {
		return nil, ErrNilSubsystem
	}
	_, err := m.exchangeManager.GetExchangeByName(exchange)
	if err != nil {
		return nil, err
	}

	return dbwithdraw.GetEventsByExchange(exchange, limit)
}

// WithdrawEventByDate returns a withdrawal request by ID
func (m *WithdrawManager) WithdrawEventByDate(exchange string, start, end time.Time, limit int) ([]*withdraw.Response, error) {
	if m == nil {
		return nil, ErrNilSubsystem
	}
	_, err := m.exchangeManager.GetExchangeByName(exchange)
	if err != nil {
		return nil, err
	}

	return dbwithdraw.GetEventsByDate(exchange, start, end, limit)
}

// WithdrawalEventByExchangeID returns a withdrawal request by Exchange ID
func (m *WithdrawManager) WithdrawalEventByExchangeID(exchange, id string) (*withdraw.Response, error) {
	if m == nil {
		return nil, ErrNilSubsystem
	}
	_, err := m.exchangeManager.GetExchangeByName(exchange)
	if err != nil {
		return nil, err
	}

	return dbwithdraw.GetEventByExchangeID(exchange, id)
}

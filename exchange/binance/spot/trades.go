package spot

import (
	"context"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/binance"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	Trade struct {
		Symbol          string          `json:"symbol"`
		ID              int64           `json:"id"`
		OrderID         int64           `json:"orderId"`
		OrderListID     int64           `json:"orderListId"`
		Price           decimal.Decimal `json:"price"`
		Qty             decimal.Decimal `json:"qty"`
		QuoteQty        decimal.Decimal `json:"quoteQty"`
		Commission      decimal.Decimal `json:"commission"`
		CommissionAsset string          `json:"commissionAsset"`
		Time            int64           `json:"time"`
		IsBuyer         bool            `json:"isBuyer"`
		IsMaker         bool            `json:"isMaker"`
		IsBestMatch     bool            `json:"isBestMatch"`
	}
)

const (
	MyTradesEndPoint = "/api/v3/myTrades"
)

func (rc *RestClient) MyTrades(ctx context.Context, req *exchange.TradeReqParam) ([]Trade, error) {
	var ret []Trade
	value := binance.TradeParam(req)
	if err := rc.Request(ctx, http.MethodGet, MyTradesEndPoint, value, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch myTrades fail")
	}
	return ret, nil
}

func (rc *RestClient) Trades(ctx context.Context, req *exchange.TradeReqParam) ([]*exchange.Trade, error) {
	trades, err := rc.MyTrades(ctx, req)
	if err != nil {
		return nil, err
	}
	ret := []*exchange.Trade{}
	for i := range trades {
		trade := trades[i]
		t, err := trade.Parse()
		if err != nil {
			return nil, err
		}

		ret = append(ret, t)
	}
	return ret, nil
}

func (t *Trade) Parse() (*exchange.Trade, error) {
	s, err := ParseSymbol(t.Symbol)
	if err != nil {
		return nil, err
	}

	var side exchange.OrderSide
	if t.IsBuyer {
		side = exchange.OrderSideBuy
	} else {
		side = exchange.OrderSideSell
	}

	ret := &exchange.Trade{
		ID:          exchange.NewIntID(t.ID),
		OrderID:     exchange.NewIntID(t.OrderID),
		Symbol:      s,
		Amount:      t.Qty,
		Price:       t.Price,
		Fee:         t.Commission,
		FeeCurrency: t.CommissionAsset,
		Time:        binance.ParseTimestamp(t.Time),
		Side:        side,
		Raw:         t,
	}
	return ret, nil
}

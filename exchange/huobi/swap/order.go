package swap

import (
	"context"
	"encoding/json"
	"strings"
)

type (
	OrderReq struct {
		data map[string]interface{}
	}

	OrderResp struct {
		OrderID    int64  `json:"order_id"`
		OrderIDStr string `json:"order_id_str"`
	}

	SwapCancelReq struct {
		symbol         string
		orderIDS       []string
		clientOrderIDS []string
	}

	SwapCancelError struct {
		OrderID string `json:"order_id"`
		ErrCode int    `json:"err_code"`
		ErrMsg  string `json:"err_msg"`
	}

	SwapCancelResp struct {
		Errors    []SwapCancelError `json:"errors"`
		Successes string            `json:"successes"` //id1,id2,id3 ...
	}

	SwapOrderDetailReq struct {
		data map[string]interface{}
	}

	SwapOrderDetailResp struct {
		Symbol          string             `json:"symbol"`
		ContractCode    string             `json:"contract_code"`
		LeverRate       int                `json:"lever_rate"`
		Direction       string             `json:"direction"`
		Offset          string             `json:"offset"`
		Volume          float64            `json:"volume"`
		Price           float64            `json:"price"`
		CreatedAt       int64              `json:"created_at"`
		CanceledAt      int64              `json:"canceled_at"`
		OrderSource     string             `json:"order_source"`
		OrderPriceType  string             `json:"order_price_type"`
		MarginFrozen    float64            `json:"margin_frozen"`
		Profit          float64            `json:"profit"`
		Trades          []OrderNotifyTrade `json:"trades"`
		TotalPage       int                `json:"total_page"`
		CurrentPage     int                `json:"current_page"`
		TotalSize       int                `json:"total_size"`
		LiquidationType string             `json:"liquidation_type"`
		FeeAsset        string             `json:"fee_asset"`
		Fee             float64            `json:"fee"`
		OrderID         int64              `json:"order_id"`
		OrderIDStr      string             `json:"order_id_str"`
		ClientOrderID   interface{}        `json:"client_order_id"`
		OrderType       string             `json:"order_type"`
		Status          int                `json:"status"`
		TradeAvgPrice   float64            `json:"trade_avg_price"`
		TradeTurnOver   float64            `json:"trade_turn_over"`
		TradeVolume     float64            `json:"trade_volume"`
		IsTpsl          interface{}        `json:"is_tpsl"`
		RealProfit      float64            `json:"real_profit"`
	}
)

const (
	SwapOrderEndPoint       = "/swap-api/v1/swap_order"
	SwapCancelEndPoint      = "/swap-api/v1/swap_cancel"
	SwapOrderDetailEndPoint = "/swap-api/v1/swap_order_detail"

	OrderDirectionBuy  = "buy"
	OrderDirectionSell = "sell"
	OrderOffsetOpen    = "open"
	OrderOffsetClose   = "close"

	OrderPriceLimit  = "limit"
	OrderPriceMarket = "opponent"
)

func NewOrderReq(contractCode string, volume int, direction string, offset string, lever int, orderPriceType string) *OrderReq {
	ret := OrderReq{
		data: make(map[string]interface{}),
	}

	ret.data["contract_code"] = contractCode
	ret.data["volume"] = volume
	ret.data["direction"] = direction
	ret.data["offset"] = offset
	ret.data["lever_rate"] = lever
	ret.data["order_price_type"] = orderPriceType
	return &ret
}

func (or *OrderReq) Price(price float64) *OrderReq {
	or.data["price"] = price
	return or
}

func (or *OrderReq) Serialize() ([]byte, error) {
	return json.Marshal(or.data)
}

func NewSwapCancelReq(symbol string) *SwapCancelReq {
	return &SwapCancelReq{
		symbol: symbol,
	}
}

func (scr *SwapCancelReq) Orders(ids ...string) *SwapCancelReq {
	for _, id := range ids {
		scr.orderIDS = append(scr.orderIDS, id)
	}
	return scr
}

func (scr *SwapCancelReq) ClientOrderIDs(ids ...string) *SwapCancelReq {
	for _, id := range ids {
		scr.clientOrderIDS = append(scr.clientOrderIDS, id)
	}
	return scr
}

func (scr *SwapCancelReq) Serialize() ([]byte, error) {
	data := map[string]string{
		"contract_code": scr.symbol,
	}

	if len(scr.orderIDS) != 0 {
		data["order_id"] = strings.Join(scr.orderIDS, ",")
	}

	if len(scr.clientOrderIDS) != 0 {
		data["client_order_id"] = strings.Join(scr.clientOrderIDS, ",")
	}

	return json.Marshal(data)
}

func NewSwapOrderDetailReq(cc string, id int64) *SwapOrderDetailReq {
	return &SwapOrderDetailReq{
		data: map[string]interface{}{
			"contract_code": cc,
			"order_id":      id,
		},
	}
}

func (sdr *SwapOrderDetailReq) CreatedAt(ts int64) *SwapOrderDetailReq {
	sdr.data["created_at"] = ts
	return sdr
}

func (sdr *SwapOrderDetailReq) OrderType(ot int) *SwapOrderDetailReq {
	sdr.data["order_type"] = ot
	return sdr
}

func (sdr *SwapOrderDetailReq) PageIndex(pi int) *SwapOrderDetailReq {
	sdr.data["page_index"] = pi
	return sdr
}

func (sdr *SwapOrderDetailReq) PageSize(ps int) *SwapOrderDetailReq {
	sdr.data["page_size"] = ps
	return sdr
}

func (sdr *SwapOrderDetailReq) Serialize() ([]byte, error) {
	return json.Marshal(sdr.data)
}

func (rc *RestClient) SwapOrder(ctx context.Context, req *OrderReq) (*OrderResp, error) {
	var resp OrderResp
	if err := rc.PrivatePostReq(ctx, SwapOrderEndPoint, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (rc *RestClient) SwapCancel(ctx context.Context, req *SwapCancelReq) (*SwapCancelResp, error) {
	var resp SwapCancelResp
	if err := rc.PrivatePostReq(ctx, SwapCancelEndPoint, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (rc *RestClient) SwapOrderDetail(ctx context.Context, req *SwapOrderDetailReq) (*SwapOrderDetailResp, error) {
	var ret SwapOrderDetailResp
	if err := rc.PrivatePostReq(ctx, SwapOrderDetailEndPoint, req, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

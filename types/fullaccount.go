package types

//go:generate ffjson $GOFILE

import (
	"encoding/json"

	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

type FullAccountInfos []FullAccountInfo

type FullAccountInfo struct {
	ID          AccountID
	AccountInfo AccountInfo
}

type AccountInfo struct {
	Account              Account               `json:"account"`
	RegistrarName        String                `json:"registrar_name"`
	ReferrerName         String                `json:"referrer_name"`
	LifetimeReferrerName String                `json:"lifetime_referrer_name"`
	CashbackBalance      VestingBalance        `json:"cashback_balance"`
	Balances             AccountBalances       `json:"balances"`
	VestingBalances      VestingBalances       `json:"vesting_balances"`
	LimitOrders          LimitOrders           `json:"limit_orders"`
	CallOrders           CallOrders            `json:"call_orders"`
	SettleOrders         ForceSettlementOrders `json:"settle_orders"`
	Statistics           AccountStatistics     `json:"statistics"`
	Assets               AssetIDs              `json:"assets"`
	//Proposals            []interface{}   `json:"proposals"`
	//Withdraws            []interface{}   `json:"withdraws"`
	// Votes                Votes   `json:"votes"`
}

func (p FullAccountInfo) MarshalJSON() ([]byte, error) {
	ret := make([]interface{}, 2)
	ret[0] = p.ID.ID()
	ret[1] = p.AccountInfo

	return ffjson.Marshal(ret)
}

func (p *FullAccountInfo) UnmarshalJSON(data []byte) error {
	raw := make([]json.RawMessage, 2)
	if err := ffjson.Unmarshal(data, &raw); err != nil {
		return errors.Annotate(err, "unmarshal [raw]")
	}
	if err := ffjson.Unmarshal(raw[0], &p.ID); err != nil {
		return errors.Annotate(err, "unmarshal [id]")
	}
	if err := ffjson.Unmarshal(raw[1], &p.AccountInfo); err != nil {
		return errors.Annotate(err, "unmarshal [AccountInfo]")
	}

	return nil
}

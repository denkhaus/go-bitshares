package types

import (
	"encoding/json"

	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

type FullAccountInfos []FullAccountInfo

type FullAccountInfo struct {
	ID   GrapheneID
	Info Info
}

type Info struct {
	Account              Account          `json:"account"`
	RegistrarName        string           `json:"registrar_name"`
	ReferrerName         string           `json:"referrer_name"`
	LifetimeReferrerName string           `json:"lifetime_referrer_name"`
	CashbackBalance      VestingBalance   `json:"cashback_balance"`
	Balances             []AccountBalance `json:"balances"`
	VestingBalances      []VestingBalance `json:"vesting_balances"`
	LimitOrders          LimitOrders      `json:"limit_orders"`
	CallOrders           CallOrders       `json:"call_orders"`
	SettleOrders         SettleOrders     `json:"settle_orders"`
	//Proposals            []interface{}   `json:"proposals"`
	//Assets               []string        `json:"assets"`
	//Withdraws            []interface{}   `json:"withdraws"`
	//    Statistics           Statistics      `json:"statistics"`
	// Votes                []interface{}   `json:"votes"`
}

type AccountBalance struct {
	ID              GrapheneID `json:"id"`
	Owner           GrapheneID `json:"owner"`
	AssetType       GrapheneID `json:"asset_type"`
	Balance         UInt64     `json:"balance"`
	MaintenanceFlag bool       `json:"maintenance_flag"`
}

type VestingBalance struct {
	ID      GrapheneID             `json:"id"`
	Balance AssetAmount            `json:"balance"`
	Owner   Address                `json:"owner"`
	Policy  []VestingPolicyBalance `json:"policy"`
}

type VestingPolicyBalance struct {
	StartClaim                  Time   `json:"start_claim"`
	VestingSeconds              UInt64 `json:"vesting_seconds"`
	CoinSecondsEarned           UInt64 `json:"coin_seconds_earned"`
	CoinSecondsEarnedLastUpdate Time   `json:"coin_seconds_earned_last_update"`
}

func (p *FullAccountInfo) MarshalJSON() ([]byte, error) {
	data := make([]interface{}, 2)
	data[0] = p.ID.ID()
	data[1] = p.Info

	return ffjson.Marshal(data)
}

func (p *FullAccountInfo) UnmarshalJSON(data []byte) error {
	raw := make([]json.RawMessage, 2)
	if err := ffjson.Unmarshal(data, &raw); err != nil {
		return errors.Annotate(err, "unmarshal Info [unmarshal]")
	}

	if err := ffjson.Unmarshal(raw[0], &p.ID); err != nil {
		return errors.Annotate(err, "unmarshal Info [id]")
	}
	if err := ffjson.Unmarshal(raw[1], &p.Info); err != nil {
		return errors.Annotate(err, "unmarshal Info [FullAccountInfo]")
	}

	return nil
}

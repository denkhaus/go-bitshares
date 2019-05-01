package types

//go:generate ffjson $GOFILE

type Accounts []Account

func (p Accounts) Lookup(ID GrapheneObject) *Account {
	for _, acct := range p {
		if acct.ID.Equals(ID) {
			return &acct
		}
	}

	return nil
}

type Account struct {
	ID                            AccountID              `json:"id"`
	Name                          String                 `json:"name"`
	Statistics                    ObjectID               `json:"statistics"`
	MembershipExpirationDate      Time                   `json:"membership_expiration_date"`
	NetworkFeePercentage          UInt64                 `json:"network_fee_percentage"`
	LifetimeReferrerFeePercentage UInt64                 `json:"lifetime_referrer_fee_percentage"`
	ReferrerRewardsPercentage     UInt64                 `json:"referrer_rewards_percentage"`
	TopNControlFlags              UInt64                 `json:"top_n_control_flags"`
	WhitelistingAccounts          AccountIDs             `json:"whitelisting_accounts"`
	BlacklistingAccounts          AccountIDs             `json:"blacklisting_accounts"`
	WhitelistedAccounts           AccountIDs             `json:"whitelisted_accounts"`
	BlacklistedAccounts           AccountIDs             `json:"blacklisted_accounts"`
	Options                       AccountOptions         `json:"options"`
	Registrar                     AccountID              `json:"registrar"`
	Referrer                      AccountID              `json:"referrer"`
	LifetimeReferrer              AccountID              `json:"lifetime_referrer"`
	CashbackVB                    VestingBalanceID       `json:"cashback_vb"`
	Owner                         Authority              `json:"owner"`
	Active                        Authority              `json:"active"`
	OwnerSpecialAuthority         OwnerSpecialAuthority  `json:"owner_special_authority"`
	ActiveSpecialAuthority        ActiveSpecialAuthority `json:"active_special_authority"`
}

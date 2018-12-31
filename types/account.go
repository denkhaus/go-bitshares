package types

//go:generate ffjson $GOFILE

type Accounts []Account
type Account struct {
	ID                            GrapheneID             `json:"id"`
	Name                          String                 `json:"name"`
	Statistics                    GrapheneID             `json:"statistics"`
	MembershipExpirationDate      Time                   `json:"membership_expiration_date"`
	NetworkFeePercentage          UInt64                 `json:"network_fee_percentage"`
	LifetimeReferrerFeePercentage UInt64                 `json:"lifetime_referrer_fee_percentage"`
	ReferrerRewardsPercentage     UInt64                 `json:"referrer_rewards_percentage"`
	TopNControlFlags              UInt64                 `json:"top_n_control_flags"`
	WhitelistingAccounts          GrapheneIDs            `json:"whitelisting_accounts"`
	BlacklistingAccounts          GrapheneIDs            `json:"blacklisting_accounts"`
	WhitelistedAccounts           GrapheneIDs            `json:"whitelisted_accounts"`
	BlacklistedAccounts           GrapheneIDs            `json:"blacklisted_accounts"`
	Options                       AccountOptions         `json:"options"`
	Registrar                     GrapheneID             `json:"registrar"`
	Referrer                      GrapheneID             `json:"referrer"`
	LifetimeReferrer              GrapheneID             `json:"lifetime_referrer"`
	CashbackVB                    GrapheneID             `json:"cashback_vb"`
	Owner                         Authority              `json:"owner"`
	Active                        Authority              `json:"active"`
	OwnerSpecialAuthority         OwnerSpecialAuthority  `json:"owner_special_authority"`
	ActiveSpecialAuthority        ActiveSpecialAuthority `json:"active_special_authority"`
}

//NewAccount creates a new Account object
func NewAccount(id GrapheneID) *Account {
	acc := Account{
		ID: id,
	}
	return &acc
}

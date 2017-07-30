package objects

//easyjson:json
type UserAccount struct {
	GrapheneID
	Name                          string         `json:"name"`
	Statistics                    string         `json:"statistics"`
	MembershipExpirationDate      RFC3339Time    `json:"membership_expiration_date"`
	NetworkFeePercentage          int64          `json:"network_fee_percentage"`
	LifetimeReferrerFeePercentage int64          `json:"lifetime_referrer_fee_percentage"`
	ReferrerRewardsPercentage     int64          `json:"referrer_rewards_percentage"`
	TopNControlFlags              int64          `json:"top_n_control_flags"`
	WhitelistingAccounts          []string       `json:"whitelisting_accounts"`
	BlacklistingAccounts          []string       `json:"blacklisting_accounts"`
	WhitelistedAccounts           []string       `json:"whitelisted_accounts"`
	BlacklistedAccounts           []string       `json:"blacklisted_accounts"`
	Options                       AccountOptions `json:"options"`
	// Registrar                     GrapheneID `json:"registrar"`
	// Referrer                      GrapheneID `json:"referrer"`
	// LifetimeReferrer              GrapheneID `json:"lifetime_referrer"`
	Owner                  Authority     `json:"owner"`
	Active                 Authority     `json:"active"`
	OwnerSpecialAuthority  []interface{} `json:"owner_special_authority"`
	ActiveSpecialAuthority []interface{} `json:"active_special_authority"`
}

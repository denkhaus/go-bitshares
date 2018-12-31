package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeCommitteeMemberUpdateGlobalParameters] = func() types.Operation {
		op := &CommitteeMemberUpdateGlobalParametersOperation{}
		return op
	}
}

type ChainParameters struct {
	AllowNonMemberWhitelists         bool              `json:"allow_non_member_whitelists"`
	CountNonMemberVotes              bool              `json:"count_non_member_votes"`
	Extensions                       types.Extensions  `json:"extensions"`
	CurrentFees                      types.FeeSchedule `json:"current_fees"`
	AccountFeeScaleBitshifts         types.UInt8       `json:"account_fee_scale_bitshifts"`
	BlockInterval                    types.UInt8       `json:"block_interval"`
	MaintenanceSkipSlots             types.UInt8       `json:"maintenance_skip_slots"`
	MaxAuthorityDepth                types.UInt8       `json:"max_authority_depth"`
	MaximumAssetFeedPublishers       types.UInt8       `json:"maximum_asset_feed_publishers"`
	MaximumAssetWhitelistAuthorities types.UInt8       `json:"maximum_asset_whitelist_authorities"`
	AccountsPerFeeScale              types.UInt16      `json:"accounts_per_fee_scale"`
	LifetimeReferrerPercentOfFee     types.UInt16      `json:"lifetime_referrer_percent_of_fee"`
	MaxPredicateOpcode               types.UInt16      `json:"max_predicate_opcode"`
	MaximumAuthorityMembership       types.UInt16      `json:"maximum_authority_membership"`
	MaximumCommitteeCount            types.UInt16      `json:"maximum_committee_count"`
	MaximumWitnessCount              types.UInt16      `json:"maximum_witness_count"`
	NetworkPercentOfFee              types.UInt16      `json:"network_percent_of_fee"`
	ReservePercentOfFee              types.UInt16      `json:"reserve_percent_of_fee"`
	CashbackVestingPeriodSeconds     types.UInt32      `json:"cashback_vesting_period_seconds"`
	CommitteeProposalReviewPeriod    types.UInt32      `json:"committee_proposal_review_period"`
	WitnessPayVestingSeconds         types.UInt32      `json:"witness_pay_vesting_seconds"`
	MaximumProposalLifetime          types.UInt32      `json:"maximum_proposal_lifetime"`
	MaximumTimeUntilExpiration       types.UInt32      `json:"maximum_time_until_expiration"`
	MaximumTransactionSize           types.UInt32      `json:"maximum_transaction_size"`
	MaintenanceInterval              types.UInt32      `json:"maintenance_interval"`
	MaximumBlockSize                 types.UInt32      `json:"maximum_block_size"`
	CashbackVestingThreshold         types.Int64       `json:"cashback_vesting_threshold"`
	WitnessPayPerBlock               types.Int64       `json:"witness_pay_per_block"`
	WorkerBudgetPerDay               types.Int64       `json:"worker_budget_per_day"`
	FeeLiquidationThreshold          types.Int64       `json:"fee_liquidation_threshold"`
}

func (p ChainParameters) Marshal(enc *util.TypeEncoder) error {
	// (current_fees)
	if err := enc.Encode(p.CurrentFees); err != nil {
		return errors.Annotate(err, "encode CurrentFees")
	}
	// (block_interval)
	if err := enc.Encode(p.BlockInterval); err != nil {
		return errors.Annotate(err, "encode BlockInterval")
	}
	// (maintenance_interval)
	if err := enc.Encode(p.MaintenanceInterval); err != nil {
		return errors.Annotate(err, "encode MaintenanceInterval")
	}
	// (maintenance_skip_slots)
	if err := enc.Encode(p.MaintenanceSkipSlots); err != nil {
		return errors.Annotate(err, "encode MaintenanceSkipSlots")
	}
	// (committee_proposal_review_period)
	if err := enc.Encode(p.CommitteeProposalReviewPeriod); err != nil {
		return errors.Annotate(err, "encode CommitteeProposalReviewPeriod")
	}
	// (maximum_transaction_size)
	if err := enc.Encode(p.MaximumTransactionSize); err != nil {
		return errors.Annotate(err, "encode MaximumTransactionSize")
	}
	// (maximum_block_size)
	if err := enc.Encode(p.MaximumBlockSize); err != nil {
		return errors.Annotate(err, "encode MaximumBlockSize")
	}
	// (maximum_time_until_expiration)
	if err := enc.Encode(p.MaximumTimeUntilExpiration); err != nil {
		return errors.Annotate(err, "encode MaximumTimeUntilExpiration")
	}
	// (maximum_proposal_lifetime)
	if err := enc.Encode(p.MaximumProposalLifetime); err != nil {
		return errors.Annotate(err, "encode MaximumProposalLifetime")
	}
	// (maximum_asset_whitelist_authorities)
	if err := enc.Encode(p.MaximumAssetWhitelistAuthorities); err != nil {
		return errors.Annotate(err, "encode MaximumAssetWhitelistAuthorities")
	}
	// (maximum_asset_feed_publishers)
	if err := enc.Encode(p.MaximumAssetFeedPublishers); err != nil {
		return errors.Annotate(err, "encode MaximumAssetFeedPublishers")
	}
	// (maximum_witness_count)
	if err := enc.Encode(p.MaximumWitnessCount); err != nil {
		return errors.Annotate(err, "encode MaximumWitnessCount")
	}
	// (maximum_committee_count)
	if err := enc.Encode(p.MaximumCommitteeCount); err != nil {
		return errors.Annotate(err, "encode MaximumCommitteeCount")
	}
	// (maximum_authority_membership)
	if err := enc.Encode(p.MaximumAuthorityMembership); err != nil {
		return errors.Annotate(err, "encode MaximumAuthorityMembership")
	}
	// (reserve_percent_of_fee)
	if err := enc.Encode(p.ReservePercentOfFee); err != nil {
		return errors.Annotate(err, "encode ReservePercentOfFee")
	}
	// (network_percent_of_fee)
	if err := enc.Encode(p.NetworkPercentOfFee); err != nil {
		return errors.Annotate(err, "encode NetworkPercentOfFee")
	}
	// (lifetime_referrer_percent_of_fee)
	if err := enc.Encode(p.LifetimeReferrerPercentOfFee); err != nil {
		return errors.Annotate(err, "encode LifetimeReferrerPercentOfFee")
	}
	// (cashback_vesting_period_seconds)
	if err := enc.Encode(p.CashbackVestingPeriodSeconds); err != nil {
		return errors.Annotate(err, "encode CashbackVestingPeriodSeconds")
	}
	// (cashback_vesting_threshold)
	if err := enc.Encode(p.CashbackVestingThreshold); err != nil {
		return errors.Annotate(err, "encode CashbackVestingThreshold")
	}
	// (count_non_member_votes)
	if err := enc.Encode(p.CountNonMemberVotes); err != nil {
		return errors.Annotate(err, "encode CountNonMemberVotes")
	}
	// (allow_non_member_whitelists)
	if err := enc.Encode(p.AllowNonMemberWhitelists); err != nil {
		return errors.Annotate(err, "encode AllowNonMemberWhitelists")
	}
	// (witness_pay_per_block)
	if err := enc.Encode(p.WitnessPayPerBlock); err != nil {
		return errors.Annotate(err, "encode WitnessPayPerBlock")
	}
	// (witness_pay_vesting_seconds)
	// if err := enc.Encode(p.WitnessPayVestingSeconds); err != nil {
	// 	return errors.Annotate(err, "encode WitnessPayVWestingSeconds")
	// }
	// (worker_budget_per_day)
	if err := enc.Encode(p.WorkerBudgetPerDay); err != nil {
		return errors.Annotate(err, "encode WorkerBudgetPerDay")
	}
	// (max_predicate_opcode)
	if err := enc.Encode(p.MaxPredicateOpcode); err != nil {
		return errors.Annotate(err, "encode MaxPredicateOpcode")
	}
	// (fee_liquidation_threshold)
	if err := enc.Encode(p.FeeLiquidationThreshold); err != nil {
		return errors.Annotate(err, "encode FeeLiquidationThreshold")
	}
	// (accounts_per_fee_scale)
	if err := enc.Encode(p.AccountsPerFeeScale); err != nil {
		return errors.Annotate(err, "encode AccountsPerFeeScale")
	}
	// (account_fee_scale_bitshifts)
	if err := enc.Encode(p.AccountFeeScaleBitshifts); err != nil {
		return errors.Annotate(err, "encode AccountFeeScaleBitshifts")
	}
	// (max_authority_depth)
	if err := enc.Encode(p.MaxAuthorityDepth); err != nil {
		return errors.Annotate(err, "encode MaxAuthorityDepth")
	}
	// (extensions)
	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode Extensions")
	}

	return nil
}

type CommitteeMemberUpdateGlobalParametersOperation struct {
	types.OperationFee
	NewParameters ChainParameters `json:"new_parameters"`
}

func (p CommitteeMemberUpdateGlobalParametersOperation) Type() types.OperationType {
	return types.OperationTypeCommitteeMemberUpdateGlobalParameters
}

func (p CommitteeMemberUpdateGlobalParametersOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}

	if err := enc.Encode(p.NewParameters); err != nil {
		return errors.Annotate(err, "encode NewParameters")
	}

	return nil
}

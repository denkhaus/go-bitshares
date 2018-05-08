package objects

import (
	"encoding/json"

	"github.com/pquerna/ffjson/ffjson"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type Operation interface {
	util.TypeMarshaller
	ApplyFee(fee AssetAmount)
	Type() OperationType
}

type OperationResult interface {
}

type OperationEnvelope struct {
	Type      OperationType
	Operation Operation
}

func (p OperationEnvelope) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		p.Type,
		p.Operation,
	})
}

func (p *OperationEnvelope) UnmarshalJSON(data []byte) error {
	raw := make([]json.RawMessage, 2)
	if err := json.Unmarshal(data, &raw); err != nil {
		return errors.Annotate(err, "Unmarshal raw object")
	}

	if len(raw) != 2 {
		return errors.Errorf("Invalid operation data: %v", string(data))
	}

	if err := json.Unmarshal(raw[0], &p.Type); err != nil {
		return errors.Annotate(err, "Unmarshal OperationType")
	}

	switch p.Type {
	case OperationTypeLimitOrderCreate:
		p.Operation = &LimitOrderCreateOperation{}
		if err := ffjson.Unmarshal(raw[1], p.Operation); err != nil {
			return errors.Annotate(err, "unmarshal LimitOrderCreateOperation")
		}

	case OperationTypeTransfer:
		p.Operation = &TransferOperation{}
		if err := ffjson.Unmarshal(raw[1], p.Operation); err != nil {
			return errors.Annotate(err, "unmarshal TransferOperation")
		}

	case OperationTypeLimitOrderCancel:
		p.Operation = &LimitOrderCancelOperation{}
		if err := ffjson.Unmarshal(raw[1], p.Operation); err != nil {
			return errors.Annotate(err, "unmarshal LimitOrderCancelOperation")
		}

	case OperationTypeCallOrderUpdate:
		p.Operation = &CallOrderUpdateOperation{}
		if err := ffjson.Unmarshal(raw[1], p.Operation); err != nil {
			return errors.Annotate(err, "unmarshal CallOrderUpdateOperation")
		}

	case OperationTypeAccountCreate:
		p.Operation = &AccountCreateOperation{}
		if err := ffjson.Unmarshal(raw[1], p.Operation); err != nil {
			return errors.Annotate(err, "unmarshal AccountCreateOperation")
		}

	case OperationTypeAccountUpdate:
		util.DumpUnmarshaled("OperationTypeAccountUpdate", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeAccountWhitelist:
		util.DumpUnmarshaled("OperationTypeAccountWhitelist", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeAccountUpgrade:
		util.DumpUnmarshaled("OperationTypeAccountUpgrade", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeAccountTransfer:
		util.DumpUnmarshaled("OperationTypeAccountTransfer", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeAssetCreate:
		util.DumpUnmarshaled("OperationTypeAssetCreate", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeFillOrder:
		p.Operation = &FillOrderOperation{}
		if err := ffjson.Unmarshal(raw[1], p.Operation); err != nil {
			return errors.Annotate(err, "unmarshal FillOrderOperation")
		}

	case OperationTypeAssetUpdate:
		p.Operation = &AssetUpdateOperation{}
		if err := ffjson.Unmarshal(raw[1], p.Operation); err != nil {
			return errors.Annotate(err, "unmarshal AssetUpdateOperation")
		}

	case OperationTypeAssetIssue:
		p.Operation = &AssetIssueOperation{}
		if err := ffjson.Unmarshal(raw[1], p.Operation); err != nil {
			return errors.Annotate(err, "unmarshal AssetIssueOperation")
		}

	case OperationTypeAssetPublishFeed:
		p.Operation = &AssetPublishFeedOperation{}
		if err := ffjson.Unmarshal(raw[1], p.Operation); err != nil {
			return errors.Annotate(err, "unmarshal AssetPublishFeedOperation")
		}

	case OperationTypeAssetUpdateBitasset:
		util.DumpUnmarshaled("OperationTypeAssetUpdateBitasset", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeAssetUpdateFeedProducers:
		util.DumpUnmarshaled("OperationTypeAssetUpdateFeedProducers", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeAssetReverse:
		util.DumpUnmarshaled("OperationTypeAssetReverse", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeAssetFundFeePool:
		util.DumpUnmarshaled("OperationTypeAssetFundFeePool", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeAssetSettle:
		util.DumpUnmarshaled("OperationTypeAssetSettle", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeAssetGlobalSettle:
		util.DumpUnmarshaled("OperationTypeAssetGlobalSettle", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeWitnessCreate:
		util.DumpUnmarshaled("OperationTypeWitnessCreate", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeWitnessUpdate:
		util.DumpUnmarshaled("OperationTypeWitnessUpdate", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeProposalCreate:
		util.DumpUnmarshaled("OperationTypeProposalCreate", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeProposalUpdate:
		util.DumpUnmarshaled("OperationTypeProposalUpdate", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeProposalDelete:
		util.DumpUnmarshaled("OperationTypeProposalDelete", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeWithdrawPermissionCreate:
		util.DumpUnmarshaled("OperationTypeWithdrawPermissionCreate", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeWithdrawPermissionUpdate:
		util.DumpUnmarshaled("OperationTypeWithdrawPermissionUpdate", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeWithdrawPermissionClaim:
		util.DumpUnmarshaled("OperationTypeWithdrawPermissionClaim", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeWithdrawPermissionDelete:
		util.DumpUnmarshaled("OperationTypeWithdrawPermissionDelete", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeCommiteeMemberCreate:
		util.DumpUnmarshaled("OperationTypeCommiteeMemberCreate", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeCommiteeMemberUpdate:
		util.DumpUnmarshaled("OperationTypeCommiteeMemberUpdate", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeCommiteeMemberUpdateGlobalParameters:
		util.DumpUnmarshaled("OperationTypeCommiteeMemberUpdateGlobalParameters", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeVestingBalanceCreate:
		util.DumpUnmarshaled("OperationTypeVestingBalanceCreate", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeVestingBalanceWithdraw:
		util.DumpUnmarshaled("OperationTypeVestingBalanceWithdraw", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeWorkerCreate:
		util.DumpUnmarshaled("OperationTypeWorkerCreate", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeCustom:
		util.DumpUnmarshaled("OperationTypeCustom", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeAssert:
		util.DumpUnmarshaled("OperationTypeAssert", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeBalanceClaim:
		util.DumpUnmarshaled("OperationTypeBalanceClaim", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeOverrideTransfer:
		util.DumpUnmarshaled("OperationTypeOverrideTransfer", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeTransferToBlind:
		util.DumpUnmarshaled("OperationTypeTransferToBlind", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeBlindTransfer:
		util.DumpUnmarshaled("OperationTypeBlindTransfer", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeTransferFromBlind:
		util.DumpUnmarshaled("OperationTypeTransferFromBlind", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeAssetSettleCancel:
		util.DumpUnmarshaled("OperationTypeAssetSettleCancel", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	case OperationTypeAssetClaimFees:
		util.DumpUnmarshaled("OperationTypeAssetClaimFees", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)

	default:
		return errors.Errorf("Operation type %d not yet supported", p.Type)
	}

	return nil
}

type Operations []Operation

//implements TypeMarshaller interface
func (p Operations) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return errors.Annotate(err, "encode Operations length")
	}

	for _, op := range p {
		if err := enc.Encode(op); err != nil {
			return errors.Annotate(err, "encode Operation")
		}
	}

	return nil
}

func (p Operations) MarshalJSON() ([]byte, error) {
	env := make([]OperationEnvelope, len(p))
	for idx, op := range p {
		env[idx] = OperationEnvelope{
			Type:      op.Type(),
			Operation: op,
		}
	}

	return json.Marshal(env)
}

func (p *Operations) UnmarshalJSON(data []byte) error {
	var envs []OperationEnvelope
	if err := json.Unmarshal(data, &envs); err != nil {
		return err
	}

	ops := make(Operations, len(envs))
	for idx, env := range envs {
		ops[idx] = env.Operation.(Operation)
	}

	*p = ops
	return nil
}

func (p Operations) ApplyFees(fees []AssetAmount) error {
	if len(p) != len(fees) {
		return errors.New("count of fees must match count of operations")
	}

	for idx, op := range p {
		op.ApplyFee(fees[idx])
	}

	return nil
}

func (p Operations) Types() [][]OperationType {
	ret := make([][]OperationType, len(p))
	for idx, op := range p {
		ret[idx] = []OperationType{op.Type()}
	}

	return ret
}

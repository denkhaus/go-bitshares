package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeAssetIssue] = func() types.Operation {
		op := &AssetIssueOperation{}
		return op
	}
}

type AssetIssueOperation struct {
	types.OperationFee
	Issuer         types.GrapheneID  `json:"issuer"`
	IssueToAccount types.GrapheneID  `json:"issue_to_account"`
	AssetToIssue   types.AssetAmount `json:"asset_to_issue"`
	Memo           *types.Memo       `json:"memo"`
	Extensions     types.Extensions  `json:"extensions"`
}

func (p AssetIssueOperation) Type() types.OperationType {
	return types.OperationTypeAssetIssue
}

func (p AssetIssueOperation) MarshalFeeScheduleParams(params types.M, enc *util.TypeEncoder) error {
	if fee, ok := params["fee"]; ok {
		if err := enc.Encode(types.UInt64(fee.(float64))); err != nil {
			return errors.Annotate(err, "encode Fee")
		}
	}

	if ppk, ok := params["price_per_kbyte"]; ok {
		if err := enc.Encode(types.UInt32(ppk.(float64))); err != nil {
			return errors.Annotate(err, "encode PricePerKByte")
		}
	}

	return nil
}

func (p AssetIssueOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.Issuer); err != nil {
		return errors.Annotate(err, "encode issuer")
	}

	if err := enc.Encode(p.AssetToIssue); err != nil {
		return errors.Annotate(err, "encode asset to issue")
	}

	if err := enc.Encode(p.IssueToAccount); err != nil {
		return errors.Annotate(err, "encode issue to account")
	}

	if err := enc.Encode(p.Memo != nil); err != nil {
		return errors.Annotate(err, "encode have memo")
	}

	if err := enc.Encode(p.Memo); err != nil {
		return errors.Annotate(err, "encode memo")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}

package operations

//go:generate ffjson   $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &AssetIssueOperation{}
	types.OperationMap[op.Type()] = op
}

type AssetIssueOperation struct {
	Issuer         types.GrapheneID  `json:"issuer"`
	IssueToAccount types.GrapheneID  `json:"issue_to_account"`
	AssetToIssue   types.AssetAmount `json:"asset_to_issue"`
	Fee            types.AssetAmount `json:"fee"`
	Memo           *types.Memo       `json:"memo"`
	Extensions     types.Extensions  `json:"extensions"`
}

func (p *AssetIssueOperation) ApplyFee(fee types.AssetAmount) {
	p.Fee = fee
}

func (p AssetIssueOperation) Type() types.OperationType {
	return types.OperationTypeAssetIssue
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

//NewAssetIssueOperation creates a new AssetIssueOperation
func NewAssetIssueOperation() *AssetIssueOperation {
	tx := AssetIssueOperation{
		Extensions: types.Extensions{},
	}
	return &tx
}

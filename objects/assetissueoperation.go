package objects

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type AssetIssueOperation struct {
	Issuer         GrapheneID  `json:"issuer"`
	IssueToAccount GrapheneID  `json:"issue_to_account"`
	AssetToIssue   AssetAmount `json:"asset_to_issue"`
	Fee            AssetAmount `json:"fee"`
	Extensions     Extensions  `json:"extensions"`
}

//implements Operation interface
func (p *AssetIssueOperation) ApplyFee(fee AssetAmount) {
	p.Fee = fee
}

//implements Operation interface
func (p AssetIssueOperation) Type() OperationType {
	return OperationTypeAssetIssue
}

//TODO: validate encode order!
//implements Operation interface
func (p AssetIssueOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode operation id")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.Issuer); err != nil {
		return errors.Annotate(err, "encode issuer")
	}

	if err := enc.Encode(p.IssueToAccount); err != nil {
		return errors.Annotate(err, "encode issue to account")
	}

	if err := enc.Encode(p.AssetToIssue); err != nil {
		return errors.Annotate(err, "encode asset to issue")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}

//NewAssetIssueOperation creates a new AssetIssueOperation
func NewAssetIssueOperation() *AssetIssueOperation {
	tx := AssetIssueOperation{
		//	Extensions: Extensions{},
	}
	return &tx
}

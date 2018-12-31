package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/denkhaus/logging"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeAssert] = func() types.Operation {
		op := &AssertOperation{}
		return op
	}
}

// struct account_name_eq_lit_predicate
// {
//    account_id_type account_id;
//    string          name;

//    bool validate()const;
// };

// struct asset_symbol_eq_lit_predicate
// {
//    asset_id_type   asset_id;
//    string          symbol;

//    bool validate()const;

// };

// struct block_id_predicate
// {
//    block_id_type id;
//    bool validate()const{ return true; }
// };

// typedef static_variant<
//    account_name_eq_lit_predicate,
//    asset_symbol_eq_lit_predicate,
//    block_id_predicate
// > predicate;

type AssertOperation struct {
	types.OperationFee
	// asset                      fee;
	// account_id_type            fee_paying_account;
	// vector<predicate>          predicates;
	// flat_set<account_id_type>  required_auths;
	// extensions_type extensions;
}

func (p AssertOperation) Type() types.OperationType {
	return types.OperationTypeAssert
}

func (p AssertOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}
	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}
	//(fee)(fee_paying_account)(predicates)(required_auths)(extensions)
	logging.Warnf("%s is not implemented", p.Type().OperationName())
	return nil
}

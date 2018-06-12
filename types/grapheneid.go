package types

import (
	"strconv"
	"strings"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

type GrapheneObject interface {
	util.TypeMarshaller
	ID() string
	String() string
	Type() ObjectType
	Space() SpaceType
	Equals(id GrapheneObject) bool
	Valid() bool
}

type GrapheneObjects []GrapheneObject

func (p GrapheneObjects) ToStrings() []string {
	ids := make([]string, len(p))
	for idx, o := range p {
		ids[idx] = o.ID()
	}

	return ids
}

type GrapheneIDs []GrapheneID

func (p GrapheneIDs) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return errors.Annotate(err, "encode length")
	}

	for _, ex := range p {
		if err := enc.Encode(ex); err != nil {
			return errors.Annotate(err, "encode GrapheneID")
		}
	}

	return nil
}

type GrapheneID struct {
	id         string
	spaceType  SpaceType
	objectType ObjectType
	instance   UInt64
}

func (p GrapheneID) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(p.instance)); err != nil {
		return errors.Annotate(err, "encode instance")
	}

	return nil
}

func (p GrapheneID) MarshalJSON() ([]byte, error) {
	return ffjson.Marshal(p.id)
}

func (p *GrapheneID) UnmarshalJSON(s []byte) error {
	var val string
	if err := ffjson.Unmarshal(s, &val); err != nil {
		return errors.Annotate(err, "Unmarshal")
	}

	if err := p.FromString(val); err != nil {
		return errors.Annotate(err, "FromString")
	}

	return nil
}

func (p GrapheneID) Equals(o GrapheneObject) bool {
	return p.id == o.ID()
}

//String, GrapheneID implements Stringer
func (p GrapheneID) String() string {
	return p.id
}

//ID returns the objects ID
func (p GrapheneID) ID() string {
	return p.id
}

//Type returns the objects ObjectType
func (p GrapheneID) Type() ObjectType {
	if !p.Valid() {
		if err := p.FromString(p.id); err != nil {
			panic(errors.Annotate(err, "from string").Error())
		}
	}
	return p.objectType
}

//Space returns the objects SpaceType
func (p GrapheneID) Space() SpaceType {
	if !p.Valid() {
		if err := p.FromString(p.id); err != nil {
			panic(errors.Annotate(err, "from string").Error())
		}
	}
	return p.spaceType
}

//NewGrapheneID creates an new GrapheneID object
func NewGrapheneID(id string) *GrapheneID {
	gid := &GrapheneID{
		spaceType:  SpaceTypeUndefined,
		objectType: ObjectTypeUndefined,
	}

	if err := gid.FromString(id); err != nil {
		panic(err.Error())
	}

	return gid
}

func (p GrapheneID) Valid() bool {
	return p.id != "" &&
		p.spaceType != SpaceTypeUndefined &&
		p.objectType != ObjectTypeUndefined
}

func (p *GrapheneID) FromRawData(in interface{}) error {
	o, ok := in.(map[string]interface{})
	if !ok {
		return errors.New("input is not map[string]interface{}")
	}

	if id, ok := o["id"]; ok {
		return p.FromString(id.(string))
	}

	return errors.New("input is no graphene object")
}

func (p *GrapheneID) FromString(in string) error {
	if len(in) == 0 {
		return nil
	}

	parts := strings.Split(in, ".")

	if len(parts) == 3 {
		p.id = in
		space, err := strconv.Atoi(parts[0])
		if err != nil {
			return errors.Errorf("unable to parse GrapheneID [space] from %s", in)
		}

		p.spaceType = SpaceType(space)

		typ, err := strconv.Atoi(parts[1])
		if err != nil {
			return errors.Errorf("unable to parse GrapheneID [type] from %s", in)
		}

		switch p.spaceType {
		case SpaceTypeProtocol:
			switch ObjectType(typ) {
			case ObjectTypeBase:
				p.objectType = ObjectTypeBase
			case ObjectTypeAccount:
				p.objectType = ObjectTypeAccount
			case ObjectTypeAsset:
				p.objectType = ObjectTypeAsset
			case ObjectTypeForceSettlement:
				p.objectType = ObjectTypeForceSettlement
			case ObjectTypeCommiteeMember:
				p.objectType = ObjectTypeCommiteeMember
			case ObjectTypeWitness:
				p.objectType = ObjectTypeWitness
			case ObjectTypeLimitOrder:
				p.objectType = ObjectTypeLimitOrder
			case ObjectTypeCallOrder:
				p.objectType = ObjectTypeCallOrder
			case ObjectTypeCustom:
				p.objectType = ObjectTypeCustom
			case ObjectTypeProposal:
				p.objectType = ObjectTypeProposal
			case ObjectTypeOperationHistory:
				p.objectType = ObjectTypeOperationHistory
			case ObjectTypeWithdrawPermission:
				p.objectType = ObjectTypeWithdrawPermission
			case ObjectTypeVestingBalance:
				p.objectType = ObjectTypeVestingBalance
			case ObjectTypeWorker:
				p.objectType = ObjectTypeWorker
			case ObjectTypeBalance:
				p.objectType = ObjectTypeBalance
			}

		case SpaceTypeImplementation:
			switch ObjectType(typ) {
			case ObjectTypeGlobalProperty:
				p.objectType = ObjectTypeGlobalProperty
			case ObjectTypeDynamicGlobalProperty:
				p.objectType = ObjectTypeDynamicGlobalProperty
			case ObjectTypeAssetDynamicData:
				p.objectType = ObjectTypeAssetDynamicData
			case ObjectTypeAssetBitAssetData:
				p.objectType = ObjectTypeAssetBitAssetData
			case ObjectTypeAccountBalance:
				p.objectType = ObjectTypeAccountBalance
			case ObjectTypeAccountStatistics:
				p.objectType = ObjectTypeAccountStatistics
			case ObjectTypeTransaction:
				p.objectType = ObjectTypeTransaction
			case ObjectTypeBlockSummary:
				p.objectType = ObjectTypeBlockSummary
			case ObjectTypeAccountTransactionHistory:
				p.objectType = ObjectTypeAccountTransactionHistory
			case ObjectTypeBlindedBalance:
				p.objectType = ObjectTypeBlindedBalance
			case ObjectTypeChainProperty:
				p.objectType = ObjectTypeChainProperty
			case ObjectTypeWitnessSchedule:
				p.objectType = ObjectTypeWitnessSchedule
			case ObjectTypeBudgetRecord:
				p.objectType = ObjectTypeBudgetRecord
			case ObjectTypeSpecialAuthority:
				p.objectType = ObjectTypeSpecialAuthority
			}
		}

		inst, err := strconv.ParseUint(parts[2], 10, 64)
		if err != nil {
			return errors.Errorf("unable to parse GrapheneID [instance] from %s", in)
		}

		p.instance = UInt64(inst)
		return nil
	}

	return errors.Errorf("unable to parse GrapheneID from %s", in)
}

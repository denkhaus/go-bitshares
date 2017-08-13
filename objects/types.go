package objects

import (
	"strconv"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type SpaceType int
type ObjectID string
type AssetType int

const (
	AssetTypeUndefined AssetType = -1
	AssetTypeCoreAsset AssetType = iota
	AssetTypeUIA
	AssetTypeSmartCoin
	AssetTypePredictionMarket
)

type ObjectType int

const (
	SpaceTypeUndefined SpaceType = -1
	SpaceTypeProtocol  SpaceType = iota
	SpaceTypeImplementation
)

const (
	ObjectTypeUndefined ObjectType = -1
	ObjectTypeBase      ObjectType = iota
	ObjectTypeAccount
	ObjectTypeAsset
	ObjectTypeFORCE_SETTLEMENT_OBJECT
	ObjectTypeCOMMITTEE_MEMBER_OBJECT
	ObjectTypeWITNESS_OBJECT
	ObjectTypeLIMIT_ORDER_OBJECT
	ObjectTypeCALL_ORDER_OBJECT
	ObjectTypeCUSTOM_OBJECT
	ObjectTypePROPOSAL_OBJECT
	ObjectTypeOPERATION_HISTORY_OBJECT
	ObjectTypeWITHDRAW_PERMISSION_OBJECT
	ObjectTypeVESTING_BALANCE_OBJECT
	ObjectTypeWORKER_OBJECT
	ObjectTypeBALANCE_OBJECT
	ObjectTypeGLOBAL_PROPERTY_OBJECT
	ObjectTypeDYNAMIC_GLOBAL_PROPERTY_OBJECT
	ObjectTypeASSET_DYNAMIC_DATA
	ObjectTypeAssetBitAssetData
	ObjectTypeACCOUNT_BALANCE_OBJECT
	ObjectTypeACCOUNT_STATISTICS_OBJECT
	ObjectTypeTRANSACTION_OBJECT
	ObjectTypeBLOCK_SUMMARY_OBJECT
	ObjectTypeACCOUNT_TRANSACTION_HISTORY_OBJECT
	ObjectTypeBLINDED_BALANCE_OBJECT
	ObjectTypeCHAIN_PROPERTY_OBJECT
	ObjectTypeWITNESS_SCHEDULE_OBJECT
	ObjectTypeBUDGET_RECORD_OBJECT
	ObjectTypeSPECIAL_AUTHORITY_OBJECT
)

/* func (p ObjectType) Space() SpaceType {

	switch p {
	case ObjectTypeBase:
	case ObjectTypeAccount:
	case ObjectTypeAsset:
	case ObjectTypeFORCE_SETTLEMENT_OBJECT:
	case ObjectTypeCOMMITTEE_MEMBER_OBJECT:
	case ObjectTypeWITNESS_OBJECT:
	case ObjectTypeLIMIT_ORDER_OBJECT:
	case ObjectTypeCALL_ORDER_OBJECT:
	case ObjectTypeCUSTOM_OBJECT:
	case ObjectTypePROPOSAL_OBJECT:
	case ObjectTypeOPERATION_HISTORY_OBJECT:
	case ObjectTypeWITHDRAW_PERMISSION_OBJECT:
	case ObjectTypeVESTING_BALANCE_OBJECT:
	case ObjectTypeWORKER_OBJECT:
	case ObjectTypeBALANCE_OBJECT:
		return SpaceTypeProtocol

	case ObjectTypeGLOBAL_PROPERTY_OBJECT:
	case ObjectTypeDYNAMIC_GLOBAL_PROPERTY_OBJECT:
	case ObjectTypeASSET_DYNAMIC_DATA:
	case ObjectTypeAssetBitAssetData:
	case ObjectTypeACCOUNT_BALANCE_OBJECT:
	case ObjectTypeACCOUNT_STATISTICS_OBJECT:
	case ObjectTypeTRANSACTION_OBJECT:
	case ObjectTypeBLOCK_SUMMARY_OBJECT:
	case ObjectTypeACCOUNT_TRANSACTION_HISTORY_OBJECT:
	case ObjectTypeBLINDED_BALANCE_OBJECT:
	case ObjectTypeCHAIN_PROPERTY_OBJECT:
	case ObjectTypeWITNESS_SCHEDULE_OBJECT:
	case ObjectTypeBUDGET_RECORD_OBJECT:
	case ObjectTypeSPECIAL_AUTHORITY_OBJECT:
		return SpaceTypeImplementation
	}

	return SpaceTypeUndefined
}

func (p ObjectType) Type() int {
	var typ = 0

	switch p {
	case ObjectTypeBase:
		typ = 1
		break
	case ObjectTypeAccount:
		typ = 2
		break
	case ObjectTypeAsset:
		typ = 3
		break
	case ObjectTypeFORCE_SETTLEMENT_OBJECT:
		typ = 4
		break
	case ObjectTypeCOMMITTEE_MEMBER_OBJECT:
		typ = 5
		break
	case ObjectTypeWITNESS_OBJECT:
		typ = 6
		break
	case ObjectTypeLIMIT_ORDER_OBJECT:
		typ = 7
		break
	case ObjectTypeCALL_ORDER_OBJECT:
		typ = 8
		break
	case ObjectTypeCUSTOM_OBJECT:
		typ = 9
		break
	case ObjectTypePROPOSAL_OBJECT:
		typ = 10
		break
	case ObjectTypeOPERATION_HISTORY_OBJECT:
		typ = 11
		break
	case ObjectTypeWITHDRAW_PERMISSION_OBJECT:
		typ = 12
		break
	case ObjectTypeVESTING_BALANCE_OBJECT:
		typ = 13
		break
	case ObjectTypeWORKER_OBJECT:
		typ = 14
		break
	case ObjectTypeBALANCE_OBJECT:
		typ = 15
		break
	case ObjectTypeGLOBAL_PROPERTY_OBJECT:
		typ = 0
		break
	case ObjectTypeDYNAMIC_GLOBAL_PROPERTY_OBJECT:
		typ = 1
		break
	case ObjectTypeASSET_DYNAMIC_DATA:
		typ = 3
		break
	case ObjectTypeAssetBitAssetData:
		typ = 4
		break
	case ObjectTypeACCOUNT_BALANCE_OBJECT:
		typ = 5
		break
	case ObjectTypeACCOUNT_STATISTICS_OBJECT:
		typ = 6
		break
	case ObjectTypeTRANSACTION_OBJECT:
		typ = 7
		break
	case ObjectTypeBLOCK_SUMMARY_OBJECT:
		typ = 8
		break
	case ObjectTypeACCOUNT_TRANSACTION_HISTORY_OBJECT:
		typ = 9
		break
	case ObjectTypeBLINDED_BALANCE_OBJECT:
		typ = 10
		break
	case ObjectTypeCHAIN_PROPERTY_OBJECT:
		typ = 11
		break
	case ObjectTypeWITNESS_SCHEDULE_OBJECT:
		typ = 12
		break
	case ObjectTypeBUDGET_RECORD_OBJECT:
		typ = 13
		break
	case ObjectTypeSPECIAL_AUTHORITY_OBJECT:
		typ = 14
	}

	return typ
}
*/
//GenericObjectID is used to return the generic object type in the form space.type.0.
//
// Not to be confused with {@link GrapheneObject#getObjectId()}, which will return
// the full object id in the form space.type.id.
//
// @return: The generic object type
//
/*
func (p ObjectType) GenericObjectID() string {
	return fmt.Sprintf("%d.%d.0", p.Space(), p.Type())
}
*/
type UInt64 uint64

func (p *UInt64) UnmarshalJSON(s []byte) error {
	str := string(s)

	if len(str) > 0 && str != "null" {
		q, err := util.SafeUnquote(str)
		if err != nil {
			return errors.Annotate(err, "parse UInt64 ->[unquote]")
		}

		*(*uint64)(p), err = strconv.ParseUint(q, 10, 64)
		if err != nil {
			return errors.Annotate(err, "parse UInt64 ->[parse]")
		}

		return nil
	}

	return errors.Errorf("unable to unmarshal UInt64 from %s", str)
}

type Float64 float64

func (p *Float64) UnmarshalJSON(s []byte) error {
	str := string(s)

	if len(str) > 0 && str != "null" {
		q, err := util.SafeUnquote(str)
		if err != nil {
			return errors.Annotate(err, "parse Float64 ->[unquote]")
		}

		*(*float64)(p), err = strconv.ParseFloat(q, 64)
		if err != nil {
			return errors.Annotate(err, "parse Float64 ->[parse]")
		}

		return nil
	}

	return errors.Errorf("unable to unmarshal Float64 from %s", str)
}

package objects

import (
	"fmt"
	"strconv"
	"strings"
)

type GrapheneObject interface {
	Id() ObjectID
	Type() ObjectType
}

type GrapheneID struct {
	ID       ObjectID `json:"id"`
	space    SpaceType
	typ      ObjectType
	instance int64
}

func (p GrapheneID) Id() ObjectID {
	return ObjectID(fmt.Sprintf("%d.%d.%d", p.space, p.typ, p.instance))
}

func (p GrapheneID) Type() ObjectType {

	switch p.space {
	case SpaceTypeProtocol:
		switch p.typ {
		case 1:
			return ObjectTypeBASE_OBJECT
		case 2:
			return ObjectTypeACCOUNT_OBJECT
		case 3:
			return ObjectTypeASSET_OBJECT
		case 4:
			return ObjectTypeFORCE_SETTLEMENT_OBJECT
		case 5:
			return ObjectTypeCOMMITTEE_MEMBER_OBJECT
		case 6:
			return ObjectTypeWITNESS_OBJECT
		case 7:
			return ObjectTypeLIMIT_ORDER_OBJECT
		case 8:
			return ObjectTypeCALL_ORDER_OBJECT
		case 9:
			return ObjectTypeCUSTOM_OBJECT
		case 10:
			return ObjectTypePROPOSAL_OBJECT
		case 11:
			return ObjectTypeOPERATION_HISTORY_OBJECT
		case 12:
			return ObjectTypeWITHDRAW_PERMISSION_OBJECT
		case 13:
			return ObjectTypeVESTING_BALANCE_OBJECT
		case 14:
			return ObjectTypeWORKER_OBJECT
		case 15:
			return ObjectTypeBALANCE_OBJECT
		}

	case SpaceTypeImplementation:
		switch p.typ {
		case 0:
			return ObjectTypeGLOBAL_PROPERTY_OBJECT
		case 1:
			return ObjectTypeDYNAMIC_GLOBAL_PROPERTY_OBJECT
		case 3:
			return ObjectTypeASSET_DYNAMIC_DATA
		case 4:
			return ObjectTypeASSET_BITASSET_DATA
		case 5:
			return ObjectTypeACCOUNT_BALANCE_OBJECT
		case 6:
			return ObjectTypeACCOUNT_STATISTICS_OBJECT
		case 7:
			return ObjectTypeTRANSACTION_OBJECT
		case 8:
			return ObjectTypeBLOCK_SUMMARY_OBJECT
		case 9:
			return ObjectTypeACCOUNT_TRANSACTION_HISTORY_OBJECT
		case 10:
			return ObjectTypeBLINDED_BALANCE_OBJECT
		case 11:
			return ObjectTypeCHAIN_PROPERTY_OBJECT
		case 12:
			return ObjectTypeWITNESS_SCHEDULE_OBJECT
		case 13:
			return ObjectTypeBUDGET_RECORD_OBJECT
		case 14:
			return ObjectTypeSPECIAL_AUTHORITY_OBJECT
		}
	}

	return ObjectTypeUndefined
}

func NewGrapheneID(id string) *GrapheneID {

	gid := GrapheneID{
		ID: ObjectID(id),
	}

	parts := strings.Split(id, ".")
	if len(parts) == 3 {
		space, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(fmt.Sprintf("unable to parse GrapheneID [space] from %s", id))
		}

		gid.space = SpaceType(space)

		typ, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(fmt.Sprintf("unable to parse GrapheneID [type] from %s", id))
		}

		gid.typ = ObjectType(typ)

		inst, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			panic(fmt.Sprintf("unable to parse GrapheneID [instance] from %s", id))
		}

		gid.instance = inst
	}

	return &gid
}

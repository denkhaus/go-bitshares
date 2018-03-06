package objects

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type ObjectID string

type AssetType int

const (
	AssetTypeUndefined AssetType = -1
	AssetTypeCoreAsset AssetType = iota
	AssetTypeUIA
	AssetTypeSmartCoin
	AssetTypePredictionMarket
)

type SpaceType Int8

const (
	SpaceTypeUndefined SpaceType = -1
	SpaceTypeProtocol  SpaceType = iota
	SpaceTypeImplementation
)

type OperationType Int8

const (
	OperationTypeTransfer OperationType = iota
	OperationTypeLimitOrderCreate
	OperationTypeLimitOrderCancel
	OperationTypeCallOrderUpdate
	OperationTypeFILL_ORDER_OPERATION
	OperationTypeACCOUNT_CREATE_OPERATION
	OperationTypeACCOUNT_UPDATE_OPERATION
	OperationTypeACCOUNT_WHITELIST_OPERATION
	OperationTypeACCOUNT_UPGRADE_OPERATION
	OperationTypeACCOUNT_TRANSFER_OPERATION
	OperationTypeASSET_CREATE_OPERATION
	OperationTypeASSET_UPDATE_OPERATION
	OperationTypeASSET_UPDATE_BITASSET_OPERATION
	OperationTypeASSET_UPDATE_FEED_PRODUCERS_OPERATION
	OperationTypeASSET_ISSUE_OPERATION
	OperationTypeASSET_RESERVE_OPERATION
	OperationTypeASSET_FUND_FEE_POOL_OPERATION
	OperationTypeASSET_SETTLE_OPERATION
	OperationTypeASSET_GLOBAL_SETTLE_OPERATION
	OperationTypeASSET_PUBLISH_FEED_OPERATION
	OperationTypeWITNESS_CREATE_OPERATION
	OperationTypeWITNESS_UPDATE_OPERATION
	OperationTypePROPOSAL_CREATE_OPERATION
	OperationTypePROPOSAL_UPDATE_OPERATION
	OperationTypePROPOSAL_DELETE_OPERATION
	OperationTypeWITHDRAW_PERMISSION_CREATE_OPERATION
	OperationTypeWITHDRAW_PERMISSION_UPDATE_OPERATION
	OperationTypeWITHDRAW_PERMISSION_CLAIM_OPERATION
	OperationTypeWITHDRAW_PERMISSION_DELETE_OPERATION
	OperationTypeCOMMITTEE_MEMBER_CREATE_OPERATION
	OperationTypeCOMMITTEE_MEMBER_UPDATE_OPERATION
	OperationTypeCOMMITTEE_MEMBER_UPDATE_GLOBAL_PARAMETERS_OPERATION
	OperationTypeVESTING_BALANCE_CREATE_OPERATION
	OperationTypeVESTING_BALANCE_WITHDRAW_OPERATION
	OperationTypeWORKER_CREATE_OPERATION
	OperationTypeCUSTOM_OPERATION
	OperationTypeASSERT_OPERATION
	OperationTypeBALANCE_CLAIM_OPERATION
	OperationTypeOVERRIDE_TRANSFER_OPERATION
	OperationTypeTRANSFER_TO_BLIND_OPERATION
	OperationTypeBLIND_TRANSFER_OPERATION
	OperationTypeTRANSFER_FROM_BLIND_OPERATION
	OperationTypeASSET_SETTLE_CANCEL_OPERATION
	OperationTypeASSET_CLAIM_FEES_OPERATION
)

type ObjectType Int8

const (
	ObjectTypeUndefined                 ObjectType = -1
	ObjectTypeBase                      ObjectType = 1
	ObjectTypeAccount                   ObjectType = 2
	ObjectTypeAsset                     ObjectType = 3
	ObjectTypeForceSettlement           ObjectType = 4
	ObjectTypeCommiteeMember            ObjectType = 5
	ObjectTypeWitness                   ObjectType = 6
	ObjectTypeLimitOrder                ObjectType = 7
	ObjectTypeCallOrder                 ObjectType = 8
	ObjectTypeCustom                    ObjectType = 9
	ObjectTypeProposal                  ObjectType = 10
	ObjectTypeOperationHistory          ObjectType = 11
	ObjectTypeWithdrawPermission        ObjectType = 12
	ObjectTypeVestingBalance            ObjectType = 13
	ObjectTypeWorker                    ObjectType = 14
	ObjectTypeBalance                   ObjectType = 15
	ObjectTypeGlobalProperty            ObjectType = 1
	ObjectTypeDynamicGlobalProperty     ObjectType = 2
	ObjectTypeAssetDynamicData          ObjectType = 3
	ObjectTypeAssetBitAssetData         ObjectType = 4
	ObjectTypeAccountBalance            ObjectType = 5
	ObjectTypeAccountStatistics         ObjectType = 6
	ObjectTypeTransaction               ObjectType = 7
	ObjectTypeBlockSummary              ObjectType = 8
	ObjectTypeAccountTransactionHistory ObjectType = 9
	ObjectTypeBlindedBalance            ObjectType = 10
	ObjectTypeChainProperty             ObjectType = 11
	ObjectTypeWitnessSchedule           ObjectType = 12
	ObjectTypeBudgetRecord              ObjectType = 13
	ObjectTypeSpecialAuthority          ObjectType = 14
)

type Rate float64

func (p Rate) Inverse() Rate {
	return 1 / p
}

func (p Rate) Value() float64 {
	return float64(p)
}

func unmarshalUInt(data []byte) (uint64, error) {
	if len(data) == 0 {
		return 0, errors.New("unmarshalUInt: empty input")
	}

	var (
		res uint64
		err error
		len = len(data)
	)

	if data[0] == '"' && data[len-1] == '"' {
		data := data[1 : len-1]
		res, err = strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			return 0, errors.Errorf("unmarshalUInt: unable to parse input %v", data)
		}
	} else if err := json.Unmarshal(data, &res); err != nil {
		return 0, errors.Errorf("unmarshalUInt: unable to unmarshal input %v", data)
	}

	return res, nil
}

func unmarshalInt(data []byte) (int64, error) {
	if len(data) == 0 {
		return 0, errors.New("unmarshalInt: empty input")
	}

	var (
		res int64
		err error
		len = len(data)
	)

	if data[0] == '"' && data[len-1] == '"' {
		data := data[1 : len-1]
		res, err = strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return 0, errors.Errorf("unmarshalInt: unable to parse input %v", data)
		}
	} else if err := json.Unmarshal(data, &res); err != nil {
		return 0, errors.Errorf("unmarshalInt: unable to unmarshal input %v", data)
	}

	return res, nil
}

func unmarshalFloat(data []byte) (float64, error) {
	if len(data) == 0 {
		return 0, errors.New("unmarshalFloat: empty input")
	}

	var (
		res float64
		err error
		len = len(data)
	)

	if data[0] == '"' && data[len-1] == '"' {
		data := data[1 : len-1]
		res, err = strconv.ParseFloat(string(data), 64)
		if err != nil {
			return 0, errors.Errorf("unmarshalFloat: unable to parse input %v", data)
		}
	} else if err := json.Unmarshal(data, &res); err != nil {
		return 0, errors.Errorf("unmarshalFloat: unable to unmarshal input %v", data)
	}

	return res, nil
}

type UInt uint

func (num *UInt) UnmarshalJSON(data []byte) error {
	v, err := unmarshalUInt(data)
	if err != nil {
		return err
	}

	*num = UInt(v)
	return nil
}

func (num UInt) Marshal(enc *util.TypeEncoder) error {
	return enc.EncodeNumber(uint(num))
}

type UInt8 uint8

func (num *UInt8) UnmarshalJSON(data []byte) error {
	v, err := unmarshalUInt(data)
	if err != nil {
		return err
	}

	*num = UInt8(v)
	return nil
}

func (num UInt8) Marshal(enc *util.TypeEncoder) error {
	return enc.EncodeNumber(uint8(num))
}

type UInt16 uint16

func (num *UInt16) UnmarshalJSON(data []byte) error {
	v, err := unmarshalUInt(data)
	if err != nil {
		return err
	}

	*num = UInt16(v)
	return nil
}

func (num UInt16) Marshal(enc *util.TypeEncoder) error {
	return enc.EncodeNumber(uint16(num))
}

type UInt32 uint32

func (num *UInt32) UnmarshalJSON(data []byte) error {
	v, err := unmarshalUInt(data)
	if err != nil {
		return err
	}

	*num = UInt32(v)
	return nil
}

func (num UInt32) Marshal(enc *util.TypeEncoder) error {
	return enc.EncodeNumber(uint32(num))
}

type UInt64 uint64

func (num *UInt64) UnmarshalJSON(data []byte) error {
	v, err := unmarshalUInt(data)
	if err != nil {
		return err
	}

	*num = UInt64(v)
	return nil
}

func (num UInt64) Marshal(enc *util.TypeEncoder) error {
	return enc.EncodeNumber(uint64(num))
}

type Int8 int8

func (num *Int8) UnmarshalJSON(data []byte) error {
	v, err := unmarshalInt(data)
	if err != nil {
		return err
	}

	*num = Int8(v)
	return nil
}

func (num Int8) Marshal(enc *util.TypeEncoder) error {
	return enc.EncodeNumber(int8(num))
}

type Int16 int16

func (num *Int16) UnmarshalJSON(data []byte) error {
	v, err := unmarshalInt(data)
	if err != nil {
		return err
	}

	*num = Int16(v)
	return nil
}

func (num Int16) Marshal(enc *util.TypeEncoder) error {
	return enc.EncodeNumber(int16(num))
}

type Int32 int32

func (num *Int32) UnmarshalJSON(data []byte) error {
	v, err := unmarshalInt(data)
	if err != nil {
		return err
	}

	*num = Int32(v)
	return nil
}

func (num Int32) Marshal(enc *util.TypeEncoder) error {
	return enc.EncodeNumber(int32(num))
}

type Int64 int64

func (num *Int64) UnmarshalJSON(data []byte) error {
	v, err := unmarshalInt(data)
	if err != nil {
		return err
	}

	*num = Int64(v)
	return nil
}

func (num Int64) Marshal(enc *util.TypeEncoder) error {
	return enc.EncodeNumber(int64(num))
}

type Float32 float32

func (num *Float32) UnmarshalJSON(data []byte) error {
	v, err := unmarshalFloat(data)
	if err != nil {
		return err
	}

	*num = Float32(v)
	return nil
}

func (num Float32) Marshal(enc *util.TypeEncoder) error {
	return enc.EncodeNumber(float32(num))
}

type Float64 float64

func (num *Float64) UnmarshalJSON(data []byte) error {
	v, err := unmarshalFloat(data)
	if err != nil {
		return err
	}

	*num = Float64(v)
	return nil
}

func (num Float64) Marshal(enc *util.TypeEncoder) error {
	return enc.EncodeNumber(float64(num))
}

const TimeFormat = `"2006-01-02T15:04:05"`

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(t.Time.Format(TimeFormat)), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	parsed, err := time.ParseInLocation(TimeFormat, string(data), time.UTC)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

func (t Time) Marshal(enc *util.TypeEncoder) error {
	return enc.Encode(uint32(t.Time.Unix()))
}

func (t Time) Add(dur time.Duration) Time {
	return Time{t.Time.Add(dur)}
}

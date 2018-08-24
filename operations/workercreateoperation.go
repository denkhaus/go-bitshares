package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &WorkerCreateOperation{}
	types.OperationMap[op.Type()] = op
}

type WorkerCreateOperation struct {
	types.OperationFee
	DailyPay      types.UInt64            `json:"daily_pay"`
	Initializer   types.WorkerInitializer `json:"initializer"`
	Name          string                  `json:"name"`
	Owner         types.GrapheneID        `json:"owner"`
	URL           string                  `json:"url"`
	WorkBeginDate types.Time              `json:"work_begin_date"`
	WorkEndDate   types.Time              `json:"work_end_date"`
}

func (p WorkerCreateOperation) Type() types.OperationType {
	return types.OperationTypeWorkerCreate
}

func (p WorkerCreateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.Owner); err != nil {
		return errors.Annotate(err, "encode Owner")
	}

	if err := enc.Encode(p.WorkBeginDate); err != nil {
		return errors.Annotate(err, "encode WorkBeginDate")
	}

	if err := enc.Encode(p.WorkEndDate); err != nil {
		return errors.Annotate(err, "encode WorkEndDate")
	}

	if err := enc.Encode(p.DailyPay); err != nil {
		return errors.Annotate(err, "encode DailyPay")
	}

	if err := enc.Encode(p.Name); err != nil {
		return errors.Annotate(err, "encode Name")
	}

	if err := enc.Encode(p.URL); err != nil {
		return errors.Annotate(err, "encode URL")
	}

	if err := enc.Encode(p.Initializer); err != nil {
		return errors.Annotate(err, "encode Initializer")
	}

	return nil
}

//NewWorkerCreateOperation creates a new WorkerCreateOperation
func NewWorkerCreateOperation() *WorkerCreateOperation {
	tx := WorkerCreateOperation{}
	return &tx
}

package types

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"

	"github.com/juju/errors"
)

//go:generate ffjson $GOFILE

type DynamicGlobalProperties struct {
	ID                             DynamicGlobalPropertyID `json:"id"`
	CurrentWitness                 WitnessID               `json:"current_witness"`
	LastBudgetTime                 Time                    `json:"last_budget_time"`
	Time                           Time                    `json:"time"`
	NextMaintenanceTime            Time                    `json:"next_maintenance_time"`
	AccountsRegisteredThisInterval int                     `json:"accounts_registered_this_interval"`
	DynamicFlags                   int                     `json:"dynamic_flags"`
	HeadBlockID                    String                  `json:"head_block_id"`
	RecentSlotsFilled              String                  `json:"recent_slots_filled"`
	HeadBlockNumber                UInt32                  `json:"head_block_number"`
	LastIrreversibleBlockNum       UInt32                  `json:"last_irreversible_block_num"`
	CurrentAslot                   int64                   `json:"current_aslot"`
	WitnessBudget                  int64                   `json:"witness_budget"`
	RecentlyMissedCount            int64                   `json:"recently_missed_count"`
}

func (p DynamicGlobalProperties) RefBlockNum() UInt16 {
	return UInt16(p.HeadBlockNumber)
}

func (p DynamicGlobalProperties) RefBlockPrefix() (UInt32, error) {
	rawBlockID, err := hex.DecodeString(p.HeadBlockID.String())
	if err != nil {
		return 0, errors.Annotatef(err, "DecodeString HeadBlockID: %v", p.HeadBlockID)
	}

	if len(rawBlockID) < 8 {
		return 0, errors.Errorf("invalid HeadBlockID: %v", p.HeadBlockID)
	}

	rawPrefix := rawBlockID[4:8]

	var prefix uint32
	if err := binary.Read(bytes.NewReader(rawPrefix), binary.LittleEndian, &prefix); err != nil {
		return 0, errors.Annotatef(err, "failed to read block prefix: %v", rawPrefix)
	}

	return UInt32(prefix), nil
}

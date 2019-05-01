package types

type Info struct {
	ActiveCommitteeMembers CommitteeMemberIDs `json:"active_committee_members"`
	ActiveWitnesses        WitnessIDs         `json:"active_witnesses"`
	ChainID                String             `json:"chain_id"`
	HeadBlockAge           String             `json:"head_block_age"`
	HeadBlockID            Buffer             `json:"head_block_id"`
	HeadBlockNum           UInt64             `json:"head_block_num"`
	NextMaintenanceTime    String             `json:"next_maintenance_time"`
	Participation          String             `json:"participation"`
}

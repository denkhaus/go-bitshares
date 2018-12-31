package types

//go:generate ffjson $GOFILE
type CommiteeMember struct {
	ID                    GrapheneID `json:"id"`
	CommiteeMemberAccount GrapheneID `json:"committee_member_account"`
	TotalVotes            UInt64     `json:"total_votes"`
	URL                   String     `json:"url"`
	VoteID                VoteID     `json:"vote_id"`
}

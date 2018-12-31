package operations

func (suite *operationsAPITest) Test_ProposalDeleteOperation() {
	op := ProposalDeleteOperation{}

	suite.OpSamplesTest(&op)
}

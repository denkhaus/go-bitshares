package operations

func (suite *operationsAPITest) Test_BalanceClaimOperation() {
	op := BalanceClaimOperation{}

	suite.OpSamplesTest(&op)
}

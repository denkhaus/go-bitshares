package operations

func (suite *operationsAPITest) Test_WitnessUpdateOperation() {
	op := WitnessUpdateOperation{}
	suite.OpSamplesTest(&op)
}

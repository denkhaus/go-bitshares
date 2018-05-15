package operations

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/tests"
	"github.com/denkhaus/bitshares/types"

	"github.com/denkhaus/bitshares/util"
	"github.com/stretchr/testify/suite"

	// importing this initializes sample data fetching
	_ "github.com/denkhaus/bitshares/gen/samples"
)

//Note: operation tests may fail for now cause extensions marshalling is questionable.
type operationsAPITest struct {
	suite.Suite
	TestAPI api.BitsharesAPI
	RefTx   *types.Transaction
}

func (suite *operationsAPITest) SetupTest() {
	api := api.New(tests.WsFullApiUrl, tests.RpcApiUrl)

	if err := api.Connect(); err != nil {
		suite.Fail(err.Error(), "Connect")
	}

	api.OnError(func(err error) {
		suite.Fail(err.Error(), "OnError")
	})

	suite.TestAPI = api

	tx := types.NewTransaction()
	tx.RefBlockNum = 34294
	tx.RefBlockPrefix = 3707022213

	if err := tx.Expiration.UnmarshalJSON([]byte(`"2016-04-06T08:29:27"`)); err != nil {
		suite.Fail(err.Error(), "Unmarshal expiration")
	}

	suite.RefTx = tx
}

func (suite *operationsAPITest) TearDown() {
	if err := suite.TestAPI.Close(); err != nil {
		suite.Fail(err.Error(), "Close")
	}
}

func (suite *operationsAPITest) Test_SerializeEmptyTransaction() {

	tx := types.NewTransaction()
	if err := tx.Expiration.UnmarshalJSON([]byte(`"2016-04-06T08:29:27"`)); err != nil {
		suite.Fail(err.Error(), "Unmarshal expiration")
	}

	suite.compareTransaction(tx)
}

func (suite *operationsAPITest) Test_SerializeTransaction() {
	hex, err := suite.TestAPI.SerializeTransaction(suite.RefTx)
	if err != nil {
		suite.Fail(err.Error(), "SerializeTransaction")
	}

	suite.NotNil(hex)
	suite.Equal("f68585abf4dce7c80457000000", hex)
}

func (suite *operationsAPITest) compareTransaction(tx *types.Transaction) {
	var buf bytes.Buffer
	enc := util.NewTypeEncoder(&buf)
	if err := enc.Encode(tx); err != nil {
		suite.Fail(err.Error(), "Encode")
	}

	ref, err := suite.TestAPI.SerializeTransaction(tx)
	if err != nil {
		suite.Fail(err.Error(), "SerializeTransaction")
	}

	test := hex.EncodeToString(buf.Bytes())

	// util.Dump("ref", ref)
	// util.Dump("test", test)

	suite.Equal(ref, test)
}

func TestOperations(t *testing.T) {
	testSuite := new(operationsAPITest)
	suite.Run(t, testSuite)
	testSuite.TearDown()
}

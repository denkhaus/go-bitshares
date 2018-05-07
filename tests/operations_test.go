package tests

import (
	"bytes"
	"encoding/hex"
	"testing"
	"time"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/objects"

	"github.com/denkhaus/bitshares/util"
	"github.com/stretchr/testify/suite"
)

//Note: operation tests may fail for now cause extensions marshalling is questionable.
type operationsAPITest struct {
	suite.Suite
	TestAPI api.BitsharesAPI
	RefTx   *objects.Transaction
}

func (suite *operationsAPITest) SetupTest() {
	api := api.New(wsTestApiUrl, rpcApiUrl)

	if err := api.Connect(); err != nil {
		suite.Fail(err.Error(), "Connect")
	}

	api.OnError(func(err error) {
		suite.Fail(err.Error(), "OnError")
	})

	suite.TestAPI = api

	tx := objects.NewTransaction()
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

func (suite *operationsAPITest) Test_SerializeTransaction() {
	hex, err := suite.TestAPI.SerializeTransaction(suite.RefTx)
	if err != nil {
		suite.Fail(err.Error(), "SerializeTransaction")
	}

	suite.NotNil(hex)
	suite.Equal("f68585abf4dce7c80457000000", hex)
}

func (suite *operationsAPITest) Test_CallOrderUpdateOperation() {
	time.Sleep(1 * time.Second)

	op := objects.CallOrderUpdateOperation{
		Extensions:     objects.Extensions{},
		FundingAccount: *objects.NewGrapheneID("1.2.29"),
		DeltaCollateral: objects.AssetAmount{
			Amount: 100000000,
			Asset:  *objects.NewGrapheneID("1.3.0"),
		},
		DeltaDebt: objects.AssetAmount{
			Amount: 10000,
			Asset:  *objects.NewGrapheneID("1.3.22"),
		},
		Fee: objects.AssetAmount{
			Amount: 100,
			Asset:  *objects.NewGrapheneID("1.3.0"),
		},
	}

	suite.RefTx.Operations = objects.Operations{
		objects.Operation(&op),
	}

	suite.compareTransaction(suite.RefTx)
}

func (suite *operationsAPITest) Test_LimitOrderCreateOperation() {
	time.Sleep(1 * time.Second)

	op := objects.LimitOrderCreateOperation{
		Extensions: objects.Extensions{},
		Seller:     *objects.NewGrapheneID("1.2.29"),
		FillOrKill: false,
		Fee: objects.AssetAmount{
			Amount: 100,
			Asset:  *objects.NewGrapheneID("1.3.0"),
		},
		AmountToSell: objects.AssetAmount{
			Amount: 100000,
			Asset:  *objects.NewGrapheneID("1.3.0"),
		},
		MinToReceive: objects.AssetAmount{
			Amount: 10000,
			Asset:  *objects.NewGrapheneID("1.3.105"),
		},
	}

	if err := op.Expiration.UnmarshalJSON([]byte(`"2016-05-18T09:22:05"`)); err != nil {
		suite.Fail(err.Error(), "Unmarshal expiration")
	}

	suite.RefTx.Operations = objects.Operations{
		objects.Operation(&op),
	}

	suite.compareTransaction(suite.RefTx)
}

func (suite *operationsAPITest) Test_LimitOrderCancelOperation() {
	time.Sleep(1 * time.Second)

	op := objects.LimitOrderCancelOperation{
		Extensions:       objects.Extensions{},
		Order:            *objects.NewGrapheneID("1.7.123"),
		FeePayingAccount: *objects.NewGrapheneID("1.2.456"),
		Fee: objects.AssetAmount{
			Amount: 1000,
			Asset:  *objects.NewGrapheneID("1.3.789"),
		},
	}

	suite.RefTx.Operations = objects.Operations{
		objects.Operation(&op),
	}

	suite.compareTransaction(suite.RefTx)
}

func (suite *operationsAPITest) compareTransaction(tx *objects.Transaction) {
	var buf bytes.Buffer
	enc := util.NewTypeEncoder(&buf)
	if err := enc.Encode(tx); err != nil {
		suite.Fail(err.Error(), "Encode")
	}

	ref, err := suite.TestAPI.SerializeTransaction(suite.RefTx)
	if err != nil {
		suite.Fail(err.Error(), "SerializeTransaction")
	}

	test := hex.EncodeToString(buf.Bytes())
	suite.Equal(ref, test)
}

func TestOperations(t *testing.T) {
	testSuite := new(operationsAPITest)
	suite.Run(t, testSuite)
	testSuite.TearDown()
}

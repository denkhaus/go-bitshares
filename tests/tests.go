package tests

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

const (
	WsFullApiUrl = "wss://node.market.rudex.org"
	//WsFullApiUrl = "wss://bitshares.openledger.info/ws"
	WsTestApiUrl  = "wss://node.testnet.bitshares.eu/ws"
	RpcFullApiUrl = "http://localhost:8095"
	RpcTestApiUrl = "http://localhost:8094"
)

var (
	ChainIDBitSharesFull = "4018d7844c78f6a6c41c6a552b898022310fc5dec06da467ee7905a8dad512c8"
	ChainIDBitSharesTest = "39f5e2ede1f8bc1a3a54a7914414e3779e33193f1f5693510e73cb7a87617447"
	UserID1              = types.NewGrapheneID("1.2.282")  //xeroc user account
	UserID2              = types.NewGrapheneID("1.2.253")  //stan user account
	UserID3              = types.NewGrapheneID("1.2.0")    //committee-account user account
	UserID4              = types.NewGrapheneID("1.2.1751") //denkhaus user account
	AssetCNY             = types.NewGrapheneID("1.3.113")  //cny asset
	AssetBTS             = types.NewGrapheneID("1.3.0")    //bts asset
	AssetUSD             = types.NewGrapheneID("1.3.121")  // usd asset
	AssetTEST            = types.NewGrapheneID("1.3.0")    // test asset
	AssetPEGFAKEUSD      = types.NewGrapheneID("1.3.22")   // test asset
	AssetBTC             = types.NewGrapheneID("1.3.103")
	AssetSILVER          = types.NewGrapheneID("1.3.105")
	AssetGOLD            = types.NewGrapheneID("1.3.106")
	AssetEUR             = types.NewGrapheneID("1.3.120")
	AssetOBITS           = types.NewGrapheneID("1.3.562")
	AssetOpenETH         = types.NewGrapheneID("1.3.850")
	AssetOpenLTC         = types.NewGrapheneID("1.3.859")
	AssetOpenBTC         = types.NewGrapheneID("1.3.861")
	AssetOpenSTEEM       = types.NewGrapheneID("1.3.973")
	AssetOpenUSDT        = types.NewGrapheneID("1.3.1042")
	AssetYOYOW           = types.NewGrapheneID("1.3.1093")
	AssetRUBEL           = types.NewGrapheneID("1.3.1325")
	AssetHERO            = types.NewGrapheneID("1.3.1362")

	SettleOrder1      = types.NewGrapheneID("1.4.1655")       // random SettleOrder ObjectID
	CommiteeMember1   = types.NewGrapheneID("1.5.15")         // random CommiteeMember ObjectID
	LimitOrder1       = types.NewGrapheneID("1.7.75961600")   // random LimitOrder ObjectID
	CallOrder1        = types.NewGrapheneID("1.8.4582")       // random CallOrder ObjectID
	OperationHistory1 = types.NewGrapheneID("1.11.187698971") // random OperationHistory ObjectID
	Balance1          = types.NewGrapheneID("1.15.1")         // random Balance ObjectID
	BitAssetDataCNY   = types.NewGrapheneID("2.4.13")         // cny bitasset data id

	TestAccount1UserName      = "denk-haus"
	TestAccount1Password      = "denkhaus-testnet"
	TestAccount1PubKeyActive  = "TEST5zzvbDtkbUVU1gFFsKqCE55U7JbjTp6mTh1usFv7KGgXL7HDQk"
	TestAccount1PrivKeyActive = "5Hx8KiHLnc3pDLkwe2jujkTTJev72n3Qx7xtyaRNBsJDuejzh9u"
	TestAccount1PubKeyOwner   = "TEST5yXqEBShUgcVm7Mve8Fg4RzQ2ftPpmo77aMbz884eX9aeGVvwD"
	TestAccount1PrivKeyOwner  = "5JyuWmopuyxFyvM9xm8fxXyujzfVnsg9cvE6z3pcib5NW1Av4rP"
	TestAccount1PrivKeyMemo   = "TEST5zzvbDtkbUVU1gFFsKqCE55U7JbjTp6mTh1usFv7KGgXL7HDQk"
	TestAccount1ID            = types.NewGrapheneID("1.2.3464")

	TestAccount2UserName = "denk-baum"
	TestAccount2Password = "denkhaus-testnet"
	TestAccount2PubKey   = "TEST5Z3vsgH6xj6HLXcsU38yo4TyoZs9AUzpfbaXbuxsAYPbutWvEP"
	TestAccount2PrivKey  = "5KRZv3ZmkcE71K9KwEKG6pV6pyufkMQgCJrCu8xKLf2y7R7J8YK"
	TestAccount2ID       = types.NewGrapheneID("1.2.3496")

	TestAccount3UserName      = "bs-test"
	TestAccount3Password      = "denkhaus-test"
	TestAccount3PubKeyActive  = "BTS5shffTjVoT4J8Zrj3f2mQJw4UVKrnbx5FWYhVgov45EpBf2NYi"
	TestAccount3PrivKeyActive = "5JTge2oTwFqfNPhUrrm6upheByG2VXvaXBAqWdDUvK2CsygMG3Z"
	TestAccount3PubKeyOwner   = "BTS56fy8qpkLzNoguGMPgPNkkznxnx2woEg1qPq7E6gF2SeGSRyK5"
	TestAccount3PrivKeyOwner  = "5JqmjeakPoTz3ComQ7Jgg11jHxywfkJHZPhMJoBomZLrZSfRAvr"
	TestAccount3ID            = types.NewGrapheneID("1.2.391614")
)

func CompareTransactions(api api.BitsharesAPI, tx *types.Transaction, debug bool) (string, string, error) {
	var buf bytes.Buffer
	enc := util.NewTypeEncoder(&buf)
	if err := enc.Encode(tx); err != nil {
		return "", "", errors.Annotate(err, "Encode")
	}

	api.SetDebug(debug)
	ref, err := api.SerializeTransaction(tx)
	if err != nil {
		return "", "", errors.Annotate(err, "SerializeTransaction")
	}

	test := hex.EncodeToString(buf.Bytes())
	return ref, test, nil
}

func NewTestAPI(t *testing.T, wsAPIEndpoint, rpcAPIEndpoint string) api.BitsharesAPI {
	api := api.New(wsAPIEndpoint, rpcAPIEndpoint)
	if err := api.Connect(); err != nil {
		assert.FailNow(t, err.Error(), "Connect")
	}

	api.OnError(func(err error) {
		assert.FailNow(t, err.Error(), "OnError")
	})

	return api
}

func CreateRefTransaction(t *testing.T) *types.Transaction {
	tx := types.NewTransaction()
	tx.RefBlockNum = 34294
	tx.RefBlockPrefix = 3707022213

	if err := tx.Expiration.UnmarshalJSON([]byte(`"2030-04-06T08:29:27"`)); err != nil {
		assert.FailNow(t, err.Error(), "Unmarshal expiration")
	}

	return tx
}

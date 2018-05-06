package tests

import "github.com/denkhaus/bitshares/objects"

const (
	wsFullApiUrl = "wss://node.market.rudex.org"
	//wsFullApiUrl = "wss://bitshares.openledger.info/ws"
	wsTestApiUrl = "wss://node.testnet.bitshares.eu/ws"
	rpcApiUrl    = "http://localhost:8095"
)

var (
	ChainIDBitSharesFull = "4018d7844c78f6a6c41c6a552b898022310fc5dec06da467ee7905a8dad512c8"
	ChainIDBitSharesTest = "39f5e2ede1f8bc1a3a54a7914414e3779e33193f1f5693510e73cb7a87617447"
	UserID1              = objects.NewGrapheneID("1.2.282")  //xeroc user account
	UserID2              = objects.NewGrapheneID("1.2.253")  //stan user account
	UserID3              = objects.NewGrapheneID("1.2.0")    //committee-account user account
	UserID4              = objects.NewGrapheneID("1.2.1751") //denkhaus user account
	AssetCNY             = objects.NewGrapheneID("1.3.113")  //cny asset
	AssetBTS             = objects.NewGrapheneID("1.3.0")    //bts asset
	AssetUSD             = objects.NewGrapheneID("1.3.121")  // usd asset
	AssetTEST            = objects.NewGrapheneID("1.3.0")    // test asset
	AssetPEGFAKEUSD      = objects.NewGrapheneID("1.3.22")   // test asset
	AssetBTC             = objects.NewGrapheneID("1.3.103")
	AssetSILVER          = objects.NewGrapheneID("1.3.105")
	AssetGOLD            = objects.NewGrapheneID("1.3.106")
	AssetEUR             = objects.NewGrapheneID("1.3.120")
	AssetOBITS           = objects.NewGrapheneID("1.3.562")
	AssetOpenETH         = objects.NewGrapheneID("1.3.850")
	AssetOpenLTC         = objects.NewGrapheneID("1.3.859")
	AssetOpenBTC         = objects.NewGrapheneID("1.3.861")
	AssetOpenSTEEM       = objects.NewGrapheneID("1.3.973")
	AssetOpenUSDT        = objects.NewGrapheneID("1.3.1042")
	AssetYOYOW           = objects.NewGrapheneID("1.3.1093")
	AssetRUBEL           = objects.NewGrapheneID("1.3.1325")
	AssetHERO            = objects.NewGrapheneID("1.3.1362")

	BitAssetDataCNY = objects.NewGrapheneID("2.4.13")       // cny bitasset data id
	LimitOrder1     = objects.NewGrapheneID("1.7.22765740") // random LimitOrder ObjectID
	CallOrder1      = objects.NewGrapheneID("1.8.4582")     // random CallOrder ObjectID
	SettleOrder1    = objects.NewGrapheneID("1.4.1655")     // random SettleOrder ObjectID

	TestAccount1UserName      = "denk-haus"
	TestAccount1Password      = "denkhaus-testnet"
	TestAccount1PubKeyActive  = "TEST5zzvbDtkbUVU1gFFsKqCE55U7JbjTp6mTh1usFv7KGgXL7HDQk"
	TestAccount1PrivKeyActive = "5Hx8KiHLnc3pDLkwe2jujkTTJev72n3Qx7xtyaRNBsJDuejzh9u"
	TestAccount1PubKeyOwner   = "TEST5yXqEBShUgcVm7Mve8Fg4RzQ2ftPpmo77aMbz884eX9aeGVvwD"
	TestAccount1PrivKeyOwner  = "5JyuWmopuyxFyvM9xm8fxXyujzfVnsg9cvE6z3pcib5NW1Av4rP"
	TestAccount1PrivKeyMemo   = "TEST5zzvbDtkbUVU1gFFsKqCE55U7JbjTp6mTh1usFv7KGgXL7HDQk"
	TestAccount1ID            = objects.NewGrapheneID("1.2.3464")

	TestAccount2UserName = "denk-baum"
	TestAccount2Password = "denkhaus-testnet"
	TestAccount2PubKey   = "TEST5Z3vsgH6xj6HLXcsU38yo4TyoZs9AUzpfbaXbuxsAYPbutWvEP"
	TestAccount2PrivKey  = "5KRZv3ZmkcE71K9KwEKG6pV6pyufkMQgCJrCu8xKLf2y7R7J8YK"
	TestAccount2ID       = objects.NewGrapheneID("1.2.3496")

	TestAccount3UserName      = "bs-test"
	TestAccount3Password      = "denkhaus-test"
	TestAccount3PubKeyActive  = "BTS5shffTjVoT4J8Zrj3f2mQJw4UVKrnbx5FWYhVgov45EpBf2NYi"
	TestAccount3PrivKeyActive = "5JTge2oTwFqfNPhUrrm6upheByG2VXvaXBAqWdDUvK2CsygMG3Z"
	TestAccount3PubKeyOwner   = "BTS56fy8qpkLzNoguGMPgPNkkznxnx2woEg1qPq7E6gF2SeGSRyK5"
	TestAccount3PrivKeyOwner  = "5JqmjeakPoTz3ComQ7Jgg11jHxywfkJHZPhMJoBomZLrZSfRAvr"
	TestAccount3ID            = objects.NewGrapheneID("1.2.391614")
)

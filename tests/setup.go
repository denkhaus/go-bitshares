package tests

import "github.com/denkhaus/bitshares/objects"

const (
	wsFullApiUrl = "wss://bitshares.openledger.info/ws"
	wsTestApiUrl = "wss://node.testnet.bitshares.eu/ws"
)

var (
	ChainIDBitShares = "4018d7844c78f6a6c41c6a552b898022310fc5dec06da467ee7905a8dad512c8"
	UserID1          = objects.NewGrapheneID("1.2.282")      //xeroc user account
	UserID2          = objects.NewGrapheneID("1.2.253")      //stan user account
	UserID3          = objects.NewGrapheneID("1.2.0")        //committee-account user account
	AssetCNY         = objects.NewGrapheneID("1.3.113")      //cny asset
	AssetBTS         = objects.NewGrapheneID("1.3.0")        //bts asset
	AssetUSD         = objects.NewGrapheneID("1.3.121")      // usd asset
	AssetTEST        = objects.NewGrapheneID("1.3.0")        // test asset
	BitAssetDataCNY  = objects.NewGrapheneID("2.4.13")       //cny bitasset data id
	LimitOrder1      = objects.NewGrapheneID("1.7.22765740") // random LimitOrder ObjectID
	CallOrder1       = objects.NewGrapheneID("1.8.4582")     // random CallOrder ObjectID
	SettleOrder1     = objects.NewGrapheneID("1.4.1655")     // random SettleOrder ObjectID

	TestAccount1UserName      = "denk-haus"
	TestAccount1Password      = "denkhaus-testnet"
	TestAccount1PubKeyActive  = "TEST5zzvbDtkbUVU1gFFsKqCE55U7JbjTp6mTh1usFv7KGgXL7HDQk"
	TestAccount1PrivKeyActive = "5Hx8KiHLnc3pDLkwe2jujkTTJev72n3Qx7xtyaRNBsJDuejzh9u"
	TestAccount1PubKeyOwner   = "TEST5yXqEBShUgcVm7Mve8Fg4RzQ2ftPpmo77aMbz884eX9aeGVvwD"
	TestAccount1PrivKeyOwner  = "5JyuWmopuyxFyvM9xm8fxXyujzfVnsg9cvE6z3pcib5NW1Av4rP"
	TestAccount1ID            = objects.NewGrapheneID("1.2.3464")

	TestAccount2UserName = "denk-baum"
	TestAccount2Password = "denkhaus-testnet"
	TestAccount2PubKey   = "TEST5Z3vsgH6xj6HLXcsU38yo4TyoZs9AUzpfbaXbuxsAYPbutWvEP"
	TestAccount2PrivKey  = "5KRZv3ZmkcE71K9KwEKG6pV6pyufkMQgCJrCu8xKLf2y7R7J8YK"
	TestAccount2ID       = objects.NewGrapheneID("1.2.3496")
)

package tests

import (
	"testing"
	"time"

	"github.com/denkhaus/bitshares"
	"github.com/denkhaus/bitshares/util"
	"github.com/denkhaus/logging"
	"github.com/stretchr/testify/suite"
	"gopkg.in/cheggaaa/pb.v2"
)

const (
	ErrMsgNotEnoughSamples = "not enough incoming samples during test window"

	SubscribeToPendingTransactionsMsgs     = 60
	SubscribeToPendingTransactionsDuration = 20 * time.Second

	SetSubscribeCallbackSubscriberID = 35
	SetSubscribeCallbackMsgs         = 3
	SetSubscribeCallbackDuration     = 60 * time.Second

	SubscribeToMarketMsgs     = 2
	SubscribeToMarketDuration = 90 * time.Second

	SetBlockAppliedCallbackMsgs = 3
	SetBlockAppliedDuration     = 12 * time.Second
)

type subscribeTest struct {
	suite.Suite
	TestAPI bitshares.WebsocketAPI
}

func (suite *subscribeTest) SetupTest() {
	suite.TestAPI = NewWebsocketTestAPI(
		suite.T(),
		WsFullApiUrl,
	)
}

func (suite *subscribeTest) TearDown() {
	if err := suite.TestAPI.Close(); err != nil {
		suite.FailNow(err.Error(), "Close")
	}
}

func (suite *subscribeTest) Test_SubscribeToPendingTransactions() {
	bar := pb.StartNew(SubscribeToPendingTransactionsMsgs).Prefix("wait for transactions")

	err := suite.TestAPI.SubscribeToPendingTransactions(func(in interface{}) error {
		logging.DDump("in", in)
		bar.Increment()
		return nil
	})

	if err != nil {
		suite.FailNow(err.Error(), "SubscribeToPendingTransactions")
	}

	suite.Condition(func() bool {
		return util.WaitForCondition(SubscribeToPendingTransactionsDuration, func() bool {
			return int(bar.Get()) >= SubscribeToPendingTransactionsMsgs
		})
	}, ErrMsgNotEnoughSamples)

	bar.Finish()

	if err := suite.TestAPI.CancelAllSubscriptions(); err != nil {
		suite.FailNow(err.Error(), "CancelAllSubscriptions")
	}
}

func (suite *subscribeTest) Test_SubscribeToMarket() {
	bar := pb.StartNew(SubscribeToMarketMsgs).Prefix("wait for market data")

	err := suite.TestAPI.SubscribeToMarket(AssetBTS, AssetCNY,
		func(in interface{}) error {
			logging.DDump("in", in)
			bar.Increment()
			return nil
		})

	if err != nil {
		suite.FailNow(err.Error(), "SubscribeToMarket")
	}

	suite.Condition(func() bool {
		return util.WaitForCondition(SubscribeToMarketDuration, func() bool {
			return int(bar.Get()) >= SubscribeToMarketMsgs
		})
	}, ErrMsgNotEnoughSamples)

	bar.Finish()

	if err := suite.TestAPI.UnsubscribeFromMarket(AssetBTS, AssetCNY); err != nil {
		suite.FailNow(err.Error(), "UnsubscribeFromMarket")
	}

}

func (suite *subscribeTest) Test_SubscribeToBlockApplied() {
	bar := pb.StartNew(SetBlockAppliedCallbackMsgs).Prefix("wait for block data")

	err := suite.TestAPI.SubscribeToBlockApplied(func(blockID string) error {
		logging.DDump("blockID", blockID)
		bar.Increment()
		return nil
	})

	if err != nil {
		suite.FailNow(err.Error(), "SubscribeToBlockApplied")
	}

	suite.Condition(func() bool {
		return util.WaitForCondition(SetBlockAppliedDuration, func() bool {
			return int(bar.Get()) >= SetBlockAppliedCallbackMsgs
		})
	}, ErrMsgNotEnoughSamples)

	bar.Finish()

	if err := suite.TestAPI.CancelAllSubscriptions(); err != nil {
		suite.FailNow(err.Error(), "CancelAllSubscriptions")
	}
}

// func (suite *subscribeTest) Test_SetSubscribeCallback() {

// 	fmt.Printf("SetSubscribeCallback: wait %s for %d incoming notifications\n",
// 		SetSubscribeCallbackDuration, SetSubscribeCallbackMsgs)
// 	if err := suite.TestAPI.SetSubscribeCallback(SetSubscribeCallbackSubscriberID, false); err != nil {
// 		suite.FailNow(err.Error(), "SetSubscribeCallback")
// 	}

// 	_, err := suite.TestAPI.CallWsAPI(suite.TestAPI.DatabaseAPIID(), "get_objects", []interface{}{"2.1.0"})
// 	if err != nil {
// 		suite.FailNow(err.Error(), "CallAPI->get_objects")
// 	}

// 	bar := pb.StartNew(SetSubscribeCallbackMsgs).Prefix("wait for data")
// 	err = suite.TestAPI.OnSubscribe(SetSubscribeCallbackSubscriberID, func(msg interface{}) error {
// 		bar.Increment()
// 		return nil
// 	})

// 	if err != nil {
// 		log.Fatal(errors.Annotate(err, "OnSubscribe"))
// 	}

// 	suite.Condition(func() bool {
// 		return util.WaitForCondition(SetSubscribeCallbackDuration, func() bool {
// 			return int(bar.Get()) >= SetSubscribeCallbackMsgs
// 		})
// 	}, ErrMsgNotEnoughSamples)

// 	bar.Finish()

// 	fmt.Println("Call CancelAllSubscriptions")
// 	if err := suite.TestAPI.CancelAllSubscriptions(); err != nil {
// 		suite.FailNow(err.Error(), "CancelAllSubscriptions")
// 	}
// }

func TestSubscribe(t *testing.T) {
	testSuite := new(subscribeTest)
	suite.Run(t, testSuite)
	testSuite.TearDown()
}

package examples

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/objects"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/stretchr/testify/suite"
)

const (
	testURL = "wss://bitshares.openledger.info/ws"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type bitsharesAPITest struct {
	suite.Suite
	TestAPI                          api.BitsharesAPI
	AssetCNY                         *objects.GrapheneID
	AssetBTS                         *objects.GrapheneID
	AssetUSD                         *objects.GrapheneID
	SetSubscribeCallbackSubscriberID int
	SetSubscribeCallbackDuration     time.Duration
	SetSubscribeCallbackMsgs         int
	SubscribeToMarketSubscriberID    int
	SubscribeToMarketMsgs            int
	SubscribeToMarketDuration        time.Duration
}

func (suite *bitsharesAPITest) SetupTest() {
	suite.AssetCNY = objects.NewGrapheneID("1.3.113") //cny asset
	suite.AssetBTS = objects.NewGrapheneID("1.3.0")   //bts asset
	suite.AssetUSD = objects.NewGrapheneID("1.3.121") // usd asset

	suite.SetSubscribeCallbackSubscriberID = 5
	suite.SetSubscribeCallbackMsgs = 8
	suite.SetSubscribeCallbackDuration = 30 * time.Second

	suite.SubscribeToMarketSubscriberID = 4
	suite.SubscribeToMarketMsgs = 3
	suite.SubscribeToMarketDuration = 90 * time.Second

	api := api.New(testURL)
	if err := api.Connect(); err != nil {
		suite.Fail(err.Error(), "Connect")
	}

	api.OnError(func(err error) {
		suite.Fail(err.Error(), "OnError")
	})

	suite.TestAPI = api
}

func (suite *bitsharesAPITest) TearDown() {
	if err := suite.TestAPI.Close(); err != nil {
		suite.Fail(err.Error(), "Close")
	}
}

func (suite *bitsharesAPITest) Test_SubscribeToMarket() {
	fmt.Printf("SubscribeToMarket: wait %s for %d incoming notifications\n",
		suite.SubscribeToMarketDuration, suite.SubscribeToMarketMsgs)

	if err := suite.TestAPI.SubscribeToMarket(suite.SubscribeToMarketSubscriberID,
		suite.AssetBTS, suite.AssetCNY); err != nil {
		suite.Fail(err.Error(), "SubscribeToMarket")
	}

	bar := pb.StartNew(suite.SubscribeToMarketMsgs).Prefix("wait for data")
	err := suite.TestAPI.OnNotify(suite.SubscribeToMarketSubscriberID, func(msg interface{}) error {
		bar.Increment()
		return nil
	})

	if err != nil {
		suite.Fail(err.Error(), "OnNotify")
	}

	suite.Condition(func() bool {
		return util.WaitForCondition(suite.SubscribeToMarketDuration, func() bool {
			return int(bar.Get()) >= suite.SubscribeToMarketMsgs
		})
	}, "not enough incomming notifications during test window")

	bar.Finish()

	fmt.Println("Call UnsubscribeFromMarket")
	if err := suite.TestAPI.UnsubscribeFromMarket(suite.AssetBTS, suite.AssetCNY); err != nil {
		suite.Fail(err.Error(), "UnsubscribeFromMarket")
	}
}

func (suite *bitsharesAPITest) Test_SetSubscribeCallback() {

	fmt.Printf("SetSubscribeCallback: wait %s for %d incoming notifications\n",
		suite.SetSubscribeCallbackDuration, suite.SetSubscribeCallbackMsgs)
	if err := suite.TestAPI.SetSubscribeCallback(suite.SetSubscribeCallbackSubscriberID, false); err != nil {
		suite.Fail(err.Error(), "SetSubscribeCallback")
	}

	_, err := suite.TestAPI.CallAPI(suite.TestAPI.DatabaseApiID(), "get_objects", []interface{}{"2.1.0"})
	if err != nil {
		suite.Fail(err.Error(), "CallAPI->get_objects")
	}

	bar := pb.StartNew(suite.SetSubscribeCallbackMsgs).Prefix("wait for data")
	err = suite.TestAPI.OnNotify(suite.SetSubscribeCallbackSubscriberID, func(msg interface{}) error {
		bar.Increment()
		return nil
	})

	if err != nil {
		log.Fatal(errors.Annotate(err, "OnNotify"))
	}

	suite.Condition(func() bool {
		return util.WaitForCondition(suite.SetSubscribeCallbackDuration, func() bool {
			return int(bar.Get()) >= suite.SetSubscribeCallbackMsgs
		})
	}, "not enough incomming notifications during test window")

	bar.Finish()

	fmt.Println("Call CancelAllSubscriptions")
	if err := suite.TestAPI.CancelAllSubscriptions(); err != nil {
		suite.Fail(err.Error(), "CancelAllSubscriptions")
	}
}

func TestSubscribe(t *testing.T) {
	testSuite := new(bitsharesAPITest)
	suite.Run(t, testSuite)
	testSuite.TearDown()
}

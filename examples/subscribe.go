package main

import (
	"log"
	"time"

	"github.com/denkhaus/bitshares/api"
	"github.com/juju/errors"
)

const (
	WebsocketURL = "wss://bitshares.openledger.info/ws"
	SubscriberID = 4
)

func main() {

	api := api.New(WebsocketURL)
	if err := api.Connect(); err != nil {
		log.Fatal(errors.Annotate(err, "connect"))
	}

	defer api.Close()

	if err := api.SetSubscribeCallback(SubscriberID, false); err != nil {
		log.Fatal(errors.Annotate(err, "SetSubscribeCallback"))
	}

	_, err := api.CallAPI(2, "get_objects", []interface{}{"2.1.0"})
	if err != nil {
		log.Fatal(errors.Annotate(err, "get_objects"))
	}

	err = api.OnNotify(SubscriberID, func(msg interface{}) error {
		log.Println("on notify", msg)
		return nil
	})

	if err != nil {
		log.Fatal(errors.Annotate(err, "on incomming message"))
	}

	time.Sleep(30 * time.Second)
}

package main

import (
	"log"
	"time"

	"github.com/denkhaus/bitshares/api"
	"github.com/juju/errors"
)

const (
	WebsocketURL   = "wss://bitshares.openledger.info/ws"
	SubscrNotifyID = 4
)

func main() {

	api := api.New(WebsocketURL)

	if err := api.Connect(); err != nil {
		log.Fatal(errors.Annotate(err, "connect"))
	}

	defer api.Close()

	if err := api.SetSubscribeCallback(SubscrNotifyID, false); err != nil {
		log.Fatal(errors.Annotate(err, "SetSubscribeCallback"))
	}

	_, err := api.Call(2, "get_objects", []interface{}{"2.1.0"})
	if err != nil {
		log.Fatal(errors.Annotate(err, "get_objects"))
	}

	time.Sleep(30 * time.Second)
}

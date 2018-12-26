package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/denkhaus/bitshares/config"
	"github.com/denkhaus/bitshares/types"
	"github.com/juju/errors"
)

const (
	memoSrc       = "My secret memo works!"
	pubKeyAString = "BTS5zzvbDtkbUVU1gFFsKqCE55U7JbjTp6mTh1usFv7KGgXL7HDQk"
	privKeyAWif   = "5Hx8KiHLnc3pDLkwe2jujkTTJev72n3Qx7xtyaRNBsJDuejzh9u"
	pubKeyBString = "BTS5Z3vsgH6xj6HLXcsU38yo4TyoZs9AUzpfbaXbuxsAYPbutWvEP"
	privKeyBWif   = "5KRZv3ZmkcE71K9KwEKG6pV6pyufkMQgCJrCu8xKLf2y7R7J8YK"
)

func main() {
	config.SetCurrentConfig(config.ChainIDBTS)

	pubKeyA, err := types.NewPublicKeyFromString(pubKeyAString)
	if err != nil {
		log.Fatal(errors.Annotate(err, "NewPublicKeyFromString [key A]"))
	}

	pubKeyB, err := types.NewPublicKeyFromString(pubKeyBString)
	if err != nil {
		log.Fatal(errors.Annotate(err, "NewPublicKeyFromString [key B]"))
	}

	memo := types.Memo{
		From:  *pubKeyA,
		To:    *pubKeyB,
		Nonce: types.UInt64(rand.Int63()),
	}

	privKeyA, err := types.NewPrivateKeyFromWif(privKeyAWif)
	if err != nil {
		log.Fatal(errors.Annotate(err, "NewPrivateKeyFromWif [key A]"))
	}

	if err := memo.Encrypt(privKeyA, memoSrc); err != nil {
		log.Fatal(errors.Annotate(err, "Encrypt"))
	}

	privKeyB, err := types.NewPrivateKeyFromWif(privKeyBWif)
	if err != nil {
		log.Fatal(errors.Annotate(err, "NewPrivateKeyFromWif [key B]"))
	}

	memoDst, err := memo.Decrypt(privKeyB)
	if err != nil {
		log.Fatal(errors.Annotate(err, "Encrypt"))
	}

	if memoSrc != memoDst {
		log.Fatalf("decryption error: memo is %q", memoDst)
	}

	fmt.Printf("decrypted memo is %q\n", memoDst)

}

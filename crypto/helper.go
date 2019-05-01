package crypto

import (
	"github.com/denkhaus/bitshares/config"
	"github.com/denkhaus/bitshares/types"
	"github.com/juju/errors"
)

//SignWithKeys signs a given transaction with given private keys.
func SignWithKeys(keys types.PrivateKeys, tx *types.SignedTransaction) error {
	signer := NewTransactionSigner(tx)
	if err := signer.Sign(keys, config.Current()); err != nil {
		return errors.Annotate(err, "Sign")
	}

	return nil
}

//VerifySignedTransaction verifies a signed transaction against all available keys in keyBag.
//If all required keys are found the function returns true, otherwise false.
func VerifySignedTransaction(keyBag *KeyBag, tx *types.SignedTransaction) (bool, error) {
	signer := NewTransactionSigner(tx)
	verified, err := signer.Verify(keyBag, config.Current())
	if err != nil {
		return false, errors.Annotate(err, "Verify")
	}

	return verified, nil
}

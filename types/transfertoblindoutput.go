package types

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type StealthConfirmation struct {

	//    struct memo_data
	//    {
	// 	  optional<public_key_type> from;
	// 	  asset                     amount;
	// 	  fc::sha256                blinding_factor;
	// 	  fc::ecc::commitment_type  commitment;
	// 	  uint32_t                  check = 0;
	//    };
	//  (from)(amount)(blinding_factor)(commitment)(check)

	OneTimeKey    PublicKey  `json:"one_time_key"`
	To            *PublicKey `json:"to,omitempty"`
	EncryptedMemo Buffer     `json:"encrypted_memo"`
}

func (p StealthConfirmation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.OneTimeKey); err != nil {
		return errors.Annotate(err, "encode OneTimeKey")
	}
	if err := enc.Encode(p.To != nil); err != nil {
		return errors.Annotate(err, "encode has To")
	}
	if err := enc.Encode(p.To); err != nil {
		return errors.Annotate(err, "encode To")
	}
	if err := enc.Encode(p.EncryptedMemo); err != nil {
		return errors.Annotate(err, "encode EncryptedMemo")
	}
	return nil
}

type TransferToBlindOutputs []TransferToBlindOutput

func (p TransferToBlindOutputs) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return errors.Annotate(err, "encode length")
	}

	for _, o := range p {
		if err := enc.Encode(o); err != nil {
			return errors.Annotate(err, "encode Output")
		}
	}

	return nil
}

type TransferToBlindOutput struct {
	Commitment          Buffer               `json:"commitment"`
	Owner               Authority            `json:"owner"`
	RangeProof          Buffer               `json:"range_proof"`
	StealthConfirmation *StealthConfirmation `json:"stealth_memo,omitempty"`
}

func (p TransferToBlindOutput) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.Commitment); err != nil {
		return errors.Annotate(err, "encode Commitment")
	}
	if err := enc.Encode(p.RangeProof); err != nil {
		return errors.Annotate(err, "encode RangeProof")
	}
	if err := enc.Encode(p.Owner); err != nil {
		return errors.Annotate(err, "encode Owner")
	}
	if err := enc.Encode(p.StealthConfirmation != nil); err != nil {
		return errors.Annotate(err, "encode has StealthConfirmation")
	}
	if err := enc.Encode(p.StealthConfirmation); err != nil {
		return errors.Annotate(err, "encode StealthConfirmation")
	}

	return nil
}

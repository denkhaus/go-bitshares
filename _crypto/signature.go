package crypto

import (
	"bytes"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/asuleymanov/golos-go/rfc6979"
	"github.com/btcsuite/btcd/btcec"
	"github.com/juju/errors"
)

type Signature struct {
	btcec.Signature
	i byte
}

func NewSignature(r *big.Int, s *big.Int, i byte) *Signature {
	sig := Signature{
		i: i,
		Signature: btcec.Signature{
			R: r,
			S: s,
		},
	}

	return &sig
}

func (p *Signature) ToHex() string {
	return hex.EncodeToString(p.Bytes())
}

func (p *Signature) Bytes() []byte {
	buf := bytes.Buffer{}
	buf.WriteByte(p.i)
	buf.Write(p.R.Bytes())
	buf.Write(p.S.Bytes())

	return buf.Bytes()
}

func (p *Signature) IsCanonical() bool {
	sig := p.Bytes()

	if ((sig[0] & 0x80) != 0) || (sig[0] == 0) ||
		((sig[1] & 0x80) != 0) ||
		((sig[32] & 0x80) != 0) || (sig[32] == 0) ||
		((sig[33] & 0x80) != 0) {
		return false
	}

	return true
}

func Sign(data []byte, privKey *btcec.PrivateKey) (*Signature, error) {
	digest := sha256.Sum256(data)
	return signBufferSha256(digest[:], privKey)
}

func recoverPubKey(curve *btcec.KoblitzCurve, sig *btcec.Signature, msg []byte, iter int, doChecks bool) (*btcec.PublicKey, error) {
	// 1.1 x = (n * i) + r
	Rx := new(big.Int).Mul(curve.Params().N,
		new(big.Int).SetInt64(int64(iter/2)))
	Rx.Add(Rx, sig.R)
	if Rx.Cmp(curve.Params().P) != -1 {
		return nil, errors.New("calculated Rx is larger than curve P")
	}

	// convert 02<Rx> to point R. (step 1.2 and 1.3). If we are on an odd
	// iteration then 1.6 will be done with -R, so we calculate the other
	// term when uncompressing the point.
	Ry, err := decompressPoint(curve, Rx, iter%2 == 1)
	if err != nil {
		return nil, err
	}

	// 1.4 Check n*R is point at infinity
	if doChecks {
		nRx, nRy := curve.ScalarMult(Rx, Ry, curve.Params().N.Bytes())
		if nRx.Sign() != 0 || nRy.Sign() != 0 {
			return nil, errors.New("n*R does not equal the point at infinity")
		}
	}

	// 1.5 calculate e from message using the same algorithm as ecdsa
	// signature calculation.
	e := hashToInt(msg, curve)

	// Step 1.6.1:
	// We calculate the two terms sR and eG separately multiplied by the
	// inverse of r (from the signature). We then add them to calculate
	// Q = r^-1(sR-eG)
	invr := new(big.Int).ModInverse(sig.R, curve.Params().N)

	// first term.
	invrS := new(big.Int).Mul(invr, sig.S)
	invrS.Mod(invrS, curve.Params().N)
	sRx, sRy := curve.ScalarMult(Rx, Ry, invrS.Bytes())

	// second term.
	e.Neg(e)
	e.Mod(e, curve.Params().N)
	e.Mul(e, invr)
	e.Mod(e, curve.Params().N)
	minuseGx, minuseGy := curve.ScalarBaseMult(e.Bytes())

	// TODO: this would be faster if we did a mult and add in one
	// step to prevent the jacobian conversion back and forth.
	Qx, Qy := curve.Add(sRx, sRy, minuseGx, minuseGy)

	return &btcec.PublicKey{
		Curve: curve,
		X:     Qx,
		Y:     Qy,
	}, nil
}

func calcPubKeyRecoveryParam(curve *btcec.KoblitzCurve, buf_sha256 []byte, ecsignature *btcec.Signature, q *btcec.PublicKey) (int, error) {
	for i := 0; i < 4; i++ {
		pubKeyRecov, err := recoverPubKey(curve, ecsignature, buf_sha256, i, true)
		if err != nil {
			return -1, errors.Annotate(err, "recoverPubKey")
		}

		pub := pubKeyRecov.SerializeCompressed()
		digest := sha256.Sum256(pub)
		fmt.Println("recovered pubkey", hex.EncodeToString(digest[:]))
		// 1.6.2 Verify Q
		if pubKeyRecov.IsEqual(q) {
			return i, nil
		}
	}

	return -1, errors.New("Unable to find valid recovery factor")
}

func signBufferSha256(buf_sha256 []byte, private_key *btcec.PrivateKey) (*Signature, error) {
	if len(buf_sha256) != 32 {
		return nil, errors.New("buf_sha256: 32 byte buffer required")
	}

	privKeyECDSA := private_key.ToECDSA()
	nonce := 0

	for {
		nonce++
		r, s, err := rfc6979.SignECDSA(privKeyECDSA, buf_sha256, sha256.New, nonce)
		if err != nil {
			return nil, errors.Annotate(err, "SignECDSA")
		}

		sig := &btcec.Signature{R: r, S: s}
		der := sig.Serialize()
		lenR := der[3]
		lenS := der[5+lenR]

		if lenR == 32 && lenS == 32 {
			i, err := calcPubKeyRecoveryParam(btcec.S256(), buf_sha256, sig, private_key.PubKey())
			if err != nil {
				return nil, errors.Annotate(err, "calcPubKeyRecoveryParam")
			}

			i += 4  // compressed
			i += 27 // compact  //  24 or 27 :( forcing odd-y 2nd key candidate)

			return NewSignature(sig.R, sig.S, byte(i)), nil
		}

		if nonce%10 == 0 {
			fmt.Print("WARN: ", nonce, " attempts to find canonical signature")
		}
	}
}

func decompressPoint(curve *btcec.KoblitzCurve, x *big.Int, ybit bool) (*big.Int, error) {
	// TODO: This will probably only work for secp256k1 due to
	// optimizations.

	// Y = +-sqrt(x^3 + B)
	x3 := new(big.Int).Mul(x, x)
	x3.Mul(x3, x)
	x3.Add(x3, curve.Params().B)

	// now calculate sqrt mod p of x2 + B
	// This code used to do a full sqrt based on tonelli/shanks,
	// but this was replaced by the algorithms referenced in
	// https://bitcointalk.org/index.php?topic=162805.msg1712294#msg1712294
	y := new(big.Int).Exp(x3, curve.QPlus1Div4(), curve.Params().P)

	if ybit != isOdd(y) {
		y.Sub(curve.Params().P, y)
	}
	if ybit != isOdd(y) {
		return nil, fmt.Errorf("ybit doesn't match oddness")
	}
	return y, nil
}

func isOdd(a *big.Int) bool {
	return a.Bit(0) == 1
}

func hashToInt(hash []byte, c elliptic.Curve) *big.Int {
	orderBits := c.Params().N.BitLen()
	orderBytes := (orderBits + 7) / 8
	if len(hash) > orderBytes {
		hash = hash[:orderBytes]
	}

	ret := new(big.Int).SetBytes(hash)
	excess := len(hash)*8 - orderBits
	if excess > 0 {
		ret.Rsh(ret, uint(excess))
	}
	return ret
}

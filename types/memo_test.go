package types

import (
	"crypto/ecdsa"
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/denkhaus/bitshares/util"
	"github.com/stretchr/testify/assert"
)

var sharedSecrets = [][]interface{}{
	[]interface{}{"5JYWCqDpeVrefVaFxJfDc3mzQ67dtsfhU7zcB7AMJYuTH57VsoE", "GPH56EzLTXkis55hBsompVXmSdnayG3afDNFmsCLohPh6rSNzkzhs", "fbb2fef5a3a115887df84c694e8ac5c9bf998c89d0c22438c18fd018f2529460"},
	[]interface{}{"5JKhu9ZKydGFz7yGURocDVEepSY9fk2VRGAA8Xnb9wwFWa8yTWy", "GPH818iy2auxecLxhWTtW219w2VAfBYHxeHaeRASoTFLnsZo1DJ63", "4a52093355abeb31cef02ee1cbdf0661d982d52ad8fe39c68957e3ae03f3bda9"},
	[]interface{}{"5KKmTkFCNnedj6hbyRYJwcaMnc4TkuwrPsJDqR2Bj9ShHkfdgQ3", "GPH78SdnBpqhEHxxzwZeKoFEXV6PviymWzBF7ev29pZcTCF8ynJAo", "500e67a07f53d49b88db635c64e4b0a2414168c7054118d40001e86f1abce131"},
	[]interface{}{"5KLBuZtagfmGqhDTEPSM84TXKxKfzNyGaxRKCgdcocEU7Nusw49", "GPH5nYv9AusGXgHyMBbSBV4HyEAmhzXqLNRPvUpKmNpFo5soho95o", "febfa7ad6c48bb0ab976c6416da24017b93a58e4e699dba76fc590b4b1ac0d26"},
	[]interface{}{"5JrVxMdeBZJvWqV4SmyFq9psQ4Dg8cFXtSWDiL7V5gUJC133xC2", "GPH7HxVNixmh33R44Kr2uJERbhvzkaLen8su4juqyFe2FW2U2cCXA", "82ef43913f83dd3ff0b4f06bcd8801a06c9f046b44b054e0a9ad042c28e5bdba"},
	[]interface{}{"5JfEonXJ4H2kSP4V9NzC3uTRtTpLx4wVgDvf5AWN1KKTV6CZ4x7", "GPH6ZtaoP6skA433YGNNJcPGnsgx15psKRBwAy83tw7XWsDy8hso3", "b1ad058e9cc48e305fb46f07736409a55692c67d3507aad6a051b35459ec2f93"},
	[]interface{}{"5J5UDLdk9XjvcbzNY5AQoUB2pttsvN7FtQFyyFZXUUsHFAp9iQd", "GPH51wPrJXWLcX6iNPAoZ9sGk4fHXk6krQgTX1jfuyxtKuhoEan83", "82fcc73de1331913945f6ce6d0207864bbc7cffee10ed3533ab32629cb759323"},
	[]interface{}{"5JNZpagkR8wWsW3n4hHqFUQVkAu3HhJ9kU1criuruFpAwoaesHs", "GPH5SCy1teB91pNYetxEwV8vRyMApsy8aG61wsi8z2B4Zb6kfnqUf", "704097d0c270e93f0ce5fa91049bb0aa2f38ccfc4bdc38840176abbb98337c0c"},
	[]interface{}{"5KKBRfgTgATU5SmF2uy7ewi7BbDDJCtmf3x9CeYziF14uj8YHMM", "GPH7pUa1fp4NtGaRDmZF6TeanHw7zELUp1eWxZasRE3zY4xYKdbhV", "114aba4ab84ea225bbf4b60aaf6d467d3b206ff8a94d531a5a6031ad90c874dd"},
	[]interface{}{"5Jg7muALcVxncN32LyGMDK8zut2b1Sw3VJA1xjZE5ght7DRM9ac", "GPH5Vj6uR2iKmrB2DcFyqNzperycD3a32BBYkefzKYCHoGnXemwWS", "60928672da8e9a7dc0f783f2bf8aaf1b206b9bbd85f0a61b638e0b99f5f8ea56"},
}

// MaxSharedKeyLength returns the maximum length of the shared key the
// public key can produce.
func MaxSharedKeyLength(pub *ecdsa.PublicKey) int {
	return (pub.Curve.Params().BitSize + 7) / 8
}

func sharedSecret(priv *btcec.PrivateKey, pub *ecdsa.PublicKey, skLen, macLen int) (sk []byte, err error) {
	if priv.PublicKey.Curve != pub.Curve {
		return nil, ErrInvalidCurve
	}

	if skLen+macLen > MaxSharedKeyLength(pub) {
		return nil, ErrSharedKeyTooBig
	}

	x, _ := pub.Curve.ScalarMult(pub.X, pub.Y, priv.D.Bytes())
	if x == nil {
		return nil, ErrSharedKeyIsPointAtInfinity
	}

	sk = make([]byte, skLen+macLen)
	skBytes := x.Bytes()
	copy(sk[len(sk)-len(skBytes):], skBytes)
	return sk, nil
}

func Test_SharedSecrets(t *testing.T) {
	for _, tst := range sharedSecrets {
		priv, err := util.GetPrivateKey(tst[0].(string))
		if err != nil {
			assert.FailNow(t, err.Error(), "GetPrivateKey")
		}

		pub, err := NewPublicKey(tst[1].(string))
		if err != nil {
			assert.FailNow(t, err.Error(), "GetPublicKey")
		}

		sec2, err := sharedSecret(priv, pub.ToECDSA(), 16, 16)
		if err != nil {
			assert.FailNow(t, err.Error(), "sharedSecret")
		}

		sec1, err := hex.DecodeString(tst[2].(string))
		if err != nil {
			assert.FailNow(t, err.Error(), "DecodeString")
		}

		assert.Equal(t, sec1, sec2)
	}
}

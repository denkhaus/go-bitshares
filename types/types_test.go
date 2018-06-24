package types

import (
	"bytes"
	"testing"

	"github.com/denkhaus/bitshares/util"
	"github.com/stretchr/testify/assert"
)

func TestBuffer_MarshalUnmarshal(t *testing.T) {

	b1 := Buffer("TestMarshalUnmarshal")

	var buf bytes.Buffer
	enc := util.NewTypeEncoder(&buf)

	if err := b1.Marshal(enc); err != nil {
		assert.FailNow(t, err.Error(), "Marshal")
	}

	dec := util.NewTypeDecoder(&buf)
	var b2 Buffer
	if err := b2.Unmarshal(dec); err != nil {
		assert.FailNow(t, err.Error(), "Unmarshal")
	}

	assert.Equal(t, b1.Bytes(), b2.Bytes())

}

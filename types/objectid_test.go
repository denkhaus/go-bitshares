package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ObjectID(t *testing.T) {
	id := NewObjectID("1.5.29874")
	assert.Equal(t, id.SpaceType(), SpaceTypeProtocol)
	assert.Equal(t, id.ObjectType(), ObjectTypeCommitteeMember)
	assert.Equal(t, id.Instance(), UInt64(29874))
}

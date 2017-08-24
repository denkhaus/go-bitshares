package objects

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type GrapheneObject interface {
	util.TypeMarshaller
	Id() ObjectID
	Type() ObjectType
}

type GrapheneObjects []GrapheneObject

func (p GrapheneObjects) ToObjectIDs() []ObjectID {
	ids := []ObjectID{}
	for _, o := range p {
		ids = append(ids, o.Id())
	}
	return ids
}

type GrapheneID struct {
	id         ObjectID
	spaceType  SpaceType
	objectType ObjectType
	instance   UInt64
}

//implements TypeMarshaller interface
func (p GrapheneID) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(p.instance)); err != nil {
		return errors.Annotate(err, "encode GrapheneID instance")
	}

	return nil
}

func (p GrapheneID) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.id)
}

func (p *GrapheneID) UnmarshalJSON(s []byte) error {
	str := string(s)

	if len(str) > 0 && str != "null" {
		q, err := util.SafeUnquote(str)
		if err != nil {
			return errors.Annotate(err, "SafeUnquote")
		}

		if err := p.FromString(q); err != nil {
			return errors.Annotate(err, "FromString")
		}

		return nil
	}

	return errors.Errorf("unable to unmarshal GrapheneID from %s", str)
}

//Id returns the objects ObjectID
func (p GrapheneID) Id() ObjectID {
	return p.id
}

//Type returns the objects ObjectType
func (p GrapheneID) Type() ObjectType {
	if !p.valid() {
		if err := p.FromString(string(p.id)); err != nil {
			panic(errors.Annotate(err, "from string").Error())
		}
	}
	return p.objectType
}

//NewGrapheneID creates an new GrapheneID object
func NewGrapheneID(id ObjectID) *GrapheneID {
	gid := &GrapheneID{
		spaceType:  SpaceTypeUndefined,
		objectType: ObjectTypeUndefined,
	}

	if err := gid.FromString(string(id)); err != nil {
		panic(err.Error())
	}

	return gid
}

func (p GrapheneID) String() string {
	return string(p.Id())
}

func (p GrapheneID) valid() bool {
	return p.spaceType != SpaceTypeUndefined &&
		p.objectType != ObjectTypeUndefined
}

func (p *GrapheneID) FromRawData(in interface{}) error {
	o, ok := in.(map[string]interface{})
	if !ok {
		return errors.New("input is not map[string]interface{}")
	}

	if id, ok := o["id"]; ok {
		return p.FromString(id.(string))
	}

	return errors.New("input is no graphene object")
}

func (p *GrapheneID) FromObjectID(in ObjectID) error {
	return p.FromString(string(in))
}

func (p *GrapheneID) FromString(in string) error {
	parts := strings.Split(in, ".")

	if len(parts) == 3 {
		p.id = ObjectID(in)
		space, err := strconv.Atoi(parts[0])
		if err != nil {
			return errors.Errorf("unable to parse GrapheneID [space] from %s", in)
		}

		p.spaceType = SpaceType(space)

		typ, err := strconv.Atoi(parts[1])
		if err != nil {
			return errors.Errorf("unable to parse GrapheneID [type] from %s", in)
		}

		p.objectType = ObjectType(typ)

		inst, err := strconv.ParseUint(parts[2], 10, 64)
		if err != nil {
			return errors.Errorf("unable to parse GrapheneID [instance] from %s", in)
		}

		p.instance = UInt64(inst)
		return nil
	}

	return errors.Errorf("unable to parse GrapheneID from %s", in)
}

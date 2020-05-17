package right

import (
	"github.com/likecoin/iscn-ipld/plugin/block"
	"github.com/likecoin/iscn-ipld/plugin/block/time_period"
)

const (
	// SchemaName of right
	SchemaName = "right"
)

// Register registers the schema of right block
func Register() {
	block.RegisterIscnObjectFactory(
		block.CodecRight,
		SchemaName,
		[]block.CodecFactoryFunc{
			newSchemaV1,
		},
	)
}

// ==================================================
// base
// ==================================================

// base is the base struct for right (codec 0x02BD)
type base struct {
	*block.Base
}

func newBase(version uint64, schema []block.Data) (*base, error) {
	blockBase, err := block.NewBase(
		block.CodecRight,
		SchemaName,
		version,
		schema,
	)
	if err != nil {
		return nil, err
	}

	return &base{
		Base: blockBase,
	}, nil
}

// ==================================================
// schemaV1
// ==================================================

// schemaV1 represents a right V1
type schemaV1 struct {
	*base
}

var _ block.IscnObject = (*schemaV1)(nil)

func newSchemaV1() (block.Codec, error) {
	schema := []block.Data{
		block.NewCid("holder", true, block.CodecEntity),
		block.NewString("type", true), // TODO use filterd string??
		block.NewCid("terms", true, 0),
		block.NewObject("period", false, timeperiod.SchemaV1Prototype),
		block.NewString("territory", false),
	}

	timePeriodBase, err := newBase(1, schema)
	if err != nil {
		return nil, err
	}

	return &schemaV1{
		base: timePeriodBase,
	}, nil
}

// SchemaV1Prototype creates a prototype for schemaV1
func SchemaV1Prototype() block.Codec {
	res, _ := newSchemaV1()
	return res
}

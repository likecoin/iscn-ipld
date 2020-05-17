package stakeholder

import (
	"github.com/likecoin/iscn-ipld/plugin/block"
)

const (
	// SchemaName of stakeholder
	SchemaName = "stakeholder"
)

// Register registers the schema of stakeholder block
func Register() {
	block.RegisterIscnObjectFactory(
		block.CodecStakeholder,
		SchemaName,
		[]block.CodecFactoryFunc{
			newSchemaV1,
		},
	)
}

// ==================================================
// base
// ==================================================

// base is the base struct for stakeholder (codec 0x02D1)
type base struct {
	*block.Base
}

func newBase(version uint64, schema []block.Data) (*base, error) {
	blockBase, err := block.NewBase(
		block.CodecStakeholder,
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

// schemaV1 represents an stakeholder V1
type schemaV1 struct {
	*base
}

var _ block.IscnObject = (*schemaV1)(nil)

func newSchemaV1() (block.Codec, error) {
	ty := NewType()

	schema := []block.Data{
		ty,
		block.NewCid("stakeholder", true, block.CodecEntity),
		block.NewNumber("sharing", true, block.Uint32T),
		NewFootprint(ty),
	}

	stakeholderBase, err := newBase(1, schema)
	if err != nil {
		return nil, err
	}

	return &schemaV1{
		base: stakeholderBase,
	}, nil
}

// SchemaV1Prototype creates a prototype for schemaV1
func SchemaV1Prototype() block.Codec {
	res, _ := newSchemaV1()
	return res
}

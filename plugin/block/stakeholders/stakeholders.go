package stakeholders

import (
	"github.com/likecoin/iscn-ipld/plugin/block"
	"github.com/likecoin/iscn-ipld/plugin/block/stakeholder"
)

const (
	// SchemaName of stakeholders
	SchemaName = "stakeholders"
)

// Register registers the schema of stakeholders block
func Register() {
	block.RegisterIscnObjectFactory(
		block.CodecStakeholders,
		SchemaName,
		[]block.CodecFactoryFunc{
			newSchemaV1,
		},
	)
}

// ==================================================
// base
// ==================================================

// base is the base struct for stakeholders (codec 0x0266)
type base struct {
	*block.Base
}

func newBase(version uint64, schema []block.Data) (*base, error) {
	blockBase, err := block.NewBase(
		block.CodecStakeholders,
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

// schemaV1 represents an stakeholders V1
type schemaV1 struct {
	*base
}

var _ block.IscnObject = (*schemaV1)(nil)

func newSchemaV1() (block.Codec, error) {
	prototype := block.NewObject("_", true, stakeholder.SchemaV1Prototype)

	schema := []block.Data{
		block.NewDataArray("stakeholders", true, prototype),
	}

	stakeholdersBase, err := newBase(1, schema)
	if err != nil {
		return nil, err
	}

	return &schemaV1{
		base: stakeholdersBase,
	}, nil
}

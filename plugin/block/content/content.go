package content

import (
	"github.com/likecoin/iscn-ipld/plugin/block"
)

const (
	// SchemaName of content
	SchemaName = "content"
)

// Register registers the schema of content block
func Register() {
	block.RegisterIscnObjectFactory(
		block.CodecContent,
		SchemaName,
		[]block.CodecFactoryFunc{
			newSchemaV1,
		},
	)
}

// ==================================================
// base
// ==================================================

// base is the base struct for content (codec 0x0267)
type base struct {
	*block.Base
}

func newBase(version uint64, schema []block.Data) (*base, error) {
	blockBase, err := block.NewBase(
		block.CodecContent,
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

// schemaV1 represents a content V1
type schemaV1 struct {
	*base

	version *block.Number
	parent  *block.Cid
}

var _ block.IscnObject = (*schemaV1)(nil)

func newSchemaV1() (block.Codec, error) {
	version := block.NewNumber("version", true, block.Uint64T)
	parent := block.NewCid("parent", false, block.CodecContent)

	schema := []block.Data{
		block.NewString("type", true),
		version,
		parent,
		block.NewString("source", false), // TODO URL
		block.NewString("edition", false),
		block.NewString("fingerprint", true), // TODO HashURL
		block.NewString("title", true),
		block.NewString("description", false),
		block.NewDataArray("tags", false, block.NewString("_", false)),
	}

	contentBase, err := newBase(1, schema)
	if err != nil {
		return nil, err
	}

	obj := schemaV1{
		base:    contentBase,
		version: version,
		parent:  parent,
	}
	contentBase.SetValidator(obj.Validate)

	return &obj, nil
}

// Validate the data
func (o *schemaV1) Validate() error {
	return block.ValidateParent(o.version, o.parent)
}

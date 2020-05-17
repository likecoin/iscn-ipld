package kernel

import (
	"fmt"

	"github.com/likecoin/iscn-ipld/plugin/block"
)

const (
	// SchemaName of ISCN kernel
	SchemaName = "iscn"
)

// Register registers the schema of ISCN kernel block
func Register() {
	block.RegisterIscnObjectFactory(
		block.CodecISCN,
		SchemaName,
		[]block.CodecFactoryFunc{
			newSchemaV1,
		},
	)
}

// ==================================================
// base
// ==================================================

// base is the base struct for ISCN kernel (codec 0x0264)
type base struct {
	*block.Base

	id *ID
}

func newBase(version uint64, schema []block.Data, id *ID) (*base, error) {
	blockBase, err := block.NewBase(
		block.CodecISCN,
		SchemaName,
		version,
		schema,
	)
	if err != nil {
		return nil, err
	}

	return &base{
		Base: blockBase,
		id:   id,
	}, nil
}

// github.com/ipfs/go-block-format.Block interface

// Loggable returns a map the type of IPLD Link
func (b *base) Loggable() map[string]interface{} {
	l := b.Base.Loggable()
	l["id"] = b.id.GetID()
	return l
}

// String is a helper for output
func (b *base) String() string {
	return fmt.Sprintf("<%s (v%d): %s>", b.GetName(), b.GetVersion(), b.id.GetID())
}

// ==================================================
// schemaV1
// ==================================================

// schemaV1 represents an ISCN kernel V1
type schemaV1 struct {
	*base

	version *block.Number
	parent  *block.Cid
}

var _ block.IscnObject = (*schemaV1)(nil)

func newSchemaV1() (block.Codec, error) {
	id := NewID()
	version := block.NewNumber("version", true, block.Uint64T)
	parent := block.NewCid("parent", false, block.CodecISCN)

	schema := []block.Data{
		id,
		block.NewTimestamp("timestamp", true),
		version,
		parent,
		block.NewCid("stakeholders", true, block.CodecStakeholders),
		block.NewCid("content", true, block.CodecContent),
	}

	iscnKernelBase, err := newBase(1, schema, id)
	if err != nil {
		return nil, err
	}

	obj := schemaV1{
		base:    iscnKernelBase,
		version: version,
		parent:  parent,
	}
	iscnKernelBase.SetValidator(obj.Validate)

	return &obj, nil
}

// Validate the data
func (o *schemaV1) Validate() error {
	return block.ValidateParent(o.version, o.parent)
}

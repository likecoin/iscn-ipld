package entity

import (
	"github.com/ipfs/go-cid"

	blocks "github.com/ipfs/go-block-format"
	node "github.com/ipfs/go-ipld-format"
	iscn "github.com/likecoin/iscn-ipld/plugin/block"
)

const (
	// SchemaName of entity
	SchemaName = "entity"
)

// Register registers the schema of entity block
func Register() {
	iscn.RegisterIscnObjectFactory(
		iscn.CodecEntity,
		SchemaName,
		[]iscn.CodecFactoryFunc{
			newSchemaV1,
		},
	)
}

// ==================================================
// base
// ==================================================

// base is the base struct for entity (codec 0x0267)
type base struct {
	*iscn.Base
}

func newBase(version uint64, schema []iscn.Data) (*base, error) {
	blockBase, err := iscn.NewBase(
		iscn.CodecEntity,
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

// schemaV1 represents an entity V1
type schemaV1 struct {
	*base
}

var _ iscn.IscnObject = (*schemaV1)(nil)

func newSchemaV1() (iscn.Codec, error) {
	schema := []iscn.Data{
		iscn.NewString("id", true), // TODO llc://id
		iscn.NewString("name", false),
		iscn.NewString("description", false),
	}

	entityBase, err := newBase(1, schema)
	if err != nil {
		return nil, err
	}

	return &schemaV1{
		base: entityBase,
	}, nil
}

// github.com/ipfs/go-ipld-format.DecodeBlockFunc

// BlockDecoder takes care of the entity IPLD objects
func BlockDecoder(block blocks.Block) (node.Node, error) {
	return iscn.Decode(block.RawData(), block.Cid())
}

// Package function

// DecodeData decodes the raw bytes to entity IPLD objects
func DecodeData(rawData []byte, c cid.Cid) (node.Node, error) {
	return iscn.Decode(rawData, c)
}

// NewEntityBlock creates an entity IPLD object
func NewEntityBlock(version uint64, data map[string]interface{}) (iscn.IscnObject, error) {
	return iscn.Encode(iscn.CodecEntity, version, data)
}

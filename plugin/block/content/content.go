package content

import (
	"github.com/ipfs/go-cid"

	blocks "github.com/ipfs/go-block-format"
	node "github.com/ipfs/go-ipld-format"
	iscn "github.com/likecoin/iscn-ipld/plugin/block"
)

const (
	// SchemaName of content
	SchemaName = "content"
)

// Register registers the schema of content block
func Register() {
	iscn.RegisterIscnObjectFactory(
		iscn.CodecContent,
		SchemaName,
		[]iscn.CodecFactoryFunc{
			newSchemaV1,
		},
	)
}

// ==================================================
// base
// ==================================================

// base is the base struct for content (codec 0x0777)
type base struct {
	*iscn.Base
}

func newBase(version uint64, schema []iscn.Data) (*base, error) {
	blockBase, err := iscn.NewBase(
		iscn.CodecContent,
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

// schemaV1 represents an content V1
type schemaV1 struct {
	*base
}

var _ iscn.IscnObject = (*schemaV1)(nil)

func newSchemaV1() (iscn.Codec, error) {
	version := iscn.NewNumber("version", true, iscn.Uint64T)
	schema := []iscn.Data{
		iscn.NewString("type", true),
		version,
		iscn.NewParent("parent", iscn.CodecContent, version),
		iscn.NewString("source", false), // TODO URL
		iscn.NewString("edition", false),
		iscn.NewString("fingerprint", true), // TODO HashURL
		iscn.NewString("title", true),
		iscn.NewString("description", false),
		iscn.NewDataArray("tags", false, iscn.NewString("_", false)),
	}

	contentBase, err := newBase(1, schema)
	if err != nil {
		return nil, err
	}

	return &schemaV1{
		base: contentBase,
	}, nil
}

// github.com/ipfs/go-ipld-format.DecodeBlockFunc

// BlockDecoder takes care of the content IPLD objects
func BlockDecoder(block blocks.Block) (node.Node, error) {
	return iscn.Decode(block.RawData(), block.Cid())
}

// Package function

// DecodeData decodes the raw bytes to content IPLD objects
func DecodeData(rawData []byte, c cid.Cid) (node.Node, error) {
	return iscn.Decode(rawData, c)
}

// NewContentBlock creates an content IPLD object
func NewContentBlock(version uint64, data map[string]interface{}) (iscn.IscnObject, error) {
	return iscn.Encode(iscn.CodecContent, version, data)
}

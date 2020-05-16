package kernel

import (
	"fmt"

	"github.com/ipfs/go-cid"

	blocks "github.com/ipfs/go-block-format"
	node "github.com/ipfs/go-ipld-format"
	iscn "github.com/likecoin/iscn-ipld/plugin/block"
)

const (
	// SchemaName of ISCN kernel
	SchemaName = "iscn"
)

//
// // Block (iscn-kernel, codec 0x0264), represents an ISCN kernel
// type Block struct {
// 	ID           []byte   `json:"id" yaml:"id"`
// 	Timestamp    string   `json:"timestamp" yaml:"timestamp"`
// 	Version      uint32   `json:"version" yaml:"version"`
// 	Parent       *cid.Cid `json:"parent" yaml:"parent"`
// 	Rights       cid.Cid  `json:"rights" yaml:"rights"`
// 	Stakeholders cid.Cid  `json:"stakeholders" yaml:"stakeholders"`
// 	Content      cid.Cid  `json:"content" yaml:"content"`
//
// 	cid     *cid.Cid
// 	rawdata []byte
// }

// Register registers the schema of ISCN kernel block
func Register() {
	iscn.RegisterIscnObjectFactory(
		iscn.CodecISCN,
		SchemaName,
		[]iscn.CodecFactoryFunc{
			newSchemaV1,
		},
	)
}

// ==================================================
// base
// ==================================================

// base is the base struct for ISCN kernel (codec 0x0264)
type base struct {
	*iscn.Base

	id *ID
}

func newBase(version uint64, schema []iscn.Data, id *ID) (*base, error) {
	blockBase, err := iscn.NewBase(
		iscn.CodecISCN,
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
}

var _ iscn.IscnObject = (*schemaV1)(nil)

func newSchemaV1() (iscn.Codec, error) {
	id := NewID()
	version := iscn.NewNumber("version", true, iscn.Uint64T)
	schema := []iscn.Data{
		id,
		iscn.NewTimestamp("timestamp", true),
		version,
		iscn.NewParent("parent", iscn.CodecISCN, version),
		iscn.NewCid("content", true, iscn.CodecContent),
	}

	iscnKernelBase, err := newBase(1, schema, id)
	if err != nil {
		return nil, err
	}

	return &schemaV1{
		base: iscnKernelBase,
	}, nil
}

// github.com/ipfs/go-ipld-format.DecodeBlockFunc

// BlockDecoder takes care of the ISCN kernel IPLD objects
func BlockDecoder(block blocks.Block) (node.Node, error) {
	return iscn.Decode(block.RawData(), block.Cid())
}

// Package function

// DecodeData decodes the raw bytes to ISCN kernel IPLD objects
func DecodeData(rawData []byte, c cid.Cid) (node.Node, error) {
	return iscn.Decode(rawData, c)
}

// NewIscnKernelBlock creates an ISCN kernel IPLD object
func NewIscnKernelBlock(version uint64, data map[string]interface{}) (iscn.IscnObject, error) {
	return iscn.Encode(iscn.CodecISCN, version, data)
}

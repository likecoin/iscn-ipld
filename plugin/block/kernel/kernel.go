package kernel

import (
	"fmt"

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

type factoryFunc []func() (iscn.IscnObject, error)

var factory factoryFunc = factoryFunc{
	newSchemaV1,
}

// ==================================================
// base
// ==================================================

// base is the base struct for ISCN kernel (iscn-kernel, codec 0x0264)
type base struct {
	*iscn.Base

	id string
}

func newBase(version uint64, schema []iscn.Data) (*base, error) {
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
	}, nil
}

// github.com/ipfs/go-block-format.Block interface

// Loggable returns a map the type of IPLD Link
func (b *base) Loggable() map[string]interface{} {
	l := b.Base.Loggable()
	l["id"] = b.id
	return l
}

// String is a helper for output
func (b *base) String() string {
	return fmt.Sprintf("<%s (v%d): %s>", b.GetName(), b.GetVersion(), b.id)
}

// ==================================================
// schemaV1
// ==================================================

// schemaV1 represents an ISCN kernel V1
type schemaV1 struct {
	*base
}

var _ iscn.IscnObject = (*schemaV1)(nil)

func newSchemaV1() (iscn.IscnObject, error) {
	schema := []iscn.Data{
		NewID(),
	}

	iscnKernelBase, err := newBase(1, schema)
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
	version, data, err := iscn.DecodeData(block)
	if err != nil {
		return nil, err
	}

	if version > (uint64)(len(factory)) {
		return nil, fmt.Errorf("<%s (v%d)> is not implemented", SchemaName, version)
	}
	version--

	obj, err := factory[version]()
	if err != nil {
		return nil, err
	}

	if err := obj.Decode(data); err != nil {
		return nil, err
	}

	return obj, nil
}

// Package function

// NewIscnKernelBlock creates an ISCN kernel IPLD object
func NewIscnKernelBlock(version uint64, data map[string]interface{}) (iscn.IscnObject, error) {
	if version > (uint64)(len(factory)) {
		return nil, fmt.Errorf("<%s (v%d)> is not implemented", SchemaName, version)
	}
	version--

	obj, err := factory[version]()
	if err != nil {
		return nil, err
	}

	if err := obj.SetData(data); err != nil {
		return nil, err
	}

	if err := obj.Encode(); err != nil {
		return nil, err
	}

	return obj, nil
}

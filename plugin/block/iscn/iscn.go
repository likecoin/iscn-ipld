package iscn

import (
	"log"

	goblocks "github.com/ipfs/go-block-format"
	cid "github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	node "github.com/ipfs/go-ipld-format"
	iscnblocks "github.com/likecoin/iscn-ipld/plugin/block"
	mh "github.com/multiformats/go-multihash"
)

// Block (iscn-block, codec 0x0264), represents an ISCN block header
type Block struct {
	Version int    `json:"version"`
	ID      string `json:"id"`
	Owner   string `json:"owner"`
	Edition int    `json:"edition"`
	Hash    string `json:"hash"`

	cid     *cid.Cid
	rawdata []byte
}

var _ node.Node = (*Block)(nil)

// github.com/ipfs/go-block-format.Block interface

// Cid returns the CID of the block header
func (b *Block) Cid() cid.Cid {
	log.Println("iscn-Cid")
	return *b.cid
}

// Loggable returns a map the type of IPLD Link
func (*Block) Loggable() map[string]interface{} {
	log.Println("iscn-Loggable")
	return map[string]interface{}{
		"type": "iscn-block",
	}
}

// RawData returns the binary of the CBOR encode of the block header
func (b *Block) RawData() []byte {
	log.Println("iscn-RawData")
	return b.rawdata
}

// String is a helper for output
func (*Block) String() string {
	log.Println("iscn-String")
	//TODO: return "<ISCN-block> `id`"
	return "iscn-block"
}

// node.Resolver interface

// Resolve resolves a path through this node, stopping at any link boundary
// and returning the object found as well as the remaining path to traverse
func (*Block) Resolve(path []string) (interface{}, []string, error) {
	log.Println("iscn-Resolve")
	return nil, nil, nil
}

// Tree lists all paths within the object under 'path', and up to the given depth.
// To list the entire object (similar to `find .`) pass "" and -1
func (*Block) Tree(path string, depth int) []string {
	log.Println("iscn-Tree")
	return nil
}

// node.Node interface

// Copy will go away. It is here to comply with the Node interface.
func (*Block) Copy() node.Node {
	log.Println("iscn-Copy")
	panic("dont use this yet")
}

// Links is a helper function that returns all links within this object
// HINT: Use `ipfs refs <cid>`
func (*Block) Links() []*node.Link {
	log.Println("iscn-Links")
	return nil
}

// ResolveLink is a helper function that allows easier traversal of links through blocks
func (*Block) ResolveLink(path []string) (*node.Link, []string, error) {
	log.Println("iscn-ResoveLink")
	return nil, []string{}, nil
}

// Size will go away. It is here to comply with the Node interface.
func (*Block) Size() (uint64, error) {
	log.Println("iscn-Size")
	return 0, nil
}

// Stat will go away. It is here to comply with the Node interface.
func (*Block) Stat() (*node.NodeStat, error) {
	log.Println("iscn-Stat")
	return &node.NodeStat{}, nil
}

// github.com/ipfs/go-ipld-format.DecodeBlockFunc

// BlockDecoder takes care of the iscn-block IPLD objects (ISCN block headers)
func BlockDecoder(block goblocks.Block) (node.Node, error) {
	log.Println("iscn-BlockDecoder")

	n, err := cbor.DecodeBlock(block)
	if err != nil {
		log.Printf("Cannot not decode block: %s", err)
		return nil, err
	}
	return n, nil
}

// Package function

// NewISCNBlock creates a iscn-block IPLD object
func NewISCNBlock(m map[string]interface{}) (*Block, error) {
	//TODO: validation code go here

	rawdata, err := cbor.DumpObject(m)
	if err != nil {
		log.Printf("Fail to marshal object: %s", err)
		return nil, err
	}

	c, err := cid.Prefix{
		Codec:    iscnblocks.CodecISCN,
		Version:  1,
		MhType:   mh.SHA2_256,
		MhLength: -1,
	}.Sum(rawdata)
	if err != nil {
		log.Printf("Fail to create CID: %s", err)
		return nil, err
	}

	block := Block{
		cid:     &c,
		rawdata: rawdata,
	}
	return &block, nil
}

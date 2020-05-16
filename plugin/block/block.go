package block

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"gitlab.com/c0b/go-ordered-json"

	blocks "github.com/ipfs/go-block-format"
	cbor "github.com/ipfs/go-ipld-cbor"
	node "github.com/ipfs/go-ipld-format"
	mh "github.com/multiformats/go-multihash"
)

// DecodeData decodes the raw IPLD data back to data object
func DecodeData(block blocks.Block) (uint64, map[string]interface{}, error) {
	data := map[string]interface{}{}
	if err := cbor.DecodeInto(block.RawData(), &data); err != nil {
		return 0, nil, err
	}

	v, ok := data[ContextKey]
	if !ok {
		return 0, nil, fmt.Errorf("Invalid ISCN IPLD object, missing context")
	}

	version, ok := v.(uint64)
	if !ok {
		return 0, nil, fmt.Errorf("Context: 'uint64' is expected but '%T' is found", v)
	}

	return version, data, nil
}

// ==================================================
// IscnObject
// ==================================================

// IscnObject is the interface for the data object of ISCN object
type IscnObject interface {
	node.Node

	GetName() string
	GetVersion() uint64
	GetCustom() map[string]interface{}

	GetBytes(string) ([]byte, error)

	SetData(map[string]interface{}) error

	Encode() error
	Decode(map[string]interface{}) error
	ToJSON() (string, error)
}

// ==================================================
// Base
// ==================================================

// Base is the basic block of all kind of ISCN objects
type Base struct {
	codec   uint64
	name    string
	version uint64
	isSet   bool
	obj     map[string]interface{}
	data    map[string]Data
	keys    []string
	custom  map[string]interface{}

	cid     *cid.Cid
	rawData []byte
}

var _ IscnObject = (*Base)(nil)

// NewBase creates the basic block of an ISCN object
func NewBase(codec uint64, name string, version uint64, schema []Data) (*Base, error) {
	// Create the base
	b := &Base{
		codec:   codec,
		name:    name,
		version: version,
		isSet:   false,
		data:    map[string]Data{},
		keys:    []string{},
	}

	// Set "context" data
	context := NewContext(name)
	err := context.Set(version)
	if err != nil {
		return nil, err
	}

	b.data[context.GetKey()] = context
	b.keys = append(b.keys, context.GetKey())

	// Setup schema
	for _, data := range schema {
		key := data.GetKey()
		b.keys = append(b.keys, key)
		b.data[key] = data
	}

	return b, nil
}

// GetName returns the name of the the ISCN object
func (b *Base) GetName() string {
	return b.name
}

// GetVersion returns the schema version of the ISCN object
func (b *Base) GetVersion() uint64 {
	return b.version
}

// GetCustom returns the custom data
func (b *Base) GetCustom() map[string]interface{} {
	return b.custom
}

// GetBytes returns the value of 'key' as byte slice
func (b *Base) GetBytes(key string) ([]byte, error) {
	value, ok := b.obj[key]
	if !ok {
		return nil, fmt.Errorf("\"%s\" is not found", key)
	}

	res, ok := value.([]byte)
	if !ok {
		return nil, fmt.Errorf("The value of \"%s\" is not '[]byte'", key)
	}

	return res, nil
}

// SetData sets and validates the data
func (b *Base) SetData(data map[string]interface{}) error {
	if b.isSet {
		return fmt.Errorf("Data has already set")
	}

	// Save the data object
	b.obj = data

	// Set and validate the data
	for key, handler := range b.data {
		// Skip context property
		if key == ContextKey {
			continue
		}

		d, ok := data[key]
		if !ok {
			return fmt.Errorf("The property \"%s\" is required", key)
		}

		err := handler.Set(d)
		if err != nil {
			return err
		}

		delete(data, key)
	}

	// Save the custom data
	b.custom = data

	b.isSet = true
	return nil
}

// Encode the ISCN object to CBOR serialized data
func (b *Base) Encode() error {
	// Extract all data from data handlers
	m := map[string]interface{}{}
	for _, handler := range b.data {
		handler.Encode(&m)
	}

	// Merge the custom data
	for key, value := range b.custom {
		m[key] = value
	}

	// CBOR-ise the data
	rawData, err := cbor.DumpObject(m)
	if err != nil {
		return err
	}

	c, err := cid.Prefix{
		Codec:    b.codec,
		Version:  1,
		MhType:   mh.SHA2_256,
		MhLength: -1,
	}.Sum(rawData)
	if err != nil {
		return err
	}

	b.cid = &c
	b.rawData = rawData

	return nil
}

// Decode the data back to ISCN object
func (b *Base) Decode(data map[string]interface{}) error {
	// Remove the context property as it is processed by base ISCN object
	delete(data, ContextKey)

	b.obj = map[string]interface{}{}
	for key, handler := range b.data {
		// Skip context property
		if key == ContextKey {
			continue
		}

		d, ok := data[key]
		if !ok {
			return fmt.Errorf("The property \"%s\" is required", key)
		}

		if err := handler.Decode(d, &b.obj); err != nil {
			return err
		}

		delete(data, key)
	}

	// Save the custom data
	b.custom = data

	b.isSet = true
	return nil
}

// ToJSON output a JSON string for the ISCN object
func (b *Base) ToJSON() (string, error) {
	om := ordered.NewOrderedMap()
	for _, key := range b.keys {
		d := b.data[key]
		if err := d.ToJSON(om); err != nil {
			return "", err
		}
	}

	for key, value := range b.custom {
		om.Set(key, value)
	}

	res, err := om.MarshalJSON()
	return string(res), err
}

// github.com/ipfs/go-block-format.Block interface

// Cid returns the CID of the block header
func (b *Base) Cid() cid.Cid {
	return *b.cid
}

// Loggable returns a map the type of IPLD Link
func (b *Base) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"type":    b.name,
		"version": b.version,
	}
}

// RawData returns the binary of the CBOR encode of the block header
func (b *Base) RawData() []byte {
	return b.rawData
}

// String is a helper for output
func (b *Base) String() string {
	return fmt.Sprintf("<%s (v%d)>", b.GetName(), b.GetVersion())
}

// node.Resolver interface

// Resolve resolves a path through this node, stopping at any link boundary
// and returning the object found as well as the remaining path to traverse
func (b *Base) Resolve(path []string) (interface{}, []string, error) {
	if len(path) == 0 {
		return b, nil, nil
	}

	_, _ = path[0], path[1:]

	// switch first {
	// // TODO: Link to nested object
	// case "parent":
	// 	return nil, nil, nil
	// case "rights":
	// 	return nil, nil, nil
	// case "stakeholders":
	// 	return nil, nil, nil
	// case "content":
	// 	return nil, nil, nil
	// }
	//
	// if len(path) != 1 {
	// 	return nil, nil, fmt.Errorf("Unexpected path elements past %s", first)
	// }
	//
	// switch first {
	// case "context":
	// 	return fmt.Sprintf("https://xxx/%d", b.Context), nil, nil
	// case "id":
	// 	return b.ID, nil, nil
	// case "timestamp":
	// 	return b.Timestamp, nil, nil
	// case "version":
	// 	return b.Version, nil, nil
	// }

	return nil, nil, nil
}

// Tree lists all paths within the object under 'path', and up to the given depth.
// To list the entire object (similar to `find .`) pass "" and -1
func (*Base) Tree(path string, depth int) []string {
	if path != "" || depth == 0 {
		return nil
	}

	return []string{
		"context",
		"id",
		"timestamp",
		"version",
		"parent",
		"rights",
		"stakeholders",
		"content",
	}
}

// node.Node interface

// Copy will go away. It is here to comply with the Node interface.
func (*Base) Copy() node.Node {
	panic("dont use this yet")
}

// Links is a helper function that returns all links within this object
// HINT: Use `ipfs refs <cid>`
func (*Base) Links() []*node.Link {
	// TODO: return link objects
	return []*node.Link{}
}

// ResolveLink is a helper function that allows easier traversal of links through blocks
func (b *Base) ResolveLink(path []string) (*node.Link, []string, error) {
	obj, rest, err := b.Resolve(path)
	if err != nil {
		return nil, nil, err
	}

	if lnk, ok := obj.(*node.Link); ok {
		return lnk, rest, nil
	}

	return nil, nil, fmt.Errorf("resolved item was not a link")
}

// Size will go away. It is here to comply with the Node interface.
func (*Base) Size() (uint64, error) {
	return 0, nil
}

// Stat will go away. It is here to comply with the Node interface.
func (*Base) Stat() (*node.NodeStat, error) {
	return &node.NodeStat{}, nil
}

package block

import (
	"fmt"
	"log"

	"github.com/ipfs/go-cid"
	"gitlab.com/c0b/go-ordered-json"

	cbor "github.com/ipfs/go-ipld-cbor"
	node "github.com/ipfs/go-ipld-format"
	mh "github.com/multiformats/go-multihash"
)

// ==================================================
// IscnObject
// ==================================================

// IscnObject is the interface for the data object of ISCN object
type IscnObject interface {
	node.Node

	GetName() string
	GetVersion() uint64
	GetCustom() map[string]interface{}

	GetArray(string) ([]interface{}, error)
	GetBytes(string) ([]byte, error)
	GetInt32(string) (int32, error)
	GetUint32(string) (uint32, error)
	GetInt64(string) (int64, error)
	GetUint64(string) (uint64, error)
	GetString(string) (string, error)
	GetCid(string) (cid.Cid, error)

	MarshalJSON() ([]byte, error)
}

// ==================================================
// CoDec
// ==================================================

// Codec is the interface for the CODEC of ISCN object
type Codec interface {
	IscnObject

	SetData(map[string]interface{}) error

	Encode() error
	Decode(map[string]interface{}) error
}

// CodecFactoryFunc returns a factory function to create ISCN object
type CodecFactoryFunc func() (Codec, error)

type codecFactory map[string][]CodecFactoryFunc

var factory codecFactory = codecFactory{}

// RegisterIscnObjectFactory registers an array of ISCN object factory functions
func RegisterIscnObjectFactory(schemaName string, factories []CodecFactoryFunc) {
	factory[schemaName] = factories
}

// Encode the data to specific ISCN object and version
func Encode(
	schemaName string,
	version uint64,
	data map[string]interface{},
) (IscnObject, error) {
	schemas, ok := factory[schemaName]
	if !ok {
		return nil, fmt.Errorf("\"%s\" is not registered", schemaName)
	}

	if version > (uint64)(len(schemas)) {
		return nil, fmt.Errorf("<%s (v%d)> is not implemented", schemaName, version)
	}
	version--

	obj, err := schemas[version]()
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

// Decode decodes the raw IPLD data back to data object and
// the CID is used for verify whether the object is consist
func Decode(schemaName string, rawData []byte, c cid.Cid) (node.Node, error) {
	data := map[string]interface{}{}
	if err := cbor.DecodeInto(rawData, &data); err != nil {
		return nil, err
	}

	v, ok := data[ContextKey]
	if !ok {
		return nil, fmt.Errorf("Invalid ISCN IPLD object, missing context")
	}

	version, ok := v.(uint64)
	if !ok {
		return nil, fmt.Errorf("Context: 'uint64' is expected but '%T' is found", v)
	}

	schemas, ok := factory[schemaName]
	if !ok {
		return nil, fmt.Errorf("\"%s\" is not registered", schemaName)
	}

	if version > (uint64)(len(schemas)) {
		return nil, fmt.Errorf("<%s (v%d)> is not implemented", schemaName, version)
	}
	version--

	obj, err := schemas[version]()
	if err != nil {
		return nil, err
	}

	if err := obj.Decode(data); err != nil {
		return nil, err
	}

	// Encode one more time to retrieve CID
	if err := obj.Encode(); err != nil {
		return nil, err
	}

	// Verify the CID
	if !obj.Cid().Equals(c) {
		current, err := obj.Cid().StringOfBase('z')
		if err != nil {
			return nil, fmt.Errorf("Cannot retrieve current CID")
		}

		expected, err := c.StringOfBase('z')
		if err != nil {
			return nil, fmt.Errorf("Cannot retrieve expected CID")
		}

		return nil, fmt.Errorf("Cid \"%s\" is not matched: expected \"%s\"", current, expected)
	}

	return obj, nil
}

// ==================================================
// Base
// ==================================================

// Base is the basic block of all kind of ISCN objects
type Base struct {
	codec   uint64
	name    string
	version uint64
	obj     map[string]interface{}
	data    map[string]Data
	keys    []string
	custom  map[string]interface{}

	cid     *cid.Cid
	rawData []byte
}

var _ Codec = (*Base)(nil)

// NewBase creates the basic block of an ISCN object
func NewBase(codec uint64, name string, version uint64, schema []Data) (*Base, error) {
	// Create the base
	b := &Base{
		codec:   codec,
		name:    name,
		version: version,
		data:    map[string]Data{},
		keys:    []string{},
		custom:  map[string]interface{}{},
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

// GetArray returns the value of 'key' as array
func (b *Base) GetArray(key string) ([]interface{}, error) {
	value, ok := b.obj[key]
	if !ok {
		return nil, fmt.Errorf("\"%s\" is not found", key)
	}

	res, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("The value of \"%s\" is not '[]interface{]'", key)
	}

	return res, nil
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

// GetInt32 returns the value of 'key' as int32
func (b *Base) GetInt32(key string) (int32, error) {
	value, ok := b.obj[key]
	if !ok {
		return 0, fmt.Errorf("\"%s\" is not found", key)
	}

	res, ok := value.(int32)
	if !ok {
		return 0, fmt.Errorf("The value of \"%s\" is not 'int32'", key)
	}

	return res, nil
}

// GetUint32 returns the value of 'key' as int32
func (b *Base) GetUint32(key string) (uint32, error) {
	value, ok := b.obj[key]
	if !ok {
		return 0, fmt.Errorf("\"%s\" is not found", key)
	}

	res, ok := value.(uint32)
	if !ok {
		return 0, fmt.Errorf("The value of \"%s\" is not 'uint32'", key)
	}

	return res, nil
}

// GetInt64 returns the value of 'key' as int64
func (b *Base) GetInt64(key string) (int64, error) {
	value, ok := b.obj[key]
	if !ok {
		return 0, fmt.Errorf("\"%s\" is not found", key)
	}

	res, ok := value.(int64)
	if !ok {
		return 0, fmt.Errorf("The value of \"%s\" is not 'int64'", key)
	}

	return res, nil
}

// GetUint64 returns the value of 'key' as int64
func (b *Base) GetUint64(key string) (uint64, error) {
	value, ok := b.obj[key]
	if !ok {
		return 0, fmt.Errorf("\"%s\" is not found", key)
	}

	res, ok := value.(uint64)
	if !ok {
		return 0, fmt.Errorf("The value of \"%s\" is not 'uint64'", key)
	}

	return res, nil
}

// GetString returns the value of 'key' as string
func (b *Base) GetString(key string) (string, error) {
	value, ok := b.obj[key]
	if !ok {
		return "", fmt.Errorf("\"%s\" is not found", key)
	}

	res, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("The value of \"%s\" is not '[]byte'", key)
	}

	return res, nil
}

// GetCid returns the value of 'key' as Cid
func (b *Base) GetCid(key string) (cid.Cid, error) {
	value, ok := b.obj[key]
	if !ok {
		return cid.Undef, fmt.Errorf("\"%s\" is not found", key)
	}

	res, ok := value.(cid.Cid)
	if !ok {
		return cid.Undef, fmt.Errorf("The value of \"%s\" is not 'Cid'", key)
	}

	return res, nil
}

// MarshalJSON convert the block to JSON format
func (b *Base) MarshalJSON() ([]byte, error) {
	om := ordered.NewOrderedMap()
	for _, key := range b.keys {
		handler := b.data[key]

		if key != ContextKey { // Context key does not exist in b.obj
			_, exist := b.obj[key]
			if !exist {
				if handler.IsRequired() {
					return nil, fmt.Errorf("Unknown error: key %q should be exist", key)
				}
				continue
			}
		}

		if err := handler.ToJSON(om); err != nil {
			return nil, err
		}
	}

	for key, value := range b.custom {
		om.Set(key, value)
	}

	return om.MarshalJSON()
}

// SetData sets and validates the data
func (b *Base) SetData(data map[string]interface{}) error {
	// Save the data object
	b.obj = data

	// Set the data
	for key, handler := range b.data {
		// Skip context property
		if key == ContextKey {
			continue
		}

		d, ok := data[key]
		if !ok {
			if handler.IsRequired() {
				return fmt.Errorf("The property \"%s\" is required", key)
			}

			continue
		}

		err := handler.Set(d)
		if err != nil {
			return err
		}
	}

	// Validate the data
	for _, handler := range b.data {
		if err := handler.Validate(); err != nil {
			return err
		}
	}

	// Save the custom data
	for key, value := range data {
		_, exist := b.data[key]
		if !exist {
			b.custom[key] = value
		}
	}

	return nil
}

// Encode the ISCN object to CBOR serialized data
func (b *Base) Encode() error {
	// Extract all data from data handlers
	m := map[string]interface{}{}
	for _, handler := range b.data {
		if handler.GetKey() != ContextKey { // Context key does not exist in b.obj
			_, exist := b.obj[handler.GetKey()]
			if !exist {
				if handler.IsRequired() {
					return fmt.Errorf("Unknown error: key %q should be exist", handler.GetKey())
				}
				continue
			}
		}

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

	c, err := cid.V1Builder{
		Codec:  b.codec,
		MhType: mh.SHA2_256,
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
			if handler.IsRequired() {
				return fmt.Errorf("The property \"%s\" is required", key)
			}

			continue
		}

		if err := handler.Decode(d, &b.obj); err != nil {
			return err
		}

		delete(data, key)
	}

	// Save the custom data
	b.custom = data

	return nil
}

// github.com/ipfs/go-block-format.Block interface

// Cid returns the CID of the block header
func (b *Base) Cid() cid.Cid {
	return *(b.cid)
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

	first, rest := path[0], path[1:]

	if data, ok := b.data[first]; ok {
		if first != ContextKey { // Context key does not exist in b.obj
			_, exist := b.obj[first]
			if !exist {
				if data.IsRequired() {
					return nil, nil, fmt.Errorf("Unknown error: key %q should be exist", first)
				}
				return nil, nil, fmt.Errorf("no such link")
			}
		}
		return data.Resolve(rest)
	}

	// TODO: custom parameters

	return nil, nil, fmt.Errorf("no such link")
}

// Tree lists all paths within the object under 'path', and up to the given depth.
// To list the entire object (similar to `find .`) pass "" and -1
func (*Base) Tree(path string, depth int) []string {
	log.Println("Tree is not implemented")
	return nil
}

// node.Node interface

// Copy will go away. It is here to comply with the Node interface.
func (*Base) Copy() node.Node {
	panic("dont use this yet")
}

// Links is a helper function that returns all links within this object
// HINT: Use `ipfs refs <cid>`
func (b *Base) Links() []*node.Link {
	links := []*node.Link{}
	for _, d := range b.data {
		switch c := d.(type) {
		case *Cid:
			if link, err := c.Link(); err == nil {
				links = append(links, link)
			}
		}
	}
	return links
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

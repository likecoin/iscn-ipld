package block

import (
	"fmt"

	"gitlab.com/c0b/go-ordered-json"
)

const (
	// ContextKey is the key of context of ISCN object
	ContextKey = "context"
)

// ==================================================
// Data
// ==================================================

// Data is the interface for the data property handler
type Data interface {
	IsLink() bool
	IsNested() bool

	Set(interface{}) error
	GetKey() string

	Encode(*map[string]interface{}) error
	ToJSON(*ordered.OrderedMap) error
	Decode(interface{}, *map[string]interface{}) error
}

// ==================================================
// DataBase
// ==================================================

// DataBase is the base struct for handling data property
type DataBase struct {
	isLink   bool
	isNested bool
	key      string
}

// NewDataBase creates a base struct for handling data property
func NewDataBase(isLink bool, isNested bool, key string) *DataBase {
	return &DataBase{
		isLink:   isLink,
		isNested: isNested,
		key:      key,
	}
}

// IsLink returns whether the data object is a link
func (b *DataBase) IsLink() bool {
	return b.isLink
}

// IsNested returns whether the data object is a nested object
func (b *DataBase) IsNested() bool {
	return b.isNested
}

// GetKey returns the key of the data property
func (b *DataBase) GetKey() string {
	return b.key
}

// ==================================================
// Context
// ==================================================

// Context is a data type for the context of ISCN object
type Context struct {
	*DataBase

	schema  string
	version uint64
}

var _ Data = (*Context)(nil)

// NewContext creates a context of ISCN object with schema name
func NewContext(schema string) *Context {
	return &Context{
		DataBase: NewDataBase(false, false, ContextKey),
		schema:   schema,
	}
}

// Set the value of Context
func (d *Context) Set(data interface{}) error {
	if version, ok := data.(uint64); ok {
		d.version = version
		return nil
	}

	return fmt.Errorf("Context: 'uint64' is expected but '%T' is found", data)
}

// Encode Context
func (d *Context) Encode(m *map[string]interface{}) error {
	(*m)[d.GetKey()] = d.version
	return nil
}

// ToJSON prepares the data for MarshalJSON
func (d *Context) ToJSON(om *ordered.OrderedMap) error {
	// TODO use the real schema path
	om.Set(d.GetKey(), fmt.Sprintf("schema/%s-v%d", d.schema, d.version))
	return nil
}

// Decode Context
func (d *Context) Decode(data interface{}, m *map[string]interface{}) error {
	if err := d.Set(data); err != nil {
		return err
	}

	(*m)[d.GetKey()] = d.version
	return nil
}

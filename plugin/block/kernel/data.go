package kernel

import (
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/likecoin/iscn-ipld/plugin/block"
	"gitlab.com/c0b/go-ordered-json"
)

// ==================================================
// ID
// ==================================================

// ID is a data handler for the ISCN ID
type ID struct {
	*block.DataBase

	id []byte
}

var _ block.Data = (*ID)(nil)

// NewID creates an ISCN ID
func NewID() *ID {
	return &ID{
		DataBase: block.NewDataBase("id", true),
	}
}

// Prototype creates a prototype ID
func (d *ID) Prototype() block.Data {
	return &ID{
		DataBase: d.DataBase.Prototype(),
	}
}

// GetID returns the human readable ID
func (d *ID) GetID() string {
	return fmt.Sprintf("1/%s", base58.Encode(d.id))
}

// Set the value of ID
func (d *ID) Set(data interface{}) error {
	if id, ok := data.([]byte); ok {
		if len(id) != 32 {
			return fmt.Errorf("ID: should length 32 but %d is found", len(id))
		}

		d.id = id
		return nil
	}

	return fmt.Errorf("ID: '[]byte' is expected but '%T' is found", data)
}

// Encode ID
func (d *ID) Encode(m *map[string]interface{}) error {
	(*m)[d.GetKey()] = d.id
	return nil
}

// Decode ID
func (d *ID) Decode(data interface{}, m *map[string]interface{}) error {
	if err := d.Set(data); err != nil {
		return err
	}

	(*m)[d.GetKey()] = d.id
	return nil
}

// ToJSON prepares the data for MarshalJSON
func (d *ID) ToJSON(om *ordered.OrderedMap) error {
	om.Set(d.GetKey(), d.GetID())
	return nil
}

// Resolve resolves the value
func (d *ID) Resolve(path []string) (interface{}, []string, error) {
	if len(path) != 0 {
		return nil, nil, fmt.Errorf("Unexpected path elements past %s", path[0])
	}

	return d.GetID(), nil, nil
}

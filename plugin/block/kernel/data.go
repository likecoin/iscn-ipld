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

// ID is a data type for the ISCN ID
type ID struct {
	*block.DataBase

	id []byte
}

var _ block.Data = (*ID)(nil)

// NewID creates an ISCN ID
func NewID() *ID {
	return &ID{
		DataBase: block.NewDataBase(false, false, "id"),
	}
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

// ToJSON prepares the data for MarshalJSON
func (d *ID) ToJSON(om *ordered.OrderedMap) error {
	om.Set(d.GetKey(), fmt.Sprintf("1/%s", base58.Encode(d.id)))
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

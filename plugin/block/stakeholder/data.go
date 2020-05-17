package stakeholder

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/likecoin/iscn-ipld/plugin/block"
	"gitlab.com/c0b/go-ordered-json"
)

// ==================================================
// Type
// ==================================================

const (
	creator     = "Creator"
	contributor = "Contributor"
	editor      = "Editor"
	publisher   = "Publisher"
	footprint   = "FootprintStakeholder"
	escrow      = "Escrow"
)

// Type is a data handler for the type of stakeholder
type Type struct {
	*block.String
}

var _ block.Data = (*Type)(nil)

// NewType creates a stakeholder type data handler
func NewType() *Type {
	return &Type{
		String: block.NewStringWithFilter(
			"type",
			true,
			[]string{
				creator,
				contributor,
				editor,
				publisher,
				footprint,
				escrow,
			},
		),
	}
}

// Prototype creates a prototype Type
func (d *Type) Prototype() block.Data {
	return NewType()
}

// ==================================================
// Footprint
// ==================================================

// Footprint is a data handler for the footprint link to the underlying work
type Footprint struct {
	*block.DataBase

	ty      *Type
	handler block.Data
}

var _ block.Data = (*Footprint)(nil)

// NewFootprint creates a footprint data handler
func NewFootprint(ty *Type) *Footprint {
	return &Footprint{
		DataBase: block.NewDataBase("footprint", false),
		ty:       ty,
	}
}

// Prototype creates a protype Footprint
func (d *Footprint) Prototype() block.Data {
	panic("Footprint do not support prototyping")
}

// Set the value of link of footprint
func (d *Footprint) Set(data interface{}) error {
	if d.handler != nil {
		return fmt.Errorf("Footprint: re-create handler")
	}

	switch data.(type) {
	case cid.Cid:
		d.handler = block.NewCid(d.GetKey(), d.IsRequired(), block.CodecISCN)
	case string:
		// TODO URL handler
		d.handler = block.NewString(d.GetKey(), d.IsRequired())
	default:
		return fmt.Errorf("Footprint: link is expected but '%T' is found", data)
	}

	return d.handler.Set(data)
}

// Validate the data
func (d *Footprint) Validate() error {
	if d.isFootprint() {
		if d.handler == nil {
			return fmt.Errorf("Footprint: missing footprint")
		}
	} else {
		if d.handler != nil {
			return fmt.Errorf("Footprint: should not be set as " +
				"this is not a footprint stakeholder")
		}
	}

	return nil
}

// Encode Footprint
func (d *Footprint) Encode(m *map[string]interface{}) error {
	if d.isFootprint() {
		return d.handler.Encode(m)
	}

	return nil
}

// Decode Footprint
func (d *Footprint) Decode(data interface{}, m *map[string]interface{}) error {
	if d.handler != nil {
		return fmt.Errorf("Footprint: re-create handler")
	}

	switch data.(type) {
	case []uint8:
		d.handler = block.NewCid(d.GetKey(), d.IsRequired(), block.CodecISCN)
	case string:
		// TODO URL handler
		d.handler = block.NewString(d.GetKey(), d.IsRequired())
	default:
		return fmt.Errorf("Footprint: link is expected but '%T' is found", data)
	}

	return d.handler.Decode(data, m)
}

// ToJSON prepares the data for MarshalJSON
func (d *Footprint) ToJSON(om *ordered.OrderedMap) error {
	if d.isFootprint() {
		return d.handler.ToJSON(om)
	}

	return nil
}

// Resolve resolves the link
func (d *Footprint) Resolve(path []string) (interface{}, []string, error) {
	if d.isFootprint() {
		return d.handler.Resolve(path)
	}

	return nil, nil, fmt.Errorf("no such link")
}

func (d *Footprint) isFootprint() bool {
	return d.ty.Get() == footprint
}

package timeperiod

import (
	"fmt"

	"github.com/likecoin/iscn-ipld/plugin/block"
)

const (
	// SchemaName of time period
	SchemaName = "timeperiod"
)

// Register registers the schema of time period block
func Register() {
	block.RegisterIscnObjectFactory(
		block.CodecTimePeriod,
		SchemaName,
		[]block.CodecFactoryFunc{
			newSchemaV1,
		},
	)
}

// ==================================================
// base
// ==================================================

// base is the base struct for time period (codec 0x033F)
type base struct {
	*block.Base
}

func newBase(version uint64, schema []block.Data) (*base, error) {
	blockBase, err := block.NewBase(
		block.CodecTimePeriod,
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

// schemaV1 represents a time period V1
type schemaV1 struct {
	*base

	from *block.Timestamp
	to   *block.Timestamp
}

var _ block.IscnObject = (*schemaV1)(nil)

func newSchemaV1() (block.Codec, error) {
	from := block.NewTimestamp("from", false)
	to := block.NewTimestamp("to", false)

	schema := []block.Data{
		from,
		to,
	}

	timePeriodBase, err := newBase(1, schema)
	if err != nil {
		return nil, err
	}

	obj := schemaV1{
		base: timePeriodBase,
		from: from,
		to:   to,
	}
	timePeriodBase.SetValidator(obj.Validate)

	return &obj, nil
}

// SchemaV1Prototype creates a prototype for schemaV1
func SchemaV1Prototype() block.Codec {
	res, _ := newSchemaV1()
	return res
}

// Validate the data
func (o *schemaV1) Validate() error {
	if !o.from.IsDefined() && !o.to.IsDefined() {
		return fmt.Errorf("At least \"from\" or \"to\" exists")
	}

	return nil
}

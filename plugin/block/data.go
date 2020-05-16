package block

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"

	"github.com/ipfs/go-cid"
	"gitlab.com/c0b/go-ordered-json"

	node "github.com/ipfs/go-ipld-format"
)

// ==================================================
// Data
// ==================================================

// Data is the interface for the data property handler
type Data interface {
	Prototype() Data

	IsRequired() bool

	Set(interface{}) error
	GetKey() string

	Validate() error

	Encode(*map[string]interface{}) error
	Decode(interface{}, *map[string]interface{}) error
	ToJSON(*ordered.OrderedMap) error

	Resolve(path []string) (interface{}, []string, error)
}

// ==================================================
// DataBase
// ==================================================

// DataBase is the base struct for handling data property
type DataBase struct {
	isRequired bool

	key string
}

// NewDataBase creates a base struct for handling data property
func NewDataBase(key string, isRequired bool) *DataBase {
	return &DataBase{
		isRequired: isRequired,
		key:        key,
	}
}

// Prototype creates a prototype DataBase
func (b *DataBase) Prototype() *DataBase {
	return &DataBase{
		isRequired: b.isRequired,
		key:        b.key,
	}
}

// IsRequired checks whether the data handler is required
func (b *DataBase) IsRequired() bool {
	return b.isRequired
}

// GetKey returns the key of the data property
func (b *DataBase) GetKey() string {
	return b.key
}

// Validate the data
func (b *DataBase) Validate() error {
	return nil
}

// ==================================================
// DataArray
// ==================================================

// DataArray is an array of data handler
type DataArray struct {
	*DataBase

	array     []Data
	prototype Data
}

var _ Data = (*DataArray)(nil)

// NewDataArray creates an array of data handler
func NewDataArray(key string, isRequired bool, prototype Data) *DataArray {
	return &DataArray{
		DataBase:  NewDataBase(key, isRequired),
		array:     []Data{},
		prototype: prototype,
	}
}

// Prototype creates a prototype DataArray
func (d *DataArray) Prototype() Data {
	return &DataArray{
		DataBase:  d.DataBase.Prototype(),
		array:     []Data{},
		prototype: d.prototype.Prototype(),
	}
}

// Set the value of data handler array
func (d *DataArray) Set(data interface{}) error {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(data)
		for i := 0; i < s.Len(); i++ {
			elem := d.prototype.Prototype()
			if err := elem.Set(s.Index(i).Interface()); err != nil {
				return fmt.Errorf("(Index %d) %s", i, err.Error())
			}
			d.array = append(d.array, elem)
		}

		return nil
	}

	return fmt.Errorf("DataArray: an array is expected but '%T' is found", data)
}

// Encode DataArray
func (d *DataArray) Encode(m *map[string]interface{}) error {
	placeholder := map[string]interface{}{}
	res := []interface{}{}
	for i, data := range d.array {
		if err := data.Encode(&placeholder); err != nil {
			return fmt.Errorf("(Index %d) %s", i, err.Error())
		}
		res = append(res, placeholder[data.GetKey()])
	}

	(*m)[d.GetKey()] = res
	return nil
}

// Decode DataArray
func (d *DataArray) Decode(data interface{}, m *map[string]interface{}) error {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		placeholder := map[string]interface{}{}
		res := []interface{}{}
		s := reflect.ValueOf(data)
		for i := 0; i < s.Len(); i++ {
			elem := d.prototype.Prototype()
			if err := elem.Decode(s.Index(i).Interface(), &placeholder); err != nil {
				return fmt.Errorf("(Index %d) %s", i, err.Error())
			}

			res = append(res, placeholder[elem.GetKey()])
			d.array = append(d.array, elem)
		}

		(*m)[d.GetKey()] = res
		return nil
	}

	return fmt.Errorf("DataArray: an array is expected but '%T' is found", data)
}

// ToJSON prepares the data for MarshalJSON
func (d *DataArray) ToJSON(om *ordered.OrderedMap) error {
	placeholder := ordered.NewOrderedMap()
	res := []interface{}{}
	for i, data := range d.array {
		if err := data.ToJSON(placeholder); err != nil {
			return fmt.Errorf("(Index %d) %s", i, err.Error())
		}
		res = append(res, placeholder.Get(data.GetKey()))
	}

	om.Set(d.GetKey(), res)
	return nil
}

// Resolve resolves the value
func (d *DataArray) Resolve(path []string) (interface{}, []string, error) {
	if len(path) == 0 {
		res := []interface{}{}
		for i, data := range d.array {
			value, _, err := data.Resolve(path)
			if err != nil {
				return nil, nil, fmt.Errorf("(Index %d) %s", i, err.Error())
			}
			res = append(res, value)
		}
		return res, nil, nil
	}

	first, rest := path[0], path[1:]
	index, err := strconv.ParseUint(first, 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf("Unexpected path elements past %s", path[0])
	}

	if index >= uint64(len(d.array)) {
		return nil, nil, fmt.Errorf("index %d does not exist", index)
	}

	return d.array[index].Resolve(rest)
}

// ==================================================
// Number
// ==================================================

// NumberType is a enum type for the type of a number
type NumberType int

const (
	// Int32T represents int32
	Int32T NumberType = iota

	// Uint32T represents uint32
	Uint32T

	// Int64T represents int64
	Int64T

	// Uint64T represents uint64
	Uint64T
)

// Number is a data handler for the number
type Number struct {
	*DataBase

	number []byte
	ty     NumberType

	i32 int32
	u32 uint32
	i64 int64
	u64 uint64
}

var _ Data = (*Number)(nil)

// NewNumber creates a number data handler
func NewNumber(key string, isRequired bool, ty NumberType) *Number {
	return &Number{
		DataBase: NewDataBase(key, isRequired),
		ty:       ty,
	}
}

// Prototype creates a prototype Number
func (d *Number) Prototype() Data {
	return &Number{
		DataBase: d.DataBase.Prototype(),
		ty:       d.ty,
	}
}

// GetType returns the type of the number
func (d *Number) GetType() NumberType {
	return d.ty
}

// GetInt32 returns an int32 value
func (d *Number) GetInt32() (int32, error) {
	if d.GetType() != Int32T {
		return 0, fmt.Errorf("Number: is not 'int32'")
	}
	return d.i32, nil
}

// GetUint32 returns an uint32 value
func (d *Number) GetUint32() (uint32, error) {
	if d.GetType() != Uint32T {
		return 0, fmt.Errorf("Number: is not 'uint32'")
	}
	return d.u32, nil
}

// GetInt64 returns an int64 value
func (d *Number) GetInt64() (int64, error) {
	if d.GetType() != Int64T {
		return 0, fmt.Errorf("Number: is not 'int64'")
	}
	return d.i64, nil
}

// GetUint64 returns an uint64 value
func (d *Number) GetUint64() (uint64, error) {
	if d.GetType() != Uint64T {
		return 0, fmt.Errorf("Number: is not 'uint64'")
	}
	return d.u64, nil
}

// Set the value of number
func (d *Number) Set(data interface{}) error {
	switch d.GetType() {
	case Int32T:
		var value int32
		switch v := data.(type) {
		case int:
			value = int32(v)
		case int8:
			value = int32(v)
		case int16:
			value = int32(v)
		case int32:
			value = v
		case int64:
			if v < math.MinInt32 || math.MaxInt32 < v {
				return fmt.Errorf("Number: 'int32' is expected but 'int64' is found")
			}
			value = int32(v)
		case uint:
			if v > math.MaxInt32 {
				return fmt.Errorf("Number: 'int32' is expected but 'uint' is found")
			}
			value = int32(v)
		case uint8:
			value = int32(v)
		case uint16:
			value = int32(v)
		case uint32:
			if v > math.MaxInt32 {
				return fmt.Errorf("Number: 'int32' is expected but 'uint32' is found")
			}
			value = int32(v)
		case uint64:
			if v > math.MaxInt32 {
				return fmt.Errorf("Number: 'int32' is expected but 'uint64' is found")
			}
			value = int32(v)
		default:
			return fmt.Errorf("Number: 'int32' is expected but '%T' is found", data)
		}

		buffer := make([]byte, binary.MaxVarintLen32)
		n := binary.PutVarint(buffer, int64(value))
		d.number = buffer[:n]
		d.i32 = value
	case Uint32T:
		var value uint32
		switch v := data.(type) {
		case int:
			if v < 0 {
				return fmt.Errorf("Number: 'uint32' is expected but 'int' is found")
			}
			value = uint32(v)
		case int8:
			if v < 0 {
				return fmt.Errorf("Number: 'uint32' is expected but 'int8' is found")
			}
			value = uint32(v)
		case int16:
			if v < 0 {
				return fmt.Errorf("Number: 'uint32' is expected but 'int16' is found")
			}
			value = uint32(v)
		case int32:
			if v < 0 {
				return fmt.Errorf("Number: 'uint32' is expected but 'int32' is found")
			}
			value = uint32(v)
		case int64:
			if v < 0 || math.MaxUint32 < v {
				return fmt.Errorf("Number: 'uint32' is expected but 'int64' is found")
			}
			value = uint32(v)
		case uint:
			if v > math.MaxUint32 {
				return fmt.Errorf("Number: 'uint32' is expected but 'uint' is found")
			}
			value = uint32(v)
		case uint8:
			value = uint32(v)
		case uint16:
			value = uint32(v)
		case uint32:
			value = v
		case uint64:
			if v > math.MaxUint32 {
				return fmt.Errorf("Number: 'uint32' is expected but 'uint64' is found")
			}
			value = uint32(v)
		default:
			return fmt.Errorf("Number: 'uint32' is expected but '%T' is found", data)
		}

		buffer := make([]byte, binary.MaxVarintLen32)
		n := binary.PutUvarint(buffer, uint64(value))
		d.number = buffer[:n]
		d.u32 = value
	case Int64T:
		var value int64
		switch v := data.(type) {
		case int:
			value = int64(v)
		case int8:
			value = int64(v)
		case int16:
			value = int64(v)
		case int32:
			value = int64(v)
		case int64:
			value = v
		case uint:
			value = int64(v)
		case uint8:
			value = int64(v)
		case uint16:
			value = int64(v)
		case uint32:
			value = int64(v)
		case uint64:
			if v > math.MaxInt64 {
				return fmt.Errorf("Number: 'int64' is expected but 'uint64' is found")
			}
			value = int64(v)
		default:
			return fmt.Errorf("Number: 'int64' is expected but '%T' is found", data)
		}

		buffer := make([]byte, binary.MaxVarintLen64)
		n := binary.PutVarint(buffer, value)
		d.number = buffer[:n]
		d.i64 = value
	case Uint64T:
		var value uint64
		switch v := data.(type) {
		case int:
			if v < 0 {
				return fmt.Errorf("Number: 'uint64' is expected but 'int' is found")
			}
			value = uint64(v)
		case int8:
			if v < 0 {
				return fmt.Errorf("Number: 'uint64' is expected but 'int8' is found")
			}
			value = uint64(v)
		case int16:
			if v < 0 {
				return fmt.Errorf("Number: 'uint64' is expected but 'int16' is found")
			}
			value = uint64(v)
		case int32:
			if v < 0 {
				return fmt.Errorf("Number: 'uint64' is expected but 'int32' is found")
			}
			value = uint64(v)
		case int64:
			if v < 0 {
				return fmt.Errorf("Number: 'uint64' is expected but 'int64' is found")
			}
			value = uint64(v)
		case uint:
			value = uint64(v)
		case uint8:
			value = uint64(v)
		case uint16:
			value = uint64(v)
		case uint32:
			value = uint64(v)
		case uint64:
			value = v
		default:
			return fmt.Errorf("Number: 'uint64' is expected but '%T' is found", data)
		}

		buffer := make([]byte, binary.MaxVarintLen64)
		n := binary.PutUvarint(buffer, value)
		d.number = buffer[:n]
		d.u64 = value
	}

	return nil
}

// Encode Number
func (d *Number) Encode(m *map[string]interface{}) error {
	(*m)[d.GetKey()] = d.number
	return nil
}

// Decode Number
func (d *Number) Decode(data interface{}, m *map[string]interface{}) error {
	number, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("Unknown error during decoding number: "+
			"'[]byte' is expected but '%T' is found",
			data,
		)
	}

	var err error
	r := bytes.NewReader(number)
	switch d.GetType() {
	case Int32T:
		value, err := binary.ReadVarint(r)
		if err == nil {
			d.i32 = int32(value)
			(*m)[d.GetKey()] = d.i32
		}
	case Uint32T:
		value, err := binary.ReadUvarint(r)
		if err == nil {
			d.u32 = uint32(value)
			(*m)[d.GetKey()] = d.u32
		}
	case Int64T:
		value, err := binary.ReadVarint(r)
		if err == nil {
			d.i64 = int64(value)
			(*m)[d.GetKey()] = d.i64
		}
	case Uint64T:
		value, err := binary.ReadUvarint(r)
		if err == nil {
			d.u64 = uint64(value)
			(*m)[d.GetKey()] = d.u64
		}
	}

	if err == nil {
		d.number = number
	}

	return err
}

// ToJSON prepares the data for MarshalJSON
func (d *Number) ToJSON(om *ordered.OrderedMap) error {
	switch d.GetType() {
	case Int32T:
		om.Set(d.GetKey(), d.i32)
	case Uint32T:
		om.Set(d.GetKey(), d.u32)
	case Int64T:
		om.Set(d.GetKey(), d.i64)
	case Uint64T:
		om.Set(d.GetKey(), d.u64)
	}
	return nil
}

// Resolve resolves the value
func (d *Number) Resolve(path []string) (interface{}, []string, error) {
	if len(path) != 0 {
		return nil, nil, fmt.Errorf("Unexpected path elements past %s", path[0])
	}

	switch d.GetType() {
	case Int32T:
		return d.i32, nil, nil
	case Uint32T:
		return d.u32, nil, nil
	case Int64T:
		return d.i64, nil, nil
	case Uint64T:
		return d.u64, nil, nil
	}

	return nil, nil, fmt.Errorf("Number: unknown error")
}

// ==================================================
// String
// ==================================================

// String is a data handler for the string
type String struct {
	*DataBase

	value string
}

var _ Data = (*String)(nil)

// NewString creates a string data handler
func NewString(key string, isRequired bool) *String {
	return &String{
		DataBase: NewDataBase(key, isRequired),
	}
}

// Prototype creates a prototype String
func (d *String) Prototype() Data {
	return &String{
		DataBase: d.DataBase.Prototype(),
	}
}

// Set the value of String
func (d *String) Set(data interface{}) error {
	if value, ok := data.(string); ok {
		d.value = value
		return nil
	}

	return fmt.Errorf("String: 'string' is expected but '%T' is found", data)
}

// Encode String
func (d *String) Encode(m *map[string]interface{}) error {
	(*m)[d.GetKey()] = d.value
	return nil
}

// Decode String
func (d *String) Decode(data interface{}, m *map[string]interface{}) error {
	if err := d.Set(data); err != nil {
		return err
	}

	(*m)[d.GetKey()] = d.value
	return nil
}

// ToJSON prepares the data for MarshalJSON
func (d *String) ToJSON(om *ordered.OrderedMap) error {
	om.Set(d.GetKey(), d.value)
	return nil
}

// Resolve resolves the value
func (d *String) Resolve(path []string) (interface{}, []string, error) {
	if len(path) != 0 {
		return nil, nil, fmt.Errorf("Unexpected path elements past %s", path[0])
	}

	return d.value, nil, nil
}

// ==================================================
// Context
// ==================================================

const (
	// ContextKey is the key of context of ISCN object
	ContextKey = "context"
)

// Context is a data handler for the context of ISCN object
type Context struct {
	*Number

	schema  string
	version uint64
}

var _ Data = (*Context)(nil)

// NewContext creates a context of ISCN object with schema name
func NewContext(schema string) *Context {
	return &Context{
		Number: NewNumber(ContextKey, true, Uint64T),
		schema: schema,
	}
}

// Set the value of Context
func (d *Context) Set(data interface{}) error {
	err := d.Number.Set(data)
	if err != nil {
		return fmt.Errorf("Context: 'uint64' is expected but '%T' is found", data)
	}

	version, err := d.GetUint64()
	if err != nil {
		return fmt.Errorf("Context: %s", err)
	}

	d.version = version
	return nil
}

// Encode Context
func (d *Context) Encode(m *map[string]interface{}) error {
	(*m)[d.GetKey()] = d.version
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

// ToJSON prepares the data for MarshalJSON
func (d *Context) ToJSON(om *ordered.OrderedMap) error {
	om.Set(d.GetKey(), d.getSchemaURL())
	return nil
}

// Resolve resolves the value
func (d *Context) Resolve(path []string) (interface{}, []string, error) {
	if len(path) != 0 {
		return nil, nil, fmt.Errorf("Unexpected path elements past %s", path[0])
	}

	return d.getSchemaURL(), nil, nil
}

func (d *Context) getSchemaURL() string {
	// TODO use the real schema path
	return fmt.Sprintf("schema/%s-v%d", d.schema, d.version)
}

// ==================================================
// Cid
// ==================================================

// Cid is a data handler for IPFS CID
type Cid struct {
	*DataBase

	codec uint64
	c     []byte
}

var _ Data = (*Cid)(nil)

// NewCid creates a IPFS CID data handler
func NewCid(key string, isRequired bool, codec uint64) *Cid {
	return &Cid{
		DataBase: NewDataBase(key, isRequired),
		codec:    codec,
	}
}

// Prototype creates a prototype Cid
func (d *Cid) Prototype() Data {
	return &Cid{
		DataBase: d.DataBase.Prototype(),
		codec:    d.codec,
	}
}

// Link returns a link object for IPLD
func (d *Cid) Link() (*node.Link, error) {
	_, c, err := cid.CidFromBytes(d.c)
	if err != nil {
		return nil, err
	}

	return &node.Link{Cid: c}, nil
}

// Set the value of IPFS CID
func (d *Cid) Set(data interface{}) error {
	if c, ok := data.(cid.Cid); ok {
		if c.Type() != d.codec {
			return fmt.Errorf(
				"Cid: Codec '0x%x' is expected but '0x%x' is found",
				d.codec,
				c.Type())
		}

		d.c = c.Bytes()
		return nil
	}

	return fmt.Errorf("Cid: 'cid.Cid' is expected but '%T' is found", data)
}

// Encode Cid
func (d *Cid) Encode(m *map[string]interface{}) error {
	(*m)[d.GetKey()] = d.c
	return nil
}

// Decode Cid
func (d *Cid) Decode(data interface{}, m *map[string]interface{}) error {
	c, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("Unknown error during decoding Cid: "+
			"'[]byte' is expected but '%T' is found",
			data,
		)
	}

	_, value, err := cid.CidFromBytes(c)
	if err != nil {
		return err
	}

	if value.Type() != d.codec {
		return fmt.Errorf(
			"Cid: Codec '0x%x' is expected but '0x%x' is found",
			d.codec,
			value.Type())
	}

	d.c = c
	(*m)[d.GetKey()] = value
	return nil
}

// ToJSON prepares the data for MarshalJSON
func (d *Cid) ToJSON(om *ordered.OrderedMap) error {
	_, c, err := cid.CidFromBytes(d.c)
	if err != nil {
		return err
	}

	value, err := c.StringOfBase('z')
	if err != nil {
		return err
	}

	link := map[string]string{
		"/": fmt.Sprintf("/ipfs/%s", value),
	}

	om.Set(d.GetKey(), link)
	return nil
}

// Resolve resolves the link
func (d *Cid) Resolve(path []string) (interface{}, []string, error) {
	link, err := d.Link()
	if err != nil {
		return nil, nil, err
	}

	return link, path, nil
}

// ==================================================
// Parent
// ==================================================

// Parent is a data handler for IPFS CID which links to previous version
type Parent struct {
	*Cid

	version *Number
}

var _ Data = (*Parent)(nil)

// NewParent creates a parent IPFS CID data handler
func NewParent(key string, codec uint64, version *Number) *Parent {
	return &Parent{
		Cid:     NewCid(key, false, codec),
		version: version,
	}
}

// Set the value of parent IPFS CID
func (d *Parent) Set(data interface{}) error {
	if data != nil {
		return d.Cid.Set(data)
	}
	return nil
}

// Validate the data
func (d *Parent) Validate() error {
	hasParent, err := d.hasParent()
	if err != nil {
		return err
	}

	if hasParent {
		if d.Cid.c == nil {
			return fmt.Errorf("Parent: missing  as version > 1")
		}
	} else {
		if d.Cid.c != nil {
			return fmt.Errorf("Parent: should not be set as version <= 1")
		}
	}

	return nil
}

// Encode Parent
func (d *Parent) Encode(m *map[string]interface{}) error {
	hasParent, err := d.hasParent()
	if err != nil {
		return err
	}

	if hasParent {
		return d.Cid.Encode(m)
	}

	return nil
}

// Decode Parent
func (d *Parent) Decode(data interface{}, m *map[string]interface{}) error {
	return d.Cid.Decode(data, m)
}

// ToJSON prepares the data for MarshalJSON
func (d *Parent) ToJSON(om *ordered.OrderedMap) error {
	hasParent, err := d.hasParent()
	if err != nil {
		return err
	}

	if hasParent {
		return d.Cid.ToJSON(om)
	}

	return nil
}

// Resolve resolves the link
func (d *Parent) Resolve(path []string) (interface{}, []string, error) {
	hasParent, err := d.hasParent()
	if err != nil {
		return nil, nil, err
	}

	if hasParent {
		return d.Cid.Resolve(path)
	}

	return nil, nil, fmt.Errorf("no such link")
}

func (d *Parent) hasParent() (bool, error) {
	switch d.version.GetType() {
	case Int32T:
		value, _ := d.version.GetInt32()
		return value > 1, nil
	case Uint32T:
		value, _ := d.version.GetUint32()
		return value > 1, nil
	case Int64T:
		value, _ := d.version.GetInt64()
		return value > 1, nil
	case Uint64T:
		value, _ := d.version.GetUint64()
		return value > 1, nil
	}

	return false, fmt.Errorf("Parent: cannot retrieve version information")
}

// ==================================================
// Timestamp
// ==================================================

const (
	// TimestampPattern is the regexp for specific ISO 8601 datetime string
	TimestampPattern = `^[0-9]{4}` + `-` +
		`(?:1[0-2]|0[1-9])` + `-` +
		`(?:3[01]|0[1-9]|[12][0-9])` + `T` +
		`(?:2[0-3]|[01][0-9])` + `:` +
		`(?:[0-5][0-9])` + `:` +
		`(?:[0-5][0-9])` +
		`(?:Z|[+-](?:2[0-3]|[01][0-9]):(?:[0-5][0-9]))$`
)

// Timestamp is a data handler for a ISO 8601 timestamp string
type Timestamp struct {
	*DataBase

	ts string
}

var _ Data = (*Timestamp)(nil)

// NewTimestamp creates a ISO 8601 timestamp string handler
func NewTimestamp(key string, isRequired bool) *Timestamp {
	return &Timestamp{
		DataBase: NewDataBase(key, isRequired),
	}
}

// Prototype creates a prototype Timestamp
func (d *Timestamp) Prototype() Data {
	return &Timestamp{
		DataBase: d.DataBase.Prototype(),
	}
}

// Set the value of timestamp string
func (d *Timestamp) Set(data interface{}) error {
	if ts, ok := data.(string); ok {
		matched, err := regexp.MatchString(TimestampPattern, ts)
		if err != nil {
			return err
		}

		if !matched {
			return fmt.Errorf("Timestamp: string must in pattern " +
				"YYYY-MM-DDTHH:MM:SS(Z|±HH:MM)")
		}

		d.ts = ts
		return nil
	}

	return fmt.Errorf("Timestamp: 'string' is expected but '%T' is found", data)
}

// Encode Timestamp
func (d *Timestamp) Encode(m *map[string]interface{}) error {
	(*m)[d.GetKey()] = d.ts
	return nil
}

// Decode Timestamp
func (d *Timestamp) Decode(data interface{}, m *map[string]interface{}) error {
	ts, ok := data.(string)
	if !ok {
		return fmt.Errorf("Unknown error during decoding Timestamp: "+
			"'string' is expected but '%T' is found",
			data,
		)
	}

	matched, err := regexp.MatchString(TimestampPattern, ts)
	if err != nil {
		return err
	}

	if !matched {
		return fmt.Errorf("Timestamp: string must in pattern " +
			"YYYY-MM-DDTHH:MM:SS(Z|±HH:MM)")
	}

	d.ts = ts
	(*m)[d.GetKey()] = ts
	return nil
}

// ToJSON prepares the data for MarshalJSON
func (d *Timestamp) ToJSON(om *ordered.OrderedMap) error {
	om.Set(d.GetKey(), d.ts)
	return nil
}

// Resolve resolves the value
func (d *Timestamp) Resolve(path []string) (interface{}, []string, error) {
	if len(path) != 0 {
		return nil, nil, fmt.Errorf("Unexpected path elements past %s", path[0])
	}

	return d.ts, nil, nil
}

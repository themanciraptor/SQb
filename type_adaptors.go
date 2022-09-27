package sqb

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"time"
)

// Used to allow nullable fields in scanPair receivers
func nullType(k reflect.Kind) reflect.Type {
	switch k {
	case reflect.String:
		return reflect.TypeOf(NullString{})
	case reflect.Int64:
		return reflect.TypeOf(NullInt64{})
	case reflect.Int32:
		return reflect.TypeOf(NullInt32{})
	case reflect.Float64:
		return reflect.TypeOf(NullFloat64{})
	case reflect.Bool:
		return reflect.TypeOf(NullBool{})
	case reflect.Struct:
		return reflect.TypeOf(NullTime{})
	}

	return reflect.TypeOf(nil)
}

/*
	Wrapper around the sql NullTypes to allow us to scan directly to a receiver.
	We can still check the valid boolean.
*/

// Should not be instantiated directly, call NewNullString()
type NullString struct {
	String *string
	ns     *sql.NullString
	Valid  bool
}

func (n *NullString) Scan(value interface{}) error {
	err := n.ns.Scan(value)
	if err != nil {
		return err
	}

	*n.String = n.ns.String
	n.Valid = n.ns.Valid

	return nil
}

func (n NullString) Value() (driver.Value, error) {
	return n.ns.Value()
}

func NewNullString(receiver *string) *NullString {
	ns := &NullString{
		ns:     &sql.NullString{},
		Valid:  false,
		String: receiver,
	}

	return ns
}

// Should not be instantiated directly, call NewNullInt64()
type NullInt64 struct {
	Int64 *int64
	ns    *sql.NullInt64
	Valid bool
}

func (n *NullInt64) Scan(value interface{}) error {
	err := n.ns.Scan(value)
	if err != nil {
		return err
	}

	*n.Int64 = n.ns.Int64
	n.Valid = n.ns.Valid

	return nil
}

func (n NullInt64) Value() (driver.Value, error) {
	return n.ns.Value()
}

func NewNullInt64(receiver *int64) *NullInt64 {
	ns := &NullInt64{
		ns:    &sql.NullInt64{},
		Valid: false,
		Int64: receiver,
	}

	return ns
}

// Should not be instantiated directly, call NewNullInt32()
type NullInt32 struct {
	Int32 *int32
	ns    *sql.NullInt32
	Valid bool
}

func (n *NullInt32) Scan(value interface{}) error {
	err := n.ns.Scan(value)
	if err != nil {
		return err
	}

	*n.Int32 = n.ns.Int32
	n.Valid = n.ns.Valid

	return nil
}

func (n NullInt32) Value() (driver.Value, error) {
	return n.ns.Value()
}

func NewNullInt32(receiver *int32) *NullInt32 {
	ns := &NullInt32{
		ns:    &sql.NullInt32{},
		Valid: false,
		Int32: receiver,
	}

	return ns
}

// Should not be instantiated directly, call NewNullFloat64()
type NullFloat64 struct {
	Float64 *float64
	ns      *sql.NullFloat64
	Valid   bool
}

func (n *NullFloat64) Scan(value interface{}) error {
	err := n.ns.Scan(value)
	if err != nil {
		return err
	}

	*n.Float64 = n.ns.Float64
	n.Valid = n.ns.Valid

	return nil
}

func (n NullFloat64) Value() (driver.Value, error) {
	return n.ns.Value()
}

func NewNullFloat64(receiver *float64) *NullFloat64 {
	ns := &NullFloat64{
		ns:      &sql.NullFloat64{},
		Valid:   false,
		Float64: receiver,
	}

	return ns
}

// Should not be instantiated directly, call NewNullBool()
type NullBool struct {
	Bool  *bool
	ns    *sql.NullBool
	Valid bool
}

func (n *NullBool) Scan(value interface{}) error {
	err := n.ns.Scan(value)
	if err != nil {
		return err
	}

	*n.Bool = n.ns.Bool
	n.Valid = n.ns.Valid

	return nil
}

func (n NullBool) Value() (driver.Value, error) {
	return n.ns.Value()
}

func NewNullBool(receiver *bool) *NullBool {
	ns := &NullBool{
		ns:    &sql.NullBool{},
		Valid: false,
		Bool:  receiver,
	}

	return ns
}

// Should not be instantiated directly, call NewNullTime()
type NullTime struct {
	Time  *time.Time
	ns    *sql.NullTime
	Valid bool
}

func (n *NullTime) Scan(value interface{}) error {
	err := n.ns.Scan(value)
	if err != nil {
		return err
	}

	*n.Time = n.ns.Time
	n.Valid = n.ns.Valid

	return nil
}

func (n NullTime) Value() (driver.Value, error) {
	return n.ns.Value()
}

func NewNullTime(receiver *time.Time) *NullTime {
	ns := &NullTime{
		ns:    &sql.NullTime{},
		Valid: false,
		Time:  receiver,
	}

	return ns
}

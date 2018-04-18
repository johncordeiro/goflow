package types

import (
	"time"

	"github.com/nyaruka/goflow/utils"
)

// XDate is a date
type XDate struct {
	baseXPrimitive

	native time.Time
}

// NewXDate creates a new date
func NewXDate(value time.Time) XDate {
	return XDate{native: value}
}

// Reduce returns the primitive version of this type (i.e. itself)
func (x XDate) Reduce() XPrimitive { return x }

// ToXText converts this type to text
func (x XDate) ToXText() XText { return NewXText(utils.DateToISO(x.Native())) }

// ToXBoolean converts this type to a bool
func (x XDate) ToXBoolean() XBoolean { return NewXBoolean(!x.Native().IsZero()) }

// ToXJSON is called when this type is passed to @(json(...))
func (x XDate) ToXJSON() XText { return MustMarshalToXText(utils.DateToISO(x.Native())) }

// Native returns the native value of this type
func (x XDate) Native() time.Time { return x.native }

// Compare compares this date to another
func (x XDate) Compare(other XDate) int {
	switch {
	case x.Native().Before(other.Native()):
		return -1
	case x.Native().After(other.Native()):
		return 1
	default:
		return 0
	}
}

// MarshalJSON is called when a struct containing this type is marshaled
func (x XDate) MarshalJSON() ([]byte, error) {
	return x.Native().MarshalJSON()
}

// UnmarshalJSON is called when a struct containing this type is unmarshaled
func (x *XDate) UnmarshalJSON(data []byte) error {
	nativePtr := &x.native
	return nativePtr.UnmarshalJSON(data)
}

// XDateZero is the zero time value
var XDateZero = NewXDate(time.Time{})
var _ XPrimitive = XDateZero

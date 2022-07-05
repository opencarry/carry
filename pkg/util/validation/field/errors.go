/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package field

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	utilerrors "github.com/opencarry/carry/pkg/util/errors"
	"github.com/opencarry/carry/pkg/util/sets"
)

// Error is an implementation of the 'error' interface, which represents a
// field-level validation error.
type Error struct {
	Type     ErrorType
	Field    string
	BadValue interface{}
	Detail   string
}

var _ error = &Error{}

// Error implements the error interface.
func (v *Error) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.ErrorBody())
}

type omitValueType struct{}

var omitValue = omitValueType{}

// ErrorBody returns the error message without the field name.  This is useful
// for building nice-looking higher-level error reporting.
func (v *Error) ErrorBody() string {
	var s string
	switch {
	case v.Type == ErrorTypeRequired:
		s = v.Type.String()
	case v.Type == ErrorTypeForbidden:
		s = v.Type.String()
	case v.Type == ErrorTypeTooLong:
		s = v.Type.String()
	case v.Type == ErrorTypeInternal:
		s = v.Type.String()
	case v.BadValue == omitValue:
		s = v.Type.String()
	default:
		value := v.BadValue
		valueType := reflect.TypeOf(value)
		if value == nil || valueType == nil {
			value = "null"
		} else if valueType.Kind() == reflect.Ptr {
			if reflectValue := reflect.ValueOf(value); reflectValue.IsNil() {
				value = "null"
			} else {
				value = reflectValue.Elem().Interface()
			}
		}
		switch t := value.(type) {
		case int64, int32, float64, float32, bool:
			// use simple printer for simple types
			s = fmt.Sprintf("%s: %v", v.Type, value)
		case string:
			s = fmt.Sprintf("%s: %q", v.Type, t)
		case fmt.Stringer:
			// anything that defines String() is better than raw struct
			s = fmt.Sprintf("%s: %s", v.Type, t.String())
		default:
			// fallback to raw struct
			// TODO: internal types have panic guards against json.Marshalling to prevent
			// accidental use of internal types in external serialized form.  For now, use
			// %#v, although it would be better to show a more expressive output in the future
			s = fmt.Sprintf("%s: %#v", v.Type, value)
		}
	}
	if len(v.Detail) != 0 {
		s += fmt.Sprintf(": %s", v.Detail)
	}
	return s
}

// ErrorType is a machine-readable value providing more detail about why
// a field is invalid.  These values are expected to match 1-1 with
// CauseType in api/types.go.
type ErrorType string

const (
	// ErrorTypeNotFound is used to report failure to find a requested value
	// (e.g. looking up an ID).  See NotFound().
	ErrorTypeNotFound ErrorType = "field_value_not_found"
	// ErrorTypeRequired is used to report required values that are not
	// provided (e.g. empty strings, null values, or empty arrays).  See
	// Required().
	ErrorTypeRequired ErrorType = "field_value_required"
	// ErrorTypeDuplicate is used to report collisions of values that must be
	// unique (e.g. unique IDs).  See Duplicate().
	ErrorTypeDuplicate ErrorType = "field_value_duplicate"
	// ErrorTypeInvalid is used to report malformed values (e.g. failed regex
	// match, too long, out of bounds).  See Invalid().
	ErrorTypeInvalid ErrorType = "field_value_invalid"
	// ErrorTypeNotSupported is used to report unknown values for enumerated
	// fields (e.g. a list of valid values).  See NotSupported().
	ErrorTypeNotSupported ErrorType = "field_value_not_supported"
	// ErrorTypeForbidden is used to report valid (as per formatting rules)
	// values which would be accepted under some conditions, but which are not
	// permitted by the current conditions (such as security policy).  See
	// Forbidden().
	ErrorTypeForbidden ErrorType = "field_value_forbidden"
	// ErrorTypeTooLong is used to report that the given value is too long.
	// This is similar to ErrorTypeInvalid, but the error will not include the
	// too-long value.  See TooLong().
	ErrorTypeTooLong ErrorType = "field_value_too_long"
	// ErrorTypeTooMany is used to report "too many". This is used to
	// report that a given list has too many items. This is similar to FieldValueTooLong,
	// but the error indicates quantity instead of length.
	ErrorTypeTooMany ErrorType = "field_value_too_many"
	// ErrorTypeInternal is used to report other errors that are not related
	// to user input.  See InternalError().
	ErrorTypeInternal ErrorType = "internal_error"
	// ErrorTypeTypeInvalid is for the value did not match the schema type for that field
	ErrorTypeTypeInvalid ErrorType = "field_value_type_invalid"
)

// String converts a ErrorType into its corresponding canonical error message.
func (t ErrorType) String() string {
	switch t {
	case ErrorTypeNotFound:
		return "not found"
	case ErrorTypeRequired:
		return "required value"
	case ErrorTypeDuplicate:
		return "duplicate value"
	case ErrorTypeInvalid:
		return "invalid value"
	case ErrorTypeNotSupported:
		return "unsupported value"
	case ErrorTypeForbidden:
		return "forbidden"
	case ErrorTypeTooLong:
		return "too long"
	case ErrorTypeTooMany:
		return "too many"
	case ErrorTypeInternal:
		return "internal error"
	case ErrorTypeTypeInvalid:
		return "invalid value"
	default:
		panic(fmt.Sprintf("unrecognized validation error: %q", string(t)))
	}
}

// TypeInvalid returns a *Error indicating "type is invalid"
func TypeInvalid(field *Path, value interface{}, detail string) *Error {
	return &Error{ErrorTypeTypeInvalid, field.String(), value, detail}
}

// NotFound returns a *Error indicating "value not found".  This is
// used to report failure to find a requested value (e.g. looking up an ID).
func NotFound(field *Path, value interface{}) *Error {
	return &Error{ErrorTypeNotFound, field.String(), value, ""}
}

// Required returns a *Error indicating "value required".  This is used
// to report required values that are not provided (e.g. empty strings, null
// values, or empty arrays).
func Required(field *Path, detail string) *Error {
	return &Error{ErrorTypeRequired, field.String(), "", detail}
}

// Duplicate returns a *Error indicating "duplicate value".  This is
// used to report collisions of values that must be unique (e.g. names or IDs).
func Duplicate(field *Path, value interface{}) *Error {
	return &Error{ErrorTypeDuplicate, field.String(), value, ""}
}

// Invalid returns a *Error indicating "invalid value".  This is used
// to report malformed values (e.g. failed regex match, too long, out of bounds).
func Invalid(field *Path, value interface{}, detail string) *Error {
	return &Error{ErrorTypeInvalid, field.String(), value, detail}
}

// NotSupported returns a *Error indicating "unsupported value".
// This is used to report unknown values for enumerated fields (e.g. a list of
// valid values).
func NotSupported(field *Path, value interface{}, validValues []string) *Error {
	detail := ""
	if len(validValues) > 0 {
		quotedValues := make([]string, len(validValues))
		for i, v := range validValues {
			quotedValues[i] = strconv.Quote(v)
		}
		detail = "supported values: " + strings.Join(quotedValues, ", ")
	}
	return &Error{ErrorTypeNotSupported, field.String(), value, detail}
}

// Forbidden returns a *Error indicating "forbidden".  This is used to
// report valid (as per formatting rules) values which would be accepted under
// some conditions, but which are not permitted by current conditions (e.g.
// security policy).
func Forbidden(field *Path, detail string) *Error {
	return &Error{ErrorTypeForbidden, field.String(), "", detail}
}

// TooLong returns a *Error indicating "too long".  This is used to
// report that the given value is too long.  This is similar to
// Invalid, but the returned error will not include the too-long
// value.
func TooLong(field *Path, value interface{}, maxLength int) *Error {
	return &Error{ErrorTypeTooLong, field.String(), value, fmt.Sprintf("must have at most %d bytes", maxLength)}
}

// TooLongMaxLength returns a *Error indicating "too long".  This is used to
// report that the given value is too long.  This is similar to
// Invalid, but the returned error will not include the too-long
// value. If maxLength is negative, no max length will be included in the message.
func TooLongMaxLength(field *Path, value interface{}, maxLength int) *Error {
	var msg string
	if maxLength >= 0 {
		msg = fmt.Sprintf("may not be longer than %d", maxLength)
	} else {
		msg = "value is too long"
	}
	return &Error{ErrorTypeTooLong, field.String(), value, msg}
}

// TooMany returns a *Error indicating "too many". This is used to
// report that a given list has too many items. This is similar to TooLong,
// but the returned error indicates quantity instead of length.
func TooMany(field *Path, actualQuantity, maxQuantity int) *Error {
	var msg string

	if maxQuantity >= 0 {
		msg = fmt.Sprintf("must have at most %d items", maxQuantity)
	} else {
		msg = "has too many items"
	}

	var actual interface{}
	if actualQuantity >= 0 {
		actual = actualQuantity
	} else {
		actual = omitValue
	}

	return &Error{ErrorTypeTooMany, field.String(), actual, msg}
}

// InternalError returns a *Error indicating "internal error".  This is used
// to signal that an error was found that was not directly related to user
// input.  The errs argument must be non-nil.
func InternalError(field *Path, err error) *Error {
	return &Error{ErrorTypeInternal, field.String(), nil, err.Error()}
}

// ErrorList holds a set of Errors.  It is plausible that we might one day have
// non-field errors in this same umbrella package, but for now we don't, so
// we can keep it simple and leave ErrorList here.
type ErrorList []*Error

// NewErrorTypeMatcher returns an errors.Matcher that returns true
// if the provided error is a Error and has the provided ErrorType.
func NewErrorTypeMatcher(t ErrorType) utilerrors.Matcher {
	return func(err error) bool {
		if e, ok := err.(*Error); ok {
			return e.Type == t
		}
		return false
	}
}

// ToAggregate converts the ErrorList into an errors.Aggregate.
func (list ErrorList) ToAggregate() utilerrors.Aggregate {
	if len(list) == 0 {
		return nil
	}
	errs := make([]error, 0, len(list))
	errorMsgs := sets.NewString()
	for _, err := range list {
		msg := fmt.Sprintf("%v", err)
		if errorMsgs.Has(msg) {
			continue
		}
		errorMsgs.Insert(msg)
		errs = append(errs, err)
	}
	return utilerrors.NewAggregate(errs)
}

func fromAggregate(agg utilerrors.Aggregate) ErrorList {
	errs := agg.Errors()
	list := make(ErrorList, len(errs))
	for i := range errs {
		list[i] = errs[i].(*Error)
	}
	return list
}

// Filter removes items from the ErrorList that match the provided fns.
func (list ErrorList) Filter(fns ...utilerrors.Matcher) ErrorList {
	err := utilerrors.FilterOut(list.ToAggregate(), fns...)
	if err == nil {
		return nil
	}
	// FilterOut takes an Aggregate and returns an Aggregate
	return fromAggregate(err.(utilerrors.Aggregate))
}

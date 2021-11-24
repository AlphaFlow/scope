package util

import (
	"reflect"
	"strings"

	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
)

func IsBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// LookupForStructFieldTag returns the string value associated with a field and tag in a model, if it exists.
func LookupForStructFieldTag(model interface{}, field string, tagName string) (value string, ok bool) {
	_, t := StructValueAndType(model)

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)

		if sf.Name == field {
			return sf.Tag.Lookup(tagName)
		}
	}

	return "", false
}

// ValuesForStructTag returns a list of string values associated with a tag in a model.
func ValuesForStructTag(model interface{}, tagName string) []string {
	_, t := StructValueAndType(model)

	output := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		value := strings.Split(sf.Tag.Get(tagName), ",")[0]

		if value != "-" {
			output = append(output, value)
		}
	}

	return output
}

// FieldWithJsonTagValue returns the first field with a specific JSON tag.  It parses the JSON tag to match against the
// primary key in the json tag, unwinding custom json tag options.
// See documentation at https://golang.org/pkg/encoding/json/#Marshal
func FieldWithJsonTagValue(model interface{}, value string) (name string, ok bool) {
	v, t := StructValueAndType(model)

	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		if v, ok := f.Tag.Lookup("json"); ok {

			primaryTag := strings.Split(v, ",")[0]
			if primaryTag != "-" && primaryTag == value {
				return f.Name, true
			}
		}
	}

	return "", false
}

// StructValueAndType returns the reflected value and type of a model.
func StructValueAndType(model interface{}) (reflect.Value, reflect.Type) {
	v := reflect.ValueOf(model)
	v = reflect.Indirect(v)
	t := v.Type()

	return v, t
}

// InterfaceSlice takes an interface representing a known slice, and converts it into a []interface{}.
func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

// GetFieldByName gets the value of a field on the given structPtr.
func GetFieldByName(structPtr interface{}, fieldName string) reflect.Value {
	return reflect.ValueOf(structPtr).Elem().FieldByName(fieldName)
}

// UuidMust returns a new valid nulls uuid, or panics if its unable to.
func UuidMust() nulls.UUID {
	return nulls.NewUUID(uuid.Must(uuid.NewV4()))
}

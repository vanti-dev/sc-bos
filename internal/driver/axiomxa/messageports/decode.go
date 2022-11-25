package messageports

import (
	"bytes"
	"encoding"
	"fmt"
	"reflect"
	"strconv"
)

const Separator = ","

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "messageports: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Pointer {
		return "messageports: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "messageports: Unmarshal(nil " + e.Type.String() + ")"
}

// Unmarshal parses the message ports encoded data and stores the separate results in values pointed to by the entries of dst.
// If any entry in dst is nil or not a pointer then returns InvalidUnmarshalError.
//
// Unmarshal does not touch values in dst that do not have corresponding data.
// If unmarshalling into an interface value, Unmarshal stores a string in the interface value.
//
// Unmarshal splits data values using Separator.
func Unmarshal(data []byte, dst ...any) error {
	return UnmarshalSep(Separator, data, dst...)
}

// UnmarshalSep is like Unmarshal but using the given separator.
func UnmarshalSep(separator string, data []byte, dst ...any) error {
	// some simple validation so we don't fill some values and not others
	for _, v := range dst {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Pointer || rv.IsNil() {
			return &InvalidUnmarshalError{Type: reflect.TypeOf(v)}
		}
	}
	words := bytes.Split(data, []byte(separator))
	n := len(words)
	if n > len(dst) {
		n = len(dst)
	}

	for i := 0; i < n; i++ {
		v := dst[i]
		word := words[i]
		if u, ok := v.(encoding.BinaryUnmarshaler); ok {
			if err := u.UnmarshalBinary(words[i]); err != nil {
				return err
			}
			continue
		}
		if u, ok := v.(encoding.TextUnmarshaler); ok {
			if err := u.UnmarshalText(words[i]); err != nil {
				return err
			}
			continue
		}

		rv := reflect.ValueOf(dst[i])
		re := rv.Elem()
		switch re.Kind() {
		default:
			re.SetString(string(word))
		case reflect.Interface:
			re.Set(reflect.ValueOf(string(word)))
		case reflect.Bool:
			n, err := strconv.ParseBool(string(word))
			if err != nil {
				return fmt.Errorf("bool@%d", i)
			}
			re.SetBool(n)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(string(word), 10, 64)
			if err != nil || re.OverflowInt(n) {
				return fmt.Errorf("number@%d", i)
			}
			re.SetInt(n)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			n, err := strconv.ParseUint(string(word), 10, 64)
			if err != nil || re.OverflowUint(n) {
				return fmt.Errorf("number@%d", i)
			}
			re.SetUint(n)
		case reflect.Float32, reflect.Float64:
			n, err := strconv.ParseFloat(string(word), re.Type().Bits())
			if err != nil || re.OverflowFloat(n) {
				return fmt.Errorf("number@%d", i)
			}
			re.SetFloat(n)
		}
	}
	return nil
}

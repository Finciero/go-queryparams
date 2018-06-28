// Package query implements simple decoder that reads an URL query string and decodes its into a struct.
//
// A simple example:
//
// 	type Options struct {
// 		Foo   int `q:"foo"`
// 		Bar   string `q:"bar"`
// 	}
//
//  var p Options
// 	if err := query.NewDecoder("foo=2&bar=baz").Decode(&p); err != nil {
//		return err
//  }
//
// 	fmt.Print(p.Foo) // will output: 2
// 	fmt.Print(p.Bar) // will output: baz
//
// based on
package query

import (
	"encoding"
	"net/url"
	"reflect"
	"runtime"
	"strconv"
)

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "query: Decode(nil)"
	}
	if e.Type.Kind() != reflect.Ptr {
		return "query: Decode(non-pointer " + e.Type.String() + ")"
	}
	return "query: Decode(nil " + e.Type.String() + ")"
}

// UnimplementerError error types that are not implemented yet.
type UnimplementerError struct {
	Type reflect.Type
}

func (e *UnimplementerError) Error() string {
	return "query: " + e.Type.String() + " is not supported yet."
}

// A Decoder reads and decodes URL query strings.
type Decoder struct {
	q string
}

// NewDecoder returns a new decoder that read the given string.
func NewDecoder(s string) *Decoder {
	return &Decoder{s}
}

// Decode reads the query string from its input and stores it in the value pointed by v.
// Note that v should specify with the a "q" tag every exportable field that
// has a value in the query string.
func (d *Decoder) Decode(v interface{}) error {
	vals, err := url.ParseQuery(d.q)
	if err != nil || len(vals) == 0 {
		return err
	}
	return d.unmarshal(vals, v)
}

func (d *Decoder) unmarshal(src url.Values, v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}
	err = d.values(src, rv.Elem(), rv.Elem().Type())
	return
}

func (d *Decoder) values(src url.Values, dst reflect.Value, dstType reflect.Type) error {
	var (
		vals []string
		ok   bool
	)

	for i := 0; i < dst.NumField(); i++ {
		ft, fv := dstType.Field(i), dst.Field(i)

		if vals, ok = src[ft.Tag.Get(tagKey)]; !ok {
			continue
		}

		var addr = fv.Addr()
		if fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				fv.Set(reflect.New(fv.Type().Elem()))
			}
			addr = fv
			fv = fv.Elem()
		}

		if u, ok := addr.Interface().(encoding.TextUnmarshaler); ok {
			if vals[0] != "" {
				if err := u.UnmarshalText([]byte(vals[0])); err != nil {
					return err
				}
			}
			continue
		}

		switch fv.Kind() {
		case reflect.Slice, reflect.Array:
			n := len(vals)
			if fv.Kind() == reflect.Slice {
				fv.Set(reflect.MakeSlice(fv.Type(), n, n))
			}
			for j := 0; j < fv.Len() && j < n; j++ {
				if err := value(vals[j], fv.Index(j).Addr()); err != nil {
					return err
				}
			}
		case reflect.Struct, reflect.Map:
			return &UnimplementerError{reflect.TypeOf(fv)}
		default:
			if err := value(vals[0], addr); err != nil {
				return err
			}
		}
	}

	return nil
}

// dst must be a pointer in order to use this function
func value(src string, dst reflect.Value) (err error) {
	el := dst.Elem()
	switch el.Kind() {
	case reflect.String:
		el.SetString(src)
	case reflect.Bool:
		err = setBool(src, dst)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		err = setUint(src, dst)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = setInt(src, dst)
	case reflect.Float32, reflect.Float64:
		err = setFloat(src, dst)
	default:
		err = &UnimplementerError{reflect.TypeOf(el)}
	}
	return
}

func setUint(src string, dst reflect.Value) error {
	el := dst.Elem()
	if val, err := strconv.ParseUint(src, 10, el.Type().Bits()); err == nil {
		el.SetUint(val)
	} else {
		return err
	}
	return nil
}

func setInt(src string, dst reflect.Value) error {
	el := dst.Elem()
	if val, err := strconv.ParseInt(src, 10, el.Type().Bits()); err == nil {
		el.SetInt(val)
	} else {
		return err
	}
	return nil
}

func setBool(src string, dst reflect.Value) error {
	if src == "" {
		dst.Elem().SetBool(true)
		return nil
	}

	val, err := strconv.ParseBool(src)
	if err != nil {
		return err
	}
	dst.Elem().SetBool(val)
	return nil
}

func setFloat(src string, dst reflect.Value) error {
	el := dst.Elem()
	if val, err := strconv.ParseFloat(src, el.Type().Bits()); err == nil {
		el.SetFloat(val)
	} else {
		return err
	}
	return nil
}

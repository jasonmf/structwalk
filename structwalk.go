// structwalk provides a mechanism to recursively visit each field in a struct.
package structwalk

import (
	"reflect"
	"unicode"
)

// Recurse recursively decends through a struct value and its members. At
// each field it provides the reflect value and struct field it encountered to
// a user-defined function. If that function returns an error the error is
// passed backed back up the call stack. If the function returns true recursion
// will descend into that field. Along they way, pointer and interfaces fields
// are dereferenced.
//
// Where possible, recursion is based on the type, not the value. When a
// field's type is an interface it must try to step through the value. If
// multiple values of the same type have an interface field with values of
// different types the structure of the output may differ.
func Recurse(val interface{}, fn func(reflect.Value, reflect.StructField, string) (bool, error)) error {
	v, t := Prederef(val)
	return recurse(v, t, "", fn)
}

func recurse(val reflect.Value, t reflect.Type, name string, fn func(reflect.Value, reflect.StructField, string) (bool, error)) error {
	numFields := t.NumField()
field:
	for i := 0; i < numFields; i++ {
		f := t.Field(i)
		for _, r := range f.Name {
			if !unicode.IsUpper(r) {
				continue field
			}
			break
		}
		n := name + f.Name

		var v reflect.Value
		tk := f.Type.Kind()
		if val.IsValid() {
			v = val.Field(i)
		}
		vk := v.Kind()
		for vk == reflect.Ptr || tk == reflect.Interface {
			if tk == reflect.Interface {
				if !v.IsValid() {
					continue field
				}
				v = v.Elem()
				f.Type = v.Type()
			} else {
				f.Type = f.Type.Elem()
				v = v.Elem()
			}
			vk = v.Kind()
			tk = f.Type.Kind()
		}
		descend, err := fn(v, f, n)
		if err != nil {
			return err
		}
		if !descend {
			continue
		}
		if err := recurse(v, f.Type, n+".", fn); err != nil {
			return err
		}
	}
	return nil
}

// Prederef will iteratively dereference val as long as it's a pointer or
// interface type.
func Prederef(val interface{}) (reflect.Value, reflect.Type) {
	v := reflect.ValueOf(val)
	t := v.Type()
	vk := v.Kind()
	tk := t.Kind()
	for vk == reflect.Ptr || tk == reflect.Interface {
		if tk == reflect.Interface {
			v = v.Elem()
			t = v.Type()
		} else {
			t = t.Elem()
			v = v.Elem()

		}
		vk = v.Kind()
		tk = t.Kind()
	}
	return v, t
}

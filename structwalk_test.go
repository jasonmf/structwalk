package structwalk

import (
	"reflect"
	"testing"
	"time"
)

var (
	testTime, _  = time.Parse(time.RFC3339, "2019-02-28T14:40:34-08:00")
	typeTimeTime = reflect.TypeOf(time.Time{})
)

type Top struct {
	Name string
	Time time.Time
	Mid1 interface{}
	Mid2 Middle
	Mid3 *Middle
}

type Middle struct {
	Name   *string
	Bottom []int
}

func strPtr(s string) *string {
	return &s
}

func TestFlattenNames(t *testing.T) {
	t.Parallel()
	top := &Top{
		Name: "top",
		Time: testTime,
		Mid1: &Middle{
			Name:   strPtr("mid1"),
			Bottom: []int{1, 2, 3},
		},
		Mid2: Middle{
			Name:   strPtr("mid2"),
			Bottom: []int{4, 5, 6},
		},
		Mid3: &Middle{
			Name:   strPtr("mid3"),
			Bottom: []int{7, 8, 9},
		},
	}
	got, err := FlattenNames(top)
	if err != nil {
		t.Fatalf("want no error, got %q", err)
	}
	want := []string{"Name", "Time", "Mid1.Name", "Mid1.Bottom", "Mid2.Name", "Mid2.Bottom", "Mid3.Name", "Mid3.Bottom"}
	if len(want) != len(got) {
		t.Fatalf("got %q, want %q", got, want)
	}
	for i, f := range got {
		if want[i] != f {
			t.Fatalf("for field %d got %q, want %q", i, f, want[i])
		}
	}
}

func TestFlattenValues(t *testing.T) {
	t.Parallel()
	top := &Top{
		Name: "top",
		Time: testTime,
		Mid1: &Middle{
			Name:   strPtr("mid1"),
			Bottom: []int{1, 2, 3},
		},
		Mid2: Middle{
			Name:   strPtr("mid2"),
			Bottom: []int{4, 5, 6},
		},
		Mid3: &Middle{
			Name:   strPtr("mid3"),
			Bottom: []int{7, 8, 9},
		},
	}
	got, err := FlattenValues(top)
	if err != nil {
		t.Fatalf("want no error, got %q", err)
	}
	if len(got) != 8 {
		t.Fatalf("got %d values, want 8", len(got))
	}

	// check the strings
	for i, v := range map[int]string{
		0: "top",
		2: "mid1",
		4: "mid2",
		6: "mid3",
	} {
		s, ok := got[i].(string)
		if !ok {
			t.Fatalf("for field %d got %T, want string", i, got[i])
		}
		if s != v {
			t.Fatalf("for field %d got %q, want %q", i, s, v)
		}
	}

	if tt, ok := got[1].(time.Time); !ok {
		t.Fatalf("for field 1 got %T, want time.Time", got[1])
	} else {
		if !tt.Equal(testTime) {
			t.Fatalf("for field 1 got %v, want %v", tt, testTime)
		}
	}

	for i, v := range map[int][]int{
		3: []int{1, 2, 3},
		5: []int{4, 5, 6},
		7: []int{7, 8, 9},
	} {
		s, ok := got[i].([]int)
		if !ok {
			t.Fatalf("for field %d got %T, want []int", i, got[i])
		}
		if !sameIntSlice(s, v) {
			t.Fatalf("for field %d got %v, want %v", i, s, v)
		}
	}
}

func sameIntSlice(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func FlattenNames(val interface{}) ([]string, error) {
	fields := []string{}
	fp := &fields
	v, t := Prederef(val)
	if err := recurse(v, t, "", func(v reflect.Value, sf reflect.StructField, name string) (bool, error) {
		switch {
		case sf.Type == typeTimeTime:
			*fp = append(*fp, name)
			return false, nil
		case sf.Type.Kind() == reflect.Struct:
			return true, nil
		default:
			*fp = append(*fp, name)
			return false, nil
		}

	}); err != nil {
		return nil, err
	}
	return fields, nil
}

func FlattenValues(val interface{}) ([]interface{}, error) {
	values := []interface{}{}
	vp := &values
	v, t := Prederef(val)
	if err := recurse(v, t, "", func(v reflect.Value, sf reflect.StructField, name string) (bool, error) {
		switch {
		case sf.Type == typeTimeTime:
			var intf interface{}
			if v.IsValid() {
				intf = v.Interface()
			}
			*vp = append(*vp, intf)
			return false, nil
		case sf.Type.Kind() == reflect.Struct:
			return true, nil
		default:
			var intf interface{}
			if v.IsValid() {
				intf = v.Interface()
			}
			*vp = append(*vp, intf)
			return false, nil
		}

	}); err != nil {
		return nil, err
	}
	return values, nil
}

package query

import (
	"reflect"
	"testing"
)

func TestDecode_ArgumentTypes(t *testing.T) {
	dec := NewDecoder("foo=2&bar=baz&q=k1:k2")

	t.Run("v=nil", func(t *testing.T) {
		got := dec.Decode(nil)
		exp := &InvalidUnmarshalError{reflect.TypeOf(nil)}

		if !reflect.DeepEqual(got, exp) {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("v=non-pointer", func(t *testing.T) {
		var test = struct{ Foo uint }{}
		got := dec.Decode(test)
		exp := &InvalidUnmarshalError{reflect.TypeOf(test)}

		if !reflect.DeepEqual(got, exp) {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("v=nil pointer", func(t *testing.T) {
		var test struct{ Foo uint }
		got := dec.Decode(test)
		exp := &InvalidUnmarshalError{reflect.TypeOf(test)}

		if !reflect.DeepEqual(got, exp) {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("v=pointer", func(t *testing.T) {
		var test = new(struct{ Foo uint })
		got := dec.Decode(test)

		if !reflect.DeepEqual(got, nil) {
			t.Fatalf("exp: %v\ngot: %v", nil, got)
		}
	})
}

func TestDecode_OneLevel(t *testing.T) {
	t.Run("field=integer", func(t *testing.T) {
		var test struct {
			Numeric int `q:"numeric"`
		}
		ok(t, NewDecoder("numeric=2").Decode(&test))
		exp := 2
		got := test.Numeric
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=pointer to integer", func(t *testing.T) {
		var test struct {
			Numeric *int `q:"numeric"`
		}
		ok(t, NewDecoder("numeric=2").Decode(&test))
		exp := 2
		got := *test.Numeric
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=float64", func(t *testing.T) {
		var test struct {
			Float float64 `q:"float"`
		}
		ok(t, NewDecoder("float=3.45").Decode(&test))
		exp := 3.45
		got := test.Float
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=pointer to float64", func(t *testing.T) {
		var test struct {
			Float *float64 `q:"float"`
		}
		ok(t, NewDecoder("float=3.45").Decode(&test))
		exp := 3.45
		got := *test.Float
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=string", func(t *testing.T) {
		var test struct {
			Text string `q:"text"`
		}
		ok(t, NewDecoder("text=this is a text&foo=bar").Decode(&test))
		exp := "this is a text"
		got := test.Text
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=pointer to string", func(t *testing.T) {
		var test struct {
			Text *string `q:"text"`
		}
		ok(t, NewDecoder("text=this is a text&foo=bar").Decode(&test))
		exp := "this is a text"
		got := *test.Text
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=bool", func(t *testing.T) {
		var test struct {
			Bv1 bool `q:"bv1"`
			Bv2 bool `q:"bv2"`
			Bv3 bool `q:"bv3"`
		}
		ok(t, NewDecoder("bv1&bv2=true&bv3=1").Decode(&test))
		exp := true
		got := test.Bv1
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
		got = test.Bv2
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
		got = test.Bv3
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=pointer to bool", func(t *testing.T) {
		var test struct {
			Bv1 *bool `q:"bv1"`
			Bv2 *bool `q:"bv2"`
			Bv3 *bool `q:"bv3"`
		}
		ok(t, NewDecoder("bv1&bv2=true&bv3=1").Decode(&test))
		exp := true
		got := *test.Bv1
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
		got = *test.Bv2
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
		got = *test.Bv3
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=array (same lenght)", func(t *testing.T) {
		var test struct {
			Array [2]int `q:"arr"`
		}
		ok(t, NewDecoder("arr=1&arr=2").Decode(&test))
		exp := [...]int{1, 2}
		got := test.Array
		if !reflect.DeepEqual(exp, got) {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=array (greater lenght)", func(t *testing.T) {
		var test struct {
			Array [4]int `q:"arr"`
		}
		ok(t, NewDecoder("arr=1&arr=2").Decode(&test))
		exp := [...]int{1, 2, 0, 0}
		got := test.Array
		if !reflect.DeepEqual(exp, got) {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=array (lesser lenght)", func(t *testing.T) {
		var test struct {
			Array [1]int `q:"arr"`
		}
		ok(t, NewDecoder("arr=1&arr=2").Decode(&test))
		exp := [...]int{1}
		got := test.Array
		if !reflect.DeepEqual(exp, got) {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=pointer to array", func(t *testing.T) {
		var test struct {
			Array *[1]int `q:"arr"`
		}
		ok(t, NewDecoder("arr=1&arr=2").Decode(&test))
		exp := [...]int{1}
		got := *test.Array
		if !reflect.DeepEqual(exp, got) {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=slice", func(t *testing.T) {
		var test struct {
			Slice []int `q:"arr"`
		}
		ok(t, NewDecoder("arr=1&arr=2").Decode(&test))
		exp := []int{1, 2}
		got := test.Slice
		if !reflect.DeepEqual(exp, got) {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=pointer to slice", func(t *testing.T) {
		var test struct {
			Slice *[]int `q:"arr"`
		}
		ok(t, NewDecoder("arr=1&arr=2").Decode(&test))
		exp := []int{1, 2}
		got := *test.Slice
		if !reflect.DeepEqual(exp, got) {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})
}

func TestDecode_OneLevelOverrides(t *testing.T) {
	t.Run("field=integer", func(t *testing.T) {
		test := struct {
			Numeric int `q:"numeric"`
		}{
			Numeric: 3,
		}
		ok(t, NewDecoder("numeric=2").Decode(&test))
		exp := 2
		got := test.Numeric
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=pointer to integer", func(t *testing.T) {
		exp := 2
		got := 3

		test := struct {
			Numeric *int `q:"numeric"`
		}{
			Numeric: &got,
		}

		ok(t, NewDecoder("numeric=2").Decode(&test))
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=float64", func(t *testing.T) {
		test := struct {
			Float float64 `q:"float"`
		}{
			Float: 4.56,
		}

		ok(t, NewDecoder("float=3.45").Decode(&test))
		exp := 3.45
		got := test.Float
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=pointer to float64", func(t *testing.T) {
		exp := 3.45
		got := 4.56

		test := struct {
			Float *float64 `q:"float"`
		}{
			Float: &got,
		}

		ok(t, NewDecoder("float=3.45").Decode(&test))
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=string", func(t *testing.T) {
		test := struct {
			Text string `q:"text"`
		}{
			Text: "this should be overwritten",
		}

		ok(t, NewDecoder("text=this is a text&foo=bar").Decode(&test))
		exp := "this is a text"
		got := test.Text
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=pointer to string", func(t *testing.T) {
		exp := "this is a text"
		got := "this should be overwritten"

		test := struct {
			Text *string `q:"text"`
		}{
			Text: &got,
		}

		ok(t, NewDecoder("text=this is a text&foo=bar").Decode(&test))
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=bool", func(t *testing.T) {
		test := struct {
			Bv1 bool `q:"bv1"`
			Bv2 bool `q:"bv2"`
			Bv3 bool `q:"bv3"`
		}{}

		ok(t, NewDecoder("bv1&bv2=true&bv3=1").Decode(&test))
		exp := true
		got := test.Bv1
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
		got = test.Bv2
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
		got = test.Bv3
		if exp != got {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=pointer to bool", func(t *testing.T) {
		var bv1, bv2, bv3 bool
		test := struct {
			Bv1 *bool `q:"bv1"`
			Bv2 *bool `q:"bv2"`
			Bv3 *bool `q:"bv3"`
		}{Bv1: &bv1, Bv2: &bv2, Bv3: &bv3}

		ok(t, NewDecoder("bv1&bv2=true&bv3=1").Decode(&test))
		if !bv1 {
			t.Fatalf("exp: %v\ngot: %v", true, bv1)
		}
		if !bv2 {
			t.Fatalf("exp: %v\ngot: %v", true, bv2)
		}
		if !bv3 {
			t.Fatalf("exp: %v\ngot: %v", true, bv3)
		}
	})

	t.Run("field=array (same lenght)", func(t *testing.T) {
		arr := [...]int{4, 5}

		test := struct {
			Array [2]int `q:"arr"`
		}{
			Array: arr,
		}

		ok(t, NewDecoder("arr=1&arr=2").Decode(&test))
		exp := [...]int{1, 2}
		got := test.Array
		if !reflect.DeepEqual(exp, got) {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=array (greater lenght)", func(t *testing.T) {
		arr := [...]int{4, 5, 6, 7}

		test := struct {
			Array [4]int `q:"arr"`
		}{
			Array: arr,
		}

		ok(t, NewDecoder("arr=1&arr=2").Decode(&test))
		exp := [...]int{1, 2, 6, 7}
		got := test.Array
		if !reflect.DeepEqual(exp, got) {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=array (lesser lenght)", func(t *testing.T) {
		arr := [...]int{1}
		test := struct {
			Array [1]int `q:"arr"`
		}{
			Array: arr,
		}
		ok(t, NewDecoder("arr=1&arr=2").Decode(&test))
		exp := [...]int{1}
		got := test.Array
		if !reflect.DeepEqual(exp, got) {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=pointer to array", func(t *testing.T) {
		got := [...]int{5}
		exp := [...]int{1}
		test := struct {
			Array *[1]int `q:"arr"`
		}{
			Array: &got,
		}
		ok(t, NewDecoder("arr=1&arr=2").Decode(&test))
		if !reflect.DeepEqual(exp, got) {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=slice", func(t *testing.T) {
		test := struct {
			Slice []int `q:"arr"`
		}{
			Slice: []int{4, 5, 6},
		}
		ok(t, NewDecoder("arr=1&arr=2").Decode(&test))
		exp := []int{1, 2}
		got := test.Slice
		if !reflect.DeepEqual(exp, got) {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})

	t.Run("field=pointer to slice", func(t *testing.T) {
		got := []int{4, 5, 6}
		exp := []int{1, 2}
		test := struct {
			Slice *[]int `q:"arr"`
		}{
			Slice: &got,
		}
		ok(t, NewDecoder("arr=1&arr=2").Decode(&test))
		if !reflect.DeepEqual(exp, got) {
			t.Fatalf("exp: %v\ngot: %v", exp, got)
		}
	})
}

func ok(t testing.TB, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

package md5hash

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"iter"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/ettle/strcase"
)

type table interface {
	TableName() string
}

// Query returns SQL SELECT statement to get the hash of table data based on the fields
// provided in the struct `v`.
//
// Format of Query is `SELECT md5_chain(<col1> || <col2> || ...) FROM <table_name>`
func Query[T table](v T) (string, error) {
	n := v.TableName()
	if n == "" {
		return "", errors.New("table name is empty")
	}

	it, err := getItems(reflect.ValueOf(v), true)
	if err != nil {
		return "", err
	}

	if len(it) == 0 {
		return "", errors.New("type contains no hashable fields")
	}

	arr := make([]string, 0, len(it))
	for _, i := range it {
		arr = append(arr, "\""+i.name+"\"") // quoted because of non lowercased column names
	}
	return fmt.Sprintf("SELECT md5_chain(%s) FROM %s", strings.Join(arr, "||"), n), nil

}

// Value returns the concatenated fields for the given struct that can be used to compute the hash values.
// If the passed value is not a struct, it's string value is returned.
func Value(x any) (string, error) {
	v := reflect.ValueOf(x)
	v, _ = elem(v) // return empty string for nil values
	//if !ok {
	//	return "", nil
	//}

	if v.Kind() != reflect.Struct {
		if fn, ok := x.(fmt.Stringer); ok {
			return fn.String(), nil
		}
		return strValue(v), nil
	}
	it, err := getItems(v, false)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	for _, i := range it {
		b.WriteString(i.value)
	}

	return b.String(), nil
}

// elem returns value of pointers and true.
// It otherwise creates zero values of nil pointers and returns false.
//
// It differs from reflect.Indirect in that, it creates zero value
// of the ELEMENT if v is a nil pointer instead of creating zero value of pointers.
func elem(v reflect.Value) (reflect.Value, bool) {
	if v.Type().Kind() == reflect.Ptr {
		if v.IsNil() {
			return reflect.Zero(v.Type().Elem()), false
		}
		v = v.Elem()
	}
	return v, true
}

func getItems(v reflect.Value, ignoreNilPtr bool) ([]item, error) {
	v, ok := elem(v)
	if !ok && !ignoreNilPtr {
		return nil, nil
	}

	var it []item
	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)

		if f.Anonymous {
			innerIt, err := getItems(v.Field(i), ignoreNilPtr)
			if err != nil {
				return nil, err
			}
			it = append(it, innerIt...)
			continue
		}

		if !f.IsExported() {
			continue
		}

		md, ok := f.Tag.Lookup("hash")
		if !ok {
			it = append(it, item{name: strcase.ToSnake(f.Name), value: strValue(v.Field(i)), pos: len(it)})
			continue
		}

		mdd := strings.Split(md, ",")
		if len(mdd) != 2 {
			return nil, fmt.Errorf("invalid hash tag, should be in format <index>,<name>")
		}
		pos, err := strconv.Atoi(mdd[0])
		if err != nil {
			return nil, fmt.Errorf("invalid struct tag for field %v", f.Name)
		}
		it = append(it, item{
			name:  mdd[1],
			value: strValue(v.Field(i)),
			pos:   pos,
		})
	}

	slices.SortStableFunc(it, func(a, b item) int {
		return a.pos - b.pos
	})

	return it, nil
}

func strValue(v reflect.Value) string {

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		if bytes, ok := v.Interface().([]byte); ok {
			return string(bytes)
		}
		fallthrough
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}

// HashValue returns the md5 hash value of the object v
func HashValue(v any) (string, error) {
	h := md5.New()
	key, err := Value(v)
	if err != nil {
		return "", err
	}
	if _, err = io.WriteString(h, key); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// ChainHashValue is the same as ChainHashValueIter
func ChainHashValue[T any](rows []T) (string, error) {
	return ChainHashValueIter(slices.Values(rows))
}

// ChainHashValueIter returns the running md5 hash value of list of items
//
// f(x) = f(x1, f(x2, f(x3, ...)))
func ChainHashValueIter[T any](it iter.Seq[T]) (string, error) {
	prev := ""
	h := md5.New()

	for i := range it {
		if _, err := io.WriteString(h, prev); err != nil {
			return "", err
		}

		v, err := Value(i)
		if err != nil {
			return "", err
		}

		if _, err = io.WriteString(h, v); err != nil {
			return "", err
		}
		prev = fmt.Sprintf("%x", h.Sum(nil))
		h.Reset()
	}

	return prev, nil
}

type item struct {
	// name is the db column name
	name string
	// value is the string representation to be used for hashing
	value string
	// pos is the position of this field when concating values.
	pos int
}

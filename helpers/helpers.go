package helpers

import (
	"errors"
	"reflect"
)

// ObjectKeys takes an interface of type struct and returns
// a slice containing all the key names for the given struct
func ObjectKeys(i interface{}) (r []string, err error) {
	v := reflect.Indirect(reflect.ValueOf(i))
	if v.Kind() != reflect.Struct {
		err = errors.New("this function only accepts struct types or pointer-to-struct types")
		return r, err
	}

	for i := 0; i < v.NumField(); i++ {
		ft := v.Type().Field(i)
		fv := v.Field(i)

		if fv.Kind() == reflect.Ptr {
			fv = fv.Elem()
		}

		r = append(r, ft.Name)
	}

	return r, err
}

// ObjectKeysFlatten takes an interface of type struct and returns a slice
// containing all the key names for the given struct and it's nested structs
func ObjectKeysFlatten(i interface{}) (r []string, err error) {
	v := reflect.Indirect(reflect.ValueOf(i))
	k := v.Kind()

	if k != reflect.Struct {
		err = errors.New("this function only accepts struct types or pointer-to-struct types")
		return r, err
	}

	for i := 0; i < v.NumField(); i++ {
		ft := v.Type().Field(i)

		fv := v.Field(i)

		if fv.Kind() == reflect.Ptr {
			fv = fv.Elem()
		}

		r = append(r, ft.Name)

		if fv.Kind() == reflect.Struct {
			fi := fv.Interface()
			sl, _ := ObjectKeysFlatten(fi)
			r = append(r, sl...)
		}
	}

	return r, err
}

// Get takes an interface of type struct and a fieldname name string.
// It returns the value from the structs field with the corresponding name.
func Get(i interface{}, fn string) (interface{}, error) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	// check the interfaces kind - we only wany structs and pointers to structs
	if v.Kind() != reflect.Struct {
		return nil, errors.New("The provided interface must be a struct type")
	}
	// check the provided fieldname exists
	sf, ok := t.FieldByName(fn)
	if !ok {
		return nil, errors.New("Sorry, your key doesn't match any fields in the provided struct type")
	}
	// use the index position of the corresponding fieldname to provide the value
	fv := v.Field(sf.Index[0])
	if !fv.CanInterface() {
		return nil, errors.New("the value you're trying to fetch can not be made into an interface")
	}
	return fv.Interface(), nil
}

// Set takes an interface of type struct and a key name string.
// It sets the value at the given field or returns an error if not.
func Set(i interface{}, k string, val interface{}) (err error) {
	t := reflect.TypeOf(i)
	v := reflect.Indirect(reflect.ValueOf(i))
	// check the interfaces kind - we only wany structs and pointers to structs
	if v.Kind() != reflect.Struct {
		return errors.New("The provided interface must be a struct type")
	}
	sf, ok := t.FieldByName(k)
	if !ok {
		return errors.New("Sorry, your key doesn't match any fields in the provided struct type")
	}
	fv := v.Field(sf.Index[0])
	if !fv.CanInterface() {
		return errors.New("Sorry, we can't interface the field value you're trying to set")
	}
	if !fv.CanSet() {
		return errors.New("Sorry, we can't set the value of that field")
	}
	fv.Set(reflect.ValueOf(val))
	return nil
}

// GetVals takes an interface and returns every value from that interface
func GetVals(i interface{}) (r []interface{}, err error) {

	// get the field names
	names, err := ObjectKeys(i)
	if err != nil {
		return []interface{}{}, err
	}

	if len(names) > 0 {
		for _, fn := range names {
			val, err := Get(i, fn)
			if err != nil {
				return []interface{}{}, err
			}
			r = append(r, val)
		}
	}

	return r, err
}

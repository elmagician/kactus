package interfaces

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"go.uber.org/zap"
)

var (
	// ErrInvalidPathElement indicates an element from interface structure is not valid to go into array/slice.
	ErrInvalidPathElement = errors.New("path element is not an integer on Array/Slice")
)

// GetFieldFromPath help to recover data from interface
// It will go through the interface using a `.` separated string list
// when element is a structure, the path is expected to be a camel cased key (a call to strconv.ToCamel is done)
// when element is an array or a slice, the path is expected as an int
// when element is a map, the path is expected as a key string in the slice
// Warning: we currently do not support map key witch are not strings
func GetFieldFromPath(obj interface{}, path string) (reflect.Value, bool) {
	splittedPath := strings.Split(path, ".")
	ref := reflect.ValueOf(obj)

	for idx, p := range splittedPath {
		if ref.Kind() == reflect.Ptr {
			if ref.IsNil() {
				return ref, idx+1 != len(splittedPath)
			}

			ref = ref.Elem()
		}

		//nolint:exhaustive
		switch ref.Kind() {
		case reflect.Map:
			pVal := reflect.ValueOf(p)
			if !hasKey(pVal, ref.MapKeys()) {
				return ref, false
			}

			ref = ref.MapIndex(pVal)
			ref = reflect.ValueOf(ref.Interface()) // ensure we got correct reflect kind for map value
		case reflect.Slice, reflect.Array:
			intP, err := strconv.Atoi(p)
			if err != nil {
				log.Fatal(ErrInvalidPathElement.Error())
			}

			ref = reflect.ValueOf(ref.Index(intP).Interface())
		case reflect.Struct:
			p = strcase.ToCamel(p)
			ref = ref.FieldByName(p)
		default:
			return ref, false
		}
	}

	return ref, true
}

func GenerateFieldList(obj interface{}) (res []string) {
	log.Debug("generating fields for", zap.Reflect("object", obj))

	ref, ok := obj.(reflect.Value)
	if !ok {
		log.Debug("Getting reflect.Value for object")

		ref = reflect.ValueOf(obj)
	}

	//nolint:exhaustive
	switch ref.Kind() {
	case reflect.Interface:
		log.Debug("RECURSION: looking for ref as interface value")
		return GenerateFieldList(ref.Interface())
	case reflect.Ptr:
		log.Debug("object is a Pointer")

		if ref.IsNil() {
			log.Debug("EXITING: nil pointer.")
			return nil
		}

		log.Debug("RECURSION: looking for fields from pointer value.")

		res = GenerateFieldList(ref.Elem())
	case reflect.Map:
		log.Debug("object is a map")

		for _, key := range ref.MapKeys() {
			log.Debug("RECURSION: looking for fields from map key.", zap.String("key", key.String()))

			fields := GenerateFieldList(ref.MapIndex(key))
			res = append(res, appendPath(key.String(), fields)...)

			log.Debug("updated result", zap.Reflect("new fields", res))
		}
	case reflect.Slice, reflect.Array:
		log.Debug("object is a slice/array")

		for i := 0; i < ref.Len(); i++ {
			log.Debug("RECURSION: looking for fields from index.", zap.Int("ind", i))

			fields := GenerateFieldList(ref.Index(i))
			res = append(res, appendPath(strconv.Itoa(i), fields)...)

			log.Debug("updated result", zap.Reflect("new fields", res))
		}
	case reflect.Struct:
		log.Debug("object is a structure")

		for i := 0; i < ref.NumField(); i++ {
			name := ref.Type().Field(i).Name
			log.Debug("RECURSION: looking for fields from struct field.", zap.String("field", name))

			fields := GenerateFieldList(ref.Field(i))

			res = append(res, appendPath(name, fields)...)

			log.Debug("updated result", zap.Reflect("new fields", res))
		}
	}

	log.Debug("returning new res", zap.Reflect("res", res))

	return res
}

func hasKey(key reflect.Value, keys []reflect.Value) bool {
	for _, k := range keys {
		if k.IsValid() && k.Interface() == key.Interface() {
			return true
		}
	}

	return false
}

func appendPath(base string, elements []string) []string {
	var res []string

	if len(elements) == 0 {
		return []string{base}
	}

	for _, e := range elements {
		res = append(res, base+"."+e)
	}

	return res
}

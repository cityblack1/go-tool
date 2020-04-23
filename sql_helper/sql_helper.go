package sql_helper

import (
	"reflect"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

func GetDbFieldsName(v interface{}) (ret []string) {
	var val reflect.Value
	if reflect.ValueOf(v).Kind() == reflect.Ptr {
		val = reflect.ValueOf(v).Elem()
	} else {
		val = reflect.ValueOf(v)
	}
	for i := 0; i < val.Type().NumField(); i++ {
		tag := strings.Split(val.Type().Field(i).Tag.Get("db"), ",")[0]
		if tag != "" && tag != "-" {
			ret = append(ret, tag)
		}
	}
	return ret
}

func GetDbFieldsAddr(v interface{}) (ret []interface{}) {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		panic(errors.New("GetDbFieldsAddr accept pointer only!"))
	}
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.Type().NumField(); i++ {
		tag := val.Type().Field(i).Tag.Get("db")
		if tag != "" && tag != "-" {
			ret = append(ret, val.Field(i).Addr().Interface())
		}
	}
	return ret
}

func GetDbEqualMap(v interface{}) squirrel.Eq {
	eq := make(squirrel.Eq)
	keys := GetDbFieldsName(v)
	vals := GetDbFieldsAddr(v)
	for i := range keys {
		v := reflect.ValueOf(vals[i])
		if vals[i] == nil || v.IsNil() {
			continue
		}
		if v.Elem().Kind() == reflect.Ptr && v.Elem().IsNil() {
			continue
		}

		eq[keys[i]] = vals[i]
	}
	return eq
}

func getElem(v interface{}) interface{} {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		return getElem(val.Elem().Interface())
	} else {
		return val.Interface()
	}
}

func GetDbElemEqualMap(v interface{}) squirrel.Eq {
	eq := GetDbEqualMap(v)
	for k, v := range eq {
		eq[k] = getElem(v)
	}
	return eq
}

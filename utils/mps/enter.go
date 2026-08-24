package mps

import "reflect"

func Struct2Map(data any, t string) (mp map[string]any) {
	mp = make(map[string]any)
	val := reflect.ValueOf(data)
	typ := val.Type()
	fieldsCount := val.NumField()
	for i := 0; i < fieldsCount; i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)
		tagValue := fieldType.Tag.Get(t)

		if tagValue == "" || tagValue == "-" {
			continue
		}

		kind := fieldVal.Kind()
		if kind == reflect.Ptr || kind == reflect.Slice || kind == reflect.Map ||
			kind == reflect.Chan || kind == reflect.Func || kind == reflect.Interface {
			if fieldVal.IsNil() {
				continue
			}
		}

		if kind == reflect.Ptr {
			mp[tagValue] = fieldVal.Elem().Interface()
		} else {
			mp[tagValue] = fieldVal.Interface()
		}
	}
	return
}

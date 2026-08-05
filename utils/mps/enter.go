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

		if tagValue == "" {
			continue
		}
		if tagValue == "-" {
			continue
		}
		if fieldVal.IsNil() {
			continue
		}
		if fieldVal.Kind() == reflect.Ptr {
			mp[tagValue] = fieldVal.Elem().Interface()
		} else {
			mp[tagValue] = fieldVal.Interface()
		}
	}
	return
}

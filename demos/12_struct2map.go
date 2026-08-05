package main

import (
	"fmt"
	"reflect"
)

func main() {
	username := "zhaoliu"
	likeTags := []string{"js", "go"}

	var cr = UserInfoUpdateRequest{
		Username: &username,
		LikeTags: &likeTags,
	}

	mp1 := Struct2Map(cr, "s-u")
	mp2 := Struct2Map(cr, "s-u-c")
	fmt.Printf("%+v \n %+v \n", mp1, mp2)
}

type UserInfoUpdateRequest struct {
	Username    *string   `json:"username" s-u:"username"`
	Nickname    *string   `json:"nickname" s-u:"nickname"`
	Avatar      *string   `json:"avatar" s-u:"avatar"`
	Abstract    *string   `json:"abstract" s-u:"abstract"`
	LikeTags    *[]string `json:"likeTags" s-u-c:"like_tags"`
	OpenCollect *bool     `json:"openCollect" s-u-c:"open_collect"`
	OpenFollow  *bool     `json:"openFollow" s-u-c:"open_follow"`
	OpenFans    *bool     `json:"openFans" s-u-c:"open_fans"`
	HomeStyleID *uint     `json:"homeStyleId" s-u-c:"home_style_id"`
}

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

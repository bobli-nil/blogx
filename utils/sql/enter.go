package sql

import "fmt"

func ConvertSliceSql(list []uint) (s string) {
	s += "("
	for i, u := range list {
		s += fmt.Sprintf("%d", u)
		if i != len(list)-1 {
			s += ","
		}
	}
	s += ")"
	return
}

func ConvertSliceOrderSql(list []uint) (s string) {
	for i, u := range list {
		if i == len(list)-1 {
			s += fmt.Sprintf("id = %d desc", u)
		} else {
			s += fmt.Sprintf("id = %d desc,", u)
		}
	}
	return
}

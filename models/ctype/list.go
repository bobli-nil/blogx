package ctype

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type List []string

func (j *List) Scan(value any) error {
	val, ok := value.([]uint8)
	if ok {
		*j = strings.Split(string(val), ",")
		return nil
	}
	return errors.New("类型错误")
}

func (j List) Value() (driver.Value, error) {
	res := strings.Join(j, ",")
	return res, nil
}

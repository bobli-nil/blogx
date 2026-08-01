package file

import (
	"blogx_server/global"
	"blogx_server/utils"
	"errors"
	"strings"
)

func ImageSuffixJudge(filename string) (suffix string, err error) {
	_list := strings.Split(filename, ".")
	if len(_list) < 2 {
		err = errors.New("文件名不合法")
		return
	}
	suffix = _list[len(_list)-1]
	if !utils.InList(suffix, global.Conf.Upload.WhiteList) {
		err = errors.New("文件后缀不正确")
		return
	}
	return
}

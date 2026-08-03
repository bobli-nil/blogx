package main

import (
	"blogx_server/global"
	"fmt"
	"image/color"

	"github.com/mojocn/base64Captcha"
)

func main() {
	driver := base64Captcha.NewDriverString(80, 240, 10, 2, 4, "1234567890qwertyuioplkjhgfdsazxcvbnm", &color.RGBA{R: 240, G: 240, B: 240, A: 255}, nil, nil)
	captcha := base64Captcha.NewCaptcha(driver, global.Store)
	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		fmt.Println("错误", err)
		return
	}
	fmt.Println(id)
	fmt.Println(b64s)
	fmt.Println(answer)
}

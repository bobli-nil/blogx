package captcha_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"fmt"
	"image/color"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

type CaptchaApi struct{}

type CaptchaResponse struct {
	CaptchaID string `json:"captchaId"`
	Captcha   string `json:"captcha"`
}

func (CaptchaApi) CaptchaCreateView(c *gin.Context) {
	driver := base64Captcha.NewDriverString(80, 240, 10, 2, 4, "1234567890qwertyuioplkjhgfdsazxcvbnm", &color.RGBA{R: 240, G: 240, B: 240, A: 255}, nil, nil)
	captcha := base64Captcha.NewCaptcha(driver, global.Store)
	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	fmt.Println(answer)
	global.Store.Set(id, answer)
	res.OkWithData(&CaptchaResponse{
		CaptchaID: id,
		Captcha:   b64s,
	}, c)

}

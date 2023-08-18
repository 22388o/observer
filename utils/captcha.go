package utils

import (
	"github.com/mojocn/base64Captcha"
	"log"
)

var captchaDrive *base64Captcha.DriverDigit
var captchaInstance *base64Captcha.Captcha
var captchaStore base64Captcha.Store

func InitCapcha() {
	captchaDrive = &base64Captcha.DriverDigit{Length: 6,
		Height:   80,
		Width:    240,
		MaxSkew:  0.7,
		DotCount: 100}
	captchaStore = base64Captcha.DefaultMemStore
	captchaInstance = base64Captcha.NewCaptcha(captchaDrive, captchaStore)
}
func GenCaptcha(onlyAnwser bool) (id, answer, b64sImg string) {
	//id, b64s, err := captchaInstance.Generate()
	var content string
	id, content, answer = captchaDrive.GenerateIdQuestionAnswer()
	err := captchaStore.Set(id, answer)
	if err != nil {
		log.Fatalln(err)
	}
	if !onlyAnwser {

		item, err := captchaDrive.DrawCaptcha(content)
		if err != nil {
			log.Fatalln(err)
		}
		b64sImg = item.EncodeB64string()
	}
	return
}
func VerifyCaptcha(id, value string) bool {
	return captchaInstance.Verify(id, value, true)
}

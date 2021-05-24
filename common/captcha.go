package common

import (
	"errors"
	"github.com/dchest/captcha"
)

type Captcha struct {
	CaptchaId string `json:"captchaId"`
	ImageUrl  string `json:"imageUrl"`
}

var GetCaptcha = func(len int) (*Captcha, error) {
	d := struct {
		CaptchaId string
	}{
		captcha.NewLen(len),
	}
	var captcha Captcha
	if d.CaptchaId != "" {
		captcha.CaptchaId = d.CaptchaId
		captcha.ImageUrl = "/captcha/"+  d.CaptchaId
		return &captcha, nil
	}
	return &captcha, errors.New("验证码生成失败")
}

var VerifyCaptcha = func(captchaId, captchaValue string) error {
	if captchaId == "" || captchaValue == "" {
		return errors.New("验证码不存在")
	} else {
		if captcha.VerifyString(captchaId, captchaValue) {
			return nil
		} else {
			return errors.New("验证码错误")
		}
	}
}

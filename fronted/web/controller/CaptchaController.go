package controller

import (
	"bytes"
	"github.com/dchest/captcha"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"net/http"
	"time"
)

type CaptchaController struct {
	Ctx iris.Context
}

func (c *CaptchaController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/{captchaId}", "Get")
}

func (c *CaptchaController) Get()  {
	captchaId := c.Ctx.Params().Get("captchaId")

	if c.Ctx.URLParam("t") != "" {
		captcha.Reload(captchaId)
	}
	c.responseCaptchaImage(captchaId, 200, 50)
}

func (c *CaptchaController) responseCaptchaImage(id string, width, height int) error {
	w := c.Ctx.ResponseWriter()
	r := c.Ctx.Request()

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	var content bytes.Buffer
	ext := ".png"
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
		captcha.WriteImage(&content, id, width, height)
	case ".wav":
		w.Header().Set("Content-Type", "audio/x-wav")
		captcha.WriteAudio(&content, id, "zh")
	default:
		return captcha.ErrNotFound
	}

	download := false
	if download {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeContent(w, r, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
	return nil
}


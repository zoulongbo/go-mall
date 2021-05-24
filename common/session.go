package common

import (
	"github.com/kataras/iris/v12/sessions"
	"time"
)

var Session *sessions.Sessions

func SessionRegister() {
	Session = sessions.New(sessions.Config{
		Cookie:  "hello world",
		Expires: 60 * time.Minute,
	})
}

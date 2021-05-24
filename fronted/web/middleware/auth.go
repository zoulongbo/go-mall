package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/zoulongbo/go-mall/common"
)

func AuthLogin(ctx iris.Context) {

	uidSign := ctx.GetCookie("sign")
	uid := ctx.GetCookie("uid")
	if uidSign == "" || uid == ""{
		ctx.Application().Logger().Debug("必须先登录! cookie不足")
		ctx.Redirect("/user/login")
		return
	}
	uidByte, err := common.AesDePwdCode(uidSign)
	if err != nil {
		ctx.Application().Logger().Debug("必须先登录!sign解密失败")
		ctx.Redirect("/user/login")
		return
	}
	if string(uidByte) != uid {
		ctx.Application().Logger().Debug("必须先登录!uid匹配失败 | uid:" + uid + ", sign:" +uidSign)
		ctx.Redirect("/user/login")
		return
	}
	ctx.Application().Logger().Debug("已经登陆")
	ctx.Next()
}

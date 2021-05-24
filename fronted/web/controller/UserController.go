package controller

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/zoulongbo/go-mall/common"
	"github.com/zoulongbo/go-mall/models"
	"github.com/zoulongbo/go-mall/services"
	"github.com/zoulongbo/go-mall/tool"
	"strconv"
)

type UserController struct {
	Ctx iris.Context
}

func (u *UserController) GetRegister() mvc.View  {
	return mvc.View{
		Name:"user/register.html",
	}
}

func (u *UserController) PostRegister()  {

	var (
		nickname = u.Ctx.FormValue("nickname")
		username = u.Ctx.FormValue("username")
		password = u.Ctx.FormValue("password")
		captchaId = u.Ctx.FormValue("captchaId")
		userCaptcha = u.Ctx.FormValue("userCaptcha")
	)
	err := common.VerifyCaptcha(captchaId, userCaptcha)
	if err != nil {
		u.Ctx.Redirect("/user/register")
	}
	//TODO 缺表单验证 ozzo-validation
	user := &models.User{
		Nickname:     nickname,
		Username:     username,
		HashPassword: password,
	}
	_, err = services.NewUserService().AddUser(user)
	if err != nil {
		u.Ctx.Application().Logger().Debug(err)
		u.Ctx.Redirect("/user/error")
		return
	}
	u.Ctx.Redirect("/user/login")
	return
}

func (u *UserController) GetLogin() mvc.View  {
	capt, err := common.GetCaptcha(4)
	if err != nil {
		return mvc.View{
			Name:"shared/error.html",
			Data:iris.Map{
				"message":"验证码生成失败",
			},
		}
	}
	return mvc.View{
		Name:"user/login.html",
		Data: iris.Map{
			"captchaPath": capt.ImageUrl,
			"captchaId": capt.CaptchaId,
		},
	}
}

func (u *UserController) PostLogin() mvc.Response  {

	var (
		username = u.Ctx.FormValue("username")
		password = u.Ctx.FormValue("password")
		captchaId = u.Ctx.FormValue("captchaId")
		userCaptcha = u.Ctx.FormValue("userCaptcha")
	)
	err := common.VerifyCaptcha(captchaId, userCaptcha)
	if err != nil {
		return mvc.Response{
			Path:"/user/login",
		}
	}
	//TODO 缺表单验证 ozzo-validation
	user, isOk := services.NewUserService().CheckUser(username, password)
	if !isOk {
		return mvc.Response{
			Path:"/user/login",
		}
	}
	//写入用户ID到cookie中
	tool.GlobalCookie(u.Ctx, "uid", strconv.FormatInt(user.ID, 10))
	uidByte := []byte(strconv.FormatInt(user.ID, 10))
	uidString, err := common.AesEnPwdCode(uidByte)
	if err != nil {
		fmt.Println(err)
	}
	//写入用户浏览器
	tool.GlobalCookie(u.Ctx, "sign", uidString)
	return mvc.Response{
		Path:"/product/4",
	}
}
package routes

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/zoulongbo/go-mall/fronted/web/controller"
	"github.com/zoulongbo/go-mall/fronted/web/middleware"
)

func FrontedRegister(app *iris.Application, ctx context.Context) {
	user := app.Party("/user")
	u := mvc.New(user)
	u.Register(ctx)
	u.Handle(new(controller.UserController))


	product := app.Party("/product")
	p := mvc.New(product)
	product.Use(middleware.AuthLogin)
	p.Register(ctx)
	p.Handle(new(controller.ProductController))

	captcha := app.Party("/captcha")
	c := mvc.New(captcha)
	c.Register(ctx)
	c.Handle(new(controller.CaptchaController))

	order := app.Party("/order")
	o := mvc.New(order)
	order.Use(middleware.AuthLogin)
	o.Register(ctx)
	o.Handle(new(controller.OrderController))

}
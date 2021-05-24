package routes

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/zoulongbo/go-mall/backend/web/controller"
)

func BackendRegister(app *iris.Application, ctx context.Context) {
	example := app.Party("/example")
	b := mvc.New(example)
	b.Register(ctx)
	b.Handle(new(controller.ExampleController))

	product := app.Party("/product")
	p := mvc.New(product)
	p.Register(ctx)
	p.Handle(new(controller.ProductController))

	order := app.Party("/order")
	o := mvc.New(order)
	o.Register(ctx)
	o.Handle(new(controller.OrderController))
}
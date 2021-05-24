package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/zoulongbo/go-mall/common"
	"github.com/zoulongbo/go-mall/routes"
	"log"
)

func main()  {
	app := iris.New()
	//设置日志等级
	app.Logger().SetLevel("debug")
	//注册模版
	app.RegisterView(iris.HTML("./backend/web/views/template", ".html").Layout("shared/layout.html").Reload(true))
	//静态文件目录
	app.HandleDir("/assets", "./backend/web/views/assets")
	//出现异常跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("shared/error.html")
		err := ctx.View("shared/error.html")
		log.Println(err)
	})
	//ctx注册
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//mysql注册
	common.NewMysqlConn()
	//注册路由 控制器
	routes.BackendRegister(app, ctx)
	//启动服务
	app.Run(
		iris.Addr("localhost:8888"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
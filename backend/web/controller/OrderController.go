package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/zoulongbo/go-mall/common"
	"github.com/zoulongbo/go-mall/models"
	"github.com/zoulongbo/go-mall/services"
	"strconv"
)

type OrderController struct {
	Ctx iris.Context
}
func (o *OrderController) Get() mvc.View {
	order, err := services.NewOrderService().GetAllOrderInfo()
	if err != nil {
		o.Ctx.Application().Logger().Debug("订单信息查询失败 " + err.Error())
	}
	return mvc.View{
		Name: "order/view",
		Data: iris.Map{
			"order": order,
		},
	}
}


func (o *OrderController) GetAll() mvc.View {
	products, _ := services.NewProductService().GetAllProduct()
	return mvc.View{
		Name: "order/view",
		Data: iris.Map{
			"orders": products,
		},
	}
}

func (o *OrderController) PostUpdate() {
	product := &models.Product{}
	o.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "info"})
	if err := dec.Decode(o.Ctx.Request().Form, product); err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	err := services.NewProductService().UpdateProduct(product)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	o.Ctx.Redirect("/order/all")
}

func (o *OrderController) GetAdd() mvc.View {
	return mvc.View{
		Name: "order/add",
	}
}

func (o *OrderController) PostAdd() {
	product := &models.Product{}
	err := o.Ctx.Request().ParseForm()
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "info"})
	if err := dec.Decode(o.Ctx.Request().Form, product); err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	_, err = services.NewProductService().InsertProduct(product)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	o.Ctx.Redirect("/order/all")
}


func (o *OrderController) GetManager() mvc.View {
	idString := o.Ctx.URLParam("id")
	id ,err := strconv.ParseInt(idString, 10, 16)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	product ,err := services.NewProductService().GetProductById(id)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	return mvc.View{
		Name: "order/manager",
		Data:iris.Map{
			"order":product,
		},
	}
}

func (o *OrderController) GetDelete() {
	idString := o.Ctx.URLParam("id")
	id ,err := strconv.ParseInt(idString, 10, 16)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	res := services.NewProductService().DeleteProductById(id)
	if !res {
		o.Ctx.Application().Logger().Debug("删除失败 id:" + idString)
	}
	o.Ctx.Redirect("/order/all")
}

package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/zoulongbo/go-mall/common"
	"github.com/zoulongbo/go-mall/models"
	"github.com/zoulongbo/go-mall/services"
	"strconv"
)

type ProductController struct {
	Ctx iris.Context
}

func (p *ProductController) GetAll() mvc.View {
	products, _ := services.NewProductService().GetAllProduct()
	return mvc.View{
		Name: "product/view",
		Data: iris.Map{
			"products": products,
		},
	}
}

func (p *ProductController) PostUpdate() {
	product := &models.Product{}
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "info"})
	if err := dec.Decode(p.Ctx.Request().Form, product); err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	err := services.NewProductService().UpdateProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}

func (p *ProductController) GetAdd() mvc.View {
	return mvc.View{
		Name: "product/add",
	}
}

func (p *ProductController) PostAdd() {
	product := &models.Product{}
	err := p.Ctx.Request().ParseForm()
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "info"})
	if err := dec.Decode(p.Ctx.Request().Form, product); err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	_, err = services.NewProductService().InsertProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}


func (p *ProductController) GetManager() mvc.View {
	idString := p.Ctx.URLParam("id")
	id ,err := strconv.ParseInt(idString, 10, 16)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product ,err := services.NewProductService().GetProductById(id)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	return mvc.View{
		Name: "product/manager",
		Data:iris.Map{
			"product":product,
		},
	}
}

func (p *ProductController) GetDelete() {
	idString := p.Ctx.URLParam("id")
	id ,err := strconv.ParseInt(idString, 10, 16)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	res := services.NewProductService().DeleteProductById(id)
	if !res {
		p.Ctx.Application().Logger().Debug("删除失败 id:" + idString)
	}
	p.Ctx.Redirect("/product/all")
}

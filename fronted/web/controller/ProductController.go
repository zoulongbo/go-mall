package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/zoulongbo/go-mall/models"
	"github.com/zoulongbo/go-mall/services"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type ProductController struct {
	Ctx iris.Context
}

var (
	//生成的Html保存目录
	htmlOutPath = "./fronted/web/html/"
	//静态文件模版目录
	templatePath = "./fronted/web/views/template/"
)

func (p *ProductController) GetGenerateHtml() {
	productString := p.Ctx.URLParam("productId")
	productId,err:=strconv.Atoi(productString)
	if err !=nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	//1.获取模版
	contentTmp,err:=template.ParseFiles(filepath.Join(templatePath,"product.html"))
	if err !=nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//2.获取html生成路径
	fileName:=filepath.Join(htmlOutPath, productString + "-product.html")

	//3.获取模版渲染数据
	product,err:= services.NewProductService().GetProductById(int64(productId))
	if err !=nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//4.生成静态文件
	generateStaticHtml(p.Ctx,contentTmp,fileName,product)
}

//生成html静态文件
func generateStaticHtml(ctx iris.Context,template *template.Template,fileName string,product *models.Product)  {
	//1.判断静态文件是否存在
	if exist(fileName) {
		err:=os.Remove(fileName)
		if err !=nil {
			ctx.Application().Logger().Error(err)
		}
	}
	//2.生成静态文件
	file,err := os.OpenFile(fileName,os.O_CREATE|os.O_WRONLY,os.ModePerm)
	if err !=nil {
		ctx.Application().Logger().Error(err)
	}
	defer file.Close()
	template.Execute(file,&product)
}


//判断文件是否存在
func exist(fileName string) bool  {
	_,err:=os.Stat(fileName)
	return err==nil || os.IsExist(err)
}

func (p *ProductController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/{productId}", "GetDetail")
}


//商品详情
func (p *ProductController) GetDetail() interface{}  {
	id := p.Ctx.Params().Get("productId")

	fileName:=filepath.Join(htmlOutPath, id + "-product.html")
	if exist(fileName) {
		contents, err := ioutil.ReadFile(fileName)
		if err == nil {
			return mvc.Response{
				ContentType:"text/html",
				Content:contents,
			}
		}
	}
	productId, _ := strconv.ParseInt(id, 10, 16)
	product, err := services.NewProductService().GetProductById(productId)
	if err != nil {
		return mvc.View{
			Name:"shared/error",
			Data:iris.Map{
				"message":"商品数据跑路了" + err.Error(),
			},
		}
	}
	return mvc.View{
		Layout:"product/layout.html",
		Name:"product/view.html",
		Data:iris.Map{
			"product":product,
		},
	}
}


//下单
func (p *ProductController) GetOrder() mvc.View  {
	productIdString := p.Ctx.URLParam("productId")
	userIdString := p.Ctx.GetCookie("uid")
	userId, err := strconv.Atoi(userIdString)
	productId, err := strconv.Atoi(productIdString)
	if err != nil {
		p.Ctx.Application().Logger().Debug("商品id类型错误")
	}
	productService := services.NewProductService()
	product, err := productService.GetProductById(int64(productId))
	if err != nil {
		p.Ctx.Application().Logger().Debug("商品id类型错误")
	}
	//有库存再走下单
	if product.ProductNum > 0 {
		//缺锁 以防超卖
		product.ProductNum -= 1
		err := productService.UpdateProduct(product)
		if err != nil {
			p.Ctx.Application().Logger().Debug("商品库存变更失败")
		}
		//创建订单
		order := &models.Order{
			UserId:      int64(userId),
			ProductId:   int64(productId),
			OrderStatus: models.OrderStatusSuccess,
		}
		orderId, err := services.NewOrderService().InsertOrder(order)
		if err != nil {
			p.Ctx.Application().Logger().Debug("订单创建失败")
		}
		return mvc.View{
			Layout:"product/layout.html",
			Name:"product/result.html",
			Data:iris.Map{
				"message":"订单创建成功",
				"orderId":orderId,
			},
		}

	}
	return mvc.View{
		Layout:"product/layout.html",
		Name:"product/result.html",
		Data:iris.Map{
			"message":"商品库存不足",
			"orderId":"",
		},
	}
}


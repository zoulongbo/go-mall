package controller

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/zoulongbo/go-mall/models"
	"github.com/zoulongbo/go-mall/rabbitMQ"
	"strconv"
)

type OrderController struct {
	Ctx iris.Context
}

//下单
func (o *OrderController) Get() []byte  {
	productIdString := o.Ctx.URLParam("productId")
	userIdString := o.Ctx.GetCookie("uid")
	userId, err := strconv.ParseInt(userIdString, 10, 64)
	productId, err := strconv.ParseInt(productIdString, 10, 64)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	//入队数据处理
	message := models.NewMessage(userId, productId)
	byteMessage, err := json.Marshal(message)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	//初始化 mq
	rabbitMQ := rabbitMQ.NewRabbitMQSimple("orderAdd")
	err = rabbitMQ.PublishSimple(string(byteMessage))
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	return []byte("true")
	//
	/*return mvc.View{
		Layout:"product/layout.html",
		Name:"product/result.html",
		Data:iris.Map{
			"message":"订单创建成功",
			"orderId":orderId,
		},
	}*/

	/*return mvc.View{
		Layout:"product/layout.html",
		Name:"product/result.html",
		Data:iris.Map{
			"message":"商品库存不足",
			"orderId":"",
		},
	}*/
}


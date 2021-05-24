package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zoulongbo/go-mall/common"
	"github.com/zoulongbo/go-mall/models"
	"github.com/zoulongbo/go-mall/rabbitMQ"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

//存信息
type AccessControl struct {
	sourceArray map[int]interface{}
	//数据量大的时候 map可能报错 程序终止
	*sync.RWMutex
}
//数据存储全局变量
var accessControl = &AccessControl{
	sourceArray: make(map[int]interface{}),
}

//最好是内网ip
var hostArray = []string{"127.0.0.1","127.0.0.1","127.0.0.1","127.0.0.1"}
//本机ip
var localHost = ""
//数量控制接口服务ip 或负载内网ip
var GetOneIp = "127.0.0.1"

var GetOnePort = "8084"

var port = ":8083"

var hashConsistent = *common.NewConsistent()


var rabbitMqValidate *rabbitMQ.RabbitMQ

func Auth(rw http.ResponseWriter, req *http.Request) error  {
	log.Println("注册验证器")
	//添加基于cookie的权限验证
	err := CheckUserInfo(req)
	if err != nil {
		return err
	}
	return nil
}

func Check(rw http.ResponseWriter, req *http.Request)  {
	log.Println("执行验证器")
	queryForm, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil && len(queryForm["productId"]) <= 0 {
		rw.Write([]byte("false"))
		return
	}
	productIdString := queryForm["productId"][0]
	fmt.Println(productIdString)
	userCookie, err:= req.Cookie("uid")
	if err != nil {
		rw.Write([]byte("false"))
		return
	}
	//分布式权限验证
	right := accessControl.GetDistributedRight(req)
	if right == false {
		rw.Write([]byte("false"))
		return
	}
	//获取数量控制权限 防止超卖
	hostUrl := "http://" + GetOneIp + ":" + GetOnePort + "/getOne"
	responseValidate, validateBody, err := GetCurl(hostUrl, req)
	if err != nil {
		rw.Write([]byte("false"))
		return
	}
	//判断数量控制接口请求状态
	if responseValidate.StatusCode == http.StatusOK {
		if string(validateBody) == "true" {
			//整合下单
			productId, err := strconv.ParseInt(productIdString, 10, 64)
			if err != nil {
				rw.Write([]byte("false"))
				return
			}
			userId, err := strconv.ParseInt(userCookie.Value, 10, 64)
			if err != nil {
				rw.Write([]byte("false"))
				return
			}
			//创建消息体
			message := models.NewMessage(userId, productId)
			//类型转化
			byteMessage, err := json.Marshal(message)
			if err != nil {
				rw.Write([]byte("false"))
				return
			}
			err = rabbitMqValidate.PublishSimple(string(byteMessage))
			if err != nil {
				rw.Write([]byte("false"))
				return
			}
			rw.Write([]byte("true"))
			return
		}
	}
	rw.Write([]byte("false"))
	return
}


func CheckRight(rw http.ResponseWriter, req *http.Request)  {
	//分布式权限验证
	right := accessControl.GetDistributedRight(req)
	if right == false {
		rw.Write([]byte("false"))
		return
	}
	rw.Write([]byte("true"))
	return
}

func (m *AccessControl) GetNewRecord(uid int) interface{}  {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	data := m.sourceArray[uid]
	return data
}

func (m *AccessControl) SetNewRecord(uid int, record interface{}){
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	m.sourceArray[uid] = record
}

func (m *AccessControl) GetDistributedRight(req *http.Request) bool {
	uid, err := req.Cookie("uid")
	if err != nil {
		return false
	}
	//uid判断下 具体在哪个机器上
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}
	//是否为本机 如果是本机, 走本机代码 ，如果不是 代理
	if hostRequest == localHost {
		return m.GetDataFromMap(uid.Value)
	} else {
		return m.GetDataFromOtherMap(hostRequest, req)
	}
}

//获取其他机器的map 处理业务
func (m *AccessControl) GetDataFromOtherMap(host string, req *http.Request) bool  {
	hostUrl := "http://"+host+":"+port+"/checkRight"
	response, body, err := GetCurl(hostUrl, req)
	if err != nil {
		return false
	}
	if response.StatusCode == http.StatusOK {
		if string(body) == "true" {
			return true
		}
	}
	return false
}
//获取本机的map 处理业务
func (m *AccessControl) GetDataFromMap(uid string) (isOk bool)  {
	uidInt, err:= strconv.Atoi(uid)
	if err != nil {
		return false
	}
	data := m.GetNewRecord(uidInt)
	if data != nil {
		return true
	}
	return
}

func main()  {
	//一致性hash算法 匹配上层可能来自LB的请求
	for _, host := range hostArray {
		hashConsistent.Add(host)
	}
	localIp, err := common.GetEntranceIp()
	if err != nil {
		fmt.Println(err)
	}
	if localIp == localHost {
		fmt.Println(localHost)
	}
	rabbitMqValidate = rabbitMQ.NewRabbitMQSimple("orderAdd")
	defer rabbitMqValidate.Destroy()

	//1、过滤器
	filter := common.NewFilter()
	//注册拦截器
	filter.RegisterFilterUri("/check", Auth)
	filter.RegisterFilterUri("/checkRight", Auth)
	//2、启动服务
	http.HandleFunc("/check", filter.Handle(Check))
	http.HandleFunc("/checkRight", filter.Handle(CheckRight))
	//启动服务
	http.ListenAndServe(port, nil)
}

//自定义逻辑判断
func checkInfo(checkStr string, signStr string) bool {
	if checkStr == signStr {
		return true
	}
	return false
}


func CheckUserInfo(r *http.Request) error {
	//获取Uid，cookie
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		return errors.New("用户UID Cookie 获取失败！")
	}
	//获取用户加密串
	signCookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("用户加密串 Cookie 获取失败！")
	}

	//对信息进行解密
	signByte, err := common.AesDePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("加密串已被篡改！")
	}

	if checkInfo(uidCookie.Value, string(signByte)) {
		return nil
	}
	return errors.New("身份校验失败！")
	//return nil
}

//模拟请求
func GetCurl(hostUrl string,request *http.Request)(response *http.Response,body []byte,err error)  {
	//获取Uid
	uidPre,err := request.Cookie("uid")
	if err !=nil {
		return
	}
	//获取sign
	uidSign,err:=request.Cookie("sign")
	if err !=nil {
		return
	}

	//模拟接口访问，
	client :=&http.Client{}
	req,err:= http.NewRequest("GET",hostUrl,nil)
	if err !=nil {
		return
	}

	//手动指定，排查多余cookies
	cookieUid :=&http.Cookie{Name:"uid",Value:uidPre.Value,Path:"/"}
	cookieSign :=&http.Cookie{Name:"sign",Value:uidSign.Value,Path:"/"}
	//添加cookie到模拟的请求中
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)

	//获取返回结果
	response,err =client.Do(req)
	defer response.Body.Close()
	if err !=nil {
		return
	}
	body,err =ioutil.ReadAll(response.Body)
	return
}

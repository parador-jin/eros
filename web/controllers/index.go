package controllers

import (
	"Eros/comm"
	"Eros/models"
	"Eros/services"
	"fmt"
	"github.com/kataras/iris/v12"
)

type IndexController struct {
	Ctx            iris.Context
	ServiceUser    services.UserService
	ServiceBlackip services.BlackipService
	ServiceCode    services.CodeService
	ServiceUserday services.UserdayService
	ServiceResult  services.ResultService
	ServiceGift    services.GiftService
}

// http://localhost:8080/
func (c *IndexController) Get() string {
	c.Ctx.Header("Content-Type", "text/html")
	return "Welcome to Go抽奖系统, <a href='/public/index.html'>开始抽奖</a>"
}

// http://localhost:8080/gifts
func (c *IndexController) GetGifts() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	dataList := c.ServiceGift.GetAll(true)
	list := make([]models.LtGift, 0)
	for _, data := range dataList {
		if data.SysStatus == 0 {
			list = append(list, data)
		}
	}
	rs["gifts"] = list
	return rs
}

// http://localhost:8080/new/prize
func (c *IndexController) GetNewPrize() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	gifts := c.ServiceGift.GetAll(true)
	giftIds := []int{}
	for _, data := range gifts {
		// 虚拟券或者实物奖才需要放到外部榜单中展示
		if data.Gtype > 1 {
			giftIds = append(giftIds, data.Id)
		}
	}
	list := c.ServiceResult.GetNewPrize(50, giftIds)
	rs["prize_list"] = list
	return rs
}

// http://localhost:8080/login
func (c *IndexController) GetLogin() {
	uid := comm.Random(100000)
	loginUser := models.ObjLoginuser{
		Uid:      uid,
		Username: fmt.Sprintf("admin-%d", uid),
		Now:      comm.NowUnix(),
		Ip:       comm.ClientIP(c.Ctx.Request()),
		Sign:     "",
	}
	comm.SetLoginUser(c.Ctx.ResponseWriter(), &loginUser)
	comm.Redirect(c.Ctx.ResponseWriter(), "/public/index.html?from=login")
}

// http://localhost:8080/logout
func (c *IndexController) GetLogout() {
	comm.SetLoginUser(c.Ctx.ResponseWriter(), nil)
	comm.Redirect(c.Ctx.ResponseWriter(), "public/index.html?from=logout")
}

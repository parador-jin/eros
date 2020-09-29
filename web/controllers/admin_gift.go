package controllers

import (
	"Eros/comm"
	"Eros/models"
	"Eros/services"
	"Eros/web/utils"
	"Eros/web/viewmodels"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"log"
	"time"
)

type AdminGiftController struct {
	Ctx         iris.Context
	ServiceGift services.GiftService
}

// http://localhost:8080/admin/gift
func (c *AdminGiftController) Get() mvc.Result {
	dataList := c.ServiceGift.GetAll(true)
	for i, giftInfo := range dataList {
		// 奖品发放的计划数据
		prizeData := make([][2]int, 0)
		err := json.Unmarshal([]byte(giftInfo.PrizeData), &prizeData)
		if err != nil || len(prizeData) < 1 {
			dataList[i].PrizeData = "[]"
		} else {
			newPd := make([]string, len(prizeData))
			for index, pd := range prizeData {
				ct := comm.FormatFromUnixTime(int64(pd[0]))
				newPd[index] = fmt.Sprintf("【%s】: %d", ct, pd[1])
			}
			str, err := json.Marshal(newPd)
			if err == nil && len(str) > 0 {
				dataList[i].PrizeData = string(str)
			} else {
				dataList[i].PrizeData = "[]"
			}
		}
		num := utils.GetGiftPoolNum(giftInfo.Id)
		dataList[i].Title = fmt.Sprintf("【%d】%s", num, dataList[i].Title)
	}
	return mvc.View{
		Name: "admin/gift.html",
		Data: iris.Map{
			"Title":    "管理后台",
			"Channel":  "gift",
			"Datalist": dataList,
			"Total":    len(dataList),
		},
		Layout: "admin/layout.html",
	}
}

// http://localhost:8080/admin/gift/edit
func (c *AdminGiftController) GetEdit() mvc.Result {
	id := c.Ctx.URLParamIntDefault("id", 0)
	giftInfo := viewmodels.ViewGift{}
	if id > 0 {
		data := c.ServiceGift.Get(id, true)
		giftInfo.Id = data.Id
		giftInfo.Displayorder = data.Displayorder
		giftInfo.Gdata = data.Gdata
		giftInfo.Gtype = data.Gtype
		giftInfo.Img = data.Img
		giftInfo.PrizeCode = data.PrizeCode
		giftInfo.PrizeNum = data.PrizeNum
		giftInfo.PrizeTime = data.PrizeTime
		giftInfo.TimeBegin = comm.FormatFromUnixTime(int64(data.TimeBegin))
		giftInfo.TimeEnd = comm.FormatFromUnixTime(int64(data.TimeEnd))
		giftInfo.Title = data.Title
	}
	return mvc.View{
		Name: "admin/giftEdit.html",
		Data: iris.Map{
			"Title":   "管理后台",
			"Channel": "gift",
			"info":    giftInfo,
		},
		Layout: "admin/layout.html",
	}
}

// http://localhost:8080/admin/gift/save
func (c *AdminGiftController) PostSave() mvc.Result {
	data := viewmodels.ViewGift{}
	err := c.Ctx.ReadForm(&data)
	if err != nil {
		log.Println("admin_gift.PostSave ReadForm error is ", err)
		return mvc.Response{
			Text: fmt.Sprintf("ReadForm转换异常，error is ", err),
		}
	}
	giftInfo := models.LtGift{}
	giftInfo.Id = data.Id
	giftInfo.Title = data.Title
	giftInfo.PrizeTime = data.PrizeTime
	giftInfo.PrizeNum = data.PrizeNum
	giftInfo.PrizeCode = data.PrizeCode
	giftInfo.Img = data.Img
	giftInfo.Gtype = data.Gtype
	giftInfo.Gdata = data.Gdata
	giftInfo.Displayorder = data.Displayorder
	t1, err1 := comm.ParseTime(data.TimeBegin)
	t2, err2 := comm.ParseTime(data.TimeEnd)
	if err1 != nil || err2 != nil {
		return mvc.Response{
			Text: fmt.Sprintf("开始时间、结束时间的格式不正确，err1 is %s, err2 is %s", err1, err2),
		}
	}
	giftInfo.TimeBegin = int(t1.Unix())
	giftInfo.TimeEnd = int(t2.Unix())
	if giftInfo.Id > 0 {
		// 数据更新
		dataInfo := c.ServiceGift.Get(giftInfo.Id, true)
		if dataInfo != nil && dataInfo.Id > 0 {
			if dataInfo.PrizeNum != giftInfo.PrizeNum {
				// 奖品数量发生变化
				giftInfo.LeftNum = dataInfo.LeftNum - dataInfo.PrizeNum - giftInfo.PrizeNum
				if giftInfo.LeftNum < 0 || giftInfo.PrizeNum <= 0 {
					giftInfo.LeftNum = 0
				}
				// 奖品总数发生了变化
				utils.ResetGiftPrizeData(&giftInfo, c.ServiceGift)
			}
			if dataInfo.PrizeTime != giftInfo.PrizeTime {
				// 发奖周期发生变化
				utils.ResetGiftPrizeData(&giftInfo, c.ServiceGift)
			}
			giftInfo.SysUpdated = int(time.Now().Unix())
			c.ServiceGift.Update(&giftInfo, []string{""})
		} else {
			giftInfo.Id = 0
		}
	}
	if giftInfo.Id == 0 {
		giftInfo.LeftNum = giftInfo.PrizeNum
		giftInfo.SysIp = comm.ClientIP(c.Ctx.Request())
		giftInfo.SysCreated = int(time.Now().Unix())
		c.ServiceGift.Create(&giftInfo)
		// 新的奖品，更新奖品的发奖计划
		utils.ResetGiftPrizeData(&giftInfo, c.ServiceGift)
	}
	return mvc.Response{
		Path: "/admin/gift",
	}
}

// http://localhost:8080/admin/gift/delete
func (c *AdminGiftController) GetDelete() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		c.ServiceGift.Delete(id)
	}
	return mvc.Response{
		Path: "/admin/gift",
	}
}

// http://localhost:8080/admin/gift/reset
func (c *AdminGiftController) GetReset() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		c.ServiceGift.Update(&models.LtGift{Id: id, SysStatus: 0}, []string{"sys_status"})
	}
	return mvc.Response{
		Path: "/admin/gift",
	}
}

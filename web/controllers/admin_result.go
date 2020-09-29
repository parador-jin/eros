package controllers

import (
	"Eros/models"
	"Eros/services"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type AdminResultController struct {
	Ctx           iris.Context
	ServiceResult services.ResultService
}

// http://localhost:8080/admin/result
func (c *AdminResultController) Get() mvc.Result {
	giftId := c.Ctx.URLParamIntDefault("gift_id", 0)
	uid := c.Ctx.URLParamIntDefault("uid", 0)
	page := c.Ctx.URLParamIntDefault("page", 1)
	size := 100
	pagePrev := ""
	pageNext := ""
	// 数据列表
	var dataList []models.LtResult
	if giftId > 0 {
		dataList = c.ServiceResult.SearchByGift(giftId, page, size)
	} else if uid > 0 {
		dataList = c.ServiceResult.SearchByUid(uid, page, size)
	} else {
		dataList = c.ServiceResult.GetAll(page, size)
	}
	total := (page-1)*size + len(dataList)
	// 数据总数
	if len(dataList) >= size {
		if giftId > 0 {
			total = int(c.ServiceResult.CountByGift(giftId))
		} else if uid > 0 {
			total = int(c.ServiceResult.CountByUid(uid))
		} else {
			total = int(c.ServiceResult.CountAll())
		}
		pageNext = fmt.Sprintf("%d", page+1)
	}
	if page > 1 {
		pagePrev = fmt.Sprintf("%d", page-1)
	}
	return mvc.View{
		Name: "admin/result.html",
		Data: iris.Map{
			"Title":    "管理后台",
			"Channel":  "result",
			"GiftId":   giftId,
			"Uid":      uid,
			"Datalist": dataList,
			"Total":    total,
			"PagePrev": pagePrev,
			"PageNext": pageNext,
		},
		Layout: "admin/layout.html",
	}
}

// http://localhost:8080/admin/result/delete
func (c *AdminResultController) GetDelete() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		c.ServiceResult.Delete(id)
	}
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/result"
	}
	return mvc.Response{
		Path: refer,
	}
}

// http://localhsot:8080/admin/result/cheat
func (c *AdminResultController) GetCheat() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		c.ServiceResult.Update(&models.LtResult{Id: id, SysStatus: 2}, []string{"sys_status"})
	}
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/result"
	}
	return mvc.Response{
		Path: refer,
	}
}

// http://localhsot:8080/admin/result/reset
func (c *AdminResultController) GetReset() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		c.ServiceResult.Update(&models.LtResult{Id: id, SysStatus: 0}, []string{"sys_status"})
	}
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/result"
	}
	return mvc.Response{
		Path: refer,
	}
}

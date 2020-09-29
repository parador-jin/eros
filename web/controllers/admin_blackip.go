package controllers

import (
	"Eros/comm"
	"Eros/models"
	"Eros/services"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type AdminBlackIPController struct {
	Ctx            iris.Context
	ServiceBlackIp services.BlackipService
}

// http://localhost:8080/admin/blackip
func (c *AdminBlackIPController) Get() mvc.Result {
	page := c.Ctx.URLParamIntDefault("page", 1)
	size := 100
	pagePrev := ""
	pageNext := ""
	// 数据列表
	dataList := c.ServiceBlackIp.GetAll(page, size)
	total := (page-1)*size + len(dataList)
	if len(dataList) >= size {
		total = int(c.ServiceBlackIp.CountAll())
		pageNext = fmt.Sprintf("%d", page+1)
	}
	if page > 1 {
		pagePrev = fmt.Sprintf("%d", page-1)
	}
	return mvc.View{
		Name: "admin/blackip.html",
		Data: iris.Map{
			"Title":    "管理后台",
			"Channel":  "blackip",
			"Datalist": dataList,
			"Total":    total,
			"Now":      comm.NowUnix(),
			"PagePrev": pagePrev,
			"PageNext": pageNext,
		},
		Layout: "admin/layout.html",
	}
}

// http://localhost:8080/admin/blackip/black?id=1&time=0
func (c *AdminBlackIPController) GetBlack() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	t := c.Ctx.URLParamIntDefault("time", 0)
	if err == nil {
		if t > 0 {
			t = t*86400 + comm.NowUnix()
		}
		c.ServiceBlackIp.Update(&models.LtBlackip{Id: id, Blacktime: t, SysUpdated: comm.NowUnix()}, []string{"blacktime"})
	}
	return mvc.Response{
		Path: "/admin/blackip",
	}
}

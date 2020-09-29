package controllers

import (
	"Eros/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type AdminController struct {
	Ctx            iris.Context
	ServiceUser    services.UserService
	ServiceUserday services.UserdayService
	ServiceGift    services.GiftService
	ServiceCode    services.CodeService
	ServiceResult  services.ResultService
}

// http://localhost:8080/admin
func (c *AdminController) Get() mvc.Result {
	return mvc.View{
		Name: "admin/index.html",
		Data: iris.Map{
			"Title":   "管理后台",
			"Channel": "",
		},
		Layout: "admin/layout.html",
	}
}

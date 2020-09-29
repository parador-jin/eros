package routes

import (
	"Eros/bootstrap"
	"Eros/services"
	"Eros/web/controllers"
	"Eros/web/middleware"
	"github.com/kataras/iris/v12/mvc"
)

func Configure(b *bootstrap.Bootstrapper) {
	userService := services.NewUserService()
	giftService := services.NewGiftService()
	userdayService := services.NewUserdayService()
	resultService := services.NewResultService()
	blackipService := services.NewBlackipService()
	codeService := services.NewCodeService()
	// 首页
	index := mvc.New(b.Party("/"))
	index.Register(userdayService, userService, giftService, resultService, blackipService, codeService)
	index.Handle(new(controllers.IndexController))
	// 后台管理首页
	admin := mvc.New(b.Party("/admin"))
	admin.Router.Use(middleware.BasicAuth)
	admin.Register(userdayService, userService, giftService, resultService, blackipService, codeService)
	admin.Handle(new(controllers.AdminController))
	// 商品页
	adminGift := admin.Party("/gift")
	adminGift.Register(giftService)
	adminGift.Handle(new(controllers.AdminGiftController))
	// 优惠券页
	adminCode := admin.Party("/code")
	adminCode.Register(codeService, giftService)
	adminCode.Handle(new(controllers.AdminCodeController))
	// 中奖记录页
	adminResult := admin.Party("/result")
	adminResult.Register(resultService)
	adminResult.Handle(new(controllers.AdminResultController))
	// 用户管理页
	adminUser := admin.Party("/user")
	adminUser.Register(userService)
	adminUser.Handle(new(controllers.AdminUserController))
	// 黑名单管理页
	adminBlackIP := admin.Party("/blackip")
	adminBlackIP.Register(blackipService)
	adminBlackIP.Handle(new(controllers.AdminBlackIPController))
}

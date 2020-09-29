package main

import (
	"Eros/bootstrap"
	"Eros/conf"
	"Eros/web/middleware/identity"
	"Eros/web/routes"
	"fmt"
	"github.com/kataras/iris/v12"
)

func newApp() *bootstrap.Bootstrapper {
	// 初始化应用
	app := bootstrap.New("Go抽奖系统", "Eros")
	app.Bootstrap()
	app.Configure(identity.Configure, routes.Configure)

	return app

}

func main() {
	app := newApp()
	app.Listen(fmt.Sprintf(":%d", conf.Port),
		iris.WithoutBanner,
		iris.WithoutServerError(iris.ErrServerClosed),
	)
}

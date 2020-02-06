package main

import (
	"github.com/Ggkd/conf"
	"github.com/Ggkd/secProxy/router"
	_ "github.com/Ggkd/secProxy/service"
	"github.com/gin-gonic/gin"
)

func main()  {
	app := gin.New()
	app.Use(gin.Logger())
	app.Use(gin.Recovery())
	// 加载路由
	router.Router(app)
	addr := conf.Config.Host.Ip + ":" + conf.Config.Host.Port
	app.Run(addr)
}
package router

import (
	"github.com/Ggkd/secProxy/controller"
	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine)  {
	router.GET("/secInfo", controller.SecInfo)
	router.GET("/secKill", controller.SecKill)
}
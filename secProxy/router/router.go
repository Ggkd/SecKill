package router

import (
	"github.com/Ggkd/secProxy/controller"
	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine)  {
	router.Any("/secInfo", controller.SecInfo)
	router.Any("/secKill", controller.SecKill)
}
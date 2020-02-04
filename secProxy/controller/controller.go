package controller

import (
	"github.com/Ggkd/conf"
	"github.com/Ggkd/secProxy/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)


// 获取秒杀商品信息
func SecInfo(c *gin.Context) {
	resp := new(conf.Result)
	resp.Code = 1
	resp.Msg = "success"
	defer c.JSON(http.StatusOK, resp)
	productId := c.Query("product_id")
	if productId == "" {
		resp.Code = 0
		resp.Msg = "product id is null"
		return
	}
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		resp.Code = 0
		resp.Msg = "product id isn't valid"
		return
	}
	data, err := service.SecInfoService(productIdInt)
	if err != nil {
		resp.Code = 0
		resp.Msg = err.Error()
		return
	}
	resp.Data = data
}

func SecKill(c *gin.Context) {

}
package controller

import (
	"github.com/Ggkd/conf"
	"github.com/Ggkd/secProxy/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)


// 获取秒杀商品信息
func SecInfo(c *gin.Context) {
	resp := new(conf.Result)
	resp.Code = 1
	resp.Msg = "success"
	defer c.JSON(http.StatusOK, resp)
	productId := c.Query("product_id")
	if productId == "" {
		data := service.SecInfoList()
		resp.Data = data
	} else {
		productIdInt, err := strconv.Atoi(productId)
		if err != nil {
			resp.Code = 0
			resp.Msg = "product id isn't valid"
			return
		}
		data := service.SecInfoById(productIdInt)
		resp.Data = data
	}
}


func SecKill(c *gin.Context) {
	resp := new(conf.Result)
	resp.Code = 1
	resp.Msg = "success"
	defer c.JSON(http.StatusOK, resp)
	productId := c.Query("product_id")
	secTime := c.Query("sec_time")
	source := c.Query("source")
	nance := c.Query("nance")
	authCode := c.Query("auth_code")
	userId, _ := c.Cookie("user_id")
	userCookieAuth, _ := c.Cookie("user_cookie_auth")
	userAddr := strings.Split(c.Request.RemoteAddr, ":")[0]
	userRefer := c.Request.Referer()
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		resp.Code = 0
		resp.Msg = "product id isn't valid"
		return
	}
	secTimeInt, err := strconv.Atoi(secTime)
	if err != nil {
		resp.Code = 0
		resp.Msg = "secTime isn't valid"
		return
	}
	req := &conf.ReqSecKill{}
	req.ProductId = productIdInt
	req.Source = source
	req.Nance = nance
	req.SecTime = secTimeInt
	req.UserId = userId
	req.AuthCode = authCode
	req.UserAuth = userCookieAuth
	req.UserAddr = userAddr
	req.UserRefer = userRefer
	accessTime := time.Now()
	req.AccessTime = accessTime
	data, code, msg := service.SecKill(req)
	resp.Code = code
	if msg != nil {
		resp.Msg = msg.Error()
	}
	resp.Data = data
}
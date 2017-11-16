package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	pay "github.com/yuntifree/wifi-server/proto/pay"
)

const (
	payName = "go.micro.srv.pay"
)

func payHandler(c *gin.Context) {
	action := c.Param("action")
	log.Printf("action:%s", action)
	switch action {
	case "wx_pay":
		wxPay(c)
	default:
		c.JSON(http.StatusOK, gin.H{"errno": errAction, "desc": "unknown action"})
	}
	return
}

type payReq struct {
	Id    int64 `json:"id"`
	Price int64 `json:"price"`
}

func wxPay(c *gin.Context) {
	var p payReq
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	wid := getCookieInt(c, "wid")
	uid := getCookieInt(c, "u")
	var req pay.WxPayRequest
	req.Wid = int64(wid)
	req.Uid = int64(uid)
	req.Item = p.Id
	req.Price = p.Price
	ip := c.ClientIP()
	callback := strings.Replace(c.Request.RequestURI, "wx_pay", "wx_pay_callback", -1)
	req.Clientip = ip
	req.Callback = callback
	cl := pay.NewPayClient(payName, client.DefaultClient)
	rsp, err := cl.WxPay(context.Background(), &req)
	if err != nil {
		log.Printf("getBusinessInfo Info failed:%d %v", wid, err)
		c.JSON(http.StatusOK, gin.H{"errno": errParam, "desc": "illegal param"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": 0, "paySign": rsp.Sign,
		"package": rsp.Pack, "timeStamp": rsp.Ts,
		"nonceStr": rsp.Nonce, "signType": rsp.Signtype})
}

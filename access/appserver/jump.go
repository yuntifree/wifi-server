package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/yuntifree/wifi-server/accounts"
)

const (
	wxHost  = "http://wx.yunxingzh.com/GetWeixinCode/get-weixin-code.html"
	logHost = "http://wxdev.yunxingzh.com/weixin/login"
)

func jumpHandler(c *gin.Context) {
	url := fmt.Sprintf("%s?appid=%s&scope=snsapi_base&state=list&redirect_uri=%s",
		wxHost, accounts.DgWxAppid, url.QueryEscape(logHost))
	log.Printf("url:%s", url)
	c.Header("Cache-Control", "no-store")
	c.Redirect(http.StatusMovedPermanently, url)
}

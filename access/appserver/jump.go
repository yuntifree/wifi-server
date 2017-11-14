package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuntifree/components/weixin"
	"github.com/yuntifree/wifi-server/accounts"
)

const (
	wxHost = "http://wx.yunxingzh.com/"
)

func jumpHandler(c *gin.Context) {
	filename := c.Param("filename")
	log.Printf("filename:%s", filename)
	echostr := wxHost + filename
	redirect := wxHost + "weixin/login" + "?echostr=" + echostr
	wx := weixin.WxInfo{Appid: accounts.DgWxAppid,
		Appkey: accounts.DgWxAppkey}
	dst := wx.GenRedirect(redirect)
	c.Redirect(http.StatusMovedPermanently, dst)
}

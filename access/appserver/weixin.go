package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	wx "github.com/yuntifree/wifi-server/proto/wx"
)

const (
	wxName  = "go.micro.srv.wx"
	homeDst = "http://wxdev.yunxingzh.com/static/bussinesstest/#/home"
)

func loginHandler(c *gin.Context) {
	code := c.Query("code")
	var req wx.LoginRequest
	req.Code = code
	cl := wx.NewWxClient(wxName, client.DefaultClient)
	rsp, err := cl.Login(context.Background(), &req)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	{
		var co http.Cookie
		co.Name = "u"
		co.Value = fmt.Sprintf("%d", rsp.Uid)
		co.Path = "/"
		co.Domain = domain
		http.SetCookie(c.Writer, &co)
	}
	{
		var co http.Cookie
		co.Name = "s"
		co.Value = rsp.Token
		co.Path = "/"
		co.Domain = domain
		http.SetCookie(c.Writer, &co)
	}
	c.Header("Cache-Control", "no-store")
	c.Redirect(http.StatusMovedPermanently, homeDst)
}

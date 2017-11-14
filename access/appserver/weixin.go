package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	wx "github.com/yuntifree/wifi-server/proto/wx"
)

const (
	wxName = "go.micro.srv.wx"
)

func loginHandler(c *gin.Context) {
	code := c.Query("code")
	log.Printf("code:%s", code)
	var req wx.LoginRequest
	req.Code = code
	cl := wx.NewWxClient(wxName, client.DefaultClient)
	rsp, err := cl.Login(context.Background(), &req)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	var co http.Cookie
	co.Name = "u"
	co.Value = fmt.Sprintf("%d", rsp.Uid)
	http.SetCookie(c.Writer, &co)
	echostr := c.Query("echostr")
	c.Redirect(http.StatusMovedPermanently, echostr)
}

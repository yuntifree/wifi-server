package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	hall "github.com/yuntifree/wifi-server/proto/hall"
	phone "github.com/yuntifree/wifi-server/proto/phone"
)

const (
	hallName = "go.micro.srv.hall"
	domain   = "asr.yunxingzh.com"
)

func hallHandler(c *gin.Context) {
	action := c.Param("action")
	log.Printf("action:%s", action)
	switch action {
	case "login":
		loginHall(c)
	default:
		c.JSON(http.StatusOK, gin.H{"errno": errAction, "desc": "unknown action"})
	}
	return
}

type loginHallRequest struct {
	Phone string `json:"phone"`
	Code  int64  `json:"code"`
	Park  int64  `json:"park"`
}

func loginHall(c *gin.Context) {
	var req loginHallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": errParam, "desc": err.Error()})
		return
	}
	log.Printf("loginHall req:%+v", req)
	var check phone.CheckRequest
	check.Phone = req.Phone
	check.Code = req.Code
	cl := phone.NewPhoneClient(phoneName, client.DefaultClient)
	_, err := cl.CheckCode(context.Background(), &check)
	if err != nil {
		log.Printf("loginHall CheckCode failed:%v", err)
		c.JSON(http.StatusOK, gin.H{"errno": errCheckCode, "desc": "验证码错误"})
		return
	}
	var login hall.LoginRequest
	login.Phone = req.Phone
	login.Park = req.Park
	hl := hall.NewHallClient(hallName, client.DefaultClient)
	res, err := hl.Login(context.Background(), &login)
	if err != nil {
		log.Printf("loginHall Login failed:%v", err)
		c.JSON(http.StatusOK, gin.H{"errno": errCheckCode, "desc": "验证码错误"})
		return
	}
	{
		var co http.Cookie
		co.Name = "wid"
		co.Path = "/"
		co.Domain = domain
		co.Value = fmt.Sprintf("%d", res.Wid)
		http.SetCookie(c.Writer, &co)
	}
	{
		var co http.Cookie
		co.Name = "phone"
		co.Path = "/"
		co.Domain = domain
		co.Value = req.Phone
		http.SetCookie(c.Writer, &co)
	}
	c.JSON(http.StatusOK, gin.H{"errno": errOK})
	return
}

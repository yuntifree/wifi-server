package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	verify "github.com/yuntifree/wifi-server/proto/verify"
)

const (
	verifyName = "go.micro.srv.verify"
)

type checkRequest struct {
	Data checkData `json:"data"`
}

type checkData struct {
	WlanUsermac string `json:"wlanusermac"`
	WlanAcname  string `json:"wlanacname"`
	WlanApmac   string `json:"wlanapmac"`
}

func checkLoginHandler(c *gin.Context) {
	var p checkRequest
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	var req verify.CheckRequest
	req.Usermac = p.Data.WlanUsermac
	req.Acname = p.Data.WlanAcname
	req.Apmac = p.Data.WlanApmac
	cl := verify.NewVerifyClient(verifyName, client.DefaultClient)
	rsp, err := cl.CheckLogin(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": 0, "data": map[string]interface{}{
		"autologin": rsp.Autologin, "img": rsp.Img, "wxappid": rsp.Wxappid,
		"wxsecret": rsp.Wxsecret, "wxshopid": rsp.Wxshopid, "wxauthurl": rsp.Wxauthurl,
		"taobao": rsp.Taobao, "logintype": rsp.Logintype, "dst": rsp.Dst,
		"cover": rsp.Cover}})
}

type codeRequest struct {
	Data codeData `json:"data"`
}

type codeData struct {
	Phone      string `json:"phone"`
	WlanAcname string `json:"wlanacname"`
	WlanApmac  string `json:"wlanapmac"`
}

func getCodeHandler(c *gin.Context) {
	var p codeRequest
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	var req verify.CodeRequest
	req.Phone = p.Data.Phone
	req.Acname = p.Data.WlanAcname
	req.Apmac = p.Data.WlanApmac
	cl := verify.NewVerifyClient(verifyName, client.DefaultClient)
	_, err := cl.GetCheckCode(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": 0})
}

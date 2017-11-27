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
	Data verify.CheckRequest `json:"data"`
}

func checkLoginHandler(c *gin.Context) {
	var p checkRequest
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	cl := verify.NewVerifyClient(verifyName, client.DefaultClient)
	rsp, err := cl.CheckLogin(context.Background(), &p.Data)
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
	Data verify.CodeRequest `json:"data"`
}

func getCodeHandler(c *gin.Context) {
	var p codeRequest
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	cl := verify.NewVerifyClient(verifyName, client.DefaultClient)
	_, err := cl.GetCheckCode(context.Background(), &p.Data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": 0})
}

type portalRequest struct {
	Data verify.PortalLoginRequest `json:"data"`
}

func portalLoginHandler(c *gin.Context) {
	var p portalRequest
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	cl := verify.NewVerifyClient(verifyName, client.DefaultClient)
	rsp, err := cl.PortalLogin(context.Background(), &p.Data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": 0, "data": map[string]interface{}{
		"uid": rsp.Uid, "token": rsp.Token, "portaldir": rsp.Portaldir,
		"portaltype": rsp.Portaltype, "adtype": rsp.Adtype,
		"cover": rsp.Cover, "dst": rsp.Dst}})
}

type oneClickRequest struct {
	Data verify.OneClickRequest `json:"data"`
}

func oneClickLoginHandler(c *gin.Context) {
	var p oneClickRequest
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	cl := verify.NewVerifyClient(verifyName, client.DefaultClient)
	rsp, err := cl.OneClickLogin(context.Background(), &p.Data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": 0, "data": map[string]interface{}{
		"uid": rsp.Uid, "token": rsp.Token, "portaldir": rsp.Portaldir,
		"portaltype": rsp.Portaltype, "adtype": rsp.Adtype,
		"cover": rsp.Cover, "dst": rsp.Dst}})
}

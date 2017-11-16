package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	business "github.com/yuntifree/wifi-server/proto/business"
)

const (
	businessName = "go.micro.srv.business"
)

func businessHandler(c *gin.Context) {
	action := c.Param("action")
	log.Printf("action:%s", action)
	switch action {
	case "info":
		getBusinessInfo(c)
	default:
		c.JSON(http.StatusOK, gin.H{"errno": errAction, "desc": "unknown action"})
	}
	return
}

func getBusinessInfo(c *gin.Context) {
	id, err := c.Cookie("wid")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": errParam, "desc": "illegal param"})
		return
	}
	wid, err := strconv.Atoi(id)
	var req business.Request
	req.Wid = int64(wid)
	cl := business.NewBusinessClient(businessName, client.DefaultClient)
	rsp, err := cl.Info(context.Background(), &req)
	if err != nil {
		log.Printf("getBusinessInfo Info failed:%d %v", wid, err)
		c.JSON(http.StatusOK, gin.H{"errno": errParam, "desc": "illegal param"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": 0, "payed": rsp.Payed,
		"expire": rsp.Expire, "items": rsp.Items})
}

package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	phone "github.com/yuntifree/wifi-server/proto/phone"
)

const (
	phoneName = "go.micro.srv.phone"
)

func phoneHandler(c *gin.Context) {
	action := c.Param("action")
	log.Printf("action:%s", action)
	switch action {
	case "getcode":
		getPhoneCode(c)
	default:
		c.JSON(http.StatusOK, gin.H{"errno": 101, "desc": "unknown action"})
	}
	return
}

func getPhoneCode(c *gin.Context) {
	var req phone.GetRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
	}
	if req.Phone == "" {
		c.JSON(http.StatusOK, gin.H{"errno": 2, "desc": "illegal param"})
	}
	log.Printf("getPhoneCode req:%+v", req)
	cl := phone.NewPhoneClient(phoneName, client.DefaultClient)
	_, err := cl.GetCode(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"errno": 0})
	}
}

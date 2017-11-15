package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	park "github.com/yuntifree/wifi-server/proto/park"
)

const (
	parkName = "go.micro.srv.park"
)

func parkHandler(c *gin.Context) {
	action := c.Param("action")
	log.Printf("action:%s", action)
	switch action {
	case "get":
		getPark(c)
	default:
		c.JSON(http.StatusOK, gin.H{"errno": 101, "desc": "unknown action"})
	}
	return
}

func getPark(c *gin.Context) {
	var req park.GetRequest
	cl := park.NewParkClient(parkName, client.DefaultClient)
	rsp, err := cl.Get(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
	} else {
		log.Printf("rsp:%+v", rsp)
		if len(rsp.Infos) > 0 {
			c.JSON(http.StatusOK, gin.H{"errno": 0, "infos": rsp.Infos})
		} else {
			c.JSON(http.StatusOK, gin.H{"errno": 102,
				"desc": "empty response"})
		}
	}
}

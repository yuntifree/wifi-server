package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	trial "github.com/yuntifree/wifi-server/proto/trial"
)

const (
	trialName = "go.micro.srv.trial"
)

func trialHandler(c *gin.Context) {
	action := c.Param("action")
	log.Printf("action:%s", action)
	switch action {
	case "info":
		getTrialInfo(c)
	case "apply":
		applyTrial(c)
	default:
		c.JSON(http.StatusOK, gin.H{"errno": errAction, "desc": "unknown action"})
	}
	return
}

func getTrialInfo(c *gin.Context) {
	wid := getCookieInt(c, "wid")
	var req trial.Request
	req.Wid = int64(wid)
	cl := trial.NewTrialClient(trialName, client.DefaultClient)
	rsp, err := cl.Info(context.Background(), &req)
	if err != nil {
		log.Printf("getTrialInfo Info failed:%d %v", wid, err)
		c.JSON(http.StatusOK, gin.H{"errno": errParam, "desc": "illegal param"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": 0, "used": rsp.Used})
}

func applyTrial(c *gin.Context) {
	wid := getCookieInt(c, "wid")
	var req trial.Request
	req.Wid = int64(wid)
	cl := trial.NewTrialClient(trialName, client.DefaultClient)
	_, err := cl.Apply(context.Background(), &req)
	if err != nil {
		log.Printf("applyTrial Apply failed:%d %v", wid, err)
		c.JSON(http.StatusOK, gin.H{"errno": errParam, "desc": "illegal param"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": 0})
}

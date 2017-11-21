package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	feedback "github.com/yuntifree/wifi-server/proto/feedback"
)

const (
	feedbackName = "go.micro.srv.feedback"
)

func feedbackHandler(c *gin.Context) {
	action := c.Param("action")
	log.Printf("action:%s", action)
	switch action {
	case "add":
		addFeedback(c)
	default:
		c.JSON(http.StatusOK, gin.H{"errno": 101, "desc": "unknown action"})
	}
	return
}

func addFeedback(c *gin.Context) {
	var json feedback.Request
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	cl := feedback.NewFeedbackClient(feedbackName, client.DefaultClient)
	rsp, err := cl.Add(context.Background(), &json)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
	} else {
		log.Printf("rsp:%+v", rsp)
		if rsp.Code == 0 {
			c.JSON(http.StatusOK, gin.H{"errno": 0})
		} else {
			c.JSON(http.StatusOK, gin.H{"errno": 1,
				"desc": fmt.Sprintf("failed code:%d", rsp.Code)})
		}
	}
}

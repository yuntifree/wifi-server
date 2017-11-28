package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	verify "github.com/yuntifree/wifi-server/proto/verify"
)

func checkToken(c *gin.Context) error {
	var req verify.TokenRequest
	req.Uid = int64(getCookieInt(c, "u"))
	req.Token = getCookieStr(c, "s")
	cl := verify.NewVerifyClient(verifyName, client.DefaultClient)
	_, err := cl.CheckToken(context.Background(), &req)
	if err == nil {
		if c.Keys == nil {
			c.Keys = make(map[string]interface{})
		}
		c.Keys["uid"] = req.Uid
	}
	return err
}

package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getCookieStr(c *gin.Context, name string) string {
	v, err := c.Cookie(name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": errParam, "desc": "illegal param"})
		panic(fmt.Sprintf("get cookie %s failed:%v", name, err))
	}
	return v
}

func getCookieInt(c *gin.Context, name string) int {
	v, err := c.Cookie(name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": errParam, "desc": "illegal param"})
		panic(fmt.Sprintf("get cookie %s failed:%v", name, err))
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": errParam, "desc": "illegal param"})
		panic(fmt.Sprintf("get cookie %s %s type failed:%v", name, v, err))
	}
	return n
}

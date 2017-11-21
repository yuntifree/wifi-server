package main

import (
	"context"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"github.com/yuntifree/components/weixin"
	"github.com/yuntifree/wifi-server/accounts"
	pay "github.com/yuntifree/wifi-server/proto/pay"
)

const (
	payName  = "go.micro.srv.pay"
	succRsp  = "<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>"
	failRsp  = "<xml><return_code><![CDATA[FAIL]]></return_code><return_msg><![CDATA[SERVER ERROR]]></return_msg></xml>"
	succCode = "SUCCESS"
)

func payHandler(c *gin.Context) {
	action := c.Param("action")
	log.Printf("action:%s", action)
	switch action {
	case "wx_pay":
		wxPay(c)
	case "wx_pay_callback":
		wxPayCB(c)
	default:
		c.JSON(http.StatusOK, gin.H{"errno": errAction, "desc": "unknown action"})
	}
	return
}

type payReq struct {
	Id    int64 `json:"id"`
	Price int64 `json:"price"`
}

func wxPay(c *gin.Context) {
	var p payReq
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 1, "desc": err.Error()})
		return
	}
	wid := getCookieInt(c, "wid")
	uid := getCookieInt(c, "u")
	var req pay.WxPayRequest
	req.Wid = int64(wid)
	req.Uid = int64(uid)
	req.Item = p.Id
	req.Price = p.Price
	ip := c.ClientIP()
	callback := strings.Replace(c.Request.RequestURI, "wx_pay", "wx_pay_callback", -1)
	req.Clientip = ip
	req.Callback = callback
	cl := pay.NewPayClient(payName, client.DefaultClient)
	rsp, err := cl.WxPay(context.Background(), &req)
	if err != nil {
		log.Printf("wxPay WxPay failed:%d %v", wid, err)
		c.JSON(http.StatusOK, gin.H{"errno": errParam, "desc": "illegal param"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": 0, "paySign": rsp.Sign,
		"package": rsp.Pack, "timeStamp": rsp.Ts,
		"nonceStr": rsp.Nonce, "signType": rsp.Signtype})
}

func wxPayCB(c *gin.Context) {
	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("wxPayCB read body failed:%v", err)
		c.Data(http.StatusOK, "application/xml", []byte(failRsp))
		return
	}
	var notify weixin.NotifyRequest
	err = xml.Unmarshal(buf, &notify)
	if err != nil {
		log.Printf("wxPayCB Unmarshal xml failed:%v", err)
		c.Data(http.StatusOK, "application/xml", []byte(failRsp))
		return
	}
	if notify.ReturnCode != succCode || notify.ResultCode != succCode {
		log.Printf("wxPayCB failed response:%+v", notify)
		c.Data(http.StatusOK, "application/xml", []byte(succRsp))
		return
	}

	wx := weixin.WxPay{MerID: accounts.WxMerID,
		MerKey: accounts.WxMerKey}
	if !wx.VerifyNotify(notify) {
		log.Printf("wxPayCB VerifyNotify failed:%+v", notify)
		c.Data(http.StatusOK, "application/xml", []byte(failRsp))
		return
	}

	c.Data(http.StatusOK, "application/xml", []byte(succRsp))
	var req pay.WxCBRequest
	req.Oid = notify.OutTradeNO
	req.Fee = notify.TotalFee
	cl := pay.NewPayClient(payName, client.DefaultClient)
	_, err = cl.WxPayCB(context.Background(), &req)
	if err != nil {
		log.Printf("wxPayCB WxPayCB failed:%s %v", req.Oid, err)
		return
	}
}

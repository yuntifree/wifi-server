package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"github.com/yuntifree/components/strutil"
	"github.com/yuntifree/components/weixin"
	"github.com/yuntifree/wifi-server/accounts"
	pay "github.com/yuntifree/wifi-server/proto/pay"
)

const (
	payName  = "go.micro.srv.pay"
	succRsp  = "<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>"
	failRsp  = "<xml><return_code><![CDATA[FAIL]]></return_code><return_msg><![CDATA[SERVER ERROR]]></return_msg></xml>"
	succCode = "SUCCESS"
	payHost  = "http://wxdev.yunxingzh.com"
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
	if err := c.BindJSON(&p); err != nil {
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
	log.Printf("callback:%s", callback)
	req.Clientip = ip
	req.Callback = payHost + callback
	cl := pay.NewPayClient(payName, client.DefaultClient)
	rsp, err := cl.WxPay(context.Background(), &req)
	if err != nil {
		log.Printf("wxPay WxPay failed:%d %v", wid, err)
		c.JSON(http.StatusOK, gin.H{"errno": errParam, "desc": "illegal param"})
		return
	}
	log.Printf("rsp:%v", rsp)
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

func getJsapiSign(c *gin.Context) {
	var req pay.TicketRequest
	cl := pay.NewPayClient(payName, client.DefaultClient)
	rsp, err := cl.GetTicket(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errno": errParam, "desc": err.Error()})
		return
	}
	noncestr := genNonce()
	ts := time.Now().Unix()
	url := c.Request.Referer()
	pos := strings.Index(url, "#")
	if pos != -1 {
		url = url[:pos]
	}

	ori := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		rsp.Ticket, noncestr, ts, url)
	sign := strutil.Sha1(ori)
	log.Printf("origin:%s sign:%s\n", ori, sign)
	out := fmt.Sprintf("var wx_cfg={\"debug\":false, \"appId\":\"%s\",\"timestamp\":%d,\"nonceStr\":\"%s\",\"signature\":\"%s\",\"jsApiList\":[],\"jsapi_ticket\":\"%s\"};",
		accounts.DgWxAppid, ts, noncestr, sign, rsp.Ticket)
	log.Printf("out:%s", out)
	c.Data(http.StatusOK, "", []byte(out))
}

func genNonce() string {
	nonce := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var res []byte
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < 12; i++ {
		ch := nonce[r.Int31n(int32(len(nonce)))]
		res = append(res, ch)
	}
	return string(res)
}

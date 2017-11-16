package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	micro "github.com/micro/go-micro"
	"github.com/yuntifree/components/strutil"
	"github.com/yuntifree/components/weixin"
	"github.com/yuntifree/wifi-server/accounts"
	"github.com/yuntifree/wifi-server/dbutil"
	pay "github.com/yuntifree/wifi-server/proto/pay"
)

const (
	tradeType = "JSAPI"
	succCode  = "SUCCESS"
	signType  = "MD5"
)

var db *sql.DB

//Server server implement
type Server struct{}

//WxPay weixin pay
func (s *Server) WxPay(ctx context.Context, req *pay.WxPayRequest,
	rsp *pay.WxPayResponse) error {
	log.Printf("WxPay request:%+v", req)
	oid := weixin.GenOrderID(req.Uid)
	id, err := recordOrderInfo(db, oid, req)
	if err != nil {
		log.Printf("WxPay recordOrderInfo failed:%v", err)
		return err
	}
	openid, err := getUserOpenid(db, req.Uid)
	if err != nil {
		log.Printf("getUserOpenid failed:%d %v", req.Uid, err)
		return err
	}

	var rq weixin.UnifyOrderReq
	rq.Appid = accounts.DgWxAppid
	rq.Body = "上网费"
	rq.MchID = accounts.WxMerID
	rq.NonceStr = strutil.GenSalt()
	rq.Openid = openid
	rq.TradeType = tradeType
	rq.SpbillCreateIP = req.Clientip
	rq.TotalFee = req.Price
	rq.OutTradeNO = oid
	rq.NotifyURL = req.Callback

	wx := weixin.WxPay{MerID: accounts.WxMerID,
		MerKey: accounts.WxMerKey}
	resp, err := wx.UnifyPayRequest(rq)
	if err != nil {
		log.Printf("WxPay UnifyPayRequest failed:%v", err)
		return err
	}
	log.Printf("resp:%+v", resp)
	if resp.ReturnCode != succCode || resp.ResultCode != succCode {
		log.Printf("WxPay UnifyPayRequest failed msg:%s", resp.ReturnMsg)
		return fmt.Errorf("pay failed:%s", resp.ReturnMsg)
	}

	now := time.Now().Unix()

	m := make(map[string]interface{})
	m["appId"] = resp.Appid
	m["nonceStr"] = resp.NonceStr
	m["package"] = "prepay_id=" + resp.PrepayID
	m["signType"] = "MD5"
	m["timeStamp"] = now
	sign := wx.CalcSign(m)

	recordPrepayid(db, id, resp.PrepayID)

	rsp.Pack = "prepay_id=" + resp.PrepayID
	rsp.Nonce = resp.NonceStr
	rsp.Ts = now
	rsp.Sign = sign
	rsp.Signtype = signType

	return nil
}

func recordOrderInfo(db *sql.DB, oid string, req *pay.WxPayRequest) (int64, error) {
	res, err := db.Exec(`INSERT INTO orders(oid, uid, wid, item, price, ctime)
	VALUES (?, ?, ?, ?, ?, NOW())`, oid, req.Uid, req.Wid, req.Item,
		req.Price)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func getUserOpenid(db *sql.DB, uid int64) (string, error) {
	var openid string
	err := db.QueryRow("SELECT username FROM users WHERE uid = ?", uid).
		Scan(&openid)
	return openid, err
}

func recordPrepayid(db *sql.DB, id int64, prepayid string) {
	_, err := db.Exec("UPDATE orders SET prepayid = ? WHERE id = ?",
		prepayid, id)
	if err != nil {
		log.Printf("recordPrepayid failed:%d %s %v", id, prepayid, err)
	}
}

func main() {
	var err error
	db, err = dbutil.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	service := micro.NewService(
		micro.Name("go.micro.srv.pay"),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
	)

	service.Init()

	pay.RegisterPayHandler(service.Server(), new(Server))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

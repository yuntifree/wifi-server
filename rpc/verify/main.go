package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	micro "github.com/micro/go-micro"
	"github.com/yuntifree/wifi-server/dbutil"
	verify "github.com/yuntifree/wifi-server/proto/verify"
	"golang.org/x/net/context"
)

const (
	adImg     = "http://img.yunxingzh.com/115cebf5-2ad3-458f-bc2c-48c667eacd52.png"
	wxAppid   = "wx0898ab51f688ee64"
	wxSecret  = "bf430af449b70efc04f11964bc5968a3"
	wxShopid  = "3535655"
	wxAuthURL = "http://wx.yunxingzh.com/auth"
)

var db *sql.DB

//Server server implement
type Server struct{}

//CheckLogin check login
func (s *Server) CheckLogin(ctx context.Context, req *verify.CheckRequest,
	rsp *verify.CheckResponse) error {
	rsp.Autologin = isAutoMac(db, req.Usermac, req.Apmac)
	rsp.Img = adImg
	rsp.Wxappid = wxAppid
	rsp.Wxsecret = wxSecret
	rsp.Wxshopid = wxShopid
	rsp.Wxauthurl = wxAuthURL
	rsp.Taobao = 0
	rsp.Logintype = 1

	return nil
}

func isAutoMac(db *sql.DB, usermac, apmac string) int64 {
	var phone string
	err := db.QueryRow(`SELECT phone FROM user_mac WHERE mac = ?`, usermac).
		Scan(&phone)
	if err != nil || phone == "" {
		return 0
	}

	var park int64
	err = db.QueryRow(`SELECT park FROM wifi_account WHERE phone = ?`, phone).
		Scan(&park)
	if err != nil || park == 0 {
		return 0
	}

	var epark int64
	err = db.QueryRow(`SELECT park FROM ap_info WHERE mac = ?`, apmac).
		Scan(&epark)
	if err != nil || epark == 0 {
		return 0
	}
	if epark == park {
		return 1
	}
	return 0
}

func main() {
	var err error
	db, err = dbutil.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	service := micro.NewService(
		micro.Name("go.micro.srv.verify"),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
	)

	service.Init()

	verify.RegisterVerifyHandler(service.Server(), new(Server))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

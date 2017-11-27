package main

import (
	"database/sql"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	micro "github.com/micro/go-micro"
	uuid "github.com/satori/go.uuid"
	"github.com/yuntifree/components/weixin"
	"github.com/yuntifree/wifi-server/accounts"
	"github.com/yuntifree/wifi-server/dbutil"
	wx "github.com/yuntifree/wifi-server/proto/wx"
	"golang.org/x/net/context"
)

var db *sql.DB

//Server server implement
type Server struct{}

//Login weixin login
func (s *Server) Login(ctx context.Context, req *wx.LoginRequest, rsp *wx.LoginResponse) error {
	wx := weixin.WxInfo{Appid: accounts.DgWxAppid,
		Appkey: accounts.DgWxAppkey}
	uinfo, err := wx.GetCodeToken(req.Code)
	if err != nil {
		return err
	}
	token := genSalt()
	res, err := db.Exec(`INSERT IGNORE INTO users(username, token, ctime) VALUES 
	(?, ?, NOW())`, uinfo.Openid, token)
	if err != nil {
		log.Printf("insert user record failed:%s %v", uinfo.Openid, err)
		return err
	}
	uid, err := res.LastInsertId()
	if err != nil {
		log.Printf("get insert id failed:%v", err)
		return err
	}

	if uid == 0 {
		err = db.QueryRow("SELECT uid FROM users WHERE username = ?", uinfo.Openid).
			Scan(&uid)
		if err != nil {
			log.Printf("scan username failed:%s %v", uinfo.Openid, err)
			return err
		}
		_, err = db.Exec(`UPDATE users SET token = ? WHERE uid = ?`, token, uid)
		if err != nil {
			log.Printf("update token failed:%v", err)
			return err
		}
	}
	rsp.Uid = uid
	rsp.Token = token
	return nil

}

func genSalt() string {
	u := uuid.NewV4()
	return strings.Join(strings.Split(u.String(), "-"), "")
}

func main() {
	var err error
	db, err = dbutil.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	service := micro.NewService(
		micro.Name("go.micro.srv.wx"),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
	)

	service.Init()

	wx.RegisterWxHandler(service.Server(), new(Server))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

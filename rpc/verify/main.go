package main

import (
	"Server/zte"
	"database/sql"
	"errors"
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

//GetCheckCode get check code
func (s *Server) GetCheckCode(ctx context.Context, req *verify.CodeRequest,
	rsp *verify.CodeResponse) error {
	if !checkPhonePark(db, req.Phone, req.Apmac) {
		return errors.New("请先到公众号开通上网服务")
	}
	var stype uint
	stype = getAcSys(db, req.Acname)
	if isExceedCodeFrequency(db, req.Phone, stype) {
		log.Printf("GetCheckCode isExceedCodeFrequency phone:%s", req.Phone)
		return errors.New("超过频率限制")
	}
	code, err := zte.Register(req.Phone, true, stype)
	if err != nil {
		log.Printf("GetCheckCode Register failed:%v", err)
		return errors.New("账号注册失败")
	}
	log.Printf("recordZteCode phone:%s code:%s type:%d", req.Phone, code, stype)
	recordZteCode(db, req.Phone, code, stype)
	return nil
}

func checkPhonePark(db *sql.DB, phone, apmac string) bool {
	var park, epark int64
	err := db.QueryRow(`SELECT park FROM wifi_account WHERE phone = ?`, phone).
		Scan(&park)
	if err != nil {
		return false
	}
	err = db.QueryRow(`SELECT park FROM ap_info WHERE mac = ?`, apmac).
		Scan(&epark)
	if err != nil {
		return false
	}
	if park == epark {
		return true
	}
	return false
}

func recordZteCode(db *sql.DB, phone, code string, stype uint) {
	if code == "" {
		return
	}
	_, err := db.Exec(`INSERT INTO zte_code(phone, code, type, ctime, mtime) 
		VALUES (?, ?, ?, NOW(), NOW()) ON DUPLICATE KEY UPDATE code = ?, 
		mtime = NOW()`,
		phone, code, stype, code)
	if err != nil {
		log.Printf("recordZteCode query failed:%s %s %d %v", phone, code, stype, err)
	}
}

func getAcSys(db *sql.DB, acname string) uint {
	var stype uint
	err := db.QueryRow("SELECT type FROM ac_info WHERE name = ?", acname).
		Scan(&stype)
	if err != nil {
		log.Printf("getAcSys query failed:%v", err)
	}
	return stype
}

func isExceedCodeFrequency(db *sql.DB, phone string, stype uint) bool {
	var flag int
	err := db.QueryRow(`SELECT IF(NOW() > DATE_ADD(mtime, INTERVAL 5 MINUTE), 
				0, 1) FROM zte_code WHERE phone = ? AND type = ?`,
		phone, stype).Scan(&flag)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("isExceedCodeFrequency query failed:%v", err)
		return false
	}
	if flag > 0 {
		return true
	}
	return false
}

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

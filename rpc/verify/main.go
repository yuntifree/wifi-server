package main

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	micro "github.com/micro/go-micro"
	"github.com/yuntifree/components/zte"
	"github.com/yuntifree/wifi-server/dbutil"
	verify "github.com/yuntifree/wifi-server/proto/verify"
	"golang.org/x/net/context"
)

const (
	adImg       = "http://img.yunxingzh.com/115cebf5-2ad3-458f-bc2c-48c667eacd52.png"
	wxAppid     = "wx0898ab51f688ee64"
	wxSecret    = "bf430af449b70efc04f11964bc5968a3"
	wxShopid    = "3535655"
	wxAuthURL   = "http://wx.yunxingzh.com/auth"
	testAcname  = "2043.0769.200.00"
	testAcip    = "120.197.159.10"
	testUserip  = "10.96.72.28"
	testUsermac = "f45c89987347"
	portalDir   = "http://api.yunxingzh.com/portal0406201704201946/index0406.html"
	tstUid      = 137
	tstToken    = "6ba9ac5a422d4473b337d57376dd3488"
)

var db *sql.DB

//Server server implement
type Server struct{}

//OneClickLogin one-click login
func (s *Server) OneClickLogin(ctx context.Context, req *verify.OneClickRequest,
	rsp *verify.LoginResponse) error {
	var phone string
	usermac := strings.Replace(req.Wlanusermac, ":", "", -1)
	err := db.QueryRow(`SELECT phone FROM user_mac WHERE mac = ? LIMIT 1`, usermac).
		Scan(&phone)
	if err != nil {
		log.Printf("OneClickLogin query phone failed:%v", err)
		return err
	}
	stype := getAcSys(db, req.Wlanacname)
	bitmap := getWifiBitmap(db, phone)
	err = checkZteReg(db, bitmap, stype, phone)
	if err != nil {
		log.Printf("checkZteReg failed:%v", err)
		return err
	}
	if !isTestClickParam(req) {
		flag, err := zteLogin(phone, req.Wlanuserip,
			req.Wlanusermac, req.Wlanacip, req.Wlanacname, stype)
		if !flag {
			log.Printf("OneClickLogin zteLogin retry failed, phone:%s ",
				phone)
			return errors.New("登录失败")
		}
		if err != nil {
			return err
		}
	}

	var uid int64
	var token string
	err = db.QueryRow(`SELECT uid, token FROM users u, wifi_account w WHERE 
	u.wifi = w.id AND w.phone = ?`, phone).Scan(&uid, &token)
	if err != nil {
		return err
	}
	rsp.Uid = tstUid
	rsp.Token = tstToken
	rsp.Portaldir = portalDir
	rsp.Portaltype = 1
	return nil

}

func getWifiBitmap(db *sql.DB, phone string) uint {
	var bitmap uint
	err := db.QueryRow("SELECT bitmap FROM wifi_account WHERE phone = ?", phone).
		Scan(&bitmap)
	if err != nil {
		log.Printf("getUserBitmap failed:%v", err)
	}
	return bitmap
}

func checkZteReg(db *sql.DB, bitmap, stype uint, phone string) error {
	if bitmap&(1<<stype) == 0 {
		code, err := zte.Register(phone, true, stype)
		if err != nil {
			log.Printf("PortalLogin zte register failed:%v", err)
			return err
		}
		recordZteCode(db, phone, code, stype)
		updateWifiBitmap(db, phone, (1 << stype))
	}
	return nil
}

func updateWifiBitmap(db *sql.DB, phone string, bitmap uint) {
	_, err := db.Exec("UPDATE wifi_account SET bitmap = bitmap | ? WHERE phone = ?",
		bitmap, phone)
	if err != nil {
		log.Printf("updateUserBitmap failed, phone:%s %v", phone, err)
	}
}

//PortalLogin portal login
func (s *Server) PortalLogin(ctx context.Context, req *verify.PortalLoginRequest,
	rsp *verify.LoginResponse) error {
	if !checkPhonePark(db, req.Phone, req.Wlanapmac) {
		return errors.New("请先到公众号开通上网服务")
	}
	stype := getAcSys(db, req.Wlanacname)
	if !checkZteCode(db, req.Phone, req.Code, stype) {
		log.Printf("PortalLogin checkZteCode failed, phone:%s code:%s stype:%d",
			req.Phone, req.Code, stype)
		return errors.New("验证码错误")

	}
	if !isTestParam(req) {
		flag, err := zteLogin(req.Phone, req.Wlanuserip,
			req.Wlanusermac, req.Wlanacip, req.Wlanacname, stype)
		if !flag {
			log.Printf("PortalLogin zteLogin retry failed, phone:%s code:%s",
				req.Phone, req.Code)
			return errors.New("登录失败")
		}
		if err != nil {
			return err
		}
	}
	var uid int64
	var token string
	err := db.QueryRow(`SELECT uid, token FROM users u, wifi_account w WHERE 
	u.wifi = w.id AND w.phone = ?`, req.Phone).Scan(&uid, &token)
	if err != nil {
		return err
	}
	recordUserMac(db, req.Wlanusermac, req.Phone)
	rsp.Uid = uid
	rsp.Token = token
	rsp.Portaldir = portalDir
	return nil
}

func recordUserMac(db *sql.DB, mac, phone string) {
	mac = strings.Replace(mac, ":", "", -1)
	_, err := db.Exec(`INSERT INTO user_mac(mac, phone, ctime) 
	VALUES (?, ?, NOW()) ON DUPLICATE KEY UPDATE mac = ?`,
		mac, phone, mac)
	if err != nil {
		log.Printf("recordUserMac failed mac:%s phone:%s err:%v",
			mac, phone, err)
	}
}

func zteLogin(phone, userip, usermac, acip, acname string, stype uint) (bool, error) {
	flag, err := zte.Loginnopass(phone, userip, usermac, acip, acname, stype)
	if flag {
		return true, nil
	}
	log.Printf("PortalLogin zte loginnopass failed, phone:%s stype:%d",
		phone, stype)
	return flag, err
}

func isTestClickParam(info *verify.OneClickRequest) bool {
	if info.Wlanacip == testAcip &&
		info.Wlanuserip == testUserip && info.Wlanusermac == testUsermac {
		return true
	}
	return false
}

func isTestParam(info *verify.PortalLoginRequest) bool {
	if info.Wlanacip == testAcip &&
		info.Wlanuserip == testUserip && info.Wlanusermac == testUsermac {
		return true
	}
	return false
}

func checkZteCode(db *sql.DB, phone, code string, stype uint) bool {
	var eCode string
	err := db.QueryRow("SELECT code FROM zte_code WHERE type = ? AND phone = ?",
		stype, phone).Scan(&eCode)
	if err != nil {
		log.Printf("checkZteCode query failed:%s %s %v", phone, code, err)
		return false
	}
	if eCode == code {
		return true
	}
	return false
}

//GetCheckCode get check code
func (s *Server) GetCheckCode(ctx context.Context, req *verify.CodeRequest,
	rsp *verify.CodeResponse) error {
	if !checkPhonePark(db, req.Phone, req.Wlanapmac) {
		log.Printf("GetCheckCode checkPhonePark failed:%+v", req)
		return errors.New("请先到公众号开通上网服务")
	}
	var stype uint
	stype = getAcSys(db, req.Wlanacname)
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
	log.Printf("park:%d real park:%d", park, epark)
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
	rsp.Autologin = isAutoMac(db, req.Wlanusermac, req.Wlanapmac)
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

//CheckToken check user token
func (s *Server) CheckToken(ctx context.Context, req *verify.TokenRequest,
	rsp *verify.TokenResponse) error {
	var etoken string
	err := db.QueryRow(`SELECT token FROM users WHERE uid = ?`, req.Uid).
		Scan(&etoken)
	if err != nil {
		log.Printf("CheckToken query token failed:%d %v", req.Uid, err)
		return err
	}
	if req.Token != etoken {
		return errors.New("token验证失败")
	}
	return nil
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

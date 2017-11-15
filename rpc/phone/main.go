package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
	micro "github.com/micro/go-micro"
	"github.com/yuntifree/components/sms"
	"github.com/yuntifree/wifi-server/accounts"
	"github.com/yuntifree/wifi-server/dbutil"
	phone "github.com/yuntifree/wifi-server/proto/phone"
)

const (
	mastercode = 251653
)

var db *sql.DB

//Server server implement
type Server struct{}

//CheckCode check sms code
func (s *Server) CheckCode(ctx context.Context, req *phone.CheckRequest, rsp *phone.CheckResponse) error {
	if req.Code == mastercode {
		return nil
	}

	var realcode, id int64
	err := db.QueryRow(`SELECT id, code FROM phone_code WHERE phone = ? AND 
	used = 0 AND etime > NOW() ORDER BY id DESC LIMIT 1`, req.Phone).
		Scan(&id, &realcode)
	if err != nil {
		log.Printf("CheckCode query failed:%s %v", req.Phone, err)
		return err
	}

	if realcode == req.Code {
		_, err = db.Exec("UPDATE phone_code SET used = 1 WHERE id = ?", id)
		if err != nil {
			log.Printf("update phone_code used failed:%s %v", req.Phone, err)
		}
		return nil
	}
	return fmt.Errorf("illegal phone code:%s %d-%d", req.Phone, req.Code,
		realcode)
}

func sendSMS(phone string, code int) int {
	yp := sms.Yunpian{Apikey: accounts.YPSMSApikey,
		TplID: accounts.YPSMSTplID}
	return yp.Send(phone, code)
}

//GetCode send sms code
func (s *Server) GetCode(ctx context.Context, req *phone.GetRequest, rsp *phone.GetResponse) error {
	var code int
	err := db.QueryRow(`SELECT code FROM phone_code WHERE phone = ?
	AND used = 0 AND etime > NOW() AND
	timestampdiff(second, stime, now()) < 300 ORDER BY id DESC LIMIT 1`,
		req.Phone).Scan(&code)
	if err != nil {
		code = genCode()
		_, err := db.Exec(`INSERT INTO phone_code(phone, code, ctime,
		stime, etime) VALUES (?, ?, NOW(), NOW(), DATE_ADD(NOW(), INTERVAL 5 MINUTE))`,
			req.Phone, code)
		if err != nil {
			log.Printf("insert into phone_code failed:%s %v", req.Phone, err)
			return err
		}
		ret := sendSMS(req.Phone, code)
		if ret != 0 {
			log.Printf("send sms code failed:%s %d", req.Phone, ret)
			return fmt.Errorf("send sms failed:%d", ret)
		}
		return nil
	}
	if code > 0 {
		ret := sendSMS(req.Phone, code)
		if ret != 0 {
			log.Printf("send sms failed:%s %d", req.Phone, ret)
			return fmt.Errorf("send sms failed:%d", ret)
		}
	}
	return nil
}

func genCode() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return int(r.Int31n(1e6))
}

func main() {
	var err error
	db, err = dbutil.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	service := micro.NewService(
		micro.Name("go.micro.srv.phone"),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
	)

	service.Init()

	phone.RegisterPhoneHandler(service.Server(), new(Server))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

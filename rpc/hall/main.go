package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	micro "github.com/micro/go-micro"
	"github.com/yuntifree/wifi-server/dbutil"
	hall "github.com/yuntifree/wifi-server/proto/hall"
	"golang.org/x/net/context"
)

var db *sql.DB

//Server server implement
type Server struct{}

//Login login hall
func (s *Server) Login(ctx context.Context, req *hall.LoginRequest,
	rsp *hall.LoginResponse) error {
	res, err := db.Exec(`INSERT IGNORE INTO wifi_account(phone, park, ctime)
	VALUES (?, ?, NOW())`, req.Phone, req.Park)
	if err != nil {
		log.Printf("Login insert failed:%s %v", req.Phone, err)
		return err
	}
	wid, err := res.LastInsertId()
	if err != nil {
		log.Printf("Login get insert id failed:%v", err)
		return err
	}
	if wid == 0 {
		err = db.QueryRow("SELECT id FROM wifi_account WHERE phone = ?",
			req.Phone).Scan(&wid)
		if err != nil {
			log.Printf("Login scan failed:%s %v", req.Phone, err)
			return err
		}
	}
	rsp.Wid = wid
	return nil
}

func main() {
	var err error
	db, err = dbutil.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	service := micro.NewService(
		micro.Name("go.micro.srv.hall"),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
	)

	service.Init()

	hall.RegisterHallHandler(service.Server(), new(Server))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

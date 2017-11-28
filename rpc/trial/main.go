package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	micro "github.com/micro/go-micro"
	"github.com/yuntifree/wifi-server/dbutil"
	trial "github.com/yuntifree/wifi-server/proto/trial"
	"golang.org/x/net/context"
)

const (
	trialDays = 2
	payedDays = 1
	trialBit  = 0x1
	payBit    = 0x2
)

var db *sql.DB

//Server server implement
type Server struct{}

//Info get trial info
func (s *Server) Info(ctx context.Context, req *trial.Request,
	rsp *trial.InfoResponse) error {
	var bitmap int64
	err := db.QueryRow("SELECT bitmap FROM wifi_account WHERE id = ?",
		req.Wid).Scan(&bitmap)
	if err != nil {
		log.Printf("Info query failed:%d %v", req.Wid, err)
		return err
	}
	rsp.Used = bitmap & trialBit
	return nil
}

//Apply apply for trial
func (s *Server) Apply(ctx context.Context, req *trial.Request,
	rsp *trial.ApplyResponse) error {
	var bitmap, eflag int64
	err := db.QueryRow("SELECT bitmap, IF(etime > NOW(), 1, 0) FROM wifi_account WHERE id = ?",
		req.Wid).Scan(&bitmap, &eflag)
	if err != nil {
		log.Printf("Apply query failed:%d %v", req.Wid, err)
		return err
	}
	if bitmap&trialBit == trialBit {
		return nil
	}
	_, err = db.Exec(`INSERT INTO trial(wid, ctime, etime) VALUES (?, 
	NOW(), DATE_ADD(NOW(), INTERVAL 2 DAY))`, req.Wid)
	if err != nil {
		log.Printf("Apply insert trial failed:%d %v", req.Wid, err)
		return err
	}
	if (bitmap&payBit == payBit) && eflag == 1 {
		_, err = db.Exec(`UPDATE wifi_account SET bitmap = bitmap | 0x1, 
	etime = DATE_ADD(etime, INTERVAL ? DAY) WHERE id = ?`,
			payedDays, req.Wid)
		if err != nil {
			log.Printf("Apply update wifi_account failed:%d %v", req.Wid, err)
			return err
		}
		return nil
	}
	_, err = db.Exec(`UPDATE wifi_account SET bitmap = bitmap | 0x1, 
	etime = DATE_ADD(NOW(), INTERVAL ? DAY) WHERE id = ?`,
		trialDays, req.Wid)
	if err != nil {
		log.Printf("Apply update wifi_account failed:%d %v", req.Wid, err)
		return err
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
		micro.Name("go.micro.srv.trial"),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
	)

	service.Init()

	trial.RegisterTrialHandler(service.Server(), new(Server))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

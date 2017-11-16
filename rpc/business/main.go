package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	micro "github.com/micro/go-micro"
	"github.com/yuntifree/wifi-server/dbutil"
	business "github.com/yuntifree/wifi-server/proto/business"
	"golang.org/x/net/context"
)

const (
	payedBit = 0x2
)

var db *sql.DB

//Server server implement
type Server struct{}

//Info get business info
func (s *Server) Info(ctx context.Context, req *business.Request,
	rsp *business.InfoResponse) error {
	var bitmap int64
	var etime string
	err := db.QueryRow("SELECT bitmap, etime FROM wifi_account WHERE id = ?",
		req.Wid).Scan(&bitmap, &etime)
	if err != nil {
		log.Printf("Info query failed:%d %v", req.Wid, err)
		return err
	}
	rsp.Payed = bitmap & payedBit
	if rsp.Payed == payedBit {
		rsp.Expire = etime
	}
	rsp.Items = getItems(db)
	return nil
}

func getItems(db *sql.DB) []*business.Item {
	rows, err := db.Query(`SELECT id, title, price FROM items WHERE deleted = 0
	AND online = 1 ORDER BY id`)
	if err != nil {
		log.Printf("getItems failed:%v", err)
		return nil
	}
	defer rows.Close()
	var infos []*business.Item
	for rows.Next() {
		var it business.Item
		err = rows.Scan(&it.Id, &it.Title, &it.Price)
		if err != nil {
			continue
		}
		infos = append(infos, &it)
	}
	return infos
}

func main() {
	var err error
	db, err = dbutil.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	service := micro.NewService(
		micro.Name("go.micro.srv.business"),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
	)

	service.Init()

	business.RegisterBusinessHandler(service.Server(), new(Server))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

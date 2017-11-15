package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	micro "github.com/micro/go-micro"
	"github.com/yuntifree/wifi-server/dbutil"
	park "github.com/yuntifree/wifi-server/proto/park"
)

var db *sql.DB

//Server server implement
type Server struct{}

//Get return online park
func (s *Server) Get(ctx context.Context, req *park.GetRequest, rsp *park.GetResponse) error {
	rows, err := db.Query(`SELECT id, name, address FROM park WHERE online = 1 AND 
	deleted = 0 ORDER BY id`)
	if err != nil {
		return err
	}
	defer rows.Close()
	var infos []*park.Info
	for rows.Next() {
		var info park.Info
		err = rows.Scan(&info.Id, &info.Name, &info.Address)
		if err != nil {
			continue
		}
		infos = append(infos, &info)
	}
	rsp.Infos = infos
	return nil
}

func main() {
	var err error
	db, err = dbutil.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	service := micro.NewService(
		micro.Name("go.micro.srv.park"),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
	)

	service.Init()

	park.RegisterParkHandler(service.Server(), new(Server))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

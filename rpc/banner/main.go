package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	micro "github.com/micro/go-micro"
	"github.com/yuntifree/wifi-server/dbutil"
	banner "github.com/yuntifree/wifi-server/proto/banner"
)

var db *sql.DB

//Server server implement
type Server struct{}

//Get return online banners
func (s *Server) Get(ctx context.Context, req *banner.GetRequest, rsp *banner.GetResponse) error {
	rows, err := db.Query(`SELECT id, img, dst FROM banner WHERE online = 1 AND 
	deleted = 0 ORDER BY id`)
	if err != nil {
		return err
	}
	defer rows.Close()
	var infos []*banner.Info
	for rows.Next() {
		var info banner.Info
		err = rows.Scan(&info.Id, &info.Img, &info.Dst)
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
		micro.Name("go.micro.srv.banner"),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
	)

	service.Init()

	banner.RegisterBannerHandler(service.Server(), new(Server))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

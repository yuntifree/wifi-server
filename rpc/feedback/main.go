package main

import (
	"database/sql"
	"log"
	"time"

	micro "github.com/micro/go-micro"
	"github.com/yuntifree/wifi-server/dbutil"
	feedback "github.com/yuntifree/wifi-server/proto/feedback"
	context "golang.org/x/net/context"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

//Server server implement
type Server struct{}

//Add add feedback record
func (s *Server) Add(ctx context.Context, req *feedback.Request, rsp *feedback.Response) error {
	log.Printf("Add request:%+v", req)
	_, err := db.Exec("INSERT INTO feedback(phone, content, ctime) VALUES (?, ?, NOW())",
		req.Phone, req.Content)
	if err != nil {
		log.Printf("Add failed:%v", err)
		rsp.Code = 1
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
		micro.Name("go.micro.srv.feedback"),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
	)

	service.Init()

	feedback.RegisterFeedbackHandler(service.Server(), new(Server))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

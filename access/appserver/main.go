package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/feedback/:action", feedbackHandler)
	router.POST("/banner/:action", bannerHandler)
	router.POST("/park/:action", parkHandler)
	router.POST("/phone/:action", phoneHandler)
	router.POST("/hall/:action", hallHandler)
	router.POST("/trial/:action", trialHandler)
	router.POST("/business/:action", businessHandler)
	router.POST("/pay/:action", payHandler)
	router.POST("/check_login", checkLoginHandler)
	router.POST("/get_check_code", getCodeHandler)
	router.POST("/portal_login", portalLoginHandler)
	router.POST("/one_click_login", oneClickLoginHandler)
	router.GET("/jumpwx", jumpHandler)
	router.GET("/weixin/login", loginHandler)
	router.GET("/pay/get_jsapi_sign", getJsapiSign)
	router.Static("/static/", "/data/wifi/html")

	srv := &http.Server{
		Addr:    ":9898",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen:%s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

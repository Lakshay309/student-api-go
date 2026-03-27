package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Lakshay309/student-api-go/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad()
	// database setup
	// setup router
	router:=http.NewServeMux()

	router.HandleFunc("GET /",func(w http.ResponseWriter,r *http.Request){
		w.Write([]byte("Welcome to student api"))
	})
	// setup server
	server:=http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}
	fmt.Printf("server started %s",cfg.HTTPServer.Addr)
	err:=server.ListenAndServe()
	if err!=nil{
		log.Fatal("failed to start server");
	}
}

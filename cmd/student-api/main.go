package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Lakshay309/student-api-go/internal/config"
	"github.com/Lakshay309/student-api-go/internal/http/handlers/student"
	"github.com/Lakshay309/student-api-go/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()
	// database setup
	storage,err:=sqlite.New(cfg)
	if err!=nil{
		log.Fatal(err);
	}
	slog.Info("storage initilized",slog.String("env",cfg.Env),slog.String("version","1.0.0"))

	// setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/student", student.New(storage))
	router.HandleFunc("GET /api/student/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/student", student.GetList(storage))

	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("server started ", slog.String("Address", cfg.Addr))
	// fmt.Printf("server started %s",cfg.HTTPServer.Addr)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()
	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("server shutdown successfully")
}

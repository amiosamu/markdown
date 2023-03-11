package main

import (
	"context"
	"github.com/amiosamu/markdown/pkg/database"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	log.Println("starting server...")

	db, err := database.InitDB()

	if err != nil {
		log.Fatalf("unable to init database: %v\n", err)
	}

	router, err := inject(db)

	if err != nil {
		log.Fatalf("unable to inject data sources %v\n", err)
	}
	//
	//engine := gin.Default()
	//
	//handler.NewHandler(&handler.Config{
	//	Engine: engine,
	//})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to init server: %v\n", err)
		}
	}()
	log.Printf("Listening on port %v\n", srv.Addr)
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.Close(); err != nil {
		log.Fatalf("a problem occured gracefully shuttding down data sources: %v\n", err)
	}

	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

}

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/maragudk/env"

	"github.com/maragudk/sqlite-app/http"
	"github.com/maragudk/sqlite-app/sql"
)

func main() {
	log := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)
	log.Println("Starting")

	_ = env.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	db := sql.NewDatabase(sql.NewDatabaseOptions{
		Log:                   log,
		URL:                   env.GetStringOrDefault("DATABASE_URL", "file:app.db"),
		MaxOpenConnections:    5,
		MaxIdleConnections:    5,
		ConnectionMaxLifetime: time.Hour,
		ConnectionMaxIdleTime: time.Hour,
	})

	if err := db.Connect(); err != nil {
		log.Fatalln("Error connecting to database:", err)
	}

	s := http.NewServer(http.NewServerOptions{
		Database: db,
		Host:     env.GetStringOrDefault("HOST", ""),
		Log:      log,
		Port:     env.GetIntOrDefault("PORT", 8080),
	})

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		if err := s.Start(); err != nil {
			log.Fatalln("Error starting server:", err)
		}
		wg.Done()
	}()

	<-ctx.Done()
	log.Println("Stopping")

	if err := s.Stop(); err != nil {
		log.Fatalln("Error stopping server:", err)
	}

	wg.Wait()
}

package main

import (
	"context"
	"log"
	"os"

	"github.com/maragudk/env"
	"github.com/maragudk/migrate"

	"github.com/maragudk/litefs-app/sql"
)

func main() {
	log := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)

	_ = env.Load()

	db := sql.NewDatabase(sql.NewDatabaseOptions{
		Log:                log,
		URL:                env.GetStringOrDefault("DATABASE_URL", "file:app.db"),
		MaxOpenConnections: 1,
		MaxIdleConnections: 1,
	})

	if err := db.Connect(); err != nil {
		log.Fatalln("Error connecting to database:", err)
	}

	fsys := os.DirFS("sql/migrations")

	if len(os.Args) < 2 {
		log.Fatalln("up or down?")
	}

	switch os.Args[1] {
	case "up":
		if err := migrate.Up(context.Background(), db.DB.DB, fsys); err != nil {
			log.Fatalln(err)
		}
		log.Println("Migrated up")
	case "down":
		if err := migrate.Down(context.Background(), db.DB.DB, fsys); err != nil {
			log.Fatalln(err)
		}
		log.Println("Migrated down")
	default:
		log.Fatalln("unknown command " + os.Args[1])
	}
}

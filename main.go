package main

import (
	"database/sql"
	"log"

	"github.com/abd-rakhman/url-app/api"
	db "github.com/abd-rakhman/url-app/db/sqlc"
	"github.com/abd-rakhman/url-app/utils"
	_ "github.com/lib/pq"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	Store := db.NewStore(conn)

	server := api.NewServer(Store)

	err = server.Run(config.ServerAddress)
	if err != nil {
		log.Fatal(err)
	}
}

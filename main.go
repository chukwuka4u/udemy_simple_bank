package main

import (
	"database/sql"
	"log"

	"simplebank/util"

	"simplebank/api"
	db "simplebank/db/sqlc"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server, error := api.NewServer(config, store)
	if error != nil {
		log.Fatal("cannot create server", error)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}

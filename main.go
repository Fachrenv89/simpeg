package main

import (
	"database/sql"
	"os"

	"github.com/SemmiDev/simpeg/api"
	db "github.com/SemmiDev/simpeg/db/mysql"
	"github.com/SemmiDev/simpeg/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	tokenMaker, err := util.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create jwt maker")
	}

	store := db.NewStore(conn)
	runGinServer(config, tokenMaker, store)
}

func runGinServer(
	config util.Config,
	tokenMaker util.Maker,
	store db.Store,
) {
	server, err := api.NewServer(config, store, tokenMaker)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}

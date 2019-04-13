package main

import (
	"delonakc.com/api/config"
	"delonakc.com/api/database"
	"delonakc.com/api/redis"
	"delonakc.com/api/router"
	"delonakc.com/api/runner"
)

var EnvMode string

func main() {
	conf := config.Get(EnvMode)

	//redis
	redis.Start(conf)

	//mongodb
	db := database.New(conf.DataBase.Host, conf.DataBase.Port, conf.DataBase.Db)

	r := router.NewRouter()

	runner.Run(r, db, conf)

	defer func() {
		db.DB.Close()
		redis.Rd.Close()
	}()
}
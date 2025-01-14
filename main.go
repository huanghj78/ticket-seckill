package main

import (
	"log"
	"ticket-seckill/conf"
	"ticket-seckill/infra/cache"
	"ticket-seckill/infra/db"
	"ticket-seckill/router"
)

func main() {
	conf.Init("./conf.ini")
	db := db.Init()
	defer db.Close()

	cache := cache.Init()
	defer cache.Close()
	app := router.Init()
	if err := app.Run(":" + conf.Conf.Server.Port); err != nil {
		log.Fatalf("App run failed, err:%v", err)
	}
}

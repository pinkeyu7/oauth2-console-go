package main

import (
	"flag"
	"oauth2-console-go/api"
	"oauth2-console-go/config"
	"oauth2-console-go/pkg/logr"
	"oauth2-console-go/pkg/valider"
	"oauth2-console-go/route"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

var port string

// @title Oauth console API
// @version 1.0
// @description Oauth console API
// @termsOfService https://github.com/pinkeyu7/oauth2-console-go
// @license.name MIT
// @license.url
func main() {
	// init http port
	flag.StringVar(&port, "port", "8080", "Initial port number")
	flag.Parse()

	// init config
	config.InitEnv()

	// init logger
	logr.InitLogger()

	// init validation
	valider.Init()

	// init driver
	_ = api.InitXorm()
	_ = api.InitRedis()
	_ = api.InitRedisCluster()

	// init gin router
	r := route.Init()

	// start server
	err := r.Run(":" + port)
	if err != nil {
		log.Println(err)
	}
}

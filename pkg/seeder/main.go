package main

import (
	"fmt"
	"oauth2-console-go/pkg/logr"
	"oauth2-console-go/pkg/seeder/seed"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"xorm.io/xorm"
)

func main() {
	var err error
	remoteBranch := os.Getenv("REMOTE_BRANCH")

	logger := logr.NewLogger()
	if remoteBranch == "" {
		// load env
		err = godotenv.Load()

		if err != nil {
			logger.Debug(err.Error())
		}
	}

	dsn := "%s:%s@(%s:%s)/%s?parseTime=true"

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	engine, err := xorm.NewEngine("mysql", fmt.Sprintf(dsn, dbUser, dbPassword, dbHost, dbPort, dbName))
	if err != nil {
		logger.Error(err.Error())
	}

	engine.TZLocation, _ = time.LoadLocation("UTC")
	engine.DatabaseTZ, _ = time.LoadLocation("UTC")

	gofakeit.Seed(time.Now().Unix())

	// Create Sys Account
	sysAccountSeeds := seed.AllSysAccount()
	run(engine, sysAccountSeeds)
}

func run(engine *xorm.Engine, channelSeeds []seed.Seed) {
	logger := logr.NewLogger()
	for _, seed := range channelSeeds {
		logger.Info(seed.Name)
		err := seed.Run(engine)
		if err != nil {
			logger.Error(seed.Name + " Failed")
			logger.Error(err.Error())
		}
	}
}

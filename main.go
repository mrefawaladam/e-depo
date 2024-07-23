package main

import (
	"context"
	"database/sql"

	usecases "e-depo/src/app/usecases"

	"e-depo/src/infra/config"

	postgres "e-depo/src/infra/persistence/postgres"
	// "e-depo/src/infra/persistence/redis"

	userRepo "e-depo/src/infra/persistence/postgres/user"

	"e-depo/src/interface/rest"

	ms_log "e-depo/src/infra/log"

	userUC "e-depo/src/app/usecases/user"

	_ "github.com/joho/godotenv/autoload"

	"github.com/sirupsen/logrus"
)

func main() {

	ctx := context.Background()

	conf := config.Make()

	isProd := false
	if conf.App.Environment == "PRODUCTION" {
		isProd = true
	}

	m := make(map[string]interface{})
	m["env"] = conf.App.Environment
	m["service"] = conf.App.Name
	logger := ms_log.NewLogInstance(
		ms_log.LogName(conf.Log.Name),
		ms_log.IsProduction(isProd),
		ms_log.LogAdditionalFields(m))

	postgresdb, err := postgres.New(conf.SqlDb, logger)
	// redisClient, err := redis.NewRedisClient(conf.Redis, logger)

	// redisServe := redisServe.NewServRedis(redisClient)

	defer func(l *logrus.Logger, sqlDB *sql.DB, dbName string) {
		err := sqlDB.Close()
		if err != nil {
			l.Errorf("error closing sql database %s: %s", dbName, err)
		} else {
			l.Printf("sql database %s successfuly closed.", dbName)
		}
	}(logger, postgresdb.Conn.DB, postgresdb.Conn.DriverName())

	userRepository := userRepo.NewUserRepository(postgresdb.Conn)

	httpServer, err := rest.New(
		conf.Http,
		isProd,
		logger,
		usecases.AllUseCases{
			UserUC: userUC.NewUserUseCase(userRepository),
		},
	)
	if err != nil {
		panic(err)
	}
	httpServer.Start(ctx)

}

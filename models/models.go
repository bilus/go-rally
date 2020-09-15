package models

import (
	log "github.com/sirupsen/logrus"

	"rally/adapter"

	r "github.com/go-redis/redis"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v5"
)

// DB is a connection to your database to be used
// throughout your application.
var DB *pop.Connection

var Redis adapter.Redis

func init() {
	var err error
	env := envy.Get("GO_ENV", "development")
	DB, err = pop.Connect(env)
	if err != nil {
		log.Fatal(err)
	}
	pop.Debug = env == "development"

	redisOptions, err := r.ParseURL(envy.Get("OPENREDIS_URL", "redis://localhost:6379"))
	if err != nil {
		log.Fatal(err)
	}
	rds := r.NewClient(redisOptions)
	Redis = adapter.NewRedis(rds)
}

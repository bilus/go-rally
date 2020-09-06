package models

import (
	log "github.com/sirupsen/logrus"

	"rally/redis"

	r "github.com/go-redis/redis"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v5"
)

// DB is a connection to your database to be used
// throughout your application.
var DB *pop.Connection

var Redis redis.Store

func init() {
	var err error
	env := envy.Get("GO_ENV", "development")
	DB, err = pop.Connect(env)
	if err != nil {
		log.Fatal(err)
	}
	pop.Debug = env == "development"

	// TODO
	rds := r.NewClient(&r.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	Redis = redis.NewStore(rds)
}

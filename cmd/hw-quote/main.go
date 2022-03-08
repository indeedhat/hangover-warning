package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/indeedhat/hangover-warning/internal/controllers"
	"github.com/indeedhat/hangover-warning/internal/env"
	"github.com/indeedhat/hangover-warning/internal/store"
)

func main() {
	db, err := store.Connect()
	if err != nil {
		panic(err)
	}

	if err = store.Migrate(db); err != nil {
		panic(err)
	}

	router := gin.Default()
	router.Use(sessions.Sessions("hwarning", sessions.NewCookieStore([]byte(env.Get(env.Secret)))))

	_ = controllers.NewQuotes(router, db)
	_ = controllers.NewViews(router, db)

	router.Run(env.GetFallback(env.BindAddress, ":8080"))
}

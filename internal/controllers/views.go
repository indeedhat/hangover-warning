package controllers

import (
	"net/http"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/indeedhat/hangover-warning/internal/env"
	"gorm.io/gorm"
)

type Views struct {
	db *gorm.DB
}

// NewQuotes controller
func NewViews(router *gin.Engine, db *gorm.DB) Views {
	views := Views{db}

	views.routes(router)

	return views
}

// routes sets up all the routes related to quotes
func (vew Views) routes(router *gin.Engine) {
	static := router.Group("/", gzip.Gzip(gzip.DefaultCompression))
	static.Static("/assets", "./web/assets")
	static.Static("/uploads", "./web/uploads")

	viewsConfig := goview.DefaultConfig
	viewsConfig.Root = "./web/views"
	viewsConfig.DisableCache = true

	router.HTMLRender = ginview.New(viewsConfig)
	public := router.Group("/")
	{
		public.GET("login", vew.Login())
		public.POST("login", vew.HandleLogin())
	}

	userLocal := router.Group("/", authRequired)
	{
		userLocal.GET("", vew.Index())
	}
}

// List teh quotes on a list
func (vew Views) Index() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index", nil)
	}
}

// Login page
func (vew Views) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login", nil)
	}
}

func (vew Views) HandleLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		password := ctx.PostForm("password")

		if password != env.Get(env.Secret) {
			ctx.Redirect(http.StatusFound, "/login")
			return
		}

		// Save the username in the session
		session.Set("loggedin", "true")
		if err := session.Save(); err != nil {
			ctx.Redirect(http.StatusFound, "/login")
			return
		}

		ctx.Redirect(http.StatusFound, "/")
	}
}

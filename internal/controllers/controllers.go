package controllers

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	errBadInput      = gin.H{"outcome": false, "message": "Bad input"}
	errSlugInUse     = gin.H{"outcome": false, "message": "Slug in use"}
	errRequestFailed = gin.H{"outcome": false, "message": "Request Failed"}
	errNotFound      = gin.H{"outcome": false, "message": "Not Found"}
	errNone          = gin.H{"outcome": true, "message": nil}
)

// authRequired is a simple middleware to check the session
func authRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("loggedin")
	if user == nil {
		// Abort the request with the appropriate error code
		c.Redirect(http.StatusFound, "/login")
		return
	}
	// Continue down the chain to handler etc
	c.Next()
}

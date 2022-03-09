package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/indeedhat/hangover-warning/internal/store"
	"gorm.io/gorm"
)

type Quotes struct {
	db *gorm.DB
}

// NewQuotes controller
func NewQuotes(router *gin.Engine, db *gorm.DB) Quotes {
	quotes := Quotes{db}

	quotes.routes(router)

	return quotes
}

// routes sets up all the routes related to quotes
func (quo Quotes) routes(router *gin.Engine) {
	lists := router.Group("/api/quotes", authRequired)
	{
		lists.GET("", quo.List())
		lists.POST("", quo.Create())
		lists.PUT("/:id", quo.Update())
		lists.DELETE("/:id", quo.Delete())
	}
}

// List teh quotes on a list
func (quo Quotes) List() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		quotes := store.ListQuotes(quo.db)

		ctx.JSON(http.StatusOK, gin.H{
			"outcome": true,
			"message": nil,
			"quotes":  quotes,
		})
	}
}

// Create a new quory on a liste
func (quo Quotes) Create() gin.HandlerFunc {
	type formInput struct {
		Text   string `form:"text" binding:"required"`
		Person string `form:"person" binding:"required"`
		Date   string `form:"date"`
	}

	return func(ctx *gin.Context) {
		var (
			input formInput
		)

		if err := ctx.ShouldBind(&input); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errBadInput)
			return
		}

		if list := store.CreateQuote(quo.db, input.Text, input.Person, &input.Date); list == nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errRequestFailed)
			return
		}

		ctx.JSON(http.StatusOK, errNone)
	}
}

// Update the text on an quory
func (quo Quotes) Update() gin.HandlerFunc {
	type formInput struct {
		Text   string `form:"text" binding:"required"`
		Person string `form:"person" binding:"required"`
		Date   string `form:"date"`
	}

	return func(ctx *gin.Context) {
		var (
			input formInput
			id    = ctx.Param("id")
		)

		if err := ctx.ShouldBind(&input); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errBadInput)
			return
		}

		existing := store.FindOuote(quo.db, id)
		if existing == nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, errNotFound)
			return
		}

		existing.Text = input.Text
		existing.Person = input.Person
		existing.Date = &input.Date
		if tx := quo.db.Save(existing); tx.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errRequestFailed)
			return
		}

		ctx.JSON(http.StatusOK, errNone)
	}
}

// Delete an quory
func (quo Quotes) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var id = ctx.Param("id")

		existing := store.FindOuote(quo.db, id)
		if existing == nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, errNotFound)
			return
		}

		if err := quo.db.Delete(existing).Error; err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errRequestFailed)
			return
		}

		ctx.JSON(http.StatusOK, errNone)
	}
}

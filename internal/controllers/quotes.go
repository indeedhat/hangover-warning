package controllers

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"

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

	var (
		errBadFile = gin.H{
			"outcome": false,
			"message": "File must be jpeg or png",
		}

		getMimeType = func(file *multipart.FileHeader) string {
			fh, err := file.Open()
			if err != nil {
				return ""
			}

			buff := make([]byte, 512)
			_, err = fh.Read(buff)
			if err != nil {
				return ""
			}

			return http.DetectContentType(buff)
		}
	)

	return func(ctx *gin.Context) {
		var input formInput

		if err := ctx.ShouldBind(&input); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errBadInput)
			return
		}

		file, err := ctx.FormFile("image")
		if err != nil {
			mimeType := getMimeType(file)
			if mimeType != "image/jpeg" && mimeType != "image/x-png" {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, errBadFile)
				return
			}
		}

		err = quo.db.Transaction(func(tx *gorm.DB) error {
			quote := store.CreateQuote(tx, input.Text, input.Person, &input.Date)
			if quote == nil {
				return errors.New("")
			}

			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			saveName := fmt.Sprintf("/uploads/%d%s", quote.ID, filepath.Ext(file.Filename))
			err = ctx.SaveUploadedFile(file, path.Join(cwd, "web", saveName))
			if err != nil {
				return err
			}

			quote.Image = &saveName
			return tx.Save(quote).Error
		})

		if err != nil {
			log.Print(err)
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

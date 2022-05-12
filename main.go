package main

import (
	"github.com/gin-gonic/gin"

	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Handler struct {
	db *gorm.DB
}

func newHandler(db *gorm.DB) *Handler {
	return &Handler{db}
}

func main() {
	// set mode for server
	gin.SetMode(gin.DebugMode)

	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Book{})

	handler := newHandler(db)
	r := gin.New()

	r.GET("/ping", pong)

	r.GET("/books", handler.listBookHandler)
	r.POST("/book", handler.addBookHandler)
	r.DELETE("/book/:id", handler.removeBookHandler)

	r.Run()
}

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func (h *Handler) listBookHandler(ctx *gin.Context) {
	var books []Book

	if result := h.db.Find(&books); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, books)
}

func (h *Handler) addBookHandler(ctx *gin.Context) {
	var book Book

	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if result := h.db.Create(&book); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, book)
}

func (h *Handler) removeBookHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	if result := h.db.Delete(&Book{}, id); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func pong(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"pong": "pong",
	})
}

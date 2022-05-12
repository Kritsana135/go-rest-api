package main

import (
	"github.com/gin-gonic/gin"

	"net/http"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET("/ping", pong)

	r.GET("/books", listBookHandler)
	r.POST("/book", addBookHandler)
	r.DELETE("/book/:id", removeBookHandler)

	r.Run()
}

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books = []Book{
	{ID: "1", Title: "เพราะเป็นวัยรุ่นจึงเจ็บปวด", Author: "คิมรันโด"},
	{ID: "2", Title: "แค่โอบกอดตัวเองให้เป็น", Author: "คิดมาก"},
	{ID: "3", Title: "นี่เราใช้ชีวิตยากเกินไปหรือเปล่านะ", Author: "ฮาวัน (Ha Wan)"},
}

func listBookHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, books)
}

func addBookHandler(ctx *gin.Context) {
	var book Book

	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	books = append(books, book)

	ctx.JSON(http.StatusCreated, book)
}

func removeBookHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	for i, a := range books {
		if a.ID == id {
			books = append(books[:i], books[i+1:]...)
			break
		}
	}

	ctx.Status(http.StatusNoContent)
}

func pong(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"pong": "pong",
	})
}

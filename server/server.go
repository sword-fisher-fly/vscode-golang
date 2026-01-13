package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}

var bookDB = map[string]Book{
	"1": {ID: "1", Title: "《Go语言实战》", Author: "William Kennedy", Price: 89.0},
	"2": {ID: "2", Title: "《Gin框架入门与实战》", Author: "张三", Price: 59.0},
	"3": {ID: "3", Title: "《REST API 设计指南》", Author: "李四", Price: 69.0},
}

func main() {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.GET("/books/:id", getBookByID)
	}

	err := r.Run(":8088")
	if err != nil {
		panic("Lauch http server failed: " + err.Error())
	}
}

func getBookByID(c *gin.Context) {
	bookID := c.Param("id")

	if bookID == "1" {
		// Mock handle delay
		time.Sleep(10 * time.Second)
	}

	book, exists := bookDB[bookID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Book not found",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hit success",
		"data":    book,
	})
}

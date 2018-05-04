package main

import (
	"fmt"
	"net/http"
	"os"
	// "github.com/gin-gonic/gin"
	// _ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/favicon.ico", NotFoundHandler)

	http.ListenAndServe(port, nil)

	// if port == "" {
	// 	log.Fatal("$PORT must be set")
	// }

	// router := gin.New()
	// router.Use(gin.Logger())
	// router.LoadHTMLGlob("templates/*.tmpl.html")
	// router.Static("/static", "static")

	// router.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.tmpl.html", nil)ush
	// })

	// router.Run(":" + port)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print(w, "Hello There")
}
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print(w, "Hello There")
}

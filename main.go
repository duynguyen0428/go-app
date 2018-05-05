package main

import (
	"encoding/json"
	"net/http"
	"os"
	// "github.com/gin-gonic/gin"
	// _ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/favicon.ico", FaviconHandler)

	http.ListenAndServe(":"+port, nil)

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
	// fmt.Print(w, "Hello There")
	user := User{"test@mail.com", "123456"}
	data, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(data)
}
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./favicon.ico")
}

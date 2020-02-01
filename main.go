package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	log.Println(GenerateKey("www"))
	r := gin.Default()
	r.GET("/:shorten_url", Redirect)
	r.POST("/delete", DealWithDeleteUrlRequest)
	r.POST("/add", DealWithAddUrlRequest)
	r.Run(":8088")
}

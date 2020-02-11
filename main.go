package main

import (
	"github.com/gin-gonic/gin"
	KGS "github.com/jscode017/go_key_generator_service"
	loadbalancer "github.com/jscode017/go_simple_load_balancer"
	//"log"
)

func main() {
	go KGS.GenerateKeysToRedis()
	r1 := gin.Default()
	r1.GET("/:shorten_url", Redirect)
	r1.POST("/delete", DealWithDeleteUrlRequest)
	r1.POST("/add", DealWithAddUrlRequest)
	go r1.Run(":8087")
	r2 := gin.Default()
	r2.GET("/:shorten_url", Redirect)
	r2.POST("/delete", DealWithDeleteUrlRequest)
	r2.POST("/add", DealWithAddUrlRequest)
	go r1.Run(":8088")

	lb := loadbalancer.NewLoadBalancer([]string{"localhost:8087", "localhost:8088"}, "localhost:8085", 10, 2)
	lb.Run()
}

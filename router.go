package main

import (
	"CloudPricingAPI/server"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/price", server.TestHandler)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

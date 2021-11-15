package main

import (
	"drone-bot/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routes.UseDronebotRouter(r)
	_ = r.Run(":8080")
}

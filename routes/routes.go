package routes

import (
	"drone-bot/routes/api"
	"github.com/gin-gonic/gin"
)

func UseDronebotRouter(r *gin.Engine) {
	repoapi := r.Group("api")
	{
		repoapi.PUT("repo", api.RepoPutHandler)
		//repoapi.GET("repo", api.RepoGetHandler)
	}
}

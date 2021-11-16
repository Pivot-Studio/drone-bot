package routes

import (
	"drone-bot/routes/api"
	"github.com/gin-gonic/gin"
)

func UseDronebotRouter(r *gin.Engine) {
	botapi := r.Group("/api")
	{
		botrepo := botapi.Group("/repo")
		{
			botrepo.PUT("", api.RepoPutHandler)
			botrepo.GET("/:ID", api.RepoGetHandler)
			botrepo.DELETE("/:ID", api.RepoDeleteHandler)
		}
		botapi.POST("/bot", api.PluginHandler)
	}
}

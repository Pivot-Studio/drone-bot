package api

import (
	"net/http"

	"github.com/Chronostasys/raft"
	"github.com/Chronostasys/raft/kvraft"
	"github.com/gin-gonic/gin"
)

type Message struct {
	Repourl  string `form:"repourl" binding:"required"`
	Bot_hook string `form:"bot_hook" binding:"required"`
}

func RepoPutHandler(ctx *gin.Context) {
	message := Message{}
	if err := ctx.ShouldBind(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ends := []string{"kv-0.kv-hs.kvrf.svc.cluster.local:8888", "kv-1.kv-hs.kvrf.svc.cluster.local:8888", "kv-2.kv-hs.kvrf.svc.cluster.local:8888"}
	rpcends := raft.MakeRPCEnds(ends)
	client := kvraft.MakeClerk(rpcends)

	if client.Get(message.Repourl) != "" {

	}
	client.Put(message.Repourl, message.Bot_hook)

	ctx.JSON(http.StatusOK, gin.H{
		"status":             "200",
		"putresult-repourl":  message.Repourl,
		"putresult-bot_hook": message.Bot_hook,
	})
	return
}

func RepoGetHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "200",
	})
	return
}

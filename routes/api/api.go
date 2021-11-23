package api

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/Chronostasys/raft"
	"github.com/Chronostasys/raft/kvraft"
	"github.com/gin-gonic/gin"
)

type Message struct {
	Repourl string `json:"repourl"`
	Bothook string `json:"bothook"`
}

type PluginMessage struct {
	Title   string `json:"title"`
	Repourl string `json:"repourl"`
	Author  string `json:"author"`
	Branch  string `json:"branch"`
	Message string `json:"message"`
	Githash string `json:"githash"`
}

var client *kvraft.Clerk

func init() {
	ends := []string{"kv-0.kv-hs.kvrf.svc.cluster.local:8888", "kv-1.kv-hs.kvrf.svc.cluster.local:8888", "kv-2.kv-hs.kvrf.svc.cluster.local:8888"}
	rpcends := raft.MakeRPCEnds(ends)
	client = kvraft.MakeClerk(rpcends)
}

func RepoPutHandler(ctx *gin.Context) {
	message := Message{}
	if err := ctx.ShouldBindJSON(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// check url valid
	if !(strings.Index(message.Repourl, "https://") == 0 &&
		strings.Index(message.Bothook, "https://") == 0) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid url"})
		return
	}

	client.Put(message.Repourl, message.Bothook)

	ctx.JSON(http.StatusOK, gin.H{
		"repourl": message.Repourl,
		"bothook": message.Bothook,
	})

}

func RepoGetHandler(ctx *gin.Context) {
	repo := ctx.Query("repo")

	ctx.JSON(http.StatusOK, gin.H{
		"repourl": repo,
		"bothook": client.Get(repo),
	})
}

func RepoDeleteHandler(ctx *gin.Context) {

	repourl := ctx.Query("repo")

	client.Put(repourl, "")

	ctx.JSON(http.StatusOK, gin.H{
		"repourl": repourl,
	})
}

func PluginHandler(ctx *gin.Context) {
	plugin_message := PluginMessage{}
	if err := ctx.ShouldBindJSON(&plugin_message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	bot_hook := client.Get(plugin_message.Repourl)

	if bot_hook == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "cannot find repourl-bot_hook"})
		return
	}

	if post_err := PostString2bot(
		plugin_message.Repourl, // FIXME too many arguments
		plugin_message.Message,
		bot_hook,
		plugin_message.Author,
		plugin_message.Branch,
		plugin_message.Githash,
		plugin_message.Title); post_err != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "send request to bot error:" + post_err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"repourl_applied": plugin_message.Repourl, // don't use "-" in name, use "_" instead. search snake case for details
		"bothook_get":     bot_hook,
		"message_send":    plugin_message.Message,
	})
}

func PostString2bot(repourl string, message string, bot_hook string, author string, branch string, githash string, title string) string {

	requestBody := fmt.Sprintf(`
		{
			"msg_type": "post",
			"content": {
				"post": {
					"zh_cn": {
						"title": "%s",
						"content": [
							[{
									"tag": "text",
									"text": "commit信息: %s "
							}],
							[{
									"tag": "text",
									"text": "触发者: %s "
							}],
							[{
									"tag": "text",
									"text": "分支: %s "
							}],
							[{
									"tag": "text",
									"text": "Githash: %s "
							}],
							[{
									"tag": "a",
									"text": "仓库链接",
									"href": "%s"
							}]
						]
					}
				}
			}
		}
	`, title, message, author, branch, githash, repourl)

	var jsonStr = []byte(requestBody)

	req, _ := http.NewRequest("POST", bot_hook, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return err.Error()
	}

	defer resp.Body.Close()
	return ""
}

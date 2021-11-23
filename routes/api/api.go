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

var (
	plugin_title   string
	plugin_repourl string
	plugin_author  string
	plugin_branch  string
	plugin_message string
	plugin_githash string
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
	plugin_message := Message{}
	if err := ctx.ShouldBindJSON(&plugin_message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"plugin_message": err.Error()})
		return
	}
	// check url valid
	if !(strings.Index(plugin_message.Repourl, "https://") == 0 &&
		strings.Index(plugin_message.Bothook, "https://") == 0) {
		ctx.JSON(http.StatusBadRequest, gin.H{"plugin_message": "invalid url"})
		return
	}

	client.Put(plugin_message.Repourl, plugin_message.Bothook)

	ctx.JSON(http.StatusOK, gin.H{
		"repourl": plugin_message.Repourl,
		"bothook": plugin_message.Bothook,
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
	plugin_messages := PluginMessage{}
	if err := ctx.ShouldBindJSON(&plugin_messages); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"plugin_message": err.Error()})
		return
	}

	bot_hook := client.Get(plugin_messages.Repourl)

	if bot_hook == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"plugin_message": "cannot find repourl-bot_hook"})
		return
	}
	plugin_message = plugin_messages.Message
	plugin_author = plugin_messages.Author
	plugin_branch = plugin_messages.Branch
	plugin_title = plugin_messages.Title
	plugin_githash = plugin_messages.Githash
	plugin_repourl = plugin_messages.Repourl

	if post_err := PostString2bot(bot_hook); post_err != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"plugin_message": "send request to bot error:" + post_err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"repourl_applied": plugin_messages.Repourl, // don't use "-" in name, use "_" instead. search snake case for details
		"bothook_get":     bot_hook,
		"message_send":    plugin_messages.Message,
	})
}

func PostString2bot(bot_hook string) string {

	requestBody := fmt.Sprintf(`
	{
		"config": {
		  "wide_screen_mode": true
		},
		"elements": [
		  {
			"tag": "div",
			"text": {
			  "content": "**commit信息**:%s\n**触发者**:%s\n**分支**:%s\n**Githash**:%s\n[仓库链接](%s)",
			  "tag": "lark_md"
			}
		  }
		],
		"header": {
		  "template": "turquoise",
		  "title": {
			"content": "%s",
			"tag": "plain_text"
		  }
		}
	  }
	`, plugin_message, plugin_author, plugin_branch, plugin_githash, plugin_repourl, plugin_title)

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

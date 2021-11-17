package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/Chronostasys/raft"
	"github.com/Chronostasys/raft/kvraft"
	"github.com/gin-gonic/gin"
)

type Message struct {
	Repourl  string `form:"repourl" binding:"required"`
	Bot_hook string `form:"bot_hook" binding:"required"`
}

type PluginMessage struct {
	Repourl string `form:"repourl" binding:"required"`
	Author  string `form:"author" binding:"required"`
	Branch  string `form:"branch" binding:"required"`
	Message string `form:"message" binding:"required"`
	Githash string `form:"githash" binding:"required"`
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

	id := Rand()
	for ; client.Get(strconv.Itoa(id)) != ""; id = Rand() {

	}
	client.Put(message.Repourl, message.Bot_hook)
	client.Put(strconv.Itoa(id), message.Repourl)

	ctx.JSON(http.StatusOK, gin.H{
		"id":       id,
		"repourl":  message.Repourl,
		"bot_hook": message.Bot_hook,
	})

	return
}

func RepoGetHandler(ctx *gin.Context) {
	id_string := ctx.Param("ID")

	ends := []string{"kv-0.kv-hs.kvrf.svc.cluster.local:8888", "kv-1.kv-hs.kvrf.svc.cluster.local:8888", "kv-2.kv-hs.kvrf.svc.cluster.local:8888"}
	rpcends := raft.MakeRPCEnds(ends)
	client := kvraft.MakeClerk(rpcends)

	id, _ := strconv.Atoi(id_string)
	ctx.JSON(http.StatusOK, gin.H{
		"id":      id,
		"repourl": client.Get(id_string),
		"bothook": client.Get(client.Get(id_string)),
	})
	return
}

func RepoDeleteHandler(ctx *gin.Context) {
	id := ctx.Param("ID")

	ends := []string{"kv-0.kv-hs.kvrf.svc.cluster.local:8888", "kv-1.kv-hs.kvrf.svc.cluster.local:8888", "kv-2.kv-hs.kvrf.svc.cluster.local:8888"}
	rpcends := raft.MakeRPCEnds(ends)
	client := kvraft.MakeClerk(rpcends)

	repourl := client.Get(id)
	bot_hook := client.Get(repourl)

	client.Put(repourl, "")
	client.Put(id, "")

	ctx.JSON(http.StatusOK, gin.H{
		"id":       id,
		"repourl":  repourl,
		"bot_hook": bot_hook,
	})
	return
}

func PluginHandler(ctx *gin.Context) {
	plugin_message := PluginMessage{}
	if err := ctx.ShouldBind(&plugin_message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ends := []string{"kv-0.kv-hs.kvrf.svc.cluster.local:8888", "kv-1.kv-hs.kvrf.svc.cluster.local:8888", "kv-2.kv-hs.kvrf.svc.cluster.local:8888"}
	rpcends := raft.MakeRPCEnds(ends)
	client := kvraft.MakeClerk(rpcends)

	bot_hook := client.Get(plugin_message.Repourl)

	if bot_hook == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "cannot find repourl-bot_hook"})
		return
	}

	if post_err := PostString2bot(plugin_message.Repourl, plugin_message.Message, bot_hook, plugin_message.Author, plugin_message.Branch, plugin_message.Githash); post_err != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "send request to bot error:" + post_err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"repourl-applied": plugin_message.Repourl,
		"bot_hook-get":    bot_hook,
		"message-send":    plugin_message.Message,
	})
	return
}

func PostString2bot(repourl string, message string, bot_hook string, author string, branch string, githash string) string {

	requestBody := fmt.Sprintf(`
		{
			"msg_type": "post",
			"content": {
				"post": {
					"zh_cn": {
						"title": "新的commit信息",
						"content": [
							[{
									"tag": "text",
									"text": "commit author: %s "
							}],
							[{
									"tag": "text",
									"text": "commit branch: %s "
							}],
							[{
									"tag": "text",
									"text": "commit githash: %s "
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
	`, author, branch, githash, repourl)

	var jsonStr = []byte(requestBody)

	req, _ := http.NewRequest("POST", bot_hook, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return err.Error()
	}

	defer resp.Body.Close()
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
	return ""
}

func Rand() int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(100)
}

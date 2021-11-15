package main

import (
	"bytes"
	"drone-bot/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func main() {
	r := gin.Default()
	routes.UseDronebotRouter(r)
	_ = r.Run(":8080")
}

func PostString(name string) {

	requestBody := fmt.Sprintf(`{
		"msg_type": "text",
    	"content": {
        	"text": "新更新提醒,%s"
		}
	}`, name)

	var jsonStr = []byte(requestBody)

	url := "https://open.feishu.cn/open-apis/bot/v2/hook/f65de4e6-d14d-4a75-9c2e-bc5c17c29da9"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

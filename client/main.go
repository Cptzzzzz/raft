package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"raft/global"
	"strings"
	"time"
)

const url = "http://127.0.0.1:80/client/operate"

func main() {
	var body []byte
	args := global.OperateArgs{
		Operator: 0,
		Key:      "11",
		Value:    "22",
	}
	body, err := json.Marshal(&args)
	if err != nil {
		panic(err)
	}
	reqBody := strings.NewReader(string(body))

	client := &http.Client{Timeout: time.Second * 5}
	request, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var reply global.OperateReply
	err = json.Unmarshal(body, &reply)
	if err != nil {
		panic(err)
	}
	fmt.Println(reply)
}

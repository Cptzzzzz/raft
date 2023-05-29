package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"log"
	"net/http"
	"os"
	"raft/global"
	"reflect"
	"strings"
	"time"
)

const url = "http://127.0.0.1:80/client/operate"

func main() {
	//tryProxy()
	fmt.Println(string(make([]byte, 0)))
}

func exampleRequest() {
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
func tryProxy() {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		os.Exit(1)
	}
	// setup a http client
	httpTransport := &http.Transport{
		Dial: dialer.Dial,
	}
	httpClient := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	if resp, err := httpClient.Post("http://192.168.103.137:8080/result", "", nil); err != nil {
		log.Fatalln(err)
	} else {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("%s\n", body)
	}

}
func getIp() {
	url := "192.168.103.137:8000"
	res := strings.Split(url, ":")
	fmt.Println(res[0])
}
func compare() {
	a := make([]global.Command, 10)
	b := make([]global.Command, 10)
	for i := 0; i < 10; i++ {
		a[i] = global.Command{
			Operator: i,
			Key:      fmt.Sprintf("%d", i),
			Value:    fmt.Sprintf("%d", i),
		}
		b[i] = global.Command{
			Operator: i,
			Key:      fmt.Sprintf("%d", i),
			Value:    fmt.Sprintf("%d", i),
		}
	}
	fmt.Println(reflect.DeepEqual(a, b))
}

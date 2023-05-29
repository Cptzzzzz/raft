package lib

import (
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"net/http"
	"os"
	"strings"
)

var HttpClient *http.Client

func InitNetwork() {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		os.Exit(1)
	}
	httpTransport := &http.Transport{
		Dial: dialer.Dial,
	}
	HttpClient = &http.Client{Transport: httpTransport}
}

func SendClientRequest(url string, data []byte) []byte {
	url = "http://" + url
	reqBody := strings.NewReader(string(data))
	request, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return nil
	}
	resp, err := HttpClient.Do(request)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	return body
}

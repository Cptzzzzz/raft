package raft

import (
	"io"
	"net/http"
	"raft/global"
	"strings"
)

func SendRequest(url string, data []byte) []byte {
	url = "http://" + url
	reqBody := strings.NewReader(string(data))
	request, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return nil
	}
	resp, err := global.HttpClient.Do(request)
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

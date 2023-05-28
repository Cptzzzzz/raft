package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"raft/global"
	"raft/network"
	"raft/raft"
	"time"
)

func main() {
	setup()
	global.Router = gin.Default()
	network.InitRoutes()
	err := global.Router.Run("0.0.0.0:80")
	if err != nil {
		return
	}
}

func setup() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	global.HttpClient = &http.Client{Timeout: time.Second * 5}
	global.Peers = viper.GetStringSlice("hosts")
	raft.DPrintf("length of peers: %d", len(global.Peers))
	global.JudgeHost = viper.GetString("judge")
	raft.DPrintf("me: %d", viper.GetInt("me"))
	global.Me = viper.GetInt("me")
	for index, host := range viper.GetStringSlice("hosts") {
		raft.DPrintf("node: %d, host: %s", index, host)
	}
	raft.Rf = raft.Make(viper.GetInt("me"))
}

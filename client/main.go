package main

import (
	"fmt"
	"raft/client/lib"
)

func main() {
	lib.ReadConfig()
	lib.InitNetwork()
	fmt.Println(lib.GetLeader())
}

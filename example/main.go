package main

import (
	"fmt"
	"raft/global"
	"reflect"
	"strings"
)

func main() {
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

package main

import (
	"fmt"
	"raft/global"
	"reflect"
)

func main() {
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

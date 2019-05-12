package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	var test []map[string]string
	test = append(test, map[string]string{"A": "AAA", "B": "BBB"})
	test = append(test, map[string]string{"GO": "GOGOHA", "HA": "HALOUTE"})
	var jsonlist []string
	for i := 0; i < len(test); i++ {
		j, _ := json.Marshal(test[i])
		jsonlist = append(jsonlist, string(j))
	}
	b, _ := json.Marshal(jsonlist)
	fmt.Println(string(b))
}

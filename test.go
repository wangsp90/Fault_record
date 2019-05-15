package main

import (
	//"encoding/json"
	"fmt"
	//"reflect"
)

type Recdata struct {
	Id         int
	Name       string
	Recordtext string
	Recorder   string
	Createtime string
	Filenames  []string
}

func main() {
	var test Recdata
	test.Recorder = "Ning"
	test.Name = "test1"
	test.Recordtext = "jdklfklsaj fldjds fdsajlf sadfj"
	test.Createtime = "2019-05-11 18:27:19"
	if test.Filenames == nil {
		fmt.Println(test.Createtime)
	}
}

package cfg

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Cfginfo struct {
	Http string
	Db   string
}

//done
func Getcfg() (cfginfo Cfginfo) {
	buf, errOpen := ioutil.ReadFile("cfg/cfg.json")
	if errOpen != nil {
		log.Println("配置文件加载失败...")
		return
	}
	errjson := json.Unmarshal(buf, &cfginfo)
	if errjson != nil {
		log.Println("error:", errjson)
		return
	}
	return
}

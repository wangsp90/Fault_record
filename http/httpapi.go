package main

import (
	"../db"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
)

type Cfginfo struct {
	Http string
	Db   string
}

//done
func Getcfg() (cfginfo Cfginfo) {
	buf, errOpen := ioutil.ReadFile("../cfg.json")
	if errOpen != nil {
		log.Println("配置文件加载失败...")
		return
	}

	errjson := json.Unmarshal(buf, &cfginfo)
	if errjson != nil {
		fmt.Println("error:", errjson)
		return
	}
	return
}

func main() {
	cfg := Getcfg()
	http.HandleFunc("/list", List)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/insert", Insertdata)
	http.HandleFunc("/search", Searchdata)
	http.HandleFunc("/deldata", Deldate)
	log.Fatal(http.ListenAndServe(cfg.Http, nil))
}

//done
func List(w http.ResponseWriter, r *http.Request) {
	Mydb := db.ConnectDatabase()
	defer Mydb.Close()
	RecordList := db.Getdata(Mydb)
	var jsonlist []string
	for i := 0; i < len(RecordList); i++ {
		j, _ := json.Marshal(RecordList[i])
		jsonlist = append(jsonlist, string(j))
	}
	data, _ := json.Marshal(jsonlist)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is update Fault Record."))
}

//done
func Deldate(w http.ResponseWriter, r *http.Request) {
	Mydb := db.ConnectDatabase()
	defer Mydb.Close()
	var rec map[string]string
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
	}
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	for k, v := range rec {
		if k == "id" {
			db.Deldata(Mydb, v)
		}
	}
}

func Searchdata(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is Search Fault Record."))
}

func Insertdata(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is Insert a new Fault Record."))
}

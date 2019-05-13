package httpserver

import (
	"encoding/json"
	// "fmt"
	_ "github.com/go-sql-driver/mysql"
	// "io/ioutil"
	"cfg"
	"log"
	"mydb"
	"net/http"
)

func Server(cfg cfg.Cfginfo) {
	http.HandleFunc("/list", List)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/insert", Insertdata)
	http.HandleFunc("/search", Searchdata)
	http.HandleFunc("/deldata", Deldate)
	log.Fatal(http.ListenAndServe(cfg.Http, nil))
}

//done
func List(w http.ResponseWriter, r *http.Request) {
	Mydb := mydb.ConnectDatabase()
	defer Mydb.Close()
	RecordList := mydb.Getdata(Mydb)
	var jsonlist []string
	for i := 0; i < len(RecordList); i++ {
		j, _ := json.Marshal(RecordList[i])
		jsonlist = append(jsonlist, string(j))
	}
	msg := make(map[string][]string)
	msg["msg"] = jsonlist
	data, _ := json.Marshal(msg)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is update Fault Record."))
}

//done
func Deldate(w http.ResponseWriter, r *http.Request) {
	Mydb := mydb.ConnectDatabase()
	defer Mydb.Close()
	var rec map[string]string
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
	}
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	var result string
	for k, v := range rec {
		if k == "id" {
			result = mydb.Deldata(Mydb, v)
		}
	}
	msg := make(map[string]string)
	msg["msg"] = result
	data, _ := json.Marshal(msg)
	w.Write(data)
}

func Searchdata(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is Search Fault Record."))
}

func Insertdata(w http.ResponseWriter, r *http.Request) {
	Mydb := mydb.ConnectDatabase()
	defer Mydb.Close()
	var rec map[string]string
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
	}
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	mydb.Insertdata(Mydb, rec)
	w.Write([]byte("This is Insert a new Fault Record."))
}

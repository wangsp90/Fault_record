package main

import (
	// "encoding/json"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type DBdetail struct {
	Id         int
	Theme      string
	Reason     string
	Recorder   string
	Lasttime   string
	Starttime  string
	Endtime    string
	Process    string
	Solution   string
	Appearance string
	Effect     string
	Filesname  string
}

func main() {
	http.HandleFunc("/list", List)
	http.HandleFunc("/search", Search)
	http.Handle("/", http.FileServer(http.Dir("./")))
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func List(w http.ResponseWriter, r *http.Request) {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	res, err := httpClient.Get("http://127.0.0.1:8000/api/v1/list")
	if err != nil {
		log.Printf("something wrong here: %s \n", err.Error())
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Println(err)
	}
	log.Println(string(body))
	w.Write(body)
}

func Search(w http.ResponseWriter, r *http.Request) {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	if r.Body == nil {
		w.Write([]byte("Nodata!"))
	}
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	log.Println(string(b))

	data := bytes.NewBuffer(b)
	res, err := httpClient.Post("http://127.0.0.1:8000/api/v1/search", "application/json;charset=utf-8", data)
	if err != nil {
		log.Printf("something wrong here: %s \n", err.Error())
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Println(err)
	}
	log.Println(string(body))
	w.Write(body)
}

package httpserver

import (
	"encoding/json"
	// "fmt"
	_ "github.com/go-sql-driver/mysql"
	// "io/ioutil"
	"cfg"
	"io"
	"log"
	"mydb"
	"net/http"
	"os"
)

func Uploadindex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(tpl))
}

const tpl = `<html>
<head>
<title>上传文件</title>
</head>
<body>
<form enctype="multipart/form-data" action="/api/v1/upload" method="post">
 <input type="file" name="uploadfile" multiple/>
 <input type="hidden" name="token" value="{...{.}...}"/>
 <input type="submit" value="上传" />
</form>
</body>
</html>`

func Server(cfg cfg.Cfginfo) {
	http.HandleFunc("/api/v1/list", Getlist)
	http.HandleFunc("/api/v1/update", Update)
	http.HandleFunc("/api/v1/insert", Insertdata)
	http.HandleFunc("/api/v1/search", Searchtheme)
	http.HandleFunc("/api/v1/deldata", Deldate)
	http.HandleFunc("/api/v1/getdetail", Getdetail)
	//文件服务器
	http.HandleFunc("/upload", Uploadindex)
	http.Handle("/api/v1/fileserver", http.StripPrefix("/api/v1/fileserver", http.FileServer(http.Dir("./files/"))))
	http.HandleFunc("/api/v1/upload", Upload)
	log.Fatal(http.ListenAndServe(cfg.Http, nil))
}

//done
func Getlist(w http.ResponseWriter, r *http.Request) {
	Mydb := mydb.ConnectDatabase()
	defer Mydb.Close()
	RecordList := mydb.Getlist(Mydb)
	b, err := json.Marshal(RecordList)
	if err != nil {
		log.Println(err)
	}
	msg := make(map[string]string)
	log.Println(string(b))
	msg["msg"] = string(b)
	data, _ := json.Marshal(msg)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func Insertdata(w http.ResponseWriter, r *http.Request) {
	Mydb := mydb.ConnectDatabase()
	defer Mydb.Close()
	var rec mydb.DBdetail
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
	}
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	b := mydb.Insertdata(Mydb, rec)
	if err != nil {
		log.Println(err)
	}
	msg := make(map[string]string)
	msg["msg"] = string(b)
	data, _ := json.Marshal(msg)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

//done
func Update(w http.ResponseWriter, r *http.Request) {
	Mydb := mydb.ConnectDatabase()
	defer Mydb.Close()
	var rec mydb.DBdetail
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
	}
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	log.Println(rec)
	var result string
	result = mydb.Updatedata(Mydb, rec)
	msg := make(map[string]string)
	msg["msg"] = result
	data, _ := json.Marshal(msg)
	w.Write(data)
	w.Write([]byte("This is update Fault Record."))
}

//done
func Deldate(w http.ResponseWriter, r *http.Request) {
	Mydb := mydb.ConnectDatabase()
	defer Mydb.Close()
	var rec mydb.DBdetail
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
	}
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	result := mydb.Deldata(Mydb, rec)
	msg := make(map[string]string)
	msg["msg"] = result
	data, _ := json.Marshal(msg)
	w.Write(data)
}

//done 只搜索name
func Searchtheme(w http.ResponseWriter, r *http.Request) {
	Mydb := mydb.ConnectDatabase()
	defer Mydb.Close()
	var rec mydb.DBdetail
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
	}
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	mydb.Searchtheme(Mydb, rec)
	//RecordList := mydb.Searchtheme(Mydb, rec)
	// var jsonlist []string
	// for i := 0; i < len(RecordList); i++ {
	// 	j, _ := json.Marshal(RecordList[i])
	// 	jsonlist = append(jsonlist, string(j))
	// }
	// msg := make(map[string][]string)
	// msg["msg"] = jsonlist
	// data, _ := json.Marshal(msg)
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// w.Write(data)
}

func Getdetail(w http.ResponseWriter, r *http.Request) {
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
	mydb.Getdetail(Mydb, rec)
	Detail := mydb.Getdetail(Mydb, rec)
	jsonlist, _ := json.Marshal(Detail)
	msg := make(map[string]string)
	msg["msg"] = string(jsonlist)
	data, _ := json.Marshal(msg)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	files := r.MultipartForm.File["uploadfile"]
	log.Println(files)
	var filenames []string
	for i := 0; i < len(files); i++ {
		f, err := os.Create("./files/" + files[i].Filename)
		filenames = append(filenames, files[i].Filename)
		defer f.Close()
		if err != nil {
			log.Println(err)
			return
		}
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			log.Println(err)
		}
		io.Copy(f, file)
	}
	msg := make(map[string][]string)
	msg["filenames"] = filenames
	data, _ := json.Marshal(msg)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

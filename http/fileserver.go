//处理上传单个文件

package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func filesayHelloHandler(w http.ResponseWriter, r *http.Request) {

	// 根据字段名获取表单文件
	formFile, header, err := r.FormFile("uploadfile")
	if err != nil {
		log.Printf("Get form file failed: %s\n", err)
		return
	}
	defer formFile.Close()
	// 创建保存文件
	destFile, err := os.Create("./upload/" + header.Filename)
	if err != nil {
		log.Printf("Create failed: %s\n", err)
		return
	}
	defer destFile.Close()

	// 读取表单文件，写入保存文件
	_, err = io.Copy(destFile, formFile)
	if err != nil {
		log.Printf("Write file failed: %s\n", err)
		return
	}

	//fmt.Fprintf(w, "Hello world!\n") //这个写入到w的是输出到客户端的
}

func mainfile() {
	http.HandleFunc("/", sayHelloHandler) //   设置访问路由
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//处理上传多个文件

func filessayHelloHandler(w http.ResponseWriter, r *http.Request) {
	//设置内存大小
	r.ParseMultipartForm(32 << 20)
	//获取上传的文件组
	files := r.MultipartForm.File["uploadfile"]
	len := len(files)
	for i := 0; i < len; i++ {
		//打开上传文件
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		//创建上传目录
		os.Mkdir("./upload", os.ModePerm)
		//创建上传文件
		cur, err := os.Create("./upload/" + files[i].Filename)

		defer cur.Close()
		if err != nil {
			log.Fatal(err)
		}
		io.Copy(cur, file)
		fmt.Println(files[i].Filename) //输出上传的文件名
	}

	//fmt.Fprintf(w, "Hello world!\n") //这个写入到w的是输出到客户端的
}

func mainfiles() {
	http.HandleFunc("/", sayHelloHandler) //	设置访问路由
	log.Fatal(http.ListenAndServe(":8080", nil))
}

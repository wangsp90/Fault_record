package db

//package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	//"os"
	"encoding/json"
	"strings"
	"time"
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

//done
func Getdata(db *sql.DB) (FaultRecord []map[string]string) {
	rows, errselect := db.Query("SELECT * from fault_record where isdel=0;")
	if errselect != nil {
		log.Fatalln(errselect)
	}
	defer rows.Close()
	//通过使用指定参数去获取数据库查询结果
	//通过使用column方法来获取数据库的查询结果，并传递给Map
	cols, errcols := rows.Columns()
	if errcols != nil {
		log.Fatalln(errcols)
	}
	//定义一个vals用于存放数据库查询的值，但是并非直接传递过来
	//数据库内容还是用Next()和Scan()获取
	vals := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range vals {
		scans[i] = &vals[i]
	}

	for rows.Next() {
		errrows := rows.Scan(scans...)
		if errrows != nil {
			log.Fatalln(errrows)
		}
		//fmt.Println(vals)
		row := make(map[string]string)
		for k, v := range vals {
			key := cols[k]
			row[key] = string(v)
		}
		FaultRecord = append(FaultRecord, row)
	}
	return
}

//done
func Insertdata(db *sql.DB, Rnewdata map[string]string) {
	//使用Prepare，可实现传递参数进行操作
	stmt, err := db.Prepare("INSERT INTO fault_record (name, recordtext, recorder, createtime) VALUES(?, ?, ?, ?);")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for k, v := range Rnewdata {
		switch {
		case k == "name":
			Rnewdata["name"] = v
		case k == "recordtext":
			Rnewdata["recordtext"] = v
		case k == "recorder":
			Rnewdata["recorder"] = v
		case k == "createtime":
			Rnewdata["createtime"] = v
		}
	}
	res, err := stmt.Exec(Rnewdata["name"], Rnewdata["recordtext"], Rnewdata["recorder"], Rnewdata["createtime"])
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID=%d, affected=%d\n", lastId, rowCnt)
}

//done
func Updatedata(db *sql.DB, Rnewdata map[string]string) {
	var execsql []string
	var values []string
	haed := "UPDATE fault_record SET"
	end := "WHERE id=?"
	for k, v := range Rnewdata {
		if k != "id" {
			values = append(values, (k + "='" + v + "'"))
		}
	}
	execsql = append(execsql, haed, strings.Join(values, ","))
	execsql = append(execsql, end)
	fmt.Println(strings.Join(execsql, " ") + ";")
	stmt, err := db.Prepare(strings.Join(execsql, " ") + ";")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	res, err := stmt.Exec([]byte(Rnewdata["id"]))
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID=%d, affected=%d\n", lastId, rowCnt)
}

//done
func Deldata(db *sql.DB, Rid interface{}) {
	stmt, err := db.Prepare("UPDATE fault_record SET isdel=1 WHERE id=?;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(Rid)
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID=%d, affected=%d\n", lastId, rowCnt)
}

func Searchdata(db *sql.DB, Rnewdata map[string]string) {

}

func ConnectDatabase() (db *sql.DB) {
	cfg := Getcfg()
	db, err := sql.Open("mysql", cfg.Db)
	if err != nil {
		panic(err)
	}
	return
}

func main() {
	time.Now()
	cfg := Getcfg()
	fmt.Println(cfg.Db)
	db := new(sql.DB)
	db, err := sql.Open("mysql", cfg.Db)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	//查询数据库
	// n := Getdata(db)
	// fmt.Println(n)
	//插入记录
	// testnewmap := map[string]string{"name": "故障台账录入测试", "recordtext": "This is insert DB test, Good Luck!", "recorder": "Duke", "createtime": "2019-05-12 12:17:28"}
	// Insertdata(db, testnewmap)
	//更新数据
	// testmap := map[string]string{"id": "7", "recordtext": "This is Update test, Good Luck!", "recorder": "Duke"}
	// Updatedata(db, testmap)
	//删除数据
	//Deldata(db, 8)

}

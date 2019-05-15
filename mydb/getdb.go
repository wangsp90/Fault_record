package mydb

//package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	//"os"
	"cfg"
	"strconv"
	"strings"
	"time"
)

//done
func Getlist(db *sql.DB) (FaultRecord []map[string]string) {
	rows, errselect := db.Query("SELECT id,name,recorder,createtime from fault_record where isdel=0;")
	if errselect != nil {
		log.Fatalln(errselect)
	}
	defer rows.Close()
	//数据库内容还是用Next()和Scan()获取
	//通过使用column方法来获取数据库的查询结果，并传递给Map
	cols, errcols := rows.Columns()
	if errcols != nil {
		log.Fatalln(errcols)
	}
	// //定义一个vals用于存放数据库查询的值，但是并非直接传递过来
	// //数据库内容还是用Next()和Scan()获取
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
	stmt1, err := db.Prepare("INSERT INTO fault_record (name, recordtext, recorder, createtime) VALUES(?, ?, ?, ?);")
	stmt2, err := db.Prepare("INSERT INTO files (recordid, filename) VALUES(?, ?);")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt1.Close()
	defer stmt2.Close()
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
		case k == "filenames":
			Rnewdata["filenames"] = v
		}
	}
	res1, err := stmt1.Exec(Rnewdata["name"], Rnewdata["recordtext"], Rnewdata["recorder"], Rnewdata["createtime"])
	if err != nil {
		log.Fatal(err)
	}
	recordlastId, err := res1.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	recordrowCnt, err := res1.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	res2, err := stmt2.Exec(recordlastId, Rnewdata["filenames"])
	if err != nil {
		log.Fatal(err)
	}
	filelastId, err := res2.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	filerowCnt, err := res2.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("记录保存成功ID=%d, affected=%d\n", recordlastId, recordrowCnt)
	fmt.Printf("文件保存成功ID=%d, affected=%d\n", filelastId, filerowCnt)
}

//done
func Updatedata(db *sql.DB, Rnewdata map[string]string) string {
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
	result := "ID=" + strconv.FormatInt(lastId, 10) + "," + "affected=" + strconv.FormatInt(rowCnt, 10)
	return result
}

//done
func Deldata(db *sql.DB, Rid interface{}) string {
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
	result := "ID=" + strconv.FormatInt(lastId, 10) + "," + "affected=" + strconv.FormatInt(rowCnt, 10)
	return result
}

func Searchdata(db *sql.DB, Rnewdata map[string]string) (FaultRecord []map[string]string) {
	for k, v := range Rnewdata {
		switch {
		case k == "starttime":
			Rnewdata["starttime"] = v
		case k == "endtime":
			Rnewdata["endtime"] = v
		case k == "keyword":
			Rnewdata["keyword"] = v
		}
	}
	sql1 := "SELECT id,name,recorder,createtime from fault_record where isdel=0 and "
	sql2 := "createtime>=\"" + Rnewdata["starttime"] + "\" and createtime<=\"" + Rnewdata["endtime"] + "\" and name like \"%" + Rnewdata["keyword"] + "%\";"
	searchsql := sql1 + sql2
	rows, errselect := db.Query(searchsql)
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

func ConnectDatabase() (db *sql.DB) {
	config := cfg.Getcfg()
	db, err := sql.Open("mysql", config.Db)
	if err != nil {
		panic(err)
	}
	return
}

func main() {
	time.Now()
	config := cfg.Getcfg()
	fmt.Println(config.Db)
	db := new(sql.DB)
	db, err := sql.Open("mysql", config.Db)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	//查询数据库
	// n := Getlist(db)
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

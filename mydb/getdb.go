package mydb

//package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	//"os"
	"cfg"
	"reflect"
	"strconv"
	"strings"
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

//done
func Getlist(db *sql.DB) []DBdetail {
	rows, errselect := db.Query("SELECT id,theme,recorder,starttime,effect from fault_record where isdel=0;")
	if errselect != nil {
		log.Fatalln(errselect)
	}
	defer rows.Close()
	var row DBdetail
	var FaultRecord []DBdetail
	for rows.Next() {
		rows.Scan(&row.Id, &row.Theme, &row.Recorder, &row.Starttime, &row.Effect)
		FaultRecord = append(FaultRecord, row)
	}
	return FaultRecord
}

//done
func Insertdata(db *sql.DB, Rnewdata DBdetail) string {
	stmt1, err := db.Prepare("INSERT INTO fault_record (theme, reason, recorder, lasttime,starttime,endtime,process,solution,appearance,effect,filesname) VALUES(?, ?, ?, ?,?,?,?,?,?,?,?);")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt1.Close()
	res1, err := stmt1.Exec(Rnewdata.Theme, Rnewdata.Reason, Rnewdata.Recorder, Rnewdata.Lasttime, Rnewdata.Starttime, Rnewdata.Endtime, Rnewdata.Process, Rnewdata.Solution, Rnewdata.Appearance, Rnewdata.Effect, Rnewdata.Filesname)
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
	fmt.Printf("记录保存成功ID=%d, affected=%d\n", recordlastId, recordrowCnt)
	b := "记录保存成功ID=" + strconv.FormatInt(recordlastId, 10) + ";" + "affected=" + strconv.FormatInt(recordrowCnt, 10) + "."
	return b
}

//done
func Updatedata(db *sql.DB, Rnewdata DBdetail) string {
	var execsql []string
	var values []string
	head := "UPDATE fault_record SET"
	end := "WHERE id=?"
	RnewType := reflect.TypeOf(Rnewdata)
	RewValues := reflect.ValueOf(Rnewdata)
	for i := 0; i < RewValues.NumField(); i++ {
		FieldName := RnewType.Field(i).Name
		FieldValue := RewValues.Field(i).String()
		if FieldName != "Id" {
			values = append(values, (strings.ToLower(FieldName) + "='" + FieldValue + "'"))
		}
	}
	execsql = append(execsql, head, strings.Join(values, ","), end)
	fmt.Println(strings.Join(execsql, " ") + ";")
	stmt1, err := db.Prepare(strings.Join(execsql, " ") + ";")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt1.Close()
	log.Println(strings.Join(execsql, " ") + ";")
	log.Println(Rnewdata.Id)
	res1, err := stmt1.Exec(Rnewdata.Id)
	recordlastId, err := res1.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	recordrowCnt, err := res1.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("记录保存成功ID=%d, affected=%d\n", recordlastId, recordrowCnt)
	result := "ID=" + strconv.FormatInt(recordlastId, 10) + "," + "affected=" + strconv.FormatInt(recordrowCnt, 10)
	return result
}

//done
func Deldata(db *sql.DB, Rnewdata DBdetail) string {
	stmt, err := db.Prepare("UPDATE fault_record SET isdel=1 WHERE id=?;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(Rnewdata.Id)
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

func Searchtheme(db *sql.DB, Rnewdata DBdetail) []DBdetail {
	head := "SELECT id,theme,recorder,starttime from fault_record where isdel=0 and "
	sqlline := []string{("starttime>=\"" + Rnewdata.Starttime + "\" and endtime<=\"" + Rnewdata.Endtime + "\"")}
	if Rnewdata.Theme != "" {
		sqlline = append(sqlline, ("theme like \"%" + Rnewdata.Theme + "%\""))
	}
	if Rnewdata.Effect != "" {
		sqlline = append(sqlline, ("effect like \"%" + Rnewdata.Effect + "%\""))
	}
	if Rnewdata.Recorder != "" {
		sqlline = append(sqlline, ("recorder like \"%" + Rnewdata.Recorder + "%\""))
	}
	full := head + (strings.Join(sqlline, " and ")) + ";"
	log.Println(full)
	rows, errselect := db.Query(head + (strings.Join(sqlline, " and ")) + ";")
	if errselect != nil {
		log.Fatalln(errselect)
	}
	defer rows.Close()
	var row DBdetail
	var FaultRecord []DBdetail
	for rows.Next() {
		rows.Scan(&row.Id, &row.Theme, &row.Recorder, &row.Starttime)
		FaultRecord = append(FaultRecord, row)
	}
	log.Println(FaultRecord)
	return FaultRecord
	// for k, v := range Rnewdata {
	// 	switch {
	// 	case k == "starttime":
	// 		Rnewdata["starttime"] = v
	// 	case k == "endtime":
	// 		Rnewdata["endtime"] = v
	// 	case k == "keyword":
	// 		Rnewdata["keyword"] = v
	// 	}
	// }
	// sql1 := "SELECT id,theme,recorder,starttime from fault_record where isdel=0 and "
	// sql2 := "starttime>=\"" + Rnewdata["starttime"] + "\" and starttime<=\"" + Rnewdata["endtime"] + "\" and theme like \"%" + Rnewdata["keyword"] + "%\";"
	// searchsql := sql1 + sql2
	// rows, errselect := db.Query(searchsql)
	// if errselect != nil {
	// 	log.Fatalln(errselect)
	// }
	// defer rows.Close()
	// //通过使用指定参数去获取数据库查询结果
	// //通过使用column方法来获取数据库的查询结果，并传递给Map
	// cols, errcols := rows.Columns()
	// if errcols != nil {
	// 	log.Fatalln(errcols)
	// }
	// //定义一个vals用于存放数据库查询的值，但是并非直接传递过来
	// //数据库内容还是用Next()和Scan()获取
	// vals := make([][]byte, len(cols))
	// scans := make([]interface{}, len(cols))
	// for i := range vals {
	// 	scans[i] = &vals[i]
	// }

	// for rows.Next() {
	// 	errrows := rows.Scan(scans...)
	// 	if errrows != nil {
	// 		log.Fatalln(errrows)
	// 	}
	// 	//fmt.Println(vals)
	// 	row := make(map[string]string)
	// 	for k, v := range vals {
	// 		key := cols[k]
	// 		row[key] = string(v)
	// 	}
	// 	FaultRecord = append(FaultRecord, row)
	// }
	// return
}

func Getdetail(db *sql.DB, Rnewdata map[string]string) map[string]string {
	execsql := []string{"SELECT a.id,a.name,a.recordtext,a.recorder,a.createtime, b.filename FROM fault_record a,upload_files b where a.id=b.recordid and b.recordid="}
	for k, v := range Rnewdata {
		if k == "id" {
			execsql = append(execsql, v, ";")
		}
	}
	rows, errselect := db.Query(strings.Join(execsql, ""))
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
	row := make(map[string]string)
	for rows.Next() {
		errrows := rows.Scan(scans...)
		if errrows != nil {
			log.Fatalln(errrows)
		}
		//fmt.Println(vals)
		for k, v := range vals {
			key := cols[k]
			row[key] = string(v)
		}
	}
	return row
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

package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

func file(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
  //文件根目录
	fileHost := "/mnt/thunder/file"
  //访问记录存入数据库
	db, _ := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/IpData")

	if len(request.Form) == 0 {
		showDir(fileHost, "/", db, writer, request)
	} else {
		file := ""
		for k, v := range request.Form {
			switch k {
			case "file":
				file = strings.Join(v, "")
				break
			}
		}
		if strings.Split(file, "/")[len(strings.Split(file, "/"))-1] == "" {
			showDir(fileHost, file, db, writer, request)
		} else {

			op, err := os.Open(fileHost + file)
			if err != nil {
				fmt.Fprintln(writer, err.Error())
			}
			defer op.Close()

			var number int
			var name string
			db.QueryRow("select * from fileIP where name = ?;", file).Scan(&number, &name)
			res, _ := db.Prepare("update fileIP set number = ? where name=?;")
			number += 1
			res.Exec(number, name)
			http.ServeContent(writer, request, file, time.Now(), op)
		}
	}
	db.Close()
}

func showDir(path string, dir string, db *sql.DB, writer http.ResponseWriter, request *http.Request) {
	dirs, err := ioutil.ReadDir(path + dir)

	var dirss []string
	var file []string
	var downNumber []string

	if err != nil {
		fmt.Fprint(writer, "dir null")
		return
	}
	for _, n := range dirs {
		number := 0
		var name string
		var fileName bytes.Buffer
		fileName.WriteString(dir)
		fileName.WriteString(n.Name())
		if n.IsDir() {
			dirss = append(dirss, n.Name())
		} else {
			err := db.QueryRow("select * from fileIP where name = ?;", fileName.String()).Scan(&number, &name)
			if err != nil {
				res, _ := db.Prepare("insert into fileIP(name,number)values (?,?)")
				res.Exec(fileName.String(), 0)

			}
			file = append(file, n.Name())
			downNumber = append(downNumber, strconv.Itoa(number))
		}
	}

	fmt.Fprint(writer,"<html><body>")
	for _, n := range dirss {
		fmt.Fprint(writer, "<img src=\"/images/folder.png\"><a href='/file?file="+dir+n+"/'>"+n+"<br/>")
	}
	for i, n := range file {
		fmt.Fprint(writer, "<img src=\"/images/document.png\"><a href='/file?file="+dir+n+"'>"+n+"</a>下载数次："+downNumber[i]+"<br/>")
	}
	fmt.Fprint(writer,"</body></html>")

	res, _ := db.Query("select name from fileIP")
	for res.Next() {
		var sqlName string
		res.Scan(&sqlName)
		_, err := os.Stat(path + sqlName)
		if os.IsNotExist(err) {
			db.Exec("delete from fileIP where name =?", sqlName)
		}
	}
}

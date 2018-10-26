package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var sqlHost = "root:rehtt946@tcp(127.0.0.1:3306)/IpData"	//数据库链接
var fileHost = "/mnt/A/"	//A盘地址
var fileHost_b = "/mnt/B/"	//B盘地址

func main(){
	
	var run []func()
	run = append(run, backup)	//定时任务要执行的函数
	go setTime(time.Hour*12, run)	//设置循环定时任务
	
	http.HandleFunc("/file", file)
	http.Handle("/Dfd4Fsgoeekcd9flsfkfsd/", http.StripPrefix("/file1/", http.FileServer(http.Dir(fileHost))))
	http.Handle("/Dfd4Fsgoeekcd2flsfkfsd/", http.StripPrefix("/file2/", http.FileServer(http.Dir(fileHost_b))))
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func file(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	db, _ := sql.Open("mysql", sqlHost)

	if len(request.Form) == 0 {
		showDir(fileHost, "/", db, writer, request)
		showDir(fileHost_b, "/", db, writer, request)
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

			//op, err := os.Open(fileHost + file)
			//if err != nil {
			//	fmt.Fprintln(writer, err.Error())
			//}
			//defer op.Close()
			filee := strings.Split(file, "//")
			var true error
			if filee[0] == "2" {
				_, true = os.Stat(fileHost_b + file)
			} else {
				_, true = os.Stat(fileHost + file)
			}

			if true != nil {
				fmt.Fprint(writer, 404)
				return
			}
			var number int
			var name string
			db.QueryRow("select * from fileIP where name = ?;", file).Scan(&number, &name)
			res, _ := db.Prepare("update fileIP set number = ? where name=?;")
			number += 1
			res.Exec(number, name)
			if filee[0] == "2" {
				http.Redirect(writer, request, "/file1/"+file, http.StatusFound)
			} else {
				http.Redirect(writer, request, "/file2/"+file, http.StatusFound)
			}

			//http.ServeContent(writer, request, file, time.Now(), op)
		}
	}
	db.Close()
}

func showDir(path string, dir string, db *sql.DB, writer http.ResponseWriter, request *http.Request) {
	dirs, err := ioutil.ReadDir(path + dir)

	file2 := ""
	if path == fileHost_b {
		file2 = "2/"
	}

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
			err := db.QueryRow("select * from fileIP where name = ?;", file2+fileName.String()).Scan(&number, &name)
			if err != nil {
				res, _ := db.Prepare("insert into fileIP(name,number)values (?,?)")
				res.Exec(file2+fileName.String(), 0)

			}
			file = append(file, n.Name())
			downNumber = append(downNumber, strconv.Itoa(number))
		}

	}

	fmt.Fprint(writer, "<html><body>")
	for _, n := range dirss {
		fmt.Fprint(writer, "<img src=\"/images/folder.png\"><a href='/file?file="+file2+dir+n+"/'>"+n+"<br/>")
	}
	for i, n := range file {
		fmt.Fprint(writer, "<img src=\"/images/document.png\"><a href='/file?file="+file2+dir+n+"'>"+n+"</a>下载数次："+downNumber[i]+"<br/>")
	}
	fmt.Fprint(writer, "</body></html>")

	res, _ := db.Query("select name from fileIP")
	for res.Next() {
		var sqlName string
		res.Scan(&sqlName)
		filee := strings.Split(sqlName, "//")
		if filee[0] == "2" {
			sqlName = ""
			sqlName = "/" + filee[1]
		}
		_, err := os.Stat(path + sqlName)
		if os.IsNotExist(err) {
			db.Exec("delete from fileIP where name =?", file2+sqlName)
		}
	}
}

//移动热度文件，将下载数量低于平均值的文件移动到B盘，高于平均值的文件移动到A盘
func backup() {
	var fileList []string
	var dirList []string
	filepath.Walk(fileHost+"/",
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				//fmt.Println("dir:", path)
				dirList = append(dirList, path)
				return nil
			}
			//fmt.Println("file:", path)
			fileList = append(fileList, path)
			return nil
		})

	db, _ := sql.Open("mysql", sqlHost)
	res, _ := db.Query("select name,number from fileIP")
	var file1 []string
	var file1_n []string
	file1_nn := 0
	var file1_m []int

	var file2 []string
	var file2_n []string
	file2_nn := 0
	var file2_m []int

	for res.Next() {
		var file string
		var number string
		res.Scan(&file, &number)
		name := strings.Split(file, "//")
		if name[0] == "2" {
			file2 = append(file2, name[1])
			file2_n = append(file2_n, number)
			i, _ := strconv.Atoi(number)
			file2_nn += i
		}
		file1 = append(file1, file)
		file1_n = append(file1_n, number)
		i, _ := strconv.Atoi(number)
		file1_nn += i
	}

	i := file1_nn / len(file1_n)
	for o, v := range file1_n {
		n, _ := strconv.Atoi(v)
		if n < i {
			file1_m = append(file1_m, o)
		}
	}
	i = file2_nn / len(file2_n)
	for o, v := range file2_n {
		n, _ := strconv.Atoi(v)
		if n < i {
			file2_m = append(file2_m, o)
		}
	}

	for _, v := range file1_m {
		moveFile(fileHost, file1[v], fileHost_b)
	}
	for _, v := range file2_m {
		moveFile(fileHost_b, file2[v], fileHost)
	}

}

//移动文件
func moveFile(f string, p string, t string) {

	_, err := os.Stat(t + p)
	if err != nil {
		if !os.IsExist(err) {
			o := strings.Split(p, "/")
			o[len(o)-1] = ""
			file := strings.Join(o, "/")
			os.MkdirAll(t+file, os.ModePerm)
		}
	}
	os.Rename(f+p, t+p)

}

//定时器
func setTime(t time.Duration, run []func()) {
	for {
		now := time.Now()
		// 计算下一个零点
		next := now.Add(t)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), next.Second(), 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		//执行函数
		for _, v := range run {
			v()
		}
	}
}

package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {

	path := flag.String("path", "./", "指定文件夹地址")
	port := flag.String("port", "8080", "端口")
	host := flag.String("host", "0.0.0.0", "监听地址")
	flag.Parse()
	http.Handle("/", http.FileServer(http.Dir(*path)))
	fmt.Println(*host + ":" + *port)
	err := http.ListenAndServe(*host+":"+*port, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

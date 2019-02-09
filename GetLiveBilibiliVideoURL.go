package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type json1 struct {
	PLAYURLRES struct {
		DATA struct {
			DURL []struct {
				URL string `json:"url"`
			} `json:"durl"`
		} `json:"data"`
	} `json:"playUrlres"`
}

func main() {
	fmt.Println(GetVideoUrl("输入网页直播地址"))
}

/*
输入网页直播地址
返回视频地址数组 
 */
func GetVideoUrl(url string) []string {
	h5:=strings.Split(url,"/")
	for i,v:=range h5{
		if v=="h5" {
			url="https://live.bilibili.com/"+h5[i+1]
			break
		}
	}
	h := make(map[string]string)
	h["Host"] = "live.bilibili.com"
	h["Accept"] = "text/html"
	h["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/18.17763"

	txt, err := Get(url,nil, h)
	if err != nil {
		fmt.Println(err)
	}
	defer txt.Body.Close()
	res, _ := ioutil.ReadAll(txt.Body)
	matched := regexp.MustCompile("<script>window(.*?)</script>")
	jso := matched.FindString(string(res))

	sp1 := strings.Split(jso, "NEPTUNE_IS_MY_WAIFU__=")[1]
	sp2 := strings.Split(sp1, "</")[0]

	var jsonn json1
	json.Unmarshal([]byte(sp2), &jsonn)
	var videoUrl []string
	for _,v:=range jsonn.PLAYURLRES.DATA.DURL {
		videoUrl=append(videoUrl, v.URL)
	}
	return videoUrl
}
func Get(url string, params map[string]string, headers map[string]string) (*http.Response, error) {
	//new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return nil, errors.New("new request is fail ")
	}
	//add params
	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}
	//add headers
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	//http client
	client := &http.Client{}
	log.Printf("Go GET URL : %s \n", req.URL.String())
	return client.Do(req)
}

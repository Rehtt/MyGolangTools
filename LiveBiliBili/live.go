package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var in_url = flag.String("u", "", "直播地址 获取视频地址")
var get_title = flag.String("t", "", "直播地址 获取视频标题")

type json1 struct {
	PLAYURLRES struct {
		DATA struct {
			DURL []struct {
				URL string `json:"url"`
			} `json:"durl"`
		} `json:"data"`
	} `json:"playUrlres"`
	BASEINFORES struct {
		DATA struct {
			TITLE     string `json:"title"`
			LIVE_TIME string `json:"live_time"`
		} `json:"data"`
	} `json:"baseInfoRes"`
}

func main() {
	flag.Parse()
	url := *in_url
	title := *get_title
	if url != "" {
		_, video := GetVideoUrl(url)
		//var res string
		//for _, v := range video {
		//	res = res + v + ",,,"
		//}
		if video != nil {
			fmt.Print(video[0])
			return
		}
	} else if title != "" {
		file_name, _ := GetVideoUrl(title)
		if file_name != "" {
			fmt.Println(file_name)
			return
		}
	}

	fmt.Println("error")

}

func GetVideoUrl(url string) (file_name string, video []string) {
	h5 := strings.Split(url, "/")
	for i, v := range h5 {
		if v == "h5" {
			url = "https://live.bilibili.com/" + h5[i+1]
			break
		}
	}
	h := make(map[string]string)
	h["Host"] = "live.bilibili.com"
	h["Accept"] = "text/html"
	h["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/18.17763"

	txt, err := Get(url, nil, h)
	if err != nil {
		//log.Println(err)
		return "", nil
	}
	defer txt.Body.Close()
	res, _ := ioutil.ReadAll(txt.Body)
	matched := regexp.MustCompile("<script>window.__NEPTUNE_IS_MY_WAIFU__=(.*?)</script>")
	if !matched.MatchString(string(res)) {
		//log.Println(err)

		return "", nil
	}
	jso := matched.FindString(string(res))
	sp1 := strings.Split(jso, "NEPTUNE_IS_MY_WAIFU__=")[1]
	sp2 := strings.Split(sp1, "</")[0]

	var jsonn json1
	json.Unmarshal([]byte(sp2), &jsonn)
	var videoUrl []string
	for _, v := range jsonn.PLAYURLRES.DATA.DURL {
		httpUrl := "http://" + strings.Split(v.URL, "://")[1]
		videoUrl = append(videoUrl, httpUrl)
	}
	//fmt.Println(jsonn.ROOMINITRES.DATA.ROOM_ID)
	return jsonn.BASEINFORES.DATA.TITLE + "_" + jsonn.BASEINFORES.DATA.LIVE_TIME, videoUrl
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
	//log.Printf("Go GET URL : %s \n", req.URL.String())
	return client.Do(req)
}

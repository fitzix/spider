package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/manifoldco/promptui"
)

const (
	URL = "https://connectkeyword.aliexpress.com/lenoIframeJson.htm"
)

var (
	keyMap       = make(map[string]int)
	resultKeyMap = make(map[string]int)
	writeChan    = make(chan keySearchResult)
)

type keySearchItem struct {
	Keywords  string `json:"keywords"`
	Count     string `json:"count"`
	CatName   string `json:"catName"`
	IsGeneral bool   `json:"isGeneral"`
}

type keySearchResult struct {
	KeyWordDTOs []keySearchItem `json:"keyWordDTOs"`
}

func main() {
	searchKeyPrompt := promptui.Prompt{Label: "输入查询关键字(按回车确定)",}
	startKey, err := searchKeyPrompt.Run()
	if err != nil {
		log.Fatalf("关键词有误: %s", err)
	}

	searchCountPrompt := promptui.Prompt{Label: "输入关键词结果数(如输入10000则查找结果集在10000以内的条目)", Validate: func(s string) error {
		if _, err := strconv.Atoi(s); err != nil {
			return errors.New("请输入数字")
		}
		return nil
	}}
	searchCountStr, err := searchCountPrompt.Run()

	if err != nil {
		log.Fatalf("关键词筛选数量有误: %s", err)
	}

	searchCount, _ := strconv.Atoi(searchCountStr)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(1)
	}()

	go writeToFile(searchCount)

	search(startKey, searchCount)
}

func search(keyword string, count int) {
	if _, ok := keyMap[keyword]; ok {
		return
	}

	log.Printf("start search keyword: %s", keyword)
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatalf("创建http client err: %s", err)
	}
	q := req.URL.Query()
	q.Add("keyword", keyword)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Referer", "https://www.aliexpress.com/")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36")
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("get data from url err: %s", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		log.Printf("read data from url err: %s", err)
		return
	}
	_ = resp.Body.Close()
	var ret keySearchResult

	if len(body) < 12 {
		log.Printf("empty search result")
		return
	}

	if err := json.Unmarshal(body[12:], &ret); err != nil {
		log.Printf("parse data err: %s", err)
		return
	}
	keyMap[keyword] = 1
	writeChan <- ret
	for _, v := range ret.KeyWordDTOs {
		if v.IsGeneral {
			search(v.Keywords, count)
		}
	}
}

func writeToFile(count int) {
	fileName := fmt.Sprintf("%d-aliexpress-keywords.csv", time.Now().Unix())
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("create file err: %s", err)
	}

	if _, err := file.WriteString("keywords,count\n"); err != nil {
		log.Fatalf("write header to csv file err: %s", err)
	}

	resultKeyMap := make(map[string]int)

	for {
		list := <-writeChan
		for _, v := range list.KeyWordDTOs {
			if !v.IsGeneral {
				continue
			}

			countInt, err := strconv.Atoi(strings.ReplaceAll(v.Count, ",", ""))
			if err != nil || countInt > count {
				continue
			}

			if c, ok := resultKeyMap[v.Keywords]; ok && c == countInt {
				continue
			}

			if _, err := file.WriteString(fmt.Sprintf("%s,%d\n", v.Keywords, countInt)); err != nil {
				log.Printf("write file err: %s", err)
			}
			resultKeyMap[v.Keywords] = countInt
		}
	}
}

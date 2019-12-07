package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fitzix/colly/browser"
	"github.com/gobuffalo/packr/v2"
	"github.com/gocolly/colly"
	"github.com/gorilla/websocket"
)

type amazonComment struct {
	Name        string `json:"name"`
	Star        string `json:"star"`
	ReviewTitle string `json:"reviewTitle"`
	ReviewLink  string `json:"reviewLink"`
	Date        string `json:"date"`
	Content     string `json:"content"`
	HelpfulVote string `json:"helpfulVote"`
	Sku         string `json:"SKU"`
	Asin        string `json:"ASIN"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {

	box := packr.New("public", "./public")

	html, _ := box.Find("index.html")

	http.HandleFunc("/parse", ws)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(html)
	})

	http.HandleFunc("/s", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		log.Println(r.FormValue("name"))
	})

	log.Println("listening on", ":7171")
	browser.Open("http://127.0.0.1:7171")
	if err := http.ListenAndServe(":7171", nil); err != nil {
		log.Fatal(err)
	}
}

func parse(url string, maxDep int, msg, csv chan amazonComment) {
	c := colly.NewCollector(colly.Async(true), colly.MaxDepth(maxDep))

	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "*amazon.*",
		Parallelism: 2,
		Delay:       2 * time.Second,
	})

	// count links
	c.OnHTML("#cm_cr-pagination_bar > ul > li.a-last > a", func(e *colly.HTMLElement) {
		_ = e.Request.Visit(e.Attr("href"))
	})

	c.OnHTML("#cm_cr-review_list", func(e *colly.HTMLElement) {
		e.ForEach("div[data-hook=review]", func(i int, e *colly.HTMLElement) {

			p := amazonComment{
				Name:        e.ChildText("div > div > div:nth-child(1) > a > div.a-profile-content > span"),
				ReviewTitle: e.ChildText("div > div > div:nth-child(2) > a.a-size-base.a-link-normal.review-title.a-color-base.review-title-content.a-text-bold > span"),
				ReviewLink:  e.Request.AbsoluteURL(e.ChildAttr("div > div > div:nth-child(2) > a.a-size-base.a-link-normal.review-title.a-color-base.review-title-content.a-text-bold", "href")),
				Content:     e.ChildText("div > div > div.a-row.a-spacing-small.review-data > span > span"),
			}

			starStr := e.ChildAttr("div > div > div:nth-child(2) > a:nth-child(1)", "title")
			if startArr := strings.Split(starStr, " "); len(startArr) > 0 {
				p.Star = startArr[0]
			}

			commentDateStr := e.ChildText("div > div > span.review-date")
			if commentDateArr := strings.Split(commentDateStr, " "); len(commentDateArr) > 0 {
				p.Date = commentDateArr[0]
			}

			helpfulStr := e.ChildText("div.a-row.a-spacing-none > div > div.a-row.review-comments > div > span.cr-vote > div.a-row.a-spacing-small > span")
			if helpfulArr := strings.Split(helpfulStr, " "); len(helpfulArr) > 0 {
				p.HelpfulVote = helpfulArr[0]
			}

			skuStr := e.ChildText("div > div > div.a-row.a-spacing-mini.review-data.review-format-strip > a")
			p.Sku = strings.TrimPrefix(skuStr, "Style: ")

			asinStr := e.ChildAttr("div > div > div.a-row.a-spacing-mini.review-data.review-format-strip > a", "href")
			if asinArr := strings.Split(asinStr, "/"); len(asinArr) > 4 {
				p.Asin = asinArr[3]
			}

			msg <- p
			csv <- p
		})
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
	})

	if err := c.Visit(url); err != nil {
		log.Println(err)
	}
	c.Wait()
}

func ws(w http.ResponseWriter, r *http.Request) {
	URL := r.URL.Query().Get("url")
	if URL == "" {
		log.Println("missing URL argument")
		return
	}
	maxDepth := 1
	maxDepthStr := r.URL.Query().Get("max")
	if maxDepthStr != "" {
		if max, err := strconv.Atoi(maxDepthStr); err == nil {
			maxDepth = max
		}
	}

	msgChan := make(chan amazonComment, 100)
	csvChan := make(chan amazonComment, 100)
	go parse(URL, maxDepth, msgChan, csvChan)
	go writeToCsv(csvChan)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		comment := <-msgChan
		comment.Content = ""
		comment.ReviewLink = ""
		resp, err := json.Marshal(comment)
		if err != nil {
			log.Println(err)
			continue
		}
		// Write message back to browser
		if err = conn.WriteMessage(websocket.TextMessage, resp); err != nil {
			return
		}
	}
}

func writeToCsv(msg chan amazonComment) {
	fileName := fmt.Sprintf("%d.csv", time.Now().Unix())
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("create file err: %s", err)
	}
	defer file.Close()
	csvWriter := csv.NewWriter(file)
	if err := csvWriter.Write([]string{"买家姓名", "标题链接", "评分", "评论时间", "SKU", "ASIN", "评论标题", "评论支持数", "评论内容"}); err != nil {
		log.Fatalf("write file err: %s", err)
	}

	ticker := time.NewTicker(time.Second * 30)

	for {
		select {
		case comment := <-msg:
			if err := csvWriter.Write([]string{comment.Name, comment.ReviewLink, comment.Star, comment.Date, comment.Sku, comment.Asin, comment.ReviewTitle, comment.HelpfulVote, comment.Content}); err != nil {
				log.Println("write csv file err: ", err)
			}
		case <-ticker.C:
			csvWriter.Flush()
		}
	}

}

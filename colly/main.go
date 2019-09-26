package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gocolly/colly"
)

func x() {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("#cm_cr-pagination_bar > ul > li.a-last > a", func(e *colly.HTMLElement) {
		_ = e.Request.Visit(e.Attr("href"))
	})

	c.OnHTML("#cm_cr-review_list", func(e *colly.HTMLElement) {
		e.ForEach("div[data-hook=review]", func(i int, e *colly.HTMLElement) {
			log.Printf("name ==== %s", e.ChildText("div > div > div:nth-child(1) > a > div.a-profile-content > span"))
			log.Printf("start ==== %s", e.ChildAttr("div > div > div:nth-child(2) > a:nth-child(1)", "title"))
			log.Printf("review title ==== %s", e.ChildText("div > div > div:nth-child(2) > a.a-size-base.a-link-normal.review-title.a-color-base.review-title-content.a-text-bold > span"))
			log.Printf("date ==== %s", e.ChildText("div > div > span.review-date"))
			log.Printf("content ==== %s", e.ChildText("div > div > div.a-row.a-spacing-small.review-data > span > span"))
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	_ = c.Visit("https://www.amazon.com/Joycuff-Inspirational-Bracelets-Encouragement-Personalized/product-reviews/B07CZFLS33/ref=cm_cr_getr_d_paging_btm_prev_1?ie=UTF8&reviewerType=all_reviews&pageNumber=1")
}

type pageInfo struct {
	Name string
	Star string
	ReviewTitle string
	Date string
	Content string
	StatusCode int
}

func handler(w http.ResponseWriter, r *http.Request) {
	URL := r.URL.Query().Get("url")
	if URL == "" {
		log.Println("missing URL argument")
		return
	}
	log.Println("visiting", URL)

	c := colly.NewCollector()

	p := &pageInfo{}

	// count links
	// c.OnHTML("#cm_cr-pagination_bar > ul > li.a-last > a", func(e *colly.HTMLElement) {
	// 	_ = e.Request.Visit(e.Attr("href"))
	// })

	c.OnHTML("#cm_cr-review_list", func(e *colly.HTMLElement) {
		e.ForEach("div[data-hook=review]", func(i int, e *colly.HTMLElement) {
			p.Name = e.ChildText("div > div > div:nth-child(1) > a > div.a-profile-content > span")
			p.Star = e.ChildAttr("div > div > div:nth-child(2) > a:nth-child(1)", "title")
			p.ReviewTitle = e.ChildText("div > div > div:nth-child(2) > a.a-size-base.a-link-normal.review-title.a-color-base.review-title-content.a-text-bold > span")
			p.Date = e.ChildText("div > div > span.review-date")
			p.Content = e.ChildText("div > div > div.a-row.a-spacing-small.review-data > span > span")
		})
	})

	// extract status code
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		p.StatusCode = r.StatusCode
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
		p.StatusCode = r.StatusCode
	})

	_ = c.Visit(URL)

	// dump results
	b, err := json.Marshal(p)
	if err != nil {
		log.Println("failed to serialize response:", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(b)
}

func main() {
	// example usage: curl -s 'http://127.0.0.1:7171/?url=http://go-colly.org/'
	addr := ":7171"

	http.HandleFunc("/", handler)

	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

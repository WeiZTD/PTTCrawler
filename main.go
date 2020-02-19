package main

import (
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
)

type articles struct {
	Url        string `gorm:"unique;not null"`
	Board      string `gorm:"PRIMARY_KEY"`
	Title      string
	Author     string
	Contains   string
	Reply      string
	Date       time.Time
	Updated_at time.Time
}

func main() {
	//database Init
	db, err := dbInit()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	for {
		getArticles(db)
		time.Sleep(5 * time.Minute)
	}
}

func getArticles(db *gorm.DB) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0"),
		colly.AllowedDomains("www.ptt.cc"),
		colly.Async(true),
	)
	//set cookie for Gossiping board. we're always older than 18 on the internet for sure.
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "over18=1")
	})
	c.Limit(&colly.LimitRule{
		Delay: 1 * time.Second,
	})

	//we don't want to crawl article that been deleted.
	c.OnHTML("div[class=r-ent]", func(e *colly.HTMLElement) {
		if !strings.Contains(e.ChildText("div.title"), "(本文已被刪除)") {
			//visit article url.
			url := "https://www.ptt.cc" + e.ChildAttr("div.title > a", "href")
			e.Request.Visit(url)
		}
	})

	c.OnHTML("div[id=main-content]", func(e *colly.HTMLElement) {
		article := articles{}
		article.Url = e.ChildAttr("span.f2 > a", "href")
		metaline := e.ChildTexts("div.article-metaline > span.article-meta-value")
		if len(metaline) > 2 && article.Url != "" {
			AllContains := strings.Split(e.Text, "\n")
			for i, v := range AllContains {
				if strings.Contains(v, "※ 文章網址: ") {
					//split article contains and reply by ※ 文章網址
					article.Contains = strings.Join(AllContains[1:i-2], "\n")
					article.Reply = strings.Join(AllContains[i+1:], "\n")
					break
				}
			}
			if article.Contains != "" {
				//check if article is already existed in database by url
				dumpcated, err := dbCheckExistArticle(article.Url, db)
				if err != nil {
					log.Fatal(err, article.Url)
				}
				if !dumpcated {
					//create new record
					article.Board = e.ChildText("div.article-metaline-right > span.article-meta-value")
					article.Author = metaline[0]
					article.Title = metaline[1]
					loc, _ := time.LoadLocation("Local")
					formatdate, err := time.ParseInLocation("Mon Jan 2 15:04:05 2006", metaline[2], loc)
					if err != nil {
						log.Fatal(err, article.Url)
					}
					article.Date = formatdate
					if err := dbCreateArticle(&article, db); err != nil {
						log.Fatal(err)
					}
				} else {
					//update contains, reply, updated_at columns
					if err := dbUpdateArticle(&article, db); err != nil {
						log.Fatal(err)
					}
				}
			}

		}
	})

	log.Println("updating articles")
	c.Visit("https://www.ptt.cc/bbs/Soft_Job/index.html")
	c.Visit("https://www.ptt.cc/bbs/Tech_Job/index.html")
	c.Visit("https://www.ptt.cc/bbs/TaichungBun/index.html")
	c.Visit("https://www.ptt.cc/bbs/Neihu/index.html")
	c.Visit("https://www.ptt.cc/bbs/Gossiping/index.html")
	c.Wait()
	log.Println("articles update completed. gopher will sleep for five minutes.")

}

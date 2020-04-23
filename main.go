package main

import (
	"flag"
	"github.com/MasterJoyHunan/fastmysql"
	"github.com/gocolly/colly"
	"golang.org/x/text/encoding/simplifiedchinese"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var tryAgain = make(map[string]int)
var lock sync.WaitGroup

func main() {
	flag.Parse()
	fastmysql.Setup()
	c := colly.NewCollector(
		colly.Async(true),
	)

	// 限制速度
	if err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       200 * time.Millisecond,
		RandomDelay: 5 * time.Second,
	}); err != nil {
		fastmysql.Logger.Panic("set colly limit error :", err)
	}

	fastmysql.Db.Exec("truncate region")

	// 随机Agent
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
	})

	// 发生错误
	c.OnError(func(rep *colly.Response, err error) {
		lock.Add(1)
		go ReVisit(rep.Request)
	})

	// 爬取省级
	c.OnHTML(".provincetr td", func(e *colly.HTMLElement) {
		text, _ := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(e.Text))
		e.Text = string(text)
		cityLink := e.ChildAttr("a", "href")

		region := Region{
			Id:       strings.Trim(cityLink, ".html"),
			ParentId: "0",
			Name:     string(text),
		}
		fastmysql.Db.Create(&region)
		e.Request.Visit(cityLink)
	})

	// 爬取市级
	c.OnHTML(".citytr", func(e *colly.HTMLElement) {
		code := e.DOM.Children().First().Text()
		gbkName := e.DOM.Children().Last().Text()
		link, _ := e.DOM.Children().First().Find("a").Attr("href")
		name, _ := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(gbkName))
		region := Region{
			Id:       code[:4],
			ParentId: code[:2],
			Name:     string(name),
		}
		fastmysql.Db.Create(&region)
		e.Request.Visit(link)
	})

	// 爬取县级
	c.OnHTML(".countytr", func(e *colly.HTMLElement) {
		code := e.DOM.Children().First().Text()
		gbkName := e.DOM.Children().Last().Text()
		link, _ := e.DOM.Children().First().Find("a").Attr("href")
		name, _ := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(gbkName))
		region := Region{
			Id:       code[:6],
			ParentId: code[:4],
			Name:     string(name),
		}
		fastmysql.Db.Create(&region)
		e.Request.Visit(link)
	})

	// 爬取乡镇级
	c.OnHTML(".towntr", func(e *colly.HTMLElement) {
		code := e.DOM.Children().First().Text()
		gbkName := e.DOM.Children().Last().Text()
		name, _ := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(gbkName))
		region := Region{
			Id:       code[:9],
			ParentId: code[:6],
			Name:     string(name),
		}
		fastmysql.Db.Create(&region)

	})

	// 开始爬取
	c.Visit("http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2019/index.html")
	//c.Visit("http://www.bb.vcc.qq")
	c.Wait()
	lock.Wait()
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func ReVisit(r *colly.Request) {
	defer lock.Done()
	url := r.URL.String()
	_, ok := tryAgain[url]
	if ok {
		//	if res >= 3 {
		//		fastmysql.Logger.Error(url, "重试3次后,无法连接")
		return
		//	}
	}
	tryAgain[url]++
	time.Sleep(6 * time.Minute)
	r.Visit(url)
}

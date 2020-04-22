package main

import (
	"flag"
	"github.com/MasterJoyHunan/fastmysql"
	"github.com/gocolly/colly"
	"golang.org/x/text/encoding/simplifiedchinese"
	"strings"
)

func main() {
	flag.Parse()
	fastmysql.Setup()
	c := colly.NewCollector(
		colly.Async(true),
	)

	//fastmysql.Db.Exec("truncate region")

	// 发生错误
	c.OnError(func(_ *colly.Response, err error) {
		fastmysql.Logger.Error(err)
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
			Id:       strings.Trim(code, "0"),
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
			Id:       strings.Trim(code, "0"),
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
			Id:       strings.Trim(code, "0"),
			ParentId: code[:6],
			Name:     string(name),
		}
		fastmysql.Db.Create(&region)

	})

	// 开始爬取
	//c.Visit("http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2019/index.html")
	c.Visit("http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2019/14/09/140923.html")

	c.Wait()
}

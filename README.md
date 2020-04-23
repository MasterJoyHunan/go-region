### 由GO语言编写的爬虫爬取国家统计局省市区编码,并存储MYSQL

* [国家统计局省市区统计地址](http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm)

注意：国家统计局官网有反爬虫系统，爬取速度不能过快，否则会出现请求重定向，导致爬取失败
所以添加如下参数：
```go
// 限制速度
if err := c.Limit(&colly.LimitRule{
    DomainGlob:  "*", 
    Parallelism: 2,
    Delay:       200 * time.Millisecond, 
    RandomDelay: 5 * time.Second,
}); err != nil {
    fastmysql.Logger.Panic("set colly limit error :", err)
}
```
这样导致的后果也很明显，没有错误，但是爬取速度也变的非常慢。好在这是一劳永逸的工作，爬取之后就完事了。
项目是爬取2019年的数据，如需修改，修改爬取地址（页面结构都是一样的，无需修改 OnHtml 回调）
```go
c.Visit("http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2019/index.html")
```
爬取结果不包含港澳台

tips: 在实战项目中，省市区数据应该存储在 redis 中最好
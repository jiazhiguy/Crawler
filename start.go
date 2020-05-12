package main

import(
	"fmt"
	"github.com/gocolly/colly/v2"
	"time"	
)
type Movie	struct{
	Name string
	Star string
	Evaluation string
	Actor string
	Quote string
	P_url string
	M_ulr string
}
func main() {
	t :=time.Now()
	movies := &[]Movie{}
	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(2),
	)
	detailCollector := c.Clone()
	detailCollector.Async=true

	c.OnHTML("ol.grid_view", func(e *colly.HTMLElement) {
		GetFromEl(e,"li",movies)
	})
	c.OnHTML("div.paginator",func(e *colly.HTMLElement){
		e.ForEach("a[href]",func(_ int, el *colly.HTMLElement){
			sub :=el.Attr("href")
			nextUrl:= e.Request.AbsoluteURL(sub)
			detailCollector.Visit(nextUrl)
		})
	})

	detailCollector.OnHTML("ol.grid_view", func(e *colly.HTMLElement) {
		GetFromEl(e,"li",movies)
	})
	detailCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://movie.douban.com/top250")
	c.Wait()
	detailCollector.Wait()

	
	fmt.Println("num of movies:",len(*movies))
	fmt.Println("cost time:",time.Since(t))
	fmt.Printf("%+v",(*movies)[1:25])

}
func GetFromEl(e *colly.HTMLElement,el string,storage *[]Movie){
	e.ForEach("li",func(_ int, el *colly.HTMLElement){
		movie :=Movie{}
		name :=el.ChildText("div.item>div.info>div.hd>a[href]>span")
		star :=el.ChildText("div.item>div.info>div.bd>div.star>span:nth-of-type(2)")
		evaluation :=el.ChildText("div.item>div.info>div.bd>div.star>span:nth-of-type(4)")
		actor :=el.ChildText("div.item>div.info>div.bd>p:nth-of-type(1)")
		quote :=el.ChildText("div.item>div.info>div.bd>p:nth-of-type(2)>span")
		m_ulr :=el.ChildAttr("div.item>div.info>div.hd>a[href]","href")
		movie.P_url=el.ChildAttr("div.item>div.pic>a[href]>img","src")
		movie.Actor=actor
		movie.Star=star
		movie.Evaluation=evaluation
		movie.Quote=quote
		movie.M_ulr=m_ulr
		movie.Name=name
		*storage = append (*storage,movie)
		// fmt.Printf("-----%+v----\n",star)
	})	
}
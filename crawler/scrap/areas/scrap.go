package areas

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/usk81/minimum-wage-ja/crawler/scrap"
	"github.com/usk81/minimum-wage-ja/domain"
	"golang.org/x/text/encoding/unicode"
)

const (
	uri = "https://www.mhlw.go.jp/stf/seisakunitsuite/bunya/koyou_roudou/roudoukijun/minimumichiran/"
)

var (
	encoding = unicode.UTF8
)

func Conf() scrap.Conf {
	return scrap.Conf{
		URL:      uri,
		Encoding: encoding,
		Func:     Scraper,
		Name:     "areas",
	}
}

func Scraper(doc *goquery.Document) (result map[string]scrap.Wages, err error) {
	result = map[string]scrap.Wages{}
	trs := doc.Find("#HID1 > div.section > div:nth-child(1) > div > div > table > tbody").
		First().
		Find("tr")

	trs.Each(func(i int, s *goquery.Selection) {
		tds := s.Find("td")
		prefecture := scrap.StripSpace(tds.Eq(0).Text())
		if !scrap.StringInSlice(prefecture, []string{"全国加重平均額", "都道府県名"}) {
			w, _ := scrap.Atoi(scrap.StripSpace(tds.Eq(1).Text()))
			d, _ := scrap.DateFromWareki(tds.Eq(3).Text())

			var wp scrap.Wages
			wp.SetRegional(domain.Wage{
				Daily:         0,
				Hourly:        w,
				ImplementedAt: d,
			})
			result[prefecture] = wp
		}
	})
	return
}

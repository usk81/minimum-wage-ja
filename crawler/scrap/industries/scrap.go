package industries

import (
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/usk81/minimum-wage-ja/crawler/scrap"
	"github.com/usk81/minimum-wage-ja/domain"
	"golang.org/x/text/encoding/japanese"
)

const (
	uri = "https://www.mhlw.go.jp/www2/topics/seido/kijunkyoku/minimum/minimum-19.htm"
)

var (
	encoding     = japanese.ShiftJIS
	regRegional  = regexp.MustCompile(`(\d+)\((.+)\)`)
	regSpace     = regexp.MustCompile(`\r\n|\r|\n|\s`)
	regDelimiter = regexp.MustCompile(`\,|，|､`)
)

func Conf() scrap.Conf {
	return scrap.Conf{
		URL:      uri,
		Encoding: encoding,
		Func:     Scraper,
		Name:     "industries",
	}
}

func Scraper(doc *goquery.Document) (result map[string]scrap.Wages, err error) {
	var prefecture string

	result = map[string]scrap.Wages{}
	trs := doc.Find("#contentsInner > div > div > div > div > div > table > tbody").
		First().
		Find("tr")

	trs.Each(func(i int, s *goquery.Selection) {
		tds := s.Find("td")
		if s.Children().Eq(0).Is("th") {
			prefecture = strings.TrimSpace(s.Find("th").First().Text())
			lw, ld, _ := getRegional(regSpace.ReplaceAllString(tds.Eq(0).Text(), ""))
			d, _ := scrap.DateFromWareki(tds.Eq(4).Text())
			wp := result[prefecture]
			wp.SetRegional(domain.Wage{
				Daily:         0,
				Hourly:        lw,
				ImplementedAt: ld,
			})
			name := scrap.Standardize(tds.Eq(1).Text())
			wp.SetIndustry(domain.Wage{
				Name:          name,
				Daily:         crop(tds.Eq(2).Text()),
				Hourly:        crop(tds.Eq(3).Text()),
				ImplementedAt: d,
			}, name)

			result[prefecture] = wp
		} else {
			d, _ := scrap.DateFromWareki(tds.Eq(3).Text())
			wp := result[prefecture]
			name := scrap.Standardize(tds.Eq(0).Text())
			wp.SetIndustry(domain.Wage{
				Name:          name,
				Daily:         crop(tds.Eq(1).Text()),
				Hourly:        crop(tds.Eq(2).Text()),
				ImplementedAt: d,
			}, name)
			result[prefecture] = wp
		}
	})
	return result, nil
}

func crop(s string) int {
	s = strings.Trim(s, " （※）")
	s = regDelimiter.ReplaceAllString(s, "")
	if d, err := scrap.Atoi(s); err == nil {
		return d
	}
	return 0
}

func getRegional(s string) (wage int, date time.Time, err error) {
	ss := regRegional.FindStringSubmatch(s)
	if wage, err = scrap.Atoi(ss[1]); err != nil {
		return
	}

	if date, err = scrap.DateFromWareki(ss[2]); err != nil {
		return
	}
	return
}

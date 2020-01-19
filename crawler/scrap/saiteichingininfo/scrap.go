package saiteichingininfo

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/usk81/minimum-wage-ja/crawler/scrap"
	"github.com/usk81/minimum-wage-ja/domain"
	"golang.org/x/text/encoding/unicode"
)

const (
	uri = "http://pc.saiteichingin.info/check/?p=%d"
)

var (
	jst      *time.Location
	encoding = unicode.UTF8
	regWage  = regexp.MustCompile(`[0-9]+`)
	regDate  = regexp.MustCompile(`[0-9]{4}/[0-9]{2}/[0-9]{2}`)
)

func init() {
	var err error
	if jst, err = time.LoadLocation("Asia/Tokyo"); err != nil {
		jst = time.FixedZone("Asia/Tokyo", 9*60*60)
	}
}

func Conf(pid int) scrap.Conf {
	return scrap.Conf{
		URL:      fmt.Sprintf(uri, pid),
		Encoding: encoding,
		Func:     Scraper,
		Name:     fmt.Sprintf("saiteichingin-%d", pid),
	}
}

func Scraper(doc *goquery.Document) (result map[string]scrap.Wages, err error) {
	result = map[string]scrap.Wages{}
	prefecture := doc.Find("h1.cont_tit02").First().Text()
	if prefecture == "" {
		err = errors.New("fail to get prefecture name")
		return
	}

	wp := result[prefecture]

	// area
	aw := doc.Find("ul.cfx").Find("li").Eq(1)
	wh, err := scrap.Atoi(aw.Find(".moneyL").First().Text())
	arrDate := regDate.FindStringSubmatch(aw.Text())
	var tt time.Time
	if len(arrDate) > 0 {
		tt, err = time.ParseInLocation("2006/01/02", arrDate[0], jst)
		if err != nil {
			return
		}
	}
	wp.SetRegional(domain.Wage{
		Daily:         0,
		Hourly:        wh,
		ImplementedAt: tt,
	})

	// industries
	idoc := doc.Find(".industryTableArea").First().Children()
	l := idoc.Length()
	if l == 0 {
		err = errors.New("Not found categorized wage")
		return
	}
	for i := 0; i < l; i = i + 2 {
		name := scrap.Standardize(idoc.Eq(i).Text())
		table := idoc.Eq(i + 1)
		rawDateStr := table.Find(".date").First().Text()

		if name == "" {
			err = errors.New("Can't get category name")
			return
		}

		var daily int
		if daily, err = getWage(table.Find(".day").First().Text()); err != nil {
			err = fmt.Errorf("parse daily wage. %w", err)
			return
		}

		var hourly int
		if hourly, err = getWage(table.Find(".time").First().Text()); err != nil {
			err = fmt.Errorf("parse hourly wage. %w", err)
			return
		}

		var ti time.Time
		arrDate = regDate.FindStringSubmatch(rawDateStr)
		if len(arrDate) > 0 {
			ti, err = time.ParseInLocation("2006/01/02", arrDate[0], jst)
			if err != nil {
				return
			}
		}
		wp.SetIndustry(domain.Wage{
			Name:          name,
			Daily:         daily,
			Hourly:        hourly,
			ImplementedAt: ti,
		}, name)

	}
	result[prefecture] = wp

	return
}

func getWage(s string) (w int, err error) {
	ss := regWage.FindStringSubmatch(s)
	if len(ss) > 0 && ss[0] != "" {
		return scrap.Atoi(ss[0])
	}
	return
}

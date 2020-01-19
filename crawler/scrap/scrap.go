package scrap

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/usk81/bone"
	"github.com/usk81/minimum-wage-ja/domain"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// Func is function that scraps a website to get minimum wages
type Func func(doc *goquery.Document) (result map[string]Wages, err error)

// Conf is config to scrap a website to get minimum wages
type Conf struct {
	Name     string
	URL      string
	Checksum string
	Encoding encoding.Encoding
	Func     Func
}

type Wages struct {
	Regional   domain.Wage
	Industries map[string]domain.Wage
}

// SetRegional sets 地域最低賃金
func (ws *Wages) SetRegional(r domain.Wage) {
	r.Name = "地域最低賃金"
	r.Regional = true
	ws.Regional = r
}

// SetIndustry sets 特定最低賃金
func (ws *Wages) SetIndustry(r domain.Wage, name string) {
	if ws.Industries == nil {
		ws.Industries = map[string]domain.Wage{}
	}
	r.Regional = false
	ws.Industries[name] = r
}

func Do(c *http.Client, conf Conf) (result map[string]Wages, skip bool, checksum string, err error) {
	resp, err := c.Get(conf.URL)
	if err != nil {
		return nil, false, "", err
	}
	defer resp.Body.Close()

	if err = bone.CheckResponse(resp); err != nil {
		return nil, false, "", err
	}

	var body io.Reader
	if conf.Encoding == unicode.UTF8 {
		body = resp.Body
	} else {
		// 情報取得先がShift-JISなのでUTF-8に変換かけてからReaderを取得
		body = transform.NewReader(resp.Body, conf.Encoding.NewDecoder())
	}

	bb, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}

	checksum, skip = DiffChecksum(bb, checksum)
	if skip {
		return
	}

	// ioutil.ReadAllの時点で消失してしまうので、[]byte から復元してgoqueryに渡す
	doc, err := goquery.NewDocumentFromReader(ioutil.NopCloser(bytes.NewBuffer(bb)))
	if err != nil {
		return
	}
	w, err := conf.Func(doc)
	return w, false, checksum, err
}

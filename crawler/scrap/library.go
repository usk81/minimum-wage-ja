package scrap

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/width"
)

var (
	jst          *time.Location
	regWareki    = regexp.MustCompile(`(.{2})(元|[0-9０-９]+)年([0-9０-９]{1,2})月([0-9０-９]{1,2})日`)
	regWareki2   = regexp.MustCompile(`([A-Z])(\d+)\.(\d+)\.(\d+)`)
	regDelimiter = regexp.MustCompile(`\,|，|､`)
)

func init() {
	var err error
	if jst, err = time.LoadLocation("Asia/Tokyo"); err != nil {
		jst = time.FixedZone("Asia/Tokyo", 9*60*60)
	}
}

// Atoi parses string and converts to type int
func Atoi(s string) (int, error) {
	s = strings.Replace(width.Narrow.String(s), ",", "", -1)
	return strconv.Atoi(s)
}

// DiffChecksum compares request body and a checksum string
func DiffChecksum(r []byte, target string) (cs string, result bool) {
	cs = Sha256Sum(r)
	return cs, cs == target
}

// Sha256Sum creates checksum by sha256
func Sha256Sum(data []byte) string {
	bytes := sha256.Sum256(data)
	return hex.EncodeToString(bytes[:])
}

// StripSpace strips 半角/全角 space from request string
//   先頭・末尾だけでなくすべてのスペースを除外
func StripSpace(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}

// DateFromWareki transforms 和暦 string to time.Time (西暦)
//   pattern1. 令和元年5月1日 -> 2019-05-01
//   pattern2. R1.05.01     -> 2019-05-01
func DateFromWareki(s string) (tt time.Time, err error) {
	var ss []string
	if regWareki.MatchString(s) {
		ss = regWareki.FindStringSubmatch(s)
	} else if regWareki2.MatchString(s) {
		ss = regWareki2.FindStringSubmatch(s)
	} else {
		err = fmt.Errorf("%s is invalid 年号 format", s)
		return
	}

	if len(ss) < 5 {
		err = fmt.Errorf("fail to parse date")
		return
	}

	var y, m, d int
	if y, err = YearFromGengo(ss[1], ss[2]); err != nil {
		return
	}
	if m, err = Atoi(ss[3]); err != nil {
		return
	}
	if d, err = Atoi(ss[4]); err != nil {
		return
	}

	tt = time.Date(y, time.Month(m), d, 0, 0, 0, 0, jst)
	return
}

// YearFromGengo transforms 元号 to year
func YearFromGengo(g, y string) (year int, err error) {
	if y == "元" {
		year = 1
	} else {
		if year, err = Atoi(y); err != nil {
			return
		}
	}

	switch g {
	case "H", "平成":
		year = 1988 + year
	case "R", "令和":
		year = 2018 + year
	default:
		return 0, fmt.Errorf("%s is invalid gengo", g)
	}
	return
}

// StringInSlice checks if a value exists in an slice
func StringInSlice(needle string, haystacks []string) bool {
	for _, hay := range haystacks {
		if hay == needle {
			return true
		}
	}
	return false
}

// Standardize fixes orthographical variants
func Standardize(s string) string {
	// 半角括弧を全角括弧に統一
	r := strings.Replace(strings.Replace(s, "(", "（", -1), ")", "）", -1)
	// 区切り文字を句読点に統一
	r = regDelimiter.ReplaceAllString(r, "、")

	// ゆらぎ、誤字、脱字の補完
	list := []struct {
		good  string
		wrong string
	}{
		// 半角中点を全角中点に統一
		{
			good:  "･",
			wrong: "・",
		},
		// 最低賃金ってどんな業種だよ！
		{
			good:  "",
			wrong: "最低賃金",
		},
		{
			good:  "附属品小売業",
			wrong: "付属品小売業",
		},
		{
			good:  "機械器具",
			wrong: "器械器具",
		},
		{
			good:  "その他の鉄鋼業",
			wrong: "その他鉄鋼業",
		},
		{
			good:  "金属製錬",
			wrong: "金属精錬",
		},
	}
	for _, l := range list {
		r = strings.Replace(r, l.wrong, l.good, -1)
	}
	return r
}

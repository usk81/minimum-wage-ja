package crawler

import (
	"net/http"
	"strings"
	"time"

	"github.com/usk81/minimum-wage-ja/domain"
	"github.com/usk81/minimum-wage-ja/interface/repository"

	"github.com/usk81/minimum-wage-ja/crawler/scrap"
	"github.com/usk81/minimum-wage-ja/crawler/scrap/areas"
	"github.com/usk81/minimum-wage-ja/crawler/scrap/industries"
	"github.com/usk81/minimum-wage-ja/crawler/scrap/saiteichingininfo"
)

func Run(conn repository.Connector) (err error) {
	c := &http.Client{
		Timeout: 30 * time.Second,
	}

	repo, err := conn.Connect(nil)
	if err != nil {
		return
	}

	if err = Do(c, repo, areas.Conf()); err != nil {
		return
	}

	if err = Do(c, repo, industries.Conf()); err != nil {
		return
	}

	for i := 0; i < 47; i++ {
		if err = Do(c, repo, saiteichingininfo.Conf(i)); err != nil {
			return
		}
	}

	return err
}

func Do(c *http.Client, repo repository.Wage, conf scrap.Conf) (err error) {
	rc := repo.Checksum()
	cs, err := rc.Get(conf.Name)
	if err != nil {
		return
	}
	conf.Checksum = cs

	result, skip, cs, err := scrap.Do(c, conf)
	if err != nil || skip {
		return
	}

	ps, err := repo.Prefecture().FindAll()
	if err != nil {
		return
	}

	// Update wages
	for k, v := range result {
		for _, p := range ps {
			if TrimPrefecture(p.Name) == TrimPrefecture(k) {
				v.Regional.PrefectureID = p.ID
				if err = update(repo, v.Regional); err != nil {
					return
				}

				for _, i := range v.Industries {
					i.PrefectureID = p.ID
					if err = update(repo, i); err != nil {
						return
					}
				}
			}
		}
	}

	// Update checksum
	if err = rc.Set(conf.Name, cs); err != nil {
		return
	}
	return nil
}

func update(repo repository.Wage, w domain.Wage) (err error) {
	if (domain.Wage{}) == w {
		return nil
	}
	r, err := repo.Find(w.Name, w.PrefectureID)
	if err != nil {
		return
	}
	if r == nil {
		// create
		w.ID = domain.ID()
		if err = repo.Set(w); err != nil {
			return
		}
	} else if w.ImplementedAt.After(r.ImplementedAt) {
		w.ID = r.ID
		if err = repo.Set(w); err != nil {
			return
		}
	}
	return
}

// TrimPrefecture trims 都府県
//   e.g. 東京都 -> 東京
func TrimPrefecture(s string) string {
	if s == "京都" || s == "京都府" {
		return "京都"
	}
	return strings.TrimRight(s, "都府県")
}

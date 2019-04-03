package yin

import (
	"net/http"
	"os"
	"path"
	"strings"
)

type ClientConfig struct {
	Directory             string
	BaseHref              string
	SinglePageApplication bool
}

func ServeClient(c ClientConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := Res(w, r)
		url := r.URL.Path
		if c.BaseHref != "" {
			url = strings.TrimPrefix(r.URL.Path, "/"+c.BaseHref)
		}
		serveThis := path.Join(c.Directory, url)

		file, _ := os.Stat(serveThis)
		if file != nil {
			res.SendFile(serveThis)
			return
		}

		if c.SinglePageApplication == false {
			file, _ = os.Stat(serveThis + ".html")
			if file != nil {
				res.SendFile(serveThis + ".html")
				return
			}

			file, _ = os.Stat(serveThis + "/index.html")
			if file != nil {
				res.SendFile(serveThis + "/index.html")
				return
			}
		}

		if c.SinglePageApplication == true {
			res.SendFile(c.Directory + "/index.html")
		} else {
			res.SendStatus(http.StatusNotFound)
		}
	}
}

package api

import (
	"fmt"
	"net/http"
	"net/url"
)

func (a *ApiServer) Index(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(r.URL.String())
	u.Path = fmt.Sprintf("/%s/", a.LatestVersion)
	http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
}

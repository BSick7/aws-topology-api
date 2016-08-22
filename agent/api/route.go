package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

type RouterItem interface {
	Register(router *mux.Router)
}

type Route struct {
	Method      string
	Path        string
	Name        string
	HandlerFunc http.HandlerFunc
}

func (r Route) Register(router *mux.Router) {
	router.Methods(r.Method).Path(r.Path).Name(r.Name).Handler(r.HandlerFunc)
}

type RouteGroup struct {
	Path      string
	Subroutes []RouterItem
}

func (r RouteGroup) Register(router *mux.Router) {
	sr := router.PathPrefix(r.Path).Subrouter()
	for _, ri := range r.Subroutes {
		ri.Register(sr)
	}
}

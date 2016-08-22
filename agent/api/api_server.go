package api

import (
	"github.com/BSick7/aws-topology-api/types"
	"github.com/gorilla/mux"
)

type ApiServer struct {
	Config        *types.AgentConfig
	Router        *mux.Router
	LatestVersion string
	AppVersion    string
}

func NewApiServer(c *types.AgentConfig, appVersion string) *ApiServer {
	return &ApiServer{
		Router:        mux.NewRouter().StrictSlash(false),
		Config:        c,
		LatestVersion: "v1",
		AppVersion:    appVersion,
	}
}

func (a *ApiServer) Register() {
	var items []RouterItem = []RouterItem{
		Route{"GET", "/", "index", a.Index},
		RouteGroup{"/v1", []RouterItem{
			Route{"GET", "/", "v1", a.Root("v1")},
		}},
	}

	for _, item := range items {
		item.Register(a.Router)
	}
}

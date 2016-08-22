package agent

import (
	"github.com/BSick7/aws-topology-api/agent/api"
	"github.com/BSick7/aws-topology-api/types"
	"net/http"
)

func Run(c *types.AgentConfig, appVersion string) error {
	s := api.NewApiServer(c, appVersion)
	s.Register()
	return http.ListenAndServe(c.Bind.String(), s.Router)
}

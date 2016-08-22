package agent

import (
	"github.com/BSick7/aws-topology-api/agent/api"
	"github.com/BSick7/aws-topology-api/types"
	"github.com/gorilla/handlers"
	"net/http"
	"os"
	"log"
)

func Run(c *types.AgentConfig, appVersion string) error {
	s := api.NewApiServer(c, appVersion)
	s.Register()
	log.Printf("Starting server %s...", c.Bind.String())
	return http.ListenAndServe(c.Bind.String(), handlers.LoggingHandler(os.Stdout, s.Router))
}

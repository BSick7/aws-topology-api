package api

import (
	"encoding/json"
	"net/http"
)

func (a *ApiServer) Root(apiVersion string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"name":        "AWS Topology API",
			"api_version": apiVersion,
			"app_version": a.AppVersion,
		})
	}
}

package api

import (
	"encoding/json"
	"github.com/BSick7/aws-topology-api/data"
	"github.com/BSick7/aws-topology-api/services"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-multierror"
	"log"
	"net/http"
)

func (a *ApiServer) Vpc(w http.ResponseWriter, r *http.Request) {
	out := json.NewEncoder(w)
	d := map[string]interface{}{}

	vars := mux.Vars(r)
	vpcId := vars["vpc-id"]

	broker := services.NewBroker(session.New())
	if err := broker.Init(); err != nil {
		log.Println(err)
		d["errors"] = toErrorStringSlice(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else if topo, err := data.GetVpcTopology(broker, vpcId); err != nil {
		log.Println(err)
		d["errors"] = toErrorStringSlice(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		d["vpc"] = topo.Vpc
		d["resources"] = topo.Resources
	}

	out.Encode(d)
}

func toErrorStringSlice(err error) []string {
	if err == nil {
		return []string{}
	}
	merr, ok := err.(*multierror.Error)
	if !ok {
		return []string{err.Error()}
	}
	errs := []string{}
	for _, werr := range merr.WrappedErrors() {
		errs = append(errs, werr.Error())
	}
	return errs
}

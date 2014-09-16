package api

import (
	"errors"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zenazn/goji/web"
	"github.com/samuel/go-zookeeper/zk"


	conf "github.com/QubitProducts/bamboo/configuration"
	service "github.com/QubitProducts/bamboo/services/service"
)

type ServiceAPI struct {
	Config    *conf.Configuration
	Zookeeper *zk.Conn
}

func (d *ServiceAPI) All(w http.ResponseWriter, r *http.Request) {
	services, err := service.All(d.Zookeeper, d.Config.DomainMapping.Zookeeper)

	if err != nil {
		fmt.Println(err)
		responseError(w, err.Error())
		return
	}

	responseJSON(w, services)
}

func (d *ServiceAPI) Create(w http.ResponseWriter, r *http.Request) {
	serviceModel, err := extractServiceModel(r)

	if err != nil {
		responseError(w, err.Error())
		return
	}

	_, err2 := service.Create(d.Zookeeper, d.Config.DomainMapping.Zookeeper, serviceModel.Id, serviceModel.Acl)
	if err2 != nil {
		responseError(w, "Marathon ID might already exist")
		return
	}

	responseJSON(w, serviceModel)
}

func (d *ServiceAPI) Put(c web.C, w http.ResponseWriter, r *http.Request) {
	identifier := c.URLParams["id"]
	serviceModel, err := extractServiceModel(r)
	if err != nil {
		responseError(w, err.Error())
		return
	}

	_, err1 := service.Put(d.Zookeeper, d.Config.DomainMapping.Zookeeper, identifier, serviceModel.Acl)
	if err1 != nil {
		responseError(w, err1.Error())
		return
	}

	responseJSON(w, serviceModel)
}


func (d *ServiceAPI) Delete(c web.C, w http.ResponseWriter, r *http.Request) {
	identifier := c.URLParams["id"]
	err := service.Delete(d.Zookeeper, d.Config.DomainMapping.Zookeeper, identifier)
	if err != nil {
		responseError(w, err.Error())
		return
	}

	responseJSON(w, new(map[string]string))
}


func extractServiceModel(r *http.Request) (service.Service, error) {
	var serviceModel service.Service
	payload := make([]byte, r.ContentLength)
	r.Body.Read(payload)
	defer r.Body.Close()

	err := json.Unmarshal(payload, &serviceModel)
	if err != nil {
		return serviceModel, errors.New("Unable to decode JSON request")
	}

	return serviceModel, nil
}


func responseError(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusBadRequest)
}

func responseJSON(w http.ResponseWriter, data interface {}) {
	w.Header().Set("Content-Type", "application/json")
	bites, _ := json.Marshal(data)
	w.Write(bites)
}
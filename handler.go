package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/Financial-Times/gourmet/log"
)

//go:embed "tpl/index.tpl"
var indexTemplate string

type HealthHandler struct {
	healthStore *HealthStore
	logger      *log.StructuredLogger
}

func NewHealthHandler(hs *HealthStore, l *log.StructuredLogger) *HealthHandler {
	return &HealthHandler{
		healthStore: hs,
		logger:      l,
	}
}

func (hh *HealthHandler) Status(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Header.Get("Accept") == "application/json":
		hh.StatusJSON(w, r)
	default:
		hh.StatusHTML(w, r)
	}
}

func (hh *HealthHandler) StatusJSON(w http.ResponseWriter, r *http.Request) {
	data := hh.healthStore.GetAll()

	categories := make(map[string][]HealthCheckStatus, 0)

	for _, i := range data {
		for _, t := range i.Service.Tags {
			_, exists := categories[t]
			if !exists {
				categories[t] = make([]HealthCheckStatus, 0)
			}
			categories[t] = append(categories[t], i)
		}
	}
	resp, err := json.Marshal(categories)
	if err != nil {
		hh.logger.Error("could not marshal json", log.WithError(err))
		fmt.Fprintf(w, "500")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(resp))
}

func (hh *HealthHandler) StatusHTML(w http.ResponseWriter, r *http.Request) {

	data := hh.healthStore.GetAll()

	categories := make(map[string][]HealthCheckStatus, 0)

	for _, i := range data {
		for _, t := range i.Service.Tags {
			_, exists := categories[t]
			if !exists {
				categories[t] = make([]HealthCheckStatus, 0)
			}
			categories[t] = append(categories[t], i)
		}
	}

	tmpl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		hh.logger.Error("could not parse template", log.WithError(err))
		fmt.Fprintf(w, "500")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, categories)
}

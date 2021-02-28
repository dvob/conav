package conav

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func NewCaseReportHanlder(cr CaseReporter) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/cases", CantonHandler(cr))
	r.HandleFunc("/cases/{municipality}", MunicipalityHandler(cr))
	return r
}

func CantonHandler(cr CaseReporter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		report, err := cr.GetCaseReport(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(report)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		return
	}
}

func MunicipalityHandler(cr CaseReporter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		municipality, _ := vars["municipality"]
		report, err := cr.GetCaseReportByMunicipality(r.Context(), municipality)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(report)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		return
	}
}

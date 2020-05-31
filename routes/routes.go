package routes

import (
	"encoding/json"
	"net/http"

	"github.com/emanpicar/currency-api/envelope"
	"github.com/emanpicar/currency-api/logger"
	"github.com/gorilla/mux"
)

type (
	Router interface {
		ServeHTTP(http.ResponseWriter, *http.Request)
	}

	routeHandler struct {
		envelopeManager envelope.Manager
		router          *mux.Router
	}
)

func NewRouter(envelopeManager envelope.Manager) Router {
	routeHandler := &routeHandler{envelopeManager: envelopeManager}

	return routeHandler.newRouter(mux.NewRouter())
}

func (rh *routeHandler) newRouter(router *mux.Router) *mux.Router {
	rh.registerRoutes(router)

	return router
}

func (rh *routeHandler) registerRoutes(router *mux.Router) {
	// router.HandleFunc("/api/authenticate", rh.authenticate).Methods("POST")
	router.HandleFunc("/rates/latest", rh.getLatestRates).Methods(http.MethodGet)
	router.HandleFunc("/rates/analyze", rh.getAnalyzedRates).Methods(http.MethodGet)
	router.HandleFunc("/rates/{cubeTime}", rh.getRatesByDate).Methods(http.MethodGet)

	rh.router = router
}

func (rh *routeHandler) getLatestRates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result, err := rh.envelopeManager.GetLatestRates()
	rh.encodeError(err, w)

	w.Write([]byte(result))
}

func (rh *routeHandler) getRatesByDate(w http.ResponseWriter, r *http.Request) {
	logger.Log.Infof("Getting rates by date: %v", mux.Vars(r)["cubeTime"])

	w.Header().Set("Content-Type", "application/json")
	result, err := rh.envelopeManager.GetRatesByDate(mux.Vars(r)["cubeTime"])
	rh.encodeError(err, w)

	w.Write([]byte(result))
}

func (rh *routeHandler) getAnalyzedRates(w http.ResponseWriter, r *http.Request) {
	logger.Log.Infoln("Getting all products")

	w.Header().Set("Content-Type", "application/json")
	result, err := rh.envelopeManager.GetAnalyzedRates()
	rh.encodeError(err, w)

	rh.encodeError(json.NewEncoder(w).Encode(result), w)
}

func (rh *routeHandler) encodeError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

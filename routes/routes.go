package routes

import (
	"encoding/json"
	"net/http"

	"github.com/emanpicar/currency-api/auth"
	"github.com/emanpicar/currency-api/entities/jsondata"
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
		authManager     auth.Manager
		router          *mux.Router
	}
)

func NewRouter(envelopeManager envelope.Manager, authManager auth.Manager) Router {
	routeHandler := &routeHandler{envelopeManager: envelopeManager, authManager: authManager}

	return routeHandler.newRouter(mux.NewRouter())
}

func (rh *routeHandler) newRouter(router *mux.Router) *mux.Router {
	rh.registerRoutes(router)

	return router
}

func (rh *routeHandler) registerRoutes(router *mux.Router) {
	router.HandleFunc("/api/auth", rh.authenticate).Methods(http.MethodPost).Name("Auth")
	router.HandleFunc("/rates/latest", rh.authMiddleware(rh.getLatestRates)).Methods(http.MethodGet).Name("RatesLatest")
	router.HandleFunc("/rates/analyze", rh.authMiddleware(rh.getAnalyzedRates)).Methods(http.MethodGet).Name("RatesAnalyze")
	router.HandleFunc("/rates/{cubeTime:[0-9]{4}-[0-9]{2}-[0-9]{2}}", rh.authMiddleware(rh.getRatesByDate)).Methods(http.MethodGet).Name("RatesByDate")

	rh.router = router
}

func (rh *routeHandler) authenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, err := rh.authManager.Authenticate(r.Body)
	if err != nil {
		rh.badRequest(err, w)
		return
	}

	rh.encodeError(json.NewEncoder(w).Encode(data), w)
}

func (rh *routeHandler) getLatestRates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result, err := rh.envelopeManager.GetLatestRates()
	if err != nil {
		rh.badRequest(err, w)
		return
	}

	w.Write([]byte(result))
}

func (rh *routeHandler) getRatesByDate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result, err := rh.envelopeManager.GetRatesByDate(mux.Vars(r)["cubeTime"])
	if err != nil {
		rh.badRequest(err, w)
		return
	}

	w.Write([]byte(result))
}

func (rh *routeHandler) getAnalyzedRates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result, err := rh.envelopeManager.GetAnalyzedRates()
	if err != nil {
		rh.badRequest(err, w)
		return
	}

	rh.encodeError(json.NewEncoder(w).Encode(result), w)
}

func (rh *routeHandler) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := rh.authManager.ValidateRequest(r)
		if err != nil {
			rh.badRequest(err, w)
			return
		}

		next(w, r)
	})
}

func (rh *routeHandler) encodeError(err error, w http.ResponseWriter) {
	if err != nil {
		logger.Log.Warnf("Error occurred: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (rh *routeHandler) badRequest(err error, w http.ResponseWriter) {
	logger.Log.Warnf("Error occurred: %v", err)
	w.WriteHeader(http.StatusBadRequest)
	rh.encodeError(json.NewEncoder(w).Encode(&jsondata.ResponseMessage{err.Error()}), w)
}

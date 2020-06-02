package routes

import (
	"testing"

	"github.com/gorilla/mux"
)

func Test_routeHandler_registerRoutes(t *testing.T) {
	dummyRouter := mux.NewRouter()
	type args struct {
		router *mux.Router
	}
	type routeDetails struct {
		routeName string
		routePath string
	}
	tests := []struct {
		name         string
		rh           *routeHandler
		args         args
		routeDetails routeDetails
	}{
		struct {
			name         string
			rh           *routeHandler
			args         args
			routeDetails routeDetails
		}{
			name:         "Validate Auth route",
			rh:           &routeHandler{},
			args:         args{router: dummyRouter},
			routeDetails: routeDetails{"Auth", "/api/auth"},
		},
		struct {
			name         string
			rh           *routeHandler
			args         args
			routeDetails routeDetails
		}{
			name:         "Validate RatesLatest route",
			rh:           &routeHandler{},
			args:         args{router: dummyRouter},
			routeDetails: routeDetails{"RatesLatest", "/rates/latest"},
		},
		struct {
			name         string
			rh           *routeHandler
			args         args
			routeDetails routeDetails
		}{
			name:         "Validate RatesAnalyze route",
			rh:           &routeHandler{},
			args:         args{router: dummyRouter},
			routeDetails: routeDetails{"RatesAnalyze", "/rates/analyze"},
		},
		struct {
			name         string
			rh           *routeHandler
			args         args
			routeDetails routeDetails
		}{
			name:         "Validate RatesByDate route",
			rh:           &routeHandler{},
			args:         args{router: dummyRouter},
			routeDetails: routeDetails{"RatesByDate", "/rates/{cubeTime:[0-9]{4}-[0-9]{2}-[0-9]{2}}"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.rh.registerRoutes(tt.args.router)
			router := dummyRouter.GetRoute(tt.routeDetails.routeName)

			if path, err := router.GetPathTemplate(); err != nil || path != tt.routeDetails.routePath {
				t.Errorf("Invalid declared route %v != %v", path, tt.routeDetails.routePath)
			}
		})
	}
}

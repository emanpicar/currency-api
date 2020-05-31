package main

import (
	"fmt"
	"net/http"

	"github.com/emanpicar/currency-api/db"
	"github.com/emanpicar/currency-api/envelope"
	"github.com/emanpicar/currency-api/logger"
	"github.com/emanpicar/currency-api/routes"
	"github.com/emanpicar/currency-api/settings"
)

func main() {
	logger.Init(settings.GetLogLevel())
	logger.Log.Infoln("Initializing Currency API")

	dbManager := db.NewManager()
	envelopeManager := envelope.NewManager(dbManager)

	envelopeManager.UpsertInitialData()
	// envelopeManager.GetLatestRates()
	// envelopeManager.GetAnalyzedRates()

	logger.Log.Fatal(http.ListenAndServeTLS(
		fmt.Sprintf("%v:%v", settings.GetServerHost(), settings.GetServerPort()),
		settings.GetServerPublicKey(),
		settings.GetServerPrivateKey(),
		routes.NewRouter(envelopeManager),
	))
}

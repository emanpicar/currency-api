package db

import (
	"fmt"

	"github.com/emanpicar/currency-api/entities/dbdata"
	"github.com/emanpicar/currency-api/logger"
	"github.com/emanpicar/currency-api/settings"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type (
	Manager interface {
		BatchFirstOrCreate(dbEnvelopeList *[]dbdata.Envelope)
		GetLatestRates() (*dbdata.Envelope, error)
		GetRatesByDate(cubeTime string) (*dbdata.Envelope, error)
		GetAnalyzedRates() (*dbdata.QuantitativeExchangeRate, error)
	}

	dbHandler struct {
		database *gorm.DB
	}
)

func NewManager() Manager {
	dbHandler := &dbHandler{}
	dbHandler.connect(gorm.Open)
	dbHandler.migrateTables()

	return dbHandler
}

func (dbHandler *dbHandler) connect(openConnection func(dialect string, args ...interface{}) (db *gorm.DB, err error)) {
	logger.Log.Infoln("Establishing connection to DB")

	var err error
	dbHandler.database, err = openConnection("postgres", fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=disable",
		settings.GetDBHost(), settings.GetDBPort(), settings.GetDBUser(), settings.GetDBName(), settings.GetDBPass(),
	))

	if err != nil {
		logger.Log.Fatalln(err)
	}

	logger.Log.Infoln("Successfully connected to DB")
}

func (dbHandler *dbHandler) migrateTables() {
	dbHandler.database.AutoMigrate(&dbdata.Envelope{})
	dbHandler.database.AutoMigrate(&dbdata.Cube{}).AddForeignKey("envelope_id", "envelopes(id)", "CASCADE", "CASCADE")
}

func (dbHandler *dbHandler) BatchFirstOrCreate(dbEnvelopeList *[]dbdata.Envelope) {
	for _, envelope := range *dbEnvelopeList {
		dbHandler.database.FirstOrCreate(&envelope, dbdata.Envelope{CubeTime: envelope.CubeTime})
	}
}

func (dbHandler *dbHandler) GetLatestRates() (*dbdata.Envelope, error) {
	env := &dbdata.Envelope{}
	err := dbHandler.database.Set("gorm:auto_preload", true).Order("cube_time desc").First(env).Error
	if err != nil {
		return nil, err
	}

	return env, nil
}

func (dbHandler *dbHandler) GetRatesByDate(cubeTime string) (*dbdata.Envelope, error) {
	env := &dbdata.Envelope{}
	err := dbHandler.database.Set("gorm:auto_preload", true).Where(&dbdata.Envelope{CubeTime: cubeTime}).First(env).Error
	if err != nil {
		return nil, err
	}

	return env, nil
}

func (dbHandler *dbHandler) GetAnalyzedRates() (*dbdata.QuantitativeExchangeRate, error) {
	result := &dbdata.QuantitativeExchangeRate{RatesAnalyze: []dbdata.RatesAnalyze{}}

	env := &dbdata.Envelope{}
	err := dbHandler.database.Set("gorm:auto_preload", false).First(env).Error
	if err != nil {
		return nil, err
	}

	result.Base = env.SenderName

	err = dbHandler.database.Raw("SELECT currency, min(rate) AS min, max(rate) AS max, avg(rate) AS avg FROM cubes GROUP BY currency").Scan(&result.RatesAnalyze).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

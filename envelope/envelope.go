package envelope

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/emanpicar/currency-api/db"
	"github.com/emanpicar/currency-api/entities/dbdata"
	"github.com/emanpicar/currency-api/entities/jsondata"
	"github.com/emanpicar/currency-api/entities/xmldata"
	"github.com/emanpicar/currency-api/logger"
	"github.com/emanpicar/currency-api/settings"
)

type (
	Manager interface {
		UpsertInitialData()
		GetLatestRates() (string, error)
		GetRatesByDate(cubeTime string) (string, error)
		GetAnalyzedRates() (*jsondata.QuantitativeExchangeRate, error)
	}

	Envelope struct {
		dbManager db.Manager
	}
)

func NewManager(dbManager db.Manager) Manager {
	return &Envelope{dbManager}
}

func (e *Envelope) UpsertInitialData() {
	logger.Log.Infoln("Upserting initial data started")
	env, err := e.downloadXMLData()
	if err != nil {
		logger.Log.Warnf("Unable to download xml data %v", err)
		env = e.useDemoData()
	}

	dbEnvelopeList := e.convertXMLtoDBEntities(env)
	e.dbManager.BatchFirstOrCreate(&dbEnvelopeList)

	logger.Log.Infoln("Upserting initial data completed")
}

func (e *Envelope) GetLatestRates() (string, error) {
	logger.Log.Infoln("Request on getting latest rates started")

	envelope, err := e.dbManager.GetLatestRates()
	if err != nil {
		return "", err
	}

	jsonResult := e.sortRatesThenToString(envelope)
	logger.Log.Infof("Latest rates data available, sender name: %v", envelope.SenderName)

	return jsonResult, nil
}

func (e *Envelope) GetRatesByDate(cubeTime string) (string, error) {
	logger.Log.Infof("Request on getting rates by date: %v started", cubeTime)

	envelope, err := e.dbManager.GetRatesByDate(cubeTime)
	if err != nil {
		return "", err
	}

	jsonResult := e.sortRatesThenToString(envelope)
	logger.Log.Infof("%v - rates data available, sender name: %v", cubeTime, envelope.SenderName)

	return jsonResult, nil
}

func (e *Envelope) GetAnalyzedRates() (*jsondata.QuantitativeExchangeRate, error) {
	logger.Log.Infoln("Request on getting analyzed rates started")

	analyzedResult, err := e.dbManager.GetAnalyzedRates()
	if err != nil {
		return nil, err
	}

	jsonResult := &jsondata.QuantitativeExchangeRate{
		Base:         analyzedResult.Base,
		RatesAnalyze: make(map[string]jsondata.RatesAnalyze),
	}

	for _, cube := range analyzedResult.RatesAnalyze {
		jsonResult.RatesAnalyze[cube.Currency] = jsondata.RatesAnalyze{
			Min: cube.Min,
			Max: cube.Max,
			Avg: cube.Avg,
		}
	}
	logger.Log.Infof("Analyzed rates data available, sender name: %v", analyzedResult.Base)

	return jsonResult, nil
}

// Json objects won't maintain order, to preserve order use Arrays
// Solution is to build the string data manually
func (e *Envelope) sortRatesThenToString(envelope *dbdata.Envelope) string {
	result := ""
	initialFormat := `{"base": "%v", "rates": %v}`
	baseFormat := "{%v}"
	delimiterFormat := ", "
	dataFormat := `"%v": "%v"`
	ratesHolder := ""

	sort.Slice(envelope.Cube, func(i, j int) bool {
		return envelope.Cube[i].Rate < envelope.Cube[j].Rate
	})

	for index, cube := range envelope.Cube {
		if index+1 >= len(envelope.Cube) {
			ratesHolder += fmt.Sprintf(dataFormat, cube.Currency, cube.Rate)
		} else {
			ratesHolder += fmt.Sprintf(dataFormat, cube.Currency, cube.Rate) + delimiterFormat
		}
	}

	result = fmt.Sprintf(initialFormat, envelope.SenderName, fmt.Sprintf(baseFormat, ratesHolder))

	return result
}

func (e *Envelope) downloadXMLData() (*xmldata.Envelope, error) {
	logger.Log.Infoln("Starting to download xml data")

	resp, err := http.Get(settings.GetXMLDataURLPath())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid response status: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil && data == nil {
		return nil, fmt.Errorf("Unable to read xml data: %v", err.Error())
	}

	env := &xmldata.Envelope{}
	err = xml.Unmarshal(data, env)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse xml data: %v", err.Error())
	}

	logger.Log.Infof("Successfully downloaded xml data. Status code: %v", resp.StatusCode)

	return env, nil
}

func (e *Envelope) useDemoData() *xmldata.Envelope {
	logger.Log.Warnf("Starting insertion of xml demo data")
	env := &xmldata.Envelope{}

	data, err := ioutil.ReadFile(settings.GetXMLDataFilePath())
	if err != nil {
		logger.Log.Fatalln("Unable to read demo data: %v", err)
	}

	if err = xml.Unmarshal(data, env); err != nil {
		logger.Log.Fatalln("Unable to parse demo data: %v", err)
	}

	logger.Log.Warnf("Currently using xml demo data")

	return env
}

func (e *Envelope) convertXMLtoDBEntities(xmlEnvelope *xmldata.Envelope) []dbdata.Envelope {
	dbEnvelopeList := []dbdata.Envelope{}

	for _, cube1 := range xmlEnvelope.Cube.Cube {
		dbCubeList := []dbdata.Cube{}

		for _, cube2 := range cube1.Cube {
			dbCubeList = append(dbCubeList, dbdata.Cube{
				Currency: cube2.Currency,
				Rate:     cube2.Rate,
			})
		}

		dbEnvelopeList = append(dbEnvelopeList, dbdata.Envelope{
			SenderName: xmlEnvelope.Sender.Name,
			CubeTime:   cube1.Time,
			Cube:       dbCubeList,
		})
	}

	return dbEnvelopeList
}

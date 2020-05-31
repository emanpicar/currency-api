package jsondata

type (
	QuantitativeExchangeRate struct {
		Base         string                  `json:"base"`
		RatesAnalyze map[string]RatesAnalyze `json:"rates_analyze"`
	}

	RatesAnalyze struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
		Avg float64 `json:"avg"`
	}

	ResponseMessage struct {
		Message string `json:"message"`
	}
)

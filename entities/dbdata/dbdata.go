package dbdata

import (
	"github.com/jinzhu/gorm"
)

type (
	Envelope struct {
		gorm.Model
		SenderName string `gorm:"type:varchar(255)"`
		CubeTime   string `gorm:"type:varchar(100)"`
		Cube       []Cube `gorm:"foreignkey:EnvelopeID"`
	}

	Cube struct {
		gorm.Model
		EnvelopeID uint
		Currency   string  `gorm:"type:varchar(10)"`
		Rate       float64 `gorm:"type:decimal(20,8)"`
	}

	QuantitativeExchangeRate struct {
		Base         string
		RatesAnalyze []RatesAnalyze
	}

	RatesAnalyze struct {
		Currency string
		Min      float64
		Max      float64
		Avg      float64
	}
)

func (Envelope) TableName() string {
	return "envelopes"
}

func (Cube) TableName() string {
	return "cubes"
}

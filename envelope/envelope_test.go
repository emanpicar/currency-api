package envelope

import (
	"errors"
	"reflect"
	"testing"

	"github.com/emanpicar/currency-api/entities/dbdata"
	"github.com/emanpicar/currency-api/entities/jsondata"
)

var (
	throwErrorInGetLatestRate, throwErrorInGetRateByDate, throwErrorInAnalyzedRate bool
	mockEnvelopeExpectedResult                                                     string          = `{"base": "Mock Sender", "rates": {"PHP": "50.999", "HPH": "999.5"}}`
	mockEnvelopeResult                                                             dbdata.Envelope = dbdata.Envelope{SenderName: "Mock Sender", Cube: []dbdata.Cube{
		dbdata.Cube{Currency: "PHP", Rate: 50.999},
		dbdata.Cube{Currency: "HPH", Rate: 999.50},
	}}
	mockAnalyzedResult jsondata.QuantitativeExchangeRate = jsondata.QuantitativeExchangeRate{
		Base:         "Mock Sender",
		RatesAnalyze: map[string]jsondata.RatesAnalyze{"PHP": jsondata.RatesAnalyze{50.555, 60.666, 55.555}},
	}
)

type (
	MockDBHandler struct{}
)

func (m MockDBHandler) BatchFirstOrCreate(dbEnvelopeList *[]dbdata.Envelope) {}
func (m MockDBHandler) GetLatestRates() (*dbdata.Envelope, error) {
	if throwErrorInGetLatestRate {
		return nil, errors.New("Record not found")
	}

	return &mockEnvelopeResult, nil
}
func (m MockDBHandler) GetRatesByDate(cubeTime string) (*dbdata.Envelope, error) {
	if throwErrorInGetRateByDate {
		return nil, errors.New("Record not found")
	}

	return &mockEnvelopeResult, nil
}
func (m MockDBHandler) GetAnalyzedRates() (*dbdata.QuantitativeExchangeRate, error) {
	if throwErrorInAnalyzedRate {
		return nil, errors.New("Record not found")
	}

	return &dbdata.QuantitativeExchangeRate{
		Base:         "Mock Sender",
		RatesAnalyze: []dbdata.RatesAnalyze{dbdata.RatesAnalyze{"PHP", 50.555, 60.666, 55.555}},
	}, nil
}

func TestEnvelope_GetLatestRates(t *testing.T) {
	tests := []struct {
		name    string
		e       *Envelope
		want    string
		wantErr bool
	}{
		struct {
			name    string
			e       *Envelope
			want    string
			wantErr bool
		}{
			name:    "Records found",
			e:       &Envelope{dbManager: &MockDBHandler{}},
			want:    mockEnvelopeExpectedResult,
			wantErr: false,
		},
		struct {
			name    string
			e       *Envelope
			want    string
			wantErr bool
		}{
			name:    "Records not found",
			e:       &Envelope{dbManager: &MockDBHandler{}},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			throwErrorInGetLatestRate = tt.wantErr
			got, err := tt.e.GetLatestRates()
			if (err != nil) != tt.wantErr {
				t.Errorf("Envelope.GetLatestRates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Envelope.GetLatestRates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvelope_GetRatesByDate(t *testing.T) {
	type args struct {
		cubeTime string
	}
	tests := []struct {
		name    string
		e       *Envelope
		args    args
		want    string
		wantErr bool
	}{
		struct {
			name    string
			e       *Envelope
			args    args
			want    string
			wantErr bool
		}{
			name:    "Records found",
			e:       &Envelope{dbManager: &MockDBHandler{}},
			args:    args{cubeTime: "2020-06-01"},
			want:    mockEnvelopeExpectedResult,
			wantErr: false,
		},
		struct {
			name    string
			e       *Envelope
			args    args
			want    string
			wantErr bool
		}{
			name:    "Records not found",
			e:       &Envelope{dbManager: &MockDBHandler{}},
			args:    args{cubeTime: "9999-66-11"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			throwErrorInGetRateByDate = tt.wantErr
			got, err := tt.e.GetRatesByDate(tt.args.cubeTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("Envelope.GetRatesByDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Envelope.GetRatesByDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvelope_GetAnalyzedRates(t *testing.T) {
	tests := []struct {
		name    string
		e       *Envelope
		want    *jsondata.QuantitativeExchangeRate
		wantErr bool
	}{
		struct {
			name    string
			e       *Envelope
			want    *jsondata.QuantitativeExchangeRate
			wantErr bool
		}{
			name:    "Records found",
			e:       &Envelope{dbManager: &MockDBHandler{}},
			want:    &mockAnalyzedResult,
			wantErr: false,
		},
		struct {
			name    string
			e       *Envelope
			want    *jsondata.QuantitativeExchangeRate
			wantErr bool
		}{
			name:    "Records not found",
			e:       &Envelope{dbManager: &MockDBHandler{}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			throwErrorInAnalyzedRate = tt.wantErr
			got, err := tt.e.GetAnalyzedRates()
			if (err != nil) != tt.wantErr {
				t.Errorf("Envelope.GetAnalyzedRates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Envelope.GetAnalyzedRates() = %v, want %v", got, tt.want)
			}
		})
	}
}

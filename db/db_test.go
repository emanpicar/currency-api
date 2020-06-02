package db

import (
	"database/sql"
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/emanpicar/currency-api/entities/dbdata"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
	gormDB  *gorm.DB
)

func beforeEach() {
	var err error
	mockDB, mockSQL, err = sqlmock.New()
	if err != nil {
		panic(err)
	}

	gormDB, err = gorm.Open("postgres", mockDB)
	if err != nil {
		panic(err)
	}
}

func afterEach() {
	mockDB.Close()
	gormDB.Close()
}

func Test_dbHandler_GetLatestRates(t *testing.T) {
	beforeEach()
	defer afterEach()

	tests := []struct {
		name          string
		dbHandler     *dbHandler
		want          *dbdata.Envelope
		wantErr       bool
		expectedQuery string
	}{
		struct {
			name          string
			dbHandler     *dbHandler
			want          *dbdata.Envelope
			wantErr       bool
			expectedQuery string
		}{
			name:          "Get latest rates - Success",
			dbHandler:     &dbHandler{database: gormDB},
			want:          &dbdata.Envelope{SenderName: "Dummy Sender", CubeTime: "2020-06-02"},
			wantErr:       false,
			expectedQuery: `SELECT \* FROM \"envelopes\" (.+) ORDER BY cube_time desc(.+) LIMIT 1`,
		},
		struct {
			name          string
			dbHandler     *dbHandler
			want          *dbdata.Envelope
			wantErr       bool
			expectedQuery string
		}{
			name:          "Get latest rates - Failed",
			dbHandler:     &dbHandler{database: gormDB},
			want:          nil,
			wantErr:       true,
			expectedQuery: `SELECT \* FROM \"envelopes\" (.+) ORDER BY cube_time desc(.+) LIMIT 1`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedHeader := []string{"sender_name", "cube_time"}
			if tt.wantErr {
				mockSQL.ExpectQuery(tt.expectedQuery).WillReturnRows(sqlmock.NewRows(nil))
			} else {
				mockSQL.ExpectQuery(tt.expectedQuery).WillReturnRows(sqlmock.NewRows(expectedHeader).
					AddRow("Dummy Sender", "2020-06-02"))
			}

			got, err := tt.dbHandler.GetLatestRates()
			if (err != nil) != tt.wantErr {
				t.Errorf("dbHandler.GetLatestRates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err = mockSQL.ExpectationsWereMet(); err != nil {
				t.Errorf("mockSQL.ExpectationsWereMet() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dbHandler.GetLatestRates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dbHandler_GetRatesByDate(t *testing.T) {
	beforeEach()
	defer afterEach()
	dummyCubeTime := "2010-12-22"

	type args struct {
		cubeTime string
	}
	tests := []struct {
		name          string
		dbHandler     *dbHandler
		args          args
		want          *dbdata.Envelope
		wantErr       bool
		expectedQuery string
	}{
		struct {
			name          string
			dbHandler     *dbHandler
			args          args
			want          *dbdata.Envelope
			wantErr       bool
			expectedQuery string
		}{
			name:          "Get rates by date - Success",
			dbHandler:     &dbHandler{database: gormDB},
			args:          args{cubeTime: dummyCubeTime},
			want:          &dbdata.Envelope{SenderName: "Dummy Sender", CubeTime: dummyCubeTime},
			wantErr:       false,
			expectedQuery: `SELECT \* FROM \"envelopes\" WHERE (.+)\"envelopes\"\.\"cube_time\" = (.+) LIMIT 1`,
		},
		struct {
			name          string
			dbHandler     *dbHandler
			args          args
			want          *dbdata.Envelope
			wantErr       bool
			expectedQuery string
		}{
			name:          "Get rates by date - Failed",
			dbHandler:     &dbHandler{database: gormDB},
			args:          args{cubeTime: dummyCubeTime},
			want:          nil,
			wantErr:       true,
			expectedQuery: `SELECT \* FROM \"envelopes\" WHERE (.+)\"envelopes\"\.\"cube_time\" = (.+) LIMIT 1`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedHeader := []string{"sender_name", "cube_time"}
			if tt.wantErr {
				mockSQL.ExpectQuery(tt.expectedQuery).WillReturnRows(sqlmock.NewRows(nil))
			} else {
				mockSQL.ExpectQuery(tt.expectedQuery).WillReturnRows(sqlmock.NewRows(expectedHeader).
					AddRow("Dummy Sender", dummyCubeTime))
			}

			got, err := tt.dbHandler.GetRatesByDate(tt.args.cubeTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("dbHandler.GetRatesByDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err = mockSQL.ExpectationsWereMet(); err != nil {
				t.Errorf("mockSQL.ExpectationsWereMet() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dbHandler.GetRatesByDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dbHandler_GetAnalyzedRates(t *testing.T) {
	beforeEach()
	defer afterEach()

	tests := []struct {
		name             string
		dbHandler        *dbHandler
		want             *dbdata.QuantitativeExchangeRate
		wantErr          bool
		expected1stQuery string
		expected2ndQuery string
	}{
		struct {
			name             string
			dbHandler        *dbHandler
			want             *dbdata.QuantitativeExchangeRate
			wantErr          bool
			expected1stQuery string
			expected2ndQuery string
		}{
			name:      "Get analyzed rates - Success",
			dbHandler: &dbHandler{database: gormDB},
			want: &dbdata.QuantitativeExchangeRate{
				Base:         "Dummy Sender",
				RatesAnalyze: []dbdata.RatesAnalyze{dbdata.RatesAnalyze{"PHP", 50.555, 60.555, 55.555}},
			},
			wantErr:          false,
			expected1stQuery: `SELECT \* FROM \"envelopes\" (.+) LIMIT 1`,
			expected2ndQuery: `SELECT currency, min\(rate\) AS min, max\(rate\) AS max, avg\(rate\) AS avg FROM cubes GROUP BY currency`,
		},
		struct {
			name             string
			dbHandler        *dbHandler
			want             *dbdata.QuantitativeExchangeRate
			wantErr          bool
			expected1stQuery string
			expected2ndQuery string
		}{
			name:             "Get analyzed rates - Failed",
			dbHandler:        &dbHandler{database: gormDB},
			want:             nil,
			wantErr:          true,
			expected1stQuery: `SELECT \* FROM \"envelopes\" (.+) LIMIT 1`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				mockSQL.ExpectQuery(tt.expected1stQuery).WillReturnRows(sqlmock.NewRows(nil))
			} else {
				mockSQL.ExpectQuery(tt.expected1stQuery).WillReturnRows(sqlmock.NewRows([]string{"sender_name"}).
					AddRow("Dummy Sender"))
				mockSQL.ExpectQuery(tt.expected2ndQuery).
					WillReturnRows(sqlmock.NewRows([]string{"currency", "min", "max", "avg"}).
						AddRow("PHP", 50.555, 60.555, 55.555))
			}

			got, err := tt.dbHandler.GetAnalyzedRates()
			if (err != nil) != tt.wantErr {
				t.Errorf("dbHandler.GetAnalyzedRates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err = mockSQL.ExpectationsWereMet(); err != nil {
				t.Errorf("mockSQL.ExpectationsWereMet() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dbHandler.GetAnalyzedRates() = %v, want %v", got, tt.want)
			}
		})
	}
}

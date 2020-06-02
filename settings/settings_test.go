package settings

import (
	"os"
	"testing"
)

func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		struct {
			name string
			want string
		}{
			name: "Loglevel error",
			want: "error",
		},
		struct {
			name string
			want string
		}{
			name: "Loglevel panic",
			want: "panic",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("LOG_LEVEL", tt.want)
			if got := GetLogLevel(); got != tt.want {
				t.Errorf("GetLogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDBHost(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		struct {
			name string
			want string
		}{
			name: "DBHost 1",
			want: "9.9.9.9",
		},
		struct {
			name string
			want string
		}{
			name: "DBHost 2",
			want: "10.12.12.45",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("DB_HOST", tt.want)
			if got := GetDBHost(); got != tt.want {
				t.Errorf("GetDBHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

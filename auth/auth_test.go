package auth

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	mockTokenResult              = "dummyTokenGeneratedByMockJWTHandler"
	throwUnexpectedSigningMethod = "throwUnexpectedSigningMethod"
)

type (
	jwtHandlerMock struct{}
)

func (j *jwtHandlerMock) generateJwtToken(mapClaims jwt.MapClaims) (string, error) {
	return mockTokenResult, nil
}

func (j *jwtHandlerMock) parseJwtToken(tokenString string) error {
	if throwUnexpectedSigningMethod == tokenString {
		return errors.New("Unexpected signing method")
	}

	return nil
}

func Test_authHandler_Authenticate(t *testing.T) {
	mockJSONUser := `{"username": "user123", "password": "pass123"}`
	mockJSONWrongUser := `{"username": "user123", "password": "wrongpass123"}`
	type args struct {
		body io.ReadCloser
	}
	tests := []struct {
		name    string
		a       *authHandler
		args    args
		want    string
		wantErr bool
	}{
		struct {
			name    string
			a       *authHandler
			args    args
			want    string
			wantErr bool
		}{
			name:    "Valid username/password",
			a:       &authHandler{&jwtHandlerMock{}},
			args:    args{body: ioutil.NopCloser(bytes.NewReader([]byte(mockJSONUser)))},
			want:    mockTokenResult,
			wantErr: false,
		},
		struct {
			name    string
			a       *authHandler
			args    args
			want    string
			wantErr bool
		}{
			name:    "Invalid username/password",
			a:       &authHandler{&jwtHandlerMock{}},
			args:    args{body: ioutil.NopCloser(bytes.NewReader([]byte(mockJSONWrongUser)))},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Authenticate(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("authHandler.Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("authHandler.Authenticate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authHandler_ValidateRequest(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name        string
		a           *authHandler
		args        args
		wantErr     bool
		errorPrefix string
	}{
		struct {
			name        string
			a           *authHandler
			args        args
			wantErr     bool
			errorPrefix string
		}{
			name:    "Valid token",
			a:       &authHandler{&jwtHandlerMock{}},
			args:    args{r: &http.Request{Header: http.Header{"Authorization": []string{"Bearer ValidToken"}}}},
			wantErr: false,
		},
		struct {
			name        string
			a           *authHandler
			args        args
			wantErr     bool
			errorPrefix string
		}{
			name:        "Missing authorization header",
			a:           &authHandler{&jwtHandlerMock{}},
			args:        args{r: &http.Request{}},
			wantErr:     true,
			errorPrefix: "An authorization header is required",
		},
		struct {
			name        string
			a           *authHandler
			args        args
			wantErr     bool
			errorPrefix string
		}{
			name:        "Cannot parse authorization header",
			a:           &authHandler{&jwtHandlerMock{}},
			args:        args{r: &http.Request{Header: http.Header{"Authorization": []string{"BearerNospacetoken"}}}},
			wantErr:     true,
			errorPrefix: "Cannot parse authorization header",
		},
		struct {
			name        string
			a           *authHandler
			args        args
			wantErr     bool
			errorPrefix string
		}{
			name:        "Unexpected signing method",
			a:           &authHandler{&jwtHandlerMock{}},
			args:        args{r: &http.Request{Header: http.Header{"Authorization": []string{fmt.Sprintf("Bearer %v", throwUnexpectedSigningMethod)}}}},
			wantErr:     true,
			errorPrefix: "Unexpected signing method",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.a.ValidateRequest(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("authHandler.ValidateRequest() error = %v, wantErr %v", err, tt.wantErr)
			}

			if (err != nil) && !strings.HasPrefix(err.Error(), tt.errorPrefix) {
				t.Errorf("authHandler.ValidateRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

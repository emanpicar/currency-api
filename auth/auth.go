package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/emanpicar/currency-api/logger"
	"github.com/emanpicar/currency-api/settings"

	jwt "github.com/dgrijalva/jwt-go"
)

type (
	Manager interface {
		Authenticate(body io.ReadCloser) (string, error)
		ValidateRequest(r *http.Request) error
	}

	authHandler struct{}

	User struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
)

func NewManager() Manager {
	return &authHandler{}
}

func (a *authHandler) Authenticate(body io.ReadCloser) (string, error) {
	var user User
	if err := json.NewDecoder(body).Decode(&user); err != nil {
		return "", err
	}

	// TODO validate user and pass in DB
	if !a.inMemoryAuthentication(user) {
		return "", errors.New("Invalid user username/password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":   user.Username,
		"authorized": true,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(settings.GetTokenSecret()))
	if err != nil {
		return "", err
	}
	logger.Log.Infoln("Successfully generated JWT token")

	return tokenString, nil
}

func (a *authHandler) ValidateRequest(r *http.Request) error {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return errors.New("An authorization header is required")
	}

	bearerToken := strings.Split(authorizationHeader, " ")
	if len(bearerToken) != 2 {
		return errors.New("Cannot parse authorization header")
	}

	token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(settings.GetTokenSecret()), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("Invalid authorization token")
	}

	return nil
}

func (a *authHandler) inMemoryAuthentication(userCreds User) bool {
	logger.Log.Infof("Authenticating user with username: %v", userCreds.Username)

	users := []User{
		User{Username: "user123", Password: "pass123"},
		User{Username: "useruser", Password: "passpass"},
	}

	for _, val := range users {
		if val.Username == userCreds.Username && val.Password == userCreds.Password {
			return true
		}
	}

	return false
}

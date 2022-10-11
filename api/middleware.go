package api

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"equiprent/internal/util/flags"
	"equiprent/internal/util/log"

	"github.com/dgrijalva/jwt-go"
)

var secret = os.Getenv("EQUIPRENT_SECRETKEY")

// middleware is applied 0->len(routes middleware) then 0->len(route middleware)
// i.e. routes first

type middleware interface {
	Handle(next http.HandlerFunc) (handled http.HandlerFunc)
}

type LogMiddleware struct {
	isDev bool
}

var logMiddleware LogMiddleware

func initMiddleware() {
	logMiddleware.isDev = *flags.Dev
}

func (l *LogMiddleware) Handle(next http.HandlerFunc) (handled http.HandlerFunc) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Logger.Debug("API Access: ", r.Method, " ", r.URL.Path)
		next(w, r)
	})
}

type AuthMiddleware struct{}

var authMiddleware AuthMiddleware

func (a *AuthMiddleware) Handle(next http.HandlerFunc) (handled http.HandlerFunc) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := checkToken(r)
		if err != nil {
			log.Logger.Debug(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			writeError(w, "unauthorized")
			return
		}
		next(w, r)
	})
}

func checkToken(r *http.Request) (token *jwt.Token, err error) {
	tokenSlice := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	if len(tokenSlice) != 2 {
		err = errors.New("bad token attempt")
		return
	}
	tokenString := tokenSlice[1]
	token, err = jwt.Parse(tokenString, jwt.Keyfunc(func(*jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}))
	if err != nil {
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("could not read claims")
		return
	}
	if !claims.VerifyAudience("https://cachegrab.org", true) || !claims.VerifyIssuer("cachegrab", true) {
		err = errors.New("bad audience or issuer")
		return
	}
	return
}

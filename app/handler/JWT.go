package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"

	"github.com/lacazethomas/goTodo/app/model"
	"github.com/lacazethomas/goTodo/config"
)

func JwtAuthentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/user/new", "/api/user/login"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path                               //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header
		log.Print(tokenHeader)

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			respondError(w, http.StatusForbidden, "Missing auth token")
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			respondError(w, http.StatusForbidden, "Invalid/Malformed auth token")
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		tk := &model.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetTokenString()), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			respondError(w, http.StatusForbidden, "Malformed authentication token")
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			respondError(w, http.StatusForbidden, "Token is not valid.")
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		fmt.Sprintf("User %", tk.UserId) //Useful for monitoring
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}

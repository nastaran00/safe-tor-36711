package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

type userKey string

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//List of endpoints that doesn't require auth
		notAuth := []string{
			"/user/new",
			"/user/login",
			"/api/quote",
			"/user/home",
			"/user/eddition",
			"/user/logout",
		}

		if needsAuth := os.Getenv("NEEDS_AUTH"); needsAuth == "yes" {
			notAuth = []string{
				"/user/new",
				"/user/login",
				"/recommendersystem",
			}
		}

		requestPath := r.URL.Path //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		var response map[string]interface{}
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //token is missing, returns with error code 403 Unauthorized
			response = message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			respond(w, response)
			return
		}
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 1 {
			response = message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			respond(w, response)
			return
		}

		tokenPart := splitted[0] //Grab the token part, what we are truly interested in
		tk := &token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_PASSWORD")), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual

			issue := fmt.Sprintf("Problem with authentication token: %s", err)

			if e, ok := err.(*jwt.ValidationError); ok {
				if e.Errors&jwt.ValidationErrorExpired != 0 {
					issue = fmt.Sprintf("Your token is too old! %s", err)
				}
			}

			fmt.Println(err)

			response = message(false, issue)

			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			respond(w, response)
			return
		}

		if !token.Valid { //token is invalid, maybe not signed on this server
			response = message(false, "token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			respond(w, response)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token

		ctx := context.WithValue(r.Context(), userKey("user"), tk.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}

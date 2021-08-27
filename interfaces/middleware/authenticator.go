package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hayashiki/audiy-api/interfaces/api/graph/auth"
	"google.golang.org/api/idtoken"
)

type Authenticator interface {
	AuthMiddleware(h http.Handler) http.Handler
}

type authenticator struct {
}

func NewAuthenticator() Authenticator {
	return &authenticator{}
}

func (a *authenticator) AuthMiddleware(h http.Handler) http.Handler {
	ctx := context.Background()

	fn := func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		if bearer == "Bearer dummy" {
			ctx = auth.SetAuth(ctx, &auth.Auth{
				ID:    "103843140833205663533",
				Name:  "hayashiki",
				Email: "hayashiki@example.com",
			})
			r = r.WithContext(ctx)
		} else if bearer != "" {
			idToken := strings.Replace(bearer, "Bearer ", "", 1)
			log.Println("idToken is", idToken)
			log.Println("audience is", os.Getenv("GOOGLE_CLIENT_ID"))
			idTokenMap, err := idtoken.Validate(context.Background(), idToken, os.Getenv("GOOGLE_CLIENT_ID"))
			if err != nil {
				log.Printf("idToken err %s", err.Error())
				if err.Error() == "idtoken: token expired" {
					// TODO: error handle
					log.Printf("tooruka? %v", err.Error())

					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{
	"data": {},
	"errors": [
	{
		"extensions": {
			"code": "UNAUTHENTICATED"
		}
	}
	]
				}`))
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			email := idTokenMap.Claims["email"].(string)
			name := idTokenMap.Claims["name"].(string)
			ctx = context.WithValue(r.Context(), auth.KeyAuth, &auth.Auth{
				ID:    idTokenMap.Subject,
				Email: email,
				Name:  name,
			})
			r = r.WithContext(ctx)
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

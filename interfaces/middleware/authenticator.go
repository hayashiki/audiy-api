package middleware

import (
	"context"
	"github.com/hayashiki/audiy-api/interfaces/api/graph"
	"google.golang.org/api/idtoken"
	"log"
	"net/http"
	"strings"
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
			ctx = graph.SetAuth(ctx, &graph.Auth{
				ID:   111111,
				Name: "hayashiki",
				Email: "hayashiki@example.com",
			})
			r = r.WithContext(ctx)
		} else if bearer != "" {
			idToken := strings.Replace(bearer, "Bearer ", "", 1)
			idTokenMap, err := idtoken.Validate(context.Background(), idToken, "185245971175-sctlo6t5hkgr2mu1qnkgp3s54hju8bi2.apps.googleusercontent.com")
			log.Printf("idTokenMap %+v, %+v", idTokenMap, err)
			id := int64(idTokenMap.Claims["id"].(float64))
			email := idTokenMap.Claims["email"].(string)
			name := idTokenMap.Claims["fullname"].(string)
			ctx = context.WithValue(r.Context(), graph.KeyAuth, &graph.Auth{
				ID:   id,
				Email: email,
				Name: name,
			})
			r = r.WithContext(ctx)
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

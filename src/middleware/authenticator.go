package middleware

import (
	"context"
	fbauth "firebase.google.com/go/auth"
	"log"
	"net/http"
	"strings"

	"github.com/hayashiki/audiy-api/src/graph/auth"
)

type Authenticator interface {
	AuthMiddleware(h http.Handler) http.Handler
}

type authenticator struct {
	client *fbauth.Client
}

func NewAuthenticator(client *fbauth.Client) Authenticator {
	return &authenticator{client: client}
}

func (a *authenticator) AuthMiddleware(h http.Handler) http.Handler {
	ctx := context.Background()

	fn := func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		// TODO: only debug mode
		if bearer == "Bearer dummy" {
			ctx = auth.SetAuth(ctx, &auth.Auth{
				ID:    "103843140833205663533",
				Name:  "hayashiki",
				Email: "hayashiki@example.com",
			})
			r = r.WithContext(ctx)
		} else if bearer != "" {
			idToken := strings.Replace(bearer, "Bearer ", "", 1)
			//idTokenMap, err := idtoken.Validate(context.Background(), idToken, os.Getenv("GOOGLE_CLIENT_ID"))

			// idTokenMap &{1640053393 https://securetoken.google.com/bulb-audiy bulb-audiy 1640234698 1640231098 JWDxVdaSIpN31vrJc9AAsnsmLXA2 JWDxVdaSIpN31vrJc9AAsnsmLXA2 {google.com  map[email:[masayuki.hayashida@bulbcorp.jp] google.com:[103843140833205663533]]} map[auth_time:1.640053393e+09 email:masayuki.hayashida@bulbcorp.jp email_verified:true firebase:map[identities:map[email:[masayuki.hayashida@bulbcorp.jp] google.com:[103843140833205663533]] sign_in_provider:google.com] name:masayuki hayashida picture:https://lh3.googleusercontent.com/a-/AOh14GhXgZoEHOEZacTv4q3KkDluaDqdlP1wa43mpTnj=s96-c user_id:JWDxVdaSIpN31vrJc9AAsnsmLXA2]}
			idTokenMap, err := a.client.VerifyIDToken(ctx, idToken)
			log.Println("idTokenMap", idTokenMap)

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

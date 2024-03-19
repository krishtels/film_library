package router

import (
	"errors"
	"log"
	"net/http"

	"film-library/src/internal/tools"
	"github.com/golang-jwt/jwt/v5"
)

func NewAuthMiddleware(key string, adminOnly bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("jwt")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					log.Printf("ERROR: no auth cookie\n")
					tools.Unauthorized(w, r)
					return
				}

				log.Printf("ERROR: failed to get auth cookie")
				tools.InternalServerError(w, r)
				return
			}

			token := cookie.Value
			uc, err := tools.ParseUserClaims(token, key)
			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					log.Printf("ERROR: token expired")
					tools.Unauthorized(w, r)
					return
				}
				if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
					log.Printf("ERROR: jwt signature is invalid\n")
					tools.Unauthorized(w, r)
					return
				}
				if errors.Is(err, tools.ErrUnknownClaimsType) {
					log.Printf("ERROR: unknown token claims\n")
					tools.Unauthorized(w, r)
					return
				}

				log.Printf("ERROR: failed to parse claims err=%s", err.Error())
				tools.InternalServerError(w, r)
				return
			}

			log.Printf("INFO: authenticated user.go=%s user_id=%d", uc.Username, uc.ID)
			if adminOnly {
				if !uc.IsAdmin {
					log.Printf("ERROR: permission denied")
					tools.Forbidden(w, r)
					return
				}

				log.Printf("INFO: admin request")
			}

			next.ServeHTTP(w, r)
		})
	}
}

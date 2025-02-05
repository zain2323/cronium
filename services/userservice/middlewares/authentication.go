package middlewares

import (
	"context"
	"github.com/zain2323/cronium/services/userservice/utils"
	"net/http"
	"strings"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// getting request context
		ctx := r.Context()

		// stripping out the token
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(tokenString, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString = tokenParts[1]
		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// adding user_id to the context if token is verified
		context.WithValue(ctx, "user_id", claims["user_id"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

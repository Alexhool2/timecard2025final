package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/alexhool2/TimeCard/config"
	"github.com/alexhool2/TimeCard/types"
	"github.com/alexhool2/TimeCard/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userKey contextKey = "user"

func CreateJWT(secret []byte, userID int, role string) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationsInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"role":      role,
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Primeiro, procure o token no cookie
		cookie, err := r.Cookie("auth_token")
		var token string
		if err == nil {
			token = cookie.Value
		}

		// Se o token não estiver no cookie, procure no cabeçalho Authorization
		if token == "" {
			authHeader := r.Header.Get("Authorization")
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				token = authHeader[7:]
			}
		}

		// Se ainda não encontrar o token, retorne 401
		if token == "" {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}

		// Valide o token
		user, err := validateToken(token)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: %v", err))
			return
		}

		// Adicione o usuário ao contexto
		ctx := context.WithValue(r.Context(), userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateToken(t string) (*types.User, error) {

	token, err := jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, err := strconv.Atoi(claims["userID"].(string))
		if err != nil {
			return nil, fmt.Errorf("invalid userID in token")
		}
		role := claims["role"].(string)
		return &types.User{ID: userID, Role: types.Role(role)}, nil
	}
	return nil, fmt.Errorf("invalid token claims")
}

// checks if user is admin

func AdminOnly(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, err := GetUserFromContext(r.Context())
		if err != nil || user.Role != "admin" {
			http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetUserFromContext(ctx context.Context) (*types.User, error) {
	user, ok := ctx.Value(userKey).(*types.User)
	if !ok || user == nil {
		return nil, fmt.Errorf("user not found in context")
	}
	return user, nil
}

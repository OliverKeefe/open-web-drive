package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type JWTRequest struct {
	Alg    string `json:"alg"`
	Typ    string `json:"typ"`
	Sub    string `json:"sub"`
	Claims string `json:"name"`
	Admin  string `json:"admin"`
	Iat    string `json:"iat"`
}

type Claims struct {
	Claims interface{}
}

type OAuth2 interface {
	Authenticator(ctx context.Context, rawJWT string) (*Claims, error)
	ReissueJWT(ctx context.Context, refreshToken string) (string, error)
}

type UCAN interface {
	ValidateJWT(token jwt.Token) (bool, error)
}

func UserIDFromCtx(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

func ClaimsFromCtx(ctx context.Context) (map[string]interface{}, bool) {
	claims, ok := ctx.Value(claimsKey).(map[string]interface{})
	return claims, ok
}

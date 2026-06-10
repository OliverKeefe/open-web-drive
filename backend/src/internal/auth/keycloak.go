package auth

import (
	"context"
	"errors"
	"log"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type ctxStringKey string

type ctxIntKey int

const userIDKey ctxStringKey = "userID"
const claimsKey ctxIntKey = iota

type Authenticator struct {
	Issuer  string
	KeyFunc keyfunc.Keyfunc
}

func New(issuer, jwksUrl string) (*Authenticator, error) {
	kf, err := keyfunc.NewDefault([]string{jwksUrl})
	if err != nil {
		log.Printf("Failed to create JWK Set from resource at the given URL, %v", err)
		return nil, err
	}

	return &Authenticator{
		Issuer:  issuer,
		KeyFunc: kf,
	}, nil
}

func (k *Authenticator) ValidateJWT(ctx context.Context, jwtB64 string) (context.Context, bool, error) {
	claims := &jwt.RegisteredClaims{}
	kf := func(t *jwt.Token) (any, error) {
		return k.KeyFunc.Keyfunc(t)
	}
	token, err := jwt.ParseWithClaims(
		jwtB64,
		claims,
		kf,
		jwt.WithValidMethods([]string{"RS256"}),
	)
	if err != nil {
		log.Printf("failed to parse the JWT. %v", err)
		return ctx, false, err
	}

	if !token.Valid {
		log.Printf("invalid token.")
		return ctx, false, nil
	}

	if claims.Issuer != k.Issuer {
		return ctx, false, errors.New("invalid issuer")
	}

	newCtx := context.WithValue(ctx, userIDKey, claims.Subject)
	return newCtx, true, nil

}

func (k *Authenticator) ReissueJWT() (jwt.Token, error) {
	var tkn jwt.Token
	return tkn, nil
}

func UserIDFromCtx(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

func (k *Authenticator) InjectUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func (k *Authenticator) InjectClaims(ctx context.Context, claims map[string]interface{}) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

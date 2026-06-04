package app

import (
	"backend/src/internal/app/router"
	"backend/src/internal/auth"
	"backend/src/internal/database"
	"backend/src/internal/middleware"
	"context"
	"fmt"
	"log"
	"net/http"
)

func Run() error {
	ctx := context.Background()

	a, err := registerKeycloakAuth(
		"http://127.0.0.1:8080/realms/gestalt",
		"http://127.0.0.1:8080/realms/gestalt/protocol/openid-connect/certs",
	)
	if err != nil {
		return err
	}

	db, err := database.New(ctx, "DATABASE_URL")
	if err != nil {
		return err
	}

	appMux := http.NewServeMux()
	var handler http.Handler = appMux
	handler = middleware.EnableCORS(handler)
	router.RegisterFileRoutes(appMux, a, db)
	srv := &http.Server{
		Addr:    ":8081",
		Handler: handler,
	}

	log.Println("running on port: 8081...")
	return srv.ListenAndServe()
}

func registerKeycloakAuth(issuer string, jwksURL string) (*auth.Authenticator, error) {
	authenticator, err := auth.New(issuer, jwksURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create new authenticator, %v", err)
	}
	return authenticator, nil
}

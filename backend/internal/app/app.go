package app

import (
	"backend/internal/app/router"
	"backend/internal/auth"
	"backend/internal/database"
	"backend/internal/middleware"
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
	if err = router.RegisterFileRoutes(appMux, a, db); err != nil {
		return err
	}
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

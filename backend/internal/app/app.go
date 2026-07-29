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
	"os"
)

func Run() error {
	ctx := context.Background()

	keycloakIssuer := os.Getenv("KEYCLOAK_ISSUER_URL")
	if keycloakIssuer == "" {
		return fmt.Errorf("KEYCLOAK_ISSUER_URL is required")
	}
	keycloakJwks := os.Getenv("KEYCLOAK_JWKS_URL")
	if keycloakJwks == "" {
		return fmt.Errorf("KEYCLOAK_JWKS_URL is required")
	}

	a, err := registerKeycloakAuth(keycloakIssuer, keycloakJwks)
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

	port := os.Getenv("PORT")
	if port == "" {
		return fmt.Errorf("PORT is required")
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	log.Printf("running on port: %s...", port)
	return srv.ListenAndServe()
}

func registerKeycloakAuth(issuer string, jwksURL string) (*auth.Authenticator, error) {
	authenticator, err := auth.New(issuer, jwksURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create new authenticator, %v", err)
	}
	return authenticator, nil
}

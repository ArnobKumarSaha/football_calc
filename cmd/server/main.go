package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/ArnobKumarSaha/football_calc/pkg/api"
	"github.com/ArnobKumarSaha/football_calc/pkg/db"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

//go:embed static
var staticFiles embed.FS

func main() {
	root := &cobra.Command{
		Use:          "football-calc",
		SilenceUsage: true,
		RunE:         run,
	}
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(_ *cobra.Command, _ []string) error {
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL not set")
	}
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		return fmt.Errorf("ADMIN_PASSWORD not set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()

	pool, err := db.Connect(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("db connect: %w", err)
	}
	defer pool.Close()

	if err := db.Bootstrap(ctx, pool); err != nil {
		return fmt.Errorf("db bootstrap: %w", err)
	}

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("static fs: %w", err)
	}

	router := api.NewRouter(pool, adminPassword, staticFS)

	klog.Infof("starting server on :%s", port)
	return http.ListenAndServe(":"+port, router)
}

package main

import (
	"backend/models"
	"backend/pkg/config"
	"backend/pkg/db"
	"backend/pkg/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var repo *handlers.Repository

func main() {
	godotenv.Load(".env")

	var cfg config.Config
	cfg.Port, _ = strconv.Atoi(os.Getenv("PortNumber"))
	cfg.Version = os.Getenv("Version")
	cfg.Env = os.Getenv("Environment")
	cfg.Db.Dsn = os.Getenv("DbDsn")
	cfg.Db.Driver = os.Getenv("DbDriver")
	cfg.Jwt.Secret = os.Getenv("JwtSecret")

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := db.OpenDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	app := config.Application{
		Config: cfg,
		Logger: logger,
		Models: models.NewModels(db),
	}

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Running in port:", cfg.Port)

	err = srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

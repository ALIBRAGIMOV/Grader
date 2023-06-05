package main

import (
	"encoding/json"
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"grader/pkg/grader"
	graderDelivery "grader/pkg/grader/delivery"
	graderRepository "grader/pkg/grader/repo"
	graderService "grader/pkg/grader/service"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type config struct {
	Login    string
	Password string
}

func main() {
	flag.Parse()
	var err error

	err = godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	cfg := config{}
	cfg.Login = os.Getenv("GRADER_LOGIN")
	cfg.Password = os.Getenv("GRADER_PASSWORD")

	formData := url.Values{
		"username": {cfg.Login},
		"password": {cfg.Password},
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	config := &grader.Config{}
	data, err := ioutil.ReadFile("../../configs/config.json")
	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatalln("cant parse config:", err)
	}

	port := ":8080"
	r := chi.NewRouter()

	graderRepo := graderRepository.NewGraderRepo()
	graderService := graderService.NewGraderService(config, graderRepo)
	graderHandler := &graderDelivery.GraderHandler{
		GraderService: graderService,
		Logger:        logger,
		FormData:      formData,
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	graderHandler.GraderLogin()
	r.Post("/api/v1/grader/grade", graderHandler.GradeSolution)

	log.Printf("Grader start on port %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}

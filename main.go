package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/nilesh93/kubecon-rubik-challenge/helpers"
	"github.com/nilesh93/kubecon-rubik-challenge/services"
	"github.com/sirupsen/logrus"
)

type ChallengeRequest struct {
	UserEmail    string `json:"userEmail,omitempty" example:"nilesh@wso2.com"`
	ChallengeWon bool   `json:"challengeWon,omitempty" example:"true"`
}
type ChallengeResponse struct {
	Id      string `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// @title Rubik Challenge API documentation
// @version 1.0.0
// @BasePath /api/v1
func main() {

	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(cors.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Route("/healthz", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			helpers.RespondwithJSON(w, 200, "healthy")
		})
	})

	r.Route("/api/v1/", func(r chi.Router) {
		r.Post("/challenge", ChallengeRoute)
	})

	logrus.Info("http server started")
	http.ListenAndServe(":4000", r)
}

// @Summary Send Challenge Result Email
// @Tags Email
// @Accept json
// @Produce json
// @Param data body ChallengeRequest	true	"data"
// @Success 200 {object} ChallengeResponse	"Okay"
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /api/v1/challenge [post]
func ChallengeRoute(w http.ResponseWriter, r *http.Request) {
	body := ChallengeRequest{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	obj, err := ChallengeHandler(ctx, body)
	if err != nil {
		helpers.RespondWithError(w, 500, err.Error())
		return
	}
	helpers.RespondwithJSON(w, 200, obj)
}

func ChallengeHandler(ctx context.Context, req ChallengeRequest) (*ChallengeResponse, error) {

	if v, err := helpers.IsValid(req); !v {
		return nil, err
	}

	res, err := services.SendEmail(ctx, req.UserEmail, req.ChallengeWon)
	if err != nil {
		return nil, err
	}

	body := ChallengeResponse{
		Id:      res,
		Message: "Success",
	}
	return &body, nil
}

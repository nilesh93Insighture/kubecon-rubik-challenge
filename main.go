package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/nilesh93/kubecon-rubik-challenge/helpers"
	"github.com/sirupsen/logrus"
)

type SayHelloRequest struct {
	Name string `json:"name,omitempty" example:"John Doe"`
}
type SayHelloResponse struct {
	Message string `json:"message,omitempty"`
}

// @title Say Hello API
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
		r.Post("/sayHello", HelloRoute)
		r.Get("/force-panic", func(w http.ResponseWriter, r *http.Request) {
			panic("force panic")
		})
	})

	logrus.Info("http server started")
	http.ListenAndServe(":4000", r)
}

// @Summary Say Hello
// @Tags Hello
// @Accept json
// @Produce json
// @Param data body SayHelloRequest	true	"data"
// @Success 200 {object} SayHelloResponse	"Okay"
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /api/v1/sayHello [post]
func HelloRoute(w http.ResponseWriter, r *http.Request) {
	body := SayHelloRequest{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	obj, err := HelloHandler(ctx, body)
	if err != nil {
		helpers.RespondWithError(w, 500, err.Error())
		return
	}
	helpers.RespondwithJSON(w, 200, obj)
}

func HelloHandler(ctx context.Context, req SayHelloRequest) (*SayHelloResponse, error) {

	if v, err := helpers.IsValid(req); !v {
		return nil, err
	}

	body := SayHelloResponse{

		Message: fmt.Sprintf("Hello, %s!", req.Name),
	}
	return &body, nil
}

// @Summary Force Panic
// @Tags Panic
// @Accept json
// @Produce json
// @Success 200 {object} string	"Okay"
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /force-panic [get]

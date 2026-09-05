package main

import (
	handler "JuniarProject/internal/handler/product"
	repository "JuniarProject/internal/repository/product"
	service "JuniarProject/internal/service/product"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "JuniarProject/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
)

type config struct {
	address string
	db      dbConfig
}

type dbConfig struct {
	address string
}

type application struct {
	config config
	db     *pgxpool.Pool
}

func (app *application) mount() http.Handler {
	port := getEnv("SERVER_PORT", "8081")
	swaggerAddress := fmt.Sprintf("http://localhost:%s/swagger/doc.json", port)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	productRepository := repository.NewProductRepository(app.db)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(swaggerAddress),
	))

	r.Route("/products", func(r chi.Router) {
		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.GetProducts)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", productHandler.GetProductByID)
			r.Put("/", productHandler.UpdateProduct)
			r.Delete("/", productHandler.DeleteProduct)
		})

	})

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.address,
		Handler:      h,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at addr %s", app.config.address)

	return srv.ListenAndServe()
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"arm_back/internal/config"
	"arm_back/internal/handler"
	"arm_back/internal/middleware"
	"arm_back/internal/repository"
	"arm_back/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("connected to PostgreSQL")

	// Repositories
	userRepo := repository.NewUserRepository(pool)
	exerciseRepo := repository.NewExerciseRepository(pool)
	workoutRepo := repository.NewWorkoutRepository(pool)

	// Services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	exerciseService := service.NewExerciseService(exerciseRepo)
	workoutService := service.NewWorkoutService(workoutRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	exerciseHandler := handler.NewExerciseHandler(exerciseService)
	workoutHandler := handler.NewWorkoutHandler(workoutService)

	// Router
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(authService))
	{
		exercises := api.Group("/exercises")
		{
			exercises.GET("", exerciseHandler.List)
			exercises.POST("", exerciseHandler.Create)
			exercises.GET("/:id", exerciseHandler.GetByID)
			exercises.PUT("/:id", exerciseHandler.Update)
			exercises.DELETE("/:id", exerciseHandler.Delete)
			exercises.GET("/:id/comments", exerciseHandler.ListComments)
			exercises.POST("/:id/comments", exerciseHandler.CreateComment)
			exercises.PUT("/:id/comments/:commentId", exerciseHandler.UpdateComment)
			exercises.DELETE("/:id/comments/:commentId", exerciseHandler.DeleteComment)
		}

		workouts := api.Group("/workouts")
		{
			workouts.GET("", workoutHandler.List)
			workouts.POST("", workoutHandler.Create)
			workouts.GET("/:id", workoutHandler.GetByID)
			workouts.PUT("/:id", workoutHandler.Update)
			workouts.DELETE("/:id", workoutHandler.Delete)
			workouts.POST("/:id/copy", workoutHandler.Copy)
			workouts.POST("/:id/exercises", workoutHandler.AddExercise)
			workouts.DELETE("/:id/exercises/:exerciseId", workoutHandler.RemoveExercise)
			workouts.POST("/:id/exercises/:exerciseId/sets", workoutHandler.AddSet)
			workouts.PUT("/:id/exercises/:exerciseId/sets/:setId", workoutHandler.UpdateSet)
			workouts.DELETE("/:id/exercises/:exerciseId/sets/:setId", workoutHandler.DeleteSet)
		}
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	go func() {
		log.Printf("server starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server stopped")
}

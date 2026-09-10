package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hyanhsing/hyandrop/backend/internal/config"
	"github.com/hyanhsing/hyandrop/backend/internal/httpapi"
	"github.com/hyanhsing/hyandrop/backend/internal/postgres"
	"github.com/hyanhsing/hyandrop/backend/internal/room"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	ctx := context.Background()

	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer db.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	defer redisClient.Close()

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"app":    cfg.AppName,
		})
	})

	router.GET("/readyz", func(c *gin.Context) {
		checkCtx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := db.Ping(checkCtx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not_ready",
				"check":  "postgres",
				"error":  err.Error(),
			})
			return
		}

		if err := redisClient.Ping(checkCtx).Err(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not_ready",
				"check":  "redis",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
			"checks": []string{"postgres", "redis"},
		})
	})

	api := router.Group("/api")
	roomRepository := postgres.NewRoomRepository(db)
	roomService := room.NewService(roomRepository, time.Duration(cfg.RoomTTLHours)*time.Hour, time.Now)
	roomHandler := httpapi.NewRoomHandler(roomService)
	roomHandler.Register(api)

	api.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":             cfg.AppName,
			"environment":      cfg.AppEnv,
			"room_ttl_hours":   cfg.RoomTTLHours,
			"max_file_size_mb": cfg.MaxFileSizeMB,
		})
	})

	log.Printf("%s backend listening on %s", cfg.AppName, cfg.HTTPAddr)
	if err := router.Run(cfg.HTTPAddr); err != nil {
		log.Fatalf("run server: %v", err)
	}
}

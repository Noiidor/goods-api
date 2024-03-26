package main

import (
	"context"
	"fmt"
	"goods-api/internal/broker/nats"
	"goods-api/internal/cache/redis"
	"goods-api/internal/config"
	"goods-api/internal/controllers"
	"goods-api/internal/db/clickhouse"
	"goods-api/internal/db/postgres"
	"goods-api/internal/errors"
	chRepos "goods-api/internal/repos/clickhouse"
	pgRepos "goods-api/internal/repos/postgres"
	"goods-api/internal/route"
	"goods-api/internal/services"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// ТЗ специально кривое?)
func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	db, err := postgres.NewConnection(&config)
	if err != nil {
		log.Fatalf("Could not establish persistance database connection: %v", err)
	}
	defer db.Close()

	analyticsDb, err := clickhouse.NewConnection(&config)
	if err != nil {
		log.Fatalf("Could not establish analytics database connection: %v", err)
	}
	defer analyticsDb.Close()

	redisConn, err := redis.NewConnection(&config)
	if err != nil {
		log.Fatalf("Could not establish cache database connection: %v", err)
	}
	defer redisConn.Close()

	natsConn, err := nats.NewConnection(&config)
	if err != nil {
		log.Fatalf("Could not establish connection with NATS: %v", err)
	}
	defer natsConn.Drain()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errors.HandleAsyncErrors()

	goodsCache := redis.New(redisConn)

	messageBroker := nats.New(natsConn)

	goodsRepo := pgRepos.NewGoodsRepo(db)
	projectsRepo := pgRepos.NewProjectsRepo(db)
	analyticsRepo := chRepos.NewGoodsAnalytics(analyticsDb)

	goodsService := services.NewGoodsService(goodsRepo, projectsRepo, goodsCache, messageBroker)
	analyticsService := services.NewAnalyticsService(analyticsRepo, messageBroker)
	analyticsService.HandleQueue(ctx)

	goodsController := controllers.NewGoodsController(goodsService)

	projectsRepo.Insert("Auto-created project")
	if err != nil {
		log.Printf("failed to create project record, err: %v", err)
		return // Return вместо log.Fatal для того что бы defer`ы сработали
	}

	api := gin.Default()
	route.RouteGoods(api, *goodsController)

	srvr := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", config.Server.Host, config.Server.Port),
		Handler: api,
	}

	go func() {
		if err := srvr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error listening on address %v: %v", srvr.Addr, err)
		}
	}()

	<-ctx.Done()

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srvr.Shutdown(ctxTimeout); err != nil {
		log.Printf("Server shutdown fatal: %v", err)
		return
	}

	log.Print("Server succesfully exited.")
}

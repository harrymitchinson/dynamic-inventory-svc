package api

import (
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	v1 "github.com/harrymitchinson/dynamic-inventory-svc/pkg/api/v1"
	r "github.com/harrymitchinson/dynamic-inventory-svc/pkg/redis"

	"go.uber.org/zap"
)

// Builder is a container for the high level services required for running the API service.
type Builder struct {
	Redis  *redis.Client
	Logger *zap.Logger
	// Engine is optional, if not provided it with be initialised with gin.New()
	Engine *gin.Engine
}

func (builder *Builder) getNamedLogger(name string) *zap.Logger {
	return builder.Logger.With().Named(name)
}

// Setup configures the API and returns the underlying gin.Engine to allow the caller to start the server additional configuration (e.g. TLS, Ports, Sockets)
func Setup(b *Builder) *gin.Engine {
	log := b.getNamedLogger("engine")

	if b.Engine == nil {
		log.Info("initialising engine")
		b.Engine = gin.New()
	} else {
		log.Info("engine already initialised")
	}

	log.Info("adding middlewares")
	b.Engine.Use(ginzap.Ginzap(b.getNamedLogger("request"), time.RFC3339, true))
	b.Engine.Use(ginzap.RecoveryWithZap(b.getNamedLogger("panic"), true))

	log.Info("registering healthcheck routes")
	b.Engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, time.Now().UTC().Format(time.RFC3339))
	})

	b.Engine.GET("/ready", func(c *gin.Context) {
		err := b.Redis.Ping(b.Redis.Context()).Err()
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, "redis ping failed")
			return
		}
		c.JSON(http.StatusOK, time.Now().UTC().Format(time.RFC3339))
	})

	log.Info("registering api/v1 routes")
	apiV1 := b.Engine.Group("api/v1")
	apiV1Logger := b.getNamedLogger("api.v1")

	hostSvc := r.NewHostService(b.Redis, apiV1Logger)

	hostCtrl := v1.NewHostController(hostSvc)
	inventoryCtrl := v1.NewInventoryController(hostSvc)

	apiV1.Use(v1.EnvironmentMiddleware)
	apiV1.POST("hosts/:environment", hostCtrl.CreateHost)
	apiV1.GET("hosts/:environment", hostCtrl.GetHosts)
	apiV1.GET("inventory/:environment", inventoryCtrl.GetInventory)

	return b.Engine
}

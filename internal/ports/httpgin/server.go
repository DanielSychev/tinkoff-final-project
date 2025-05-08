package httpgin

import (
	"context"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
	"homework9/internal/app"
)

func ServiceRecovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic", zap.Error(err.(error)))
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func NewHTTPServer(ctx context.Context, port string, a app.App) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	s := &http.Server{Addr: port, Handler: handler}

	handler.Use(ServiceRecovery(ctx.Value("logger").(*zap.Logger)))
	handler.POST("/api/v1/ads", func(c *gin.Context) {
		CreateAd(c, a)
	})

	handler.PUT("/api/v1/ads/:id/status", func(c *gin.Context) {
		ChangeAdStatus(c, a)
	})

	handler.PUT("/api/v1/ads/:id", func(c *gin.Context) {
		UpdateAd(c, a)
	})

	handler.GET("/api/v1/ads", func(c *gin.Context) {
		ListAds(c, a)
	})

	handler.GET("api/v1/ads/:id", func(c *gin.Context) {
		GetAd(c, a)
	})

	handler.DELETE("/api/v1/ads/:id/del", func(c *gin.Context) {
		DeleteAd(c, a)
	})

	handler.POST("/api/v1/users", func(c *gin.Context) {
		CreateUser(c, a)
	})

	handler.GET("/api/v1/users/:id", func(c *gin.Context) {
		GetUser(c, a)
	})

	handler.DELETE("/api/v1/users/:id/del", func(c *gin.Context) {
		DeleteUser(c, a)
	})
	return s
}

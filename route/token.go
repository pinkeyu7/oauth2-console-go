package route

import (
	apiV1 "oauth2-console-go/api/v1"
	"oauth2-console-go/pkg/request_cache"

	"github.com/gin-gonic/gin"
)

func TokenV1(r *gin.Engine, store request_cache.CacheStore) {
	v1 := r.Group("/v1")

	v1.POST("/token", func(c *gin.Context) {
		apiV1.GetToken(c)
	})
}

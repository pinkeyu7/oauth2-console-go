package route

import (
	apiV1 "oauth2-console-go/api/v1"
	"oauth2-console-go/middleware"
	"oauth2-console-go/pkg/request_cache"
	"time"

	"github.com/gin-gonic/gin"
)

func OauthClientV1(r *gin.Engine, store request_cache.CacheStore) {
	v1Auth := r.Group("/v1/oauth/clients")
	v1Auth.Use(middleware.TokenAuth())

	// Oauth Client 列表
	v1Auth.GET("/", func(c *gin.Context) {
		apiV1.ListOauthClient(c)
	})

	// 取得 Oauth Client
	v1Auth.GET("/:id", func(c *gin.Context) {
		apiV1.GetOauthClient(c)
	})

	// 新增 Oauth Client
	v1Auth.POST("/", request_cache.CachePage(store, time.Second*1, func(c *gin.Context) {
		apiV1.AddOauthClient(c)
	}))

	// 編輯 Oauth Client
	v1Auth.PUT("/:id", request_cache.CachePage(store, time.Second*1, func(c *gin.Context) {
		apiV1.EditOauthClient(c)
	}))
}

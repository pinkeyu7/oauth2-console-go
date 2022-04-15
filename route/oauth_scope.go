package route

import (
	apiV1 "oauth2-console-go/api/v1"
	"oauth2-console-go/middleware"
	"oauth2-console-go/pkg/request_cache"
	"time"

	"github.com/gin-gonic/gin"
)

func OauthScopeV1(r *gin.Engine, store request_cache.CacheStore) {
	v1Auth := r.Group("/v1/oauth/scopes")
	v1Auth.Use(middleware.TokenAuth())

	// Oauth Scope 列表
	v1Auth.GET("/", func(c *gin.Context) {
		apiV1.ListOauthScope(c)
	})

	// 取得 Oauth Scope
	v1Auth.GET("/:id", func(c *gin.Context) {
		apiV1.GetOauthScope(c)
	})

	// 新增 Oauth Scope
	v1Auth.POST("/", request_cache.CachePage(store, time.Second*1, func(c *gin.Context) {
		apiV1.AddOauthScope(c)
	}))

	// 編輯 Oauth Scope
	v1Auth.PUT("/:id", request_cache.CachePage(store, time.Second*1, func(c *gin.Context) {
		apiV1.EditOauthScope(c)
	}))
}

package repository

import (
	"oauth2-console-go/internal/oauth/scope"

	"github.com/go-redis/redis/v7"
)

type Cache struct {
	redis *redis.ClusterClient
}

func NewCache(r *redis.ClusterClient) scope.Cache {
	return &Cache{redis: r}
}

func (c *Cache) DeleteOne(path, method string) error {
	hashKey := scope.GetScopeHashKey()
	key := scope.GetScopeKey(path, method)

	err := c.redis.HDel(hashKey, key).Err()
	return err
}

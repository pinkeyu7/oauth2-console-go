package repository

import (
	"oauth2-console-go/internal/oauth/client"

	"github.com/go-redis/redis/v7"
)

type Cache struct {
	redis *redis.ClusterClient
}

func NewCache(r *redis.ClusterClient) client.Cache {
	return &Cache{redis: r}
}

func (c *Cache) DeleteClientScopeList(clientId string) error {
	key := client.GetClientScopeListKey(clientId)

	err := c.redis.Del(key, key).Err()
	return err
}

func (c *Cache) DeleteAllClientScopeList() error {
	key := client.GetClientScopeListKey("*")
	err := c.redis.ForEachMaster(func(client *redis.Client) error {
		iter := client.Scan(0, key, 0).Iterator()
		for iter.Next() {
			err := client.Del(iter.Val()).Err()
			if err != nil {
				return err
			}
		}
		if err := iter.Err(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

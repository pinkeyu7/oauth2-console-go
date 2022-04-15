package client

import (
	"fmt"
	"oauth2-console-go/dto/apires"
	"oauth2-console-go/dto/model"
)

type Repository interface {
	Count() (int, error)
	Find(limit, offset int) ([]*apires.ListOauthClientItem, error)
	FindOne(client *model.OauthClient) (*model.OauthClient, error)
	Insert(info *model.OauthClient) error
	Update(info *model.OauthClient) error
}

type Cache interface {
	DeleteClientScopeList(clientId string) error
	DeleteAllClientScopeList() error
}

func GetClientScopeListKey(clientId string) string {
	return fmt.Sprintf("client:%s:scope_list", clientId)
}

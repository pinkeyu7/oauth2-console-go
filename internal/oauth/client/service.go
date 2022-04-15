package client

import (
	"oauth2-console-go/dto/apireq"
	"oauth2-console-go/dto/apires"
	"oauth2-console-go/internal/oauth/scope"
)

type Service interface {
	ListClient(sysAccId, page, perPage int) (*apires.ListOauthClient, error)
	GetClient(sysAccId int, clientId string, scopeRepo scope.Repository) (*apires.OauthClient, error)
	AddClient(req *apireq.AddOauthClientWithFile) error
	EditClient(clientId string, req *apireq.EditOauthClientWithFile, scopeRepo scope.Repository) error
}

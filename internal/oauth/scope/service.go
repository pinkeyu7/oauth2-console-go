package scope

import (
	"oauth2-console-go/dto/apireq"
	"oauth2-console-go/dto/apires"
	"oauth2-console-go/dto/model"
)

type Service interface {
	ListScope(sysAccId int, page, perPage int) (*apires.ListOauthScope, error)
	GetScope(sysAccId int, scopeId int) (*model.OauthScope, error)
	AddScope(req *apireq.AddOauthScope) error
	EditScope(scopeId int, req *apireq.EditOauthScope) error
}

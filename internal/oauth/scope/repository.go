package scope

import (
	"fmt"
	"oauth2-console-go/dto/model"
)

type Repository interface {
	Count() (int, error)
	Find(limit, offset int) ([]*model.OauthScope, error)
	FindScope() ([]string, error)
	FindOne(scope *model.OauthScope) (*model.OauthScope, error)
	Insert(scope *model.OauthScope) error
	Update(scope *model.OauthScope) error
}

type Cache interface {
	DeleteOne(path, method string) error
}

func GetScopeHashKey() string {
	return "scope"
}

func GetScopeKey(path, method string) string {
	return fmt.Sprintf("%s:%s", path, method)
}

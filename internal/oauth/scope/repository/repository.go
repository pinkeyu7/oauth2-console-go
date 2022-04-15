package repository

import (
	"oauth2-console-go/dto/model"
	"oauth2-console-go/internal/oauth/scope"

	"xorm.io/xorm"
)

type Repository struct {
	orm *xorm.EngineGroup
}

func NewRepository(orm *xorm.EngineGroup) scope.Repository {
	return &Repository{orm: orm}
}

func (r *Repository) Count() (int, error) {
	cnt, err := r.orm.Count(&model.OauthScope{})
	return int(cnt), err
}

func (r *Repository) Find(limit, offset int) ([]*model.OauthScope, error) {
	var err error
	scopes := make([]*model.OauthScope, 0)

	err = r.orm.Limit(limit, offset).Find(&scopes)
	if err != nil {
		return nil, err
	}

	return scopes, nil
}

func (r *Repository) FindScope() ([]string, error) {
	var err error
	scopes := make([]string, 0)

	err = r.orm.Table("oauth_scope").Where(" is_disable =  ? ", 0).Select("scope").Find(&scopes)
	if err != nil {
		return nil, err
	}

	return scopes, nil
}

func (r *Repository) FindOne(scope *model.OauthScope) (*model.OauthScope, error) {
	has, err := r.orm.Get(scope)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	return scope, nil
}

func (r *Repository) Insert(scope *model.OauthScope) error {
	_, err := r.orm.Insert(scope)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(scope *model.OauthScope) error {
	_, err := r.orm.ID(scope.Id).Update(scope)
	return err
}

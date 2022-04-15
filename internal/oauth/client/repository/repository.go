package repository

import (
	"encoding/json"
	"oauth2-console-go/dto/apires"
	"oauth2-console-go/dto/model"
	"oauth2-console-go/internal/oauth/client"

	"xorm.io/xorm"
)

type Repository struct {
	orm *xorm.EngineGroup
}

func NewRepository(orm *xorm.EngineGroup) client.Repository {
	return &Repository{orm: orm}
}

func (r *Repository) Count() (int, error) {
	cnt, err := r.orm.Count(&model.OauthClient{})
	return int(cnt), err
}

func (r *Repository) Find(limit, offset int) ([]*apires.ListOauthClientItem, error) {
	var err error
	clients := make([]*apires.ListOauthClientItem, 0)

	err = r.orm.Table("oauth_client").Limit(limit, offset).Find(&clients)
	if err != nil {
		return nil, err
	}

	return clients, nil
}

func (r *Repository) FindOne(client *model.OauthClient) (*model.OauthClient, error) {
	has, err := r.orm.Get(client)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	return client, nil
}

func (r *Repository) Insert(info *model.OauthClient) error {
	oc := model.OauthClient{
		Id:           info.Id,
		SysAccountId: info.SysAccountId,
		Name:         info.Name,
		Secret:       info.Secret,
		Domain:       info.Domain,
		Scope:        info.Scope,
		IconPath:     info.IconPath,
	}

	jsonData, err := json.Marshal(oc)
	if err != nil {
		return err
	}
	oc.Data = string(jsonData)

	_, err = r.orm.Insert(&oc)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(info *model.OauthClient) error {
	oc := model.OauthClient{
		Id:           info.Id,
		SysAccountId: info.SysAccountId,
		Name:         info.Name,
		Secret:       info.Secret,
		Domain:       info.Domain,
		Scope:        info.Scope,
		IconPath:     info.IconPath,
	}

	jsonData, err := json.Marshal(oc)
	if err != nil {
		return err
	}
	oc.Data = string(jsonData)

	_, err = r.orm.Where("id = ? ", oc.Id).Cols("sys_account_id", "name", "secret", "domain", "scope", "icon_path", "data").Update(oc)
	if err != nil {
		return err
	}

	return nil
}

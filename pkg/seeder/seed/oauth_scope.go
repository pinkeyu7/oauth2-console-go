package seed

import (
	"oauth2-console-go/dto/model"

	"xorm.io/xorm"
)

func CreateOauthScope(engine *xorm.Engine, scope, path, method, name, description string) error {
	status := false
	oauthScope := model.OauthScope{
		Scope:       scope,
		Path:        path,
		Method:      method,
		Name:        name,
		Description: description,
		IsDisable:   &status,
	}

	_, err := engine.Insert(&oauthScope)

	return err
}

func AllOauthScope() []Seed {
	return []Seed{
		{
			Name: "Create user.profile_get",
			Run: func(engine *xorm.Engine) error {
				err := CreateOauthScope(engine, "user.profile_get", "/v1/users", "GET", "取得使用者資訊", "取得使用者資訊")
				return err
			},
		},
		{
			Name: "Create address-book.list_get",
			Run: func(engine *xorm.Engine) error {
				err := CreateOauthScope(engine, "address-book.list_get", "/v1/contacts", "GET", "取得聯絡人列表", "取得聯絡人列表")
				return err
			},
		},
		{
			Name: "Create address-book.contact_post",
			Run: func(engine *xorm.Engine) error {
				err := CreateOauthScope(engine, "address-book.contact_post", "/v1/contacts", "POST", "新增聯絡人", "新增聯絡人")
				return err
			},
		},
		{
			Name: "Create address-book.contact_get",
			Run: func(engine *xorm.Engine) error {
				err := CreateOauthScope(engine, "address-book.contact_get", "/v1/contacts/:id", "GET", "取得聯絡人資訊", "取得聯絡人資訊")
				return err
			},
		},
	}
}

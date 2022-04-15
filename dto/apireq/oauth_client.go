package apireq

import (
	"mime/multipart"
	"oauth2-console-go/dto/model"
)

type ListOauthClient struct {
	AccountId int `form:"account_id" validate:"required"`
	Page      int `form:"page" validate:"required"`
	PerPage   int `form:"per_page" validate:"required"`
}

type AddOauthClient struct {
	AccountId int    `form:"account_id" validate:"required"`
	Id        string `form:"id" validate:"required"`
	Secret    string `form:"secret" validate:"required"`
	Domain    string `form:"domain" validate:"required"`
	Name      string `form:"name" validate:"required"`
}

type AddOauthClientWithFile struct {
	*AddOauthClient
	File          multipart.File
	FileName      string
	FileExtension string
}

type EditOauthClient struct {
	AccountId int    `form:"account_id" validate:"required"`
	Secret    string `form:"secret" validate:"required"`
	Domain    string `form:"domain" validate:"required"`
	Name      string `form:"name" validate:"required"`
	HasImage  *bool  `form:"has_image" validate:"required"`
	ScopeList string `form:"scope_list" validate:"required"`
}

type EditOauthClientWithFile struct {
	*EditOauthClient
	File          multipart.File
	FileName      string
	FileExtension string
	ScopeList     *model.ScopeList
}

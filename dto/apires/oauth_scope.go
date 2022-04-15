package apires

import (
	"oauth2-console-go/dto/model"
)

type ListOauthScope struct {
	List        []*model.OauthScope `json:"list"`
	CurrentPage int                 `json:"current_page"`
	PerPage     int                 `json:"per_page"`
	Total       int                 `json:"total"`
}

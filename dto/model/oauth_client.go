package model

import "time"

type OauthClient struct {
	Id           string    `xorm:"not null default '' comment('id') VARCHAR(255)" json:"id"`
	SysAccountId int       `xorm:"not null default '' comment('sys_account_id') VARCHAR(255)" json:"sys_account_id"`
	Name         string    `xorm:"not null default '' comment('name') VARCHAR(255)" json:"name"`
	Secret       string    `xorm:"not null default '' comment('secret') VARCHAR(255)" json:"secret"`
	Domain       string    `xorm:"not null default '' comment('domain') VARCHAR(255)" json:"domain"`
	Scope        string    `xorm:"not null default '' comment('scope') VARCHAR(255)" json:"scope"`
	IconPath     string    `xorm:"not null default '' comment('icon_path') VARCHAR(191)" json:"icon_path"`
	Data         string    `xorm:"not null default '' comment('data') TEXT" json:"data"`
	CreatedAt    time.Time `xorm:"not null created DATETIME" json:"created_at"`
	UpdatedAt    time.Time `xorm:"not null updated DATETIME" json:"updated_at"`
}

type OauthClientRedisCache struct {
	ClientID            []string `json:"client_id"`
	CodeChallenge       []string `json:"code_challenge"`
	CodeChallengeMethod []string `json:"code_challenge_method"`
	RedirectURI         []string `json:"redirect_uri"`
	ResponseType        []string `json:"response_type"`
	Scope               []string `json:"scope"`
	State               []string `json:"state"`
}

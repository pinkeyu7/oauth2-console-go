package apires

import (
	"oauth2-console-go/dto/model"
	"time"
)

type ListOauthClient struct {
	List        []*ListOauthClientItem `json:"list"`
	CurrentPage int                    `json:"current_page"`
	PerPage     int                    `json:"per_page"`
	Total       int                    `json:"total"`
}

type ListOauthClientItem struct {
	Id        string    `xorm:"not null pk VARCHAR(255)" json:"id"`
	Secret    string    `xorm:"not null VARCHAR(255)" json:"secret"`
	Domain    string    `xorm:"not null VARCHAR(255)" json:"domain"`
	Name      string    `xorm:"not null VARCHAR(255)" json:"name"`
	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

type OauthClient struct {
	Id           string           `xorm:"not null pk VARCHAR(255)" json:"id"`
	SysAccountId int              `xorm:"not null INT" json:"sys_account_id"`
	Name         string           `xorm:"not null VARCHAR(255)" json:"name"`
	Secret       string           `xorm:"not null VARCHAR(255)" json:"secret"`
	Domain       string           `xorm:"not null VARCHAR(255)" json:"domain"`
	Scope        string           `xorm:"not null VARCHAR(255)" json:"scope"`
	IconPath     string           `xorm:"not null VARCHAR(191)" json:"icon_path"`
	ScopeList    *model.ScopeList `json:"scope_list"`
	CreatedAt    time.Time        `xorm:"created" json:"created_at"`
	UpdatedAt    time.Time        `xorm:"updated" json:"updated_at"`
}

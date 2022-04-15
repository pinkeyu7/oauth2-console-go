package model

import "time"

type OauthScope struct {
	Id          int       `xorm:"not null pk autoincr INT(11)" json:"id"`
	Scope       string    `xorm:"not null VARCHAR(100) scope" json:"scope"`
	Path        string    `xorm:"not null VARCHAR(100) path" json:"path"`
	Method      string    `xorm:"not null VARCHAR(20) method" json:"method"`
	Name        string    `xorm:"not null VARCHAR(30) name" json:"name"`
	Description string    `xorm:"not null VARCHAR(255) description" json:"description" `
	IsDisable   *bool     `xorm:"not null TINYINT is_disable" json:"is_disable"`
	CreatedAt   time.Time `xorm:"not null DATETIME created" json:"created_at"`
	UpdatedAt   time.Time `xorm:"not null DATETIME updated" json:"updated_at"`
}

type ScopeList map[string]*ScopeCategory

type ScopeCategory struct {
	Name   string                `json:"name"`
	Items  map[string]*ScopeItem `json:"items"`
	IsAuth bool                  `json:"is_auth"`
}

type ScopeItem struct {
	Name   string `json:"name"`
	IsAuth bool   `json:"is_auth"`
}

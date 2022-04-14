package model

import "time"

type SysAccount struct {
	Id                       int       `xorm:"pk autoincr BIGINT(20)" json:"id"`
	Account                  string    `xorm:"not null default '' comment('account') VARCHAR(64)" json:"account"`
	Phone                    string    `xorm:"not null default '' comment('phone') VARCHAR(20)" json:"phone"`
	Email                    string    `xorm:"not null default '' comment('email') VARCHAR(64)" json:"email"`
	Password                 string    `xorm:"not null default '' comment('password') VARCHAR(64)" json:"password"`
	Name                     string    `xorm:"not null default '' comment('name') VARCHAR(64)" json:"name"`
	IsDisable                bool      `xorm:"not null is_disable" json:"is_disable"`
	VerifyAt                 time.Time `xorm:"comment('verify_at') DATETIME" json:"verify_at"`
	ForgotPassToken          string    `xorm:"default '' comment('forgot_pass_token') VARCHAR(64)" json:"forgot_pass_token"`
	ForgotPassTokenExpiredAt time.Time `xorm:"comment('forgot_pass_token_expired_at') DATETIME" json:"forgot_pass_token_expired_at"`
	CreatedAt                time.Time `xorm:"not null created DATETIME" json:"created_at"`
	UpdatedAt                time.Time `xorm:"not null updated DATETIME" json:"updated_at"`
}

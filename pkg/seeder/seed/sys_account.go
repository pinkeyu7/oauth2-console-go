package seed

import (
	"github.com/brianvoe/gofakeit/v4"
	"oauth2-console-go/dto/model"
	"xorm.io/xorm"
)

func CreateSysAccount(engine *xorm.Engine, account, name, email, phone string) error {
	defaultPassword := "0eb683eacea7957d8b4140ed837f1ee7fce60ba74e48839a51d6b2085938b49b"

	con := model.SysAccount{
		Account:  account,
		Phone:    phone,
		Email:    email,
		Password: defaultPassword,
		Name:     name,
	}

	_, err := engine.Insert(&con)
	return err
}

func AllSysAccount() []Seed {
	return []Seed{
		{
			Name: "Create System Account - 1",
			Run: func(engine *xorm.Engine) error {
				err := CreateSysAccount(engine, "sys_account", gofakeit.Name(), gofakeit.Email(), gofakeit.Phone())
				if err != nil {
					return err
				}
				return nil
			},
		},
	}
}

package seed

import (
	"encoding/json"
	"oauth2-console-go/dto/model"

	"xorm.io/xorm"
)

func CreateOauthClient(engine *xorm.Engine, id, name, secret, domain, scope string) error {
	con := model.OauthClient{
		Id:           id,
		SysAccountId: 0,
		Name:         name,
		Secret:       secret,
		Domain:       domain,
		Scope:        scope,
		IconPath:     "",
	}

	jsonData, _ := json.Marshal(con)
	con.Data = string(jsonData)

	_, err := engine.Insert(&con)
	return err
}

func AllOauthClient() []Seed {
	return []Seed{
		{
			Name: "Create Oauth Client - address book api",
			Run: func(engine *xorm.Engine) error {
				err := CreateOauthClient(engine, "address-book-go", "Address Book API", "address-book-secret", "http:localhost:9094", "user address-book")
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name: "Create Oauth Client - billing api",
			Run: func(engine *xorm.Engine) error {
				err := CreateOauthClient(engine, "billing-go", "Billing API", "billing-secret", "http:localhost:9094", "user")
				if err != nil {
					return err
				}
				return nil
			},
		},
	}
}

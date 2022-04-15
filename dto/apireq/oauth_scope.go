package apireq

type ListOauthScope struct {
	AccountId int `form:"account_id" validate:"required"`
	Page      int `form:"page" validate:"required"`
	PerPage   int `form:"per_page" validate:"required"`
}

type AddOauthScope struct {
	AccountId   int    `json:"account_id" validate:"required"`
	Scope       string `json:"scope" validate:"required"`
	Path        string `json:"path" validate:"required"`
	Method      string `json:"method" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type EditOauthScope struct {
	AccountId   int    `json:"account_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	IsDisable   *bool  `json:"is_disable" validate:"required"`
}
